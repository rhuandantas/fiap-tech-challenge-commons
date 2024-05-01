package repository

import (
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/gommon/log"
	"xorm.io/xorm"
	"xorm.io/xorm/names"
)

type DBConnector interface {
	GetORM() *xorm.Engine
	Close()
	SyncTables(beans ...interface{}) error
}

type MySQLConnector struct {
	engine *xorm.Engine
}

func (m *MySQLConnector) GetORM() *xorm.Engine {
	return m.engine
}

func (m *MySQLConnector) Close() {
	err := m.engine.Close()
	if err != nil {
		log.Fatal(err.Error())
	}
}

func NewMySQLConnector() DBConnector {
	// TODO put in env vars
	var (
		dbName     string
		dbPassword string
		dbUser     string
		dbPort     string
		dbHost     string
		err        error
	)

	dbHost = os.Getenv("DB_HOST")
	dbPassword = os.Getenv("DB_PASS")
	dbName = os.Getenv("DB_NAME")
	dbUser = os.Getenv("DB_USER")
	dbPort = os.Getenv("DB_PORT")
	if dbHost == "" || dbPassword == "" || dbName == "" || dbUser == "" || dbPort == "" {
		log.Fatal("make sure your db variable are configured properly")
	}

	engine, err := xorm.NewEngine("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", dbUser, dbPassword, dbHost, dbPort, dbName))
	if err != nil {
		panic(err)
	}
	engine.ShowSQL(true) // TODO it should come from env
	//engine.Logger().SetLevel(log.DEBUG)
	engine.SetMapper(names.SnakeMapper{})

	return &MySQLConnector{
		engine: engine,
	}
}

func (m *MySQLConnector) SyncTables(beans ...interface{}) error {
	if err := m.engine.Sync(beans); err != nil {
		return err
	}

	return nil
}
