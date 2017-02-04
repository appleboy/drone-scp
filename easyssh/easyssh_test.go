package easyssh

import (
	"os"
	"os/user"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetKeyFile(t *testing.T) {
	// missing file
	_, err := getKeyFile("abc")
	assert.Error(t, err)
	assert.Equal(t, "open abc: no such file or directory", err.Error())

	// wrong format
	_, err = getKeyFile("../tests/.ssh/id_rsa.pub")
	assert.Error(t, err)
	assert.Equal(t, "ssh: no key found", err.Error())

	_, err = getKeyFile("../tests/.ssh/id_rsa")
	assert.NoError(t, err)
}

func TestRunCommand(t *testing.T) {
	ssh := &MakeConfig{
		Server:  "localhost",
		User:    "drone-scp",
		Port:    "22",
		KeyPath: "../tests/.ssh/id_rsa",
	}

	output, err := ssh.Run("whoami")
	assert.Equal(t, "drone-scp\n", output)
	assert.NoError(t, err)
}

func TestSCPCommand(t *testing.T) {
	ssh := &MakeConfig{
		Server:  "localhost",
		User:    "drone-scp",
		Port:    "22",
		KeyPath: "../tests/.ssh/id_rsa",
	}

	err := ssh.Scp("../tests/a.txt")
	assert.NoError(t, err)

	u, err := user.Lookup("drone-scp")
	if err != nil {
		t.Fatalf("Lookup: %v", err)
	}

	// check file exist
	if _, err := os.Stat(u.HomeDir + "/a.txt"); os.IsNotExist(err) {
		t.Fatalf("SCP-error: %v", err)
	}
}
