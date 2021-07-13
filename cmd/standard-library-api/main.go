package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/steveyackey/example-go-api/pkg/data"
)

const basicPaths = `
---BASIC PATHS: localhost:8888 ---
GET /guitars -- all guitars
GET /guitars/<brand> -- a single guitar
GET /finish -- all finish colors
POST /finish '{"color":"Blue"}' -- adds the color blue to finishes
`

func allGuitarsHandler(w http.ResponseWriter, r *http.Request) {
	err := json.NewEncoder(w).Encode(data.Guitars)
	if err != nil {
		log.Printf("error handling guitars... %+v \n", err)
	}
}

func singleGuitarHandler(w http.ResponseWriter, r *http.Request) {
	segments := strings.Split(r.URL.Path, "/")
	requestedGuitar := segments[len(segments)-1]

	value, ok := data.Guitars[requestedGuitar]
	if !ok {
		http.Error(w, "Guitar not found", http.StatusNotFound)
		return
	}

	err := json.NewEncoder(w).Encode(value)
	if err != nil {
		log.Printf("error handling guitars... %+v \n", err)
	}

}

func finishHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		err := json.NewEncoder(w).Encode(data.Finishes)
		if err != nil {
			log.Printf("error handling colors... %+v \n", err)
		}
	case http.MethodPost:
		var c data.Finish
		err := json.NewDecoder(r.Body).Decode(&c)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		data.Finishes = append(data.Finishes, c)

		// json.MarshallIndent returns a byte slice of pretty printed json
		// Could also just use json.Marshall here
		newColors, err := json.MarshalIndent(data.Finishes, "", "    ")
		if err != nil {
			http.Error(w, "Unable to unmarshall data", http.StatusInternalServerError)
		}

		fmt.Fprintln(w, string(newColors))
		log.Printf("Added %+v", c)
		r.Body.Close()
	default:
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
	}

}

func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s - %s %s %s - Request Started", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		next.ServeHTTP(w, r)
		log.Printf("%s - %s %s %s - Request Finished", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
	})
}

// recoverPanic handles any panics that may happen and gracefully recovers, logging the error
func recoverPanic(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				log.Printf("%s\n%s", err, debug.Stack())
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func main() {

	log.Println(basicPaths)

	log.Println("Starting web server...")

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, basicPaths)
	})

	mux.HandleFunc("/guitars", recoverPanic(loggingMiddleware(allGuitarsHandler)))
	mux.HandleFunc("/guitars/", recoverPanic(loggingMiddleware(singleGuitarHandler)))
	mux.HandleFunc("/finish", recoverPanic(loggingMiddleware(finishHandler)))

	log.Fatal(http.ListenAndServe(":8888", mux))
}
