package main

import (
	"io"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {

	dl := r.FormValue("dl")
	if r.URL.Path != "/" || dl == "" {
		BadRequest(w)
		return
	}

	log.Println(dl)

	resp, err := http.Get(dl)
	if err != nil {
		log.Fatalf("http.Get => %v", err.Error())
	}

	io.Copy(w, resp.Body)

	if err != nil {
		log.Fatalf("io.Copy => %v", err.Error())
	}
}

func BadRequest(c http.ResponseWriter) {
	http.Error(c, "Sad web server is sad :(", http.StatusBadRequest)
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8081", nil)
}
