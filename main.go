package main

import (
	"os"
	"runtime"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/harborapp/harbor-api/cmd"
	"github.com/harborapp/harbor-api/config"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	cfg := &config.Config{}

	app := cli.NewApp()
	app.Name = "harbor"
	app.Version = config.Version
	app.Author = "Thomas Boerger <thomas@webhippie.de>"
	app.Usage = "A docker distribution management system"

	app.Before = func(c *cli.Context) error {
		logrus.SetOutput(os.Stdout)

		if cfg.Debug {
			logrus.SetLevel(logrus.DebugLevel)
		} else {
			logrus.SetLevel(logrus.InfoLevel)
		}

		return nil
	}

	app.Commands = []cli.Command{
		cmd.Server(cfg),
	}

	cli.HelpFlag = cli.BoolFlag{
		Name:  "help, h",
		Usage: "Show the help, so what you see now",
	}

	cli.VersionFlag = cli.BoolFlag{
		Name:  "version, v",
		Usage: "Print the current version of that tool",
	}

	app.Run(os.Args)
}
