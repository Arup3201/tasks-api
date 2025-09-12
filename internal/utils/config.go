package utils

import (
	"fmt"

	"github.com/nathanbcrocker/env"
)

const (
	PORT   = "PORT"
	DBHOST = "DBHOST"
	DBUSER = "DBUSER"
	DBPORT = "DBPORT"
	DBPASS = "DBPASS"
	DBNAME = "DBNAME"
)

const defaultPort = "8080"
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
	e := env.NewEnv()
	port, ok := e.Get(PORT)
	if !ok {
		eList.Port = defaultPort
	} else {
		eList.Port = port.Value
	}

	db_host, ok := e.Get(DBHOST)
	if !ok {
		panic("DBHOST variable missing in .env")
	} else {
		eList.DBHost = db_host.Value
	}

	db_user, ok := e.Get(DBUSER)
	if !ok {
		panic("DBUSER variable missing in .env")
	} else {
		eList.DBUser = db_user.Value
	}

	db_port, ok := e.Get(DBPORT)
	if !ok {
		eList.DBPort = defaultDBPort
		fmt.Printf("Default DB Port used: 5432")
	} else {
		eList.DBPort = db_port.Value
	}

	db_pass, ok := e.Get(DBPASS)
	if !ok {
		panic("DBPASS variable missing in .env")
	} else {
		eList.DBPass = db_pass.Value
	}

	db_name, ok := e.Get(DBNAME)
	if !ok {
		panic("DBNAME variable missing in .env")
	} else {
		eList.DBName = db_name.Value
	}
}
