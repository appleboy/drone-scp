package main

import (
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMissingAllConfig(t *testing.T) {
	var plugin Plugin

	err := plugin.Exec()

	assert.NotNil(t, err)
}

func TestMissingSSHConfig(t *testing.T) {
	plugin := Plugin{
		Config: Config{
			Host:     []string{"example.com"},
			Username: "ubuntu",
		},
	}

	err := plugin.Exec()

	assert.NotNil(t, err)
}

func TestMissingSourceConfig(t *testing.T) {
	plugin := Plugin{
		Config: Config{
			Host:     []string{"example.com"},
			Username: "ubuntu",
			Port:     "443",
			Password: "1234",
		},
	}

	err := plugin.Exec()

	assert.NotNil(t, err)
}

func TestTrimElement(t *testing.T) {
	var input, result []string

	input = []string{"1", "     ", "3"}
	result = []string{"1", "3"}

	assert.Equal(t, result, trimPath(input))

	input = []string{"1", "2"}
	result = []string{"1", "2"}

	assert.Equal(t, result, trimPath(input))
}

func TestSCPFileFromPublicKey(t *testing.T) {
	if os.Getenv("SSH_AUTH_SOCK") != "" {
		exec.Command("eval", "`ssh-agent -k`").Run()
	}

	u, err := user.Lookup("drone-scp")
	if err != nil {
		t.Fatalf("Lookup: %v", err)
	}

	plugin := Plugin{
		Config: Config{
			Host:           []string{"localhost"},
			Username:       "drone-scp",
			Port:           "22",
			KeyPath:        "tests/.ssh/id_rsa",
			Source:         []string{"tests/a.txt", "tests/b.txt"},
			Target:         []string{filepath.Join(u.HomeDir, "/test")},
			CommandTimeout: 60,
		},
	}

	err = plugin.Exec()
	assert.Nil(t, err)

	// check file exist
	if _, err := os.Stat(filepath.Join(u.HomeDir, "/test/tests/a.txt")); os.IsNotExist(err) {
		t.Fatalf("SCP-error: %v", err)
	}

	if _, err := os.Stat(filepath.Join(u.HomeDir, "/test/tests/b.txt")); os.IsNotExist(err) {
		t.Fatalf("SCP-error: %v", err)
	}

	// Test -rm flag
	plugin.Config.Source = []string{"tests/a.txt"}
	plugin.Config.Remove = true

	err = plugin.Exec()
	assert.Nil(t, err)

	// check file exist
	if _, err := os.Stat(filepath.Join(u.HomeDir, "/test/tests/b.txt")); os.IsExist(err) {
		t.Fatalf("SCP-error: %v", err)
	}
}

func TestSCPWildcardFileList(t *testing.T) {
	if os.Getenv("SSH_AUTH_SOCK") != "" {
		exec.Command("eval", "`ssh-agent -k`").Run()
	}

	u, err := user.Lookup("drone-scp")
	if err != nil {
		t.Fatalf("Lookup: %v", err)
	}

	plugin := Plugin{
		Config: Config{
			Host:           []string{"localhost"},
			Username:       "drone-scp",
			Port:           "22",
			KeyPath:        "tests/.ssh/id_rsa",
			Source:         []string{"tests/global/*"},
			Target:         []string{filepath.Join(u.HomeDir, "abc")},
			CommandTimeout: 60,
		},
	}

	err = plugin.Exec()
	assert.Nil(t, err)

	// check file exist
	if _, err := os.Stat(filepath.Join(u.HomeDir, "abc/tests/global/c.txt")); os.IsNotExist(err) {
		t.Fatalf("SCP-error: %v", err)
	}

	if _, err := os.Stat(filepath.Join(u.HomeDir, "abc/tests/global/d.txt")); os.IsNotExist(err) {
		t.Fatalf("SCP-error: %v", err)
	}
}

func TestSCPFromProxySetting(t *testing.T) {
	u, err := user.Lookup("drone-scp")
	if err != nil {
		t.Fatalf("Lookup: %v", err)
	}

	plugin := Plugin{
		Config: Config{
			Host:           []string{"localhost"},
			Username:       "drone-scp",
			Port:           "22",
			KeyPath:        "tests/.ssh/id_rsa",
			Source:         []string{"tests/global/*"},
			Target:         []string{filepath.Join(u.HomeDir, "def")},
			CommandTimeout: 60,
			Proxy: defaultConfig{
				Server:  "localhost",
				User:    "drone-scp",
				Port:    "22",
				KeyPath: "./tests/.ssh/id_rsa",
			},
		},
	}

	err = plugin.Exec()
	assert.Nil(t, err)

	// check file exist
	if _, err := os.Stat(filepath.Join(u.HomeDir, "def/tests/global/c.txt")); os.IsNotExist(err) {
		t.Fatalf("SCP-error: %v", err)
	}

	if _, err := os.Stat(filepath.Join(u.HomeDir, "def/tests/global/d.txt")); os.IsNotExist(err) {
		t.Fatalf("SCP-error: %v", err)
	}
}

// func TestSCPFileFromSSHAgent(t *testing.T) {
// 	if os.Getenv("SSH_AUTH_SOCK") == "" {
// 		exec.Command("eval", "`ssh-agent -s`").Run()
// 		exec.Command("ssh-add", "tests/.ssh/id_rsa").Run()
// 	}

// 	u, err := user.Lookup("drone-scp")
// 	if err != nil {
// 		t.Fatalf("Lookup: %v", err)
// 	}

// 	plugin := Plugin{
// 		Config: Config{
// 			Host:     []string{"localhost"},
// 			Username: "drone-scp",
// 			Port:     "22",
// 			Source:   []string{"tests/a.txt", "tests/b.txt"},
// 			Target:   []string{u.HomeDir + "/test"},
// 		},
// 	}

// 	err = plugin.Exec()
// 	assert.Nil(t, err)
// }

// func TestSCPFileFromPassword(t *testing.T) {
// 	u, err := user.Lookup("drone-scp")
// 	if err != nil {
// 		t.Fatalf("Lookup: %v", err)
// 	}

// 	plugin := Plugin{
// 		Config: Config{
// 			Host:     []string{"localhost"},
// 			Username: "drone-scp",
// 			Port:     "22",
// 			Password: "1234",
// 			Source:   []string{"tests/a.txt", "tests/b.txt"},
// 			Target:   []string{u.HomeDir + "/test"},
// 		},
// 	}

// 	err = plugin.Exec()
// 	assert.Nil(t, err)
// }

func TestIncorrectPassword(t *testing.T) {
	plugin := Plugin{
		Config: Config{
			Host:           []string{"localhost"},
			Username:       "drone-scp",
			Port:           "22",
			Password:       "123456",
			Source:         []string{"tests/a.txt", "tests/b.txt"},
			Target:         []string{"/home"},
			CommandTimeout: 60,
		},
	}

	err := plugin.Exec()
	assert.NotNil(t, err)
}

func TestNoPermissionCreateFolder(t *testing.T) {
	u, err := user.Lookup("drone-scp")
	if err != nil {
		t.Fatalf("Lookup: %v", err)
	}

	plugin := Plugin{
		Config: Config{
			Host:           []string{"localhost"},
			Username:       "drone-scp",
			Port:           "22",
			KeyPath:        "tests/.ssh/id_rsa",
			Source:         []string{"tests/a.txt", "tests/b.txt"},
			Target:         []string{"/etc/test"},
			CommandTimeout: 60,
		},
	}

	err = plugin.Exec()
	assert.NotNil(t, err)

	// check tmp file exist
	if _, err = os.Stat(filepath.Join(u.HomeDir, plugin.DestFile)); os.IsExist(err) {
		t.Fatalf("SCP-error: %v", err)
	}
}

func TestGlobList(t *testing.T) {
	// wrong patern
	paterns := []string{"[]a]", "tests/?.txt"}
	expects := []string{"tests/a.txt", "tests/b.txt"}
	assert.Equal(t, expects, globList(paterns))

	paterns = []string{"tests/*.txt", "tests/.ssh/*", "abc*"}
	expects = []string{"tests/a.txt", "tests/b.txt", "tests/.ssh/id_rsa", "tests/.ssh/id_rsa.pub"}
	assert.Equal(t, expects, globList(paterns))

	paterns = []string{"tests/?.txt"}
	expects = []string{"tests/a.txt", "tests/b.txt"}
	assert.Equal(t, expects, globList(paterns))

	// remove item which file not found.
	paterns = []string{"tests/aa.txt", "tests/b.txt"}
	expects = []string{"tests/b.txt"}
	assert.Equal(t, expects, globList(paterns))
}
