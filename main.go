package main

import (
	"os"

	"github.com/bpcoder16/baize-demo/bootstrap"

	baize "github.com/bpcoder16/baize/bootstrap"
)

func main() {
	os.Exit(baize.RunCLI(bootstrap.NewPlan()))
}
