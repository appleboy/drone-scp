package main

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/urfave/cli"
)

// Version for command line
var Version string

func main() {
	app := cli.NewApp()
	app.Name = "telegram plugin"
	app.Usage = "telegram plugin"
	app.Action = run
	app.Version = Version
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "host",
			Usage:  "Server host",
			EnvVar: "PLUGIN_HOST,SCP_HOST",
		},
		cli.StringFlag{
			Name:   "port",
			Value:  "22",
			Usage:  "Server port, default to 22",
			EnvVar: "PLUGIN_PORT,SCP_PORT",
		},
		cli.StringFlag{
			Name:   "username",
			Usage:  "Server username",
			EnvVar: "PLUGIN_USERNAME,SCP_USERNAME",
		},
		cli.StringFlag{
			Name:   "password",
			Usage:  "Password for password-based authentication",
			EnvVar: "PLUGIN_PASSWORD,SCP_PASSWORD",
		},
		cli.StringFlag{
			Name:   "path",
			Value:  "/",
			Usage:  "Target path on the server, default to '/'",
			EnvVar: "PLUGIN_DESTINATION",
		},
		cli.StringSliceFlag{
			Name:   "files",
			Usage:  "scp file list",
			EnvVar: "PLUGIN_FILE",
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
	}
	app.Run(os.Args)
}

func run(c *cli.Context) error {
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
			Host:     c.String("host"),
			Port:     c.String("port"),
			Username: c.String("username"),
			Password: c.String("password"),
			Path:     c.String("path"),
			File:     c.StringSlice("files"),
		},
	}

	return plugin.Exec()
}
