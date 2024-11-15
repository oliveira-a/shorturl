package main

import (
	"crypto/sha256"
	"encoding/base64"
	"io"
	"log"
	"net/http"
)

const Port = ":8080"
const AppBaseURL = "http://localhost"

// hash code -> origin url
var shortenedUrls map[string]string

func main() {
	shortenedUrls = make(map[string]string)

	http.HandleFunc("/url/short", handleShortenUrl)
	http.HandleFunc("/{hash}", handleUrlRedirect)

	log.Print("'/url/short' waiting for requests.\n")
	log.Print("'/{hash}' waiting for requests.\n")

	log.Fatal(http.ListenAndServe(Port, nil))
}

func handleShortenUrl(w http.ResponseWriter, req *http.Request) {
	url := req.FormValue("url")
	hash := generateHash(url)

	shortenedUrls[hash] = url
	log.Printf("Url '%s' has been stored with hash code '%s'", url, hash)

	io.WriteString(w, AppBaseURL+Port+"/"+hash)
}

func handleUrlRedirect(w http.ResponseWriter, req *http.Request) {
	// we're getting rid of the "/" here. e.g. /F2scAd -> F2scAd
	hash := req.RequestURI[1:]

	originalUrl := shortenedUrls[hash]

	log.Printf("Found original url '%s' from hash code '%s'. Redirecting...", originalUrl, hash)

	http.Redirect(w, req, originalUrl, http.StatusFound)
}

func generateHash(url string) string {
	hash := sha256.Sum256([]byte(url))

	encoded := base64.URLEncoding.EncodeToString(hash[:])

	return encoded[:6]
}
