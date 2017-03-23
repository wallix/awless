package cloud

// Resources
const (
	Region string = "region"
	//infra
	Vpc              string = "vpc"
	Subnet           string = "subnet"
	Image            string = "image"
	SecurityGroup    string = "securitygroup"
	AvailabilityZone string = "availabilityzone"
	Keypair          string = "keypair"
	Volume           string = "volume"
	Instance         string = "instance"
	InternetGateway  string = "internetgateway"
	RouteTable       string = "routetable"
	//loadbalancer
	LoadBalancer string = "loadbalancer"
	TargetGroup  string = "targetgroup"
	Listener     string = "listener"
	//database
	Database      string = "database"
	DbSubnetGroup string = "dbsubnetgroup"
	//access
	User      string = "user"
	Role      string = "role"
	Group     string = "group"
	Policy    string = "policy"
	AccessKey string = "accesskey"
	//storage
	Bucket string = "bucket"
	Object string = "storageobject"
	Acl    string = "storageacl"
	//notification
	Subscription string = "subscription"
	Topic        string = "topic"
	//queue
	Queue string = "queue"
	//dns
	Zone   string = "zone"
	Record string = "record"
)
