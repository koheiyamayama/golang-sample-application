package main

import "os"

func GetGCPProjectID() string {
	return os.Getenv("GCP_PROJECT_ID")
}

func GetDBUserName() string {
	defaultName := "root"
	if userName := os.Getenv("DATABASE_USER_NAME"); userName != "" {
		return userName
	} else {
		return defaultName
	}
}

func GetDBPassword() string {
	defaultPassword := "root"
	if password := os.Getenv("DATABASE_PASSWORD"); password != "" {
		return password
	} else {
		return defaultPassword
	}
}

func GetDBHostName() string {
	defaultHostname := "rds"
	if hostname := os.Getenv("DATABASE_HOSTNAME"); hostname != "" {
		return hostname
	} else {
		return defaultHostname
	}
}

func GetDBPort() string {
	defaultPort := "3306"
	if port := os.Getenv("DATABASE_PORT"); port != "" {
		return port
	} else {
		return defaultPort
	}
}

func GetDatabaseName() string {
	defaultDatabaseName := "google-cloud-go"
	if databaseName := os.Getenv("DATABASE_NAME"); databaseName != "" {
		return databaseName
	} else {
		return defaultDatabaseName
	}
}
