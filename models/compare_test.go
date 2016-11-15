package models

import "testing"

func TestCompareRegion(t *testing.T) {
	region1 := &Region{Id: "michigan"}
	region2 := &Region{Id: "wisconsin"}

	diff := Compare(region1, region2)
	if got, want := diff.RegionMatch, false; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}

	region1 = &Region{Id: "michigan"}
	region2 = &Region{Id: "michigan"}

	diff = Compare(region1, region2)
	if got, want := diff.RegionMatch, true; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
}

func TestCompareVpcs(t *testing.T) {
	region1 := &Region{
		Id: "michigan",
		Vpcs: []*Vpc{
			&Vpc{Id: "vpc_1"},
		},
	}

	region2 := &Region{
		Id: "michigan",
		Vpcs: []*Vpc{
			&Vpc{Id: "vpc_1"},
			&Vpc{Id: "vpc_2"},
		},
	}

	diff := Compare(region1, region2)
	compareSlice(t, diff.MissingVpcs, []string{"vpc_2"})
	if got, want := len(diff.ExtraVpcs), 0; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}

	region1 = &Region{
		Id: "michigan",
		Vpcs: []*Vpc{
			&Vpc{Id: "vpc_1"},
			&Vpc{Id: "vpc_2"},
		},
	}

	region2 = &Region{
		Id: "michigan",
		Vpcs: []*Vpc{
			&Vpc{Id: "vpc_1"},
			&Vpc{Id: "vpc_3"},
		},
	}

	diff = Compare(region1, region2)
	compareSlice(t, diff.MissingVpcs, []string{"vpc_3"})
	compareSlice(t, diff.ExtraVpcs, []string{"vpc_2"})

	region1 = &Region{
		Id: "michigan",
		Vpcs: []*Vpc{
			&Vpc{Id: "vpc_1"},
			&Vpc{Id: "vpc_2"},
		},
	}

	region2 = &Region{
		Id: "michigan",
		Vpcs: []*Vpc{
			&Vpc{Id: "vpc_1"},
			&Vpc{Id: "vpc_2"},
		},
	}

	diff = Compare(region1, region2)
	if got, want := len(diff.MissingVpcs), 0; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := len(diff.ExtraVpcs), 0; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
}

