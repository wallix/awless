package rdf

import (
	"fmt"

	"github.com/wallix/awless/cloud/properties"
)

type rdfProp struct {
	ID, RdfType, RdfsLabel, RdfsDefinedBy, RdfsDataType string
}

// Namespaces
const (
	RdfsNS     = "rdfs"
	RdfNS      = "rdf"
	CloudNS    = "cloud"
	CloudOwlNS = "cloud-owl"
	XsdNS      = "xsd"
	netNS      = "net"
	netowlNS   = "net-owl"
)

// Existing terms
var (
	RdfsLabel       = fmt.Sprintf("%s:label", RdfsNS)
	RdfsList        = fmt.Sprintf("%s:list", RdfsNS)
	RdfsDefinedBy   = fmt.Sprintf("%s:isDefinedBy", RdfsNS)
	RdfsDataType    = fmt.Sprintf("%s:Datatype", RdfsNS)
	RdfsSeeAlso     = fmt.Sprintf("%s:seeAlso", RdfsNS)
	RdfsLiteral     = fmt.Sprintf("%s:Literal", RdfsNS)
	RdfsClass       = fmt.Sprintf("%s:Class", RdfsNS)
	RdfsSubProperty = fmt.Sprintf("%s:subPropertyOf", RdfsNS)
	RdfsComment     = fmt.Sprintf("%s:comment", RdfsNS)

	RdfType     = fmt.Sprintf("%s:type", RdfNS)
	RdfProperty = fmt.Sprintf("%s:Property", RdfNS)

	XsdString   = fmt.Sprintf("%s:string", XsdNS)
	XsdBoolean  = fmt.Sprintf("%s:boolean", XsdNS)
	XsdInt      = fmt.Sprintf("%s:int", XsdNS)
	XsdDateTime = fmt.Sprintf("%s:dateTime", XsdNS)
)

// Classes
var (
	Grant = fmt.Sprintf("%s:Grant", CloudOwlNS)

	NetFirewallRule = fmt.Sprintf("%s:FirewallRule", netowlNS)
	NetRoute        = fmt.Sprintf("%s:Route", netowlNS)

	Permission = fmt.Sprintf("%s:permission", CloudNS)
	Grantee    = fmt.Sprintf("%s:grantee", CloudNS)

	NetRouteTargets          = fmt.Sprintf("%s:routeTargets", netNS)
	NetDestinationPrefixList = fmt.Sprintf("%s:routeDestinationPrefixList", netNS)
)

