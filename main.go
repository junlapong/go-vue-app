package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/gobuffalo/packr/v2"
	"github.com/gorilla/mux"
)

// Message : struct for message
type Message struct {
	Text string `json:"message"`
}

func main() {
	// Bind folder path for packaging with Packr
	staticBox := packr.New("static", "./public/static")

	router := mux.NewRouter()
	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// return index.html
		html, _ := staticBox.FindString("index.html")
		w.Write([]byte(html))
	})

	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(staticBox)))
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		html, _ := staticBox.FindString("index.html")
		w.Write([]byte(html))
	})

	var api = router.PathPrefix("/api").Subrouter()
	api.Use(loggingRequest)

	api.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{ "message" : "API not found" }`))
	})

	api.HandleFunc("/ping", ping).Methods("GET")
	api.HandleFunc("/now", now).Methods("POST")

	// Run server at port 8000
	log.Println("http://localhost:8080")
	myIP()
	log.Fatal(http.ListenAndServe(":8080", router))
}

func ping(w http.ResponseWriter, r *http.Request) {
	// Create Message JSON data
	message := Message{Text: "pong"}

	// Return JSON encoding to output
	output, err := json.Marshal(message)

	// Catch error, if it happens
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set header Content-Type
	w.Header().Set("Content-Type", "application/json")

	// Write output
	w.Write(output)
}

func now(w http.ResponseWriter, r *http.Request) {
	// Read body
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(os.Stdout, "%v\n", string(b))

	// Set header Content-Type
	w.Header().Set("Content-Type", "application/json")

	// Write output
	w.Write([]byte(`{"timestamp":"` + time.Now().Format(time.RFC3339) + `"}`))
}

type statusRecorder struct {
	http.ResponseWriter
	Status int
}

func (r *statusRecorder) writeHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}

func loggingRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		recorder := &statusRecorder{
			ResponseWriter: w,
			Status:         200,
		}
		next.ServeHTTP(recorder, r)
		log.Printf("%7s %d %s", r.Method, recorder.Status, r.URL.Path)
	})
}

func myIP() {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		os.Stderr.WriteString("Oops: " + err.Error() + "\n")
		os.Exit(1)
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				os.Stdout.WriteString("http://" + ipnet.IP.String() + ":8080\n")
			}
		}
	}
}
