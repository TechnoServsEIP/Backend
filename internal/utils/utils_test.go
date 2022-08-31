package utils

import (
	"os/exec"
	"testing"
)

func TestCreateTmpFolder(t *testing.T) {
	folderName := "testUtils"

	if err := CreateTmpFolder(folderName); err != nil {
		t.Errorf("Error when creating folder")
	}

	cmd := exec.Command("ls", folderName)

	_, err := cmd.CombinedOutput()

	if err != nil {
		t.Errorf("Error when test to ls folder")
	}
}

func TestDeleteTmpFolder(t *testing.T) {
	folderName := "testUtils"

	if err := DeleteTmpFolder(folderName); err != nil {
		t.Errorf("Error when deleting folder")
	}

	cmd := exec.Command("ls", folderName)

	_, err := cmd.CombinedOutput()

	if err == nil {
		t.Errorf("Error when test to ls folder")
	}
}
