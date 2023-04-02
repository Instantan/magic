package magic

import (
	"net/http"
)

type PageContext interface {
	Page() *Page
	Request() *http.Request
	OnCleanup(func())
}

type pageContext struct {
	page         *Page
	request      *http.Request
	epb          *epBuffer
	cleanupFuncs []func()
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

func (pctx *pageContext) OnCleanup(fn func()) {
	pctx.cleanupFuncs = append(pctx.cleanupFuncs, fn)
}

func (pctx *pageContext) Cleanup() {
	for i := range pctx.cleanupFuncs {
		pctx.cleanupFuncs[i]()
	}
}