// Properties
var (
	Actions                   = fmt.Sprintf("%s:actions", CloudNS)
	Affinity                  = fmt.Sprintf("%s:affinity", CloudNS)
	ApproximateMessageCount   = fmt.Sprintf("%s:approximateMessageCount", CloudNS)
	Architecture              = fmt.Sprintf("%s:architecture", CloudNS)
	Arn                       = fmt.Sprintf("%s:arn", CloudNS)
	Attachable                = fmt.Sprintf("%s:attachable", CloudNS)
	AutoUpgrade               = fmt.Sprintf("%s:autoUpgrade", CloudNS)
	AvailabilityZone          = fmt.Sprintf("%s:availabilityZone", CloudNS)
	AvailabilityZones         = fmt.Sprintf("%s:availabilityZones", CloudNS)
	BackupRetentionPeriod     = fmt.Sprintf("%s:backupRetentionPeriod", CloudNS)
	Bucket                    = fmt.Sprintf("%s:bucketName", CloudNS)
	CallerReference           = fmt.Sprintf("%s:callerReference", CloudNS)
	CertificateAuthority      = fmt.Sprintf("%s:certificateAuthority", CloudNS)
	Certificates              = fmt.Sprintf("%s:certificates", CloudNS)
	Charset                   = fmt.Sprintf("%s:charset", CloudNS)
	CheckHTTPCode             = fmt.Sprintf("%s:checkHTTPCode", CloudNS)
	CheckInterval             = fmt.Sprintf("%s:checkInterval", CloudNS)
	CheckPath                 = fmt.Sprintf("%s:checkPath", CloudNS)
	CheckPort                 = fmt.Sprintf("%s:checkPort", CloudNS)
	CheckProtocol             = fmt.Sprintf("%s:checkProtocol", CloudNS)
	CheckTimeout              = fmt.Sprintf("%s:checkTimeout", CloudNS)
	CIDR                      = fmt.Sprintf("%s:cidr", netNS)
	CIDRv6                    = fmt.Sprintf("%s:cidrv6", netNS)
	CipherSuite               = fmt.Sprintf("%s:cipherSuite", CloudNS)
	Class                     = fmt.Sprintf("%s:class", CloudNS)
	Cluster                   = fmt.Sprintf("%s:cluster", CloudNS)
	Comment                   = RdfsComment
	Continent                 = fmt.Sprintf("%s:continent", CloudNS)
	CopyTagsToSnapshot        = fmt.Sprintf("%s:copyTagsToSnapshot", CloudNS)
	Country                   = fmt.Sprintf("%s:country", CloudNS)
	Created                   = fmt.Sprintf("%s:created", CloudNS)
	DBSecurityGroups          = fmt.Sprintf("%s:dbSecurityGroups", CloudNS)
	DBSubnetGroup             = fmt.Sprintf("%s:dbSubnetGroup", CloudNS)
	Default                   = fmt.Sprintf("%s:default", CloudNS)
	Delay                     = fmt.Sprintf("%s:delaySeconds", CloudNS)
	Description               = fmt.Sprintf("%s:description", CloudNS)
	Encrypted                 = fmt.Sprintf("%s:encrypted", CloudNS)
	Endpoint                  = fmt.Sprintf("%s:endpoint", CloudNS)
	Engine                    = fmt.Sprintf("%s:engine", CloudNS)
	EngineVersion             = fmt.Sprintf("%s:engineVersion", CloudNS)
	Failover                  = fmt.Sprintf("%s:failover", CloudNS)
	Fingerprint               = fmt.Sprintf("%s:fingerprint", CloudNS)
	GlobalID                  = fmt.Sprintf("%s:globalID", CloudNS)
	Grants                    = fmt.Sprintf("%s:grants", CloudNS)
	HealthCheck               = fmt.Sprintf("%s:healthCheck", CloudNS)
	HealthyThresholdCount     = fmt.Sprintf("%s:healthyThresholdCount", CloudNS)
	Host                      = fmt.Sprintf("%s:host", CloudNS)
	Hypervisor                = fmt.Sprintf("%s:hypervisor", CloudNS)
	ID                        = fmt.Sprintf("%s:id", CloudNS)
	Image                     = fmt.Sprintf("%s:image", CloudNS)
	InboundRules              = fmt.Sprintf("%s:inboundRules", netNS)
	InlinePolicies            = fmt.Sprintf("%s:inlinePolicies", CloudNS)
	IOPS                      = fmt.Sprintf("%s:iops", CloudNS)
	IPType                    = fmt.Sprintf("%s:ipType", netNS)
	Key                       = fmt.Sprintf("%s:key", CloudNS)
	LatestRestorableTime      = fmt.Sprintf("%s:latestRestorableTime", CloudNS)
	Launched                  = fmt.Sprintf("%s:launched", CloudNS)
	License                   = fmt.Sprintf("%s:license", CloudNS)
	Lifecycle                 = fmt.Sprintf("%s:lifecycle", CloudNS)
	LoadBalancer              = fmt.Sprintf("%s:loadBalancer", CloudNS)
	Main                      = fmt.Sprintf("%s:main", CloudNS)
	Messages                  = fmt.Sprintf("%s:messages", CloudNS)
	Modified                  = fmt.Sprintf("%s:modified", CloudNS)
	MonitoringInterval        = fmt.Sprintf("%s:monitoringInterval", CloudNS)
	MonitoringRole            = fmt.Sprintf("%s:monitoringRole", CloudNS)
	MultiAZ                   = fmt.Sprintf("%s:multiAZ", CloudNS)
	Name                      = fmt.Sprintf("%s:name", CloudNS)
	NetworkInterfaces         = fmt.Sprintf("%s:networkInterfaces", CloudNS)
	OptionGroups              = fmt.Sprintf("%s:optionGroups", CloudNS)
	OutboundRules             = fmt.Sprintf("%s:outboundRules", netNS)
	Owner                     = fmt.Sprintf("%s:owner", CloudNS)
	ParameterGroups           = fmt.Sprintf("%s:parameterGroups", CloudNS)
	PasswordLastUsed          = fmt.Sprintf("%s:passwordLastUsed", CloudNS)
	Path                      = fmt.Sprintf("%s:path", CloudNS)
	PlacementGroup            = fmt.Sprintf("%s:placementGroup", CloudNS)
	Port                      = fmt.Sprintf("%s:port", netNS)
	PortRange                 = fmt.Sprintf("%s:portRange", netNS)
	PreferredBackupDate       = fmt.Sprintf("%s:preferredBackupDate", CloudNS)
	PreferredMaintenanceDate  = fmt.Sprintf("%s:preferredMaintenanceDate", CloudNS)
	Private                   = fmt.Sprintf("%s:private", CloudNS)
	PrivateIP                 = fmt.Sprintf("%s:privateIP", netNS)
	Profile                   = fmt.Sprintf("%s:profile", CloudNS)
	Protocol                  = fmt.Sprintf("%s:protocol", netNS)
	Public                    = fmt.Sprintf("%s:public", CloudNS)
	PublicDNS                 = fmt.Sprintf("%s:publicDNS", CloudNS)
	PublicIP                  = fmt.Sprintf("%s:publicIP", netNS)
	RecordCount               = fmt.Sprintf("%s:recordCount", CloudNS)
	Records                   = fmt.Sprintf("%s:records", CloudNS)
	Region                    = fmt.Sprintf("%s:region", CloudNS)
	RootDevice                = fmt.Sprintf("%s:rootDevice", CloudNS)
	RootDeviceType            = fmt.Sprintf("%s:rootDeviceType", CloudNS)
	Routes                    = fmt.Sprintf("%s:routes", netNS)
	Scheme                    = fmt.Sprintf("%s:scheme", netNS)
	SecondaryAvailabilityZone = fmt.Sprintf("%s:secondaryAvailabilityZone", CloudNS)
	SecurityGroups            = fmt.Sprintf("%s:securityGroups", CloudNS)
	Set                       = fmt.Sprintf("%s:set", CloudNS)
	Size                      = fmt.Sprintf("%s:size", CloudNS)
	SpotInstanceRequestId     = fmt.Sprintf("%s:spotInstanceRequestId", CloudNS)
	SSHKey                    = fmt.Sprintf("%s:sshKey", CloudNS)
	State                     = fmt.Sprintf("%s:state", CloudNS)
	Storage                   = fmt.Sprintf("%s:storage", CloudNS)
	StorageType               = fmt.Sprintf("%s:storageType", CloudNS)
	Subnet                    = fmt.Sprintf("%s:subnet", CloudNS)
	Subnets                   = fmt.Sprintf("%s:subnets", CloudNS)
	Timezone                  = fmt.Sprintf("%s:timezone", CloudNS)
	Topic                     = fmt.Sprintf("%s:topic", CloudNS)
	TrafficPolicyInstance     = fmt.Sprintf("%s:trafficPolicyInstance", CloudNS)
	TTL                       = fmt.Sprintf("%s:ttl", CloudNS)
	Type                      = fmt.Sprintf("%s:type", CloudNS)
	UnhealthyThresholdCount   = fmt.Sprintf("%s:unhealthyThresholdCount", CloudNS)
	Updated                   = fmt.Sprintf("%s:updated", CloudNS)
	Username                  = fmt.Sprintf("%s:username", CloudNS)
	Vpc                       = fmt.Sprintf("%s:vpc", CloudNS)
	Vpcs                      = fmt.Sprintf("%s:vpcs", CloudNS)
	Weight                    = fmt.Sprintf("%s:weight", CloudNS)
	Zone                      = fmt.Sprintf("%s:zone", CloudNS)
)

