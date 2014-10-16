package main

import (
	"io"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func handler(w http.ResponseWriter, r *http.Request) {
	dl, err := url.Parse(r.FormValue("url"));
	if err != nil || !dl.IsAbs() {
		BadRequest(w)
		return
	}
	log.Println("GET: " + dl.String())
	resp, err := http.Get(dl.String())
	if err != nil {
		log.Println("ERROR: http.Get => %v", err.Error())
		InternalServerError(w, "Failure fetching URL. Wrong protocol? DNS correct? Connection refused?")
		return
	}
	tokens := strings.Split(dl.Path, "/")
	fileName := tokens[len(tokens)-1];
	w.Header().Add("Content-Type", "application/octet-stream");
	w.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
	_, err = io.Copy(w, resp.Body)

	if err != nil {
		log.Println("ERROR: io.Copy  => %v", err.Error())
		return
	}
}

func BadRequest(w http.ResponseWriter) {
	http.Error(w, "Sad web server is sad :(", http.StatusBadRequest)
}

func NotFound(w http.ResponseWriter) {
	http.Error(w, "404 nothing here", http.StatusNotFound)
}

func InternalServerError(w http.ResponseWriter, message string) {
	http.Error(w, message, http.StatusInternalServerError)
}

func main() {
	http.HandleFunc("/dl", handler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		NotFound(w);
	})
	log.Fatal(http.ListenAndServeTLS(":10443", "cert.pem", "key.pem", nil))
}
