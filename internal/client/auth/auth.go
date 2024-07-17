package auth

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"
)

const TOKEN_KEY = "token"

func getConfigPath() (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		log.Print(err)
		return "", err
	}
	filePath := filepath.Join(dir, ".texttunnel", "config.json")

	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		log.Print(err)
		return "", err
	}

	return filePath, nil
}

func SaveToken(token string) error {
	filePath, err := getConfigPath()
	if err != nil {
		return err
	}

	configData := map[string]string{
		TOKEN_KEY: token,
	}

	data, err := json.Marshal(configData)
	if err != nil {
		log.Print(err)
		return err
	}

	err = os.WriteFile(filePath, data, 0600)
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}

func LoadToken() (string, error) {
	filePath, err := getConfigPath()
	if err != nil {
		return "", err
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Print(err)
		return "", err
	}

	configMap := make(map[string]string)
	err = json.Unmarshal(data, &configMap)
	if err != nil {
		log.Print(err)
		return "", err
	}

	token, ok := configMap[TOKEN_KEY]
	if !ok {
		log.Print("No token, Please login")
		return "", errors.New("No token found. Please login")
	}

	return token, nil
}
