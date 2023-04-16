package main

import (
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/appleboy/easyssh-proxy"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/ssh"
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

	assert.Equal(t, result, trimValues(input))

	input = []string{"1", "2"}
	result = []string{"1", "2"}

	assert.Equal(t, result, trimValues(input))
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

func TestSCPFileFromPublicKeyWithPassphrase(t *testing.T) {
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
			KeyPath:        "tests/.ssh/test",
			Passphrase:     "1234",
			Source:         []string{"tests/a.txt", "tests/b.txt"},
			Target:         []string{filepath.Join(u.HomeDir, "/test2")},
			CommandTimeout: 60 * time.Second,
			TarExec:        "tar",
		},
	}

	err = plugin.Exec()
	assert.Nil(t, err)

	// check file exist
	if _, err := os.Stat(filepath.Join(u.HomeDir, "/test2/tests/a.txt")); os.IsNotExist(err) {
		t.Fatalf("SCP-error: %v", err)
	}

	if _, err := os.Stat(filepath.Join(u.HomeDir, "/test2/tests/b.txt")); os.IsNotExist(err) {
		t.Fatalf("SCP-error: %v", err)
	}
}

func TestWrongFingerprint(t *testing.T) {
	u, err := user.Lookup("drone-scp")
	if err != nil {
		t.Fatalf("Lookup: %v", err)
	}

	plugin := Plugin{
		Config: Config{
			Host:           []string{"localhost"},
			Username:       "drone-scp",
			Port:           "22",
			KeyPath:        "./tests/.ssh/id_rsa",
			Source:         []string{"tests/a.txt", "tests/b.txt"},
			Target:         []string{filepath.Join(u.HomeDir, "/test2")},
			CommandTimeout: 60 * time.Second,
			TarExec:        "tar",
			Fingerprint:    "wrong",
		},
	}

	err = plugin.Exec()
	log.Println(err)
	assert.NotNil(t, err)
}

func getHostPublicKeyFile(keypath string) (ssh.PublicKey, error) {
	var pubkey ssh.PublicKey
	var err error
	buf, err := os.ReadFile(keypath)
	if err != nil {
		return nil, err
	}

	pubkey, _, _, _, err = ssh.ParseAuthorizedKey(buf)

	if err != nil {
		return nil, err
	}

	return pubkey, nil
}

