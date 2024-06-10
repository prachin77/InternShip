package handlers

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/smtp"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/pquerna/otp/totp"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

var id primitive.ObjectID
var generatedOtpcode string
var secretKey = []byte("secret-key")
var userResendEmail string
var emailOfResetPassword string

type Info struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	Username string             `json:"username"`
	Email    string             `json:"email"`
	Password string             `json:"password"`
}

// type GoogleCaptchaResponse struct {
// 	Success     bool      `json:"success"`
// 	Score       float64   `json:"score"`
// 	Action      string    `json:"action"`
// 	ChallengeTS time.Time `json:"challenge_ts"`
// 	Hostname    string    `json:"hostname"`
// 	ErrorCodes  []string  `json:"error-codes"`
// }

const (
	connectionString = "mongodb://localhost:27017"
	dbName           = "timepass"
	collName         = "task2"
)

// this is a pointer(reference) to collection in mongo db
var collection *mongo.Collection

func init() {
	clientOpt := options.Client().ApplyURI(connectionString)

	// connect to mongo db
	client, err := mongo.Connect(context.TODO(), clientOpt)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("connection to mongo db successfull ✌️✌️")
	collection = client.Database(dbName).Collection(collName)

	// collection instance
	fmt.Println("collection instance is ready")
}

func DefaultRoute(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./templates/main.html"))
	tmpl.Execute(w, nil)
}

func GetApp(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./templates/register.html"))
	tmpl.Execute(w, nil)
}

func CheckUser(info Info) bool {
	filter := bson.M{"email": info.Email}
	var existingUser Info

	err := collection.FindOne(context.TODO(), filter).Decode(&existingUser)
	if err == mongo.ErrNoDocuments {
		// No document found, email doesn't exist
		return false
	} else if err != nil {
		// Some error occurred while querying
		return false
	}

	// Email already exists
	return true
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func ComparePasswords(originalPassword string, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(originalPassword))
	return err == nil
}

func GetRegister(w http.ResponseWriter, r *http.Request) {
	// fmt.Println("get register method type : ",r.Method)
	tmpl := template.Must(template.ParseFiles("./templates/register.html"))
	tmpl.Execute(w, nil)
}

func Register(w http.ResponseWriter, r *http.Request) {
	id = primitive.NewObjectID()
	info := Info{
		ID:       id,
		Username: r.PostFormValue("username"),
		Email:    r.PostFormValue("email"),
		Password: r.PostFormValue("password"),
	}
	originalPassword := info.Password
	fmt.Println("user info : ", info)
	userexists := CheckUser(info)

	if userexists == true {
		fmt.Println("user already exists 😭😭😭")
		fmt.Fprintf(w, "<script>alert('User already exists');</script>")
		// tmpl := template.Must(template.ParseFiles("./templates/register.html"))
		// tmpl.Execute(w, nil)
	} else {
		// Hash the password
		hashedPassword, err := HashPassword(info.Password)
		if err != nil {
			log.Fatal(err)
			return
		}

		// Use the hashed password
		info.Password = hashedPassword

		if ComparePasswords(originalPassword, info.Password) {
			fmt.Println("original password = hashed password")
			fmt.Println("original password : ", originalPassword)
			fmt.Println("hashed password : ", info.Password)
			// insert user in mongo db
			userinfo, err := collection.InsertOne(context.TODO(), info)
			if err != nil {
				log.Fatal(err)
			}
			GetLogin(w, r)
			fmt.Println("user info id : ", userinfo.InsertedID)
			return
		}

	}
}

func GetLogin(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./templates/login.html"))
	tmpl.Execute(w, nil)
}

