package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"syscall"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

var command = flag.String("c", "", "command to run, if set")
var ec2Sess = ec2.New(session.Must(session.NewSession()))

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [-c command] instance-id [ssh-args]\n", os.Args[0])
		flag.PrintDefaults()
	}
}

func getInstanceHostname(instanceID string) string {
	input := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{
			aws.String(instanceID),
		},
	}

	res, err := ec2Sess.DescribeInstances(input)
	if err != nil {
		panic(err)
	}

	if len(res.Reservations) != 1 {
		log.Fatal("Reservations is not 1")
	}

	if len(res.Reservations[0].Instances) != 1 {
		log.Fatal("Instances is not 1")
	}

	return *res.Reservations[0].Instances[0].PublicDnsName
}

func main() {
	flag.Parse()
	args := flag.Args()
	instanceID := args[0]
	sshArgs := args[1:]

	hostName := getInstanceHostname(instanceID)

	sshArgv := append([]string{"ssh"}, sshArgs...)
	sshArgv = append(sshArgv, hostName)
	if *command != "" {
		sshArgv = append(sshArgv, *command)
	}
	if err := syscall.Exec("/usr/bin/ssh", sshArgv, []string{}); err != nil {
		panic(err)
	}
}
