package main

import (
	"os"

	"github.com/bpcoder16/baize/bootstrap"
	"github.com/bpcoder16/baize/bootstrap/aliyunossboot"
	"github.com/bpcoder16/baize/bootstrap/aliyunstsboot"
	esv8 "github.com/bpcoder16/baize/bootstrap/elasticsearchboot/v8"
	"github.com/bpcoder16/baize/bootstrap/mysqlboot"
	"github.com/bpcoder16/baize/bootstrap/natsboot"
	"github.com/bpcoder16/baize/bootstrap/redisboot"
)

func main() {
	os.Exit(bootstrap.RunCLI(bootstrap.Plan{
		Resources: []bootstrap.Resource{
			mysqlboot.Instances("default"),
			redisboot.Instances("default"),
			natsboot.Instances("default"),
			esv8.Instances("default"),
			aliyunossboot.Instances("default"),
			aliyunstsboot.Instances("oss"),
		},
	}))
}
