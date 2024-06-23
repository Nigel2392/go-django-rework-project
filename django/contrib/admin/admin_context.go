package admin

import (
	"net/http"

	"github.com/Nigel2392/django/core"
	"github.com/Nigel2392/django/core/assert"
	"github.com/Nigel2392/django/core/ctx"
	"github.com/Nigel2392/django/core/tpl"
	"github.com/Nigel2392/django/forms/media"
	"github.com/Nigel2392/goldcrest"
	"github.com/justinas/nosurf"
)

var _ ctx.Context = (*adminContext)(nil)
var _ tpl.RequestContext = (*adminContext)(nil)

type BreadCrumb struct {
	Title string
	URL   string
}

type Action struct {
	Icon   string
	Target string
	Title  string
	URL    string
}

type PageOptions struct {
	Request     *http.Request
	TitleFn     func() string
	SubtitleFn  func() string
	MediaFn     func() media.Media
	BreadCrumbs []BreadCrumb
	Actions     []Action
}

func (p *PageOptions) Title() string {
	if p.TitleFn == nil {
		return ""
	}
	return p.TitleFn()
}

func (p *PageOptions) Subtitle() string {
	if p.SubtitleFn == nil {
		return ""
	}
	return p.SubtitleFn()
}

func (p *PageOptions) GetBreadCrumbs() []BreadCrumb {
	var breadCrumbs = p.BreadCrumbs
	if breadCrumbs == nil {
		breadCrumbs = make([]BreadCrumb, 0)
	}

	var hooks = goldcrest.Get[RegisterBreadCrumbHookFunc](RegisterNavBreadCrumb)
	for _, hook := range hooks {
		var crumbs = hook(p.Request, AdminSite)
		breadCrumbs = append(breadCrumbs, crumbs...)
	}

	return breadCrumbs
}

func (p *PageOptions) GetActions() []Action {
	var actions = p.Actions
	if actions == nil {
		actions = make([]Action, 0)
	}

	var hooks = goldcrest.Get[RegisterNavActionHookFunc](RegisterNavAction)
	for _, hook := range hooks {
		var acts = hook(p.Request, AdminSite)
		actions = append(actions, acts...)
	}

	return actions
}

type adminContext struct {
	Page    *PageOptions
	Site    *AdminApplication
	request *http.Request
	Context ctx.Context
}

func NewContext(request *http.Request, site *AdminApplication, context ctx.Context) *adminContext {
	if context == nil {
		context = core.Context(request)
	}

	assert.False(
		site == nil,
		"Site must be provided to AdminContext",
	)

	var c = &adminContext{
		Context: context,
		Site:    site,
		request: request,
	}

	return c
}

func (c *adminContext) Get(key string) interface{} {
	switch key {
	case "site", "Site":
		return c.Site
	case "page", "Page":
		return c.Page
	}

	return c.Context.Get(key)
}

func (c *adminContext) Set(key string, value interface{}) {
	switch key {
	case "site", "Site":
		c.Site = value.(*AdminApplication)
		return
	case "page", "Page":
		c.Page = value.(*PageOptions)
	}
	c.Context.Set(key, value)
}

func (c *adminContext) SetPage(page PageOptions) {
	c.Page = &page
	c.Page.Request = c.request
}

func (c *adminContext) Request() *http.Request {
	return c.request
}

func (c *adminContext) CsrfToken() string {
	return nosurf.Token(c.request)
}
