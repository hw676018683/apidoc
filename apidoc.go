package apidoc

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/lovego/apidoc/defaults"
	"github.com/lovego/apidoc/json_doc"
	"github.com/lovego/apidoc/router"
)

var GenerateRes = generateRes

func GenDocs(r *router.R, workDir string) {
	if err := os.RemoveAll(workDir); err != nil {
		panic(err)
	}
	//merge(r)
	genDocs(r, ``, workDir)
}

// genDocs Generate doc.
// basePath is base router path.
// dirPath is dictionary path.
func genDocs(r *router.R, basePath, workDir string) {
	basePath = basePath + r.Info.Path
	if err := os.MkdirAll(workDir, 0755); err != nil {
		log.Panic(err)
	}
	indexes := make([]string, 0)
	for i := range r.Nodes {
		child := r.Nodes[i]
		if child.Info.IsEntry {
			if child.Info.Title == `` {
				log.Println(`Warning: Title is required. API: ` + r.Info.Method + ` ` + basePath + child.Info.Path)
				continue
			}
			docStr := parseEntryDoc(child, basePath)
			buf := []byte(docStr)
			fileName := child.Info.Title + `.md`
			fullPath := filepath.Join(workDir, fileName)
			if _, err := os.Stat(fullPath); err == nil {
				panic(`Error: ` + fileName + ` is exist, are you using a existing title ?`)
			}
			if err := ioutil.WriteFile(fullPath, buf, 0666); err != nil {
				log.Panic(err)
			}
			indexes = append(indexes, `### [`+child.Info.Title+`](`+fileName+`)`)
		}

		// If child router is not an entry and title is not empty,
		// then create a sub directory.
		childDir := workDir
		if !child.Info.IsEntry && child.Info.Title != `` {
			childDir = filepath.Join(workDir, child.Info.Title)
			indexes = append(indexes, `### [`+child.Info.Title+`](`+child.Info.Title+`)`)
		}
		genDocs(child, basePath, childDir)
	}
	if len(indexes) > 0 {
		indexesBuf := []byte(strings.Join(indexes, "\n"))
		if err := ioutil.WriteFile(
			filepath.Join(workDir, `README.md`), indexesBuf, 0666,
		); err != nil {
			log.Panic(err)
		}
	}
}

// merge same path group
func merge(r *router.R) {
	if r == nil || r.Info.IsEntry {
		return
	}
	path2Node := make(map[string]*router.R)
	for i := range r.Nodes {
		n := r.Nodes[i]
		if n.Info.IsEntry {
			continue
		}
		if path2Node[n.Info.Path] == nil {
			path2Node[n.Info.Path] = n
		} else {
			path2Node[n.Info.Path].Nodes = append(path2Node[n.Info.Path].Nodes, n.Nodes...)
			if n.Info.Title != `` {
				path2Node[n.Info.Path].Info.Title = n.Info.Title
			}
			r.Nodes[i] = nil
		}
		merge(n)
	}

	nodes := make([]*router.R, 0)
	for i := range r.Nodes {
		if r.Nodes[i] != nil {
			nodes = append(nodes, r.Nodes[i])
		}
	}
	r.Nodes = nodes
}

func parseEntryDoc(r *router.R, basePath string) (content string) {
	urlPath := basePath + r.Info.Path
	docs := make([]string, 0)
	// title
	title := `# ` + r.Info.Title
	docs = append(docs, title)

	// description
	if r.Info.Desc != `` {
		docs = append(docs, r.Info.Desc)
	}

	// URL
	reqUrl := `## ` + r.Info.Method + ` ` + urlPath
	docs = append(docs, reqUrl)

	// RegComments
	if len(r.Info.RegComments) > 0 {
		docs = append(docs, "\n"+`## Path 参数说明`)
		for _, o := range r.Info.RegComments {
			docs = append(docs, `- `+o.Field+`: `+o.Comment)
		}
	}
	// QueryComments
	if len(r.Info.QueryComments) > 0 {
		docs = append(docs, "\n"+`## Query 参数说明`)
		for _, o := range r.Info.QueryComments {
			docs = append(docs, `- `+o.Field+`: `+o.Comment)
		}
	}

	var hasResBody = false
	for i := range r.Info.RoundTripBodies {
		o := &r.Info.RoundTripBodies[i]
		switch o.Type {
		case router.TypeReqBody:
			docs = append(docs, "\n"+`## 请求体说明(`+r.Info.ReqContentType+`)`)
			if o.Desc != `` {
				docs = append(docs, "\n"+o.Desc)
			}
			docs = append(docs, "```json5")
			docs = append(docs, parseJsonDoc(defaults.Set(o.Body)))
			docs = append(docs, "```")
		case router.TypeResBody:
			hasResBody = true
			docs = append(docs, "\n"+`## 返回体说明`)
			if o.Desc != `` {
				docs = append(docs, "\n"+o.Desc)
			}

			res := GenerateRes(defaults.Set(o.Body))
			docs = append(docs, "```json5")
			docs = append(docs, parseJsonDoc(res))
			docs = append(docs, "```")
		case router.TypeErrResBody:
			docs = append(docs, "\n"+`## 返回错误说明`)
			if o.Desc != `` {
				docs = append(docs, "\n"+o.Desc)
			}
			docs = append(docs, "```json5")
			docs = append(docs, parseJsonDoc(defaults.Set(o.Body)))
			docs = append(docs, "```")
		}
	}

	if !hasResBody {
		res := GenerateRes(nil)
		docs = append(docs, "\n"+`## 返回体说明`)
		docs = append(docs, "```json5")
		docs = append(docs, parseJsonDoc(res))
		docs = append(docs, "```")
	}

	content = strings.Join(docs, "\n")
	return
}

func parseJsonDoc(v interface{}) string {
	const commentLineOffset = 50
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
			comment := strings.TrimSpace(res[0][1])
			if comment[0] == '*' {
				comment = `【必须】` + comment[1:]
			}
			repeatTimes := commentLineOffset - len(str)
			if repeatTimes < 1 {
				repeatTimes = 1
			}
			str += strings.Repeat(` `, repeatTimes) + `// ` + comment
			list[i] = str
		}
	}
	return strings.Join(list, "\n")
}

type resBodyTpl struct {
	Code    string      `json:"code" c:"ok 表示成功，其他表示错误代码"`
	Message string      `json:"message" c:"与code对应的描述信息"`
	Data    interface{} `json:"data"`
}

func generateRes(data interface{}) interface{} {
	return resBodyTpl{Code: "ok", Message: "success", Data: data}
}
