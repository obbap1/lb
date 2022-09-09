package main

import (
	"log"
	"os"

	"github.com/obbap1/lb.git/cli"
)

func main() {
	if err := cli.App().Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
