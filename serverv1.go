package main

import (
	"encoding/json"
	"net/http"
	"log"
	"time"
)
func trace(s string) (string, time.Time) {
    log.Println("START:", s)
    return s, time.Now()
}
func un(s string, startTime time.Time) {
    endTime := time.Now()
    log.Println("  END:", s, "ElapsedTime in seconds:", endTime.Sub(startTime))
}
func handler(w http.ResponseWriter, r *http.Request) {
	defer un(trace("query_time"))
	qp := r.URL.Query().Get("q")
	if qp != "" {
		config := getConfig()
		results := asyncHttpGets(qp, config)
		json.NewEncoder(w).Encode(results)
		json.MarshalIndent(results, "", "    ")
	}
}
func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}