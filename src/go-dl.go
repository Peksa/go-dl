package main

import (
	"bufio"
	"bytes"
	"code.google.com/p/go.crypto/bcrypt"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"crypto/tls"
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
	file, err := os.Open("users.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		str := scanner.Text()
		if strings.HasPrefix(str, "#") {
			continue
		}
		tokens := strings.Split(str, ":")
		if tokens[0] == username {
			return bcrypt.CompareHashAndPassword([]byte(tokens[1]), []byte(password)) == nil
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return false
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

func getTlsServer() *http.Server {
	return &http.Server{
		Addr: ":10443",
		TLSConfig: &tls.Config{
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_RSA_WITH_RC4_128_SHA,
			},
			MinVersion: tls.VersionTLS10,
		},
	}
}

func main() {
	srv := getTlsServer()
	http.HandleFunc("/dl", handler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		NotFound(w)
	})
	log.Fatal(srv.ListenAndServeTLS("cert.pem", "key.pem"))
}
