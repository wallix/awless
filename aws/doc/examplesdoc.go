package awsdoc

import (
	"bytes"
	"fmt"
)

func AwlessExamplesDoc(action, entity string) string {
	return exampleDoc(action + entity)
}

func exampleDoc(key string) string {
	examples, ok := cliExamplesDoc[key]
	if ok {
		var buf bytes.Buffer
		for i, ex := range examples {
			buf.WriteString(fmt.Sprintf("  %s", ex))
			if i != len(examples)-1 {
				buf.WriteByte('\n')
			}
		}
		return buf.String()
	}
	return ""
}

var cliExamplesDoc = map[string][]string{
	"attachalarm":         {},
	"attachcontainertask": {},
	"attachelasticip": {
		"attach elasticip id=eipalloc-1c517b26 instance=@redis",
	},
	"attachinstance": {},
	"attachinstanceprofile": {
		"attach instanceprofile instance=@redis name=MyProfile replace=true",
	},
	"attachinternetgateway": {
		"attach internetgateway id=igw-636c0504 vpc=vpc-1aba387c",
	},
	"attachpolicy": {
		"awless attach policy role=MyNewRole service=ec2 access=readonly",
		"awless attach policy user=jsmith service=s3 access=readonly",
	},
	"attachrole": {
		"attach role instanceprofile=MyProfile name=MyRole",
	},
	"attachroutetable": {
		"attach routetable id=rtb-306da254 subnet=@my-subnet",
	},
	"attachsecuritygroup": {
		"attach securitygroup id=sg-0714247d instance=@redis",
	},
	"attachuser": {
		"attach user name=jsmith group=AdminGroup",
	},
	"attachvolume": {
		"attach volume id=vol-123oefwejf device=/dev/sdh instance=@redis",
	},
	"authenticateregistry": {
		"awless authenticate registry",
	},
	"checkdatabase": {
		"awless check database id=@mydb state=available timeout=180",
	},
	"checkdistribution": {
		"awless check distribution id=@mydistr state=Deployed timeout=180",
	},
	"checkinstance": {
		"awless check instance id=@redis state=running timeout=180",
	},
	"checkloadbalancer": {
		"awless check loadbalancer id=@myloadb state=active timeout=180",
	},
	"checknatgateway": {
		"awless check natgateway id=@mynat state=active timeout=180",
	},
	"checkscalinggroup": {
		"awless check scalinggroup name=MyAutoScalingGroup count=3 timeout=180",
	},
	"checksecuritygroup": {
		"awless check securitygroup id=@mysshsecgroup state=unused timeout=180",
	},
	"checkvolume": {
		"awless check volume id=vol-12r1o3rp state=available timeout=180",
	},
	"copyimage": {
		"awless copy image name=my-ami-name source-id=ami-23or2or source-region=us-west-2",
	},
	"copysnapshot": {
		"awless copy snapshot source-id=efwqwdr2or source-region=us-west-2",
	},
	"createaccesskey": {
		"awless create accesskey user=jsmith no-prompt=true",
	},
	"createalarm": {
		" awless create alarm namespace=AWS/EC2 dimensions=AutoScalingGroupName:instancesScalingGroup evaluation-periods=2 metric=CPUUtilization name=scaleinAlarm operator=GreaterThanOrEqualToThreshold period=300 statistic-function=Average threshold=75",
	},
	"createappscalingpolicy": {
		" awless create appscalingpolicy dimension=ecs:service:DesiredCount name=ScaleOutPolicy resource=service/my-ecs-cluster/my-service-deployment-name service-namespace=ecs stepscaling-adjustment-type=ChangeInCapacity stepscaling-adjustments=0::+1 type=StepScaling stepscaling-aggregation-type=Average stepscaling-cooldown=60",
	},
	"createappscalingtarget": {
		"awless create appscalingtarget dimension=ecs:service:DesiredCount min-capacity=2 max-capacity=10 resource=service/my-ecs-cluster/my-service-deployment-nameource role=arn:aws:iam::519101889238:role/ecsAutoscaleRole service-namespace=ecs",
	},
	"createbucket": {
		"awless create bucket name=my-bucket-name acl=public-read",
	},
	"createcontainercluster": {
		"awless create containercluster name=mycluster",
	},
	"createdatabase": {
		"awless create database engine=postgres id=mystartup-prod-db subnetgroup=@my-dbsubnetgroup  password=notsafe dbname=mydb size=5 type=db.t2.small username=admin vpcsecuritygroups=@postgres_sg",
	},
	"createdbsubnetgroup": {
		"awless create dbsubnetgroup name=mydbsubnetgroup description=\"subnets for peps db\" subnets=[@my-firstsubnet, @my-secondsubnet]",
	},
	"createdistribution": {
		"awless create distribution origin-domain=mybucket.s3.amazonaws.com",
	},
	"createelasticip": {
		"awless create elasticip domain=vpc",
	},
	"createfunction": {},
	"creategroup": {
		"awless create name=admins",
	},
	"createimage": {
		"awless create image instance=@my-instance-name name=redis-image description='redis prod image'",
		"awless create image instance=i-0ee436a45561c04df name=redis-image no-reboot=true",
		"awless create image instance=@redis-prod name=redis-prod-image",
	},
	"createinstance": {
		"awless create instance keypair=jsmith type=t2.micro subnet=@my-subnet",
		"awless create instance image=ami-123456 keypair=jsmith",
		"awless create instance name=redis type=t2.nano keypair=jsmith userdata=/home/jsmith/data.sh",
		"", // create empty line for clarity
		"awless create instance distro=redhat type=t2.micro",
		"awless create instance distro=redhat::7.2 type=t2.micro",
		"awless create instance distro=canonical:ubuntu role=MyInfraReadOnlyRole",
		"awless create instance distro=debian:debian:jessie lock=true",
		"awless create instance distro=amazonlinux securitygroup=@my-ssh-secgroup",
		"awless create instance distro=amazonlinux:::::instance-store",
	},
	"createinstanceprofile":     {},
	"createinternetgateway":     {},
	"createkeypair":             {},
	"createlaunchconfiguration": {},
	"createlistener":            {},
	"createloadbalancer":        {},
	"createloginprofile":        {},
	"createnatgateway":          {},
	"createpolicy":              {},
	"createqueue":               {},
	"createrecord":              {},
	"createrepository":          {},
	"createrole":                {},
	"createroute":               {},
	"createroutetable":          {},
	"creates3object":            {},
	"createscalinggroup":        {},
	"createscalingpolicy":       {},
	"createsecuritygroup": {
		"awless create securitygroup vpc=@myvpc name=ssh-only description=ssh-access",
		"(... see more params at `awless update securitygroup -h`)",
	},
	"createsnapshot":            {},
	"createstack":               {},
	"createsubnet":              {},
	"createsubscription":        {},
	"createtag":                 {},
	"createtargetgroup":         {},
	"createtopic":               {},
	"createuser":                {},
	"createvolume":              {},
	"createvpc":                 {},
	"createzone":                {},
	"deleteaccesskey":           {},
	"deletealarm":               {},
	"deleteappscalingpolicy":    {},
	"deleteappscalingtarget":    {},
	"deletebucket":              {},
	"deletecontainercluster":    {},
	"deletecontainertask":       {},
	"deletedatabase":            {},
	"deletedbsubnetgroup":       {},
	"deletedistribution":        {},
	"deleteelasticip":           {},
	"deletefunction":            {},
	"deletegroup":               {},
	"deleteimage":               {},
	"deleteinstance":            {},
	"deleteinstanceprofile":     {},
	"deleteinternetgateway":     {},
	"deletekeypair":             {},
	"deletelaunchconfiguration": {},
	"deletelistener":            {},
	"deleteloadbalancer":        {},
	"deleteloginprofile":        {},
	"deletenatgateway":          {},
	"deletepolicy":              {},
	"deletequeue":               {},
	"deleterecord":              {},
	"deleterepository":          {},
	"deleterole":                {},
	"deleteroute":               {},
	"deleteroutetable":          {},
	"deletes3object":            {},
	"deletescalinggroup":        {},
	"deletescalingpolicy":       {},
	"deletesecuritygroup":       {},
	"deletesnapshot":            {},
	"deletestack":               {},
	"deletesubnet":              {},
	"deletesubscription":        {},
	"deletetag":                 {},
	"deletetargetgroup":         {},
	"deletetopic":               {},
	"deleteuser":                {},
	"deletevolume":              {},
	"deletevpc":                 {},
	"deletezone":                {},
	"detachalarm":               {},
	"detachcontainertask":       {},
	"detachelasticip":           {},
	"detachinstance":            {},
	"detachinstanceprofile":     {},
	"detachinternetgateway":     {},
	"detachpolicy":              {},
	"detachrole":                {},
	"detachroutetable":          {},
	"detachsecuritygroup":       {},
	"detachuser":                {},
	"detachvolume":              {},
	"importimage":               {},
	"startalarm":                {},
	"startcontainertask":        {},
	"startinstance":             {},
	"stopalarm":                 {},
	"stopcontainertask":         {},
	"stopinstance":              {},
	"updatebucket":              {},
	"updatecontainertask":       {},
	"updatedistribution":        {},
	"updateinstance":            {},
	"updateloginprofile":        {},
	"updatepolicy":              {},
	"updaterecord":              {},
	"updates3object":            {},
	"updatescalinggroup":        {},
	"updatesecuritygroup": {
		"awless update securitygroup id=@ssh-only inbound=authorize protocol=tcp cidr=0.0.0.0/0 portrange=26257",
		"awless update securitygroup id=@ssh-only inbound=authorize protocol=tcp securitygroup=sg-123457 portrange=8080",
	},
	"updatestack":       {},
	"updatesubnet":      {},
	"updatetargetgroup": {},
}
