package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/steveyackey/example-go-api/pkg/data"
)

const basicPaths = `
---BASIC PATHS: localhost:8888 ---
GET /guitars -- all guitars
GET /guitars/<brand> -- a single guitar
GET /finish -- all finish colors
POST /finish '{"color":"Blue"}' -- adds the color blue to finishes
`

func allGuitarsHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	err := json.NewEncoder(w).Encode(data.Guitars)
	if err != nil {
		log.Printf("error handling guitars... %+v \n", err)
	}
}

func singleGuitarHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if _, ok := data.Guitars[p.ByName("brand")]; !ok {
		w.WriteHeader(404)
		w.Write([]byte("404 - Guitar not found"))
		log.Printf("%s - %s %s %s - 404 - Guitar not found", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
	}
	err := json.NewEncoder(w).Encode(data.Guitars[p.ByName("brand")])
	if err != nil {
		log.Printf("error handling guitars... %+v \n", err)
	}
}

func addFinishColor(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var c data.FinishColor
	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	data.FinishColors = append(data.FinishColors, c)

	// json.MarshallIndent returns a byte slice of pretty printed json
	// Could also just use json.Marshall here
	newColors, err := json.MarshalIndent(data.FinishColors, "", "    ")
	if err != nil {
		http.Error(w, "Unable to unmarshall data", http.StatusInternalServerError)
	}

	fmt.Fprintln(w, string(newColors))
	log.Printf("Added %+v", c)
}

func allFinishColors(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	err := json.NewEncoder(w).Encode(data.FinishColors)
	if err != nil {
		log.Printf("error handling colors... %+v \n", err)
	}
}

func loggingMiddleware(next httprouter.Handle) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		log.Printf("%s - %s %s %s - Request Started", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		next(w, r, p)
		log.Printf("%s - %s %s %s - Request Finished", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
	})
}

func main() {

	log.Println(basicPaths)

	log.Println("Starting web server...")

	router := httprouter.New()
	router.GET("/", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		fmt.Fprintln(w, basicPaths)
	})

	router.GET("/guitars", loggingMiddleware(allGuitarsHandler))
	router.GET("/guitars/:brand", loggingMiddleware(singleGuitarHandler))
	router.GET("/finish", loggingMiddleware(allFinishColors))
	router.POST("/finish", loggingMiddleware(addFinishColor))

	log.Fatal(http.ListenAndServe(":8888", router))
}
