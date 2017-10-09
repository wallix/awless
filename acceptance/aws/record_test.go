package awsat

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/route53"
)

func TestRecord(t *testing.T) {
	t.Run("create", func(t *testing.T) {
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
}
