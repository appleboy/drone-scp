package main

import (
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/appleboy/easyssh-proxy"
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

func TestSetPasswordAndKey(t *testing.T) {
	plugin := Plugin{
		Config: Config{
			Host:     []string{"example.com"},
			Username: "ubuntu",
			Port:     "443",
			Password: "1234",
			Key:      "test",
		},
	}

	err := plugin.Exec()

	assert.NotNil(t, err)
	assert.Equal(t, errSetPasswordandKey, err)
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
		if err := exec.Command("eval", "`ssh-agent -k`").Run(); err != nil {
			t.Fatalf("exec: %v", err)
		}
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
			CommandTimeout: 60 * time.Second,
			TarExec:        "tar",
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
		if err := exec.Command("eval", "`ssh-agent -k`").Run(); err != nil {
			t.Fatalf("exec: %v", err)
		}
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
			CommandTimeout: 60 * time.Second,
			TarExec:        "tar",
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
			CommandTimeout: 60 * time.Second,
			TarExec:        "tar",
			Proxy: easyssh.DefaultConfig{
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

func TestStripComponentsFlag(t *testing.T) {
	if os.Getenv("SSH_AUTH_SOCK") != "" {
		if err := exec.Command("eval", "`ssh-agent -k`").Run(); err != nil {
			t.Fatalf("exec: %v", err)
		}
	}

	u, err := user.Lookup("drone-scp")
	if err != nil {
		t.Fatalf("Lookup: %v", err)
	}

	plugin := Plugin{
		Config: Config{
			Host:            []string{"localhost"},
			Username:        "drone-scp",
			Port:            "22",
			KeyPath:         "tests/.ssh/id_rsa",
			Source:          []string{"tests/global/*"},
			StripComponents: 2,
			Target:          []string{filepath.Join(u.HomeDir, "123")},
			CommandTimeout:  60 * time.Second,
			TarExec:         "tar",
		},
	}

	err = plugin.Exec()
	assert.Nil(t, err)

	// check file exist
	if _, err := os.Stat(filepath.Join(u.HomeDir, "123/c.txt")); os.IsNotExist(err) {
		t.Fatalf("SCP-error: %v", err)
	}

	if _, err := os.Stat(filepath.Join(u.HomeDir, "123/d.txt")); os.IsNotExist(err) {
		t.Fatalf("SCP-error: %v", err)
	}
}

func TestIgnoreList(t *testing.T) {
	if os.Getenv("SSH_AUTH_SOCK") != "" {
		if err := exec.Command("eval", "`ssh-agent -k`").Run(); err != nil {
			t.Fatalf("exec: %v", err)
		}
	}

	u, err := user.Lookup("drone-scp")
	if err != nil {
		t.Fatalf("Lookup: %v", err)
	}

	plugin := Plugin{
		Config: Config{
			Host:            []string{"localhost"},
			Username:        "drone-scp",
			Port:            "22",
			KeyPath:         "tests/.ssh/id_rsa",
			Source:          []string{"tests/global/*", "!tests/global/c.txt", "!tests/global/e.txt"},
			StripComponents: 2,
			Target:          []string{filepath.Join(u.HomeDir, "ignore")},
			CommandTimeout:  60 * time.Second,
			TarExec:         "tar",
			Debug:           true,
		},
	}

	err = plugin.Exec()
	assert.Nil(t, err)

	// check file exist
	if _, err := os.Stat(filepath.Join(u.HomeDir, "ignore/c.txt")); err == nil {
		t.Fatal("c.txt file exist")
	}

	// check file exist
	if _, err := os.Stat(filepath.Join(u.HomeDir, "ignore/e.txt")); err == nil {
		t.Fatal("c.txt file exist")
	}

	if _, err := os.Stat(filepath.Join(u.HomeDir, "ignore/d.txt")); os.IsNotExist(err) {
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
			CommandTimeout: 60 * time.Second,
			TarExec:        "tar",
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
			CommandTimeout: 60 * time.Second,
			TarExec:        "tar",
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
	assert.Equal(t, expects, globList(paterns).Source)

	paterns = []string{"tests/*.txt", "tests/.ssh/*", "abc*"}
	expects = []string{"tests/a.txt", "tests/b.txt", "tests/.ssh/id_rsa", "tests/.ssh/id_rsa.pub"}
	assert.Equal(t, expects, globList(paterns).Source)

	paterns = []string{"tests/?.txt"}
	expects = []string{"tests/a.txt", "tests/b.txt"}
	assert.Equal(t, expects, globList(paterns).Source)

	// remove item which file not found.
	paterns = []string{"tests/aa.txt", "tests/b.txt"}
	expects = []string{"tests/b.txt"}
	assert.Equal(t, expects, globList(paterns).Source)

	paterns = []string{"./tests/b.txt"}
	expects = []string{"./tests/b.txt"}
	assert.Equal(t, expects, globList(paterns).Source)

	paterns = []string{"./tests/*.txt", "!./tests/b.txt"}
	expectSources := []string{"tests/a.txt", "tests/b.txt"}
	expectIgnores := []string{"./tests/b.txt"}
	result := globList(paterns)
	assert.Equal(t, expectSources, result.Source)
	assert.Equal(t, expectIgnores, result.Ignore)
}

func TestBuildArgs(t *testing.T) {
	list := fileList{
		Source: []string{"tests/a.txt", "tests/b.txt", "tests/c.txt"},
		Ignore: []string{"tests/a.txt", "tests/b.txt"},
	}

	result := buildArgs("test.tar.gz", list)
	expects := []string{"--exclude", "tests/a.txt", "--exclude", "tests/b.txt", "-cf", "test.tar.gz", "tests/a.txt", "tests/b.txt", "tests/c.txt"}
	assert.Equal(t, expects, result)

	list = fileList{
		Source: []string{"tests/a.txt", "tests/b.txt"},
	}

	result = buildArgs("test.tar.gz", list)
	expects = []string{"-cf", "test.tar.gz", "tests/a.txt", "tests/b.txt"}
	assert.Equal(t, expects, result)
}

func TestRemoveDestFile(t *testing.T) {
	ssh := &easyssh.MakeConfig{
		Server:  "localhost",
		User:    "drone-scp",
		Port:    "22",
		KeyPath: "tests/.ssh/id_rsa",
		// io timeout
		Timeout: 1,
	}
	plugin := Plugin{
		Config: Config{
			CommandTimeout: 60 * time.Second,
		},
		DestFile: "/etc/resolv.conf",
	}

	// ssh io timeout
	err := plugin.removeDestFile(ssh)
	assert.Error(t, err)

	ssh.Timeout = 0

	// permission denied
	err = plugin.removeDestFile(ssh)
	assert.Error(t, err)
}

func TestPlugin_buildArgs(t *testing.T) {
	type fields struct {
		Repo     Repo
		Build    Build
		Config   Config
		DestFile string
	}
	type args struct {
		target string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []string
	}{
		{
			name: "default command",
			fields: fields{
				Config: Config{
					Overwrite: false,
					TarExec:   "tar",
				},
				DestFile: "foo.tar",
			},
			args: args{
				target: "foo",
			},
			want: []string{"tar", "-xf", "foo.tar", "-C", "foo"},
		},
		{
			name: "strip components",
			fields: fields{
				Config: Config{
					Overwrite:       false,
					TarExec:         "tar",
					StripComponents: 2,
				},
				DestFile: "foo.tar",
			},
			args: args{
				target: "foo",
			},
			want: []string{"tar", "-xf", "foo.tar", "--strip-components", "2", "-C", "foo"},
		},
		{
			name: "overwrite",
			fields: fields{
				Config: Config{
					TarExec:         "tar",
					StripComponents: 2,
					Overwrite:       true,
				},
				DestFile: "foo.tar",
			},
			args: args{
				target: "foo",
			},
			want: []string{"tar", "-xf", "foo.tar", "--strip-components", "2", "--overwrite", "-C", "foo"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Plugin{
				Repo:     tt.fields.Repo,
				Build:    tt.fields.Build,
				Config:   tt.fields.Config,
				DestFile: tt.fields.DestFile,
			}
			if got := p.buildArgs(tt.args.target); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Plugin.buildArgs() = %v, want %v", got, tt.want)
			}
		})
	}
}
