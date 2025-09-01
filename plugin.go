package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/appleboy/com/random"
	"github.com/appleboy/easyssh-proxy"
	"github.com/fatih/color"
)

var (
	errMissingHost           = errors.New("Error: missing server host")
	errMissingPasswordOrKey  = errors.New("Error: can't connect without a private SSH key or password")
	errMissingSourceOrTarget = errors.New("missing source or target config")
)

type (
	// Config for the plugin.
	Config struct {
		Host              []string
		Port              int
		Protocol          easyssh.Protocol
		Username          string
		Password          string
		Key               string
		Passphrase        string
		Fingerprint       string
		KeyPath           string
		Timeout           time.Duration
		CommandTimeout    time.Duration
		Target            []string
		Source            []string
		Remove            bool
		StripComponents   int
		TarExec           string
		TarTmpPath        string
		Proxy             easyssh.DefaultConfig
		Debug             bool
		Overwrite         bool
		UnlinkFirst       bool
		Ciphers           []string
		UseInsecureCipher bool
		TarDereference    bool
	}

	// Plugin values.
	Plugin struct {
		Config   Config
		DestFile string
	}

	copyError struct {
		host    string
		message string
	}
)

func (e copyError) Error() string {
	return fmt.Sprintf("error copy file to dest: %s, error message: %s\n", e.host, e.message)
}

func globList(paths []string) fileList {
	var list fileList

	for _, pattern := range paths {
		ignore := false
		pattern = strings.TrimSpace(pattern)
		if string(pattern[0]) == "!" {
			pattern = pattern[1:]
			ignore = true
		}
		matches, err := filepath.Glob(pattern)
		if err != nil {
			fmt.Printf("Glob error for %q: %s\n", pattern, err)
			continue
		}

		if ignore {
			list.Ignore = append(list.Ignore, matches...)
		} else {
			list.Source = append(list.Source, matches...)
		}
	}

	return list
}

func (p Plugin) log(host string, message ...interface{}) {
	if count := len(p.Config.Host); count == 1 {
		fmt.Printf("%s", fmt.Sprintln(message...))
	} else {
		fmt.Printf("%s: %s", host, fmt.Sprintln(message...))
	}
}

func (p *Plugin) removeDestFile(os string, ssh *easyssh.MakeConfig) error {
	p.log(ssh.Server, "remove file", p.DestFile)
	_, errStr, _, err := ssh.Run(rmcmd(os, p.DestFile), p.Config.CommandTimeout)
	if err != nil {
		return err
	}

	if errStr != "" {
		return errors.New(errStr)
	}

	return nil
}

func (p *Plugin) removeAllDestFile() error {
	for _, h := range trimValues(p.Config.Host) {
		host, port := p.hostPort(h)
		ssh := &easyssh.MakeConfig{
			Server:            host,
			User:              p.Config.Username,
			Password:          p.Config.Password,
			Port:              port,
			Protocol:          p.Config.Protocol,
			Key:               p.Config.Key,
			KeyPath:           p.Config.KeyPath,
			Passphrase:        p.Config.Passphrase,
			Timeout:           p.Config.Timeout,
			Ciphers:           p.Config.Ciphers,
			Fingerprint:       p.Config.Fingerprint,
			UseInsecureCipher: p.Config.UseInsecureCipher,
			Proxy: easyssh.DefaultConfig{
				Server:            p.Config.Proxy.Server,
				User:              p.Config.Proxy.User,
				Password:          p.Config.Proxy.Password,
				Port:              p.Config.Proxy.Port,
				Protocol:          p.Config.Proxy.Protocol,
				Key:               p.Config.Proxy.Key,
				KeyPath:           p.Config.Proxy.KeyPath,
				Passphrase:        p.Config.Proxy.Passphrase,
				Timeout:           p.Config.Proxy.Timeout,
				Ciphers:           p.Config.Proxy.Ciphers,
				Fingerprint:       p.Config.Proxy.Fingerprint,
				UseInsecureCipher: p.Config.Proxy.UseInsecureCipher,
			},
		}

		_, _, _, err := ssh.Run("ver", p.Config.CommandTimeout)
		systemType := "unix"
		if err == nil {
			systemType = "windows"
		}

		// remove tar file
		err = p.removeDestFile(systemType, ssh)
		if err != nil {
			return err
		}
	}

	return nil
}

