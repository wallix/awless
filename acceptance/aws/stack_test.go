package awsat

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

func TestStack(t *testing.T) {
	_, tplFilePath, tplClean := generateTmpFile("tpl body content")
	defer tplClean()

	_, polFilePath, polClean := generateTmpFile("policy content")
	defer polClean()

	t.Run("create", func(t *testing.T) {
		Template("create stack name=new-stack template-file="+tplFilePath+" tags=Env:Prod,Dept:IT capabilities=one,two disable-rollback=true notifications=none,ntwo on-failure=done parameters=1:pone,2:ptwo resource-types=rone,rtwo role=donjuan policy-file="+polFilePath+" rollback-monitoring-min=2 rollback-triggers=[arn1,arn2] timeout=180").Mock(&cloudformationMock{
			CreateStackFunc: func(input *cloudformation.CreateStackInput) (*cloudformation.CreateStackOutput, error) {
				return &cloudformation.CreateStackOutput{StackId: String("new-stack-id")}, nil
			}}).ExpectInput("CreateStack", &cloudformation.CreateStackInput{
			StackName:        String("new-stack"),
			TemplateBody:     String("tpl body content"),
			Capabilities:     []*string{String("one"), String("two")},
			DisableRollback:  Bool(true),
			NotificationARNs: []*string{String("none"), String("ntwo")},
			OnFailure:        String("done"),
			Parameters:       []*cloudformation.Parameter{{ParameterKey: String("1"), ParameterValue: String("pone")}, {ParameterKey: String("2"), ParameterValue: String("ptwo")}},
			ResourceTypes:    []*string{String("rone"), String("rtwo")},
			RoleARN:          String("donjuan"),
			StackPolicyBody:  String("policy content"),
			TimeoutInMinutes: Int64(180),
			Tags:             []*cloudformation.Tag{{Key: String("Env"), Value: String("Prod")}, {Key: String("Dept"), Value: String("IT")}},
			RollbackConfiguration: &cloudformation.RollbackConfiguration{
				MonitoringTimeInMinutes: Int64(2),
				RollbackTriggers: []*cloudformation.RollbackTrigger{
					{Arn: String("arn1"), Type: aws.String("AWS::CloudWatch::Alarm")},
					{Arn: String("arn2"), Type: aws.String("AWS::CloudWatch::Alarm")},
				},
			},
		}).ExpectCommandResult("new-stack-id").ExpectCalls("CreateStack").Run(t)
	})

	t.Run("update", func(t *testing.T) {
		_, polUpdateFilePath, clean := generateTmpFile("update policy content")
		defer clean()

		Template("update stack name=other-name template-file="+tplFilePath+" use-previous-template=true tags=Env:Prod,Dept:IT capabilities=one,two notifications=none,ntwo parameters=1:pone,2:ptwo resource-types=rone,rtwo role=donjuan policy-file="+polFilePath+" policy-update-file="+polUpdateFilePath+" rollback-monitoring-min=2 rollback-triggers=[arn1,arn2]").Mock(&cloudformationMock{
			UpdateStackFunc: func(input *cloudformation.UpdateStackInput) (*cloudformation.UpdateStackOutput, error) {
				return &cloudformation.UpdateStackOutput{StackId: String("any-stack-id")}, nil
			}}).ExpectInput("UpdateStack", &cloudformation.UpdateStackInput{
			StackName:                   String("other-name"),
			TemplateBody:                String("tpl body content"),
			Capabilities:                []*string{String("one"), String("two")},
			NotificationARNs:            []*string{String("none"), String("ntwo")},
			Parameters:                  []*cloudformation.Parameter{{ParameterKey: String("1"), ParameterValue: String("pone")}, {ParameterKey: String("2"), ParameterValue: String("ptwo")}},
			ResourceTypes:               []*string{String("rone"), String("rtwo")},
			RoleARN:                     String("donjuan"),
			StackPolicyBody:             String("policy content"),
			StackPolicyDuringUpdateBody: String("update policy content"),
			UsePreviousTemplate:         Bool(true),
			Tags:                        []*cloudformation.Tag{{Key: String("Env"), Value: String("Prod")}, {Key: String("Dept"), Value: String("IT")}},
			RollbackConfiguration: &cloudformation.RollbackConfiguration{
				MonitoringTimeInMinutes: Int64(2),
				RollbackTriggers: []*cloudformation.RollbackTrigger{
					{Arn: String("arn1"), Type: aws.String("AWS::CloudWatch::Alarm")},
					{Arn: String("arn2"), Type: aws.String("AWS::CloudWatch::Alarm")},
				},
			},
		}).ExpectCommandResult("any-stack-id").ExpectCalls("UpdateStack").Run(t)
	})

	t.Run("delete", func(t *testing.T) {
		Template("delete stack name=any-stack-name retain-resources=1,2").Mock(&cloudformationMock{
			DeleteStackFunc: func(input *cloudformation.DeleteStackInput) (*cloudformation.DeleteStackOutput, error) {
				return nil, nil
			}}).ExpectInput("DeleteStack", &cloudformation.DeleteStackInput{
			StackName:       String("any-stack-name"),
			RetainResources: []*string{String("1"), String("2")},
		}).ExpectCalls("DeleteStack").Run(t)
	})

	_, stackFileYMLPath, stackFileYMLClean := generateTmpFileWithName(`
Parameters:
  Test1: 1
  Test2: 2
  Test3: 3
Tags:
  Tag1: 1
  Tag2: 2
  Tag3: 3
StackPolicy:
  Statement:
  - Effect: Allow
    Resource: "*"
`, "stackfile.yml")

	defer stackFileYMLClean()

	t.Run("update", func(t *testing.T) {
		Template("update stack name=some-stack template-file="+tplFilePath+" stack-file="+stackFileYMLPath+" parameters=Test1:a,Test2:b tags=Tag1:a,Tag2:b policy-file="+polFilePath).Mock(&cloudformationMock{
			UpdateStackFunc: func(input *cloudformation.UpdateStackInput) (*cloudformation.UpdateStackOutput, error) {
				return &cloudformation.UpdateStackOutput{StackId: String("any-stack-id")}, nil
			}}).ExpectInput("UpdateStack", &cloudformation.UpdateStackInput{
			StackName:       String("some-stack"),
			TemplateBody:    String("tpl body content"),
			Parameters:      []*cloudformation.Parameter{{ParameterKey: String("Test1"), ParameterValue: String("a")}, {ParameterKey: String("Test2"), ParameterValue: String("b")}, {ParameterKey: String("Test3"), ParameterValue: String("3")}},
			Tags:            []*cloudformation.Tag{{Key: String("Tag1"), Value: String("a")}, {Key: String("Tag2"), Value: String("b")}, {Key: String("Tag3"), Value: String("3")}},
			StackPolicyBody: String("policy content"),
		}).ExpectCalls("UpdateStack").Run(t)
	})

	_, stackFileJSONPath, stackFileJSONClean := generateTmpFileWithName(`{"Parameters":{"Test1":"1","Test2":"2","Test3":"3"},"Tags":{"Tag1":"1","Tag2":"2","Tag3":"3"},"StackPolicy":{"Statement":[{"Effect":"Allow","Resource":"*"}]}}`, "stackfile.json")

	defer stackFileJSONClean()

	t.Run("update", func(t *testing.T) {
		Template("update stack name=some-stack template-file="+tplFilePath+" stack-file="+stackFileJSONPath+" parameters=Test1:a,Test2:b tags=Tag1:a,Tag2:b").Mock(&cloudformationMock{
			UpdateStackFunc: func(input *cloudformation.UpdateStackInput) (*cloudformation.UpdateStackOutput, error) {
				return &cloudformation.UpdateStackOutput{StackId: String("any-stack-id")}, nil
			}}).ExpectInput("UpdateStack", &cloudformation.UpdateStackInput{
			StackName:       String("some-stack"),
			TemplateBody:    String("tpl body content"),
			Parameters:      []*cloudformation.Parameter{{ParameterKey: String("Test1"), ParameterValue: String("a")}, {ParameterKey: String("Test2"), ParameterValue: String("b")}, {ParameterKey: String("Test3"), ParameterValue: String("3")}},
			Tags:            []*cloudformation.Tag{{Key: String("Tag1"), Value: String("a")}, {Key: String("Tag2"), Value: String("b")}, {Key: String("Tag3"), Value: String("3")}},
			StackPolicyBody: String(`{"Statement":[{"Effect":"Allow","Resource":"*"}]}`),
		}).ExpectCalls("UpdateStack").Run(t)
	})

}

func generateTmpFileWithName(content, filename string) (*os.File, string, func()) {
	dir, err := ioutil.TempDir("", "awless-at-tmpdir")
	if err != nil {
		panic(err)
	}

	tmpfn := filepath.Join(dir, filename)
	if err := ioutil.WriteFile(tmpfn, []byte(content), 0666); err != nil {
		panic(err)
	}

	file, err := os.Open(tmpfn)
	if err != nil {
		panic(err)
	}

	cleanup := func() {
		file.Close()
		if err := os.RemoveAll(dir); err != nil {
			panic(err)
		}
	}

	return file, file.Name(), cleanup
}
