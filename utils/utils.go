package utils

import (
	"bytes"
	"crypto/aes"
	"encoding/hex"
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

// EncryptToHex encrypts plaintext using AES in ECB mode and returns the ciphertext as a hexadecimal string
func EncryptToHex(plaintext []byte, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	blockSize := block.BlockSize()
	plaintext = PKCS7Padding(plaintext, blockSize)

	ciphertext := make([]byte, len(plaintext))

	for i := 0; i < len(plaintext); i += blockSize {
		block.Encrypt(ciphertext[i:i+blockSize], plaintext[i:i+blockSize])
	}

	return hex.EncodeToString(ciphertext), nil
}

// PKCS7Padding pads the plaintext to a multiple of block size using PKCS#7 padding
func PKCS7Padding(plaintext []byte, blockSize int) []byte {
	padding := blockSize - (len(plaintext) % blockSize)
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(plaintext, padText...)
}
