package repository

import (
	"context"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
)

type DBConnector interface {
	GetDB() *mongo.Database
	Close()
}

type MongoConnector struct {
	engine *mongo.Client
	db     *mongo.Database
}

func (m *MongoConnector) GetDB() *mongo.Database {
	return m.db
}

func (m *MongoConnector) Close() {
	// TODO change ctx
	err := m.engine.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err.Error())
	}
}

func NewMySQLConnector() DBConnector {
	var (
		dbURI  = os.Getenv("DB_URI")
		dbName = os.Getenv("DB_NAME")
		err    error
	)

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(dbURI).SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}

	db := client.Database(dbName)
	// Send a ping to confirm a successful connection
	var result bson.M
	if err = db.RunCommand(context.TODO(), bson.D{{"ping", 1}}).Decode(&result); err != nil {
		panic(err)
	}
	fmt.Println("You successfully connected to MongoDB!")
	return &MongoConnector{
		engine: client,
		db:     db,
	}
}
