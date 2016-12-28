package main

import (
	"os"
	"runtime"

	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"github.com/urfave/cli"
)

// Version set at compile-time
var Version = "v1.0.0-dev"

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	app := cli.NewApp()
	app.Name = "scp plugin"
	app.Usage = "scp plugin"
	app.Action = run
	app.Version = Version
	app.Flags = []cli.Flag{
		cli.StringSliceFlag{
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
			Name:   "key",
			Usage:  "ssh private key",
			EnvVar: "PLUGIN_KEY,SCP_KEY",
		},
		cli.StringSliceFlag{
			Name:   "target",
			Usage:  "Target path on the server",
			EnvVar: "PLUGIN_TARGET,SCP_TARGET",
		},
		cli.StringSliceFlag{
			Name:   "source",
			Usage:  "scp file list",
			EnvVar: "PLUGIN_SOURCE,SCP_SOURCE",
		},
		cli.BoolFlag{
			Name:   "debug",
			Usage:  "display message from command",
			EnvVar: "PLUGIN_DEBUG,SCP_DEBUG",
		},
		cli.BoolFlag{
			Name:   "rm",
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
			Name:   "env-file",
			Usage:  "source env file",
			EnvVar: "ENV_FILE",
		},
	}
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
			Key:      c.String("key"),
			Target:   c.StringSlice("target"),
			Source:   c.StringSlice("source"),
			Debug:    c.Bool("debug"),
			Remove:   c.Bool("rm"),
		},
	}

	return plugin.Exec()
}
