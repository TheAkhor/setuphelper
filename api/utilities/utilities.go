package utilities

import (
	"fmt"
	"log"
)

var (
	DatabaseObj = new(Database)
)

//Init - Initialization functions
func Init() {
	DatabaseObj.Connect()
	DatabaseObj.SetDatabaseName("MyClassRoom")
}

// PrintDebug - Send the incoming data to the log file
func PrintDebug(i ...interface{}) {
	log.Print("Debug Print - ", fmt.Sprintf("%+v\n", i))
}

//GetDatabaseStruct - Return the Struct
func GetDatabaseStruct() *Database {
	return DatabaseObj
}
