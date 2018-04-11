package awsspec

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"sync"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/wallix/awless/aws/doc"
)

// Image resolving allows to find AWS AMIs identifiers specifying what you want instead
// of an id that is specific to a region. The ami query string specification is as follows:
//
// owner:distro:variant:arch:virtualization:store
//
// Everything optional expect for the owner.
//
// As for now only the main specific owner are taken into account
// and we deal with bares machines only distribution. Here are some examples:
//
// - canonical:ubuntu:trusty
//
// - redhat:rhel:6.8
//
// - redhat::6.8
//
// - amazonlinux
//
// - suselinux:sles-12
//
// - canonical:::i386
//
// - redhat::::instance-store
//
// - centos:centos
//
// The default values are: Arch="x86_64", Virt="hvm", Store="ebs"
type ImageResolver func(*ec2.DescribeImagesInput) (*ec2.DescribeImagesOutput, error)

func EC2ImageResolver() ImageResolver {
	factory := CommandFactory.Build("createinstance")
	return ImageResolver(factory().(*CreateInstance).api.DescribeImages)
}

var DefaultImageResolverCache = new(ImageResolverCache)

type ImageResolverCache struct {
	mu    sync.Mutex
	cache map[string][]*AwsImage
}

func (r *ImageResolverCache) Store(key string, images []*AwsImage) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.cache == nil {
		r.cache = make(map[string][]*AwsImage)
	}
	r.cache[key] = images
}

func (r *ImageResolverCache) Get(key string) ([]*AwsImage, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.cache == nil {
		r.cache = make(map[string][]*AwsImage)
	}
	images, ok := r.cache[key]
	return images, ok
}

const ImageQuerySpec = "owner:distro:variant:arch:virtualization:store"

type AwsImage struct {
	Id                 string
	Owner              string
	Location           string
	Type               string
	Architecture       string
	VirtualizationType string
	Name               string
	Created            time.Time
	Hypervisor         string
	Store              string
}

type ImageQuery struct {
	Platform Platform
	Distro   Distro
}

func (q ImageQuery) String() string {
	var all []string
	all = append(all, q.Platform.Name)
	all = append(all, q.Distro.Name)
	all = append(all, q.Distro.Variant)
	all = append(all, q.Distro.Arch)
	all = append(all, q.Distro.Virt)
	all = append(all, q.Distro.Store)
	return strings.Join(all, ":")
}

type Distro struct {
	Name, Variant, Arch, Virt, Store string
}

var (
	validArchs  = []string{"i386", "x86_64"}
	validVirts  = []string{"paravirtual", "hvm"}
	validStores = []string{"ebs", "instance-store"}
)

type Platform struct {
	Name          string
	Id            string
	DistroName    string
	LatestVariant string
	MatchFunc     func(s string, d Distro) bool
}

func (resolv ImageResolver) Resolve(q ImageQuery) ([]*AwsImage, bool, error) {
	images, found := DefaultImageResolverCache.Get(q.String())
	if found {
		return images, true, nil
	}

	results := make([]*AwsImage, 0) // json empty array friendly

	var filters []*ec2.Filter
	filters = append(filters,
		&ec2.Filter{
			Name:   awssdk.String("state"),
			Values: []*string{awssdk.String("available")},
		},
		&ec2.Filter{
			Name:   awssdk.String("is-public"),
			Values: []*string{awssdk.String("true")},
		},
	)

	filters = append(filters,
		&ec2.Filter{
			Name:   awssdk.String("owner-id"),
			Values: []*string{awssdk.String(q.Platform.Id)},
		},
	)

	filters = append(filters,
		&ec2.Filter{
			Name:   awssdk.String("virtualization-type"),
			Values: []*string{awssdk.String(q.Distro.Virt)},
		},
	)

	filters = append(filters,
		&ec2.Filter{
			Name:   awssdk.String("architecture"),
			Values: []*string{awssdk.String(q.Distro.Arch)},
		},
	)

	filters = append(filters,
		&ec2.Filter{
			Name:   awssdk.String("root-device-type"),
			Values: []*string{awssdk.String(q.Distro.Store)},
		},
	)

	params := &ec2.DescribeImagesInput{
		ExecutableUsers: []*string{awssdk.String("all")},
		Filters:         filters,
	}

	amis, err := resolv(params)
	if err != nil {
		return results, false, err
	}

	for _, ami := range amis.Images {
		if !q.Platform.MatchFunc(strings.ToLower(awssdk.StringValue(ami.Name)), q.Distro) {
			continue
		}

		img := &AwsImage{
			Id:                 awssdk.StringValue(ami.ImageId),
			Owner:              q.Platform.Id,
			Location:           awssdk.StringValue(ami.ImageLocation),
			Type:               awssdk.StringValue(ami.ImageType),
			Architecture:       awssdk.StringValue(ami.Architecture),
			VirtualizationType: awssdk.StringValue(ami.VirtualizationType),
			Name:               awssdk.StringValue(ami.Name),
			Hypervisor:         awssdk.StringValue(ami.Hypervisor),
			Store:              awssdk.StringValue(ami.RootDeviceType),
		}

		img.Created, _ = time.Parse(time.RFC3339, awssdk.StringValue(ami.CreationDate))

		results = append(results, img)
	}

	sort.Slice(results, func(i, j int) bool { return results[i].Created.After(results[j].Created) })

	DefaultImageResolverCache.Store(q.String(), results)

	return results, false, nil
}

