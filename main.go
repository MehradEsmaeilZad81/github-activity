package main

import (
	"github-activity/internal/cli"
	"github-activity/internal/ghapi"
	"os"
)

func main() {
	api := ghapi.NewClient()
	exitCode := cli.Main(api)
	os.Exit(exitCode)
}
