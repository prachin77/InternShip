package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	authhandler "github.com/prachin77/task6/auth/AuthHandler"
	"github.com/prachin77/task6/webhandler"
)

func main()  {
	r := mux.NewRouter()

	fmt.Println("listening on port 8080 ....")

	r.HandleFunc("/", authhandler.DefaultRoute).Methods("GET")
	r.HandleFunc("/app", authhandler.GetApp).Methods("GET")

	r.HandleFunc("/getadminauth",authhandler.GetAdminAuth).Methods("GET")
	r.HandleFunc("/verifyadmin",authhandler.VerifyAdmin).Methods("POST")
	r.HandleFunc("/admin/getadminpanel",webhandler.GetAdminPanel).Methods("GET")
	r.HandleFunc("/admin/getadduserform",webhandler.GetAddAdminUserForm).Methods("GET")
	r.HandleFunc("/admin/adduser",webhandler.AdminAddUser).Methods("POST")
	r.HandleFunc("/admin/deleteuser/{id}",authhandler.DeleteUser).Methods("DELETE")
	r.HandleFunc("/admin/getupdateuser/{id}",webhandler.GetAdminUpdateUserForm).Methods("PUT")
	r.HandleFunc("/admin/updateuser/{id}",authhandler.UpdateUser).Methods("POST")

	r.HandleFunc("/user/auth/getlogin",authhandler.GetUserAuthLoginForm).Methods("GET")
	r.HandleFunc("/user/auth/getregister",authhandler.GetUserAuthRegisterForm).Methods("GET")
	r.HandleFunc("/user/auth/login",authhandler.UserAuthLogin).Methods("POST")
	r.HandleFunc("/user/auth/register",authhandler.UserAuthRegister).Methods("POST")


	log.Fatal(http.ListenAndServe(":8080", r))

}
