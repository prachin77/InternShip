package handlers

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/smtp"
	"time"

	"github.com/pquerna/otp/totp"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

var id primitive.ObjectID
var generatedOtpcode string

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

func VerifyOtp(w http.ResponseWriter , r* http.Request){
	userOtpValue := r.PostFormValue("otp")
	fmt.Println("user otp value : ",userOtpValue)
	fmt.Println("generate otp value : ",generatedOtpcode)
	if userOtpValue == generatedOtpcode{
		fmt.Println("otp verification successfull")
		fmt.Fprintf(w, "<script>alert('otp verification successfull');</script>")
		// return
	}else{
		fmt.Fprintf(w, "<script>alert('otp verification not successfull');</script>")
		// tmpl := template.Must(template.ParseFiles("./templates.login.html"))
		// tmpl.Execute(w,nil)
		tmpl := template.Must(template.ParseFiles("./templates/login.html"))
		tmpl.Execute(w, nil)
	}
}

// func Capthca() {
	// 1. random char generate
	// 2. return number from register function
	// 3. call captcha
// }

