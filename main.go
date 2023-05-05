package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/appleboy/easyssh-proxy"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
)

// Version set at compile-time
var (
	Version string
)

func main() {
	// Load env-file if it exists first
	if filename, found := os.LookupEnv("PLUGIN_ENV_FILE"); found {
		_ = godotenv.Load(filename)
	}

	app := cli.NewApp()
	app.Name = "Drone SCP"
	app.Usage = "Copy files and artifacts via SSH."
	app.Copyright = "Copyright (c) " + strconv.Itoa(time.Now().Year()) + " Bo-Yi Wu"
	app.Version = Version
	app.Authors = []*cli.Author{
		{
			Name:  "Bo-Yi Wu",
			Email: "appleboy.tw@gmail.com",
		},
	}
	app.Action = run
	app.Version = Version
	app.Flags = []cli.Flag{
		&cli.StringSliceFlag{
			Name:     "host",
			Aliases:  []string{"H"},
			Usage:    "connect to host",
			EnvVars:  []string{"PLUGIN_HOST", "SSH_HOST", "INPUT_HOST"},
			FilePath: ".host",
		},
		&cli.IntFlag{
			Name:    "port",
			Aliases: []string{"p"},
			Usage:   "connect to port",
			EnvVars: []string{"PLUGIN_PORT", "SSH_PORT", "INPUT_PORT"},
			Value:   22,
		},
		&cli.StringFlag{
			Name:    "username",
			Aliases: []string{"user", "u"},
			Usage:   "connect as user",
			EnvVars: []string{"PLUGIN_USERNAME", "PLUGIN_USER", "SSH_USERNAME", "INPUT_USERNAME"},
			Value:   "root",
		},
		&cli.StringFlag{
			Name:    "password",
			Aliases: []string{"P"},
			Usage:   "user password",
			EnvVars: []string{"PLUGIN_PASSWORD", "SSH_PASSWORD", "INPUT_PASSWORD"},
		},
		&cli.DurationFlag{
			Name:    "timeout",
			Usage:   "connection timeout",
			EnvVars: []string{"PLUGIN_TIMEOUT", "SSH_TIMEOUT", "INPUT_TIMEOUT"},
			Value:   30 * time.Second,
		},
		&cli.StringFlag{
			Name:    "ssh-key",
			Usage:   "private ssh key",
			EnvVars: []string{"PLUGIN_SSH_KEY", "PLUGIN_KEY", "SSH_KEY", "INPUT_KEY"},
		},
		&cli.StringFlag{
			Name:    "ssh-passphrase",
			Usage:   "The purpose of the passphrase is usually to encrypt the private key.",
			EnvVars: []string{"PLUGIN_SSH_PASSPHRASE", "PLUGIN_PASSPHRASE", "SSH_PASSPHRASE", "INPUT_PASSPHRASE"},
		},
		&cli.StringFlag{
			Name:    "key-path",
			Aliases: []string{"i"},
			Usage:   "ssh private key path",
			EnvVars: []string{"PLUGIN_KEY_PATH", "SSH_KEY_PATH", "INPUT_KEY_PATH"},
		},
		&cli.StringSliceFlag{
			Name:    "ciphers",
			Usage:   "The allowed cipher algorithms. If unspecified then a sensible",
			EnvVars: []string{"PLUGIN_CIPHERS", "SSH_CIPHERS", "INPUT_CIPHERS"},
		},
		&cli.BoolFlag{
			Name:    "useInsecureCipher",
			Usage:   "include more ciphers with use_insecure_cipher",
			EnvVars: []string{"PLUGIN_USE_INSECURE_CIPHER", "SSH_USE_INSECURE_CIPHER", "INPUT_USE_INSECURE_CIPHER"},
		},
		&cli.StringFlag{
			Name:    "fingerprint",
			Usage:   "fingerprint SHA256 of the host public key, default is to skip verification",
			EnvVars: []string{"PLUGIN_FINGERPRINT", "SSH_FINGERPRINT", "INPUT_FINGERPRINT"},
		},
		&cli.DurationFlag{
			Name:    "command.timeout",
			Usage:   "command timeout",
			EnvVars: []string{"PLUGIN_COMMAND_TIMEOUT", "SSH_COMMAND_TIMEOUT", "INPUT_COMMAND_TIMEOUT"},
			Value:   10 * time.Minute,
		},
		&cli.StringSliceFlag{
			Name:    "target",
			Aliases: []string{"t"},
			Usage:   "Target path on the server",
			EnvVars: []string{"PLUGIN_TARGET", "SSH_TARGET", "INPUT_TARGET"},
		},
		&cli.StringSliceFlag{
			Name:    "source",
			Aliases: []string{"s"},
			Usage:   "scp file list",
			EnvVars: []string{"PLUGIN_SOURCE", "SCP_SOURCE", "INPUT_SOURCE"},
		},
		&cli.BoolFlag{
			Name:    "rm",
			Aliases: []string{"r"},
			Usage:   "remove target folder before upload data",
			EnvVars: []string{"PLUGIN_RM", "SCP_RM", "INPUT_RM"},
		},
		&cli.StringFlag{
			Name:    "proxy.host",
			Usage:   "connect to host of proxy",
			EnvVars: []string{"PLUGIN_PROXY_HOST", "PROXY_SSH_HOST", "INPUT_PROXY_HOST"},
		},
		&cli.StringFlag{
			Name:    "proxy.port",
			Usage:   "connect to port of proxy",
			EnvVars: []string{"PLUGIN_PROXY_PORT", "PROXY_SSH_PORT", "INPUT_PROXY_PORT"},
			Value:   "22",
		},
		&cli.StringFlag{
			Name:    "proxy.username",
			Usage:   "connect as user of proxy",
			EnvVars: []string{"PLUGIN_PROXY_USERNAME", "PLUGIN_PROXY_USER", "PROXY_SSH_USERNAME", "INPUT_PROXY_USERNAME"},
			Value:   "root",
		},
		&cli.StringFlag{
			Name:    "proxy.password",
			Usage:   "user password of proxy",
			EnvVars: []string{"PLUGIN_PROXY_PASSWORD", "PROXY_SSH_PASSWORD", "INPUT_PROXY_PASSWORD"},
		},
		&cli.StringFlag{
			Name:    "proxy.ssh-key",
			Usage:   "private ssh key of proxy",
			EnvVars: []string{"PLUGIN_PROXY_SSH_KEY", "PLUGIN_PROXY_KEY", "PROXY_SSH_KEY", "INPUT_PROXY_KEY"},
		},
		&cli.StringFlag{
			Name:    "proxy.ssh-passphrase",
			Usage:   "The purpose of the passphrase is usually to encrypt the private key.",
			EnvVars: []string{"PLUGIN_PROXY_SSH_PASSPHRASE", "PLUGIN_PROXY_PASSPHRASE", "PROXY_SSH_PASSPHRASE", "INPUT_PROXY_PASSPHRASE"},
		},
		&cli.StringFlag{
			Name:    "proxy.key-path",
			Usage:   "ssh private key path of proxy",
			EnvVars: []string{"PLUGIN_PROXY_KEY_PATH", "PROXY_SSH_KEY_PATH", "INPUT_PROXY_KEY_PATH"},
		},
		&cli.DurationFlag{
			Name:    "proxy.timeout",
			Usage:   "proxy connection timeout",
			EnvVars: []string{"PLUGIN_PROXY_TIMEOUT", "PROXY_SSH_TIMEOUT", "INPUT_PROXY_TIMEOUT"},
		},
		&cli.StringSliceFlag{
			Name:    "proxy.ciphers",
			Usage:   "The allowed cipher algorithms. If unspecified then a sensible",
			EnvVars: []string{"PLUGIN_PROXY_CIPHERS", "PROXY_SSH_CIPHERS", "INPUT_PROXY_CIPHERS"},
		},
		&cli.BoolFlag{
			Name:    "proxy.useInsecureCipher",
			Usage:   "include more ciphers with use_insecure_cipher",
			EnvVars: []string{"PLUGIN_PROXY_USE_INSECURE_CIPHER", "PROXY_SSH_USE_INSECURE_CIPHER", "INPUT_PROXY_USE_INSECURE_CIPHER"},
		},
		&cli.StringFlag{
			Name:    "proxy.fingerprint",
			Usage:   "fingerprint SHA256 of the host public key, default is to skip verification",
			EnvVars: []string{"PLUGIN_PROXY_FINGERPRINT", "PROXY_SSH_FINGERPRINT", "PROXY_FINGERPRINT", "INPUT_PROXY_FINGERPRINT"},
		},
		&cli.IntFlag{
			Name:    "strip.components",
			Usage:   "Remove the specified number of leading path elements.",
			EnvVars: []string{"PLUGIN_STRIP_COMPONENTS", "TAR_STRIP_COMPONENTS", "INPUT_STRIP_COMPONENTS"},
		},
		&cli.StringFlag{
			Name:    "tar.exec",
			Usage:   "Alternative `tar` executable to on the dest host",
			EnvVars: []string{"PLUGIN_TAR_EXEC", "SSH_TAR_EXEC", "INPUT_TAR_EXEC"},
			Value:   "tar",
		},
		&cli.StringFlag{
			Name:    "tar.tmp-path",
			Usage:   "Temporary path for tar file on the dest host",
			EnvVars: []string{"PLUGIN_TAR_TMP_PATH", "SSH_TAR_TMP_PATH", "INPUT_TAR_TMP_PATH"},
		},
		&cli.BoolFlag{
			Name:    "debug",
			Usage:   "remove target folder before upload data",
			EnvVars: []string{"PLUGIN_DEBUG", "INPUT_DEBUG"},
		},
		&cli.BoolFlag{
			Name:    "overwrite",
			Usage:   "use --overwrite flag with tar",
			EnvVars: []string{"PLUGIN_OVERWRITE", "INPUT_OVERWRITE"},
		},
		&cli.BoolFlag{
			Name:    "unlink.first",
			Usage:   "use --unlink-first flag with tar",
			EnvVars: []string{"PLUGIN_UNLINK_FIRST", "INPUT_UNLINK_FIRST"},
		},
		&cli.BoolFlag{
			Name:    "tar.dereference",
			Usage:   "use --dereference flag with tar",
			EnvVars: []string{"PLUGIN_TAR_DEREFERENCE", "INPUT_TAR_DEREFERENCE"},
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
	plugin := Plugin{
		Config: Config{
			Host:              c.StringSlice("host"),
			Port:              c.String("port"),
			Username:          c.String("username"),
			Password:          c.String("password"),
			Passphrase:        c.String("ssh-passphrase"),
			Fingerprint:       c.String("fingerprint"),
			Timeout:           c.Duration("timeout"),
			CommandTimeout:    c.Duration("command.timeout"),
			Key:               c.String("ssh-key"),
			KeyPath:           c.String("key-path"),
			Target:            c.StringSlice("target"),
			Source:            c.StringSlice("source"),
			Remove:            c.Bool("rm"),
			Debug:             c.Bool("debug"),
			StripComponents:   c.Int("strip.components"),
			TarExec:           c.String("tar.exec"),
			TarTmpPath:        c.String("tar.tmp-path"),
			Overwrite:         c.Bool("overwrite"),
			UnlinkFirst:       c.Bool("unlink.first"),
			Ciphers:           c.StringSlice("ciphers"),
			UseInsecureCipher: c.Bool("useInsecureCipher"),
			TarDereference:    c.Bool("tar.dereference"),
			Proxy: easyssh.DefaultConfig{
				Key:               c.String("proxy.ssh-key"),
				Passphrase:        c.String("proxy.ssh-passphrase"),
				Fingerprint:       c.String("proxy.fingerprint"),
				KeyPath:           c.String("proxy.key-path"),
				User:              c.String("proxy.username"),
				Password:          c.String("proxy.password"),
				Server:            c.String("proxy.host"),
				Port:              c.String("proxy.port"),
				Timeout:           c.Duration("proxy.timeout"),
				Ciphers:           c.StringSlice("proxy.ciphers"),
				UseInsecureCipher: c.Bool("proxy.useInsecureCipher"),
			},
		},
	}

	return plugin.Exec()
}
