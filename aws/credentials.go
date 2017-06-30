package aws

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws/credentials"
)

var (
	AWSCredDir      = filepath.Join(os.Getenv("HOME"), ".aws")
	AWSCredFilepath = filepath.Join(AWSCredDir, "credentials")
)

type credentialsPrompter struct {
	profile string
	val     credentials.Value
}

func NewCredsPrompter(profile string) *credentialsPrompter {
	return &credentialsPrompter{profile: profile}
}

func (c *credentialsPrompter) Prompt() error {
	fmt.Printf("\nPlease enter access keys for profile '%s' (stored at %s):\n", c.profile, AWSCredFilepath)
	fmt.Print("AWS Access Key ID? ")
	if _, err := fmt.Scanln(&c.val.AccessKeyID); err != nil {
		return err
	}
	fmt.Print("AWS Secret Access Key? ")
	if _, err := fmt.Scanln(&c.val.SecretAccessKey); err != nil {
		return err
	}

	return nil
}

func (c *credentialsPrompter) Store() (bool, error) {
	var created bool

	if c.val.SecretAccessKey == "" {
		return created, errors.New("given empty secret access key")
	}
	if c.val.AccessKeyID == "" {
		return created, errors.New("given empty access key")
	}

	if credentialsDirMissing() {
		if err := os.MkdirAll(AWSCredDir, 0700); err != nil {
			return created, fmt.Errorf("creating '%s' : %s", AWSCredDir, err)
		}
		created = true
	}

	f, err := os.OpenFile(AWSCredFilepath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return created, fmt.Errorf("appending to '%s': %s", AWSCredFilepath, err)
	}

	if _, err := fmt.Fprintf(f, "[%s]\naws_access_key_id = %s\naws_secret_access_key = %s\n", c.profile, c.val.AccessKeyID, c.val.SecretAccessKey); err != nil {
		return created, err
	}

	return created, nil
}

func credentialsDirMissing() bool {
	_, err := os.Stat(AWSCredDir)
	return os.IsNotExist(err)
}
