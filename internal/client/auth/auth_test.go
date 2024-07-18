package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"testing"
)

type MockFileSystem struct {
	files map[string][]byte
	home  string
}

const HOME_DIR = "test/home"

func NewMockFileSystem() *MockFileSystem {
	return &MockFileSystem{files: make(map[string][]byte), home: HOME_DIR}
}

func (m *MockFileSystem) ReadFile(filename string) ([]byte, error) {
	data, exists := m.files[filename]
	if !exists {
		return nil, errors.New("file not found")
	}
	return data, nil
}

func (m *MockFileSystem) WriteFile(filename string, data []byte, perm os.FileMode) error {
	m.files[filename] = data
	return nil
}

func (m *MockFileSystem) UserHomeDir() (string, error) {
	return m.home, nil
}

func (m *MockFileSystem) MkDirAll(path string, perm os.FileMode) error {
	return nil
}

func TestSaveToken(t *testing.T) {
	fs := NewMockFileSystem()
	cm := ConfigManager{FS: fs}
	err := cm.SaveToken("Test")

	if err != nil {
		t.Fatalf("Expected nil error. Received %s", err)
	}

	expectedPath := fmt.Sprintf("%s/.texttunnel/config.json", HOME_DIR)
	data, err := fs.ReadFile(expectedPath)
	if err != nil {
		t.Fatalf("Encountered error when reading from file: %s", err)
	}

	file := make(map[string]string)
	err = json.Unmarshal(data, &file)

	if err != nil {
		t.Fatalf("Failed to unmarshal json")
	}

	token, _ := file[TOKEN_KEY]

	if token != "Test" {
		t.Fatalf("Expected Test, Received: %s", string(data))
	}
}
