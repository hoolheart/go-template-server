package database

import (
	"os"
	"database/sql"
	"encoding/json"
	"fmt"
	_"github.com/lib/pq" // import postgre driver
	_"github.com/go-sql-driver/mysql" // import mysql driver 
	//_"github.com/mattn/go-sqlite3" // import sqlite3 driver
)

// TemplateDb is the database connection to template application
var TemplateDb *sql.DB

// dbCfg is the configuration of database connection
type dbCfg struct {
	Driver string `json:"driver"` // driver name
	Host  string `json:"host"` // database host
	Port int `json:"port"` // connection port, <=0 means default
	User string `json:"user"` // login user
	Password string `json:"password"` // login password
	DbName string `json:"database"` // database name
	SslMode string `json:"sslmode"` // ssl mode
}

// init function initialize the database connection for template application
func init() {
	var cfg dbCfg//prepare configuration
	file, err := os.Open("assets/database_cfg.json")//open configure file
	if err == nil {
		defer file.Close()
		err = json.NewDecoder(file).Decode(&cfg)//parse setting from configure
	}

	if err != nil {
		fmt.Println("Default database driver: postgres and the configuration: dbname=db_template password=142857 sslmode=disable")//print info
		TemplateDb, err = sql.Open("postgres","dbname=db_template password=142857 sslmode=disable")//default parameter
	} else {
		var dbStr string//prepare database string
		if cfg.Driver == "postgres" {//postgresql
			dbStr = fmt.Sprintf("user=%s dbname=%s password=%s sslmode=%s",cfg.User,cfg.DbName,cfg.Password,cfg.SslMode)//login info
			if len(cfg.Host)>0 {
				dbStr = dbStr + fmt.Sprintf(" host=%s",cfg.Host)//host
			}
			if cfg.Port>0 {
				dbStr = dbStr + fmt.Sprintf(" port=%d",cfg.Port);//port
			}
		} else if cfg.Driver == "mysql" {//mysql
			url := cfg.Host//host
			if cfg.Port>0 {
				url += url + fmt.Sprintf(":%d",cfg.Port)//port
			}
			dbStr = fmt.Sprintf("%s:%s@%s/%s?charset=utf8", cfg.User, cfg.Password, url, cfg.DbName)//login info
		} else if cfg.Driver == "sqlite3" {//sqlite3
			dbStr = cfg.DbName//use only the databas name
		} else {
			panic("Database: Invalid driver " + cfg.Driver)
		}
		
		fmt.Println("Database driver: " + cfg.Driver + " and the configuration: " + dbStr)//print info
		TemplateDb, err = sql.Open(cfg.Driver,dbStr)
	}
	if err != nil {
		panic(err)
	}
	fmt.Println("Succeed to connect database")
}
