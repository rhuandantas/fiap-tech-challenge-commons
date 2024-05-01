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

type MongoDBConnector interface {
	GetDB() *mongo.Database
	Close()
}

type MongoConnector struct {
	engine *mongo.Client
	db     *mongo.Database
}

type MongoEnvs struct {
	DbURI  string
	DbName string
}

func (m *MongoConnector) GetDB() *mongo.Database {
	return m.db
}

func (m *MongoConnector) Close() {
	err := m.engine.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err.Error())
	}
}

func getDbEnvs() MongoEnvs {
	dbURI := os.Getenv("DB_URI")
	dbName := os.Getenv("DB_NAME")

	if dbURI == "" {
		panic("DB_URI não encontrado")
	}
	if dbName == "" {
		panic("DB_NAME não encontrado")
	}

	return MongoEnvs{
		DbURI:  dbURI,
		DbName: dbName,
	}
}

func NewMongoConnector() MongoDBConnector {
	var (
		err error
	)

	envs := getDbEnvs()

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(envs.DbURI).SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}

	db := client.Database(envs.DbName)
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
