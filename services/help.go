package services

import (
	"os"
	"io"
	"log"
	"net/http"
)

func getHelpInfo(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {//check method
		writer.WriteHeader(404)//service not found
		return
	}
	writer.Header().Add("Content-Type", "application/json")//content type
	file, err := os.Open("assets/help_groups.json")//open help info file
	if err==nil {
		defer file.Close()
		_, err = io.Copy(writer, file)//copy file into response
	}
	if err!=nil {
		log.Printf("Failed to return echo: %s", err.Error())
	}
}