var Labels = map[string]string{
	properties.Actions:                   Actions,
	properties.Affinity:                  Affinity,
	properties.ApproximateMessageCount:   ApproximateMessageCount,
	properties.Architecture:              Architecture,
	properties.Arn:                       Arn,
	properties.Attachable:                Attachable,
	properties.AutoUpgrade:               AutoUpgrade,
	properties.AvailabilityZone:          AvailabilityZone,
	properties.AvailabilityZones:         AvailabilityZones,
	properties.BackupRetentionPeriod:     BackupRetentionPeriod,
	properties.Bucket:                    Bucket,
	properties.CallerReference:           CallerReference,
	properties.CertificateAuthority:      CertificateAuthority,
	properties.Certificates:              Certificates,
	properties.Charset:                   Charset,
	properties.CheckHTTPCode:             CheckHTTPCode,
	properties.CheckInterval:             CheckInterval,
	properties.CheckPath:                 CheckPath,
	properties.CheckPort:                 CheckPort,
	properties.CheckProtocol:             CheckProtocol,
	properties.CheckTimeout:              CheckTimeout,
	properties.CIDR:                      CIDR,
	properties.CipherSuite:               CipherSuite,
	properties.Class:                     Class,
	properties.Cluster:                   Cluster,
	properties.Comment:                   Comment,
	properties.Continent:                 Continent,
	properties.CopyTagsToSnapshot:        CopyTagsToSnapshot,
	properties.Country:                   Country,
	properties.Created:                   Created,
	properties.DBSecurityGroups:          DBSecurityGroups,
	properties.DBSubnetGroup:             DBSubnetGroup,
	properties.Default:                   Default,
	properties.Delay:                     Delay,
	properties.Description:               Description,
	properties.Encrypted:                 Encrypted,
	properties.Endpoint:                  Endpoint,
	properties.Engine:                    Engine,
	properties.EngineVersion:             EngineVersion,
	properties.Failover:                  Failover,
	properties.Fingerprint:               Fingerprint,
	properties.GlobalID:                  GlobalID,
	properties.Grants:                    Grants,
	properties.HealthCheck:               HealthCheck,
	properties.HealthyThresholdCount:     HealthyThresholdCount,
	properties.Host:                      Host,
	properties.Hypervisor:                Hypervisor,
	properties.ID:                        ID,
	properties.Image:                     Image,
	properties.InboundRules:              InboundRules,
	properties.InlinePolicies:            InlinePolicies,
	properties.IOPS:                      IOPS,
	properties.IPType:                    IPType,
	properties.Key:                       Key,
	properties.LatestRestorableTime:      LatestRestorableTime,
	properties.Launched:                  Launched,
	properties.License:                   License,
	properties.Lifecycle:                 Lifecycle,
	properties.LoadBalancer:              LoadBalancer,
	properties.Main:                      Main,
	properties.Messages:                  Messages,
	properties.Modified:                  Modified,
	properties.MonitoringInterval:        MonitoringInterval,
	properties.MonitoringRole:            MonitoringRole,
	properties.MultiAZ:                   MultiAZ,
	properties.Name:                      Name,
	properties.NetworkInterfaces:         NetworkInterfaces,
	properties.OptionGroups:              OptionGroups,
	properties.OutboundRules:             OutboundRules,
	properties.Owner:                     Owner,
	properties.ParameterGroups:           ParameterGroups,
	properties.PasswordLastUsed:          PasswordLastUsed,
	properties.Path:                      Path,
	properties.PlacementGroup:            PlacementGroup,
	properties.Port:                      Port,
	properties.PreferredBackupDate:       PreferredBackupDate,
	properties.PreferredMaintenanceDate:  PreferredMaintenanceDate,
	properties.Private:                   Private,
	properties.PrivateIP:                 PrivateIP,
	properties.Profile:                   Profile,
	properties.Protocol:                  Protocol,
	properties.Public:                    Public,
	properties.PublicDNS:                 PublicDNS,
	properties.PublicIP:                  PublicIP,
	properties.Records:                   Records,
	properties.RecordCount:               RecordCount,
	properties.Region:                    Region,
	properties.RootDevice:                RootDevice,
	properties.RootDeviceType:            RootDeviceType,
	properties.Routes:                    Routes,
	properties.Scheme:                    Scheme,
	properties.SecondaryAvailabilityZone: SecondaryAvailabilityZone,
	properties.SecurityGroups:            SecurityGroups,
	properties.Set:                       Set,
	properties.Size:                      Size,
	properties.SpotInstanceRequestId:     SpotInstanceRequestId,
	properties.SSHKey:                    SSHKey,
	properties.State:                     State,
	properties.Storage:                   Storage,
	properties.StorageType:               StorageType,
	properties.Subnet:                    Subnet,
	properties.Subnets:                   Subnets,
	properties.Timezone:                  Timezone,
	properties.Topic:                     Topic,
	properties.TrafficPolicyInstance:     TrafficPolicyInstance,
	properties.TTL:                       TTL,
	properties.Type:                      Type,
	properties.UnhealthyThresholdCount:   UnhealthyThresholdCount,
	properties.Updated:                   Updated,
	properties.Username:                  Username,
	properties.Vpc:                       Vpc,
	properties.Vpcs:                      Vpcs,
	properties.Weight:                    Weight,
	properties.Zone:                      Zone,
}

