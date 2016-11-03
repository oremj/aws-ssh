package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/oremj/aws-ssh/awsutils"
)

var filterFlag []string
var verboseFlag = flag.Bool("v", false, "Print all tags")

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [-f filter]...\n", os.Args[0])
		fmt.Fprintf(os.Stderr, `Example: %s -f "tag:App=myapp" -f "tag:Type=web,admin"`+"\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Var((*stringSliceVar)(&filterFlag), "f", "filters that are passed to DescribeInstances")
}

func getInstances(filters []*ec2.Filter) []*ec2.Instance {
	input := &ec2.DescribeInstancesInput{
		Filters: filters,
	}
	res, err := awsutils.EC2Sess.DescribeInstances(input)
	if err != nil {
		panic(err)
	}

	var instances []*ec2.Instance
	for _, reservation := range res.Reservations {
		instances = append(instances, reservation.Instances...)
	}
	return instances
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

func parseFilters(filterStrs []string) ([]*ec2.Filter, error) {
	var filters []*ec2.Filter
	for _, filterStr := range filterStrs {
		filterParts := strings.SplitN(filterStr, "=", 2)
		if len(filterParts) != 2 {
			return nil, errors.New("Filter must be in the format Key=Value,Value,Value")
		}
		filters = append(filters, &ec2.Filter{
			Name:   aws.String(filterParts[0]),
			Values: aws.StringSlice(strings.Split(filterParts[1], ",")),
		})
	}

	return filters, nil
}

func formatInstance(inst *ec2.Instance) string {
	tmp := fmt.Sprintf("%s\t%s\t%s\t", *inst.InstanceId, *inst.PublicDnsName, *inst.PrivateIpAddress)
	if *verboseFlag {
		return tmp + formatTags(inst)
	}
	return tmp + getTag(inst, "Name")
}

func main() {
	flag.Parse()
	filters, err := parseFilters(filterFlag)
	if err != nil {
		log.Fatalf("Could not parse filters: %v", err)
	}
	instances := getInstances(filters)
	for _, inst := range instances {
		if *inst.PublicDnsName != "" {
			fmt.Println(formatInstance(inst))
		}
	}
}
