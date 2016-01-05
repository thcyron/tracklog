package server

import (
	"net/http"
	"time"

	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
	"github.com/thcyron/tracklog"
)

type contextKey int

const (
	ctxTitle contextKey = iota + 1
	ctxActiveTab
	ctxBreadcrumb
	ctxNoLayout
	ctxStart
	ctxData
	ctxUser
	ctxParams
)

type Context struct {
	r *http.Request
	w http.ResponseWriter
}

func NewContext(r *http.Request, w http.ResponseWriter) *Context {
	return &Context{r: r, w: w}
}

func (c *Context) Breadcrumb() *Breadcrumb {
	br, ok := context.Get(c.r, ctxBreadcrumb).(*Breadcrumb)
	if !ok || br == nil {
		br = new(Breadcrumb)
		context.Set(c.r, ctxBreadcrumb, br)
	}
	return br
}

func (c *Context) Title() string {
	title, _ := context.Get(c.r, ctxTitle).(string)
	return title
}

func (c *Context) SetTitle(title string) {
	context.Set(c.r, ctxTitle, title)
}

func (c *Context) ActiveTab() string {
	activeTab, _ := context.Get(c.r, ctxActiveTab).(string)
	return activeTab
}

func (c *Context) SetActiveTab(activeTab string) {
	context.Set(c.r, ctxActiveTab, activeTab)
}

func (c *Context) NoLayout() bool {
	noLayout, _ := context.Get(c.r, ctxNoLayout).(bool)
	return noLayout
}

func (c *Context) SetNoLayout(noLayout bool) {
	context.Set(c.r, ctxNoLayout, noLayout)
}

func (c *Context) Data() interface{} {
	return context.Get(c.r, ctxData)
}

func (c *Context) SetData(data interface{}) {
	context.Set(c.r, ctxData, data)
}

func (c *Context) User() *tracklog.User {
	user, _ := context.Get(c.r, ctxUser).(*tracklog.User)
	return user
}

func (c *Context) SetUser(user *tracklog.User) {
	context.Set(c.r, ctxUser, user)
}

func (c *Context) Params() httprouter.Params {
	params, _ := context.Get(c.r, ctxParams).(httprouter.Params)
	return params
}

func (c *Context) SetParams(ps httprouter.Params) {
	context.Set(c.r, ctxParams, ps)
}

func (c *Context) Start() time.Time {
	start, _ := context.Get(c.r, ctxStart).(time.Time)
	return start
}

func (c *Context) SetStart(start time.Time) {
	context.Set(c.r, ctxStart, start)
}
