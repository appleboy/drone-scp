package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/appleboy/drone-scp/easyssh"
)

type (
	// Repo information.
	Repo struct {
		Owner string
		Name  string
	}

	// Build information.
	Build struct {
		Event   string
		Number  int
		Commit  string
		Message string
		Branch  string
		Author  string
		Status  string
		Link    string
	}

	// Config for the plugin.
	Config struct {
		Host     string
		Port     string
		Username string
		Password string
		Key      string
		Target   string
		Source   []string
	}

	// Plugin values.
	Plugin struct {
		Repo   Repo
		Build  Build
		Config Config
	}
)

func trimPath(keys []string) []string {
	var newKeys []string

	for _, value := range keys {
		value = strings.Trim(value, " ")
		if len(value) == 0 {
			continue
		}

		newKeys = append(newKeys, value)
	}

	return newKeys
}

// Exec executes the plugin.
func (p Plugin) Exec() error {

	if len(p.Config.Host) == 0 || len(p.Config.Username) == 0 || (len(p.Config.Password) == 0 && len(p.Config.Key) == 0) {
		log.Println("missing ssh config")

		return errors.New("missing ssh config")
	}

	files := trimPath(p.Config.Source)
	src := strings.Join(files, " ")
	dest := fmt.Sprintf("%s-%s.tar", p.Repo.Name, p.Build.Commit[:7])

	// create a temporary file for the archive
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		return err
	}
	tar := filepath.Join(dir, dest)

	// run archive command
	log.Println("tar all files into " + tar)
	cmd := exec.Command("tar", "-cf", tar, src)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	// Create MakeConfig instance with remote username, server address and path to private key.
	ssh := &easyssh.MakeConfig{
		Server:   p.Config.Host,
		User:     p.Config.Username,
		Password: p.Config.Password,
		Port:     p.Config.Port,
		Key:      p.Config.Key,
	}

	// Call Scp method with file you want to upload to remote server.
	log.Println("scp file to remote server remote server.")
	err = ssh.Scp(tar)

	// Handle errors
	if err != nil {
		log.Println(err.Error())
		return err
	}

	// mkdir path
	log.Println("create remote folder " + p.Config.Target)
	_, err = ssh.Run(fmt.Sprintf("mkdir -p %s", p.Config.Target))

	if err != nil {
		log.Println(err.Error())
		return err
	}

	// untar file
	log.Println("untar remote file " + dest)
	_, err = ssh.Run(fmt.Sprintf("tar -xf %s -C %s", dest, p.Config.Target))

	if err != nil {
		log.Println(err.Error())
		return err
	}

	// remove tar file
	log.Println("remove remote file " + dest)
	_, err = ssh.Run(fmt.Sprintf("rm -rf %s", dest))

	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}
