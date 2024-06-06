package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prachin77/task2/Task2/handlers"
)

func main() {
	r := mux.NewRouter()

	fmt.Println("listening on port 8080 ....")
	r.HandleFunc("/",handlers.DefaultRoute).Methods("GET")
	r.HandleFunc("/app",handlers.GetApp).Methods("GET")

	r.HandleFunc("/getregister",handlers.GetRegister).Methods("GET")
	r.HandleFunc("/postregister",handlers.Register).Methods("POST")
	r.HandleFunc("/getlogin",handlers.GetLogin).Methods("GET")
	r.HandleFunc("/postlogin",handlers.Login).Methods("POST")
	r.HandleFunc("/verifyotp",handlers.VerifyOtp).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", r))
}
