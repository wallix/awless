package awsdriver

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

var awsPolicies = initAWSPolicies()

func LookupAWSPolicy(service, access string) (*policy, error) {
	if access != "readonly" && access != "full" {
		return nil, errors.New("looking up AWS policies: access value can only be 'readonly' or 'full'")
	}
	for _, p := range awsPolicies {
		name := strings.ToLower(p.Name)
		match := fmt.Sprintf("%s%s", strings.ToLower(service), strings.ToLower(access))
		if strings.Contains(name, match) {
			return p, nil
		}
	}

	return nil, fmt.Errorf("no existing AWS policy with service '%s' and access '%s'", service, access)
}

type policy struct {
	Name string `json:"PolicyName"`
	Id   string `json:"PolicyId"`
	Arn  string `json:"Arn"`
}

func initAWSPolicies() (all []*policy) {
	json.Unmarshal(policiesJSON, &all)
	return
}

var policiesJSON = []byte(`[
        {
            "PolicyName": "AWSDirectConnectReadOnlyAccess", 
            "PolicyId": "ANPAI23HZ27SI6FQMGNQ2", 
            "Arn": "arn:aws:iam::aws:policy/AWSDirectConnectReadOnlyAccess" 
        }, 
        {
            "PolicyName": "AmazonGlacierReadOnlyAccess", 
            "PolicyId": "ANPAI2D5NJKMU274MET4E", 
            "Arn": "arn:aws:iam::aws:policy/AmazonGlacierReadOnlyAccess" 
        }, 
        {
            "PolicyName": "AWSMarketplaceFullAccess", 
            "PolicyId": "ANPAI2DV5ULJSO2FYVPYG", 
            "Arn": "arn:aws:iam::aws:policy/AWSMarketplaceFullAccess" 
        }, 
        {
            "PolicyName": "AutoScalingConsoleReadOnlyAccess", 
            "PolicyId": "ANPAI3A7GDXOYQV3VUQMK", 
            "Arn": "arn:aws:iam::aws:policy/AutoScalingConsoleReadOnlyAccess" 
        }, 
        {
            "PolicyName": "AmazonDMSRedshiftS3Role", 
            "PolicyId": "ANPAI3CCUQ4U5WNC5F6B6", 
            "Arn": "arn:aws:iam::aws:policy/service-role/AmazonDMSRedshiftS3Role" 
        }, 
        {
            "PolicyName": "AWSQuickSightListIAM", 
            "PolicyId": "ANPAI3CH5UUWZN4EKGILO", 
            "Arn": "arn:aws:iam::aws:policy/service-role/AWSQuickSightListIAM" 
        }, 
        {
            "PolicyName": "AWSHealthFullAccess", 
            "PolicyId": "ANPAI3CUMPCPEUPCSXC4Y", 
            "Arn": "arn:aws:iam::aws:policy/AWSHealthFullAccess" 
        }, 
        {
            "PolicyName": "AmazonRDSFullAccess", 
            "PolicyId": "ANPAI3R4QMOG6Q5A4VWVG", 
            "Arn": "arn:aws:iam::aws:policy/AmazonRDSFullAccess" 
        }, 
        {
            "PolicyName": "SupportUser", 
            "PolicyId": "ANPAI3V4GSSN5SJY3P2RO", 
            "Arn": "arn:aws:iam::aws:policy/job-function/SupportUser" 
        }, 
        {
            "PolicyName": "AmazonEC2FullAccess", 
            "PolicyId": "ANPAI3VAJF5ZCRZ7MCQE6", 
            "Arn": "arn:aws:iam::aws:policy/AmazonEC2FullAccess" 
        }, 
        {
            "PolicyName": "AWSElasticBeanstalkReadOnlyAccess", 
            "PolicyId": "ANPAI47KNGXDAXFD4SDHG", 
            "Arn": "arn:aws:iam::aws:policy/AWSElasticBeanstalkReadOnlyAccess" 
        }, 
        {
            "PolicyName": "AWSCertificateManagerReadOnly", 
            "PolicyId": "ANPAI4GSWX6S4MESJ3EWC", 
            "Arn": "arn:aws:iam::aws:policy/AWSCertificateManagerReadOnly" 
        }, 
        {
            "PolicyName": "AWSQuicksightAthenaAccess", 
            "PolicyId": "ANPAI4JB77JXFQXDWNRPM", 
            "Arn": "arn:aws:iam::aws:policy/service-role/AWSQuicksightAthenaAccess" 
        }, 
        {
            "PolicyName": "AWSCodeCommitPowerUser", 
            "PolicyId": "ANPAI4UIINUVGB5SEC57G", 
            "Arn": "arn:aws:iam::aws:policy/AWSCodeCommitPowerUser" 
        }, 
        {
            "PolicyName": "AWSCodeCommitFullAccess", 
            "PolicyId": "ANPAI4VCZ3XPIZLQ5NZV2", 
            "Arn": "arn:aws:iam::aws:policy/AWSCodeCommitFullAccess" 
        }, 
        {
            "PolicyName": "IAMSelfManageServiceSpecificCredentials", 
            "PolicyId": "ANPAI4VT74EMXK2PMQJM2", 
            "Arn": "arn:aws:iam::aws:policy/IAMSelfManageServiceSpecificCredentials" 
        }, 
        {
            "PolicyName": "AmazonSQSFullAccess", 
            "PolicyId": "ANPAI65L554VRJ33ECQS6", 
            "Arn": "arn:aws:iam::aws:policy/AmazonSQSFullAccess" 
        }, 
        {
            "PolicyName": "AWSLambdaFullAccess", 
            "PolicyId": "ANPAI6E2CYYMI4XI7AA5K", 
            "Arn": "arn:aws:iam::aws:policy/AWSLambdaFullAccess" 
        }, 
        {
            "PolicyName": "AWSIoTLogging", 
            "PolicyId": "ANPAI6R6Z2FHHGS454W7W", 
            "Arn": "arn:aws:iam::aws:policy/service-role/AWSIoTLogging" 
        }, 
        {
            "PolicyName": "AmazonEC2RoleforSSM", 
            "PolicyId": "ANPAI6TL3SMY22S4KMMX6", 
            "Arn": "arn:aws:iam::aws:policy/service-role/AmazonEC2RoleforSSM" 
        }, 
        {
            "PolicyName": "AWSCloudHSMRole", 
            "PolicyId": "ANPAI7QIUU4GC66SF26WE", 
            "Arn": "arn:aws:iam::aws:policy/service-role/AWSCloudHSMRole" 
        }, 
        {
            "PolicyName": "IAMFullAccess", 
            "PolicyId": "ANPAI7XKCFMBPM3QQRRVQ", 
            "Arn": "arn:aws:iam::aws:policy/IAMFullAccess" 
        }, 
        {
            "PolicyName": "AmazonInspectorFullAccess", 
            "PolicyId": "ANPAI7Y6NTA27NWNA5U5E", 
            "Arn": "arn:aws:iam::aws:policy/AmazonInspectorFullAccess" 
        }, 
        {
            "PolicyName": "AmazonElastiCacheFullAccess", 
            "PolicyId": "ANPAIA2V44CPHAUAAECKG", 
            "Arn": "arn:aws:iam::aws:policy/AmazonElastiCacheFullAccess" 
        }, 
        {
            "PolicyName": "AWSAgentlessDiscoveryService", 
            "PolicyId": "ANPAIA3DIL7BYQ35ISM4K", 
            "Arn": "arn:aws:iam::aws:policy/AWSAgentlessDiscoveryService" 
        }, 
        {
            "PolicyName": "AWSXrayWriteOnlyAccess", 
            "PolicyId": "ANPAIAACM4LMYSRGBCTM6", 
            "Arn": "arn:aws:iam::aws:policy/AWSXrayWriteOnlyAccess" 
        }, 
        {
            "PolicyName": "AutoScalingReadOnlyAccess", 
            "PolicyId": "ANPAIAFWUVLC2LPLSFTFG", 
            "Arn": "arn:aws:iam::aws:policy/AutoScalingReadOnlyAccess" 
        }, 
        {
            "PolicyName": "AutoScalingFullAccess", 
            "PolicyId": "ANPAIAWRCSJDDXDXGPCFU", 
            "Arn": "arn:aws:iam::aws:policy/AutoScalingFullAccess" 
        }, 
        {
            "PolicyName": "AmazonEC2RoleforAWSCodeDeploy", 
            "PolicyId": "ANPAIAZKXZ27TAJ4PVWGK", 
            "Arn": "arn:aws:iam::aws:policy/service-role/AmazonEC2RoleforAWSCodeDeploy" 
        }, 
        {
            "PolicyName": "AWSMobileHub_ReadOnly", 
            "PolicyId": "ANPAIBXVYVL3PWQFBZFGW", 
            "Arn": "arn:aws:iam::aws:policy/AWSMobileHub_ReadOnly" 
        }, 
        {
            "PolicyName": "CloudWatchEventsBuiltInTargetExecutionAccess", 
            "PolicyId": "ANPAIC5AQ5DATYSNF4AUM", 
            "Arn": "arn:aws:iam::aws:policy/service-role/CloudWatchEventsBuiltInTargetExecutionAccess" 
        }, 
        {
            "PolicyName": "AmazonCloudDirectoryReadOnlyAccess", 
            "PolicyId": "ANPAICMSZQGR3O62KMD6M", 
            "Arn": "arn:aws:iam::aws:policy/AmazonCloudDirectoryReadOnlyAccess" 
        }, 
        {
            "PolicyName": "AWSOpsWorksFullAccess", 
            "PolicyId": "ANPAICN26VXMXASXKOQCG", 
            "Arn": "arn:aws:iam::aws:policy/AWSOpsWorksFullAccess" 
        }, 
        {
            "PolicyName": "AWSOpsWorksCMInstanceProfileRole", 
            "PolicyId": "ANPAICSU3OSHCURP2WIZW", 
            "Arn": "arn:aws:iam::aws:policy/AWSOpsWorksCMInstanceProfileRole" 
        }, 
        {
            "PolicyName": "AWSCodePipelineApproverAccess", 
            "PolicyId": "ANPAICXNWK42SQ6LMDXM2", 
            "Arn": "arn:aws:iam::aws:policy/AWSCodePipelineApproverAccess" 
        }, 
        {
            "PolicyName": "AWSApplicationDiscoveryAgentAccess", 
            "PolicyId": "ANPAICZIOVAGC6JPF3WHC", 
            "Arn": "arn:aws:iam::aws:policy/AWSApplicationDiscoveryAgentAccess" 
        }, 
        {
            "PolicyName": "ViewOnlyAccess", 
            "PolicyId": "ANPAID22R6XPJATWOFDK6", 
            "Arn": "arn:aws:iam::aws:policy/job-function/ViewOnlyAccess" 
        }, 
        {
            "PolicyName": "AmazonElasticMapReduceRole", 
            "PolicyId": "ANPAIDI2BQT2LKXZG36TW", 
            "Arn": "arn:aws:iam::aws:policy/service-role/AmazonElasticMapReduceRole" 
        }, 
        {
            "PolicyName": "AmazonRoute53DomainsReadOnlyAccess", 
            "PolicyId": "ANPAIDRINP6PPTRXYVQCI", 
            "Arn": "arn:aws:iam::aws:policy/AmazonRoute53DomainsReadOnlyAccess" 
        }, 
        {
            "PolicyName": "AWSOpsWorksRole", 
            "PolicyId": "ANPAIDUTMOKHJFAPJV45W", 
            "Arn": "arn:aws:iam::aws:policy/service-role/AWSOpsWorksRole" 
        }, 
        {
            "PolicyName": "ApplicationAutoScalingForAmazonAppStreamAccess", 
            "PolicyId": "ANPAIEL3HJCCWFVHA6KPG", 
            "Arn": "arn:aws:iam::aws:policy/service-role/ApplicationAutoScalingForAmazonAppStreamAccess" 
        }, 
        {
            "PolicyName": "AmazonEC2ContainerRegistryFullAccess", 
            "PolicyId": "ANPAIESRL7KD7IIVF6V4W", 
            "Arn": "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryFullAccess" 
        }, 
        {
            "PolicyName": "SimpleWorkflowFullAccess", 
            "PolicyId": "ANPAIFE3AV6VE7EANYBVM", 
            "Arn": "arn:aws:iam::aws:policy/SimpleWorkflowFullAccess" 
        }, 
        {
            "PolicyName": "AmazonS3FullAccess", 
            "PolicyId": "ANPAIFIR6V6BVTRAHWINE", 
            "Arn": "arn:aws:iam::aws:policy/AmazonS3FullAccess" 
        }, 
        {
            "PolicyName": "AWSStorageGatewayReadOnlyAccess", 
            "PolicyId": "ANPAIFKCTUVOPD5NICXJK", 
            "Arn": "arn:aws:iam::aws:policy/AWSStorageGatewayReadOnlyAccess" 
        }, 
        {
            "PolicyName": "Billing", 
            "PolicyId": "ANPAIFTHXT6FFMIRT7ZEA", 
            "Arn": "arn:aws:iam::aws:policy/job-function/Billing" 
        }, 
        {
            "PolicyName": "QuickSightAccessForS3StorageManagementAnalyticsReadOnly", 
            "PolicyId": "ANPAIFWG3L3WDMR4I7ZJW", 
            "Arn": "arn:aws:iam::aws:policy/service-role/QuickSightAccessForS3StorageManagementAnalyticsReadOnly" 
        }, 
        {
            "PolicyName": "AmazonEC2ContainerRegistryReadOnly", 
            "PolicyId": "ANPAIFYZPA37OOHVIH7KQ", 
            "Arn": "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly" 
        }, 
        {
            "PolicyName": "AmazonElasticMapReduceforEC2Role", 
            "PolicyId": "ANPAIGALS5RCDLZLB3PGS", 
            "Arn": "arn:aws:iam::aws:policy/service-role/AmazonElasticMapReduceforEC2Role" 
        }, 
        {
            "PolicyName": "DatabaseAdministrator", 
            "PolicyId": "ANPAIGBMAW4VUQKOQNVT6", 
            "Arn": "arn:aws:iam::aws:policy/job-function/DatabaseAdministrator" 
        }, 
        {
            "PolicyName": "AmazonRedshiftReadOnlyAccess", 
            "PolicyId": "ANPAIGD46KSON64QBSEZM", 
            "Arn": "arn:aws:iam::aws:policy/AmazonRedshiftReadOnlyAccess" 
        }, 
        {
            "PolicyName": "AmazonEC2ReadOnlyAccess", 
            "PolicyId": "ANPAIGDT4SV4GSETWTBZK", 
            "Arn": "arn:aws:iam::aws:policy/AmazonEC2ReadOnlyAccess" 
        }, 
        {
            "PolicyName": "AWSXrayReadOnlyAccess", 
            "PolicyId": "ANPAIH4OFXWPS6ZX6OPGQ", 
            "Arn": "arn:aws:iam::aws:policy/AWSXrayReadOnlyAccess" 
        }, 
        {
            "PolicyName": "AWSElasticBeanstalkEnhancedHealth", 
            "PolicyId": "ANPAIH5EFJNMOGUUTKLFE", 
            "Arn": "arn:aws:iam::aws:policy/service-role/AWSElasticBeanstalkEnhancedHealth" 
        }, 
        {
            "PolicyName": "AmazonElasticMapReduceReadOnlyAccess", 
            "PolicyId": "ANPAIHP6NH2S6GYFCOINC", 
            "Arn": "arn:aws:iam::aws:policy/AmazonElasticMapReduceReadOnlyAccess" 
        }, 
        {
            "PolicyName": "AWSDirectoryServiceReadOnlyAccess", 
            "PolicyId": "ANPAIHWYO6WSDNCG64M2W", 
            "Arn": "arn:aws:iam::aws:policy/AWSDirectoryServiceReadOnlyAccess" 
        }, 
        {
            "PolicyName": "AmazonVPCReadOnlyAccess", 
            "PolicyId": "ANPAIICZJNOJN36GTG6CM", 
            "Arn": "arn:aws:iam::aws:policy/AmazonVPCReadOnlyAccess" 
        }, 
        {
            "PolicyName": "CloudWatchEventsReadOnlyAccess", 
            "PolicyId": "ANPAIILJPXXA6F7GYLYBS", 
            "Arn": "arn:aws:iam::aws:policy/CloudWatchEventsReadOnlyAccess" 
        }, 
        {
            "PolicyName": "AmazonAPIGatewayInvokeFullAccess", 
            "PolicyId": "ANPAIIWAX2NOOQJ4AIEQ6", 
            "Arn": "arn:aws:iam::aws:policy/AmazonAPIGatewayInvokeFullAccess" 
        }, 
        {
            "PolicyName": "AmazonKinesisAnalyticsReadOnly", 
            "PolicyId": "ANPAIJIEXZAFUK43U7ARK", 
            "Arn": "arn:aws:iam::aws:policy/AmazonKinesisAnalyticsReadOnly" 
        }, 
        {
            "PolicyName": "AmazonMobileAnalyticsFullAccess", 
            "PolicyId": "ANPAIJIKLU2IJ7WJ6DZFG", 
            "Arn": "arn:aws:iam::aws:policy/AmazonMobileAnalyticsFullAccess" 
        }, 
        {
            "PolicyName": "AWSMobileHub_FullAccess", 
            "PolicyId": "ANPAIJLU43R6AGRBK76DM", 
            "Arn": "arn:aws:iam::aws:policy/AWSMobileHub_FullAccess" 
        }, 
        {
            "PolicyName": "AmazonAPIGatewayPushToCloudWatchLogs", 
            "PolicyId": "ANPAIK4GFO7HLKYN64ASK", 
            "Arn": "arn:aws:iam::aws:policy/service-role/AmazonAPIGatewayPushToCloudWatchLogs" 
        }, 
        {
            "PolicyName": "AWSDataPipelineRole", 
            "PolicyId": "ANPAIKCP6XS3ESGF4GLO2", 
            "Arn": "arn:aws:iam::aws:policy/service-role/AWSDataPipelineRole" 
        }, 
        {
            "PolicyName": "CloudWatchFullAccess", 
            "PolicyId": "ANPAIKEABORKUXN6DEAZU", 
            "Arn": "arn:aws:iam::aws:policy/CloudWatchFullAccess" 
        }, 
        {
            "PolicyName": "ServiceCatalogAdminFullAccess", 
            "PolicyId": "ANPAIKTX42IAS75B7B7BY", 
            "Arn": "arn:aws:iam::aws:policy/ServiceCatalogAdminFullAccess" 
        }, 
        {
            "PolicyName": "AmazonRDSDirectoryServiceAccess", 
            "PolicyId": "ANPAIL4KBY57XWMYUHKUU", 
            "Arn": "arn:aws:iam::aws:policy/service-role/AmazonRDSDirectoryServiceAccess" 
        }, 
        {
            "PolicyName": "AWSCodePipelineReadOnlyAccess", 
            "PolicyId": "ANPAILFKZXIBOTNC5TO2Q", 
            "Arn": "arn:aws:iam::aws:policy/AWSCodePipelineReadOnlyAccess" 
        }, 
        {
            "PolicyName": "ReadOnlyAccess", 
            "PolicyId": "ANPAILL3HVNFSB6DCOWYQ", 
            "Arn": "arn:aws:iam::aws:policy/ReadOnlyAccess" 
        }, 
        {
            "PolicyName": "AmazonMachineLearningBatchPredictionsAccess", 
            "PolicyId": "ANPAILOI4HTQSFTF3GQSC", 
            "Arn": "arn:aws:iam::aws:policy/AmazonMachineLearningBatchPredictionsAccess" 
        }, 
        {
            "PolicyName": "AmazonRekognitionReadOnlyAccess", 
            "PolicyId": "ANPAILWSUHXUY4ES43SA4", 
            "Arn": "arn:aws:iam::aws:policy/AmazonRekognitionReadOnlyAccess" 
        }, 
        {
            "PolicyName": "AWSCodeDeployReadOnlyAccess", 
            "PolicyId": "ANPAILZHHKCKB4NE7XOIQ", 
            "Arn": "arn:aws:iam::aws:policy/AWSCodeDeployReadOnlyAccess" 
        }, 
        {
            "PolicyName": "CloudSearchFullAccess", 
            "PolicyId": "ANPAIM6OOWKQ7L7VBOZOC", 
            "Arn": "arn:aws:iam::aws:policy/CloudSearchFullAccess" 
        }, 
        {
            "PolicyName": "AWSCloudHSMFullAccess", 
            "PolicyId": "ANPAIMBQYQZM7F63DA2UU", 
            "Arn": "arn:aws:iam::aws:policy/AWSCloudHSMFullAccess" 
        }, 
        {
            "PolicyName": "AmazonEC2SpotFleetAutoscaleRole", 
            "PolicyId": "ANPAIMFFRMIOBGDP2TAVE", 
            "Arn": "arn:aws:iam::aws:policy/service-role/AmazonEC2SpotFleetAutoscaleRole" 
        }, 
        {
            "PolicyName": "AWSCodeBuildDeveloperAccess", 
            "PolicyId": "ANPAIMKTMR34XSBQW45HS", 
            "Arn": "arn:aws:iam::aws:policy/AWSCodeBuildDeveloperAccess" 
        }, 
        {
            "PolicyName": "AmazonEC2SpotFleetRole", 
            "PolicyId": "ANPAIMRTKHWK7ESSNETSW", 
            "Arn": "arn:aws:iam::aws:policy/service-role/AmazonEC2SpotFleetRole" 
        }, 
        {
            "PolicyName": "AWSDataPipeline_PowerUser", 
            "PolicyId": "ANPAIMXGLVY6DVR24VTYS", 
            "Arn": "arn:aws:iam::aws:policy/AWSDataPipeline_PowerUser" 
        }, 
        {
            "PolicyName": "AmazonElasticTranscoderJobsSubmitter", 
            "PolicyId": "ANPAIN5WGARIKZ3E2UQOU", 
            "Arn": "arn:aws:iam::aws:policy/AmazonElasticTranscoderJobsSubmitter" 
        }, 
        {
            "PolicyName": "AWSCodeStarServiceRole", 
            "PolicyId": "ANPAIN6D4M2KD3NBOC4M4", 
            "Arn": "arn:aws:iam::aws:policy/service-role/AWSCodeStarServiceRole" 
        }, 
        {
            "PolicyName": "AWSDirectoryServiceFullAccess", 
            "PolicyId": "ANPAINAW5ANUWTH3R4ANI", 
            "Arn": "arn:aws:iam::aws:policy/AWSDirectoryServiceFullAccess" 
        }, 
        {
            "PolicyName": "AmazonDynamoDBFullAccess", 
            "PolicyId": "ANPAINUGF2JSOSUY76KYA", 
            "Arn": "arn:aws:iam::aws:policy/AmazonDynamoDBFullAccess" 
        }, 
        {
            "PolicyName": "AmazonSESReadOnlyAccess", 
            "PolicyId": "ANPAINV2XPFRMWJJNSCGI", 
            "Arn": "arn:aws:iam::aws:policy/AmazonSESReadOnlyAccess" 
        }, 
        {
            "PolicyName": "AWSWAFReadOnlyAccess", 
            "PolicyId": "ANPAINZVDMX2SBF7EU2OC", 
            "Arn": "arn:aws:iam::aws:policy/AWSWAFReadOnlyAccess" 
        }, 
        {
            "PolicyName": "AutoScalingNotificationAccessRole", 
            "PolicyId": "ANPAIO2VMUPGDC5PZVXVA", 
            "Arn": "arn:aws:iam::aws:policy/service-role/AutoScalingNotificationAccessRole" 
        }, 
        {
            "PolicyName": "AmazonMechanicalTurkReadOnly", 
            "PolicyId": "ANPAIO5IY3G3WXSX5PPRM", 
            "Arn": "arn:aws:iam::aws:policy/AmazonMechanicalTurkReadOnly" 
        }, 
        {
            "PolicyName": "AmazonKinesisReadOnlyAccess", 
            "PolicyId": "ANPAIOCMTDT5RLKZ2CAJO", 
            "Arn": "arn:aws:iam::aws:policy/AmazonKinesisReadOnlyAccess" 
        }, 
        {
            "PolicyName": "AWSCodeDeployFullAccess", 
            "PolicyId": "ANPAIONKN3TJZUKXCHXWC", 
            "Arn": "arn:aws:iam::aws:policy/AWSCodeDeployFullAccess" 
        }, 
        {
            "PolicyName": "CloudWatchActionsEC2Access", 
            "PolicyId": "ANPAIOWD4E3FVSORSZTGU", 
            "Arn": "arn:aws:iam::aws:policy/CloudWatchActionsEC2Access" 
        }, 
        {
            "PolicyName": "AWSLambdaDynamoDBExecutionRole", 
            "PolicyId": "ANPAIP7WNAGMIPYNW4WQG", 
            "Arn": "arn:aws:iam::aws:policy/service-role/AWSLambdaDynamoDBExecutionRole" 
        }, 
        {
            "PolicyName": "AmazonRoute53DomainsFullAccess", 
            "PolicyId": "ANPAIPAFBMIYUILMOKL6G", 
            "Arn": "arn:aws:iam::aws:policy/AmazonRoute53DomainsFullAccess" 
        }, 
        {
            "PolicyName": "AmazonElastiCacheReadOnlyAccess", 
            "PolicyId": "ANPAIPDACSNQHSENWAKM2", 
            "Arn": "arn:aws:iam::aws:policy/AmazonElastiCacheReadOnlyAccess" 
        }, 
        {
            "PolicyName": "AmazonAthenaFullAccess", 
            "PolicyId": "ANPAIPJMLMD4C7RYZ6XCK", 
            "Arn": "arn:aws:iam::aws:policy/AmazonAthenaFullAccess" 
        }, 
        {
            "PolicyName": "AmazonElasticFileSystemReadOnlyAccess", 
            "PolicyId": "ANPAIPN5S4NE5JJOKVC4Y", 
            "Arn": "arn:aws:iam::aws:policy/AmazonElasticFileSystemReadOnlyAccess" 
        }, 
        {
            "PolicyName": "CloudFrontFullAccess", 
            "PolicyId": "ANPAIPRV52SH6HDCCFY6U", 
            "Arn": "arn:aws:iam::aws:policy/CloudFrontFullAccess" 
        }, 
        {
            "PolicyName": "AmazonMachineLearningRoleforRedshiftDataSource", 
            "PolicyId": "ANPAIQ5UDYYMNN42BM4AK", 
            "Arn": "arn:aws:iam::aws:policy/service-role/AmazonMachineLearningRoleforRedshiftDataSource" 
        }, 
        {
            "PolicyName": "AmazonMobileAnalyticsNon-financialReportAccess", 
            "PolicyId": "ANPAIQLKQ4RXPUBBVVRDE", 
            "Arn": "arn:aws:iam::aws:policy/AmazonMobileAnalyticsNon-financialReportAccess" 
        }, 
        {
            "PolicyName": "AWSCloudTrailFullAccess", 
            "PolicyId": "ANPAIQNUJTQYDRJPC3BNK", 
            "Arn": "arn:aws:iam::aws:policy/AWSCloudTrailFullAccess" 
        }, 
        {
            "PolicyName": "AmazonCognitoDeveloperAuthenticatedIdentities", 
            "PolicyId": "ANPAIQOKZ5BGKLCMTXH4W", 
            "Arn": "arn:aws:iam::aws:policy/AmazonCognitoDeveloperAuthenticatedIdentities" 
        }, 
        {
            "PolicyName": "AWSConfigRole", 
            "PolicyId": "ANPAIQRXRDRGJUA33ELIO", 
            "Arn": "arn:aws:iam::aws:policy/service-role/AWSConfigRole" 
        }, 
        {
            "PolicyName": "AmazonAppStreamServiceAccess", 
            "PolicyId": "ANPAISBRZ7LMMCBYEF3SE", 
            "Arn": "arn:aws:iam::aws:policy/service-role/AmazonAppStreamServiceAccess" 
        }, 
        {
            "PolicyName": "AmazonRedshiftFullAccess", 
            "PolicyId": "ANPAISEKCHH4YDB46B5ZO", 
            "Arn": "arn:aws:iam::aws:policy/AmazonRedshiftFullAccess" 
        }, 
        {
            "PolicyName": "AmazonZocaloReadOnlyAccess", 
            "PolicyId": "ANPAISRCSSJNS3QPKZJPM", 
            "Arn": "arn:aws:iam::aws:policy/AmazonZocaloReadOnlyAccess" 
        }, 
        {
            "PolicyName": "AWSCloudHSMReadOnlyAccess", 
            "PolicyId": "ANPAISVCBSY7YDBOT67KE", 
            "Arn": "arn:aws:iam::aws:policy/AWSCloudHSMReadOnlyAccess" 
        }, 
        {
            "PolicyName": "SystemAdministrator", 
            "PolicyId": "ANPAITJPEZXCYCBXANDSW", 
            "Arn": "arn:aws:iam::aws:policy/job-function/SystemAdministrator" 
        }, 
        {
            "PolicyName": "AmazonEC2ContainerServiceEventsRole", 
            "PolicyId": "ANPAITKFNIUAG27VSYNZ4", 
            "Arn": "arn:aws:iam::aws:policy/service-role/AmazonEC2ContainerServiceEventsRole" 
        }, 
        {
            "PolicyName": "AmazonRoute53ReadOnlyAccess", 
            "PolicyId": "ANPAITOYK2ZAOQFXV2JNC", 
            "Arn": "arn:aws:iam::aws:policy/AmazonRoute53ReadOnlyAccess" 
        }, 
        {
            "PolicyName": "AmazonEC2ReportsAccess", 
            "PolicyId": "ANPAIU6NBZVF2PCRW36ZW", 
            "Arn": "arn:aws:iam::aws:policy/AmazonEC2ReportsAccess" 
        }, 
        {
            "PolicyName": "AmazonEC2ContainerServiceAutoscaleRole", 
            "PolicyId": "ANPAIUAP3EGGGXXCPDQKK", 
            "Arn": "arn:aws:iam::aws:policy/service-role/AmazonEC2ContainerServiceAutoscaleRole" 
        }, 
        {
            "PolicyName": "AWSBatchServiceRole", 
            "PolicyId": "ANPAIUETIXPCKASQJURFE", 
            "Arn": "arn:aws:iam::aws:policy/service-role/AWSBatchServiceRole" 
        }, 
        {
            "PolicyName": "AWSElasticBeanstalkWebTier", 
            "PolicyId": "ANPAIUF4325SJYOREKW3A", 
            "Arn": "arn:aws:iam::aws:policy/AWSElasticBeanstalkWebTier" 
        }, 
        {
            "PolicyName": "AmazonSQSReadOnlyAccess", 
            "PolicyId": "ANPAIUGSSQY362XGCM6KW", 
            "Arn": "arn:aws:iam::aws:policy/AmazonSQSReadOnlyAccess" 
        }, 
        {
            "PolicyName": "AWSMobileHub_ServiceUseOnly", 
            "PolicyId": "ANPAIUHPQXBDZUWOP3PSK", 
            "Arn": "arn:aws:iam::aws:policy/service-role/AWSMobileHub_ServiceUseOnly" 
        }, 
        {
            "PolicyName": "AmazonKinesisFullAccess", 
            "PolicyId": "ANPAIVF32HAMOXCUYRAYE", 
            "Arn": "arn:aws:iam::aws:policy/AmazonKinesisFullAccess" 
        }, 
        {
            "PolicyName": "AmazonMachineLearningReadOnlyAccess", 
            "PolicyId": "ANPAIW5VYBCGEX56JCINC", 
            "Arn": "arn:aws:iam::aws:policy/AmazonMachineLearningReadOnlyAccess" 
        }, 
        {
            "PolicyName": "AmazonRekognitionFullAccess", 
            "PolicyId": "ANPAIWDAOK6AIFDVX6TT6", 
            "Arn": "arn:aws:iam::aws:policy/AmazonRekognitionFullAccess" 
        }, 
        {
            "PolicyName": "RDSCloudHsmAuthorizationRole", 
            "PolicyId": "ANPAIWKFXRLQG2ROKKXLE", 
            "Arn": "arn:aws:iam::aws:policy/service-role/RDSCloudHsmAuthorizationRole" 
        }, 
        {
            "PolicyName": "AmazonMachineLearningFullAccess", 
            "PolicyId": "ANPAIWKW6AGSGYOQ5ERHC", 
            "Arn": "arn:aws:iam::aws:policy/AmazonMachineLearningFullAccess" 
        }, 
        {
            "PolicyName": "AdministratorAccess", 
            "PolicyId": "ANPAIWMBCKSKIEE64ZLYK", 
            "Arn": "arn:aws:iam::aws:policy/AdministratorAccess" 
        }, 
        {
            "PolicyName": "AmazonMachineLearningRealTimePredictionOnlyAccess", 
            "PolicyId": "ANPAIWMCNQPRWMWT36GVQ", 
            "Arn": "arn:aws:iam::aws:policy/AmazonMachineLearningRealTimePredictionOnlyAccess" 
        }, 
        {
            "PolicyName": "AWSConfigUserAccess", 
            "PolicyId": "ANPAIWTTSFJ7KKJE3MWGA", 
            "Arn": "arn:aws:iam::aws:policy/AWSConfigUserAccess" 
        }, 
        {
            "PolicyName": "AWSIoTConfigAccess", 
            "PolicyId": "ANPAIWWGD4LM4EMXNRL7I", 
            "Arn": "arn:aws:iam::aws:policy/AWSIoTConfigAccess" 
        }, 
        {
            "PolicyName": "SecurityAudit", 
            "PolicyId": "ANPAIX2T3QCXHR2OGGCTO", 
            "Arn": "arn:aws:iam::aws:policy/SecurityAudit" 
        }, 
        {
            "PolicyName": "AWSCodeStarFullAccess", 
            "PolicyId": "ANPAIXI233TFUGLZOJBEC", 
            "Arn": "arn:aws:iam::aws:policy/AWSCodeStarFullAccess" 
        }, 
        {
            "PolicyName": "AWSDataPipeline_FullAccess", 
            "PolicyId": "ANPAIXOFIG7RSBMRPHXJ4", 
            "Arn": "arn:aws:iam::aws:policy/AWSDataPipeline_FullAccess" 
        }, 
        {
            "PolicyName": "AmazonDynamoDBReadOnlyAccess", 
            "PolicyId": "ANPAIY2XFNA232XJ6J7X2", 
            "Arn": "arn:aws:iam::aws:policy/AmazonDynamoDBReadOnlyAccess" 
        }, 
        {
            "PolicyName": "AutoScalingConsoleFullAccess", 
            "PolicyId": "ANPAIYEN6FJGYYWJFFCZW", 
            "Arn": "arn:aws:iam::aws:policy/AutoScalingConsoleFullAccess" 
        }, 
        {
            "PolicyName": "AmazonSNSReadOnlyAccess", 
            "PolicyId": "ANPAIZGQCQTFOFPMHSB6W", 
            "Arn": "arn:aws:iam::aws:policy/AmazonSNSReadOnlyAccess" 
        }, 
        {
            "PolicyName": "AmazonElasticMapReduceFullAccess", 
            "PolicyId": "ANPAIZP5JFP3AMSGINBB2", 
            "Arn": "arn:aws:iam::aws:policy/AmazonElasticMapReduceFullAccess" 
        }, 
        {
            "PolicyName": "AmazonS3ReadOnlyAccess", 
            "PolicyId": "ANPAIZTJ4DXE7G6AGAE6M", 
            "Arn": "arn:aws:iam::aws:policy/AmazonS3ReadOnlyAccess" 
        }, 
        {
            "PolicyName": "AWSElasticBeanstalkFullAccess", 
            "PolicyId": "ANPAIZYX2YLLBW2LJVUFW", 
            "Arn": "arn:aws:iam::aws:policy/AWSElasticBeanstalkFullAccess" 
        }, 
        {
            "PolicyName": "AmazonWorkSpacesAdmin", 
            "PolicyId": "ANPAJ26AU6ATUQCT5KVJU", 
            "Arn": "arn:aws:iam::aws:policy/AmazonWorkSpacesAdmin" 
        }, 
        {
            "PolicyName": "AWSCodeDeployRole", 
            "PolicyId": "ANPAJ2NKMKD73QS5NBFLA", 
            "Arn": "arn:aws:iam::aws:policy/service-role/AWSCodeDeployRole" 
        }, 
        {
            "PolicyName": "AmazonSESFullAccess", 
            "PolicyId": "ANPAJ2P4NXCHAT7NDPNR4", 
            "Arn": "arn:aws:iam::aws:policy/AmazonSESFullAccess" 
        }, 
        {
            "PolicyName": "CloudWatchLogsReadOnlyAccess", 
            "PolicyId": "ANPAJ2YIYDYSNNEHK3VKW", 
            "Arn": "arn:aws:iam::aws:policy/CloudWatchLogsReadOnlyAccess" 
        }, 
        {
            "PolicyName": "AmazonKinesisFirehoseReadOnlyAccess", 
            "PolicyId": "ANPAJ36NT645INW4K24W6", 
            "Arn": "arn:aws:iam::aws:policy/AmazonKinesisFirehoseReadOnlyAccess" 
        }, 
        {
            "PolicyName": "AWSOpsWorksRegisterCLI", 
            "PolicyId": "ANPAJ3AB5ZBFPCQGTVDU4", 
            "Arn": "arn:aws:iam::aws:policy/AWSOpsWorksRegisterCLI" 
        }, 
        {
            "PolicyName": "AmazonDynamoDBFullAccesswithDataPipeline", 
            "PolicyId": "ANPAJ3ORT7KDISSXGHJXA", 
            "Arn": "arn:aws:iam::aws:policy/AmazonDynamoDBFullAccesswithDataPipeline" 
        }, 
        {
            "PolicyName": "AmazonEC2RoleforDataPipelineRole", 
            "PolicyId": "ANPAJ3Z5I2WAJE5DN2J36", 
            "Arn": "arn:aws:iam::aws:policy/service-role/AmazonEC2RoleforDataPipelineRole" 
        }, 
        {
            "PolicyName": "CloudWatchLogsFullAccess", 
            "PolicyId": "ANPAJ3ZGNWK2R5HW5BQFO", 
            "Arn": "arn:aws:iam::aws:policy/CloudWatchLogsFullAccess" 
        }, 
        {
            "PolicyName": "AWSElasticBeanstalkMulticontainerDocker", 
            "PolicyId": "ANPAJ45SBYG72SD6SHJEY", 
            "Arn": "arn:aws:iam::aws:policy/AWSElasticBeanstalkMulticontainerDocker" 
        }, 
        {
            "PolicyName": "AmazonElasticTranscoderFullAccess", 
            "PolicyId": "ANPAJ4D5OJU75P5ZJZVNY", 
            "Arn": "arn:aws:iam::aws:policy/AmazonElasticTranscoderFullAccess" 
        }, 
        {
            "PolicyName": "IAMUserChangePassword", 
            "PolicyId": "ANPAJ4L4MM2A7QIEB56MS", 
            "Arn": "arn:aws:iam::aws:policy/IAMUserChangePassword" 
        }, 
        {
            "PolicyName": "AmazonAPIGatewayAdministrator", 
            "PolicyId": "ANPAJ4PT6VY5NLKTNUYSI", 
            "Arn": "arn:aws:iam::aws:policy/AmazonAPIGatewayAdministrator" 
        }, 
        {
            "PolicyName": "ServiceCatalogEndUserAccess", 
            "PolicyId": "ANPAJ56OMCO72RI4J5FSA", 
            "Arn": "arn:aws:iam::aws:policy/ServiceCatalogEndUserAccess" 
        }, 
        {
            "PolicyName": "AmazonPollyReadOnlyAccess", 
            "PolicyId": "ANPAJ5FENL3CVPL2FPDLA", 
            "Arn": "arn:aws:iam::aws:policy/AmazonPollyReadOnlyAccess" 
        }, 
        {
            "PolicyName": "AmazonMobileAnalyticsWriteOnlyAccess", 
            "PolicyId": "ANPAJ5TAWBBQC2FAL3G6G", 
            "Arn": "arn:aws:iam::aws:policy/AmazonMobileAnalyticsWriteOnlyAccess" 
        }, 
        {
            "PolicyName": "AmazonEC2SpotFleetTaggingRole", 
            "PolicyId": "ANPAJ5U6UMLCEYLX5OLC4", 
            "Arn": "arn:aws:iam::aws:policy/service-role/AmazonEC2SpotFleetTaggingRole" 
        }, 
        {
            "PolicyName": "DataScientist", 
            "PolicyId": "ANPAJ5YHI2BQW7EQFYDXS", 
            "Arn": "arn:aws:iam::aws:policy/job-function/DataScientist" 
        }, 
        {
            "PolicyName": "AWSMarketplaceMeteringFullAccess", 
            "PolicyId": "ANPAJ65YJPG7CC7LDXNA6", 
            "Arn": "arn:aws:iam::aws:policy/AWSMarketplaceMeteringFullAccess" 
        }, 
        {
            "PolicyName": "AWSOpsWorksCMServiceRole", 
            "PolicyId": "ANPAJ6I6MPGJE62URSHCO", 
            "Arn": "arn:aws:iam::aws:policy/service-role/AWSOpsWorksCMServiceRole" 
        }, 
        {
            "PolicyName": "AWSConnector", 
            "PolicyId": "ANPAJ6YATONJHICG3DJ3U", 
            "Arn": "arn:aws:iam::aws:policy/AWSConnector" 
        }, 
        {
            "PolicyName": "AWSBatchFullAccess", 
            "PolicyId": "ANPAJ7K2KIWB3HZVK3CUO", 
            "Arn": "arn:aws:iam::aws:policy/AWSBatchFullAccess" 
        }, 
        {
            "PolicyName": "ServiceCatalogAdminReadOnlyAccess", 
            "PolicyId": "ANPAJ7XOUSS75M4LIPKO4", 
            "Arn": "arn:aws:iam::aws:policy/ServiceCatalogAdminReadOnlyAccess" 
        }, 
        {
            "PolicyName": "AmazonSSMFullAccess", 
            "PolicyId": "ANPAJA7V6HI4ISQFMDYAG", 
            "Arn": "arn:aws:iam::aws:policy/AmazonSSMFullAccess" 
        }, 
        {
            "PolicyName": "AWSCodeCommitReadOnly", 
            "PolicyId": "ANPAJACNSXR7Z2VLJW3D6", 
            "Arn": "arn:aws:iam::aws:policy/AWSCodeCommitReadOnly" 
        }, 
        {
            "PolicyName": "AmazonEC2ContainerServiceFullAccess", 
            "PolicyId": "ANPAJALOYVTPDZEMIACSM", 
            "Arn": "arn:aws:iam::aws:policy/AmazonEC2ContainerServiceFullAccess" 
        }, 
        {
            "PolicyName": "AmazonCognitoReadOnly", 
            "PolicyId": "ANPAJBFTRZD2GQGJHSVQK", 
            "Arn": "arn:aws:iam::aws:policy/AmazonCognitoReadOnly" 
        }, 
        {
            "PolicyName": "AmazonDMSCloudWatchLogsRole", 
            "PolicyId": "ANPAJBG7UXZZXUJD3TDJE", 
            "Arn": "arn:aws:iam::aws:policy/service-role/AmazonDMSCloudWatchLogsRole" 
        }, 
        {
            "PolicyName": "AWSApplicationDiscoveryServiceFullAccess", 
            "PolicyId": "ANPAJBNJEA6ZXM2SBOPDU", 
            "Arn": "arn:aws:iam::aws:policy/AWSApplicationDiscoveryServiceFullAccess" 
        }, 
        {
            "PolicyName": "AmazonVPCFullAccess", 
            "PolicyId": "ANPAJBWPGNOVKZD3JI2P2", 
            "Arn": "arn:aws:iam::aws:policy/AmazonVPCFullAccess" 
        }, 
        {
            "PolicyName": "AWSImportExportFullAccess", 
            "PolicyId": "ANPAJCQCT4JGTLC6722MQ", 
            "Arn": "arn:aws:iam::aws:policy/AWSImportExportFullAccess" 
        }, 
        {
            "PolicyName": "AmazonMechanicalTurkFullAccess", 
            "PolicyId": "ANPAJDGCL5BET73H5QIQC", 
            "Arn": "arn:aws:iam::aws:policy/AmazonMechanicalTurkFullAccess" 
        }, 
        {
            "PolicyName": "AmazonEC2ContainerRegistryPowerUser", 
            "PolicyId": "ANPAJDNE5PIHROIBGGDDW", 
            "Arn": "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryPowerUser" 
        }, 
        {
            "PolicyName": "AmazonMachineLearningCreateOnlyAccess", 
            "PolicyId": "ANPAJDRUNIC2RYAMAT3CK", 
            "Arn": "arn:aws:iam::aws:policy/AmazonMachineLearningCreateOnlyAccess" 
        }, 
        {
            "PolicyName": "AWSCloudTrailReadOnlyAccess", 
            "PolicyId": "ANPAJDU7KJADWBSEQ3E7S", 
            "Arn": "arn:aws:iam::aws:policy/AWSCloudTrailReadOnlyAccess" 
        }, 
        {
            "PolicyName": "AWSLambdaExecute", 
            "PolicyId": "ANPAJE5FX7FQZSU5XAKGO", 
            "Arn": "arn:aws:iam::aws:policy/AWSLambdaExecute" 
        }, 
        {
            "PolicyName": "AWSIoTRuleActions", 
            "PolicyId": "ANPAJEZ6FS7BUZVUHMOKY", 
            "Arn": "arn:aws:iam::aws:policy/service-role/AWSIoTRuleActions" 
        }, 
        {
            "PolicyName": "AWSQuickSightDescribeRedshift", 
            "PolicyId": "ANPAJFEM6MLSLTW4ZNBW2", 
            "Arn": "arn:aws:iam::aws:policy/service-role/AWSQuickSightDescribeRedshift" 
        }, 
        {
            "PolicyName": "VMImportExportRoleForAWSConnector", 
            "PolicyId": "ANPAJFLQOOJ6F5XNX4LAW", 
            "Arn": "arn:aws:iam::aws:policy/service-role/VMImportExportRoleForAWSConnector" 
        }, 
        {
            "PolicyName": "AWSCodePipelineCustomActionAccess", 
            "PolicyId": "ANPAJFW5Z32BTVF76VCYC", 
            "Arn": "arn:aws:iam::aws:policy/AWSCodePipelineCustomActionAccess" 
        }, 
        {
            "PolicyName": "AWSOpsWorksInstanceRegistration", 
            "PolicyId": "ANPAJG3LCPVNI4WDZCIMU", 
            "Arn": "arn:aws:iam::aws:policy/AWSOpsWorksInstanceRegistration" 
        }, 
        {
            "PolicyName": "AmazonCloudDirectoryFullAccess", 
            "PolicyId": "ANPAJG3XQK77ATFLCF2CK", 
            "Arn": "arn:aws:iam::aws:policy/AmazonCloudDirectoryFullAccess" 
        }, 
        {
            "PolicyName": "AWSStorageGatewayFullAccess", 
            "PolicyId": "ANPAJG5SSPAVOGK3SIDGU", 
            "Arn": "arn:aws:iam::aws:policy/AWSStorageGatewayFullAccess" 
        }, 
        {
            "PolicyName": "AmazonLexReadOnly", 
            "PolicyId": "ANPAJGBI5LSMAJNDGBNAM", 
            "Arn": "arn:aws:iam::aws:policy/AmazonLexReadOnly" 
        }, 
        {
            "PolicyName": "AmazonElasticTranscoderReadOnlyAccess", 
            "PolicyId": "ANPAJGPP7GPMJRRJMEP3Q", 
            "Arn": "arn:aws:iam::aws:policy/AmazonElasticTranscoderReadOnlyAccess" 
        }, 
        {
            "PolicyName": "AWSIoTConfigReadOnlyAccess", 
            "PolicyId": "ANPAJHENEMXGX4XMFOIOI", 
            "Arn": "arn:aws:iam::aws:policy/AWSIoTConfigReadOnlyAccess" 
        }, 
        {
            "PolicyName": "AmazonWorkMailReadOnlyAccess", 
            "PolicyId": "ANPAJHF7J65E2QFKCWAJM", 
            "Arn": "arn:aws:iam::aws:policy/AmazonWorkMailReadOnlyAccess" 
        }, 
        {
            "PolicyName": "AmazonDMSVPCManagementRole", 
            "PolicyId": "ANPAJHKIGMBQI4AEFFSYO", 
            "Arn": "arn:aws:iam::aws:policy/service-role/AmazonDMSVPCManagementRole" 
        }, 
        {
            "PolicyName": "AWSLambdaKinesisExecutionRole", 
            "PolicyId": "ANPAJHOLKJPXV4GBRMJUQ", 
            "Arn": "arn:aws:iam::aws:policy/service-role/AWSLambdaKinesisExecutionRole" 
        }, 
        {
            "PolicyName": "ResourceGroupsandTagEditorReadOnlyAccess", 
            "PolicyId": "ANPAJHXQTPI5I5JKAIU74", 
            "Arn": "arn:aws:iam::aws:policy/ResourceGroupsandTagEditorReadOnlyAccess" 
        }, 
        {
            "PolicyName": "AmazonSSMAutomationRole", 
            "PolicyId": "ANPAJIBQCTBCXD2XRNB6W", 
            "Arn": "arn:aws:iam::aws:policy/service-role/AmazonSSMAutomationRole" 
        }, 
        {
            "PolicyName": "ServiceCatalogEndUserFullAccess", 
            "PolicyId": "ANPAJIW7AFFOONVKW75KU", 
            "Arn": "arn:aws:iam::aws:policy/ServiceCatalogEndUserFullAccess" 
        }, 
        {
            "PolicyName": "AWSStepFunctionsConsoleFullAccess", 
            "PolicyId": "ANPAJIYC52YWRX6OSMJWK", 
            "Arn": "arn:aws:iam::aws:policy/AWSStepFunctionsConsoleFullAccess" 
        }, 
        {
            "PolicyName": "AWSCodeBuildReadOnlyAccess", 
            "PolicyId": "ANPAJIZZWN6557F5HVP2K", 
            "Arn": "arn:aws:iam::aws:policy/AWSCodeBuildReadOnlyAccess" 
        }, 
        {
            "PolicyName": "AmazonMachineLearningManageRealTimeEndpointOnlyAccess", 
            "PolicyId": "ANPAJJL3PC3VCSVZP6OCI", 
            "Arn": "arn:aws:iam::aws:policy/AmazonMachineLearningManageRealTimeEndpointOnlyAccess" 
        }, 
        {
            "PolicyName": "CloudWatchEventsInvocationAccess", 
            "PolicyId": "ANPAJJXD6JKJLK2WDLZNO", 
            "Arn": "arn:aws:iam::aws:policy/service-role/CloudWatchEventsInvocationAccess" 
        }, 
        {
            "PolicyName": "CloudFrontReadOnlyAccess", 
            "PolicyId": "ANPAJJZMNYOTZCNQP36LG", 
            "Arn": "arn:aws:iam::aws:policy/CloudFrontReadOnlyAccess" 
        }, 
        {
            "PolicyName": "AmazonSNSRole", 
            "PolicyId": "ANPAJK5GQB7CIK7KHY2GA", 
            "Arn": "arn:aws:iam::aws:policy/service-role/AmazonSNSRole" 
        }, 
        {
            "PolicyName": "AmazonMobileAnalyticsFinancialReportAccess", 
            "PolicyId": "ANPAJKJHO2R27TXKCWBU4", 
            "Arn": "arn:aws:iam::aws:policy/AmazonMobileAnalyticsFinancialReportAccess" 
        }, 
        {
            "PolicyName": "AWSElasticBeanstalkService", 
            "PolicyId": "ANPAJKQ5SN74ZQ4WASXBM", 
            "Arn": "arn:aws:iam::aws:policy/service-role/AWSElasticBeanstalkService" 
        }, 
        {
            "PolicyName": "IAMReadOnlyAccess", 
            "PolicyId": "ANPAJKSO7NDY4T57MWDSQ", 
            "Arn": "arn:aws:iam::aws:policy/IAMReadOnlyAccess" 
        }, 
        {
            "PolicyName": "AmazonRDSReadOnlyAccess", 
            "PolicyId": "ANPAJKTTTYV2IIHKLZ346", 
            "Arn": "arn:aws:iam::aws:policy/AmazonRDSReadOnlyAccess" 
        }, 
        {
            "PolicyName": "AmazonCognitoPowerUser", 
            "PolicyId": "ANPAJKW5H2HNCPGCYGR6Y", 
            "Arn": "arn:aws:iam::aws:policy/AmazonCognitoPowerUser" 
        }, 
        {
            "PolicyName": "AmazonElasticFileSystemFullAccess", 
            "PolicyId": "ANPAJKXTMNVQGIDNCKPBC", 
            "Arn": "arn:aws:iam::aws:policy/AmazonElasticFileSystemFullAccess" 
        }, 
        {
            "PolicyName": "ServerMigrationConnector", 
            "PolicyId": "ANPAJKZRWXIPK5HSG3QDQ", 
            "Arn": "arn:aws:iam::aws:policy/ServerMigrationConnector" 
        }, 
        {
            "PolicyName": "AmazonZocaloFullAccess", 
            "PolicyId": "ANPAJLCDXYRINDMUXEVL6", 
            "Arn": "arn:aws:iam::aws:policy/AmazonZocaloFullAccess" 
        }, 
        {
            "PolicyName": "AWSLambdaReadOnlyAccess", 
            "PolicyId": "ANPAJLDG7J3CGUHFN4YN6", 
            "Arn": "arn:aws:iam::aws:policy/AWSLambdaReadOnlyAccess" 
        }, 
        {
            "PolicyName": "AWSAccountUsageReportAccess", 
            "PolicyId": "ANPAJLIB4VSBVO47ZSBB6", 
            "Arn": "arn:aws:iam::aws:policy/AWSAccountUsageReportAccess" 
        }, 
        {
            "PolicyName": "AWSMarketplaceGetEntitlements", 
            "PolicyId": "ANPAJLPIMQE4WMHDC2K7C", 
            "Arn": "arn:aws:iam::aws:policy/AWSMarketplaceGetEntitlements" 
        }, 
        {
            "PolicyName": "AmazonEC2ContainerServiceforEC2Role", 
            "PolicyId": "ANPAJLYJCVHC7TQHCSQDS", 
            "Arn": "arn:aws:iam::aws:policy/service-role/AmazonEC2ContainerServiceforEC2Role" 
        }, 
        {
            "PolicyName": "AmazonAppStreamFullAccess", 
            "PolicyId": "ANPAJLZZXU2YQVGL4QDNC", 
            "Arn": "arn:aws:iam::aws:policy/AmazonAppStreamFullAccess" 
        }, 
        {
            "PolicyName": "AWSIoTDataAccess", 
            "PolicyId": "ANPAJM2KI2UJDR24XPS2K", 
            "Arn": "arn:aws:iam::aws:policy/AWSIoTDataAccess" 
        }, 
        {
            "PolicyName": "AmazonESFullAccess", 
            "PolicyId": "ANPAJM6ZTCU24QL5PZCGC", 
            "Arn": "arn:aws:iam::aws:policy/AmazonESFullAccess" 
        }, 
        {
            "PolicyName": "ServerMigrationServiceRole", 
            "PolicyId": "ANPAJMBH3M6BO63XFW2D4", 
            "Arn": "arn:aws:iam::aws:policy/service-role/ServerMigrationServiceRole" 
        }, 
        {
            "PolicyName": "AWSWAFFullAccess", 
            "PolicyId": "ANPAJMIKIAFXZEGOLRH7C", 
            "Arn": "arn:aws:iam::aws:policy/AWSWAFFullAccess" 
        }, 
        {
            "PolicyName": "AmazonKinesisFirehoseFullAccess", 
            "PolicyId": "ANPAJMZQMTZ7FRBFHHAHI", 
            "Arn": "arn:aws:iam::aws:policy/AmazonKinesisFirehoseFullAccess" 
        }, 
        {
            "PolicyName": "CloudWatchReadOnlyAccess", 
            "PolicyId": "ANPAJN23PDQP7SZQAE3QE", 
            "Arn": "arn:aws:iam::aws:policy/CloudWatchReadOnlyAccess" 
        }, 
        {
            "PolicyName": "AWSLambdaBasicExecutionRole", 
            "PolicyId": "ANPAJNCQGXC42545SKXIK", 
            "Arn": "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole" 
        }, 
        {
            "PolicyName": "ResourceGroupsandTagEditorFullAccess", 
            "PolicyId": "ANPAJNOS54ZFXN4T2Y34A", 
            "Arn": "arn:aws:iam::aws:policy/ResourceGroupsandTagEditorFullAccess" 
        }, 
        {
            "PolicyName": "AWSKeyManagementServicePowerUser", 
            "PolicyId": "ANPAJNPP7PPPPMJRV2SA4", 
            "Arn": "arn:aws:iam::aws:policy/AWSKeyManagementServicePowerUser" 
        }, 
        {
            "PolicyName": "AWSImportExportReadOnlyAccess", 
            "PolicyId": "ANPAJNTV4OG52ESYZHCNK", 
            "Arn": "arn:aws:iam::aws:policy/AWSImportExportReadOnlyAccess" 
        }, 
        {
            "PolicyName": "AmazonElasticTranscoderRole", 
            "PolicyId": "ANPAJNW3WMKVXFJ2KPIQ2", 
            "Arn": "arn:aws:iam::aws:policy/service-role/AmazonElasticTranscoderRole" 
        }, 
        {
            "PolicyName": "AmazonEC2ContainerServiceRole", 
            "PolicyId": "ANPAJO53W2XHNACG7V77Q", 
            "Arn": "arn:aws:iam::aws:policy/service-role/AmazonEC2ContainerServiceRole" 
        }, 
        {
            "PolicyName": "AWSDeviceFarmFullAccess", 
            "PolicyId": "ANPAJO7KEDP4VYJPNT5UW", 
            "Arn": "arn:aws:iam::aws:policy/AWSDeviceFarmFullAccess" 
        }, 
        {
            "PolicyName": "AmazonSSMReadOnlyAccess", 
            "PolicyId": "ANPAJODSKQGGJTHRYZ5FC", 
            "Arn": "arn:aws:iam::aws:policy/AmazonSSMReadOnlyAccess" 
        }, 
        {
            "PolicyName": "AWSStepFunctionsReadOnlyAccess", 
            "PolicyId": "ANPAJONHB2TJQDJPFW5TM", 
            "Arn": "arn:aws:iam::aws:policy/AWSStepFunctionsReadOnlyAccess" 
        }, 
        {
            "PolicyName": "AWSMarketplaceRead-only", 
            "PolicyId": "ANPAJOOM6LETKURTJ3XZ2", 
            "Arn": "arn:aws:iam::aws:policy/AWSMarketplaceRead-only" 
        }, 
        {
            "PolicyName": "AWSCodePipelineFullAccess", 
            "PolicyId": "ANPAJP5LH77KSAT2KHQGG", 
            "Arn": "arn:aws:iam::aws:policy/AWSCodePipelineFullAccess" 
        }, 
        {
            "PolicyName": "AWSGreengrassResourceAccessRolePolicy", 
            "PolicyId": "ANPAJPKEIMB6YMXDEVRTM", 
            "Arn": "arn:aws:iam::aws:policy/service-role/AWSGreengrassResourceAccessRolePolicy" 
        }, 
        {
            "PolicyName": "NetworkAdministrator", 
            "PolicyId": "ANPAJPNMADZFJCVPJVZA2", 
            "Arn": "arn:aws:iam::aws:policy/job-function/NetworkAdministrator" 
        }, 
        {
            "PolicyName": "AmazonWorkSpacesApplicationManagerAdminAccess", 
            "PolicyId": "ANPAJPRL4KYETIH7XGTSS", 
            "Arn": "arn:aws:iam::aws:policy/AmazonWorkSpacesApplicationManagerAdminAccess" 
        }, 
        {
            "PolicyName": "AmazonDRSVPCManagement", 
            "PolicyId": "ANPAJPXIBTTZMBEFEX6UA", 
            "Arn": "arn:aws:iam::aws:policy/AmazonDRSVPCManagement" 
        }, 
        {
            "PolicyName": "AWSXrayFullAccess", 
            "PolicyId": "ANPAJQBYG45NSJMVQDB2K", 
            "Arn": "arn:aws:iam::aws:policy/AWSXrayFullAccess" 
        }, 
        {
            "PolicyName": "AWSElasticBeanstalkWorkerTier", 
            "PolicyId": "ANPAJQDLBRSJVKVF4JMSK", 
            "Arn": "arn:aws:iam::aws:policy/AWSElasticBeanstalkWorkerTier" 
        }, 
        {
            "PolicyName": "AWSDirectConnectFullAccess", 
            "PolicyId": "ANPAJQF2QKZSK74KTIHOW", 
            "Arn": "arn:aws:iam::aws:policy/AWSDirectConnectFullAccess" 
        }, 
        {
            "PolicyName": "AWSCodeBuildAdminAccess", 
            "PolicyId": "ANPAJQJGIOIE3CD2TQXDS", 
            "Arn": "arn:aws:iam::aws:policy/AWSCodeBuildAdminAccess" 
        }, 
        {
            "PolicyName": "AmazonKinesisAnalyticsFullAccess", 
            "PolicyId": "ANPAJQOSKHTXP43R7P5AC", 
            "Arn": "arn:aws:iam::aws:policy/AmazonKinesisAnalyticsFullAccess" 
        }, 
        {
            "PolicyName": "AWSAccountActivityAccess", 
            "PolicyId": "ANPAJQRYCWMFX5J3E333K", 
            "Arn": "arn:aws:iam::aws:policy/AWSAccountActivityAccess" 
        }, 
        {
            "PolicyName": "AmazonGlacierFullAccess", 
            "PolicyId": "ANPAJQSTZJWB2AXXAKHVQ", 
            "Arn": "arn:aws:iam::aws:policy/AmazonGlacierFullAccess" 
        }, 
        {
            "PolicyName": "AmazonWorkMailFullAccess", 
            "PolicyId": "ANPAJQVKNMT7SVATQ4AUY", 
            "Arn": "arn:aws:iam::aws:policy/AmazonWorkMailFullAccess" 
        }, 
        {
            "PolicyName": "AWSMarketplaceManageSubscriptions", 
            "PolicyId": "ANPAJRDW2WIFN7QLUAKBQ", 
            "Arn": "arn:aws:iam::aws:policy/AWSMarketplaceManageSubscriptions" 
        }, 
        {
            "PolicyName": "AWSElasticBeanstalkCustomPlatformforEC2Role", 
            "PolicyId": "ANPAJRVFXSS6LEIQGBKDY", 
            "Arn": "arn:aws:iam::aws:policy/AWSElasticBeanstalkCustomPlatformforEC2Role" 
        }, 
        {
            "PolicyName": "AWSSupportAccess", 
            "PolicyId": "ANPAJSNKQX2OW67GF4S7E", 
            "Arn": "arn:aws:iam::aws:policy/AWSSupportAccess" 
        }, 
        {
            "PolicyName": "AmazonElasticMapReduceforAutoScalingRole", 
            "PolicyId": "ANPAJSVXG6QHPE6VHDZ4Q", 
            "Arn": "arn:aws:iam::aws:policy/service-role/AmazonElasticMapReduceforAutoScalingRole" 
        }, 
        {
            "PolicyName": "AWSLambdaInvocation-DynamoDB", 
            "PolicyId": "ANPAJTHQ3EKCQALQDYG5G", 
            "Arn": "arn:aws:iam::aws:policy/AWSLambdaInvocation-DynamoDB" 
        }, 
        {
            "PolicyName": "IAMUserSSHKeys", 
            "PolicyId": "ANPAJTSHUA4UXGXU7ANUA", 
            "Arn": "arn:aws:iam::aws:policy/IAMUserSSHKeys" 
        }, 
        {
            "PolicyName": "AWSIoTFullAccess", 
            "PolicyId": "ANPAJU2FPGG6PQWN72V2G", 
            "Arn": "arn:aws:iam::aws:policy/AWSIoTFullAccess" 
        }, 
        {
            "PolicyName": "AWSQuickSightDescribeRDS", 
            "PolicyId": "ANPAJU5J6OAMCJD3OO76O", 
            "Arn": "arn:aws:iam::aws:policy/service-role/AWSQuickSightDescribeRDS" 
        }, 
        {
            "PolicyName": "AWSConfigRulesExecutionRole", 
            "PolicyId": "ANPAJUB3KIKTA4PU4OYAA", 
            "Arn": "arn:aws:iam::aws:policy/service-role/AWSConfigRulesExecutionRole" 
        }, 
        {
            "PolicyName": "AmazonESReadOnlyAccess", 
            "PolicyId": "ANPAJUDMRLOQ7FPAR46FQ", 
            "Arn": "arn:aws:iam::aws:policy/AmazonESReadOnlyAccess" 
        }, 
        {
            "PolicyName": "AWSCodeDeployDeployerAccess", 
            "PolicyId": "ANPAJUWEPOMGLMVXJAPUI", 
            "Arn": "arn:aws:iam::aws:policy/AWSCodeDeployDeployerAccess" 
        }, 
        {
            "PolicyName": "AmazonPollyFullAccess", 
            "PolicyId": "ANPAJUZOYQU6XQYPR7EWS", 
            "Arn": "arn:aws:iam::aws:policy/AmazonPollyFullAccess" 
        }, 
        {
            "PolicyName": "AmazonSSMMaintenanceWindowRole", 
            "PolicyId": "ANPAJV3JNYSTZ47VOXYME", 
            "Arn": "arn:aws:iam::aws:policy/service-role/AmazonSSMMaintenanceWindowRole" 
        }, 
        {
            "PolicyName": "AmazonRDSEnhancedMonitoringRole", 
            "PolicyId": "ANPAJV7BS425S4PTSSVGK", 
            "Arn": "arn:aws:iam::aws:policy/service-role/AmazonRDSEnhancedMonitoringRole" 
        }, 
        {
            "PolicyName": "AmazonLexFullAccess", 
            "PolicyId": "ANPAJVLXDHKVC23HRTKSI", 
            "Arn": "arn:aws:iam::aws:policy/AmazonLexFullAccess" 
        }, 
        {
            "PolicyName": "AWSLambdaVPCAccessExecutionRole", 
            "PolicyId": "ANPAJVTME3YLVNL72YR2K", 
            "Arn": "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole" 
        }, 
        {
            "PolicyName": "AmazonLexRunBotsOnly", 
            "PolicyId": "ANPAJVZGB5CM3N6YWJHBE", 
            "Arn": "arn:aws:iam::aws:policy/AmazonLexRunBotsOnly" 
        }, 
        {
            "PolicyName": "AmazonSNSFullAccess", 
            "PolicyId": "ANPAJWEKLCXXUNT2SOLSG", 
            "Arn": "arn:aws:iam::aws:policy/AmazonSNSFullAccess" 
        }, 
        {
            "PolicyName": "CloudSearchReadOnlyAccess", 
            "PolicyId": "ANPAJWPLX7N7BCC3RZLHW", 
            "Arn": "arn:aws:iam::aws:policy/CloudSearchReadOnlyAccess" 
        }, 
        {
            "PolicyName": "AWSGreengrassFullAccess", 
            "PolicyId": "ANPAJWPV6OBK4QONH4J3O", 
            "Arn": "arn:aws:iam::aws:policy/AWSGreengrassFullAccess" 
        }, 
        {
            "PolicyName": "AWSCloudFormationReadOnlyAccess", 
            "PolicyId": "ANPAJWVBEE4I2POWLODLW", 
            "Arn": "arn:aws:iam::aws:policy/AWSCloudFormationReadOnlyAccess" 
        }, 
        {
            "PolicyName": "AmazonRoute53FullAccess", 
            "PolicyId": "ANPAJWVDLG5RPST6PHQ3A", 
            "Arn": "arn:aws:iam::aws:policy/AmazonRoute53FullAccess" 
        }, 
        {
            "PolicyName": "AWSLambdaRole", 
            "PolicyId": "ANPAJX4DPCRGTC4NFDUXI", 
            "Arn": "arn:aws:iam::aws:policy/service-role/AWSLambdaRole" 
        }, 
        {
            "PolicyName": "AWSLambdaENIManagementAccess", 
            "PolicyId": "ANPAJXAW2Q3KPTURUT2QC", 
            "Arn": "arn:aws:iam::aws:policy/service-role/AWSLambdaENIManagementAccess" 
        }, 
        {
            "PolicyName": "AWSOpsWorksCloudWatchLogs", 
            "PolicyId": "ANPAJXFIK7WABAY5CPXM4", 
            "Arn": "arn:aws:iam::aws:policy/AWSOpsWorksCloudWatchLogs" 
        }, 
        {
            "PolicyName": "AmazonAppStreamReadOnlyAccess", 
            "PolicyId": "ANPAJXIFDGB4VBX23DX7K", 
            "Arn": "arn:aws:iam::aws:policy/AmazonAppStreamReadOnlyAccess" 
        }, 
        {
            "PolicyName": "AWSStepFunctionsFullAccess", 
            "PolicyId": "ANPAJXKA6VP3UFBVHDPPA", 
            "Arn": "arn:aws:iam::aws:policy/AWSStepFunctionsFullAccess" 
        }, 
        {
            "PolicyName": "AmazonInspectorReadOnlyAccess", 
            "PolicyId": "ANPAJXQNTHTEJ2JFRN2SE", 
            "Arn": "arn:aws:iam::aws:policy/AmazonInspectorReadOnlyAccess" 
        }, 
        {
            "PolicyName": "AWSCertificateManagerFullAccess", 
            "PolicyId": "ANPAJYCHABBP6VQIVBCBQ", 
            "Arn": "arn:aws:iam::aws:policy/AWSCertificateManagerFullAccess" 
        }, 
        {
            "PolicyName": "PowerUserAccess", 
            "PolicyId": "ANPAJYRXTHIB4FOVS3ZXS", 
            "Arn": "arn:aws:iam::aws:policy/PowerUserAccess" 
        }, 
        {
            "PolicyName": "CloudWatchEventsFullAccess", 
            "PolicyId": "ANPAJZLOYLNHESMYOJAFU", 
            "Arn": "arn:aws:iam::aws:policy/CloudWatchEventsFullAccess" 
        }
    ]`)
