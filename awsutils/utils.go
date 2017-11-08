package awsutils

import (
	"errors"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

var AWSSess *session.Session
var EC2Sess *ec2.EC2

func init() {
	AWSSess := session.Must(session.NewSession())

	// If AWS_REGION is not set, try to detect region
	if envRegion := os.Getenv("AWS_REGION"); envRegion == "" {
		meta := ec2metadata.New(AWSSess)
		if region, _ := meta.Region(); region != "" {
			AWSSess = AWSSess.Copy(&aws.Config{Region: aws.String(region)})
		}
	}

	EC2Sess = ec2.New(AWSSess)
}

func GetInstances(filters []*ec2.Filter) []*ec2.Instance {
	input := &ec2.DescribeInstancesInput{
		Filters: filters,
	}
	res, err := EC2Sess.DescribeInstances(input)
	if err != nil {
		panic(err)
	}

	var instances []*ec2.Instance
	for _, reservation := range res.Reservations {
		instances = append(instances, reservation.Instances...)
	}
	return instances
}

func ParseFilters(filterStrs []string) ([]*ec2.Filter, error) {
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
