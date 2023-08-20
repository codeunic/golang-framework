package framework

import (
	"database/sql"
	"fmt"
	"log"
)

type Db struct {
	db     *sql.DB
	config ConfigDatabase
}

type ConfigDatabase struct {
	User     string
	DbName   string
	Password string
	Host     string
	Port     string
}

func NewDatabase(config ConfigDatabase) *Db {
	return &Db{config: config}
}

func (d *Db) Run() *sql.DB {
	if d.db != nil {
		return d.db
	}

	connStr := "user=%s dbname=%s password=%s host=%s sslmode=disable port=%s"
	db, err := sql.Open("postgres", fmt.Sprintf(connStr,
		d.config.User,
		d.config.DbName,
		d.config.Password,
		d.config.Host,
		d.config.Port,
	))

	if err != nil {
		log.Fatal(err)
	}

	//defer db.Close()

	err = db.Ping()

	if err != nil {
		log.Fatal(err)
	}

	d.db = db
	fmt.Println("[DATABASE] Successfully connected!")
	return db
}

func (d *Db) GetDb() *sql.DB {
	return d.db
}
