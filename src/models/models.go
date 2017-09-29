package models

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

type (
	App struct {
		DB     *gorm.DB
		Router *mux.Router
		Socket string
	}

	Student struct {
		ID        uint       `gorm:"primary_key" json:"id"`
		Name      string     `gorm:"type:varchar(100)" json:"name"`
		Age       int        `json:"age"`
		Rating    int        `json:"rating"`
		CreatedAt time.Time  `json:"-"`
		UpdatedAt time.Time  `json:"-"`
		DeletedAt *time.Time `json:"-"`
	}

	Message struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}
)

func (a App) InitDB() { // create DB schema
	a.DB.AutoMigrate(&Student{})
}

func (a *App) Run() {
	a.DB, _ = gorm.Open("sqlite3", "../data.db")
	defer a.DB.Close()
	a.InitDB()
	a.Router = mux.NewRouter()
	a.Router.HandleFunc("/{id:[0-9]+}", a.Get).Methods("GET")
	a.Router.HandleFunc("/", a.GetAll).Methods("GET")
	a.Router.HandleFunc("/", a.Add).Methods("POST")
	a.Router.HandleFunc("/{id:[0-9]+}", a.Update).Methods("PUT")
	a.Router.HandleFunc("/{id:[0-9]+}", a.Delete).Methods("DELETE")
	log.Printf("Started on socket %s", a.Socket)
	http.ListenAndServe(a.Socket, a.Router)
}

func (s Student) Validate() bool {
	if s.Name != "" && s.Age > 0 && s.Rating > 0 {
		return true
	}
	return false
}

func (a App) Add(w http.ResponseWriter, r *http.Request) {
	s := &Student{}
	w.Header().Set("Content-Type", "application/json")
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &s)
	if err != nil {
		j, _ := json.Marshal(Message{Error: true, Message: err.Error()})
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "%s\n", j)
		log.Printf("%v %v", r.RequestURI, err.Error())
	} else {
		a.DB.Create(&s)
		w.WriteHeader(http.StatusCreated)
		j, _ := json.Marshal(Message{Error: false, Message: "Created"})
		fmt.Fprintf(w, "%s\n", j)
	}
}

func (a *App) Get(w http.ResponseWriter, r *http.Request) {
	s := &Student{}
	w.Header().Set("Content-Type", "application/json")
	a.DB.Where("id = ? AND deleted_at IS NULL", mux.Vars(r)["id"]).First(&s)
	if s.Validate() {
		j, _ := json.Marshal(s)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "%s\n", j)
	} else {
		j, _ := json.Marshal(Message{Error: true, Message: "No content"})
		w.WriteHeader(http.StatusNoContent)
		fmt.Fprintf(w, "%s\n", j)
		log.Printf("%v No content", r.RequestURI)
	}
}

func (a *App) Update(w http.ResponseWriter, r *http.Request) {
	s := &Student{}
	w.Header().Set("Content-Type", "application/json")
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &s)
	if err != nil {
		j, _ := json.Marshal(Message{Error: true, Message: err.Error()})
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "%s\n", j)
		log.Printf("%v %v", r.RequestURI, err.Error())
		return
	} else if s.Validate() {
		id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 32)
		s.ID = uint(id)
		a.DB.Save(&s)
		w.WriteHeader(http.StatusOK)
		j, _ := json.Marshal(Message{Error: false, Message: "Updated"})
		fmt.Fprintf(w, "%s\n", j)
	} else {
		j, _ := json.Marshal(Message{Error: true, Message: "JSON validation fails"})
		fmt.Fprintf(w, "%s\n", j)
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("%v %v", r.RequestURI, err.Error())
	}
}

func (a *App) Delete(w http.ResponseWriter, r *http.Request) {
	s := &Student{}
	w.Header().Set("Content-Type", "application/json")
	a.DB.Where("id = ? AND deleted_at IS NULL", mux.Vars(r)["id"]).First(&s)
	if s.Name != "" && s.Age > 0 && s.Rating > 0 {
		a.DB.Delete(&s)
		w.WriteHeader(http.StatusOK)
		j, _ := json.Marshal(Message{Error: false, Message: "Deleted"})
		fmt.Fprintf(w, "%s\n", j)
	} else {
		j, _ := json.Marshal(Message{Error: true, Message: "Record is already deleted"})
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "%s\n", j)
		log.Printf("%v Record is already deleted", r.RequestURI)
	}
}

func (a *App) GetAll(w http.ResponseWriter, r *http.Request) {
	s := []Student{}
	w.Header().Set("Content-Type", "application/json")
	a.DB.Where("deleted_at IS NULL").Find(&s)
	if len(s) > 0 {
		j, _ := json.Marshal(s)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "%s\n", j)
	} else {
		j, _ := json.Marshal(Message{Error: true, Message: "Could not find any record"})
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "%s\n", j)
		log.Printf("%v Could not find any record", r.RequestURI)
	}
}
