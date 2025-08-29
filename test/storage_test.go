package storage

import (
	"os"
	"testing"

	"github.com/Arup3201/gotasks/internal/storage"
)

func TestCreateStorage(t *testing.T) {
	_, err := storage.CreateStorage("../data/test.json")
	if err != nil {
		t.Errorf("storage.CreateStorage failed: %v", err)
	}

	_, err = os.Open("../data/test.json")
	if err != nil {
		t.Errorf("storage.CreateStorage did not create file")
	}
}

func TestRemoveStorage(t *testing.T) {
	err := storage.RemoveStorage("../data/test.json")
	if err != nil {
		t.Errorf("storage.RemoveStorage failed: %v", err)
	}

	_, err = os.Open("../data/test.json")
	if err==nil {
		t.Errorf("storage.RemoveStorage did not remove the storage file")
	}
}
