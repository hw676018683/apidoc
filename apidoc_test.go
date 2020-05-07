package apidoc

import (
	"path/filepath"
	"runtime"

	"github.com/lovego/apidoc/router"
)

type req struct {
	Name string `json:"name" c:"名称"`
	Req  string `json:"req" c:"请求参数"`
}

type res struct {
	Name string `json:"name" doc:"名称"`
	Age  string `json:"age" doc:"年龄"`
	Res  string `json:"res" doc:"返回信息"`
}

type errorRes struct {
	Name   string `json:"name" c:"*名称"`
	Age    string `json:"age" c:"年龄" binding:"required"`
	ErrRes string `json:"errRes" c:"错误返回信息"`
}

type bodyTpl struct {
	Code string      `json:"code"`
	Msg  string      `json:"msg,omitempty"`
	Data interface{} `json:"data"`
}

func (b *bodyTpl) SetData(d interface{}) {
	b.Data = d
}

func ExampleGenDocs() {
	BaseRes = &bodyTpl{Code: "ok"}
	rootRouter := router.New(``)
	setup(rootRouter)
	GenDocs(rootRouter, filepath.Join(sourceDir(), `apidocs`))

	// Output:
}

func setup(r *router.R) {
	purchaseRouter := r.Group(`/purchases`).Title(`采购`)
	arlRouter := purchaseRouter.Group(`/arrival`).Title(`到货单`)

	arlRouter.Post(`/book`).Title(`订餐`).
		Desc(`描述信息描述信息描述信息描述信息描述信息描述信息描述信息`).
		Regex(`id:公司ID`).
		Query(`qid:公司QID;qid2:公司`).
		Req(`描述信息描述信息描述信息描述信息描述信息描述信息描述信息req1`, &req{}).
		Req(`描述信息描述信息描述信息描述信息描述信息描述信息描述信息req2`, nil).
		Res(`描述信息描述信息描述信息描述信息描述信息描述信息描述信息res1`, &res{}).
		Res(`描述信息描述信息描述信息描述信息描述信息描述信息描述信息res2`, nil).
		ErrRes(`描述信息描述信息描述信息描述信息描述信息描述信息描述信息err1`, bodyTpl{`something-wrong`, `some thing wrong`, &errorRes{}}).
		ErrRes(`描述信息描述信息描述信息描述信息描述信息描述信息描述信息err2`, bodyTpl{`something-wrong2`, `some thing wrong2`, nil})

	saleRouter := r.Group(`/sales`).Title(`销售`)
	saleOrderRouter := saleRouter.Group(`/order`).Title(`订单`)
	saleOrderRouter.Get(`/detail/(\d+)`).
		Doc(`获取订单详情`, `ID:订单ID`, `name:用户名`, nil, nil)
	saleOrderRouter.Put(`/detail/(\d+)`).
		Doc(`更新订单详情`, `ID:订单ID`, `name:用户名`, &req{}, &res{})
}

func sourceDir() string {
	_, filename, _, _ := runtime.Caller(1)
	return filepath.Dir(filename)
}
