package main

import (
	"log"
	"os"

	"benkyou/utils"
)

func main() {
	txtDir := os.Getenv("TXT_DIR")
	if txtDir == "" {
		log.Fatal("empty txt dir")
	}

	writeDir := os.Getenv("WRITE_DIR")
	if writeDir == "" {
		log.Fatal("empty write dir")
	}

	jsonFileName := os.Getenv("JSON_FILE_NAME")
	if jsonFileName == "" {
		log.Fatal("empty json filename")
	}

	err := utils.ConvertTxtToJSON(txtDir, writeDir, jsonFileName)
	if err != nil {
		log.Fatal(err)
	}
}
