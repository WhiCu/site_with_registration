package server

import (
	"log"
	"net/http"
	"os"
	"reg/pkg/mysql"
)

type server struct {
	mux *http.ServeMux
	app application
}

const (
	code     = "ui/html/code.layout.tmpl"
	regIndex = "ui/html/index.html"
	assets   = "ui/assets"
)

func New() *server {

	app := application{
		infoLog:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		errorLog: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
	}

	userDB := "user:password@/bd?parseTime=true" // ________________________________
	db, err := mysql.OpenDB(userDB)
	if err != nil {
		app.errorLog.Fatal(err)
	}
	app.infoLog.Println("db is connected")
	app.db = mysql.New(db)

	mux := http.NewServeMux()

	fileServer := http.FileServer(modFileSystem{
		fs:  http.Dir("ui/assets"),
		app: app,
	})
	mux.Handle("/assets/", http.StripPrefix("/assets/", fileServer))
	//mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir(assets))))

	mux.HandleFunc("/", app.NotFound())

	mux.HandleFunc("/assets", app.NotAllowed())

	mux.HandleFunc("/home", app.home())

	mux.HandleFunc("/home/account", app.Ok())

	mux.HandleFunc("/db", app.registration())

	return &server{
		mux: mux,
		app: app,
	}
}

func (sv *server) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	sv.app.infoLog.Println(pattern)
	sv.mux.HandleFunc(pattern, handler)
}

func (sv *server) Go() {
	sv.app.infoLog.Println("Server is listening...")
	err := http.ListenAndServe(":8080", sv.mux)
	sv.app.errorLog.Fatal(err)
}
