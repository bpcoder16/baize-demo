package main

import (
	"os"

	"github.com/bpcoder16/baize/bootstrap"
	"github.com/bpcoder16/baize/bootstrap/mysqlboot"
)

func main() {
	os.Exit(bootstrap.RunCLI(bootstrap.Plan{
		Resources: []bootstrap.Resource{
			mysqlboot.Instances("default"),
		},
	}))
}
