package webhandler

import (
	"crypto/rand"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	authhandler "github.com/prachin77/task6/auth/AuthHandler"
	"github.com/prachin77/task6/models"
)

func GenerateToken() (string, error) {
	b := make([]byte, 4)
	rand.Read(b)
	return fmt.Sprintf("%x", b), nil
}

func GetAdminPanel(w http.ResponseWriter, r *http.Request) {
	users, err := authhandler.GetAllUsers()
	if err != nil {
		log.Printf("Error retrieving users: %v", err)
		http.Error(w, "Failed to retrieve users", http.StatusInternalServerError)
		return
	}
	fmt.Println("all users = ", users)
	tmpl := template.Must(template.ParseFiles(
		"/InternShip/Task6/webtemplates/adminpanel.html",
		"/InternShip/Task6/webtemplates/usertable.html",
	))
	err = tmpl.Execute(w, users)
	if err != nil {
		log.Fatal(err)
	}
}

func GetAddAdminUserForm(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("/InternShip/Task6/webtemplates/adminadduserform.html"))
	tmpl.Execute(w, nil)
}

func AdminAddUser(w http.ResponseWriter, r *http.Request) {
	tokenValue, err := GenerateToken()
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("TOKEN STRING : ", tokenValue)
	user := models.User{
		UserName:  r.PostFormValue("username"),
		Email:     r.PostFormValue("email"),
		Password:  r.PostFormValue("password"),
		UserId:    tokenValue,
		IsUser: true,
		NoOfBlogs: 0,
	}
	fmt.Println("user value inserted by admin : ", user)

	// now add user to mysql db in db folder
	// value , err := db.AddUser(user)
	value, err := authhandler.AddUser(user)
	if err != nil {
		log.Fatal(err)
	} else {
		tmpl := template.Must(template.ParseFiles("/InternShip/Task6/webtemplates/usertable.html"))
		tmpl.Execute(w, value)
	}
}

func GetAdminUpdateUserForm(w http.ResponseWriter, r *http.Request) {
	// Extract task ID from request URL
	userid := mux.Vars(r)["id"]
	fmt.Println("user id for update function : ", userid)
	userinfo := models.User{
		UserId: userid,
	}
	tmpl := template.Must(template.ParseFiles("/InternShip/Task6/webtemplates/adminupdateuserform.html"))
	tmpl.Execute(w, userinfo)
}
