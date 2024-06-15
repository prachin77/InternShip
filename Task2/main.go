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
	r.HandleFunc("/getforgotpassword",handlers.GetForgotPassword).Methods("GET")
	r.HandleFunc("/sendmail",handlers.SendMail).Methods("POST")
	r.HandleFunc("/getresetpassword",handlers.GetResetPassword).Methods("GET")
	r.HandleFunc("/resetpassword",handlers.ResetPassword).Methods("POST")
	r.HandleFunc("/getprofilepage",handlers.GetProfilePage).Methods("GET")
	r.HandleFunc("/verify2FA",handlers.Verify2FA).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", r))
}
