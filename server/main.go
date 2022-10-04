package main

import (
	"html/template"
	"image"
	"log"
	"net/http"
	"time"

	"image/jpeg"

	"github.com/o1egl/govatar"
)

func main() {
	port := "8080"
	router := http.NewServeMux()
	router.HandleFunc("/api/v1/avatar", postAvatar)
	router.HandleFunc("/", index)

	server := &http.Server{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Addr:         ":" + port,
		Handler:      router,
	}
	log.Println("serving govatar @ http://localhost:" + port)
	log.Fatal(server.ListenAndServe())
}

func postAvatar(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.NotFound(w, r)
	} else {
		r.ParseForm()
		name := ""
		name = r.FormValue("username")
		gender := govatar.MALE
		if r.FormValue("gender") == "female" {
			gender = govatar.FEMALE
		}
		var img image.Image
		var err error
		if name == "" {
			img, err = govatar.Generate(gender)
		} else {
			img, err = govatar.GenerateForUsername(gender, name)
		}
		if err != nil {
			panic(err)
		}
		w.Header().Set("Content-Type", "image/jpeg")
		if err := jpeg.Encode(w, img, nil); err != nil {
			log.Println("unable to write image.")
		}
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("templates/index.html")
		t.Execute(w, nil)
	} else {
		http.NotFound(w, r)
	}
}
