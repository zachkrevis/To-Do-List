package middleware

import(
	"net/http"
	"log"
	"os"
	"github.com/joho/godotenv"
	"fmt"
	"encoding/json"
	"context"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/zachkrevis/to-do-golang/models"
)

var collection *mongo.Collection

func init() {
	loadTheEnv()
	createDBInstance()
}

func loadTheEnv() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading env file.")
	}
}

func createDBInstance() {
	connectionString := os.Getenv("DB_URI")
	dbName := os.Getenv("DB_NAME")
	collectionName := os.Getenv("COLLECTION_NAME")

	clientOptions := options.Client().ApplyURI(connectionString)
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	e := client.Ping(context.TODO(), nil)

	if e != nil {
		log.Fatal(e)
	}

	fmt.Println("Connected to MongoDB")

	collection = client.Database(dbName).Collection(collectionName)
	fmt.Println("collection instance created.")
}

func GetAllTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-ALlow-Origin", "*")
	payload := getAllTasks()

	json.NewEncoder(w).Encode(payload)
}

func CreateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-ALlow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-ALlow-Headers", "Content-Type")

	var task models.ToDoList

	json.NewDecoder(r.Body).Decode(&task)

	
	if(task.Task != "") {
		insertOneTask(task)
	}
	json.NewEncoder(w).Encode(task)


}

func CompleteTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-ALlow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	w.Header().Set("Access-Control-ALlow-Headers", "Content-Type")

	params := mux.Vars(r)
	completeTask(params["id"])
	json.NewEncoder(w).Encode(params["id"])



}

func UndoTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-ALlow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	w.Header().Set("Access-Control-ALlow-Headers", "Content-Type")

	params := mux.Vars(r)
	undoTask(params["id"])
	json.NewEncoder(w).Encode(params["id"])


}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-ALlow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")
	w.Header().Set("Access-Control-ALlow-Headers", "Content-Type")

	params := mux.Vars(r)
	deleteOneTask(params["id"])

}

func DeleteAllTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-ALlow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")
	w.Header().Set("Access-Control-ALlow-Headers", "Content-Type")

	count := deleteAllTasks()
	json.NewEncoder(w).Encode(count)
}

func getAllTasks() []primitive.M{
	curr, err := collection.Find(context.Background(), bson.D{{}})

	if err != nil {
		log.Fatal(err)
	}

	var results []primitive.M

	for curr.Next(context.Background()) {
		var result bson.M
		err := curr.Decode(&result)

		if err != nil {
			log.Fatal(err)
		}
		results = append(results, result)
	}

	if err := curr.Err(); err != nil {
		log.Fatal(err)
	}

	curr.Close(context.Background())
	return results
}

func insertOneTask(task models.ToDoList) {
	insertResult, err := collection.InsertOne(context.Background(), task)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted a single record. ID: ", insertResult.InsertedID)
}

func completeTask(task string) {
	id, _ := primitive.ObjectIDFromHex(task)
	filter := bson.M{"_id":id}
	update := bson.M{"$set":bson.M{"status":true}}
	result, err := collection.UpdateOne(context.Background(), filter, update)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("modified count: ", result.ModifiedCount)
}

func undoTask(task string) {
	id, _ := primitive.ObjectIDFromHex(task)
	filter := bson.M{"_id":id}
	update := bson.M{"$set":bson.M{"status":false}}
	result, err := collection.UpdateOne(context.Background(), filter, update)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("modified count: ", result.ModifiedCount)
}

func deleteOneTask(task string) {
	id, _ := primitive.ObjectIDFromHex(task)
	filter := bson.M{"_id":id}
	d, err := collection.DeleteOne(context.Background(), filter)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Deleted Document", d.DeletedCount )
}

func deleteAllTasks() int64{
	d, err := collection.DeleteMany(context.Background(), bson.D{{}})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Deleted document", d.DeletedCount)
	return d.DeletedCount
}