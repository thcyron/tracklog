package server

import (
	"errors"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/thcyron/tracklog/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

const tokenCookieName = "tracklog_token"

var tokenSigningMethod = jwt.SigningMethodHS256

func (s *Server) HandleGetSignIn(w http.ResponseWriter, r *http.Request) {
	ctx := NewContext(r, w)
	user := ctx.User()
	if user != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	s.renderSignIn(w, r, signInData{})
}

func (s *Server) HandlePostSignIn(w http.ResponseWriter, r *http.Request) {
	username, password := r.FormValue("username"), r.FormValue("password")
	if username == "" || password == "" {
		s.renderSignIn(w, r, signInData{})
		return
	}

	user, err := s.db.UserByUsername(username)
	if err != nil {
		panic(err)
	}
	if user == nil {
		s.renderSignIn(w, r, signInData{
			Username: username,
			Alert:    "Bad username/password",
		})
		return
	}
	switch bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) {
	case nil:
		break
	case bcrypt.ErrMismatchedHashAndPassword:
		s.renderSignIn(w, r, signInData{
			Username: username,
			Alert:    "Bad username/password",
		})
		return
	default:
		panic(err)
	}

	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims["user_id"] = user.ID
	token.Claims["v"] = user.PasswordVersion

	tokenString, err := token.SignedString([]byte(s.config.Server.SigningKey))
	if err != nil {
		panic(err)
	}

	cookie := &http.Cookie{
		Name:     tokenCookieName,
		Value:    tokenString,
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)

	http.Redirect(w, r, "/", http.StatusFound)
}

type signInData struct {
	Username string
	Alert    string
}

func (s *Server) renderSignIn(w http.ResponseWriter, r *http.Request, data signInData) {
	ctx := NewContext(r, w)
	ctx.SetNoLayout(true)
	ctx.SetData(data)
	s.render(w, r, "signin")
}

func (s *Server) HandlePostSignOut(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:     tokenCookieName,
		Value:    "",
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
	s.redirectToSignIn(w, r)
}

func (s *Server) redirectToSignIn(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/signin", http.StatusFound)
}

func (s *Server) userFromRequest(r *http.Request) (*models.User, error) {
	cookie, err := r.Cookie(tokenCookieName)
	if err != nil {
		return nil, nil
	}

	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("bad signing method")
		}
		return []byte(s.config.Server.SigningKey), nil
	})
	if err != nil || !token.Valid {
		return nil, nil
	}

	id, ok := token.Claims["user_id"].(float64)
	if !ok {
		return nil, nil
	}
	user, err := s.db.UserByID(int(id))
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, nil
	}

	v, ok := token.Claims["v"].(float64)
	if !ok {
		return nil, err
	}
	if int(v) != user.PasswordVersion {
		return nil, nil
	}
	return user, nil
}