func TestCompareSubnets(t *testing.T) {
	region1 := &Region{
		Id: "michigan",
		Vpcs: []*Vpc{
			&Vpc{Id: "vpc_1", Subnets: []*Subnet{
				&Subnet{Id: "sub_1", VpcId: "vpc_1"},
			}},
		},
	}

	region2 := &Region{
		Id: "michigan",
		Vpcs: []*Vpc{
			&Vpc{Id: "vpc_1", Subnets: []*Subnet{
				&Subnet{Id: "sub_1", VpcId: "vpc_1"},
				&Subnet{Id: "sub_2", VpcId: "vpc_1"},
			}},
		},
	}

	diff := Compare(region1, region2)
	compareSlice(t, diff.MissingSubnets, []string{"sub_2"})
	if got, want := len(diff.ExtraSubnets), 0; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}

	region1 = &Region{
		Id: "michigan",
		Vpcs: []*Vpc{
			&Vpc{Id: "vpc_1", Subnets: []*Subnet{
				&Subnet{Id: "sub_1", VpcId: "vpc_1"},
				&Subnet{Id: "sub_2", VpcId: "vpc_1"},
			}},
		},
	}

	region2 = &Region{
		Id: "michigan",
		Vpcs: []*Vpc{
			&Vpc{Id: "vpc_1", Subnets: []*Subnet{
				&Subnet{Id: "sub_1", VpcId: "vpc_1"},
				&Subnet{Id: "sub_3", VpcId: "vpc_1"},
			}},
		},
	}

	diff = Compare(region1, region2)
	compareSlice(t, diff.MissingSubnets, []string{"sub_3"})
	compareSlice(t, diff.ExtraSubnets, []string{"sub_2"})

	region1 = &Region{
		Id: "michigan",
		Vpcs: []*Vpc{
			&Vpc{Id: "vpc_1", Subnets: []*Subnet{
				&Subnet{Id: "sub_1", VpcId: "vpc_1"},
				&Subnet{Id: "sub_2", VpcId: "vpc_1"},
			}},
			&Vpc{Id: "vpc_2", Subnets: []*Subnet{
				&Subnet{Id: "vpc_2_sub_1", VpcId: "vpc_2"},
				&Subnet{Id: "vpc_2_sub_2", VpcId: "vpc_2"},
			}},
		},
	}

	region2 = &Region{
		Id: "michigan",
		Vpcs: []*Vpc{
			&Vpc{Id: "vpc_1", Subnets: []*Subnet{
				&Subnet{Id: "sub_1", VpcId: "vpc_1"},
				&Subnet{Id: "sub_2", VpcId: "vpc_1"},
			}},
			&Vpc{Id: "vpc_2", Subnets: []*Subnet{
				&Subnet{Id: "vpc_2_sub_1", VpcId: "vpc_2"},
				&Subnet{Id: "vpc_2_sub_3", VpcId: "vpc_2"},
			}},
		},
	}

	diff = Compare(region1, region2)
	compareSlice(t, diff.MissingSubnets, []string{"vpc_2_sub_3"})
	compareSlice(t, diff.ExtraSubnets, []string{"vpc_2_sub_2"})

	region1 = &Region{
		Id: "michigan",
		Vpcs: []*Vpc{
			&Vpc{Id: "vpc_1", Subnets: []*Subnet{
				&Subnet{Id: "sub_1", VpcId: "vpc_1"},
				&Subnet{Id: "sub_2", VpcId: "vpc_1"},
			}},
		},
	}

	region2 = &Region{
		Id: "michigan",
		Vpcs: []*Vpc{
			&Vpc{Id: "vpc_1", Subnets: []*Subnet{
				&Subnet{Id: "sub_1", VpcId: "vpc_1"},
				&Subnet{Id: "sub_2", VpcId: "vpc_1"},
			}},
		},
	}

	diff = Compare(region1, region2)
	if got, want := len(diff.MissingSubnets), 0; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := len(diff.ExtraSubnets), 0; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
}

