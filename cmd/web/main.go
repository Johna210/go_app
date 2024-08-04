package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type application struct {
	errorLog *log.Logger
	infoLog *log.Logger
}

func main() {
	// Define a new command-line flag with the name 'addr', a default value of ":4000"
	// and some short help text explaining what the flag controls. The value of the
	// flag will be stpred in the addr variable at runtime.
	addr := flag.String("addr",":4000","Http network address")

	// This reads in the command-line flag value and assigns it to the addr variable
	flag.Parse()

	// A logger for information
	infoLog := log.New(os.Stdout,"INFO\t",log.Ldate|log.Ltime )

	// A logger for an error
	errorLog := log.New(os.Stderr,"Error\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{
		errorLog: errorLog,
		infoLog: infoLog,
	}


	// router / servermux
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)

	// Create a new File Serve
	// Given relative path
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	// Use the mux.Handle Function to register the file server as the handler
	mux.Handle("/static/",http.StripPrefix("/static",fileServer))

	srv := &http.Server{
		Addr: *addr,
		ErrorLog: errorLog,
		Handler: mux,
	}

	// Server
	infoLog.Printf("Starting server on %s", *addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}
