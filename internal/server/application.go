package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"reg/pkg/mysql"
	"reg/pkg/types"
	"strconv"
	"strings"
	"text/template"
	"time"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	db       *mysql.Users
}

type status struct {
	Code int16
	Text string
}

var (
	NotAllowed = status{
		Code: 405,
		Text: "Method Not Allowed",
	}
	NotFound = status{
		Code: 404,
		Text: "Not Found",
	}
	InternalServerError = status{
		Code: 500,
		Text: "Internal Server Error",
	}
	Ok = status{
		Code: 200,
		Text: "OK",
	}
)

func (app application) home() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Println("IPv6", r.RemoteAddr)
		http.ServeFile(w, r, regIndex)

	}
}

func (app application) Ok() func(http.ResponseWriter, *http.Request) {
	return app.handleStatus(Ok, nil)
}
func (app application) NotAllowed() func(http.ResponseWriter, *http.Request) {
	return app.handleStatus(NotAllowed, nil)
}

func (app application) NotFound() func(http.ResponseWriter, *http.Request) {
	return app.handleStatus(NotFound, nil)
}
func (app application) ServerError() func(http.ResponseWriter, *http.Request) {
	return app.handleStatus(InternalServerError, nil)
}

func (app application) Code(stat status, w http.ResponseWriter) {

	w.WriteHeader(int(stat.Code))

	tmpl, err := template.ParseFiles(code)
	if err != nil {
		app.errorLog.Println(err)
		http.Error(w, "Internal Server Error", 500)
		return
	}

	err = tmpl.Execute(w, stat)
	if err != nil {
		app.errorLog.Println(err)
		http.Error(w, "Internal Server Error", 500)
		return
	}
}

func (app application) handleStatus(stat status, methods []string) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Println("IPv6", r.RemoteAddr)
		app.infoLog.Println(stat.Code, stat.Text, r.URL.String())

		if methods != nil {
			w.Header().Set("Allow", strings.Join(methods, ", "))
		}

		app.Code(stat, w)
	}
}

func (app application) registration() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Println(r.URL)
		switch r.Method {
		case http.MethodPost:
			ok, id := app.regUser(r)
			if ok {
				w.WriteHeader(int(Ok.Code))
				jsres, _ := json.Marshal(struct {
					Id int `json:"id"`
				}{id})
				w.Write(jsres)
				return
			}
			w.WriteHeader(int(NotFound.Code))
		case http.MethodGet:
			app.showUser(w, r)
		default:
			app.Code(NotAllowed, w)
		}
	}

}

func (app application) regUser(r *http.Request) (bool, int) {

	dec := json.NewDecoder(r.Body)

	var user types.User

	err := dec.Decode(&user)

	if err != nil {
		app.errorLog.Println(err)
	}
	if user.Login == "" || user.Password == "" {
		app.infoLog.Println("Плохая отправка")
		return false, 0
	}
	app.infoLog.Println("Хорошая отправка")
	user.Created = time.Now()

	app.infoLog.Println(user)
	id, _ := app.db.Insert(&user)
	return true, id
}

func (app application) showUser(w http.ResponseWriter, r *http.Request) {
	app.infoLog.Println(r.URL.String())
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.Code(NotFound, w) // Страница не найдена.
		return
	}

	u, err := app.db.Get(id)

	if err != nil {
		if errors.Is(err, types.ErrNoRecord) {
			app.Code(NotFound, w)
		} else {
			app.errorLog.Println(err)
			app.Code(InternalServerError, w)
		}
		return
	}
	fmt.Fprintf(w, "%v", u)
}
