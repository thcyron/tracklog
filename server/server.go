package server

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/context"
	"github.com/gorilla/csrf"
	"github.com/gorilla/handlers"
	"github.com/julienschmidt/httprouter"
	"github.com/thcyron/tracklog"
	"github.com/thcyron/tracklog/db"
)

type Server struct {
	config      *tracklog.Config
	db          db.DB
	handler     http.Handler
	csrfHandler func(http.Handler) http.Handler
}

func New(config *tracklog.Config, db db.DB) *Server {
	s := &Server{
		config: config,
		db:     db,
	}

	n := negroni.Classic()

	csrfHandler := csrf.Protect(
		[]byte(s.config.Server.CSRFAuthKey),
		csrf.Secure(!s.config.Server.Development),
		csrf.FieldName("_csrf"),
	)
	n.UseFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		csrfHandler(next).ServeHTTP(w, r)
	})

	n.UseFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		handlers.HTTPMethodOverrideHandler(next).ServeHTTP(w, r)
	})

	n.UseFunc(s.userAuthMiddleware)

	r := httprouter.New()
	r.NotFound = http.FileServer(http.Dir("public"))
	r.GET("/signin", s.wrapHandler(s.HandleGetSignIn))
	r.POST("/signin", s.wrapHandler(s.HandlePostSignIn))
	r.POST("/signout", s.wrapHandler(s.HandlePostSignOut))
	r.GET("/logs", s.wrapHandler(s.HandleGetLogs))
	r.POST("/logs", s.wrapHandler(s.HandlePostLog))
	r.GET("/logs/:id/download", s.wrapHandler(s.HandleDownloadLog))
	r.GET("/logs/:id", s.wrapHandler(s.HandleGetLog))
	r.DELETE("/logs/:id", s.wrapHandler(s.HandleDeleteLog))
	r.GET("/", s.wrapHandler(s.HandleDashboard))
	n.UseHandler(r)

	s.handler = n
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer context.Clear(r)
	s.handler.ServeHTTP(w, r)
}

type HandlerFunc func(w http.ResponseWriter, r *http.Request)

func (server *Server) wrapHandler(handler HandlerFunc) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		ctx := NewContext(r, w)
		ctx.SetStart(time.Now())
		ctx.SetParams(ps)
		handler(w, r)
	}
}

func (s *Server) userAuthMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	user, err := s.userFromRequest(r)
	if err != nil {
		panic(err)
	}
	if user != nil {
		ctx := NewContext(r, w)
		ctx.SetUser(user)
	}
	next(w, r)
}

type renderData struct {
	Title      string
	ActiveTab  string
	Breadcrumb *Breadcrumb
	User       *tracklog.User
	CSRFToken  string
	CSRFField  template.HTML
	Version    string
	Runtime    string
	Content    template.HTML
	Data       interface{}
}

func (server *Server) render(w http.ResponseWriter, r *http.Request, name string) {
	ctx := NewContext(r, w)

	// TODO(thcyron): Load templates initially at server startup, not on every request.
	tmpl, err := template.ParseGlob("templates/*.html")
	if err != nil {
		panic(err)
	}

	data := renderData{
		Title:      ctx.Title(),
		ActiveTab:  ctx.ActiveTab(),
		Breadcrumb: ctx.Breadcrumb(),
		User:       ctx.User(),
		CSRFToken:  csrf.Token(r),
		CSRFField:  csrf.TemplateField(r),
		Version:    tracklog.Version,
		Runtime:    fmt.Sprintf("%.0fms", time.Now().Sub(ctx.Start()).Seconds()*1000),
		Data:       ctx.Data(),
	}

	if ctx.NoLayout() {
		if err := tmpl.ExecuteTemplate(w, name+".html", data); err != nil {
			panic(err)
		}
		return
	}

	buf := new(bytes.Buffer)
	if err := tmpl.ExecuteTemplate(buf, name+".html", data); err != nil {
		panic(err)
	}
	data.Content = template.HTML(buf.String())

	if err := tmpl.ExecuteTemplate(w, "layout.html", data); err != nil {
		panic(err)
	}
}
