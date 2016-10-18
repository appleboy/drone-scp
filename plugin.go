package main

import (
	"errors"
	"log"
	"strings"
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
		BaseURL  string
		Username string
		Token    string
		Job      []string
	}

	// Plugin values.
	Plugin struct {
		Repo   Repo
		Build  Build
		Config Config
	}
)

func trimElement(keys []string) []string {
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

	if len(p.Config.BaseURL) == 0 || len(p.Config.Username) == 0 || len(p.Config.Token) == 0 {
		log.Println("missing jenkins config")

		return errors.New("missing jenkins config")
	}

	auth := &Auth{
		Username: p.Config.Username,
		Token:    p.Config.Token,
	}
	jenkins := NewJenkins(auth, p.Config.BaseURL)

	for _, value := range trimElement(p.Config.Job) {
		jenkins.trigger(value, nil)
	}

	return nil
}
