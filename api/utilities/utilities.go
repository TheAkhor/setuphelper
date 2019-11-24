package utilities

import (
	"fmt"
	"log"
)

// var (
// 	clientGlobal *mongo.Client
// )

var (
	DatabaseObj = new(Database)
)

func Init() {
	DatabaseObj.Connect()
	DatabaseObj.SetDatabaseName("MyClassRoom")
}

//"github.com/labstack/echo/middleware"

// PrintDebug - Send the incoming data to the log file
func PrintDebug(i ...interface{}) {
	log.Print("Debug Print - ", fmt.Sprintf("%+v\n", i))
}

func GetDatabaseStruct() *Database {
	return DatabaseObj
}

// func DatabaseConnect() {
// 	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
// 	client, err := mongo.NewClient(clientOptions)
// 	clientGlobal = client
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	err = clientGlobal.Connect(context.TODO())
// 	if err != nil {
// 		PrintDebug("Connection Failed")
// 		log.Fatal(err)
// 	} else {
// 		PrintDebug("Database - Connection Success")
// 	}

// 	// Check the connection
// 	if err := clientGlobal.Ping(context.Background(), nil); err != nil {
// 		PrintDebug("DatabaseConnect - Could not Connect to MongoDB!")

// 	} else {

// 		PrintDebug("DatabaseConnect - Connected MongoDB!")
// 	}

// }

// func DatabaseGetClient() *mongo.Client {
// 	return clientGlobal
// }
