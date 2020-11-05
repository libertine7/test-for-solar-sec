package main

// it simple as possible but no simpler

import (
	"encoding/json"
	"fmt"
	"github.com/go-pg/pg"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"os"
)

var db *pg.DB

type Vacancy struct {
	tableName   struct{} `sql:"public.vacancies"`
	Id          int      `sql:"id"`
	Forename    string   `sql:"forename"`
	Salarylevel int      `sql:"salarylevel"`
	Experience  string   `sql:"experience"`
	City        string   `sql:"city"`
}

func CreateTable() {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS public.vacancies(id serial PRIMARY KEY,
							forename text,
							salarylevel int,
							experience text,
							city text
						);`)
	if err != nil {
		log.Fatal(err)
	}
}

type handler func(w http.ResponseWriter, r *http.Request)

func basicAuth(hnd handler) handler {
	return func(w http.ResponseWriter, r *http.Request) {
		user, _, ok := r.BasicAuth()
		if !ok {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized.", 401)
			return
		}

		if (user != "viewer") && (user != "editor") {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized.", 401)
			return
		}

		hnd(w, r)
	}
}

func basicAuthForEditor(hnd handler) handler {
	return func(w http.ResponseWriter, r *http.Request) {
		user, _, ok := r.BasicAuth()

		if !ok {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized.", 401)
			return
		}
		if user != "editor" {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized.", 401)
			return
		}

		hnd(w, r)
	}
}

func ifError(w http.ResponseWriter, err error) bool {
	if err != nil {
		w.WriteHeader(422)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(fmt.Sprintf("{error: \"%s\" }", err)))
		return true
	}
	return false
}

func GetVacancyHandler(w http.ResponseWriter, r *http.Request) {
	var vacancies []Vacancy
	err := db.Model(&vacancies).Select()

	if ifError(w, err) {
		return
	}
	buf, err := json.Marshal(vacancies)
	if ifError(w, err) {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(buf)
}

func GetVacancyIdHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if ifError(w, err) {
		return
	}

	vacancy := &Vacancy{Id: id}
	err = db.Select(vacancy)
	if ifError(w, err) {
		return
	}
	buf, err := json.Marshal(vacancy)
	if ifError(w, err) {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(buf)
}

func PutVacancyIdHandler(w http.ResponseWriter, r *http.Request) {
	var vacancy Vacancy
	err := json.NewDecoder(r.Body).Decode(&vacancy)
	if ifError(w, err) {
		return
	}

	err = db.Insert(&vacancy)
	db.Model(&vacancy).Returning("*").Insert()
	if ifError(w, err) {
		return
	}

	buf, err := json.Marshal(vacancy)
	if ifError(w, err) {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(buf)
}

func DeleteVacancyIdHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if ifError(w, err) {
		return
	}

	vacancy := &Vacancy{Id: id}
	err = db.Delete(vacancy)
	if ifError(w, err) {
		return
	}
	buf, err := json.Marshal(vacancy)
	if ifError(w, err) {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(buf)
}

func main() {

	connStr := "postgres://postgres:mysecretpassword@172.17.0.2/postgres?sslmode=disable"
	if len(os.Args) < 2 {
		fmt.Println("need DB connection string as param, exemple:")
		fmt.Println("go run main.go postgres://postgres:mysecretpassword@172.17.0.2/postgres?sslmode=disable")
		fmt.Println("--")
		fmt.Println("to run docker if you need use:")
		fmt.Println("docker run --name some-postgres -e POSTGRES_PASSWORD=mysecretpassword -d postgres")
		fmt.Println("")
		os.Exit(0)
	} else {
		connStr = os.Args[1]
	}

	opt, err := pg.ParseURL(connStr)
	if err != nil {
		log.Fatal(err)
	}

	db = pg.Connect(opt)
	CreateTable()

	router := mux.NewRouter()
	router.HandleFunc("/vacancy", basicAuth(GetVacancyHandler)).Methods("GET")
	router.HandleFunc("/vacancy/{id}", basicAuth(GetVacancyIdHandler)).Methods("GET")
	router.HandleFunc("/vacancy", basicAuthForEditor(PutVacancyIdHandler)).Methods("PUT")
	router.HandleFunc("/vacancy/{id}", basicAuthForEditor(DeleteVacancyIdHandler)).Methods("DELETE")
	fmt.Println("HTTP Start at http://0.0.0.0:8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}
