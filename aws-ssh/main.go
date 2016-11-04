package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/oremj/aws-tools/awsutils"
)

var command string

func init() {
	flag.StringVar(&command, "c", "", "command to run, if set")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [-c command] [user@]instance-id [ssh-args]\n", os.Args[0])
		flag.PrintDefaults()
	}

}

func getInstanceHostname(instanceID string) string {
	input := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{
			aws.String(instanceID),
		},
	}

	res, err := awsutils.EC2Sess.DescribeInstances(input)
	if err != nil {
		panic(err)
	}

	if len(res.Reservations) < 1 {
		log.Fatal("Instance not found.")
	}

	if len(res.Reservations) > 1 {
		log.Fatal("Too many reservations found.")
	}

	if len(res.Reservations[0].Instances) != 1 {
		log.Fatal("Instances is not 1")
	}

	return *res.Reservations[0].Instances[0].PublicDnsName
}

func parseHost(host string) (string, string) {
	parts := strings.SplitN(host, "@", 2)
	if len(parts) == 1 {
		return "", parts[0]
	}
	return parts[0] + "@", parts[1]
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		flag.Usage()
		log.Fatal("instance-id is required")
	}

	user, instanceID := parseHost(args[0])
	sshArgs := args[1:]

	hostName := getInstanceHostname(instanceID)

	sshArgv := append([]string{"ssh"}, sshArgs...)
	sshArgv = append(sshArgv, user+hostName)
	if command != "" {
		sshArgv = append(sshArgv, command)
	}
	if err := syscall.Exec("/usr/bin/ssh", sshArgv, os.Environ()); err != nil {
		panic(err)
	}
}
