package main

import (
	"os"

	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"github.com/urfave/cli"
)

// Version set at compile-time
var Version = "v1.0.0-dev"

func main() {

	app := cli.NewApp()
	app.Name = "Drone SCP"
	app.Usage = "Copy files and artifacts via SSH."
	app.Copyright = "Copyright (c) 2017 Bo-Yi Wu"
	app.Authors = []cli.Author{
		{
			Name:  "Bo-Yi Wu",
			Email: "appleboy.tw@gmail.com",
		},
	}
	app.Action = run
	app.Version = Version
	app.Flags = []cli.Flag{
		cli.StringSliceFlag{
			Name:   "host, H",
			Usage:  "Server host",
			EnvVar: "PLUGIN_HOST,SCP_HOST",
		},
		cli.StringFlag{
			Name:   "port, P",
			Value:  "22",
			Usage:  "Server port, default to 22",
			EnvVar: "PLUGIN_PORT,SCP_PORT",
		},
		cli.StringFlag{
			Name:   "username, u",
			Usage:  "Server username",
			EnvVar: "PLUGIN_USERNAME,SCP_USERNAME",
		},
		cli.StringFlag{
			Name:   "password, p",
			Usage:  "Password for password-based authentication",
			EnvVar: "PLUGIN_PASSWORD,SCP_PASSWORD",
		},
		cli.DurationFlag{
			Name:   "timeout,t",
			Usage:  "connection timeout",
			EnvVar: "PLUGIN_TIMEOUT,SCP_TIMEOUT",
		},
		cli.StringFlag{
			Name:   "key, k",
			Usage:  "ssh private key",
			EnvVar: "PLUGIN_KEY,SCP_KEY",
		},
		cli.StringFlag{
			Name:   "key-path, i",
			Usage:  "ssh private key path",
			EnvVar: "PLUGIN_KEY_PATH,SCP_KEY_PATH",
		},
		cli.StringSliceFlag{
			Name:   "target, t",
			Usage:  "Target path on the server",
			EnvVar: "PLUGIN_TARGET,SCP_TARGET",
		},
		cli.StringSliceFlag{
			Name:   "source, s",
			Usage:  "scp file list",
			EnvVar: "PLUGIN_SOURCE,SCP_SOURCE",
		},
		cli.BoolFlag{
			Name:   "rm, r",
			Usage:  "remove target folder before upload data",
			EnvVar: "PLUGIN_RM,SCP_RM",
		},
		cli.StringFlag{
			Name:   "repo.owner",
			Usage:  "repository owner",
			EnvVar: "DRONE_REPO_OWNER",
		},
		cli.StringFlag{
			Name:   "repo.name",
			Usage:  "repository name",
			EnvVar: "DRONE_REPO_NAME",
		},
		cli.StringFlag{
			Name:   "commit.sha",
			Usage:  "git commit sha",
			EnvVar: "DRONE_COMMIT_SHA",
		},
		cli.StringFlag{
			Name:   "commit.branch",
			Value:  "master",
			Usage:  "git commit branch",
			EnvVar: "DRONE_COMMIT_BRANCH",
		},
		cli.StringFlag{
			Name:   "commit.author",
			Usage:  "git author name",
			EnvVar: "DRONE_COMMIT_AUTHOR",
		},
		cli.StringFlag{
			Name:   "commit.message",
			Usage:  "commit message",
			EnvVar: "DRONE_COMMIT_MESSAGE",
		},
		cli.StringFlag{
			Name:   "build.event",
			Value:  "push",
			Usage:  "build event",
			EnvVar: "DRONE_BUILD_EVENT",
		},
		cli.IntFlag{
			Name:   "build.number",
			Usage:  "build number",
			EnvVar: "DRONE_BUILD_NUMBER",
		},
		cli.StringFlag{
			Name:   "build.status",
			Usage:  "build status",
			Value:  "success",
			EnvVar: "DRONE_BUILD_STATUS",
		},
		cli.StringFlag{
			Name:   "build.link",
			Usage:  "build link",
			EnvVar: "DRONE_BUILD_LINK",
		},
		cli.StringFlag{
			Name:  "env-file",
			Usage: "source env file",
		},
	}

	// Override a template
	cli.AppHelpTemplate = `
________                                         ____________________________
\______ \_______  ____   ____   ____            /   _____/\_   ___ \______   \
 |    |  \_  __ \/  _ \ /    \_/ __ \   ______  \_____  \ /    \  \/|     ___/
 |    |   \  | \(  <_> )   |  \  ___/  /_____/  /        \\     \___|    |
/_______  /__|   \____/|___|  /\___  >         /_______  / \______  /____|
        \/                  \/     \/                  \/         \/
                                                            version: {{.Version}}
NAME:
   {{.Name}} - {{.Usage}}

USAGE:
   {{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}{{if .Commands}} command [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}
   {{if len .Authors}}
AUTHOR:
   {{range .Authors}}{{ . }}{{end}}
   {{end}}{{if .Commands}}
COMMANDS:
{{range .Commands}}{{if not .HideHelp}}   {{join .Names ", "}}{{ "\t"}}{{.Usage}}{{ "\n" }}{{end}}{{end}}{{end}}{{if .VisibleFlags}}
GLOBAL OPTIONS:
   {{range .VisibleFlags}}{{.}}
   {{end}}{{end}}{{if .Copyright }}
COPYRIGHT:
   {{.Copyright}}
   {{end}}{{if .Version}}
VERSION:
   {{.Version}}
   {{end}}
REPOSITORY:
    Github: https://github.com/appleboy/drone-scp
`
	app.Run(os.Args)
}

func run(c *cli.Context) error {
	if c.String("env-file") != "" {
		_ = godotenv.Load(c.String("env-file"))
	}

	plugin := Plugin{
		Repo: Repo{
			Owner: c.String("repo.owner"),
			Name:  c.String("repo.name"),
		},
		Build: Build{
			Number:  c.Int("build.number"),
			Event:   c.String("build.event"),
			Status:  c.String("build.status"),
			Commit:  c.String("commit.sha"),
			Branch:  c.String("commit.branch"),
			Author:  c.String("commit.author"),
			Message: c.String("commit.message"),
			Link:    c.String("build.link"),
		},
		Config: Config{
			Host:     c.StringSlice("host"),
			Port:     c.String("port"),
			Username: c.String("username"),
			Password: c.String("password"),
			Timeout:  c.Duration("timeout"),
			Key:      c.String("key"),
			KeyPath:  c.String("key-path"),
			Target:   c.StringSlice("target"),
			Source:   c.StringSlice("source"),
			Remove:   c.Bool("rm"),
		},
	}

	return plugin.Exec()
}
