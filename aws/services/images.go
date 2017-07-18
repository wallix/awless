package awsservices

import (
	"fmt"
	"sort"
	"strings"
	"time"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
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
// The default values are: Arch="x86_64", Virt="hvm", Store="ebs"
type ImageResolver struct {
	InfraService *Infra
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

func (r *ImageResolver) Resolve(q ImageQuery) ([]*AwsImage, error) {
	results := make([]*AwsImage, 0) // json empty array friendly

	filters := []*ec2.Filter{}

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

	amis, err := r.InfraService.EC2API.DescribeImages(params)
	if err != nil {
		return results, err
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

	return results, nil
}

var (
	Platforms = map[string]Platform{
		"canonical":       Canonical,
		"redhat":          RedHat,
		"debian":          Debian,
		"amazonlinux":     AmazonLinux,
		"suselinux":       SuseLinux,
		"microsoftserver": MicrosoftServer,
	}

	Canonical = Platform{
		Name: "canonical", Id: "099720109477", DistroName: "ubuntu", LatestVariant: "xenial",
		MatchFunc: func(s string, d Distro) bool {
			return strings.HasPrefix(s, fmt.Sprintf("%s/images/%s-ssd/%s-%s", d.Name, d.Virt, d.Name, d.Variant))
		},
	}

	RedHat = Platform{
		Name: "redhat", Id: "309956199498", DistroName: "rhel", LatestVariant: "7.3",
		MatchFunc: func(s string, d Distro) bool {
			return strings.Contains(s, fmt.Sprintf("%s-%s", d.Name, d.Variant))
		},
	}

	Debian = Platform{
		Name: "debian", Id: "379101102735", DistroName: "debian", LatestVariant: "jessie",
		MatchFunc: func(s string, d Distro) bool {
			return strings.HasPrefix(s, fmt.Sprintf("%s-%s", d.Name, d.Variant))
		},
	}

	AmazonLinux = Platform{
		Name: "amazonlinux", Id: "137112412989", LatestVariant: "hvm",
		MatchFunc: func(s string, d Distro) bool {
			return strings.HasPrefix(s, fmt.Sprintf("amzn-ami-%s", d.Variant))
		},
	}

	SuseLinux = Platform{
		Name: "suselinux", Id: "013907871322", LatestVariant: "sles-12",
		MatchFunc: func(s string, d Distro) bool {
			return strings.HasPrefix(s, fmt.Sprintf("suse-%s", d.Variant))
		},
	}

	MicrosoftServer = Platform{
		Name: "microsoftserver", Id: "801119661308", LatestVariant: "Windows_Server-2016",
		MatchFunc: func(s string, d Distro) bool {
			return strings.HasPrefix(s, fmt.Sprintf("%s", d.Variant))
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
		return q, fmt.Errorf("unsupported owner %s. Expecting: %s", splits[0], supported)
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

func contains(arr []string, s string) bool {
	for _, e := range arr {
		if e == s {
			return true
		}
	}
	return false
}
