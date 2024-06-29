package authhandler

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/prachin77/task6/models"
)

// pointer to mysql db
var db *sql.DB

func init() {
	// Load environment vari""ables from .env file
	if err := godotenv.Load("/InternShip/Task6/auth/.env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Assign loaded values to variables
	DriverName := os.Getenv("DriverName")
	DataSource := os.Getenv("DataSource")

	var err error
	db, err = sql.Open(DriverName, DataSource)
	if err != nil {
		log.Fatal(err)
	}
	// defer db.Close()

	fmt.Println("successfully connected to MySql 👍👍")
}

func GetAllUsers() ([]models.User, error) {
	var users []models.User

	// query := "SELECT * FROM userinfo"
	query := "SELECT userid, username, email, noofblogs FROM userinfo"

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.UserId, &user.UserName, &user.Email, &user.NoOfBlogs); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	if len(users) == 0 {
		return nil, fmt.Errorf("no users found")
	}

	return users, nil
}

// admin adds this user
func AddUser(user models.User) (models.User, error) {
	fmt.Println("user value in add user func = ", user)

	// prepare statment
	stmt, err := db.Prepare("INSERT INTO userinfo(username, email, password, userid, noofblogs,isuser) VALUES(?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatalf("Error preparing statement: %v", err)
		// return models.User{}, err
		return models.User{}, err
	}

	// execute statement
	_, err = stmt.Exec(user.UserName, user.Email, user.Password, user.UserId, user.NoOfBlogs,user.IsUser)
	if err != nil {
		log.Fatalf("Error executing query: %v", err)
		return models.User{}, err
	}

	return user, nil
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	// Extract task ID from request URL
	userid := mux.Vars(r)["id"]
	fmt.Println("user id for delete function : ", userid)

	query := "DELETE FROM userinfo WHERE userid = ?"
	// Prepare SQL statement
	stmt, err := db.Prepare(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	// Execute SQL statement
	_, err = stmt.Exec(userid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("user delete succesfully ")
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	// extract user id from url
	userid := mux.Vars(r)["id"]
	fmt.Println("user id for update func : ", userid)

	userinfo := models.User{
		UserName: r.PostFormValue("username"),
		Email:    r.PostFormValue("email"),
		Password: r.PostFormValue("password"),
	}
	fmt.Println("values from admin to update user : ", userinfo)

	// Prepare SQL statement to update user
	query := "UPDATE userinfo SET username=?, email=?, password=? WHERE userid=?"
	stmt, err := db.Prepare(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close() // Ensure statement is closed after function ends

	// Execute SQL statement with user data
	_, err = stmt.Exec(userinfo.UserName, userinfo.Email, userinfo.Password, userid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("User updated successfully")
	users, err := GetAllUsers()
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

func CheckUserInDb(user models.User) (bool, error) {
    var count int
    query := "SELECT COUNT(*) FROM userinfo WHERE email = ? OR username = ?"
    err := db.QueryRow(query, user.Email, user.UserName).Scan(&count)
    if err != nil {
        return false, err
    }
    return count > 0, nil
}



