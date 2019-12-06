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
	//Database - components required to access the db
	Database struct {
		Client       *mongo.Client
		Connected    bool
		databaseName string
		//context
	}

	//DatabaseModel interface
	DatabaseModel interface {
		SetID(id primitive.ObjectID)
		GetID() string
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

// GetClient - Get the mongo Client
func (d *Database) GetClient() *mongo.Client {
	return d.Client
}

// SetDatabase - Set the database name to be used when getting collections
func (d *Database) SetDatabaseName(databaseName string) {
	d.databaseName = databaseName
}

//Insert add a value to the collection
func (d *Database) Insert(collectionName string, model DatabaseModel) error {
	collection := d.GetCollection(collectionName)

	id := primitive.NewObjectID()
	PrintDebug("Primitive", id)
	model.SetID(id)

	insertResult, err := collection.InsertOne(context.TODO(), model)
	if err != nil {
		log.Fatal(err)
		return err
	}

	PrintDebug("Inserted a single document: ", insertResult.InsertedID, insertResult.InsertedID)
	return err

}

// Find - find result via the filter
func (d *Database) Find(collectionName string, filter bson.M) (*mongo.Cursor, error) {
	collection := d.GetCollection(collectionName)

	// Try and find the result and decode into the model pointer variable passed
	//results, err := collection.Find(context.TODO(), filter)
	results, err := collection.Find(context.TODO(), filter)

	PrintDebug("Database Find: ", results)

	return results, err
}

// FindOne - find one result via the filter
func (d *Database) FindOne(collectionName string, filter bson.M, result interface{}) error {
	collection := d.GetCollection(collectionName)

	// Try and find the result and decode into the model pointer variable passed
	err := collection.FindOne(context.TODO(), filter).Decode(result)

	PrintDebug("Database Find One: ", result, err)

	return err
}

//FindOneAndDelete find a filter option and remove it from the db
func (d *Database) FindOneAndDelete(collectionName string, filter bson.M) *mongo.SingleResult {
	collection := d.GetCollection(collectionName)

	options := &options.FindOneAndDeleteOptions{}
	result := collection.FindOneAndDelete(context.TODO(), filter, options)
	return result
}

//FindOneAndUpdate find one result and update
func (d *Database) FindOneAndUpdate(collectionName string, filter bson.M, update interface{}) *mongo.SingleResult {
	collection := d.GetCollection(collectionName)

	options := &options.FindOneAndUpdateOptions{}
	options.SetUpsert(true)

	result := collection.FindOneAndUpdate(context.TODO(), filter, update, options)

	return result
}

//FindOneAndReplace find one result and replace
func (d *Database) FindOneAndReplace(collectionName string, filter bson.M, update interface{}) *mongo.SingleResult {
	collection := d.GetCollection(collectionName)

	result := collection.FindOneAndReplace(context.TODO(), filter, update)

	PrintDebug("ERROR", result.Err())

	return result
}
