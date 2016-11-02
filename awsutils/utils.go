package awsutils

import (
	"os"

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
