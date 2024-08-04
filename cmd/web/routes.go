package main

import "net/http"

func (app *application) routes() *http.ServeMux{
	// router / servermux
	mux := http.NewServeMux()
	mux.HandleFunc("/",app.home)
	mux.HandleFunc("/snippet",app.showSnippet)
	mux.HandleFunc("/snippet/create",app.createSnippet)

	// Create a new File Serve
	// Given relative path
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	// Use the mux.Handle Function to register the file server as the handler
	mux.Handle("/static/",http.StripPrefix("/static",fileServer))

	return mux
}