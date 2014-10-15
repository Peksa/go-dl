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

	log.Println("GET: " + dl)

	resp, err := http.Get(dl)
	if err != nil {
		log.Println("ERROR: http.Get => %v", err.Error())
		InternalServerError(w, "Failure fetching URL. Wrong protocol? DNS correct? Connection refused?")
		return
	}

	_, err = io.Copy(w, resp.Body)

	if err != nil {
		log.Println("ERROR: io.Copy  => %v", err.Error())
		return
	}
}

func BadRequest(w http.ResponseWriter) {
	http.Error(w, "Sad web server is sad :(", http.StatusBadRequest)
}

func InternalServerError(w http.ResponseWriter, message string) {
	http.Error(w, message, http.StatusInternalServerError)
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8081", nil)
}
