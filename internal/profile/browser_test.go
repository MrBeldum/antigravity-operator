package profile_test

import (
	"os"
	"testing"

	"github.com/tiagoboas/antigravity-operator/internal/profile"
)

func TestPIDManagement(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "profile-test-*")
	if err != nil {
		t.Fatalf("falha ao criar temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Inicialmente não deve haver PID
	if pid := profile.ReadPID(tempDir); pid != 0 {
		t.Errorf("esperava PID 0, obteve: %d", pid)
	}

	// Gravar PID
	testPID := 12345
	if err := profile.SavePID(tempDir, testPID); err != nil {
		t.Fatalf("falha ao gravar PID: %v", err)
	}

	// Ler PID gravado
	read := profile.ReadPID(tempDir)
	if read != testPID {
		t.Errorf("esperava PID %d, obteve: %d", testPID, read)
	}
}

func TestCheckStatusInactivePort(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "profile-status-test-*")
	if err != nil {
		t.Fatalf("falha ao criar temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Porta 59999 provavelmente inativa
	st := profile.CheckStatusOnPort(tempDir, 59999)
	if st.IsRunning {
		t.Errorf("esperava porta 59999 inativa")
	}
}
