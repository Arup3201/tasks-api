package utils

import (
	"fmt"
	"log"
	"os"
)

const (
	PORT   = "PORT"
	DBHOST = "DBHOST"
	DBUSER = "DBUSER"
	DBPORT = "DBPORT"
	DBPASS = "DBPASS"
	DBNAME = "DBNAME"
)

const defaultPort = "80"
const defaultDBPort = "5432"

type envList struct {
	Port   string
	DBHost string
	DBUser string
	DBPort string
	DBPass string
	DBName string
}

var Config = &envList{}

func (eList *envList) Configure() {
	port, ok := os.LookupEnv(PORT)
	if !ok {
		eList.Port = defaultPort
	} else {
		eList.Port = port
	}

	db_host, ok := os.LookupEnv(DBHOST)
	if !ok {
		log.Fatal("DBHOST variable missing in environment variables")
	} else {
		eList.DBHost = db_host
	}

	db_user, ok := os.LookupEnv(DBUSER)
	if !ok {
		log.Fatal("DBUSER variable missing in environment variables")
	} else {
		eList.DBUser = db_user
	}

	db_port, ok := os.LookupEnv(DBPORT)
	if !ok {
		eList.DBPort = defaultDBPort
		fmt.Printf("Default DB Port used: 5432")
	} else {
		eList.DBPort = db_port
	}

	db_pass, ok := os.LookupEnv(DBPASS)
	if !ok {
		log.Fatal("DBPASS variable missing in environment variables")
	} else {
		eList.DBPass = db_pass
	}

	db_name, ok := os.LookupEnv(DBNAME)
	if !ok {
		log.Fatal("DBNAME variable missing in environment variables")
	} else {
		eList.DBName = db_name
	}
}
