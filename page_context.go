package magic

import (
	"net/http"
)

type PageContext interface {
	Page() *Page
	Request() *http.Request

	Mount(func())
	Unmount(func())
}

type pageContext struct {
	page    *Page
	request *http.Request
	epb     *epBuffer
	// lifetime
	mount   func()
	unmount func()
}

func newPageContext(p *Page, r *http.Request) *pageContext {
	return &pageContext{
		page:    p,
		request: r,
		epb:     &epBuffer{},
	}
}

func (pctx *pageContext) Page() *Page {
	return pctx.page
}

func (pctx *pageContext) Request() *http.Request {
	return pctx.request
}

func (pctx *pageContext) Patch(op Operation, path string, value any) {
	pctx.epb.PushPatch(op, path, value)
}

func (pctx *pageContext) Mount(fn func()) {
	pctx.mount = fn
}

func (pctx *pageContext) Unmount(fn func()) {
	pctx.unmount = fn
}

func (pctx *pageContext) runMount() {
	if pctx.mount != nil {
		pctx.mount()
	}
}

func (pctx *pageContext) runUnmount() {
	if pctx.unmount != nil {
		pctx.unmount()
	}
}
