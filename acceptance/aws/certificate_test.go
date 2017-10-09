package awsat

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/acm"
)

func TestCertificate(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		tcases := []struct {
			template            string
			expCertificateInput *acm.RequestCertificateInput
		}{
			{
				template: "create certificate domains=my.domain.1,my.domain.2,my.domain.3 validation-domains=domain.1,2",
				expCertificateInput: &acm.RequestCertificateInput{
					DomainName:              String("my.domain.1"),
					SubjectAlternativeNames: []*string{String("my.domain.2"), String("my.domain.3")},
					DomainValidationOptions: []*acm.DomainValidationOption{
						{DomainName: String("my.domain.1"), ValidationDomain: String("domain.1")},
						{DomainName: String("my.domain.2"), ValidationDomain: String("2")},
					},
				},
			},
			{
				template: "create certificate domains=my.domain.1,my.domain.2,my.domain.3 validation-domains=my.domain.1",
				expCertificateInput: &acm.RequestCertificateInput{
					DomainName:              String("my.domain.1"),
					SubjectAlternativeNames: []*string{String("my.domain.2"), String("my.domain.3")},
					DomainValidationOptions: []*acm.DomainValidationOption{
						{DomainName: String("my.domain.1"), ValidationDomain: String("my.domain.1")},
					},
				},
			},
			{
				template: "create certificate domains=my.domain.1,my.domain.2 validation-domains=domain.1,domain.2",
				expCertificateInput: &acm.RequestCertificateInput{
					DomainName:              String("my.domain.1"),
					SubjectAlternativeNames: []*string{String("my.domain.2")},
					DomainValidationOptions: []*acm.DomainValidationOption{
						{DomainName: String("my.domain.1"), ValidationDomain: String("domain.1")},
						{DomainName: String("my.domain.2"), ValidationDomain: String("domain.2")},
					},
				},
			},
		}

		for _, tcase := range tcases {

			Template(tcase.template).
				Mock(&acmMock{
					RequestCertificateFunc: func(param0 *acm.RequestCertificateInput) (*acm.RequestCertificateOutput, error) {
						return &acm.RequestCertificateOutput{CertificateArn: String("arn:my:new:certificate")}, nil
					},
				}).ExpectInput("RequestCertificate", tcase.expCertificateInput).ExpectCommandResult("arn:my:new:certificate").ExpectCalls("RequestCertificate").Run(t)
		}
	})

	t.Run("delete", func(t *testing.T) {
		Template("delete certificate arn=arn:certificate:to:delete").
			Mock(&acmMock{
				DeleteCertificateFunc: func(param0 *acm.DeleteCertificateInput) (*acm.DeleteCertificateOutput, error) {
					return nil, nil
				},
			}).ExpectInput("DeleteCertificate", &acm.DeleteCertificateInput{CertificateArn: String("arn:certificate:to:delete")}).
			ExpectCalls("DeleteCertificate").Run(t)
	})

	t.Run("check", func(t *testing.T) {
		Template("check certificate arn=arn:certificate:to:check state=issued timeout=1").Mock(&acmMock{
			DescribeCertificateFunc: func(input *acm.DescribeCertificateInput) (*acm.DescribeCertificateOutput, error) {
				return &acm.DescribeCertificateOutput{Certificate: &acm.CertificateDetail{CertificateArn: String("arn:certificate:to:check"), Status: String("issued")}}, nil
			}}).ExpectInput("DescribeCertificate", &acm.DescribeCertificateInput{CertificateArn: String("arn:certificate:to:check")}).
			ExpectCalls("DescribeCertificate").Run(t)
	})
}
