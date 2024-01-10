package hyperdrive

import (
	"os/user"
	"testing"
)

func TestIsRoot(t *testing.T) {
	// Test IsRoot with different UIDs and usernames
	u, err := user.Current()
	if u.Name == "root" {
		t.Errorf("Current User should not test with root: %v", err)
	}
	if IsRoot() {
		t.Log("Current user is root")
	} else {
		t.Log("Current user is not root")
	}
}

func TestCallingUserHomeDir(t *testing.T) {
	// Test CallingUserHomeDir with different home directories
	u, err := user.Current()
	if err != nil {
		t.Errorf("Failed to get current user: %v", err)
		return
	}

	homeDir, err := CallingUserHomeDir()
	if err != nil {
		t.Errorf("Failed to chown calling user home directory: %v", err)
		return
	}

	if homeDir != u.HomeDir {
		t.Errorf("Homedir != Current User Homedir: %v", homeDir)
	}
}
