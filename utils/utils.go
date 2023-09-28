package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type Database struct {
	DBType       string `json:"db_type"`
	Host         string `json:"host"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	DatabaseName string `json:"database_name"`
	Port         string `json:"port"`
}

type ProjectConfig struct {
	ProjectName string    `json:"project_name"`
	DB          Database  `json:"db"`
	Port        string    `json:"port"`
	SecretKey   string    `json:"secrete_key"`
	Services    []Service `json:"services"`
}

type Service struct {
	Service string `json"service"`
	URL     string `json"url"`
}

type Route struct {
	Path    string
	Target  string
	Methods []string
}

func GetProjectConfig() ProjectConfig {
	var config ProjectConfig
	mydir, _ := os.Getwd()
	filePath := filepath.Join(mydir, "config", "project_config.json")
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Error opening JSON file: %v", err)
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file) // Read the entire file content
	if err != nil {
		log.Fatalf("Error reading JSON file: %v", err)
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}
	return config
}

func HandleError(ctx *gin.Context, statusCode int, errorMessage string) {
	ctx.JSON(statusCode, gin.H{
		"error": errorMessage,
	})
	ctx.Abort()
}

func Contains(slice []string, str string) bool {
	for _, item := range slice {
		if item == str {
			return true
		}
	}
	return false
}
