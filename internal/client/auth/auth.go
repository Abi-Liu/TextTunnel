package auth

import (
	"encoding/json"
	"errors"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

const TOKEN_KEY = "token"

type FileSystem interface {
	UserHomeDir() (string, error)
	MkDirAll(path string, perm fs.FileMode) error
	WriteFile(filePath string, data []byte, perm fs.FileMode) error
	ReadFile(filePath string) ([]byte, error)
}

type OSFileSystem struct{}

func (*OSFileSystem) UserHomeDir() (string, error) {
	return os.UserHomeDir()
}

func (*OSFileSystem) MkDirAll(path string, perm fs.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (*OSFileSystem) WriteFile(filePath string, data []byte, perm fs.FileMode) error {
	return os.WriteFile(filePath, data, perm)
}

func (*OSFileSystem) ReadFile(filePath string) ([]byte, error) {
	return os.ReadFile(filePath)
}

type ConfigManager struct {
	FS FileSystem
}

func (c *ConfigManager) SaveToken(token string) error {
	homeDir, err := c.FS.UserHomeDir()
	if err != nil {
		log.Print(err)
		return err
	}

	filePath := filepath.Join(homeDir, ".texttunnel", "config.json")
	if err := c.FS.MkDirAll(filepath.Dir(filePath), 0755); err != nil {
		log.Print(err)
		return err
	}

	configMap := map[string]string{
		TOKEN_KEY: token,
	}

	data, err := json.Marshal(configMap)
	if err != nil {
		log.Print(err)
		return err
	}

	err = c.FS.WriteFile(filePath, data, 0600)
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}

func (c *ConfigManager) LoadToken() (string, error) {
	homeDir, err := c.FS.UserHomeDir()
	if err != nil {
		return "", err
	}

	filePath := filepath.Join(homeDir, ".texttunnel", "config.json")

	data, err := c.FS.ReadFile(filePath)
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
