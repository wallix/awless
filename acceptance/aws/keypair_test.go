package awsat

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/wallix/awless/console"
)

func TestKeypair(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		tmpFolder, err := ioutil.TempDir("", "tmpfolder")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(tmpFolder)
		console.GenerateSSHKeyPair = func(size int, encryptKey bool) ([]byte, []byte, error) {
			var public string
			if encryptKey {
				public = "encrypted keypair"
			} else {
				public = "unencrypted keypair"
			}
			return []byte(public), []byte{}, nil
		}
		os.Setenv("__AWLESS_KEYS_DIR", tmpFolder)

		Template("create keypair name=my-kp encrypted=true").
			Mock(&ec2Mock{
				ImportKeyPairFunc: func(param0 *ec2.ImportKeyPairInput) (*ec2.ImportKeyPairOutput, error) {
					return &ec2.ImportKeyPairOutput{KeyName: String("my-kp")}, nil
				},
			}).ExpectInput("ImportKeyPair", &ec2.ImportKeyPairInput{
			KeyName:           String("my-kp"),
			PublicKeyMaterial: []byte("encrypted keypair"),
		}).
			ExpectCommandResult("my-kp").ExpectCalls("ImportKeyPair").Run(t)
	})

	t.Run("delete", func(t *testing.T) {
		Template("delete keypair name=kp-to-delete").
			Mock(&ec2Mock{
				DeleteKeyPairFunc: func(param0 *ec2.DeleteKeyPairInput) (*ec2.DeleteKeyPairOutput, error) {
					return nil, nil
				},
			}).ExpectInput("DeleteKeyPair", &ec2.DeleteKeyPairInput{KeyName: String("kp-to-delete")}).
			ExpectCalls("DeleteKeyPair").Run(t)
	})
}