var RdfProperties = map[string]rdfProp{
	Actions:                 {ID: Actions, RdfType: RdfProperty, RdfsLabel: properties.Actions, RdfsDefinedBy: RdfsList, RdfsDataType: XsdString},
	Affinity:                {ID: Affinity, RdfType: RdfProperty, RdfsLabel: properties.Affinity, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdString},
	ApproximateMessageCount: {ID: ApproximateMessageCount, RdfType: RdfProperty, RdfsLabel: properties.ApproximateMessageCount, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdInt},
	Architecture:            {ID: Architecture, RdfType: RdfProperty, RdfsLabel: properties.Architecture, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdString},
	Arn:                     {ID: Arn, RdfType: RdfProperty, RdfsLabel: properties.Arn, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdString},
	Attachable:              {ID: Attachable, RdfType: RdfProperty, RdfsLabel: properties.Attachable, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdBoolean},
	AutoUpgrade:             {ID: AutoUpgrade, RdfType: RdfProperty, RdfsLabel: properties.AutoUpgrade, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdBoolean},
	AvailabilityZone:        {ID: AvailabilityZone, RdfType: RdfProperty, RdfsLabel: properties.AvailabilityZone, RdfsDefinedBy: RdfsClass, RdfsDataType: XsdString},
	AvailabilityZones:       {ID: AvailabilityZones, RdfType: RdfProperty, RdfsLabel: properties.AvailabilityZones, RdfsDefinedBy: RdfsList, RdfsDataType: RdfsClass},
	BackupRetentionPeriod:   {ID: BackupRetentionPeriod, RdfType: RdfProperty, RdfsLabel: properties.BackupRetentionPeriod, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdDateTime},
	Bucket:                  {ID: Bucket, RdfType: RdfProperty, RdfsLabel: properties.Bucket, RdfsDefinedBy: RdfsClass, RdfsDataType: XsdString},
	CallerReference:         {ID: CallerReference, RdfType: RdfProperty, RdfsLabel: properties.CallerReference, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdString},
	CertificateAuthority:    {ID: CertificateAuthority, RdfType: RdfProperty, RdfsLabel: properties.CertificateAuthority, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdString},
	Certificates:            {ID: Certificates, RdfType: RdfProperty, RdfsLabel: properties.Certificates, RdfsDefinedBy: RdfsList, RdfsDataType: XsdString},
	Charset:                 {ID: Charset, RdfType: RdfProperty, RdfsLabel: properties.Charset, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdString},
	CheckHTTPCode:           {ID: CheckHTTPCode, RdfType: RdfProperty, RdfsLabel: properties.CheckHTTPCode, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdString},
	CheckInterval:           {ID: CheckInterval, RdfType: RdfProperty, RdfsLabel: properties.CheckInterval, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdInt},
	CheckPath:               {ID: CheckPath, RdfType: RdfProperty, RdfsLabel: properties.CheckPath, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdString},
	CheckPort:               {ID: CheckPort, RdfType: RdfProperty, RdfsLabel: properties.CheckPort, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdString},
	CheckProtocol:           {ID: CheckProtocol, RdfType: RdfProperty, RdfsLabel: properties.CheckProtocol, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdString},
	CheckTimeout:            {ID: CheckTimeout, RdfType: RdfProperty, RdfsLabel: properties.CheckTimeout, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdInt},
	CIDR:                    {ID: CIDR, RdfType: RdfProperty, RdfsLabel: properties.CIDR, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdString},
	CipherSuite:             {ID: CipherSuite, RdfType: RdfProperty, RdfsLabel: properties.CipherSuite, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdString},
	Class:                   {ID: Class, RdfType: RdfProperty, RdfsLabel: properties.Class, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdString},
	Cluster:                 {ID: Cluster, RdfType: RdfProperty, RdfsLabel: properties.Cluster, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdString},
	Comment:                 {ID: Comment, RdfType: RdfProperty, RdfsLabel: properties.Comment, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdString},
	Continent:               {ID: Continent, RdfType: RdfProperty, RdfsLabel: properties.Continent, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdString},
	CopyTagsToSnapshot:      {ID: CopyTagsToSnapshot, RdfType: RdfProperty, RdfsLabel: properties.CopyTagsToSnapshot, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdString},
	Country:                 {ID: Country, RdfType: RdfProperty, RdfsLabel: properties.Country, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdString},
	Created:                 {ID: Created, RdfType: RdfProperty, RdfsLabel: properties.Created, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdDateTime},
	DBSecurityGroups:        {ID: DBSecurityGroups, RdfType: RdfProperty, RdfsLabel: properties.DBSecurityGroups, RdfsDefinedBy: RdfsList, RdfsDataType: XsdString},
	DBSubnetGroup:           {ID: DBSubnetGroup, RdfType: RdfProperty, RdfsLabel: properties.DBSubnetGroup, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdString},
	Default:                 {ID: Default, RdfType: RdfProperty, RdfsLabel: properties.Default, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdBoolean},
	Delay:                   {ID: Delay, RdfType: RdfProperty, RdfsLabel: properties.Delay, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdInt},
	Description:             {ID: Description, RdfType: RdfProperty, RdfsLabel: properties.Description, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdString},
	Encrypted:               {ID: Encrypted, RdfType: RdfProperty, RdfsLabel: properties.Encrypted, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdString},
	Endpoint:                {ID: Endpoint, RdfType: RdfProperty, RdfsLabel: properties.Endpoint, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdString},
	Engine:                  {ID: Engine, RdfType: RdfProperty, RdfsLabel: properties.Engine, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdString},
	EngineVersion:           {ID: EngineVersion, RdfType: RdfProperty, RdfsLabel: properties.EngineVersion, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdString},
	Failover:                {ID: Failover, RdfType: RdfProperty, RdfsLabel: properties.Failover, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdString},
	Fingerprint:             {ID: Fingerprint, RdfType: RdfProperty, RdfsLabel: properties.Fingerprint, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdString},
	GlobalID:                {ID: GlobalID, RdfType: RdfProperty, RdfsLabel: properties.GlobalID, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdString},
	Grants:                  {ID: Grants, RdfType: RdfProperty, RdfsLabel: properties.Grants, RdfsDefinedBy: RdfsList, RdfsDataType: Grant},
	HealthCheck:             {ID: HealthCheck, RdfType: RdfProperty, RdfsLabel: properties.HealthCheck, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdString},
	HealthyThresholdCount:   {ID: HealthyThresholdCount, RdfType: RdfProperty, RdfsLabel: properties.HealthyThresholdCount, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdInt},
	Host:                     {ID: Host, RdfType: RdfProperty, RdfsLabel: properties.Host, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdString},
	Hypervisor:               {ID: Hypervisor, RdfType: RdfProperty, RdfsLabel: properties.Hypervisor, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdString},
	ID:                       {ID: ID, RdfType: RdfProperty, RdfsLabel: properties.ID, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdString},
	Image:                    {ID: Image, RdfType: RdfProperty, RdfsLabel: properties.Image, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdString},
	InboundRules:             {ID: InboundRules, RdfType: RdfProperty, RdfsLabel: properties.InboundRules, RdfsDefinedBy: RdfsList, RdfsDataType: NetFirewallRule},
	InlinePolicies:           {ID: InlinePolicies, RdfType: RdfProperty, RdfsLabel: properties.InlinePolicies, RdfsDefinedBy: RdfsList, RdfsDataType: RdfsClass},
	IOPS:                     {ID: IOPS, RdfType: RdfProperty, RdfsLabel: properties.IOPS, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdInt},
	IPType:                   {ID: IPType, RdfType: RdfProperty, RdfsLabel: properties.IPType, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdString},
	Key:                      {ID: Key, RdfType: RdfProperty, RdfsLabel: properties.Key, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdString},
	LatestRestorableTime:     {ID: LatestRestorableTime, RdfType: RdfProperty, RdfsLabel: properties.LatestRestorableTime, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdDateTime},
	Launched:                 {ID: Launched, RdfType: RdfProperty, RdfsLabel: properties.Launched, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdDateTime},
	License:                  {ID: License, RdfType: RdfProperty, RdfsLabel: properties.License, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdString},
	Lifecycle:                {ID: Lifecycle, RdfType: RdfProperty, RdfsLabel: properties.Lifecycle, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdString},
	LoadBalancer:             {ID: LoadBalancer, RdfType: RdfProperty, RdfsLabel: properties.LoadBalancer, RdfsDefinedBy: RdfsClass, RdfsDataType: XsdString},
	Main:                     {ID: Main, RdfType: RdfProperty, RdfsLabel: properties.Main, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdBoolean},
	Messages:                 {ID: Messages, RdfType: RdfProperty, RdfsLabel: properties.Messages, RdfsDefinedBy: RdfsList, RdfsDataType: XsdString},
	Modified:                 {ID: Modified, RdfType: RdfProperty, RdfsLabel: properties.Modified, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdDateTime},
	MonitoringInterval:       {ID: MonitoringInterval, RdfType: RdfProperty, RdfsLabel: properties.MonitoringInterval, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdString},
	MonitoringRole:           {ID: MonitoringRole, RdfType: RdfProperty, RdfsLabel: properties.MonitoringRole, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdString},
	MultiAZ:                  {ID: MultiAZ, RdfType: RdfProperty, RdfsLabel: properties.MultiAZ, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdString},
	Name:                     {ID: Name, RdfType: RdfProperty, RdfsLabel: properties.Name, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdString},
	NetworkInterfaces:        {ID: NetworkInterfaces, RdfType: RdfProperty, RdfsLabel: properties.NetworkInterfaces, RdfsDefinedBy: RdfsList, RdfsDataType: XsdString},
	OptionGroups:             {ID: OptionGroups, RdfType: RdfProperty, RdfsLabel: properties.OptionGroups, RdfsDefinedBy: RdfsList, RdfsDataType: XsdString},
	OutboundRules:            {ID: OutboundRules, RdfType: RdfProperty, RdfsLabel: properties.OutboundRules, RdfsDefinedBy: RdfsList, RdfsDataType: NetFirewallRule},
	Owner:                    {ID: Owner, RdfType: RdfProperty, RdfsLabel: properties.Owner, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdString},
	ParameterGroups:          {ID: ParameterGroups, RdfType: RdfProperty, RdfsLabel: properties.ParameterGroups, RdfsDefinedBy: RdfsList, RdfsDataType: XsdString},
	PasswordLastUsed:         {ID: PasswordLastUsed, RdfType: RdfProperty, RdfsLabel: properties.PasswordLastUsed, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdDateTime},
	Path:                     {ID: Path, RdfType: RdfProperty, RdfsLabel: properties.Path, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdString},
	PlacementGroup:           {ID: PlacementGroup, RdfType: RdfProperty, RdfsLabel: properties.PlacementGroup, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdString},
	Port:                     {ID: Port, RdfType: RdfProperty, RdfsLabel: properties.Port, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdInt},
	PreferredBackupDate:      {ID: PreferredBackupDate, RdfType: RdfProperty, RdfsLabel: properties.PreferredBackupDate, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdString},
	PreferredMaintenanceDate: {ID: PreferredMaintenanceDate, RdfType: RdfProperty, RdfsLabel: properties.PreferredMaintenanceDate, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdString},
	Private:                  {ID: Private, RdfType: RdfProperty, RdfsLabel: properties.Private, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdString},
	PrivateIP:                {ID: PrivateIP, RdfType: RdfProperty, RdfsLabel: properties.PrivateIP, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdString},
	Profile:                  {ID: Profile, RdfType: RdfProperty, RdfsLabel: properties.Profile, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdString},
	Protocol:                 {ID: Protocol, RdfType: RdfProperty, RdfsLabel: properties.Protocol, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdString},
	Public:                   {ID: Public, RdfType: RdfProperty, RdfsLabel: properties.Public, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdBoolean},
	PublicDNS:                {ID: PublicDNS, RdfType: RdfProperty, RdfsLabel: properties.PublicDNS, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdString},
	PublicIP:                 {ID: PublicIP, RdfType: RdfProperty, RdfsLabel: properties.PublicIP, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdString},
	Records:                  {ID: Records, RdfType: RdfProperty, RdfsLabel: properties.Records, RdfsDefinedBy: RdfsList, RdfsDataType: RdfsClass},
	RecordCount:              {ID: RecordCount, RdfType: RdfProperty, RdfsLabel: properties.RecordCount, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdInt},
	Region:                   {ID: Region, RdfType: RdfProperty, RdfsLabel: properties.Region, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdString},
	RootDevice:               {ID: RootDevice, RdfType: RdfProperty, RdfsLabel: properties.RootDevice, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdString},
	RootDeviceType:           {ID: RootDeviceType, RdfType: RdfProperty, RdfsLabel: properties.RootDeviceType, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdString},
	Routes:                   {ID: Routes, RdfType: RdfProperty, RdfsLabel: properties.Routes, RdfsDefinedBy: RdfsList, RdfsDataType: NetRoute},
	Scheme:                   {ID: Scheme, RdfType: RdfProperty, RdfsLabel: properties.Scheme, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdString},
	SecondaryAvailabilityZone: {ID: SecondaryAvailabilityZone, RdfType: RdfProperty, RdfsLabel: properties.SecondaryAvailabilityZone, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdString},
	SecurityGroups:            {ID: SecurityGroups, RdfType: RdfProperty, RdfsLabel: properties.SecurityGroups, RdfsDefinedBy: RdfsList, RdfsDataType: RdfsClass},
	Set:                       {ID: Set, RdfType: RdfProperty, RdfsLabel: properties.Set, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdString},
	Size:                      {ID: Size, RdfType: RdfProperty, RdfsLabel: properties.Size, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdInt},
	SpotInstanceRequestId: {ID: SpotInstanceRequestId, RdfType: RdfProperty, RdfsLabel: properties.SpotInstanceRequestId, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdString},
	SSHKey:                {ID: SSHKey, RdfType: RdfProperty, RdfsLabel: properties.SSHKey, RdfsDefinedBy: RdfsClass, RdfsDataType: XsdString},
	State:                 {ID: State, RdfType: RdfProperty, RdfsLabel: properties.State, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdString},
	Storage:               {ID: Storage, RdfType: RdfProperty, RdfsLabel: properties.Storage, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdString},
	StorageType:           {ID: StorageType, RdfType: RdfProperty, RdfsLabel: properties.StorageType, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdString},
	Subnet:                {ID: Subnet, RdfType: RdfProperty, RdfsLabel: properties.Subnet, RdfsDefinedBy: RdfsClass, RdfsDataType: XsdString},
	Subnets:               {ID: Subnets, RdfType: RdfProperty, RdfsLabel: properties.Subnets, RdfsDefinedBy: RdfsList, RdfsDataType: RdfsClass},
	Timezone:              {ID: Timezone, RdfType: RdfProperty, RdfsLabel: properties.Timezone, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdString},
	Topic:                 {ID: Topic, RdfType: RdfProperty, RdfsLabel: properties.Topic, RdfsDefinedBy: RdfsClass, RdfsDataType: XsdString},
	TrafficPolicyInstance: {ID: TrafficPolicyInstance, RdfType: RdfProperty, RdfsLabel: properties.TrafficPolicyInstance, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdString},
	TTL:  {ID: TTL, RdfType: RdfProperty, RdfsLabel: properties.TTL, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdInt},
	Type: {ID: Type, RdfType: RdfProperty, RdfsLabel: properties.Type, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdString},
	UnhealthyThresholdCount: {ID: UnhealthyThresholdCount, RdfType: RdfProperty, RdfsLabel: properties.UnhealthyThresholdCount, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdInt},
	Updated:                 {ID: Updated, RdfType: RdfProperty, RdfsLabel: properties.Updated, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdString},
	Username:                {ID: Username, RdfType: RdfProperty, RdfsLabel: properties.Username, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdString},
	Vpc:                     {ID: Vpc, RdfType: RdfProperty, RdfsLabel: properties.Vpc, RdfsDefinedBy: RdfsClass, RdfsDataType: XsdString},
	Vpcs:                    {ID: Vpcs, RdfType: RdfProperty, RdfsLabel: properties.Vpcs, RdfsDefinedBy: RdfsList, RdfsDataType: RdfsClass},
	Weight:                  {ID: Weight, RdfType: RdfProperty, RdfsLabel: properties.Weight, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdString},
	Zone:                    {ID: Zone, RdfType: RdfProperty, RdfsLabel: properties.Zone, RdfsDefinedBy: RdfsClass, RdfsDataType: XsdString},

	//Subproperties
	PortRange: {ID: PortRange, RdfType: RdfsSubProperty, RdfsDefinedBy: RdfsLiteral, RdfsDataType: XsdString},
}
