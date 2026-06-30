package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"benkyou/model"
)

func ConvertTxtToJSON(txtFileDir string, writeDir string, jsonFileName string) error {
	fileContent, err := getFileContent(txtFileDir)
	if err != nil {
		return err
	}

	examples := getExamples(fileContent)

	err = writeToJSON(examples, writeDir, jsonFileName)
	if err != nil {
		return err
	}

	return nil
}

func writeToJSON(examples []model.Examples, writeDir string, fileName string) error {
	jsonFile, err := json.MarshalIndent(examples, "", "	")
	if err != nil {
		return err
	}

	writePath := fmt.Sprintf("%v/%v", writeDir, fileName)

	err = os.WriteFile(writePath, jsonFile, 0o644)
	if err != nil {
		return err
	}

	return nil
}

func getFileContent(fileDir string) ([]string, error) {
	rawContent, err := os.ReadFile(fileDir)
	if err != nil {
		return []string{}, err
	}

	content := strings.Split(string(rawContent), "\r\n")

	cleanContent := make([]string, 0)

	for _, v := range content {
		temp := strings.TrimSpace(v)

		if temp == "" {
			continue
		}

		cleanContent = append(cleanContent, v)

	}

	return cleanContent, nil
}

func getExamples(fileContent []string) []model.Examples {
	currentEntryIdx := 0

	examples := make([]model.Examples, 0)

	for i, v := range fileContent {
		switch {
		case i%3 == 0:
			examples = append(examples, model.Examples{
				Word: v,
			})
		case i%3 == 1:
			examples[currentEntryIdx].Reading = v
		case i%3 == 2:
			examples[currentEntryIdx].Meaning = v
			currentEntryIdx++
		}
	}

	return examples
}
