package aws

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// Image resolving allows to find AWS AMIs identifiers specifying what you want instead
// of an id that is specific to a region. The ami query string specification is as follows:
//
// OWNER:DISTRO[VARIANT]:ARCH:VIRTUALIZATION:STORE
// (with everything optional expect for the OWNER)
//
// As for now only the main specific owner are taken into account
// and we deal with bare machine only distribution. Here are some examples:
//
// - canonical:ubuntu[trusty]
//
// - redhat:rhel[6.8]
//
// - amazonlinux
//
// - suselinux:[sles-12]
//
// - canonical::i386
//
// - redhat::::instance-store
//
// The default values are: Arch="x86_64", Virt="hvm", Store="ebs"
type ImageResolver struct {
	InfraService *Infra
}

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

func (r *ImageResolver) Resolve(query string) ([]*AwsImage, error) {
	var results []*AwsImage

	q, err := parseQuery(query)
	if err != nil {
		return results, err
	}

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

	amis, err := r.InfraService.DescribeImages(params)
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

	supportedPlatforms []string
)

func init() {
	for name := range Platforms {
		supportedPlatforms = append(supportedPlatforms, name)
	}
}

type imageQuery struct {
	Platform Platform
	Distro   Distro
}

type Distro struct {
	Name, Variant, Arch, Virt, Store string
}

type Platform struct {
	Name          string
	Id            string
	DistroName    string
	LatestVariant string
	MatchFunc     func(s string, d Distro) bool
}

var distroVariant = regexp.MustCompile(`\[([^]]*)\]`)

func parseQuery(s string) (imageQuery, error) {
	splits := strings.SplitN(s, ":", 5)
	if len(splits) < 1 {
		return imageQuery{}, fmt.Errorf("invalid image query: must contains at least the owner name: %s ", strings.Join(supportedPlatforms, ","))
	}

	for i, s := range splits {
		splits[i] = strings.ToLower(s)
	}

	q := imageQuery{}

	plat, ok := Platforms[splits[0]]
	if !ok {
		return imageQuery{}, fmt.Errorf("unsupported owner/platform %s", splits[0])
	}

	q.Platform = plat
	q.Distro = Distro{Variant: q.Platform.LatestVariant, Name: q.Platform.DistroName}

	if len(splits) > 1 {
		distro := splits[1]
		matches := distroVariant.FindStringSubmatch(distro)
		if len(matches) == 2 {
			index := strings.Index(distro, "[")
			if index != -1 {
				q.Distro.Name = distro[:index]
			}
			q.Distro.Variant = matches[1]
		} else {
			q.Distro.Name = distro
		}
	}

	if len(splits) > 2 && strings.TrimSpace(splits[2]) != "" {
		q.Distro.Arch = splits[2]
	} else {
		q.Distro.Arch = defaultArch
	}

	if len(splits) > 3 && strings.TrimSpace(splits[3]) != "" {
		q.Distro.Virt = splits[3]
	} else {
		q.Distro.Virt = defaultVirt
	}

	if len(splits) > 4 && strings.TrimSpace(splits[4]) != "" {
		q.Distro.Store = splits[4]
	} else {
		q.Distro.Store = defaultStore
	}

	return q, nil
}
