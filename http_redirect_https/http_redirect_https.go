package main

import "net/http"

func main() {
	if err := http.ListenAndServe(":80", http.RedirectHandler("https://juunini.xyz", http.StatusFound)); err != nil {
		panic(err)
	}
}
