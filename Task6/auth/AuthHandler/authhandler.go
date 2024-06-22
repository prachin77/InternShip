package authhandler

import (
	"crypto/rand"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	"github.com/prachin77/task6/models"
)

type Info struct {
	IsAdmin bool
}

var (
	secretKey = []byte("secret-key")
)

func GenerateToken() (string, error) {
	b := make([]byte, 4)
	rand.Read(b)
	return fmt.Sprintf("%x", b), nil
}

func DefaultRoute(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("main.html"))
	tmpl.Execute(w, nil)
}

func SetCookie(w http.ResponseWriter, r *http.Request, IssuerName string, tokenString string) {
	fmt.Println("issuer name : ", IssuerName)
	cookie := http.Cookie{
		Name:     IssuerName,
		Value:    tokenString,
		Path:     "/",
		MaxAge:   600,
		HttpOnly: true,
		Secure:   true,
	}

	http.SetCookie(w, &cookie)

}

func GetCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie("sessiontoken")
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

func GetApp(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("/InternShip/Task6/webtemplates/homepage.html"))
	tmpl.Execute(w, nil)
}

func GetAdminAuth(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("/InternShip/Task6/auth/AuthTemplates/adminauth.html"))
	tmpl.Execute(w, nil)
}

func VerifyAdmin(w http.ResponseWriter, r *http.Request) {
	adminInfo := models.Admin{
		Email:    r.PostFormValue("email"),
		Password: r.PostFormValue("password"),
	}

	// Load environment vari""ables from .env file
	if err := godotenv.Load("/InternShip/Task6/auth/.env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	AdminEmail := os.Getenv("AdminEmail")
	AdminPassword := os.Getenv("AdminPassword")

	if adminInfo.Email == AdminEmail && adminInfo.Password == AdminPassword {
		info := Info{
			IsAdmin: true,
		}

		// GENERATE TOKEN
		token := jwt.NewWithClaims(jwt.SigningMethodHS256,
			jwt.MapClaims{
				"username": "admin",
				"exp":      time.Now().Add(time.Hour * 24).Unix(),
			})

		tokenString, err := token.SignedString(secretKey)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("token string : ", tokenString)

		// set cookie for admin
		SetCookie(w, r, "admin", tokenString)

		fmt.Println("admin authentication successfull 😑😑")
		fmt.Println("admin details : ", adminInfo)
		tmpl := template.Must(template.ParseFiles("/InternShip/Task6/webtemplates/homepage.html"))
		// tmpl.Execute(w, nil)
		// tmpl.Execute(w, adminInfo)
		tmpl.Execute(w, info)
		return
	}
}

func GetUserAuthLoginForm(w http.ResponseWriter, r *http.Request) {
	userAuthStatus := models.User{
		IsLogin: true,
	}
	tmpl := template.Must(template.ParseFiles("/InternShip/Task6/auth/AuthTemplates/userauth.html"))
	tmpl.Execute(w, userAuthStatus)
}

func GetUserAuthRegisterForm(w http.ResponseWriter, r *http.Request) {
	userAuthStatus := models.User{
		IsLogin: false,
	}
	tmpl := template.Must(template.ParseFiles("/InternShip/Task6/auth/AuthTemplates/userauth.html"))
	tmpl.Execute(w, userAuthStatus)
}

func UserAuthLogin(w http.ResponseWriter, r *http.Request) {
	userinfo := models.User{
		Email:    r.PostFormValue("email"),
		Password: r.PostFormValue("password"),
	}
	isPresent, err := CheckUserInDb(userinfo)
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if isPresent {
		fmt.Println("login user id : ", userinfo.UserId)
		fmt.Println("login username : ", userinfo.UserName)
		fmt.Println("login user email : ", userinfo.Email)
		fmt.Println("login user password : ", userinfo.Password)
		fmt.Println("login user blog count : ", userinfo.NoOfBlogs)

		tmpl := template.Must(template.ParseFiles("/InternShip/Task6/webtemplates/homepage.html"))
		tmpl.Execute(w, nil)
	} else {
		fmt.Fprintf(w, `<alert>user not found</alert>`)
	}
}

func UserAuthRegister(w http.ResponseWriter, r *http.Request) {
	userinfo := models.User{
		UserName:  r.PostFormValue("username"),
		Email:     r.PostFormValue("email"),
		Password:  r.PostFormValue("password"),
		NoOfBlogs: 0,
	}

	// Check if user is already present in db
	isPresent, err := CheckUserInDb(userinfo)
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// If user is already present, return an error or handle accordingly
	if isPresent {
		http.Error(w, "User already exists", http.StatusBadRequest)
		return
	}

	// Generate token for the user
	tokenString, err := GenerateToken()
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	userinfo.UserId = tokenString
	userinfo.IsUser = true

	// Add user to db
	user, err := AddUser(userinfo)
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	fmt.Println("Registration successful")
	fmt.Println("User registered : ", user)

	tmpl := template.Must(template.ParseFiles("/InternShip/Task6/webtemplates/homepage.html"))
	tmpl.Execute(w, userinfo)
}
