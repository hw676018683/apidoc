package router

type fieldCommentPair struct {
	Field   string
	Comment string
}

type ResBody interface {
	SetData(interface{})
}

type ResBodyTpl struct {
	Code    string      `json:"code" c:"ok 表示成功，其他表示错误代码"`
	Message string      `json:"message" c:"与code对应的描述信息"`
	Data    interface{} `json:"data"`
}

func (res *ResBodyTpl) SetData(d interface{}) {
	res.Data = d
}

const (
	TypeReqBody uint8 = iota
	TypeResBody
	TypeErrResBody
)

// TODO
type roundTripBody struct {
	Type uint8 // 请求体/成功返回体/错误返回体
	Desc string
	Body interface{}
}

type routerInfo struct {
	Path   string
	Method string

	Title          string
	Desc           string // 描述
	ReqContentType string

	RegComments   []fieldCommentPair
	QueryComments []fieldCommentPair
	// 保存请求体/成功返回体/错误返回体，数据的数组。并以此顺序生成文档。

	RoundTripBodies []roundTripBody

	//Req           interface{}
	//SucRes        interface{}
	//ErrRes        []ResBodyTpl

	IsEntry bool // 是否 api 接口
}

type R struct {
	Info  routerInfo
	Nodes []*R
}

func New(path string) *R {
	return &R{
		Info: routerInfo{
			Path:            path,
			ReqContentType:  `application/json`,
			RegComments:     make([]fieldCommentPair, 0),
			QueryComments:   make([]fieldCommentPair, 0),
			RoundTripBodies: make([]roundTripBody, 0),
		},
		Nodes: make([]*R, 0),
	}
}

func NewEntry(path string) *R {
	entry := New(path)
	entry.Info.IsEntry = true
	return entry
}

func (r *R) Group(path string) *R {
	child := New(path)
	r.Nodes = append(r.Nodes, child)
	return child
}

func (r *R) Get(path string) *R {
	child := NewEntry(path)
	child.Info.Method = `GET`
	r.Nodes = append(r.Nodes, child)
	return child
}

func (r *R) Post(path string) *R {
	child := NewEntry(path)
	child.Info.Method = `POST`
	r.Nodes = append(r.Nodes, child)
	return child
}

func (r *R) Put(path string) *R {
	child := NewEntry(path)
	child.Info.Method = `PUT`
	r.Nodes = append(r.Nodes, child)
	return child
}

func (r *R) Patch(path string) *R {
	child := NewEntry(path)
	child.Info.Method = `PATCH`
	r.Nodes = append(r.Nodes, child)
	return child
}

func (r *R) Delete(path string) *R {
	child := NewEntry(path)
	child.Info.Method = `DELETE`
	r.Nodes = append(r.Nodes, child)
	return child
}
