package awsdriver

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws/credentials"
)

var (
	AWSCredDir      = filepath.Join(os.Getenv("HOME"), ".aws")
	AWSCredFilepath = filepath.Join(AWSCredDir, "credentials")
)

type credentialsPrompter struct {
	Profile string
	Val     credentials.Value
}

func NewCredsPrompter(profile string) *credentialsPrompter {
	return &credentialsPrompter{Profile: profile}
}

func (c *credentialsPrompter) Prompt() error {
	token := "and choose a profile name"
	if c.HasProfile() {
		token = fmt.Sprintf("for profile '%s'", c.Profile)
	}
	fmt.Printf("\nPlease enter access keys %s (stored at %s):\n", token, AWSCredFilepath)

	promptUntilNonEmpty("AWS Access Key ID? ", &c.Val.AccessKeyID)
	promptUntilNonEmpty("AWS Secret Access Key? ", &c.Val.SecretAccessKey)
	if c.HasProfile() {
		promptToOverride(fmt.Sprintf("Change your profile name (or just press Enter to keep '%s')?", c.Profile), &c.Profile)
	} else {
		c.Profile = "default"
		promptToOverride("Choose a profile name (or just press Enter to have AWS 'default')? ", &c.Profile)
	}

	return nil
}

func (c *credentialsPrompter) Store() (bool, error) {
	var created bool

	if c.Val.SecretAccessKey == "" {
		return created, errors.New("given empty secret access key")
	}
	if c.Val.AccessKeyID == "" {
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

	if _, err := fmt.Fprintf(f, "\n[%s]\naws_access_key_id = %s\naws_secret_access_key = %s\n", c.Profile, c.Val.AccessKeyID, c.Val.SecretAccessKey); err != nil {
		return created, err
	}

	return created, nil
}

func (c *credentialsPrompter) HasProfile() bool {
	return strings.TrimSpace(c.Profile) != ""
}

func promptToOverride(question string, v *string) {
	fmt.Print(question)
	var override string
	fmt.Scanln(&override)
	if strings.TrimSpace(override) != "" {
		*v = override
		return
	}
}

func promptUntilNonEmpty(question string, v *string) {
	ask := func(v *string) bool {
		fmt.Print(question)
		_, err := fmt.Scanln(v)
		if err == nil && strings.TrimSpace(*v) != "" {
			return false
		}
		if err != nil {
			fmt.Printf("Error: %s. Retry please...\n", err)
		}
		return true
	}
	for ask(v) {}
	return
}

func credentialsDirMissing() bool {
	_, err := os.Stat(AWSCredDir)
	return os.IsNotExist(err)
}
