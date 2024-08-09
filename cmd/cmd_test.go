package cmd_test

// import (
// 	"bytes"
// 	"testing"

// 	"github.com/DnFreddie/backy/cmd"
// )

// func TestBackupCmd(t *testing.T) {
// 	var outputBuffer bytes.Buffer

// 	cmd.BackupCmd.Root().SetOut(&outputBuffer)

// 	cmd.BackupCmd.Root().SetArgs([]string{"backup.go", "trip.go"})

// 	err := cmd.BackupCmd.Execute()
// 	if err != nil {
// 		t.Fatalf("Error executing command: %v", err)
// 	}

// 	expectedOutput := "expected output"
// 	actualOutput := outputBuffer.String()

// 	if actualOutput != expectedOutput {
// 		t.Errorf("Expected output %q, but got %q", expectedOutput, actualOutput)
// 	}

// }
