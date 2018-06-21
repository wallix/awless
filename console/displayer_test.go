/*
Copyright 2017 WALLIX

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package console

import (
	"bytes"
	"testing"
	"time"

	"github.com/fatih/color"
	p "github.com/wallix/awless/cloud/properties"
	"github.com/wallix/awless/graph"
	"github.com/wallix/awless/graph/resourcetest"
)

func init() {
	color.NoColor = true
}

func TestSorting(t *testing.T) {
	g := createInfraGraph()
	columns := []string{"ID", "Name", "State", "Type", "PublicIP"}
	var w bytes.Buffer

	t.Run("ascending", func(t *testing.T) {
		w.Reset()
		displayer, _ := BuildOptions(
			WithRdfType("instance"),
			WithColumns(columns),
			WithFormat("csv"),
			WithSortBy("Name"),
		).SetSource(g).Build()

		expected := "ID,Name,State,Type,Public IP\n" +
			"inst_3,apache,running,t2.xlarge,\n" +
			"inst_2,django,stopped,t2.medium,\n" +
			"inst_1,redis,running,t2.micro,1.2.3.4\n"

		if err := displayer.Print(&w); err != nil {
			t.Fatal(err)
		}
		if got, want := w.String(), expected; got != want {
			t.Fatalf("got \n%q\n\nwant\n\n%q\n", got, want)
		}
	})

	t.Run("descending", func(t *testing.T) {
		w.Reset()
		displayer, _ := BuildOptions(
			WithRdfType("instance"),
			WithColumns(columns),
			WithFormat("csv"),
			WithSortBy("Name"),
			WithReverseSort(true),
		).SetSource(g).Build()

		expected := "ID,Name,State,Type,Public IP\n" +
			"inst_1,redis,running,t2.micro,1.2.3.4\n" +
			"inst_2,django,stopped,t2.medium,\n" +
			"inst_3,apache,running,t2.xlarge,\n"

		if err := displayer.Print(&w); err != nil {
			t.Fatal(err)
		}
		if got, want := w.String(), expected; got != want {
			t.Fatalf("got \n%q\n\nwant\n\n%q\n", got, want)
		}
	})
}

func TestJSONDisplays(t *testing.T) {
	g := createInfraGraph()
	var w bytes.Buffer

	t.Run("Single resource", func(t *testing.T) {
		displayer, _ := BuildOptions(
			WithRdfType("instance"),
			WithFormat("json"),
		).SetSource(g).Build()

		expected := `[{"ID": "inst_1", "Name": "redis", "PublicIP": "1.2.3.4", "State": "running", "Type": "t2.micro"},
		{"ID": "inst_2", "Name": "django", "State": "stopped", "Type": "t2.medium" },
		{"ID": "inst_3", "Name": "apache", "State": "running", "Type": "t2.xlarge"}]`

		if err := displayer.Print(&w); err != nil {
			t.Fatal(err)
		}

		compareJSON(t, w.String(), expected)
	})

	t.Run("Multi resource", func(t *testing.T) {
		t.Skip("Comparison fail: until we can order what is inside each resource")

		displayer, _ := BuildOptions(
			WithFormat("json"),
		).SetSource(g).Build()

		expected := `{"instances": [
			{ "ID": "inst_1", "Name": "redis", "PublicIP": "1.2.3.4", "State": "running", "Type": "t2.micro"},
		  { "ID": "inst_2", "Name": "django", "State": "stopped", "Type": "t2.medium" },
		  { "ID": "inst_3", "Name": "apache", "State": "running", "Type": "t2.xlarge" }
		 ], "subnets": [
		  { "ID": "sub_1", "Name": "my_subnet", "Vpc": "vpc_1" }, {"ID": "sub_2", "Vpc": "vpc_2" }
		 ], "vpcs": [
		  { "ID": "vpc_1", "NewProp": "my_value" }, { "ID": "vpc_2", "Name": "my_vpc_2" }
		 ]}`

		w.Reset()
		if err := displayer.Print(&w); err != nil {
			t.Fatal(err)
		}

		compareJSON(t, w.String(), expected)
	})
}

func TestTabularDisplays(t *testing.T) {
	g := createInfraGraph()
	columns := []string{"ID", "Name", "State", "Type", "PublicIP"}

	displayer, _ := BuildOptions(
		WithRdfType("instance"),
		WithColumns(columns),
		WithFormat("csv"),
	).SetSource(g).Build()

	expected := "ID,Name,State,Type,Public IP\n" +
		"inst_1,redis,running,t2.micro,1.2.3.4\n" +
		"inst_2,django,stopped,t2.medium,\n" +
		"inst_3,apache,running,t2.xlarge,\n"
	var w bytes.Buffer
	if err := displayer.Print(&w); err != nil {
		t.Fatal(err)
	}
	if got, want := w.String(), expected; got != want {
		t.Fatalf("got \n[%q]\n\nwant\n\n[%q]\n", got, want)
	}

	displayer, _ = BuildOptions(
		WithRdfType("instance"),
		WithColumns(columns),
		WithFormat("csv"),
		WithSortBy("Name"),
	).SetSource(g).Build()

	expected = "ID,Name,State,Type,Public IP\n" +
		"inst_3,apache,running,t2.xlarge,\n" +
		"inst_2,django,stopped,t2.medium,\n" +
		"inst_1,redis,running,t2.micro,1.2.3.4\n"

	w.Reset()
	if err := displayer.Print(&w); err != nil {
		t.Fatal(err)
	}
	if got, want := w.String(), expected; got != want {
		t.Fatalf("got \n%q\n\nwant\n\n%q\n", got, want)
	}

	columns = []string{"ID", "Name", "State", "Type", "PublicIP"}

	displayer, _ = BuildOptions(
		WithRdfType("instance"),
		WithColumns(columns),
	).SetSource(g).Build()

	expected = `|  ID ▲  |  NAME  |  STATE  |   TYPE    | PUBLIC IP |
|--------|--------|---------|-----------|-----------|
| inst_1 | redis  | running | t2.micro  | 1.2.3.4   |
| inst_2 | django | stopped | t2.medium |           |
| inst_3 | apache | running | t2.xlarge |           |
`
	w.Reset()
	if err := displayer.Print(&w); err != nil {
		t.Fatal(err)
	}
	if got, want := w.String(), expected; got != want {
		t.Fatalf("got \n%s\n\nwant\n\n%s\n", got, want)
	}

	displayer, _ = BuildOptions(
		WithRdfType("instance"),
		WithColumns(columns),
		WithSortBy("state", "id"),
	).SetSource(g).Build()

	expected = `|   ID   |  NAME  | STATE ▲ |   TYPE    | PUBLIC IP |
|--------|--------|---------|-----------|-----------|
| inst_1 | redis  | running | t2.micro  | 1.2.3.4   |
| inst_3 | apache | running | t2.xlarge |           |
| inst_2 | django | stopped | t2.medium |           |
`
	w.Reset()
	if err := displayer.Print(&w); err != nil {
		t.Fatal(err)
	}
	if got, want := w.String(), expected; got != want {
		t.Fatalf("got \n%s\n\nwant\n\n%s\n", got, want)
	}

	displayer, _ = BuildOptions(
		WithRdfType("instance"),
		WithColumns(columns),
		WithSortBy("state", "name"),
	).SetSource(g).Build()

	expected = `|   ID   |  NAME  | STATE ▲ |   TYPE    | PUBLIC IP |
|--------|--------|---------|-----------|-----------|
| inst_3 | apache | running | t2.xlarge |           |
| inst_1 | redis  | running | t2.micro  | 1.2.3.4   |
| inst_2 | django | stopped | t2.medium |           |
`
	w.Reset()
	if err := displayer.Print(&w); err != nil {
		t.Fatal(err)
	}
	if got, want := w.String(), expected; got != want {
		t.Fatalf("got \n%s\n\nwant\n\n%s\n", got, want)
	}

	columns = []string{"ID", "Name"}

	displayer, _ = BuildOptions(
		WithRdfType("instance"),
		WithColumns(columns),
		WithFormat("porcelain"),
	).SetSource(g).Build()

	expected = `inst_1
redis
inst_2
django
inst_3
apache`

	w.Reset()
	if err := displayer.Print(&w); err != nil {
		t.Fatal(err)
	}
	if got, want := w.String(), expected; got != want {
		t.Fatalf("got \n%s\n\nwant\n\n%s\n", got, want)
	}
}

func TestMultiResourcesDisplays(t *testing.T) {
	g := createInfraGraph()

	displayer, _ := BuildOptions(
		WithFormat("table"),
	).SetSource(g).Build()

	expected := `|  TYPE ▲  |  NAME/ID  | PROPERTY  |   VALUE   |
|----------|-----------|-----------|-----------|
| instance | apache    | ID        | inst_3    |
|          |           | Name      | apache    |
|          |           | State     | running   |
|          |           | Type      | t2.xlarge |
|          | django    | ID        | inst_2    |
|          |           | Name      | django    |
|          |           | State     | stopped   |
|          |           | Type      | t2.medium |
|          | redis     | ID        | inst_1    |
|          |           | Name      | redis     |
|          |           | Public IP | 1.2.3.4   |
|          |           | State     | running   |
|          |           | Type      | t2.micro  |
| subnet   | my_subnet | ID        | sub_1     |
|          |           | Name      | my_subnet |
|          |           | Vpc       | vpc_1     |
|          | sub_2     | ID        | sub_2     |
|          |           | Vpc       | vpc_2     |
| vpc      | my_vpc_2  | ID        |           |
|          |           | Name      | my_vpc_2  |
|          | vpc_1     | ID        | vpc_1     |
`
	var w bytes.Buffer
	if err := displayer.Print(&w); err != nil {
		t.Fatal(err)
	}
	if got, want := w.String(), expected; got != want {
		t.Fatalf("got \n%s\n\nwant\n\n%s\n", got, want)
	}

	displayer, _ = BuildOptions(WithColumns([]string{"ID"}), WithFormat("porcelain")).SetSource(g).Build()

	expected = `inst_1
inst_2
inst_3
sub_1
sub_2
vpc_1
vpc_2`

	w.Reset()
	if err := displayer.Print(&w); err != nil {
		t.Fatal(err)
	}
	if got, want := w.String(), expected; got != want {
		t.Fatalf("got \n%s\n\nwant\n\n%s\n", got, want)
	}

	displayer, _ = BuildOptions(
		WithFormat("porcelain"),
		WithIDsOnly(true),
	).SetSource(g).Build()

	expected = `inst_1
redis
inst_2
django
inst_3
apache
sub_1
my_subnet
sub_2
vpc_1
vpc_2
my_vpc_2`

	w.Reset()
	if err := displayer.Print(&w); err != nil {
		t.Fatal(err)
	}
	if got, want := w.String(), expected; got != want {
		t.Fatalf("got \n%s\n\nwant\n\n%s\n", got, want)
	}
}

func TestDiffDisplay(t *testing.T) {
	rootNode := graph.InitResource("region", "eu-west-1")
	diff, err := createDiff(rootNode)
	if err != nil {
		t.Fatal(err)
	}

	displayer, _ := BuildOptions(
		WithFormat("table"),
		WithRootNode(rootNode),
	).SetSource(diff).Build()

	expected := `|  TYPE ▲  |   NAME/ID    | PROPERTY |  VALUE   |
|----------|--------------|----------|----------|
| instance | + inst_4     |          |          |
|          | + inst_5     |          |          |
|          | + inst_6     |          |          |
|          | - inst_2     |          |          |
|          | redis        | ID       | + new_id |
|          |              |          | - inst_1 |
| subnet   | + new_subnet |          |          |
| vpc      | vpc_1        | Default  | - true   |
`
	var w bytes.Buffer
	if err := displayer.Print(&w); err != nil {
		t.Fatal(err)
	}
	if got, want := w.String(), expected; got != want {
		t.Errorf("got \n%s\n\nwant \n%s\n", got, want)
	}

	displayer, _ = BuildOptions(
		WithFormat("tree"),
		WithRootNode(rootNode),
	).SetSource(diff).Build()

	expected = `region, eu-west-1
	vpc, vpc_2
+		subnet, new_subnet
+			instance, inst_6
		subnet, sub_2
-			instance, inst_2
+			instance, inst_4
+			instance, inst_5
`
	w.Reset()
	if err := displayer.Print(&w); err != nil {
		t.Fatal(err)
	}
	if got, want := w.String(), expected; got != want {
		t.Fatalf("got \n%s\n\nwant\n\n%s\n", got, want)
	}
}

func TestDateLists(t *testing.T) {
	g := graph.NewGraph()
	g.AddResource(resourcetest.Region("eu-west-1").Build(),
		resourcetest.User("user1").Prop("Name", "my_username_1").Build(),
		resourcetest.User("user2").Prop("Name", "my_username_2").Prop("PasswordLastUsed", time.Unix(1482405203, 0).UTC()).Build(),
		resourcetest.User("user3").Prop("Name", "my_username_3").Prop("PasswordLastUsed", time.Unix(1481358937, 0).UTC()).Build(),
	)
	globalNow = time.Unix(1505832866, 0)
	defer func() {
		globalNow = time.Now().UTC()
	}()

	columns := []string{"ID", "Name", "PasswordLastUsed"}

	displayer, _ := BuildOptions(
		WithRdfType("user"),
		WithColumns(columns),
	).SetSource(g).Build()

	expected := `| ID ▲  |     NAME      | PASSWORDLASTUSED |
|-------|---------------|------------------|
| user1 | my_username_1 |                  |
| user2 | my_username_2 | 9 months         |
| user3 | my_username_3 | 9 months         |
`
	var w bytes.Buffer
	if err := displayer.Print(&w); err != nil {
		t.Fatal(err)
	}
	if got, want := w.String(), expected; got != want {
		t.Fatalf("got \n%s\n\nwant\n\n%s\n", got, want)
	}

	displayer, _ = BuildOptions(
		WithRdfType("user"),
		WithColumns(columns),
		WithSortBy("passwordlastused"),
	).SetSource(g).Build()

	expected = `|  ID   |     NAME      | PASSWORDLASTUSED ▲ |
|-------|---------------|--------------------|
| user1 | my_username_1 |                    |
| user2 | my_username_2 | 9 months           |
| user3 | my_username_3 | 9 months           |
`
	w.Reset()
	if err := displayer.Print(&w); err != nil {
		t.Fatal(err)
	}
	if got, want := w.String(), expected; got != want {
		t.Fatalf("got \n%s\n\nwant\n\n%s\n", got, want)
	}
}

func TestMaxWidth(t *testing.T) {
	g := createInfraGraph()
	columns := []string{"ID", "Name", "State", "Type", "PublicIP"}

	displayer, _ := BuildOptions(
		WithRdfType("instance"),
		WithColumns(columns),
		WithSortBy("state", "name"),
		WithMaxWidth(55),
	).SetSource(g).Build()

	expected := `|   ID   |  NAME  | STATE ▲ |   TYPE    | PUBLIC IP |
|--------|--------|---------|-----------|-----------|
| inst_3 | apache | running | t2.xlarge |           |
| inst_1 | redis  | running | t2.micro  | 1.2.3.4   |
| inst_2 | django | stopped | t2.medium |           |
`
	var w bytes.Buffer
	if err := displayer.Print(&w); err != nil {
		t.Fatal(err)
	}
	if got, want := w.String(), expected; got != want {
		t.Fatalf("got \n%s\n\nwant\n\n%s\n", got, want)
	}

	columns = []string{"ID", "Name", "State", "Type", "PublicIP"}

	autowrapMaxSize = 4
	tableColWidth = 4
	displayer, _ = BuildOptions(
		WithRdfType("instance"),
		WithColumns(columns),
		WithSortBy("state", "name"),
		WithMaxWidth(45),
	).SetSource(g).Build()

	expected = `|  ID  | NAME | STATE ▲ | TYPE | PUBLIC IP |
|------|------|------|------|------|
| inst | apac | runn | t2.  |      |
| _3   | he   | ing  | xlar |      |
|      |      |      | ge   |      |
| inst | redi | runn | t2.  | 1.2. |
| _1   | s    | ing  | micr | 3.4  |
|      |      |      | o    |      |
| inst | djan | stop | t2.  |      |
| _2   | go   | ped  | medi |      |
|      |      |      | um   |      |
`
	w.Reset()
	if err := displayer.Print(&w); err != nil {
		t.Fatal(err)
	}
	if got, want := w.String(), expected; got != want {
		t.Fatalf("got \n%s\n\nwant\n\n%s\n", got, want)
	}

	displayer, _ = BuildOptions(
		WithRdfType("instance"),
		WithColumns(columns),
		WithSortBy("state", "name"),
		WithMaxWidth(70),
	).SetSource(g).Build()

	expected = `|   ID   |  NAME  | STATE ▲ |   TYPE    | PUBLIC IP |
|--------|--------|---------|-----------|---------|
| inst_3 | apache | running | t2.xlarge |         |
| inst_1 | redis  | running | t2.micro  | 1.2.3.4 |
| inst_2 | django | stopped | t2.medium |         |
`
	w.Reset()
	if err := displayer.Print(&w); err != nil {
		t.Fatal(err)
	}
	if got, want := w.String(), expected; got != want {
		t.Fatalf("got \n%s\n\nwant\n\n%s\n", got, want)
	}

	columnDefs := []ColumnDefinition{
		StringColumnDefinition{Prop: "ID", Friendly: "I"},
		StringColumnDefinition{Prop: "Name", Friendly: "N"},
		StringColumnDefinition{Prop: "State", Friendly: "S"},
		StringColumnDefinition{Prop: "Type", Friendly: "T"},
		StringColumnDefinition{Prop: "PublicIP", Friendly: "P"},
	}

	builder := BuildOptions(
		WithRdfType("instance"),
		WithColumnDefinitions(columnDefs),
		WithSortBy("s", "n"),
		WithMaxWidth(40),
	)

	displayer, _ = builder.SetSource(g).Build()

	expected = `|  I   |  N   | S ▲  |  T   |  P   |
|------|------|------|------|------|
| inst | apac | runn | t2.  |      |
| _3   | he   | ing  | xlar |      |
|      |      |      | ge   |      |
| inst | redi | runn | t2.  | 1.2. |
| _1   | s    | ing  | micr | 3.4  |
|      |      |      | o    |      |
| inst | djan | stop | t2.  |      |
| _2   | go   | ped  | medi |      |
|      |      |      | um   |      |
`
	w.Reset()
	if err := displayer.Print(&w); err != nil {
		t.Fatal(err)
	}
	if got, want := w.String(), expected; got != want {
		t.Fatalf("got \n%s\n\nwant\n\n%s\n", got, want)
	}

	displayer, _ = BuildOptions(
		WithRdfType("instance"),
		WithColumnDefinitions(columnDefs),
		WithSortBy("s", "n"),
		WithMaxWidth(50),
	).SetSource(g).Build()

	w.Reset()
	if err := displayer.Print(&w); err != nil {
		t.Fatal(err)
	}
	if got, want := w.String(), expected; got != want {
		t.Fatalf("got \n%s\n\nwant\n\n%s\n", got, want)
	}

	autowrapMaxSize = 5
	tableColWidth = 5
	displayer, _ = BuildOptions(
		WithRdfType("instance"),
		WithColumnDefinitions(columnDefs),
		WithSortBy("s", "n"),
		WithMaxWidth(29),
	).SetSource(g).Build()

	expected = `|   I   |   N   |  S ▲  |
|-------|-------|-------|
| inst_ | apach | runni |
| 3     | e     | ng    |
| inst_ | redis | runni |
| 1     |       | ng    |
| inst_ | djang | stopp |
| 2     | o     | ed    |
Columns truncated to fit terminal: 'T', 'P'
`
	w.Reset()
	if err := displayer.Print(&w); err != nil {
		t.Fatal(err)
	}
	if got, want := w.String(), expected; got != want {
		t.Fatalf("got \n%s\n\nwant\n\n%s\n", got, want)
	}
}

func TestFilter(t *testing.T) {
	g := graph.NewGraph()
	g.AddResource(
		resourcetest.Subnet("sub_1").Prop(p.Name, "my_subnet").Prop(p.Vpc, "vpc_1").Prop(p.Public, true).Build(),
		resourcetest.Subnet("sub_2").Prop(p.Vpc, "vpc_2").Prop(p.Public, false).Build(),
		resourcetest.Subnet("sub_3").Prop(p.Name, "my_subnet").Prop(p.Vpc, "vpc_1").Prop(p.Public, false).Build(),
	)

	t.Run("No filter", func(t *testing.T) {
		var w bytes.Buffer
		displayer, _ := BuildOptions(
			WithRdfType("subnet"),
			WithFormat("json"),
		).SetSource(g).Build()
		expected := `[{"ID":"sub_1","Public":true,"Name":"my_subnet","Vpc":"vpc_1"},
		{"ID":"sub_2","Public":false,"Vpc":"vpc_2"},
		{"ID":"sub_3","Public":false,"Name":"my_subnet","Vpc":"vpc_1"}]`
		if err := displayer.Print(&w); err != nil {
			t.Fatal(err)
		}
		compareJSON(t, w.String(), expected)
	})
	t.Run("Filter column name", func(t *testing.T) {
		var w bytes.Buffer
		displayer, _ := BuildOptions(
			WithRdfType("subnet"),
			WithFormat("json"),
			WithFilters([]string{"Vpc=vpc_1"}),
		).SetSource(g).Build()
		expected := `[{"ID":"sub_1","Public":true,"Name":"my_subnet","Vpc":"vpc_1"},
		{"ID":"sub_3","Public":false,"Name":"my_subnet","Vpc":"vpc_1"}]`
		if err := displayer.Print(&w); err != nil {
			t.Fatal(err)
		}
		compareJSON(t, w.String(), expected)
	})
	t.Run("Filter friendly name", func(t *testing.T) {
		var w bytes.Buffer
		displayer, _ := BuildOptions(
			WithRdfType("subnet"),
			WithFormat("json"),
			WithFilters([]string{"public=false"}),
		).SetSource(g).Build()
		expected := `[{"ID":"sub_2","Public":false,"Vpc":"vpc_2"},
		{"ID":"sub_3","Public":false,"Name":"my_subnet","Vpc":"vpc_1"}]`
		if err := displayer.Print(&w); err != nil {
			t.Fatal(err)
		}
		compareJSON(t, w.String(), expected)
	})
}

func TestCompareInterface(t *testing.T) {
	if got, want := valueLowerOrEqual(interface{}(1), interface{}(4)), true; got != want {
		t.Fatalf("got %t want %t", got, want)
	}
	if got, want := valueLowerOrEqual(interface{}(1), interface{}(1)), true; got != want {
		t.Fatalf("got %t want %t", got, want)
	}
	if got, want := valueLowerOrEqual(interface{}(1), interface{}(-3)), false; got != want {
		t.Fatalf("got %t want %t", got, want)
	}
	if got, want := valueLowerOrEqual(interface{}("abc"), interface{}("bbc")), true; got != want {
		t.Fatalf("got %t want %t", got, want)
	}
	if got, want := valueLowerOrEqual(interface{}("abc"), interface{}("aac")), false; got != want {
		t.Fatalf("got %t want %t", got, want)
	}
	if got, want := valueLowerOrEqual(interface{}(1.2), interface{}(1.3)), true; got != want {
		t.Fatalf("got %t want %t", got, want)
	}
	if got, want := valueLowerOrEqual(interface{}(1.2), interface{}(1.1)), false; got != want {
		t.Fatalf("got %t want %t", got, want)
	}
}

func TestCSVDisplayWithCommaAndQuotes(t *testing.T) {
	g := graph.NewGraph()
	g.AddResource(
		resourcetest.Instance("inst_1").Prop(p.Name, "with,comma").Build(),
		resourcetest.Instance("inst_2").Prop(p.Name, "with\nlinebreak").Build(),
		resourcetest.Instance("inst_3").Prop(p.Name, "with\"quote").Build(),
	)

	displayer, _ := BuildOptions(
		WithRdfType("instance"),
		WithColumns([]string{"ID", "Name"}),
		WithFormat("csv"),
	).SetSource(g).Build()

	expected := "ID,Name\n" +
		"inst_1,\"with,comma\"\n" +
		"inst_2,\"with\nlinebreak\"\n" +
		"inst_3,\"with\"\"quote\"\n"
	var w bytes.Buffer
	if err := displayer.Print(&w); err != nil {
		t.Fatal(err)
	}
	if got, want := w.String(), expected; got != want {
		t.Fatalf("got \n[%q]\n\nwant\n\n[%q]\n", got, want)
	}
}

func createInfraGraph() *graph.Graph {
	g := graph.NewGraph()
	g.AddResource(resourcetest.Region("eu-west-1").Build(),
		resourcetest.Instance("inst_1").Prop(p.Name, "redis").Prop(p.Type, "t2.micro").Prop(p.PublicIP, "1.2.3.4").Prop(p.State, "running").Build(),
		resourcetest.Instance("inst_2").Prop(p.Name, "django").Prop(p.Type, "t2.medium").Prop(p.State, "stopped").Build(),
		resourcetest.Instance("inst_3").Prop(p.Name, "apache").Prop(p.Type, "t2.xlarge").Prop(p.State, "running").Build(),
		resourcetest.VPC("vpc_1").Build(),
		resourcetest.VPC("vpc_2").Prop(p.Name, "my_vpc_2").Build(),
		resourcetest.Subnet("sub_1").Prop(p.Name, "my_subnet").Prop(p.Vpc, "vpc_1").Build(),
		resourcetest.Subnet("sub_2").Prop(p.Vpc, "vpc_2").Build(),
	)

	return g
}

func createDiff(root *graph.Resource) (*graph.Diff, error) {
	localDiffG := graph.NewGraph()
	err := localDiffG.AddResource(
		resourcetest.Region("eu-west-1").Build(),
		resourcetest.Instance("inst_1").Prop(p.Name, "redis").Prop(p.Type, "t2.micro").Prop(p.State, "running").Build(),
		resourcetest.Instance("inst_2").Build(),
		resourcetest.Instance("inst_3").Prop(p.Name, "apache").Prop(p.Type, "t2.xlarge").Prop(p.State, "running").Build(),
		resourcetest.Subnet("sub_1").Prop(p.Name, "my_subnet").Prop(p.Vpc, "vpc_1").Build(),
		resourcetest.Subnet("sub_2").Prop(p.Vpc, "vpc_2").Build(),
		resourcetest.VPC("vpc_1").Prop(p.Default, true).Build(),
		resourcetest.VPC("vpc_2").Prop(p.Name, "my_vpc_2").Build(),
	)
	resourcetest.AddParents(localDiffG,
		"eu-west-1 -> vpc_1", "eu-west-1 -> vpc_2",
		"vpc_1 -> sub_1", "vpc_2 -> sub_2",
		"sub_1 -> inst_1", "sub_2 -> inst_2", "sub_2 -> inst_3",
	)
	if err != nil {
		panic(err)
	}

	remoteDiffG := graph.NewGraph()
	err = remoteDiffG.AddResource(
		resourcetest.Region("eu-west-1").Build(),
		resourcetest.Instance("inst_1").Prop(p.ID, "new_id").Prop(p.Name, "redis").Prop(p.Type, "t2.micro").Prop(p.State, "running").Build(),
		resourcetest.Instance("inst_3").Prop(p.Name, "apache").Prop(p.Type, "t2.xlarge").Prop(p.State, "running").Build(),
		resourcetest.Instance("inst_4").Build(),
		resourcetest.Instance("inst_5").Build(),
		resourcetest.Instance("inst_6").Build(),
		resourcetest.Subnet("sub_1").Prop(p.Name, "my_subnet").Prop(p.Vpc, "vpc_1").Build(),
		resourcetest.Subnet("sub_2").Prop(p.Vpc, "vpc_2").Build(),
		resourcetest.Subnet("new_subnet").Build(),
		resourcetest.VPC("vpc_1").Build(),
		resourcetest.VPC("vpc_2").Prop(p.Name, "my_vpc_2").Build(),
	)
	resourcetest.AddParents(remoteDiffG,
		"eu-west-1 -> vpc_1", "eu-west-1 -> vpc_2",
		"vpc_1 -> sub_1", "vpc_2 -> sub_2", "vpc_2 -> new_subnet",
		"sub_1 -> inst_1", "sub_2 -> inst_3", "sub_2 -> inst_4", "sub_2 -> inst_5", "new_subnet -> inst_6",
	)
	if err != nil {
		panic(err)
	}

	return graph.DefaultDiffer.Run(root.Id(), localDiffG, remoteDiffG)
}

func TestEmotyDisplays(t *testing.T) {
	g := graph.NewGraph()
	columns := []string{"ID", "Name", "PublicIP"}

	displayer, _ := BuildOptions(
		WithRdfType("instance"),
		WithColumns(columns),
		WithFormat("csv"),
	).SetSource(g).Build()

	expected := "ID,Name,Public IP\n"
	var w bytes.Buffer
	if err := displayer.Print(&w); err != nil {
		t.Fatal(err)
	}
	if got, want := w.String(), expected; got != want {
		t.Fatalf("got \n[%q]\n\nwant\n\n[%q]\n", got, want)
	}

	displayer, _ = BuildOptions(
		WithColumns(columns),
		WithRdfType("instance"),
		WithFormat("table"),
	).SetSource(g).Build()

	expected = "No results found.\n"
	w.Reset()
	if err := displayer.Print(&w); err != nil {
		t.Fatal(err)
	}
	if got, want := w.String(), expected; got != want {
		t.Fatalf("got \n[%q]\n\nwant\n\n[%q]\n", got, want)
	}

	g = createInfraGraph()
	columns = []string{}
	DefaultsColumnDefinitions = make(map[string][]ColumnDefinition)

	displayer, _ = BuildOptions(
		WithColumns(columns),
		WithRdfType("instance"),
		WithFormat("csv"),
	).SetSource(g).Build()

	expected = ""
	w.Reset()
	if err := displayer.Print(&w); err != nil {
		t.Fatal(err)
	}
	if got, want := w.String(), expected; got != want {
		t.Fatalf("got \n[%q]\n\nwant\n\n[%q]\n", got, want)
	}

	displayer, _ = BuildOptions(
		WithColumns(columns),
		WithRdfType("instance"),
		WithFormat("table"),
	).SetSource(g).Build()

	expected = "No columns to display.\n"
	w.Reset()
	if err := displayer.Print(&w); err != nil {
		t.Fatal(err)
	}
	if got, want := w.String(), expected; got != want {
		t.Fatalf("got \n[%q]\n\nwant\n\n[%q]\n", got, want)
	}
}

func TestNoHeadersDisplay(t *testing.T) {
	g := createInfraGraph()
	columnDefs := []ColumnDefinition{
		StringColumnDefinition{Prop: "ID"},
		StringColumnDefinition{Prop: "Name"},
		StringColumnDefinition{Prop: "State"},
		StringColumnDefinition{Prop: "Type"},
		StringColumnDefinition{Prop: "PublicIP", Friendly: "Public IP"},
	}

	displayer, _ := BuildOptions(
		WithRdfType("instance"),
		WithColumnDefinitions(columnDefs),
		WithFormat("csv"),
	).SetSource(g).Build()

	expected := "ID,Name,State,Type,Public IP\n" +
		"inst_1,redis,running,t2.micro,1.2.3.4\n" +
		"inst_2,django,stopped,t2.medium,\n" +
		"inst_3,apache,running,t2.xlarge,\n"
	var w bytes.Buffer
	if err := displayer.Print(&w); err != nil {
		t.Fatal(err)
	}
	if got, want := w.String(), expected; got != want {
		t.Fatalf("got \n[%q]\n\nwant\n\n[%q]\n", got, want)
	}

	displayer, _ = BuildOptions(
		WithRdfType("instance"),
		WithColumnDefinitions(columnDefs),
		WithFormat("csv"),
		WithNoHeaders(true),
	).SetSource(g).Build()

	expected = "inst_1,redis,running,t2.micro,1.2.3.4\n" +
		"inst_2,django,stopped,t2.medium,\n" +
		"inst_3,apache,running,t2.xlarge,\n"
	w.Reset()

	if err := displayer.Print(&w); err != nil {
		t.Fatal(err)
	}
	if got, want := w.String(), expected; got != want {
		t.Fatalf("got \n%s\n\nwant\n\n%s\n", got, want)
	}

	displayer, _ = BuildOptions(
		WithRdfType("instance"),
		WithColumnDefinitions(columnDefs),
		WithFormat("tsv"),
		WithNoHeaders(true),
	).SetSource(g).Build()

	expected = "inst_1\tredis\trunning\tt2.micro\t1.2.3.4\n" +
		"inst_2\tdjango\tstopped\tt2.medium\t\n" +
		"inst_3\tapache\trunning\tt2.xlarge\t\n"
	w.Reset()

	if err := displayer.Print(&w); err != nil {
		t.Fatal(err)
	}
	if got, want := w.String(), expected; got != want {
		t.Fatalf("got \n%s\n\nwant\n\n%s\n", got, want)
	}

	autowrapMaxSize = 20
	tableColWidth = 20
	displayer, _ = BuildOptions(
		WithRdfType("instance"),
		WithColumnDefinitions(columnDefs),
		WithFormat("table"),
		WithNoHeaders(true),
	).SetSource(g).Build()

	expected = `| inst_1 | redis  | running | t2.micro  | 1.2.3.4 |
| inst_2 | django | stopped | t2.medium |         |
| inst_3 | apache | running | t2.xlarge |         |
`
	w.Reset()

	if err := displayer.Print(&w); err != nil {
		t.Fatal(err)
	}
	if got, want := w.String(), expected; got != want {
		t.Fatalf("got \n%s\n\nwant\n\n%s\n", got, want)
	}
}