func Login(w http.ResponseWriter, r *http.Request) {
	info := Info{
		Email:    r.PostFormValue("email"),
		Password: r.PostFormValue("password"),
	}

	filter := bson.M{"email": info.Email}
	var existingUser Info
	err := collection.FindOne(context.TODO(), filter).Decode(&existingUser)
	if err != nil {
		log.Fatal(err)
	}

	if ComparePasswords(info.Password, existingUser.Password) && existingUser.Email == info.Email {
		fmt.Fprintf(w, "<script>alert('login successfull');</script>")
		fmt.Println("User authenticated!")
		fmt.Println("otp verified!")
		fmt.Println("email : ", existingUser.Email)
		fmt.Println("original password : ", info.Password)
		fmt.Println("hashed password : ", existingUser.Password)

		generatedOtpcode = GenerateOtpKey()

		// Send OTP to user's email
		err := SendOtpWithSmtp(existingUser.Email, generatedOtpcode)
		if err != nil {
			log.Println("Error sending OTP:", err)
			// Handle error appropriately
			return
		}

		tmpl := template.Must(template.ParseFiles("./templates/otp.html"))
		tmpl.Execute(w, nil)
	} else {
		fmt.Fprintf(w, "<script>alert('login not successfull');</script>")
		fmt.Println("Incorrect email or password")
		tmpl := template.Must(template.ParseFiles("./templates/login.html"))
		tmpl.Execute(w, nil)
	}
}

func GenerateOtpKey() string {
	// Generate a TOTP key
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "ExampleCorp",
		AccountName: "user@example.com",
	})
	if err != nil {
		fmt.Println("Error generating TOTP key:", err)
		return ""
	}

	// Get the current TOTP code
	code, err := totp.GenerateCode(key.Secret(), time.Now())
	if err != nil {
		fmt.Println("Error generating TOTP code:", err)
		return ""
	}

	fmt.Println("Current TOTP code:", code)
	return code
}

// Function to send OTP using net/smtp
func SendOtpWithSmtp(emailAddr, otpValue string) error {
	from := "prachinnayak07@gmail.com"
	password := "uonw bges rove omhz"
	smtpHost := "smtp.gmail.com"

	msg := []byte("Subject: OTP Verification\r\n" +
		"\r\n" +
		"Your OTP is: " + otpValue)

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err := smtp.SendMail(smtpHost+":587", auth, from, []string{emailAddr}, msg)
	if err != nil {
		return err
	}
	return nil
}

func VerifyOtp(w http.ResponseWriter, r *http.Request) {
	userOtpValue := r.PostFormValue("otp")
	fmt.Println("user otp value : ", userOtpValue)
	fmt.Println("generate otp value : ", generatedOtpcode)
	if userOtpValue == generatedOtpcode {
		fmt.Println("otp verification successfull")
		fmt.Fprintf(w, "<script>alert('otp verification successfull');</script>")
		// return
	} else {
		fmt.Fprintf(w, "<script>alert('otp verification not successfull');</script>")
		tmpl := template.Must(template.ParseFiles("./templates/login.html"))
		tmpl.Execute(w, nil)
	}
}

func GetForgotPassword(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./templates/forgotpassword.html"))
	tmpl.Execute(w, nil)
}

// for active search
// func FindMail(w http.ResponseWriter, r *http.Request){
// 	mailValue := r.PostFormValue("findmail")
// 	fmt.Println("mail value : ",mailValue)

// 	findOptions := options.Find()
// 	filter := bson.M{
// 		"email" : primitive.Regex{Pattern: mailValue, Options: "i"},
// 	}
// 	cursor , err := collection.Find(context.TODO(),filter,findOptions)
// 	if err != nil{
// 		log.Fatal(err)
// 	}
// 	defer cursor.Close(context.TODO())

// 	var users Info

// 	for cursor.Next(context.TODO()){
// 		var user Info
// 		if err := cursor.Decode(&user);err != nil{
// 			fmt.Println("error decoding info : ",err)
// 			continue
// 		}
// 		users = user
// 		break
// 	}
// 	fmt.Println("user : ",users)
// }

func CreateToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func SendForgotMail(email string, tokenString string) {
	from := "prachinnayak07@gmail.com"
	password := "uonw bges rove omhz"
	smtpHost := "smtp.gmail.com"
	userResendEmail = email
	fmt.Println("user resend email id : ", userResendEmail)

	// msg := []byte("Subject: OTP Verification\r\n" +
	// 	"MIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n" +
	// 	"<html><body>" +
	// 	"<p>Dear User,</p>" +
	// 	"<p><h3>YOUR TOKEN: </h3>" + tokenString + "</p>" +
	// 	"<p>You have requested to reset your password. Please click the link below to reset your password:</p>" +
	// 	"<button>reset password</button>" +
	// 	"</body></html>",
	// )
	msg := []byte("Subject: OTP Verification\r\n" +
		"MIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n" +
		"<html><body>" +
		"<p>Dear User,</p>" +
		"<p><h3>YOUR TOKEN: </h3>" + tokenString + "</p>" +
		"<p>You have requested to reset your password. Please click the button below to reset your password:</p>" +
		// "<button onclick=\"sendRequest()\">Reset Password</button>" +
		// "<script>" +
		// 	"function sendRequest() {" +
		// 		"var xhr = new XMLHttpRequest();" +
		// 		"xhr.open('GET', '/getresetpassword', true);" +
		// 		"xhr.send();" +
		// 	"}" +
		// "</script>" +
		"<a href='http://localhost:8080/getresetpassword'><button>Reset Password</button></a>"+
		"</body></html>",
	)

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err := smtp.SendMail(smtpHost+":587", auth, from, []string{email}, msg)
	if err != nil {
		// return
		log.Fatal(err)
	}
	// return
}

func SendMail(w http.ResponseWriter, r *http.Request) {
	mailValue := r.PostFormValue("findmail")
	fmt.Println("mail value : ", mailValue)

	// Define filter to find the email ID
	filter := bson.M{"email": mailValue}

	// Define a variable to store the found user
	var user Info

	// Execute the query to find the user
	err := collection.FindOne(context.TODO(), filter).Decode(&user)
	// if err != nil{
	// 	log.Fatal(err)
	// }
	// fmt.Println("result : ",result)
	if err == mongo.ErrNoDocuments {
		// No document found, handle accordingly
		fmt.Println("No user found with email:", mailValue)
		fmt.Fprintf(w, "<script>alert('mail id not found');</script>")
		tmpl := template.Must(template.ParseFiles("./templates/forgotpassword.html"))
		tmpl.Execute(w, nil)
		return
	} else if err != nil {
		log.Fatal(err)
		return
	}

	// User found, print the user's email ID
	fmt.Println("Found user with email:", user.Email)
	fmt.Fprintf(w, "<script>alert('mail id found');</script>")
	tokenString, err := CreateToken(user.Email)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("token string : ", tokenString)
	SendForgotMail(user.Email, tokenString)

}

func GetResetPassword(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method type : ",r.Method)
	tmpl := template.Must(template.ParseFiles("./templates/resetpassword.html"))
	tmpl.Execute(w, nil)
	// http.Redirect(w, r, "/getresetpassword", http.StatusSeeOther)
}

func ResetPassword(w http.ResponseWriter, r *http.Request) {
	newPassword := r.PostFormValue("newpassword")
	confirmPassword := r.PostFormValue("confirmpassword")
	fmt.Println("new password = ", newPassword)
	fmt.Println("confirm password = ", confirmPassword)

	if newPassword == confirmPassword {
		// Find the user by email
		filter := bson.M{"email": userResendEmail}
		var user Info
		err := collection.FindOne(context.TODO(), filter).Decode(&user)
		if err == mongo.ErrNoDocuments {
			// No user found with the email
			fmt.Fprintf(w, "<script>alert('No user found with the email');</script>")
			return
		} else if err != nil {
			log.Fatal(err)
			return
		}

		// Hash the new password
		hashedPassword, err := HashPassword(newPassword)
		if err != nil {
			log.Fatal(err)
			return
		}
		if ComparePasswords(newPassword, hashedPassword) {
			fmt.Println("original password = hashed password")
			fmt.Println("original password : ", newPassword)
			fmt.Println("hashed password : ", hashedPassword)
			// Update the user's password
			update := bson.M{"$set": bson.M{"password": hashedPassword}}
			_, err = collection.UpdateOne(context.TODO(), filter, update)
			if err != nil {
				log.Fatal(err)
				return
			}
			fmt.Println("password updated successfully")
			return
		}

	} else {
		fmt.Fprintf(w, "<script>alert('Passwords did not match');</script>")
		tmpl := template.Must(template.ParseFiles("./templates/resetpassword.html"))
		tmpl.Execute(w, nil)
	}
}