type fileList struct {
	Ignore []string
	Source []string
}

func (p *Plugin) buildTarArgs(src string) []string {
	files := globList(trimValues(p.Config.Source))
	args := []string{}
	if len(files.Ignore) > 0 {
		for _, v := range files.Ignore {
			args = append(args, "--exclude")
			args = append(args, v)
		}
	}

	if p.Config.TarDereference {
		args = append(args, "--dereference")
	}

	args = append(args, "-zcf")
	args = append(args, getRealPath(src))

	// For precise operation, adding an additional on/off option needed.
	// e.g. SCP_ACTION_WILDCARD_COMPATIBLE
	hasCommonFolder := true
	var basePrefix string
	if len(files.Source) > 0 {
		basePrefix = strings.TrimPrefix(filepath.Dir(files.Source[0]), "!")
		for i := 1; i < len(files.Source) && hasCommonFolder; i++ {
			comparePath := files.Source[i]
			if strings.HasPrefix(files.Source[i], "!") {
				comparePath = comparePath[1:]
			}
			basePrefixWithSlash := basePrefix + string(os.PathSeparator)
			for !strings.HasPrefix(comparePath, basePrefixWithSlash) {
				basePrefix = filepath.Dir(basePrefix)
				if basePrefix == "." {
					hasCommonFolder = false // if Source[i] doesn't have same prefix
					break
				}
			}
		}
	} else {
		hasCommonFolder = false
	}
	if hasCommonFolder { // if all files are in basePrefix folder, change execution position
		args = append(args, "-C", basePrefix)
		var relativePaths []string
		for _, path := range files.Source {
			ignorePrefix := ""
			if strings.HasPrefix(path, "!") {
				path = path[1:]
				ignorePrefix = "!"
			}
			rel, err := filepath.Rel(basePrefix, path)
			if err != nil {
				fmt.Printf("Error while processing relative paths")
				continue
			}
			relativePaths = append(relativePaths, ignorePrefix + rel)
		}
		args = append(args, relativePaths...) // modified argument appending
	} else {
		args = append(args, files.Source...) // original argument appending
	}
	return args
}

func (p *Plugin) buildUnTarArgs(target string) []string {
	args := []string{}

	args = append(args,
		p.Config.TarExec,
		"-zxf",
		p.DestFile,
	)

	if p.Config.StripComponents > 0 {
		args = append(args, "--strip-components")
		args = append(args, strconv.Itoa(p.Config.StripComponents))
	}

	if p.Config.Overwrite {
		args = append(args, "--overwrite")
	}

	if p.Config.UnlinkFirst {
		args = append(args, "--unlink-first")
	}

	args = append(args,
		"-C",
		target,
	)

	return args
}