func TestSCPFileFromPublicKeyWithFingerprint(t *testing.T) {
	if os.Getenv("SSH_AUTH_SOCK") != "" {
		if err := exec.Command("eval", "`ssh-agent -k`").Run(); err != nil {
			t.Fatalf("exec: %v", err)
		}
	}

	u, err := user.Lookup("drone-scp")
	if err != nil {
		t.Fatalf("Lookup: %v", err)
	}

	hostKey, err := getHostPublicKeyFile("/etc/ssh/ssh_host_rsa_key.pub")
	assert.NoError(t, err)

	plugin := Plugin{
		Config: Config{
			Host:           []string{"localhost"},
			Username:       "drone-scp",
			Port:           "22",
			KeyPath:        "./tests/.ssh/id_rsa",
			Fingerprint:    ssh.FingerprintSHA256(hostKey),
			Source:         []string{"tests/a.txt", "tests/b.txt"},
			Target:         []string{filepath.Join(u.HomeDir, "/test2")},
			CommandTimeout: 60 * time.Second,
			TarExec:        "tar",
		},
	}

	err = plugin.Exec()
	assert.Nil(t, err)

	// check file exist
	if _, err := os.Stat(filepath.Join(u.HomeDir, "/test2/tests/a.txt")); os.IsNotExist(err) {
		t.Fatalf("SCP-error: %v", err)
	}

	if _, err := os.Stat(filepath.Join(u.HomeDir, "/test2/tests/b.txt")); os.IsNotExist(err) {
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

func TestUseInsecureCipherFlag(t *testing.T) {
	u, err := user.Lookup("drone-scp")
	if err != nil {
		t.Fatalf("Lookup: %v", err)
	}

	plugin := Plugin{
		Config: Config{
			Host:              []string{"localhost"},
			Username:          "drone-scp",
			Port:              "22",
			KeyPath:           "tests/.ssh/id_rsa",
			Source:            []string{"tests/global/*"},
			StripComponents:   2,
			Target:            []string{filepath.Join(u.HomeDir, "123")},
			CommandTimeout:    60 * time.Second,
			TarExec:           "tar",
			UseInsecureCipher: true,
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
	if _, err := os.Stat(filepath.Join(u.HomeDir, "ignore/c.txt")); !os.IsNotExist(err) {
		t.Fatal("c.txt file exist")
	}

	// check file exist
	if _, err := os.Stat(filepath.Join(u.HomeDir, "ignore/e.txt")); !os.IsNotExist(err) {
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
	expects = []string{"tests/a.txt", "tests/b.txt", "tests/.ssh/id_rsa", "tests/.ssh/id_rsa.pub", "tests/.ssh/test", "tests/.ssh/test.pub"}
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

	_, _, _, err := ssh.Run("ver", plugin.Config.CommandTimeout)
	systemType := "unix"
	if err == nil {
		systemType = "windows"
	}

	// ssh io timeout
	err = plugin.removeDestFile(systemType, ssh)
	assert.Error(t, err)

	ssh.Timeout = 0

	// permission denied
	err = plugin.removeDestFile(systemType, ssh)
	assert.Error(t, err)
}

func TestPlugin_buildUnTarArgs(t *testing.T) {
	type fields struct {
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
					Overwrite:   false,
					UnlinkFirst: false,
					TarExec:     "tar",
				},
				DestFile: "foo.tar.gz",
			},
			args: args{
				target: "foo",
			},
			want: []string{"tar", "-zxf", "foo.tar.gz", "-C", "foo"},
		},
		{
			name: "strip components",
			fields: fields{
				Config: Config{
					Overwrite:       false,
					UnlinkFirst:     false,
					TarExec:         "tar",
					StripComponents: 2,
				},
				DestFile: "foo.tar.gz",
			},
			args: args{
				target: "foo",
			},
			want: []string{"tar", "-zxf", "foo.tar.gz", "--strip-components", "2", "-C", "foo"},
		},
		{
			name: "overwrite",
			fields: fields{
				Config: Config{
					TarExec:         "tar",
					StripComponents: 2,
					Overwrite:       true,
					UnlinkFirst:     false,
				},
				DestFile: "foo.tar.gz",
			},
			args: args{
				target: "foo",
			},
			want: []string{"tar", "-zxf", "foo.tar.gz", "--strip-components", "2", "--overwrite", "-C", "foo"},
		},
		{
			name: "unlink first",
			fields: fields{
				Config: Config{
					TarExec:         "tar",
					StripComponents: 2,
					Overwrite:       true,
					UnlinkFirst:     true,
				},
				DestFile: "foo.tar.gz",
			},
			args: args{
				target: "foo",
			},
			want: []string{"tar", "-zxf", "foo.tar.gz", "--strip-components", "2", "--overwrite", "--unlink-first", "-C", "foo"},
		},
		{
			name: "output folder path with space",
			fields: fields{
				Config: Config{
					TarExec:         "tar",
					StripComponents: 2,
					Overwrite:       true,
					UnlinkFirst:     true,
				},
				DestFile: "foo.tar.gz",
			},
			args: args{
				target: "foo\\ bar",
			},
			want: []string{"tar", "-zxf", "foo.tar.gz", "--strip-components", "2", "--overwrite", "--unlink-first", "-C", "foo\\ bar"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Plugin{
				Config:   tt.fields.Config,
				DestFile: tt.fields.DestFile,
			}
			if got := p.buildUnTarArgs(tt.args.target); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Plugin.buildArgs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlugin_buildTarArgs(t *testing.T) {
	type fields struct {
		Config Config
	}
	type args struct {
		src string
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
					TarExec: "tar",
				},
			},
			args: args{
				src: "foo.tar.gz",
			},
			want: []string{"-zcf", "foo.tar.gz"},
		},
		{
			name: "ignore list",
			fields: fields{
				Config: Config{
					TarExec: "tar",
					Source: []string{
						"tests/*.txt",
						"!tests/a.txt",
					},
				},
			},
			args: args{
				src: "foo.tar.gz",
			},
			want: []string{"--exclude", "tests/a.txt", "-zcf", "foo.tar.gz", "tests/a.txt", "tests/b.txt"},
		},
		{
			name: "dereference flag",
			fields: fields{
				Config: Config{
					TarExec:        "tar",
					TarDereference: true,
					Source: []string{
						"tests/*.txt",
						"!tests/a.txt",
					},
				},
			},
			args: args{
				src: "foo.tar.gz",
			},
			want: []string{"--exclude", "tests/a.txt", "--dereference", "-zcf", "foo.tar.gz", "tests/a.txt", "tests/b.txt"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Plugin{
				Config: tt.fields.Config,
			}
			if got := p.buildTarArgs(tt.args.src); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Plugin.buildTarArgs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTargetFolderWithSpaces(t *testing.T) {
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
			Target:          []string{filepath.Join(u.HomeDir, "123 456 789")},
			CommandTimeout:  60 * time.Second,
			TarExec:         "tar",
		},
	}

	err = plugin.Exec()
	assert.Nil(t, err)

	// check file exist
	if _, err := os.Stat(filepath.Join(u.HomeDir, "123 456 789", "c.txt")); os.IsNotExist(err) {
		t.Fatalf("SCP-error: %v", err)
	}

	if _, err := os.Stat(filepath.Join(u.HomeDir, "123 456 789", "d.txt")); os.IsNotExist(err) {
		t.Fatalf("SCP-error: %v", err)
	}
}

func TestHostPortString(t *testing.T) {
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
			Host:            []string{"localhost:22", "localhost:22"},
			Username:        "drone-scp",
			Port:            "8080",
			KeyPath:         "tests/.ssh/id_rsa",
			Source:          []string{"tests/global/*"},
			StripComponents: 2,
			Target:          []string{filepath.Join(u.HomeDir, "1234")},
			CommandTimeout:  60 * time.Second,
			TarExec:         "tar",
		},
	}

	err = plugin.Exec()
	assert.Nil(t, err)

	// check file exist
	if _, err := os.Stat(filepath.Join(u.HomeDir, "1234", "c.txt")); os.IsNotExist(err) {
		t.Fatalf("SCP-error: %v", err)
	}

	if _, err := os.Stat(filepath.Join(u.HomeDir, "1234", "d.txt")); os.IsNotExist(err) {
		t.Fatalf("SCP-error: %v", err)
	}
}

// Unit test for hostPort
func TestHostPort(t *testing.T) {
	p := Plugin{
		Config: Config{
			Port: "8080",
		},
	}

	// Test case 1: host string with port
	host1 := "example.com:1234"
	expectedHost1 := "example.com"
	expectedPort1 := "1234"
	actualHost1, actualPort1 := p.hostPort(host1)
	if actualHost1 != expectedHost1 || actualPort1 != expectedPort1 {
		t.Errorf("hostPort(%s) = (%s, %s); expected (%s, %s)", host1, actualHost1, actualPort1, expectedHost1, expectedPort1)
	}

	// Test case 2: host string without port
	host2 := "example.com"
	expectedHost2 := "example.com"
	expectedPort2 := "8080" // default port
	actualHost2, actualPort2 := p.hostPort(host2)
	if actualHost2 != expectedHost2 || actualPort2 != expectedPort2 {
		t.Errorf("hostPort(%s) = (%s, %s); expected (%s, %s)", host2, actualHost2, actualPort2, expectedHost2, expectedPort2)
	}
}
