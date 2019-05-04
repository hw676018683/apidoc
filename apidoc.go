package apidoc

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/lovego/apidoc/defaults"
	"github.com/lovego/apidoc/json_doc"
	"github.com/lovego/apidoc/router"
)

type ResBodyTpl struct {
	Code    string      `json:"code" c:"ok 表示成功，其他表示错误代码"`
	Message string      `json:"message" c:"与code对应的描述信息"`
	Data    interface{} `json:"data"`
}

var BaseRes = ResBodyTpl{Code: "ok", Message: "success"}

type Doc struct {
	Indexes  []string
	Contents []string
}

func NewDoc(r *router.R) *Doc {
	d := &Doc{
		Indexes:  make([]string, 0),
		Contents: make([]string, 0),
	}
	d.Indexes = append(d.Indexes, `# Index <a name="index"></a>`)

	d.Parse(r, ``, 1)
	return d
}

func (d *Doc) Create(dir, name string) {
	docs := append(append(d.Indexes, ``), d.Contents...)
	buf := []byte(strings.Join(docs, "\n"))
	if err := ioutil.WriteFile(
		filepath.Join(dir, name+".md"), buf, 0666,
	); err != nil {
		log.Panic(err)
	}
}

func (d *Doc) Parse(r *router.R, path string, level int) {
	path += r.Path
	if r.IsGroup && r.Path != `` && level < 3 {
		idx := strings.Repeat("  ", level-1) + `- `
		idx += `[` + r.Title + ` ` + r.Path + `](#` + path + `)`
		d.Indexes = append(d.Indexes, idx)
		content := strings.Repeat("#", level)
		content += r.Title + ` ` + r.Path
		d.Contents = append(d.Contents, content)
		level += 1
	}
	if len(r.Node) == 0 {
		idx, c := parseRouterDoc(r, path, level)
		d.Indexes = append(d.Indexes, idx)
		d.Contents = append(d.Contents, c)
	}
	for i := range r.Node {
		d.Parse(r.Node[i], path, level)
	}
}

var anchorNameReg = regexp.MustCompile(`[^a-zA-Z0-9]+`)

func parseRouterDoc(r *router.R, path string, level int) (idx, content string) {
	docs := make([]string, 0)
	name := r.Method + anchorNameReg.ReplaceAllStringFunc(path, func(s string) string {
		res := `-`
		return res
	})
	title := `#### ` + r.Method + ` ` + path
	if r.Title != `` {
		title += ` (` + r.Title + `)`
	}
	title += `<a name="` + name + `"></a>`
	title += ` [index](#index)`
	docs = append(docs, title)

	if len(r.RegComments) > 0 {
		docs = append(docs, `##### 正则参数说明`)
		for _, o := range r.RegComments {
			docs = append(docs, `- `+o.Field+`: `+o.Comment)
		}
	}

	if len(r.QueryComments) > 0 {
		docs = append(docs, `##### Query 参数说明`)
		for _, o := range r.QueryComments {
			docs = append(docs, `- `+o.Field+`: `+o.Comment)
		}
	}
	if r.ReqBody != nil {
		docs = append(docs, `##### Request Body`)
		docs = append(docs, "```json5")
		docs = append(docs, parseJsonDoc(defaults.Set(r.ResBody)))
		docs = append(docs, "```")
	}

	if r.ResBody != nil {

		res := BaseRes
		res.Data = defaults.Set(r.ResBody)
		docs = append(docs, `##### Response Body`)
		docs = append(docs, "```json5")
		docs = append(docs, parseJsonDoc(&res))
		docs = append(docs, "```")
	}
	docs = append(docs)
	idx = strings.Repeat(`  `, level-1) + `- `
	idx += `[` + r.Title + ` ` + r.Method + ` ` + r.Path + `](#` + name + `)`
	content = strings.Join(docs, "\n")
	return
}

func parseJsonDoc(v interface{}) string {
	data, err := json_doc.MarshalIndent(v, ``, `  `)
	if err != nil {
		log.Panic(err)
	}
	list := strings.Split(string(data), "\n")

	r := regexp.MustCompile(`@@@([\s\S]*)":`)
	for i := range list {
		res := r.FindAllStringSubmatch(list[i], -1)
		if len(res) > 0 {
			str := r.ReplaceAllString(list[i], `":`)
			str += ` // ` + res[0][1]
			list[i] = str
		}
	}
	return strings.Join(list, "\n")
}