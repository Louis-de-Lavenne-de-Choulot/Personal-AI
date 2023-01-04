package main

import (
	"net/http"
)

func Send(w http.ResponseWriter, r *http.Request) {
	// find optional parameters text, lang and fromLang
	input := r.FormValue("text")
	lang := r.FormValue("lang")
	fromLang := r.FormValue("fromlang")
	w.Write([]byte(Serv(input, lang, fromLang)))
}