func TestCompareInstances(t *testing.T) {
	region1 := &Region{
		Id: "michigan",
		Vpcs: []*Vpc{
			&Vpc{Id: "vpc_1", Subnets: []*Subnet{
				&Subnet{Id: "sub_1", Instances: []*Instance{
					&Instance{Id: "inst_1"},
				}},
			}},
		},
	}

	region2 := &Region{
		Id: "michigan",
		Vpcs: []*Vpc{
			&Vpc{Id: "vpc_1", Subnets: []*Subnet{
				&Subnet{Id: "sub_1", Instances: []*Instance{
					&Instance{Id: "inst_1"},
					&Instance{Id: "inst_2"},
				}},
			}},
		},
	}

	diff := Compare(region1, region2)
	compareSlice(t, diff.MissingInstances, []string{"inst_2"})
	if got, want := len(diff.ExtraInstances), 0; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}

	region1 = &Region{
		Id: "michigan",
		Vpcs: []*Vpc{
			&Vpc{Id: "vpc_1", Subnets: []*Subnet{
				&Subnet{Id: "sub_1", Instances: []*Instance{
					&Instance{Id: "inst_1"},
					&Instance{Id: "inst_2"},
				}},
			}},
		},
	}

	region2 = &Region{
		Id: "michigan",
		Vpcs: []*Vpc{
			&Vpc{Id: "vpc_1", Subnets: []*Subnet{
				&Subnet{Id: "sub_1", Instances: []*Instance{
					&Instance{Id: "inst_1"},
					&Instance{Id: "inst_3"},
				}},
			}},
		},
	}

	diff = Compare(region1, region2)
	compareSlice(t, diff.ExtraInstances, []string{"inst_2"})
	compareSlice(t, diff.MissingInstances, []string{"inst_3"})

	region1 = &Region{
		Id: "michigan",
		Vpcs: []*Vpc{
			&Vpc{Id: "vpc_1", Subnets: []*Subnet{
				&Subnet{Id: "sub_1", Instances: []*Instance{
					&Instance{Id: "inst_1"},
					&Instance{Id: "inst_2"},
				}},
			}},
			&Vpc{Id: "vpc_2", Subnets: []*Subnet{
				&Subnet{Id: "sub_2", Instances: []*Instance{
					&Instance{Id: "sub_2_inst_1"},
					&Instance{Id: "sub_2_inst_2"},
				}},
			}},
		},
	}

	region2 = &Region{
		Id: "michigan",
		Vpcs: []*Vpc{
			&Vpc{Id: "vpc_1", Subnets: []*Subnet{
				&Subnet{Id: "sub_1", Instances: []*Instance{
					&Instance{Id: "inst_1"},
					&Instance{Id: "inst_2"},
				}},
			}},
			&Vpc{Id: "vpc_2", Subnets: []*Subnet{
				&Subnet{Id: "sub_2", Instances: []*Instance{
					&Instance{Id: "sub_2_inst_1"},
					&Instance{Id: "sub_2_inst_3"},
				}},
			}},
		},
	}

	diff = Compare(region1, region2)
	compareSlice(t, diff.ExtraInstances, []string{"sub_2_inst_2"})
	compareSlice(t, diff.MissingInstances, []string{"sub_2_inst_3"})

	region1 = &Region{
		Id: "michigan",
		Vpcs: []*Vpc{
			&Vpc{Id: "vpc_1", Subnets: []*Subnet{
				&Subnet{Id: "sub_1", Instances: []*Instance{
					&Instance{Id: "inst_1"},
					&Instance{Id: "inst_2"},
				}},
			}},
			&Vpc{Id: "vpc_2", Subnets: []*Subnet{
				&Subnet{Id: "sub_3", Instances: []*Instance{
					&Instance{Id: "sub_3_inst_1"},
					&Instance{Id: "sub_3_inst_2"},
				}},
			}},
		},
	}

	region2 = &Region{
		Id: "michigan",
		Vpcs: []*Vpc{
			&Vpc{Id: "vpc_1", Subnets: []*Subnet{
				&Subnet{Id: "sub_1", Instances: []*Instance{
					&Instance{Id: "inst_1"},
					&Instance{Id: "inst_2"},
				}},
				&Subnet{Id: "sub_2", Instances: []*Instance{
					&Instance{Id: "sub_2_inst_1"},
					&Instance{Id: "sub_2_inst_2"},
				}},
			}},
			&Vpc{Id: "vpc_2", Subnets: []*Subnet{
				&Subnet{Id: "sub_3", Instances: []*Instance{
					&Instance{Id: "sub_3_inst_1"},
					&Instance{Id: "sub_3_inst_3"},
				}},
			}},
		},
	}

	diff = Compare(region1, region2)
	compareSlice(t, diff.ExtraInstances, []string{"sub_3_inst_2"})
	compareSlice(t, diff.MissingInstances, []string{"sub_2_inst_1", "sub_2_inst_2", "sub_3_inst_3"})

	region1 = &Region{
		Id: "michigan",
		Vpcs: []*Vpc{
			&Vpc{Id: "vpc_1", Subnets: []*Subnet{
				&Subnet{Id: "sub_1", Instances: []*Instance{
					&Instance{Id: "inst_1"},
					&Instance{Id: "inst_2"},
				}},
			}},
		},
	}

	region2 = &Region{
		Id: "michigan",
		Vpcs: []*Vpc{
			&Vpc{Id: "vpc_1", Subnets: []*Subnet{
				&Subnet{Id: "sub_1", Instances: []*Instance{
					&Instance{Id: "inst_1"},
					&Instance{Id: "inst_2"},
				}},
			}},
		},
	}

	diff = Compare(region1, region2)
	if got, want := len(diff.MissingInstances), 0; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := len(diff.ExtraInstances), 0; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
}

func compareSlice(t *testing.T, first, second []string) {
	if len(first) != len(second) {
		t.Fatalf("slice not of same length:\n%v\n%v", first, second)
	}
	for i := 0; i < len(first); i++ {
		if first[i] != second[i] {
			t.Fatalf("slice not equal:\n%v\n%v", first, second)
		}
	}
}
