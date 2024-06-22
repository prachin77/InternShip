package models

type User struct {
	UserName  string
	Email     string
	Password  string
	UserId    string
	NoOfBlogs int
	IsUser    bool
	IsLogin   bool
}

type Admin struct {
	Email    string
	Password string
}
