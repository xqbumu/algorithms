package main

import "net/http"

func main() {
	http.ListenAndServe(":6100", http.HandlerFunc(handler))
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello"))
}
