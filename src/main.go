package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"models"
	"net/http"
)

func main() {
	a := &models.App{}

	a.DB, _ = gorm.Open("sqlite3", "../data.db")
	defer a.DB.Close()
	a.InitDB()
	fmt.Printf("%T\n", a.DB)

	a.Router = mux.NewRouter()
	a.Router.HandleFunc("/{id:[0-9]+}", a.Get).Methods("GET")
	a.Router.HandleFunc("/", a.GetAll).Methods("GET")
	a.Router.HandleFunc("/", a.Add).Methods("POST")
	a.Router.HandleFunc("/{id:[0-9]+}", a.Update).Methods("PUT")
	a.Router.HandleFunc("/{id:[0-9]+}", a.Delete).Methods("DELETE")

	http.ListenAndServe(":8000", a.Router)

}
