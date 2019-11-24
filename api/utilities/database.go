package utilities

import (
	"context"
	"log"

	//_ "go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	Database struct {
		Client       *mongo.Client
		Connected    bool
		databaseName string
		//context
	}

	DatabaseModel interface {
		SetID(id primitive.ObjectID)
		//SetID(id string)
	}
)

//Connect - Connect to the database and set the Connected flag
func (d *Database) Connect() {

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.NewClient(clientOptions)
	d.Client = client
	if err != nil {
		log.Fatal(err)
	}
	err = d.Client.Connect(context.TODO())
	if err != nil {
		PrintDebug("Connection Failed")
		log.Fatal(err)
	} else {
		PrintDebug("Database - Connection Success")
		d.Connected = true
	}

	// Check the connection
	if err := d.Client.Ping(context.Background(), nil); err != nil {
		PrintDebug("DatabaseConnect - Could not Connect to MongoDB!")

	} else {

		PrintDebug("DatabaseConnect - Connected MongoDB!")
	}

}

// GetCollection - Send a Collection name to get the pointer to the collection
func (d *Database) GetCollection(collection string) *mongo.Collection {
	return d.Client.Database(d.databaseName).Collection(collection)
}

// GetClient - Get teh mongo Client
func (d *Database) GetClient() *mongo.Client {
	// if !d.Connected {
	// 	return nil
	// }
	return d.Client
}

// SetDatabase - Set the database name to be used when getting collections
func (d *Database) SetDatabaseName(databaseName string) {
	d.databaseName = databaseName
}

// SetDatabase - Set the database name to be used when getting collections
//func (d *Database) Insert(collectionName string, model interface{}) {
func (d *Database) Insert(collectionName string, model DatabaseModel) {
	collection := d.GetCollection(collectionName)

	//model.CreationTime.UpdateTimes()
	//SetTImes(model)

	insertResult, err := collection.InsertOne(context.TODO(), model)
	if err != nil {
		log.Fatal(err)
	}

	PrintDebug("Inserted a single document: ", insertResult.InsertedID, insertResult.InsertedID)
	//model.SetID(insertResult.InsertedID)
	//model.SetID(fmt.Sprintf("%v", insertResult.InsertedID))

	test := insertResult.InsertedID.(primitive.ObjectID)
	PrintDebug("Object ID TEST", test.String(), test.Hex())

	// On the model "interface" we need to type assert the response
	// This will send a primitive.ObjectID back to SetID.
	model.SetID(insertResult.InsertedID.(primitive.ObjectID))

	//model.SetID(insertResult.InsertedID))
	//return insertResult, err

}

// FindOne - find one result via the filter
func (d *Database) FindOne(collectionName string, filter bson.M, result interface{}) error {
	collection := d.GetCollection(collectionName)

	// Create an empty interface to store the results of the Decode
	// Then in the model use "type assert" to load the data
	err := collection.FindOne(context.TODO(), filter).Decode(result)

	if err != nil {
		log.Fatal(err)
	}

	PrintDebug("Database Find One: ", result, err)

	return err
}

// func SetTimes(m models.UpdateModelTimes) {
// 	m.UpdateTimes()
// }
