package main

import (
	"blog_view/router"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {

	http.HandleFunc("/", router.Index)
	http.HandleFunc("/tags", router.Tags)
	http.HandleFunc("/tags/", router.TagsList)
	http.HandleFunc("/log/", router.Log)
	http.HandleFunc("/img/", img)
	http.HandleFunc("/css/", css)
	http.HandleFunc("/js/", js)
	http.HandleFunc("/webfonts/", webfonts)

	log.Fatalln(http.ListenAndServeTLS(
		":443",
		"/etc/letsencrypt/live/juunini.xyz/fullchain.pem",
		"/etc/letsencrypt/live/juunini.xyz/privkey.pem",
		nil,
	))
}

// "/img/..."
func img(res http.ResponseWriter, req *http.Request) {
	if !strings.EqualFold(req.Method, "GET") {
		return
	}

	FileReader, err := ioutil.ReadFile("/home/juunini/blog_view/views" + req.URL.Path)
	if err != nil {
		stderr(err)
		return
	}
	switch req.URL.Path[strings.LastIndex(req.URL.Path, "."):] {
	case ".png":
		res.Header().Set("Content-Type", "image/png")
	case ".jpg":
		res.Header().Set("Content-Type", "image/jpeg")
	case ".gif":
		res.Header().Set("Content-Type", "image/gif")
	}
	if _, err := res.Write(FileReader); err != nil {
		stderr(err)
		return
	}
	return
}

// "/js/..."
func js(res http.ResponseWriter, req *http.Request) {
	if !strings.EqualFold(req.Method, "GET") {
		return
	}

	FileReader, err := ioutil.ReadFile("/home/juunini/blog_view/views" + req.URL.Path)
	if err != nil {
		stderr(err)
		return
	}

	res.Header().Set("Content-Type", "text/javascript")
	if _, err := res.Write(FileReader); err != nil {
		stderr(err)
		return
	}
	return
}

// "/css/..."
func css(res http.ResponseWriter, req *http.Request) {
	if !strings.EqualFold(req.Method, "GET") {
		return
	}

	FileReader, err := ioutil.ReadFile("/home/juunini/blog_view/views" + req.URL.Path)
	if err != nil {
		stderr(err)
		return
	}

	res.Header().Set("Content-Type", "text/css")
	if _, err := res.Write(FileReader); err != nil {
		stderr(err)
		return
	}
	return
}

// "/webfonts/..."
func webfonts(res http.ResponseWriter, req *http.Request) {
	if !strings.EqualFold(req.Method, "GET") {
		return
	}

	FileReader, err := ioutil.ReadFile("/home/juunini/blog_view/views" + req.URL.Path)
	if err != nil {
		stderr(err)
		return
	}

	switch req.URL.Path[strings.LastIndex(req.URL.Path, "."):] {
	case ".woff":
		res.Header().Set("Content-Type", "application/font-woff")
	case ".woff2":
		res.Header().Set("Content-Type", "font/woff2")
	case ".svg":
		res.Header().Set("Content-Type", "image/svg+xml")
	case ".ttf":
		res.Header().Set("Content-Type", "application/font-sfnt")
	case ".otf":
		res.Header().Set("Content-Type", "application/font-sfnt")
	case ".eot":
		res.Header().Set("Content-Type", "application/vnd.ms-fontobject")
	}
	if _, err := res.Write(FileReader); err != nil {
		stderr(err)
		return
	}
	return
}

func stderr(err error) {
	if _, err := fmt.Fprintln(os.Stderr, time.Now().String()[:19], err); err != nil {
		return
	}
}
