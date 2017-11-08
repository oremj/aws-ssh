package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/oremj/aws-tools/awsutils"
)

var (
	filterFlag  []string
	verboseFlag bool
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [-f filter]...\n", os.Args[0])
		fmt.Fprintf(os.Stderr, `Example: %s -f "tag:App=myapp" -f "tag:Type=web,admin"`+"\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.BoolVar(&verboseFlag, "v", false, "Print all tags")
	flag.Var((*awsutils.StringSliceVar)(&filterFlag), "f", "filters that are passed to DescribeInstances")
}

func getTag(inst *ec2.Instance, tagKey string) string {
	for _, tag := range inst.Tags {
		if *tag.Key == tagKey {
			return *tag.Value
		}
	}
	return ""
}

func formatTags(inst *ec2.Instance) string {
	var tags []string
	for _, tag := range inst.Tags {
		tags = append(tags, fmt.Sprintf("%s:%s", *tag.Key, *tag.Value))
	}

	return strings.Join(tags, "|")
}

func formatInstance(inst *ec2.Instance) string {
	tmp := fmt.Sprintf("%s\t%s\t%s\t", *inst.InstanceId, *inst.PublicDnsName, *inst.PrivateIpAddress)
	if verboseFlag {
		return tmp + formatTags(inst)
	}
	return tmp + getTag(inst, "Name")
}

func main() {
	flag.Parse()
	filters, err := awsutils.ParseFilters(filterFlag)
	if err != nil {
		log.Fatalf("Could not parse filters: %v", err)
	}
	instances := awsutils.GetInstances(filters)
	for _, inst := range instances {
		if *inst.PublicDnsName != "" {
			fmt.Println(formatInstance(inst))
		}
	}
}
