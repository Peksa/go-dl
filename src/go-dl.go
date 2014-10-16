package main

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func handler(w http.ResponseWriter, r *http.Request) {
	if !IsAuthorized(r) {
		NotAuthorized(w)
		return
	}

	dl, err := url.Parse(r.FormValue("url"))
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
	fileName := tokens[len(tokens)-1]
	w.Header().Add("Content-Type", "application/octet-stream")
	w.Header().Add("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, fileName))
	_, err = io.Copy(w, resp.Body)

	if err != nil {
		log.Println("ERROR: io.Copy  => %v", err.Error())
		return
	}
}

func IsAuthorized(r *http.Request) bool {
	username, password, err := ExtractCredentials(r.Header.Get("Authorization"))
	if err != nil {
		return false
	}
	return ValidateCredentials(username, password)
}

func ExtractCredentials(header string) (string, string, error) {
	if header == "" {
		return "", "", errors.New("Empty authorization string")
	}
	parts := strings.Split(header, " ")
	if len(parts) != 2 {
		return "", "", errors.New("Malformed authorization header")
	}
	creds, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", "", errors.New("Bad credential encoding")
	}
	index := bytes.Index(creds, []byte(":"))
	if parts[0] != "Basic" || index < 0 {
		return "", "", errors.New("Authorization format not HTTP Basic")
	}
	return string(creds[:index]), string(creds[index+1:]), nil
}

func ValidateCredentials(username string, password string) bool {
	return "demo" == username && "demo" == password
}

func BadRequest(w http.ResponseWriter) {
	http.Error(w, "Sad web server is sad :(", http.StatusBadRequest)
}

func NotFound(w http.ResponseWriter) {
	http.Error(w, "404 nothing here", http.StatusNotFound)
}

func NotAuthorized(w http.ResponseWriter) {
	w.Header().Set("Www-Authenticate", `Basic realm="Credentials please!"`)
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte("401 Beep, wrong username or password"))
}

func InternalServerError(w http.ResponseWriter, message string) {
	http.Error(w, message, http.StatusInternalServerError)
}

func main() {
	http.HandleFunc("/dl", handler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		NotFound(w)
	})
	log.Fatal(http.ListenAndServeTLS(":10443", "cert.pem", "key.pem", nil))
}