var (
	Platforms = map[string]Platform{
		"canonical":   Canonical,
		"redhat":      RedHat,
		"debian":      Debian,
		"amazonlinux": AmazonLinux,
		"coreos":      CoreOS,
		"centos":      CentOS,
		"suselinux":   SuseLinux,
		"windows":     Windows,
	}

	Canonical = Platform{
		Name: "canonical", Id: "099720109477", DistroName: "ubuntu", LatestVariant: "xenial",
		MatchFunc: func(s string, d Distro) bool {
			return strings.HasPrefix(s, fmt.Sprintf("%s/images/%s-ssd/%s-%s", d.Name, d.Virt, d.Name, d.Variant))
		},
	}

	RedHat = Platform{
		Name: "redhat", Id: "309956199498", DistroName: "rhel", LatestVariant: "7.5",
		MatchFunc: func(s string, d Distro) bool {
			return strings.Contains(s, fmt.Sprintf("%s-%s", d.Name, d.Variant))
		},
	}

	Debian = Platform{
		Name: "debian", Id: "379101102735", DistroName: "debian", LatestVariant: "stretch",
		MatchFunc: func(s string, d Distro) bool {
			return strings.HasPrefix(s, fmt.Sprintf("%s-%s", d.Name, d.Variant))
		},
	}

	CoreOS = Platform{
		Name: "coreos", Id: "595879546273", DistroName: "coreos", LatestVariant: "1688",
		MatchFunc: func(s string, d Distro) bool {
			return strings.HasPrefix(s, strings.ToLower(fmt.Sprintf("%s-stable-%s", d.Name, d.Variant)))
		},
	}

	CentOS = Platform{
		Name: "centos", Id: "679593333241", DistroName: "centos", LatestVariant: "7",
		MatchFunc: func(s string, d Distro) bool {
			return strings.HasPrefix(s, strings.ToLower(fmt.Sprintf("%s Linux %s", d.Name, d.Variant)))
		},
	}

	AmazonLinux = Platform{
		Name: "amazonlinux", Id: "137112412989", DistroName: "amzn", LatestVariant: "hvm",
		MatchFunc: func(s string, d Distro) bool {
			return strings.HasPrefix(s, fmt.Sprintf("%s-ami-%s", d.Name, d.Variant))
		},
	}

	SuseLinux = Platform{
		Name: "suselinux", Id: "013907871322", LatestVariant: "sles-12",
		MatchFunc: func(s string, d Distro) bool {
			return strings.HasPrefix(s, fmt.Sprintf("suse-%s", d.Variant))
		},
	}

	Windows = Platform{
		Name: "windows", Id: "801119661308", DistroName: "server", LatestVariant: "2016",
		MatchFunc: func(s string, d Distro) bool {
			return strings.HasPrefix(s, strings.ToLower(fmt.Sprintf("windows_%s-%s-english", d.Name, d.Variant)))
		},
	}

	defaultArch  = "x86_64"
	defaultVirt  = "hvm"
	defaultStore = "ebs"

	SupportedAMIOwners []string
)

func init() {
	for name := range Platforms {
		SupportedAMIOwners = append(SupportedAMIOwners, name)
	}
	awsdoc.CommandDefinitionsDoc["create.instance"] = fmt.Sprintf("Create an EC2 instance.\n\nThe `distro` param allows to resolve from the current region the official community free bare AMI according to an awless specific bare distro query format, ordering by latest first. The query string specification is the following column separated format:\n\n\t\t%s\n\nIn this query format, everything is optional expect for the 'owner'. Supported owners: %s", ImageQuerySpec, strings.Join(SupportedAMIOwners, ", "))
}

func ParseImageQuery(s string) (ImageQuery, error) {
	supported := strings.Join(SupportedAMIOwners, ", ")
	splits := strings.Split(s, ":")

	splitsCount := len(splits)

	q := ImageQuery{}

	if splitsCount < 1 {
		return q, fmt.Errorf("malformed image query '%s': missing at least one supported owner name: %s", s, supported)
	}

	if splitsCount > 6 {
		return q, fmt.Errorf("malformed image query '%s': to many tokens, expecting format: %s", s, ImageQuerySpec)
	}

	for i, s := range splits {
		splits[i] = strings.ToLower(s)
	}

	plat, ok := Platforms[splits[0]]
	if !ok {
		return q, fmt.Errorf("unsupported owner '%s'. Expecting: %s (see awless search images -h for more)", splits[0], supported)
	}

	q.Platform = plat

	if splitsCount > 1 && strings.TrimSpace(splits[1]) != "" {
		q.Distro.Name = splits[1]
	} else {
		q.Distro.Name = q.Platform.DistroName
	}

	if splitsCount > 2 && strings.TrimSpace(splits[2]) != "" {
		q.Distro.Variant = splits[2]
	} else {
		q.Distro.Variant = q.Platform.LatestVariant
	}

	if splitsCount > 3 && strings.TrimSpace(splits[3]) != "" {
		if !contains(validArchs, splits[3]) {
			return q, fmt.Errorf("image query: invalid architecture '%s' (expecting: %s)", splits[3], strings.Join(validArchs, ", "))
		}
		q.Distro.Arch = splits[3]
	} else {
		q.Distro.Arch = defaultArch
	}

	if splitsCount > 4 && strings.TrimSpace(splits[4]) != "" {
		if !contains(validVirts, splits[4]) {
			return q, fmt.Errorf("image query: invalid virtualization '%s' (expecting: %s)", splits[4], strings.Join(validVirts, ", "))
		}
		q.Distro.Virt = splits[4]
	} else {
		q.Distro.Virt = defaultVirt
	}

	if splitsCount > 5 && strings.TrimSpace(splits[5]) != "" {
		if !contains(validStores, splits[5]) {
			return q, fmt.Errorf("image query: invalid store '%s' (expecting: %s)", splits[5], strings.Join(validStores, ", "))
		}
		q.Distro.Store = splits[5]
	} else {
		q.Distro.Store = defaultStore
	}

	return q, nil
}
