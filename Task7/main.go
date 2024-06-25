package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prachin77/task7/handlers"
)

func main() {
	r := mux.NewRouter()

	fmt.Println("listening on port 8080")

	r.HandleFunc("/",handlers.DefaultRoute).Methods("GET")
	r.HandleFunc("/app",handlers.GetApp).Methods("GET")

	r.HandleFunc("/getipadd",handlers.GetIpAdd).Methods("GET")
	r.HandleFunc("/getpingform",handlers.GetPingForm).Methods("GET")
	r.HandleFunc("/ping",handlers.PingToUser).Methods("POST")
	r.HandleFunc("/getnslookupform",handlers.GetNsLookUpForm).Methods("GET")
	r.HandleFunc("/nslookup",handlers.NsLookUp).Methods("POST")
	// r.HandleFunc("/getwhoisform").Methods("GET")

	log.Fatal(http.ListenAndServe(":8080",r))
}
