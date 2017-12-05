package awsat

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/wallix/awless/aws/spec"
)

func TestDistribution(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		awsspec.CallerReferenceFunc = func() string {
			return "callerReference"
		}
		t.Run("all parameters", func(t *testing.T) {
			Template("create distribution origin-domain=my.test.domain.com certificate=arn:of:the:certificate comment='useless comment' default-file=index.go "+
				"domain-aliases=any.domain.com,other.domain.com enable=true forward-cookies=whitelist forward-queries=true https-behaviour=redirect-to-https "+
				"origin-path=/my/custom/path price-class=PriceClass_All min-ttl=42").
				Mock(&cloudfrontMock{
					CreateDistributionFunc: func(param0 *cloudfront.CreateDistributionInput) (*cloudfront.CreateDistributionOutput, error) {
						return &cloudfront.CreateDistributionOutput{Distribution: &cloudfront.Distribution{Id: String("new-distribution-id")}}, nil
					},
				}).ExpectInput("CreateDistribution", &cloudfront.CreateDistributionInput{
				DistributionConfig: &cloudfront.DistributionConfig{
					CallerReference: String("callerReference"),
					Origins: &cloudfront.Origins{
						Items: []*cloudfront.Origin{
							{
								DomainName: String("my.test.domain.com"),
								Id:         String("orig_1"),
								OriginPath: String("/my/custom/path"),
							},
						},
						Quantity: Int64(1),
					},
					ViewerCertificate: &cloudfront.ViewerCertificate{
						ACMCertificateArn: String("arn:of:the:certificate"),
						SSLSupportMethod:  String("sni-only"),
					},
					Comment:           String("useless comment"),
					DefaultRootObject: String("index.go"),
					Aliases: &cloudfront.Aliases{
						Items:    []*string{String("any.domain.com"), String("other.domain.com")},
						Quantity: Int64(2),
					},
					Enabled: Bool(true),
					DefaultCacheBehavior: &cloudfront.DefaultCacheBehavior{
						ForwardedValues: &cloudfront.ForwardedValues{
							Cookies: &cloudfront.CookiePreference{
								Forward: String("whitelist"),
							},
							QueryString: Bool(true),
						},
						ViewerProtocolPolicy: String("redirect-to-https"),
						MinTTL:               Int64(42),
						TargetOriginId:       aws.String("orig_1"),
						TrustedSigners: &cloudfront.TrustedSigners{
							Enabled:  aws.Bool(false),
							Quantity: aws.Int64(0),
						},
					},
					PriceClass: String("PriceClass_All"),
				},
			}).
				ExpectCommandResult("new-distribution-id").ExpectCalls("CreateDistribution").Run(t)
		})
		t.Run("one parameter", func(t *testing.T) {
			Template("create distribution origin-domain=my.test.domain.com").
				Mock(&cloudfrontMock{
					CreateDistributionFunc: func(param0 *cloudfront.CreateDistributionInput) (*cloudfront.CreateDistributionOutput, error) {
						return &cloudfront.CreateDistributionOutput{Distribution: &cloudfront.Distribution{Id: String("new-distribution-id")}}, nil
					},
				}).ExpectInput("CreateDistribution", &cloudfront.CreateDistributionInput{
				DistributionConfig: &cloudfront.DistributionConfig{
					Comment: aws.String("my.test.domain.com"),
					DefaultCacheBehavior: &cloudfront.DefaultCacheBehavior{
						MinTTL: aws.Int64(0),
						ForwardedValues: &cloudfront.ForwardedValues{
							Cookies:     &cloudfront.CookiePreference{Forward: aws.String("all")},
							QueryString: aws.Bool(true),
						},
						TrustedSigners: &cloudfront.TrustedSigners{
							Enabled:  aws.Bool(false),
							Quantity: aws.Int64(0),
						},
						TargetOriginId:       aws.String("orig_1"),
						ViewerProtocolPolicy: aws.String("allow-all"),
					},
					Enabled:         aws.Bool(true),
					CallerReference: String("callerReference"),
					Origins: &cloudfront.Origins{
						Items: []*cloudfront.Origin{
							{
								DomainName: String("my.test.domain.com"),
								Id:         String("orig_1"),
							},
						},
						Quantity: Int64(1),
					},
				},
			}).
				ExpectCommandResult("new-distribution-id").ExpectCalls("CreateDistribution").Run(t)
		})
	})

	t.Run("update", func(t *testing.T) {
		t.Run("already enabled", func(t *testing.T) {
			Template("update distribution id=my-distribution-to-update enable=true").Mock(&cloudfrontMock{
				GetDistributionFunc: func(input *cloudfront.GetDistributionInput) (*cloudfront.GetDistributionOutput, error) {
					return &cloudfront.GetDistributionOutput{
						ETag: String("etag-already-enabled"),
						Distribution: &cloudfront.Distribution{
							DistributionConfig: &cloudfront.DistributionConfig{
								Enabled: Bool(true),
							},
						},
					}, nil
				},
			}).ExpectInput("GetDistribution", &cloudfront.GetDistributionInput{
				Id: String("my-distribution-to-update"),
			}).ExpectCommandResult("etag-already-enabled").ExpectCalls("GetDistribution").Run(t)
		})

		t.Run("to enable", func(t *testing.T) {
			Template("update distribution id=my-distribution-to-update enable=true").Mock(&cloudfrontMock{
				GetDistributionFunc: func(input *cloudfront.GetDistributionInput) (*cloudfront.GetDistributionOutput, error) {
					return &cloudfront.GetDistributionOutput{
						Distribution: &cloudfront.Distribution{
							DistributionConfig: &cloudfront.DistributionConfig{
								Enabled: Bool(false),
							},
						},
						ETag: String("etag-id"),
					}, nil
				},
				UpdateDistributionFunc: func(input *cloudfront.UpdateDistributionInput) (*cloudfront.UpdateDistributionOutput, error) {
					return &cloudfront.UpdateDistributionOutput{
						ETag: String("etag-after-update"),
					}, nil
				},
			}).ExpectInput("GetDistribution", &cloudfront.GetDistributionInput{
				Id: String("my-distribution-to-update"),
			}).ExpectInput("UpdateDistribution", &cloudfront.UpdateDistributionInput{
				IfMatch: String("etag-id"),
				Id:      String("my-distribution-to-update"),
				DistributionConfig: &cloudfront.DistributionConfig{
					Enabled: Bool(true),
				},
			}).ExpectCommandResult("etag-after-update").ExpectCalls("GetDistribution", "UpdateDistribution").Run(t)
		})

		t.Run("already enabled and change other params", func(t *testing.T) {
			Template("update distribution id=my-distribution-to-update origin-domain=my.test.domain.com certificate=arn:of:the:certificate comment='useless comment' default-file=index.go "+
				"domain-aliases=any.domain.com,other.domain.com enable=true forward-cookies=whitelist forward-queries=true https-behaviour=redirect-to-https "+
				"origin-path=/my/custom/path price-class=PriceClass_All min-ttl=42").Mock(&cloudfrontMock{
				GetDistributionFunc: func(input *cloudfront.GetDistributionInput) (*cloudfront.GetDistributionOutput, error) {
					return &cloudfront.GetDistributionOutput{
						ETag: String("etag-already-enabled"),
						Distribution: &cloudfront.Distribution{
							DistributionConfig: &cloudfront.DistributionConfig{
								Enabled: Bool(true),
							},
						},
					}, nil
				},
				UpdateDistributionFunc: func(input *cloudfront.UpdateDistributionInput) (*cloudfront.UpdateDistributionOutput, error) {
					return &cloudfront.UpdateDistributionOutput{
						ETag: String("etag-updated-distribution"),
					}, nil
				},
			}).ExpectInput("GetDistribution", &cloudfront.GetDistributionInput{
				Id: String("my-distribution-to-update"),
			}).ExpectInput("UpdateDistribution", &cloudfront.UpdateDistributionInput{
				Id:      String("my-distribution-to-update"),
				IfMatch: String("etag-already-enabled"),
				DistributionConfig: &cloudfront.DistributionConfig{
					Origins: &cloudfront.Origins{
						Items: []*cloudfront.Origin{
							{
								DomainName: String("my.test.domain.com"),
								Id:         String("orig_1"),
								OriginPath: String("/my/custom/path"),
							},
						},
						Quantity: Int64(1),
					},
					ViewerCertificate: &cloudfront.ViewerCertificate{
						ACMCertificateArn: String("arn:of:the:certificate"),
						SSLSupportMethod:  String("sni-only"),
					},
					Comment:           String("useless comment"),
					DefaultRootObject: String("index.go"),
					Aliases: &cloudfront.Aliases{
						Items:    []*string{String("any.domain.com"), String("other.domain.com")},
						Quantity: Int64(2),
					},
					Enabled: Bool(true),
					DefaultCacheBehavior: &cloudfront.DefaultCacheBehavior{
						ForwardedValues: &cloudfront.ForwardedValues{
							Cookies: &cloudfront.CookiePreference{
								Forward: String("whitelist"),
							},
							QueryString: Bool(true),
						},
						ViewerProtocolPolicy: String("redirect-to-https"),
						MinTTL:               Int64(42),
					},
					PriceClass: String("PriceClass_All"),
				},
			}).ExpectCommandResult("etag-updated-distribution").ExpectCalls("GetDistribution", "UpdateDistribution").Run(t)
		})
	})

	t.Run("check", func(t *testing.T) {
		Template("check distribution id=my-distribution-id state=deployed timeout=1").Mock(&cloudfrontMock{
			GetDistributionFunc: func(input *cloudfront.GetDistributionInput) (*cloudfront.GetDistributionOutput, error) {
				return &cloudfront.GetDistributionOutput{
					Distribution: &cloudfront.Distribution{
						Status: String("deployed"),
					},
				}, nil
			}}).ExpectInput("GetDistribution", &cloudfront.GetDistributionInput{
			Id: String("my-distribution-id"),
		}).ExpectCalls("GetDistribution").Run(t)
	})

	t.Run("delete", func(t *testing.T) {
		Template("delete distribution id=my-distribution-to-delete").Mock(&cloudfrontMock{
			GetDistributionFunc: func(input *cloudfront.GetDistributionInput) (*cloudfront.GetDistributionOutput, error) {
				return &cloudfront.GetDistributionOutput{
					Distribution: &cloudfront.Distribution{
						Status: String("Deployed"),
						DistributionConfig: &cloudfront.DistributionConfig{
							Enabled: Bool(false),
						},
					},
					ETag: String("etag-distrib-to-delete"),
				}, nil
			},
			DeleteDistributionFunc: func(input *cloudfront.DeleteDistributionInput) (*cloudfront.DeleteDistributionOutput, error) {
				return nil, nil
			},
		}).ExpectInput("GetDistribution", &cloudfront.GetDistributionInput{
			Id: String("my-distribution-to-delete"),
		}).ExpectInput("GetDistribution", &cloudfront.GetDistributionInput{
			Id: String("my-distribution-to-delete"),
		}).ExpectInput("DeleteDistribution", &cloudfront.DeleteDistributionInput{
			Id:      String("my-distribution-to-delete"),
			IfMatch: String("etag-distrib-to-delete"),
		}).ExpectCalls("GetDistribution", "GetDistribution", "DeleteDistribution").Run(t)
	})
}
