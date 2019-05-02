package server

import (
	"encoding/hex"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	uuid "github.com/satori/go.uuid"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"short-url/db"
)

type repo interface {
	Save(key, target string)
	Find(key string) (string, bool)
	Track(key string)
	Visits(key string) int64
}

var mainTmpl *template.Template

var MainDB repo = db.NewMemory()

const BaseURL = "http://localhost:5000"

func init() {
	var err error
	mainTmpl, err = template.ParseGlob("tmpl/*.tmpl")
	if err != nil {
		log.Fatal(err)
	}
}

func ListenAndServe() error {
	r := chi.NewMux()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.StripSlashes)
	r.Use(middleware.DefaultCompress)

	r.Get("/new", new)
	r.Get("/stats/{key}", stats)
	r.Get("/{key}", visit)
	r.Post("/", create)

	return http.ListenAndServe(":5000", r)
}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}

func new(w http.ResponseWriter, r *http.Request) {
	panicIf(mainTmpl.ExecuteTemplate(w, "new.tmpl", nil))
}

func create(w http.ResponseWriter, r *http.Request) {
	panicIf(r.ParseForm())
	dest := r.Form.Get("destination")
	u, err := url.Parse(dest)
	if err != nil {
		panicIf(err)
	}
	if u.Scheme == "" {
		u.Scheme = "http"
	}
	if dest == "" {
		panic("missing destination")
	}
	var key string
	for {
		key = hex.EncodeToString(uuid.NewV4().Bytes()[:4])
		_, ok := MainDB.Find(key)
		if ok {
			continue
		}
		break
	}
	MainDB.Save(key, u.String())

	_, err = io.WriteString(w, fmt.Sprintf("Your new URL is %s/%s", BaseURL, key))
	panicIf(err)
}

func visit(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	dest, ok := MainDB.Find(key)
	if !ok {
		panic(fmt.Sprintf("could not find key %q", key))
	}
	MainDB.Track(key)
	http.Redirect(w, r, dest, http.StatusTemporaryRedirect)
}

func stats(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")

	dest, ok := MainDB.Find(key)
	if !ok {
		panic(fmt.Sprintf("key does not exist"))
	}
	_, err := io.WriteString(w, fmt.Sprintf("URL %s\n\n", dest))
	panicIf(err)
	_, err = io.WriteString(w, fmt.Sprintf("Key %q has %d views\n", key, MainDB.Visits(key)))
	panicIf(err)
}
