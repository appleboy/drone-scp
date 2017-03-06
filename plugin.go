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
	"sync"
	"time"

	"github.com/appleboy/com/random"
	"github.com/appleboy/easyssh-proxy"
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
		Host           []string
		Port           string
		Username       string
		Password       string
		Key            string
		KeyPath        string
		Timeout        time.Duration
		CommandTimeout int
		Target         []string
		Source         []string
		Remove         bool
	}

	// Plugin values.
	Plugin struct {
		Repo   Repo
		Build  Build
		Config Config
	}
)

var wg sync.WaitGroup

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

func globList(paths []string) []string {
	var newPaths []string

	for _, pattern := range paths {
		pattern = strings.Trim(pattern, " ")
		matches, err := filepath.Glob(pattern)
		if err != nil {
			log.Printf("Glob error for %q: %s\n", pattern, err)
			continue
		}

		newPaths = append(newPaths, matches...)
	}

	return newPaths
}

func (p Plugin) log(host string, message ...interface{}) {
	log.Printf("%s: %s", host, fmt.Sprintln(message...))
}

// Exec executes the plugin.
func (p Plugin) Exec() error {

	if len(p.Config.Host) == 0 || len(p.Config.Username) == 0 {
		return errors.New("missing ssh config (Host, Username)")
	}

	if len(p.Config.Source) == 0 || len(p.Config.Target) == 0 {
		return errors.New("missing source or target config")
	}

	files := globList(trimPath(p.Config.Source))
	dest := fmt.Sprintf("%s.tar", random.String(10))

	// create a temporary file for the archive
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		return err
	}
	tar := filepath.Join(dir, dest)

	// run archive command
	log.Println("tar all files into " + tar)
	args := append(append([]string{}, "-cf", getRealPath(tar)), files...)

	cmd := exec.Command("tar", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	wg.Add(len(p.Config.Host))
	errChannel := make(chan error, 1)
	finished := make(chan bool, 1)
	for _, host := range p.Config.Host {
		go func(host string) {
			// Create MakeConfig instance with remote username, server address and path to private key.
			ssh := &easyssh.MakeConfig{
				Server:   host,
				User:     p.Config.Username,
				Password: p.Config.Password,
				Port:     p.Config.Port,
				Key:      p.Config.Key,
				KeyPath:  p.Config.KeyPath,
				Timeout:  p.Config.Timeout,
			}

			// Call Scp method with file you want to upload to remote server.
			p.log(host, "scp file to server.")
			err := ssh.Scp(tar, dest)

			if err != nil {
				errChannel <- err
			}

			for _, target := range p.Config.Target {
				// remove target before upload data
				if p.Config.Remove {
					p.log(host, "Remove target folder:", target)

					_, _, _, err := ssh.Run(fmt.Sprintf("rm -rf %s", target), p.Config.CommandTimeout)

					if err != nil {
						errChannel <- err
					}
				}

				// mkdir path
				p.log(host, "create folder", target)
				_, errStr, _, err := ssh.Run(fmt.Sprintf("mkdir -p %s", target), p.Config.CommandTimeout)
				if err != nil {
					errChannel <- err
				}

				if len(errStr) != 0 {
					errChannel <- fmt.Errorf(errStr)
				}

				// untar file
				p.log(host, "untar file", dest)
				_, _, _, err = ssh.Run(fmt.Sprintf("tar -xf %s -C %s", dest, target), p.Config.CommandTimeout)

				if err != nil {
					errChannel <- err
				}
			}

			// remove tar file
			p.log(host, "remove file", dest)
			_, _, _, err = ssh.Run(fmt.Sprintf("rm -rf %s", dest), p.Config.CommandTimeout)

			if err != nil {
				errChannel <- err
			}

			wg.Done()

		}(host)
	}

	go func() {
		wg.Wait()
		close(finished)
	}()

	select {
	case <-finished:
	case err := <-errChannel:
		if err != nil {
			fmt.Println("drone-scp error: ", err)
			return err
		}
	}

	fmt.Println("Successfully executed transfer data to all host.")

	return nil
}
