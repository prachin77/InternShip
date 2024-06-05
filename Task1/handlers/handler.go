package handlers

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var id primitive.ObjectID

type Task struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	TaskValue string             `json:"taskvalue"`
	TaskDate  string             `json:"taskdate"`
	// Add other fields as needed
}

const (
	connectionString = "mongodb://localhost:27017"
	dbName           = "timepass"
	collName         = "todolist"
)

// this is a pointer(reference) to collection in mongo db
var collection *mongo.Collection

func init() {
	clientOpt := options.Client().ApplyURI(connectionString)

	// connect to mongodb
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
	// get all task from mongo db
	tasks, err := GetAllTasks(w, r)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("tasks from app page : ", tasks)
		tmpl := template.Must(template.ParseFiles("./templates/todo.html", "./templates/todocompo.html"))
		tmpl.Execute(w, tasks)
		// tmpl.Execute(w, nil)
	}
}

func GetAllTasks(w http.ResponseWriter, r *http.Request) ([]Task, error) {
	// Create a slice to hold the tasks
	var tasks []Task

	// Define a context for the operation
	ctx := context.TODO()

	// Define options to customize the query
	findOptions := options.Find()

	// Find all documents in the collection with specified options
	cursor, err := collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Iterate over the cursor and decode each document into a Task
	for cursor.Next(ctx) {
		var task Task
		if err := cursor.Decode(&task); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	// Check for errors during cursor iteration
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	// Return the slice of tasks
	return tasks, nil
}

func AddTask(w http.ResponseWriter, r *http.Request) {
	id = primitive.NewObjectID()
	task := Task{
		ID:        id,
		TaskValue: r.PostFormValue("addtask"),
		TaskDate:  time.Now().Format("2006-01-02"),
	}
	// Insert the task into the MongoDB collection
	_, err := collection.InsertOne(context.TODO(), task)
	if err != nil {
		// http.Error(w, "Failed to add task to database", http.StatusInternalServerError)
		log.Fatal(err)
		// return
	} else {
		tmpl := template.Must(template.ParseFiles("./templates/todocompo.html"))
		tmpl.Execute(w, task)
	}

}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	// Extract task ID from request URL
	taskId := mux.Vars(r)["id"]
	fmt.Println("task id : ", taskId)

	// Convert task ID string to ObjectID
	objId, err := primitive.ObjectIDFromHex(taskId)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("task id with hex converted value = ", objId)
	filter := bson.M{"_id": objId	}

	_, err = collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}

}

func SearchTask(w http.ResponseWriter, r *http.Request) {
	searchValue := r.PostFormValue("search")
	fmt.Println("search value : ", searchValue)
	if searchValue == "" {
		fmt.Println("search value is empty")
		allTaskList, err := GetAllTasks(w, r)
		if err != nil {
			http.Error(w, "Failed to retrieve all tasks", http.StatusInternalServerError)
			fmt.Println("Error retrieving all tasks:", err)
			return
		}
		fmt.Println("all tasks : ",allTaskList)
		tmpl := template.Must(template.ParseFiles("./templates/todocompo.html"))
        for _, task := range allTaskList {
            tmpl.Execute(w, task) // Pass each task to the template
        }
		return
	}

	// Define a context for the operation
	ctx := context.TODO()

	// Define options to customize the query
	findOptions := options.Find()

	// Define a filter for the search query
	filter := bson.M{"taskvalue": primitive.Regex{Pattern: searchValue, Options: "i"}}

	// Find documents in the collection that match the filter
	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		http.Error(w, "Failed to search tasks", http.StatusInternalServerError)
		fmt.Println("Error searching tasks:", err)
		return
	}
	defer cursor.Close(ctx)

	// Create a slice to hold the search results
	// var tasks []Task
	var tasks Task

	// Iterate over the cursor and decode each document into a Task
	for cursor.Next(ctx) {
		var task Task
		if err := cursor.Decode(&task); err != nil {
			fmt.Println("Error decoding task:", err)
			continue
		}
		tasks = task
		break
	}
	fmt.Println("tasks from search function : ", tasks)

	tmpl := template.Must(template.ParseFiles("./templates/todocompo.html"))
	tmpl.Execute(w, tasks)
}

func UpdateTask(w http.ResponseWriter, r *http.Request){
	// Extract task ID from request URL
	taskId := mux.Vars(r)["id"]
	fmt.Println("task id from get update : ", taskId)
	tmpl := template.Must(template.ParseFiles("./templates/updateinputbox.html"))
	// tmpl.Execute(w,nil)
	tmpl.Execute(w,taskId)
}

func UpdatePost(w http.ResponseWriter, r *http.Request){
    // Extract task ID from request URL
    taskId := mux.Vars(r)["id"]
    newUpdateValue := r.PostFormValue("newupdatevalue")
    fmt.Println("task id from post update : ", taskId)
    fmt.Println("new update value : ", newUpdateValue)

    // Convert task ID string to ObjectID
    objId, err := primitive.ObjectIDFromHex(taskId)
    if err != nil {
        log.Fatal(err)
        return
    }

    // Define a filter to find the task by ID
    filter := bson.M{"_id": objId}

    // Define an update operation to set the new task value and date
    update := bson.M{
        "$set": bson.M{
            "taskvalue": newUpdateValue,
            "taskdate":  time.Now().Format("2006-01-02"),
        },
    }

    // Perform the update operation
    _, err = collection.UpdateOne(context.TODO(), filter, update)
    if err != nil {
        log.Fatal(err)
        return
    }

    // Fetch the updated task from the database
    var updatedTask Task
    err = collection.FindOne(context.TODO(), filter).Decode(&updatedTask)
    if err != nil {
        log.Fatal(err)
        return
    }

    // Render the template with the updated task
    tmpl := template.Must(template.ParseFiles("./templates/todocompo.html"))
    tmpl.Execute(w, updatedTask)
}
