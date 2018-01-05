package awsat

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/wallix/awless/cloud/properties"
	"github.com/wallix/awless/graph"
	"github.com/wallix/awless/graph/resourcetest"
)

func TestRecord(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		t.Run("with value", func(t *testing.T) {
			Template("create record zone=/hostedzone/1234ABCD name=my.domain.com type=A value=127.0.0.1 ttl=60 comment='this is my localhost record'").
				Mock(&route53Mock{
					ChangeResourceRecordSetsFunc: func(param0 *route53.ChangeResourceRecordSetsInput) (*route53.ChangeResourceRecordSetsOutput, error) {
						return &route53.ChangeResourceRecordSetsOutput{ChangeInfo: &route53.ChangeInfo{Id: String("change-id")}}, nil
					},
				}).ExpectInput("ChangeResourceRecordSets", &route53.ChangeResourceRecordSetsInput{
				HostedZoneId: String("/hostedzone/1234ABCD"),
				ChangeBatch: &route53.ChangeBatch{
					Changes: []*route53.Change{
						{
							ResourceRecordSet: &route53.ResourceRecordSet{
								ResourceRecords: []*route53.ResourceRecord{
									{Value: String("127.0.0.1")},
								},
								Name: String("my.domain.com"),
								Type: String("A"),
								TTL:  Int64(60),
							},
							Action: String("CREATE"),
						},
					},
					Comment: String("this is my localhost record"),
				},
			}).ExpectCommandResult("change-id").ExpectCalls("ChangeResourceRecordSets").Run(t)
		})

		t.Run("with values", func(t *testing.T) {
			Template("create record zone=/hostedzone/1234ABCD name=my.domain2.com type=A values=1.2.3.4,2.3.4.5 ttl=60 comment='this is my localhost record'").
				Mock(&route53Mock{
					ChangeResourceRecordSetsFunc: func(param0 *route53.ChangeResourceRecordSetsInput) (*route53.ChangeResourceRecordSetsOutput, error) {
						return &route53.ChangeResourceRecordSetsOutput{ChangeInfo: &route53.ChangeInfo{Id: String("change-id")}}, nil
					},
				}).ExpectInput("ChangeResourceRecordSets", &route53.ChangeResourceRecordSetsInput{
				HostedZoneId: String("/hostedzone/1234ABCD"),
				ChangeBatch: &route53.ChangeBatch{
					Changes: []*route53.Change{
						{
							ResourceRecordSet: &route53.ResourceRecordSet{
								ResourceRecords: []*route53.ResourceRecord{
									{Value: String("1.2.3.4")},
									{Value: String("2.3.4.5")},
								},
								Name: String("my.domain2.com"),
								Type: String("A"),
								TTL:  Int64(60),
							},
							Action: String("CREATE"),
						},
					},
					Comment: String("this is my localhost record"),
				},
			}).ExpectCommandResult("change-id").ExpectCalls("ChangeResourceRecordSets").Run(t)
		})
	})

	t.Run("update", func(t *testing.T) {
		Template("update record zone=/hostedzone/1234ABCD name=myupdated.domain.com type=A value=127.0.0.1 ttl=60").
			Mock(&route53Mock{
				ChangeResourceRecordSetsFunc: func(param0 *route53.ChangeResourceRecordSetsInput) (*route53.ChangeResourceRecordSetsOutput, error) {
					return &route53.ChangeResourceRecordSetsOutput{ChangeInfo: &route53.ChangeInfo{Id: String("updated-id")}}, nil
				},
			}).ExpectInput("ChangeResourceRecordSets", &route53.ChangeResourceRecordSetsInput{
			HostedZoneId: String("/hostedzone/1234ABCD"),
			ChangeBatch: &route53.ChangeBatch{
				Changes: []*route53.Change{
					{
						ResourceRecordSet: &route53.ResourceRecordSet{
							ResourceRecords: []*route53.ResourceRecord{
								{Value: String("127.0.0.1")},
							},
							Name: String("myupdated.domain.com"),
							Type: String("A"),
							TTL:  Int64(60),
						},
						Action: String("UPSERT"),
					},
				},
			},
		}).ExpectCommandResult("updated-id").ExpectCalls("ChangeResourceRecordSets").Run(t)
	})

	t.Run("delete", func(t *testing.T) {
		t.Run("from awless-id", func(t *testing.T) {
			g := graph.NewGraph()
			zone := resourcetest.Zone("/hostedzone/1234ABCD").Build()
			record := resourcetest.Record("awls-rec").Prop(properties.Name, "mydeleted.domain.com").Prop(properties.Type, "A").Prop(properties.Records, []string{"1.2.3.4"}).Prop(properties.TTL, 60).Build()
			g.AddResource(zone, record)
			g.AddParentRelation(zone, record)
			Template("delete record id=awls-rec").
				Mock(&route53Mock{
					ChangeResourceRecordSetsFunc: func(param0 *route53.ChangeResourceRecordSetsInput) (*route53.ChangeResourceRecordSetsOutput, error) {
						return &route53.ChangeResourceRecordSetsOutput{ChangeInfo: &route53.ChangeInfo{Id: String("deleted-id")}}, nil
					},
				}).ExpectInput("ChangeResourceRecordSets", &route53.ChangeResourceRecordSetsInput{
				HostedZoneId: String("/hostedzone/1234ABCD"),
				ChangeBatch: &route53.ChangeBatch{
					Changes: []*route53.Change{
						{
							ResourceRecordSet: &route53.ResourceRecordSet{
								ResourceRecords: []*route53.ResourceRecord{
									{Value: String("1.2.3.4")},
								},
								Name: String("mydeleted.domain.com"),
								Type: String("A"),
								TTL:  Int64(60),
							},
							Action: String("DELETE"),
						},
					},
				},
			}).Graph(g).ExpectCommandResult("deleted-id").ExpectCalls("ChangeResourceRecordSets").Run(t)
		})

		t.Run("with all params", func(t *testing.T) {
			Template("delete record zone=/hostedzone/1234ABCD name=mydeleted.domain.com type=A value=127.0.0.1 ttl=60").
				Mock(&route53Mock{
					ChangeResourceRecordSetsFunc: func(param0 *route53.ChangeResourceRecordSetsInput) (*route53.ChangeResourceRecordSetsOutput, error) {
						return &route53.ChangeResourceRecordSetsOutput{ChangeInfo: &route53.ChangeInfo{Id: String("deleted-id")}}, nil
					},
				}).ExpectInput("ChangeResourceRecordSets", &route53.ChangeResourceRecordSetsInput{
				HostedZoneId: String("/hostedzone/1234ABCD"),
				ChangeBatch: &route53.ChangeBatch{
					Changes: []*route53.Change{
						{
							ResourceRecordSet: &route53.ResourceRecordSet{
								ResourceRecords: []*route53.ResourceRecord{
									{Value: String("127.0.0.1")},
								},
								Name: String("mydeleted.domain.com"),
								Type: String("A"),
								TTL:  Int64(60),
							},
							Action: String("DELETE"),
						},
					},
				},
			}).ExpectCommandResult("deleted-id").ExpectCalls("ChangeResourceRecordSets").Run(t)
		})
	})
}
