package server

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"path"
	"time"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/context"
	"github.com/gorilla/csrf"
	"github.com/gorilla/handlers"
	"github.com/julienschmidt/httprouter"
	"github.com/thcyron/tracklog"
	"github.com/thcyron/tracklog/pkg/config"
	"github.com/thcyron/tracklog/pkg/db"
	"github.com/thcyron/tracklog/pkg/models"
)

// DataDir points to the directory where the public/ and templates/ directories are.
var DataDir = "."

type Server struct {
	config      *config.Config
	db          db.DB
	handler     http.Handler
	csrfHandler func(http.Handler) http.Handler
	tmpl        *template.Template
}

func New(conf *config.Config, db db.DB) (*Server, error) {
	s := &Server{
		config: conf,
		db:     db,
	}

	if !s.config.Server.Development {
		tmpl, err := s.loadTemplates()
		if err != nil {
			return nil, err
		}
		s.tmpl = tmpl
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
	r.NotFound = http.FileServer(http.Dir(path.Join(DataDir, "public")))
	r.GET("/signin", s.wrapHandler(s.HandleGetSignIn))
	r.POST("/signin", s.wrapHandler(s.HandlePostSignIn))
	r.POST("/signout", s.wrapHandler(s.HandlePostSignOut))
	r.GET("/logs", s.wrapHandler(s.HandleGetLogs))
	r.POST("/logs", s.wrapHandler(s.HandlePostLog))
	r.GET("/logs/:id/download", s.wrapHandler(s.HandleDownloadLog))
	r.GET("/logs/:id", s.wrapHandler(s.HandleGetLog))
	r.PATCH("/logs/:id", s.wrapHandler(s.HandlePatchLog))
	r.DELETE("/logs/:id", s.wrapHandler(s.HandleDeleteLog))
	r.GET("/", s.wrapHandler(s.HandleDashboard))
	n.UseHandler(r)

	s.handler = n
	return s, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer context.Clear(r)
	s.handler.ServeHTTP(w, r)
}

type HandlerFunc func(w http.ResponseWriter, r *http.Request)

func (s *Server) wrapHandler(handler HandlerFunc) httprouter.Handle {
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

func (s *Server) loadTemplates() (*template.Template, error) {
	return template.ParseGlob(path.Join(DataDir, "templates/*.html"))
}

type renderData struct {
	Title      string
	ActiveTab  string
	Breadcrumb *Breadcrumb
	User       *models.User
	CSRFToken  string
	CSRFField  template.HTML
	Version    string
	Runtime    string
	Content    template.HTML
	Data       interface{}
}

func (s *Server) render(w http.ResponseWriter, r *http.Request, name string) {
	ctx := NewContext(r, w)

	tmpl := s.tmpl
	if tmpl == nil {
		t, err := s.loadTemplates()
		if err != nil {
			panic(err)
		}
		tmpl = t
	}

	data := renderData{
		Title:      ctx.Title(),
		ActiveTab:  ctx.ActiveTab(),
		Breadcrumb: ctx.Breadcrumb(),
		User:       ctx.User(),
		CSRFToken:  csrf.Token(r),
		CSRFField:  csrf.TemplateField(r),
		Version:    tracklog.Version,
		Data:       ctx.Data(),
	}
	if s.config.Server.Development {
		data.Runtime = fmt.Sprintf("%.0fms", time.Now().Sub(ctx.Start()).Seconds()*1000)
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
