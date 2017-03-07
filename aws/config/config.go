package config

var DefaultAMIUsers = []string{"ec2-user", "ubuntu", "centos", "bitnami", "admin", "root"}

var AmiPerRegion = map[string]string{
	"us-east-1":      "ami-0b33d91d",
	"us-east-2":      "ami-c55673a0",
	"us-west-1":      "ami-165a0876",
	"us-west-2":      "ami-f173cc91",
	"ca-central-1":   "ami-ebed508f",
	"eu-west-1":      "ami-70edb016",
	"eu-west-2":      "ami-f1949e95",
	"eu-central-1":   "ami-af0fc0c0",
	"ap-southeast-1": "ami-dc9339bf",
	"ap-southeast-2": "ami-1c47407f",
	"ap-northeast-1": "ami-56d4ad31",
	"ap-northeast-2": "ami-dac312b4",
	"ap-south-1":     "ami-f9daac96",
	"sa-east-1":      "ami-80086dec",
}
