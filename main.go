package main

import (
	"log"
	"net/http"
)


func middleware(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS"{
			w.WriteHeader(http.StatusOK)
		}
		next.ServeHTTP(w, r)
	})
}

func main(){
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(".")))
	mux.Handle("/assets", http.FileServer(http.Dir("./assets")))

	corsMux := middleware(mux)

	
	log.Fatal(http.ListenAndServe(":8080", corsMux))

}