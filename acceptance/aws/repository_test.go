package awsat

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/ecr"
)

func TestRepository(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		Template("create repository name=new-repo").Mock(&ecrMock{
			CreateRepositoryFunc: func(input *ecr.CreateRepositoryInput) (*ecr.CreateRepositoryOutput, error) {
				return &ecr.CreateRepositoryOutput{Repository: &ecr.Repository{RepositoryArn: String("new-repo-arn")}}, nil
			}}).
			ExpectInput("CreateRepository", &ecr.CreateRepositoryInput{
				RepositoryName: String("new-repo"),
			}).ExpectCommandResult("new-repo-arn").ExpectCalls("CreateRepository").Run(t)
	})

	t.Run("delete", func(t *testing.T) {
		Template("delete repository name=any-repo").Mock(&ecrMock{
			DeleteRepositoryFunc: func(input *ecr.DeleteRepositoryInput) (*ecr.DeleteRepositoryOutput, error) {
				return &ecr.DeleteRepositoryOutput{Repository: &ecr.Repository{RepositoryArn: String("any-repo-arn")}}, nil
			}}).
			ExpectInput("DeleteRepository", &ecr.DeleteRepositoryInput{
				RepositoryName: String("any-repo"),
			}).ExpectCalls("DeleteRepository").Run(t)
	})
}
