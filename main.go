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
			Name:   "base.url",
			Usage:  "jenkins base url",
			EnvVar: "PLUGIN_BASE_URL,JENKINS_BASE_URL",
		},
		cli.StringFlag{
			Name:   "username",
			Usage:  "jenkins username",
			EnvVar: "PLUGIN_USERNAME,JENKINS_USERNAME",
		},
		cli.StringFlag{
			Name:   "token",
			Usage:  "jenkins token",
			EnvVar: "PLUGIN_TOKEN,JENKINS_TOKEN",
		},
		cli.StringSliceFlag{
			Name:   "job",
			Usage:  "jenkins job",
			EnvVar: "PLUGIN_JOB",
		},
		cli.StringFlag{
			Name:   "format",
			Value:  "markdown",
			Usage:  "telegram message format",
			EnvVar: "PLUGIN_FORMAT",
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
			BaseURL:  c.String("base.url"),
			Username: c.String("username"),
			Token:    c.String("token"),
			Job:      c.StringSlice("job"),
		},
	}

	return plugin.Exec()
}
