package utils

import (
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	DBHOST                 = "DBHOST"
	DBUSER                 = "DBUSER"
	DBPORT                 = "DBPORT"
	DBPASS                 = "DBPASS"
	DBNAME                 = "DBNAME"
	KEYCLOAK_SERVER_URL    = "KEYCLOAK_SERVER_URL"
	KEYCLOAK_REALM_NAME    = "KEYCLOAK_REALM"
	KEYCLOAK_CLIENT_ID     = "KEYCLOAK_CLIENT_ID"
	KEYCLOAK_CLIENT_SECRET = "KEYCLOAK_CLIENT_SECRET"
	TESTING                = "TESTING"
)

const defaultPort = "8086"
const defaultDBPort = "5432"

type envList struct {
	Port                 string
	DBHost               string
	DBUser               string
	DBPort               string
	DBPass               string
	DBName               string
	KeycloakServerUrl    string
	KeycloakRealName     string
	KeycloakClientId     string
	KeycloakClientSecret string
	Testing              bool
}

var Config = &envList{}

func (eList *envList) Configure() {
	eList.Port = defaultPort

	db_host, ok := os.LookupEnv(DBHOST)
	if !ok {
		log.Fatalf("%s variable missing in environment variables", DBHOST)
	} else {
		eList.DBHost = db_host
	}

	db_user, ok := os.LookupEnv(DBUSER)
	if !ok {
		log.Fatalf("%s variable missing in environment variables", DBUSER)
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
		log.Fatalf("%s variable missing in environment variables", DBPASS)
	} else {
		eList.DBPass = db_pass
	}

	db_name, ok := os.LookupEnv(DBNAME)
	if !ok {
		log.Fatalf("%s variable missing in environment variables", DBNAME)
	} else {
		eList.DBName = db_name
	}

	KeycloakServerUrl, ok := os.LookupEnv(KEYCLOAK_SERVER_URL)
	if !ok {
		log.Fatalf("%s variable missing in environment variables", KEYCLOAK_SERVER_URL)
	} else {
		eList.KeycloakServerUrl = KeycloakServerUrl
	}

	KeycloakRealName, ok := os.LookupEnv(KEYCLOAK_REALM_NAME)
	if !ok {
		log.Fatalf("%s variable missing in environment variables", KEYCLOAK_REALM_NAME)
	} else {
		eList.KeycloakRealName = KeycloakRealName
	}

	KeycloakClientId, ok := os.LookupEnv(KEYCLOAK_CLIENT_ID)
	if !ok {
		log.Fatalf("%s variable missing in environment variables", KEYCLOAK_CLIENT_ID)
	} else {
		eList.KeycloakClientId = KeycloakClientId
	}

	KeycloakClientSecret, ok := os.LookupEnv(KEYCLOAK_CLIENT_SECRET)
	if !ok {
		log.Fatalf("%s variable missing in environment variables", KEYCLOAK_CLIENT_SECRET)
	} else {
		eList.KeycloakClientSecret = KeycloakClientSecret
	}

	IsTesting, ok := os.LookupEnv(TESTING)
	if !ok {
		eList.Testing = false
	} else {
		parsed := strings.ToLower(IsTesting)
		switch parsed {
		case "true":
			eList.Testing = true
		case "false":
			eList.Testing = false
		default:
			log.Fatalf("%s variable should be true/false", TESTING)
		}
	}
}
