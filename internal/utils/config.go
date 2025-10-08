package utils

import (
	"fmt"
	"log"
	"os"
)

const (
	DBHOST           = "DBHOST"
	DBUSER           = "DBUSER"
	DBPORT           = "DBPORT"
	DBPASS           = "DBPASS"
	DBNAME           = "DBNAME"
	KEYCLOAK_JWT_URL = "KEYCLOAK_JWT_URL"
)

const defaultPort = "8080"
const defaultDBPort = "5432"

type envList struct {
	Port             string
	DBHost           string
	DBUser           string
	DBPort           string
	DBPass           string
	DBName           string
	KEYCLOAK_JWT_URL string
}

var Config = &envList{}

func (eList *envList) Configure() {
	eList.Port = defaultPort

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

	keycloakJwtUrl, ok := os.LookupEnv(KEYCLOAK_JWT_URL)
	if !ok {
		log.Fatal("KEYCLOAK_JWT_URL variable missing in environment variables")
	} else {
		eList.KEYCLOAK_JWT_URL = keycloakJwtUrl
	}
}
