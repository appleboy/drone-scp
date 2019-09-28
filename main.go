package main

import (
	"log"
	"os"
	"time"

	"github.com/appleboy/easyssh-proxy"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"github.com/urfave/cli"
)

// Version set at compile-time
var (
	Version string
)

func main() {
	app := cli.NewApp()
	app.Name = "Drone SCP"
	app.Usage = "Copy files and artifacts via SSH."
	app.Copyright = "Copyright (c) 2019 Bo-Yi Wu"
	app.Version = Version
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
			EnvVar: "PLUGIN_HOST,SCP_HOST,SSH_HOST,HOST,INPUT_HOST",
		},
		cli.StringFlag{
			Name:   "port, P",
			Value:  "22",
			Usage:  "Server port, default to 22",
			EnvVar: "PLUGIN_PORT,SCP_PORT,SSH_PORT,PORT,INPUT_PORT",
		},
		cli.StringFlag{
			Name:   "username, u",
			Usage:  "Server username",
			EnvVar: "PLUGIN_USERNAME,PLUGIN_USER,SCP_USERNAME,SSH_USERNAME,USERNAME,INPUT_USERNAME",
		},
		cli.StringFlag{
			Name:   "password, p",
			Usage:  "Password for password-based authentication",
			EnvVar: "PLUGIN_PASSWORD,SCP_PASSWORD,SSH_PASSWORD,PASSWORD,INPUT_PASSWORD",
		},
		cli.DurationFlag{
			Name:   "timeout",
			Usage:  "connection timeout",
			EnvVar: "PLUGIN_TIMEOUT,SCP_TIMEOUT,INPUT_TIMEOUT",
		},
		cli.DurationFlag{
			Name:   "command.timeout,T",
			Usage:  "command timeout",
			EnvVar: "PLUGIN_COMMAND_TIMEOUT,SSH_COMMAND_TIMEOUT,INPUT_COMMAND_TIMEOUT",
			Value:  60 * time.Second,
		},
		cli.StringFlag{
			Name:   "key, k",
			Usage:  "ssh private key",
			EnvVar: "PLUGIN_KEY,SCP_KEY,SSH_KEY,KEY,INPUT_KEY",
		},
		cli.StringFlag{
			Name:   "key-path, i",
			Usage:  "ssh private key path",
			EnvVar: "PLUGIN_KEY_PATH,SCP_KEY_PATH,SSH_KEY_PATH,INPUT_KEY_PATH",
		},
		cli.StringSliceFlag{
			Name:   "target, t",
			Usage:  "Target path on the server",
			EnvVar: "PLUGIN_TARGET,SCP_TARGET,TARGET,INPUT_TARGET",
		},
		cli.StringSliceFlag{
			Name:   "source, s",
			Usage:  "scp file list",
			EnvVar: "PLUGIN_SOURCE,SCP_SOURCE,SOURCE,INPUT_SOURCE",
		},
		cli.BoolFlag{
			Name:   "rm, r",
			Usage:  "remove target folder before upload data",
			EnvVar: "PLUGIN_RM,SCP_RM,RM,INPUT_RM",
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
		cli.StringFlag{
			Name:   "proxy.ssh-key",
			Usage:  "private ssh key of proxy",
			EnvVar: "PLUGIN_PROXY_SSH_KEY,PLUGIN_PROXY_KEY,PROXY_SSH_KEY,PROXY_KEY,INPUT_PROXY_SSH_KEY",
		},
		cli.StringFlag{
			Name:   "proxy.key-path",
			Usage:  "ssh private key path of proxy",
			EnvVar: "PLUGIN_PROXY_KEY_PATH,PROXY_SSH_KEY_PATH,INPUT_PROXY_SSH_KEY_PATH",
		},
		cli.StringFlag{
			Name:   "proxy.username",
			Usage:  "connect as user of proxy",
			EnvVar: "PLUGIN_PROXY_USERNAME,PLUGIN_PROXY_USER,PROXY_SSH_USERNAME,PROXY_USERNAME,INPUT_PROXY_USERNAME",
			Value:  "root",
		},
		cli.StringFlag{
			Name:   "proxy.password",
			Usage:  "user password of proxy",
			EnvVar: "PLUGIN_PROXY_PASSWORD,PROXY_SSH_PASSWORD,PROXY_PASSWORD,INPUT_PROXY_PASSWORD",
		},
		cli.StringFlag{
			Name:   "proxy.host",
			Usage:  "connect to host of proxy",
			EnvVar: "PLUGIN_PROXY_HOST,PROXY_SSH_HOST,PROXY_HOST,INPUT_PROXY_HOST",
		},
		cli.StringFlag{
			Name:   "proxy.port",
			Usage:  "connect to port of proxy",
			EnvVar: "PLUGIN_PROXY_PORT,PROXY_SSH_PORT,PROXY_PORT,INPUT_PROXY_PORT",
			Value:  "22",
		},
		cli.DurationFlag{
			Name:   "proxy.timeout",
			Usage:  "proxy connection timeout",
			EnvVar: "PLUGIN_PROXY_TIMEOUT,PROXY_SSH_TIMEOUT,INPUT_PROXY_TIMEOUT",
		},
		cli.IntFlag{
			Name:   "strip.components",
			Usage:  "Remove the specified number of leading path elements.",
			EnvVar: "PLUGIN_STRIP_COMPONENTS,TAR_STRIP_COMPONENTS,INPUT_STRIP_COMPONENTS",
		},
		cli.StringFlag{
			Name:   "tar.exec",
			Usage:  "Alternative `tar` executable to on the dest host",
			EnvVar: "PLUGIN_TAR_EXEC,SCP_TAR_EXEC,INPUT_TAR_EXEC",
			Value:  "tar",
		},
		cli.StringFlag{
			Name:   "tar.tmp-path",
			Usage:  "Temporary path for tar file on the dest host",
			EnvVar: "PLUGIN_TAR_TMP_PATH,SCP_TAR_TMP_PATH",
		},
		cli.BoolFlag{
			Name:   "debug",
			Usage:  "remove target folder before upload data",
			EnvVar: "PLUGIN_DEBUG,DEBUG,INPUT_DEBUG",
		},
		cli.BoolFlag{
			Name:   "overwrite",
			Usage:  "use --overwrite flag with tar",
			EnvVar: "PLUGIN_OVERWRITE,SCP_OVERWRITE,INPUT_OVERWRITE",
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

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
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
			Host:            c.StringSlice("host"),
			Port:            c.String("port"),
			Username:        c.String("username"),
			Password:        c.String("password"),
			Timeout:         c.Duration("timeout"),
			CommandTimeout:  c.Duration("command.timeout"),
			Key:             c.String("key"),
			KeyPath:         c.String("key-path"),
			Target:          c.StringSlice("target"),
			Source:          c.StringSlice("source"),
			Remove:          c.Bool("rm"),
			Debug:           c.Bool("debug"),
			StripComponents: c.Int("strip.components"),
			TarExec:         c.String("tar.exec"),
			TarTmpPath:      c.String("tar.tmp-path"),
			Overwrite:       c.Bool("overwrite"),
			Proxy: easyssh.DefaultConfig{
				Key:      c.String("proxy.ssh-key"),
				KeyPath:  c.String("proxy.key-path"),
				User:     c.String("proxy.username"),
				Password: c.String("proxy.password"),
				Server:   c.String("proxy.host"),
				Port:     c.String("proxy.port"),
				Timeout:  c.Duration("proxy.timeout"),
			},
		},
	}

	return plugin.Exec()
}
