package main

import "os"

func GetGCPProjectID() string {
	return os.Getenv("GCP_PROJECT_ID")
}
