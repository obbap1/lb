package cli

import (
	"errors"
	"fmt"
	"os/user"
	"strings"

	"github.com/obbap1/lb.git/cgroups"
	"github.com/urfave/cli"
)

var (
	minimum int
	maximum int
	port    int
	cpu     float64
	mem     string
)

func App() *cli.App {
	// 1. Read Config
	return &cli.App{
		Name:  "lboom",
		Usage: "load balance and autoscale requests amongst processes",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:        "min",
				Value:       1,
				Usage:       "minimum number of processes to run. Must be at least 1.",
				Destination: &minimum,
			},
			&cli.IntFlag{
				Name:        "max",
				Value:       1,
				Usage:       "maximum number of processes to run. Must be at most the number of logical cores.",
				Destination: &maximum,
			},
			&cli.Float64Flag{
				Name:        "cpu",
				Value:       0,
				Usage:       "Max percentage of the CPU that this process can use. Zero means there's no limit.",
				Destination: &cpu,
			},
			&cli.StringFlag{
				Name:        "mem",
				Usage:       "Maximum memory that this process can use. Units: kb, mb, gb. Zero means there's no limit.",
				Destination: &mem,
			},
			&cli.IntFlag{
				Name:        "port",
				Value:       8080,
				Usage:       "port number to run the load balancer server on.",
				Destination: &port,
			},
		},
		Action: func(cCtx *cli.Context) error {
			// the name of the file must be sent along
			// else there is nothing to do.
			// Its best to handle args for the node process with an env file instead of sending them along
			if cCtx.NArg() < 1 {
				return fmt.Errorf("invalid number of arguments %d. Expecting at least 1", cCtx.NArg())
			}
			filePath := cCtx.Args().Get(0)
			if err := validate(filePath, mem, port, minimum, maximum, cpu); err != nil {
				return err
			}
			// Start the server

			// Initialize Cgroups
			cgroups.Start()

			return nil
		},
	}

}

func validate(filePath, mem string, port, minimum, maximum int, cpu float64) error {
	if strings.Contains(filePath, "~") {
		usr, err := user.Current()
		// if there is an error, then we can basically use the string "~"
		if err == nil {
			// Eg. ~/a/b/c/d.js becomes /users/paschal/a/b/c/d.js
			strings.Replace(filePath, "~", usr.HomeDir, -1)
		}
	}

	fileName := filePath
	// If it has directories, then the fileName is the last string after the last path
	if strings.Contains(filePath, "/") {
		fileDirs := strings.Split(filePath, "/")
		fileName = fileDirs[len(fileDirs)-1]
	}
	s := strings.Split(fileName, ".")
	if len(s) != 2 {
		return errors.New("invalid file name")
	}

	if s[1] != "js" {
		return errors.New("invalid file extension")
	}

	if len(mem) > 0 && (len(mem) < 2 || !isValidUnit(mem[len(mem)-3:])) {
		return errors.New("invalid memory limit")
	}

	if port < 1 {
		return errors.New("invalid port")
	}

	if minimum < 1 {
		return errors.New("invalid minimum number")
	}

	// if a very high maximum is set, a maximum higher than the logical cores
	// it will be replaced by number of logical cores
	if maximum < 1 {
		return errors.New("invalid maximum number")
	}

	if cpu > 100 {
		return errors.New("invalid cpu limit")
	}

	return nil
}

func isValidUnit(unit string) bool {
	switch unit {
	case "kb":
	case "mb":
	case "gb":
		return true
	default:
		return false
	}
	return false
}
