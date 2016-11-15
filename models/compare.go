package models

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/wallix/awless/utils"
)

type Diff struct {
	RegionMatch bool

	ExtraVpcs   []string
	MissingVpcs []string

	ExtraSubnets   []string
	MissingSubnets []string

	ExtraInstances   []string
	MissingInstances []string
}

func (d *Diff) Json() []byte {
	content, err := json.MarshalIndent(d, "", " ")
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(-1)
	}

	return content
}

func Compare(local *Region, remote *Region) *Diff {
	diff := &Diff{RegionMatch: true}

	if local.Id != remote.Id {
		diff.RegionMatch = false
		return diff
	}

	diff.ExtraVpcs = utils.Substraction(local.Vpcs, remote.Vpcs)
	diff.MissingVpcs = utils.Substraction(remote.Vpcs, local.Vpcs)

	localSubnets, localInstances := local.AllSubnetsAndInstances()
	remoteSubnets, remoteInstances := remote.AllSubnetsAndInstances()

	diff.ExtraSubnets = utils.Substraction(localSubnets, remoteSubnets)
	diff.MissingSubnets = utils.Substraction(remoteSubnets, localSubnets)

	diff.ExtraInstances = utils.Substraction(localInstances, remoteInstances)
	diff.MissingInstances = utils.Substraction(remoteInstances, localInstances)

	return diff
}
