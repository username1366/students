package models

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type App struct {
	DB     *gorm.DB
	Router *mux.Router
}

type Student struct {
	//gorm.Model

	ID        uint       `gorm:"primary_key"`
	Name      string     `gorm:"type:varchar(100)" json:"name"`
	Age       int        `json:"age"`
	Rating    int        `json:"rating"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-"`
}

func (a App) InitDB() {
	a.DB.AutoMigrate(&Student{})
}

func (a App) Add(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%v\n", r.Method)
	s := &Student{}
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &s)
	fmt.Printf("%v", s)

	if err != nil {
		fmt.Printf("ERROR %v", err)

	}
	a.DB.Create(&s)
}

func (a *App) Get(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%v\n", r.Method)
	s := &Student{}
	//fmt.Fprintf(w, "%v\n", mux.Vars(r)["id"])
	a.DB.Where("id = ? AND deleted_at IS NULL", mux.Vars(r)["id"]).First(&s)
	if s.Age > 0 && s.Name != "" {
		j, _ := json.Marshal(s)
		fmt.Fprintf(w, "%s\n", j)
	} else {
		fmt.Fprintf(w, "	!\n")
	}
}

func (a *App) Update(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%v\n", r.Method)
	s := &Student{}
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &s)
	if err != nil {
		fmt.Printf("ERROR %v", err)

	}
	//fmt.Fprintf(w, "%v\n", s)
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 32)
	s.ID = uint(id)
	//fmt.Printf("%T", uint(id))
	a.DB.Save(&s)
	//fmt.Fprintf(w, "{\"error\":\"false\"}", ...)
}

func (a *App) Delete(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%v\n", r.Method)
	s := &Student{}
	//fmt.Fprintf(w, "%v\n", mux.Vars(r)["id"])
	a.DB.Where("id = ?", mux.Vars(r)["id"]).First(&s)
	a.DB.Delete(&s)
}

func (a *App) GetAll(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "GETALL\n")
	s := []Student{}
	a.DB.Find(&s)
	fmt.Printf("%T\n", s)
	fmt.Printf("%v\n", s)
	//json.Marshal(v)
	j, _ := json.Marshal(s)
	fmt.Fprintf(w, "%s\n", j)
}
