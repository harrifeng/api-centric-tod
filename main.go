package main

import (
	"log"
	"net/http"
	"time"

	"github.com/ant0ine/go-json-rest/rest"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

func main() {

	i := Impl{}
	i.InitDB()
	i.InitSchema()

	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		rest.Get("/todos", i.GetAllTodos),
		rest.Post("/todos", i.PostTodo),
		rest.Get("/todos/:id", i.GetTodo),
		rest.Put("/todos/:id", i.PutTodo),
		rest.Delete("/todos/:id", i.DeleteTodo),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)

	http.Handle("/api/", http.StripPrefix("/api", api.MakeHandler()))
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("static"))))
	log.Fatal(http.ListenAndServe(":8087", nil))
}

func CommonFileServer(w rest.ResponseWriter, r *rest.Request) {
	http.FileServer(http.Dir("static")).ServeHTTP(w.(http.ResponseWriter), r.Request)
}

type Todo struct {
	Id        int64     `json:"Id"`
	Title     string    `sql:"size:1024" json:"Title"`
	Done      bool      `json:"Done"`
	CreatedAt time.Time `json:"CreatedAt"`
	UpdatedAt time.Time `json:"UpdatedAt"`
	DeletedAt time.Time `json:"-"`
}

type Impl struct {
	DB gorm.DB
}

func (i *Impl) InitDB() {
	var err error
	i.DB, err = gorm.Open("mysql", "root:wyyzmwy@tcp(10.97.31.116:23306)/resttest?charset=utf8&parseTime=True")
	if err != nil {
		log.Fatalf("Got error when connect database, the error is '%v'", err)
	}
	i.DB.LogMode(true)
}

func (i *Impl) InitSchema() {
	i.DB.AutoMigrate(&Todo{})
}

func (i *Impl) GetAllTodos(w rest.ResponseWriter, r *rest.Request) {
	todos := []Todo{}
	i.DB.Find(&todos)
	w.WriteJson(&todos)
}

func (i *Impl) GetTodo(w rest.ResponseWriter, r *rest.Request) {
	id := r.PathParam("id")
	todo := Todo{}
	if i.DB.First(&todo, id).Error != nil {
		rest.NotFound(w, r)
		return
	}
	w.WriteJson(&todo)
}

func (i *Impl) PostTodo(w rest.ResponseWriter, r *rest.Request) {
	todo := Todo{}
	if err := r.DecodeJsonPayload(&todo); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := i.DB.Save(&todo).Error; err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson(&todo)
}

func (i *Impl) PutTodo(w rest.ResponseWriter, r *rest.Request) {

	id := r.PathParam("id")
	todo := Todo{}
	if i.DB.First(&todo, id).Error != nil {
		rest.NotFound(w, r)
		return
	}

	updated := Todo{}
	if err := r.DecodeJsonPayload(&updated); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	todo.Title = updated.Title
	todo.Done = updated.Done

	if err := i.DB.Save(&todo).Error; err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson(&todo)
}

func (i *Impl) DeleteTodo(w rest.ResponseWriter, r *rest.Request) {
	id := r.PathParam("id")
	todo := Todo{}
	if i.DB.First(&todo, id).Error != nil {
		rest.NotFound(w, r)
		return
	}
	if err := i.DB.Delete(&todo).Error; err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
