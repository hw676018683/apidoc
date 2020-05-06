package router

import (
	"log"
	"reflect"
	"strings"
)

// Title set router Title.
// Don't contains '/'
// Don't start with + or - or .
func (r *R) Title(t string) *R {
	t = strings.TrimSpace(t)
	if t != `` {
		if strings.ContainsAny(t, `/`) {
			panic(`Title contains '/' : ` + t)
		}
		if t[0] == '+' || t[0] == '-' || t[0] == '.' {
			panic(`Title starts with '+-.' : ` + t)
		}
	}
	if r.Info.Title != `` {
		log.Println(`Warning: Title is reassigned. old: ` + r.Info.Title + ` new: ` + t)
	}
	r.Info.Title = t
	return r
}

// Desc set router descriptions.
func (r *R) Desc(d string) *R {
	r.Info.Desc = d
	return r
}

// ContentType set request content type.
func (r *R) ContentType(s string) *R {
	r.Info.ReqContentType = strings.TrimSpace(s)
	return r
}

// Regex set request regex parameters.
func (r *R) Regex(d string) *R {
	r.Info.RegComments = parseFieldCommentPair(d)
	return r
}

// Query set request query parameters.
func (r *R) Query(d string) *R {
	r.Info.QueryComments = parseFieldCommentPair(d)
	return r
}

// Req set request body.
func (r *R) Req(desc string, d interface{}) *R {
	if d == nil {
		return r
	}
	if reflect.TypeOf(d).Kind() != reflect.Ptr {
		panic(`Req need pointer`)
	}
	roundTripInfo := roundTripBody{
		Type: TypeReqBody,
		Desc: desc,
		Body: d,
	}
	r.Info.RoundTripBodies = append(r.Info.RoundTripBodies, roundTripInfo)
	return r
}

// Res set success response body.
func (r *R) Res(desc string, d interface{}) *R {
	if d != nil && reflect.TypeOf(d).Kind() != reflect.Ptr {
		panic(`Res need pointer`)
	}

	roundTripInfo := roundTripBody{
		Type: TypeResBody,
		Desc: desc,
		Body: d,
	}
	r.Info.RoundTripBodies = append(r.Info.RoundTripBodies, roundTripInfo)
	return r
}

// ErrRes add error response bodies.
func (r *R) ErrRes(desc, code string, msg string, data interface{}) *R {
	if data == nil {
		return r
	}
	if reflect.TypeOf(data).Kind() != reflect.Ptr {
		panic(`ErrRes need pointer`)
	}
	obj := ResBodyTpl{
		Code:    code,
		Message: msg,
		Data:    data,
	}
	roundTripInfo := roundTripBody{
		Type: TypeErrResBody,
		Desc: desc,
		Body: obj,
	}
	r.Info.RoundTripBodies = append(r.Info.RoundTripBodies, roundTripInfo)
	return r
}

// Doc provide quick set common api docs.
func (r *R) Doc(t string, reg, query string, req, res interface{}) *R {
	r.Title(t)
	r.Regex(reg)
	r.Query(query)
	if req != nil {
		r.Req(``, req)
	}
	r.Res(``, res)
	return r
}

func parseFieldCommentPair(src string) (list []fieldCommentPair) {
	list = make([]fieldCommentPair, 0)
	if src == `` {
		return
	}
	pairs := strings.Split(src, ";")
	for i := range pairs {
		item := strings.TrimSpace(pairs[i])
		if item == `` {
			continue
		}
		parts := strings.Split(item, ":")
		if len(parts) > 0 {
			p := fieldCommentPair{Field: strings.TrimSpace(parts[0])}
			if len(parts) > 1 {
				p.Comment = strings.TrimSpace(parts[1])
			}
			list = append(list, p)
		}
	}
	return
}