// Exec executes the plugin.
func (p *Plugin) Exec() error {
	if len(p.Config.Key) == 0 && len(p.Config.Password) == 0 && len(p.Config.KeyPath) == 0 {
		return errMissingPasswordOrKey
	}

	if len(p.Config.Source) == 0 || len(p.Config.Target) == 0 {
		return errMissingSourceOrTarget
	}

	hosts := trimValues(p.Config.Host)
	if len(hosts) == 0 {
		return errMissingHost
	}

	p.DestFile = random.String(10) + ".tar.gz"

	// create a temporary file for the archive
	dir := os.TempDir()
	src := filepath.Join(dir, p.DestFile)

	// show current version
	fmt.Println("drone-scp version: " + Version)
	// run archive command
	fmt.Println("tar all files into " + src)
	args := p.buildTarArgs(src)
	cmd := exec.Command(p.Config.TarExec, args...)
	if p.Config.Debug {
		fmt.Println("$", strings.Join(cmd.Args, " "))
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	wg := sync.WaitGroup{}
	wg.Add(len(p.Config.Host))
	errChannel := make(chan error)
	finished := make(chan struct{})
	for _, host := range hosts {
		go func(h string) {
			defer wg.Done()
			host, port := p.hostPort(h)
			// Create MakeConfig instance with remote username, server address and path to private key.
			ssh := &easyssh.MakeConfig{
				Server:            host,
				User:              p.Config.Username,
				Password:          p.Config.Password,
				Port:              port,
				Key:               p.Config.Key,
				KeyPath:           p.Config.KeyPath,
				Passphrase:        p.Config.Passphrase,
				Timeout:           p.Config.Timeout,
				Ciphers:           p.Config.Ciphers,
				Fingerprint:       p.Config.Fingerprint,
				UseInsecureCipher: p.Config.UseInsecureCipher,
				Proxy: easyssh.DefaultConfig{
					Server:            p.Config.Proxy.Server,
					User:              p.Config.Proxy.User,
					Password:          p.Config.Proxy.Password,
					Port:              p.Config.Proxy.Port,
					Key:               p.Config.Proxy.Key,
					KeyPath:           p.Config.Proxy.KeyPath,
					Passphrase:        p.Config.Proxy.Passphrase,
					Timeout:           p.Config.Proxy.Timeout,
					Ciphers:           p.Config.Proxy.Ciphers,
					Fingerprint:       p.Config.Proxy.Fingerprint,
					UseInsecureCipher: p.Config.Proxy.UseInsecureCipher,
				},
			}

			systemType := "unix"
			_, _, _, err := ssh.Run("ver", p.Config.CommandTimeout)
			if err == nil {
				systemType = "windows"
			}

			// upload file to the tmp path
			p.DestFile = fmt.Sprintf("%s%s", p.Config.TarTmpPath, p.DestFile)

			p.log(host, "remote server os type is "+systemType)
			// Call Scp method with file you want to upload to remote server.
			p.log(host, "scp file to server.")
			err = ssh.Scp(src, p.DestFile)
			if err != nil {
				errChannel <- copyError{host, err.Error()}
				return
			}

			for _, target := range p.Config.Target {
				target = strings.ReplaceAll(target, " ", "\\ ")
				// remove target folder before upload data
				if p.Config.Remove {
					p.log(host, "Remove target folder:", target)

					_, _, _, err := ssh.Run(rmcmd(systemType, target), p.Config.CommandTimeout)
					if err != nil {
						errChannel <- err
						return
					}
				}

				p.log(host, "create folder", target)
				_, errStr, _, err := ssh.Run(mkdircmd(systemType, target), p.Config.CommandTimeout)
				if err != nil {
					errChannel <- err
					return
				}

				if len(errStr) != 0 {
					errChannel <- fmt.Errorf("%s", errStr)
					return
				}

				// untar file
				p.log(host, "untar file", p.DestFile)
				commamd := strings.Join(p.buildUnTarArgs(target), " ")
				if p.Config.Debug {
					fmt.Println("$", commamd)
				}
				outStr, errStr, _, err := ssh.Run(commamd, p.Config.CommandTimeout)

				if outStr != "" {
					p.log(host, "output: ", outStr)
				}

				if errStr != "" {
					p.log(host, "error: ", errStr)
				}

				if err != nil {
					errChannel <- err
					return
				}
			}

			// remove tar file
			err = p.removeDestFile(systemType, ssh)
			if err != nil {
				errChannel <- err
				return
			}
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
			c := color.New(color.FgRed)
			c.Println("drone-scp error: ", err)
			var cerr copyError
			if !errors.As(err, &cerr) {
				fmt.Println("drone-scp rollback: remove all target tmp file")
				if err := p.removeAllDestFile(); err != nil {
					return err
				}
			}
			return err
		}
	}

	fmt.Println("===================================================")
	fmt.Println("✅ Successfully executed transfer data to all host")
	fmt.Println("===================================================")

	return nil
}

func (p Plugin) hostPort(host string) (string, string) {
	hosts := strings.Split(host, ":")
	port := strconv.Itoa(p.Config.Port)
	if len(hosts) > 1 &&
		(p.Config.Protocol == easyssh.PROTOCOL_TCP ||
			p.Config.Protocol == easyssh.PROTOCOL_TCP4) {
		host = hosts[0]
		port = hosts[1]
	}

	return host, port
}

func trimValues(keys []string) []string {
	var newKeys []string

	for _, value := range keys {
		value = strings.TrimSpace(value)
		if len(value) == 0 {
			continue
		}

		newKeys = append(newKeys, value)
	}

	return newKeys
}
