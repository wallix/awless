package ssh

import (
	"fmt"
	"net"
	"os"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"golang.org/x/crypto/ssh/terminal"
)

func agentAuth() (ssh.AuthMethod, error) {
	sock, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK"))
	if err != nil {
		return nil, err
	}
	return ssh.PublicKeysCallback(agent.NewClient(sock).Signers), nil
}

func privateKeyAuth(priv privateKey) (ssh.AuthMethod, error) {
	signer, err := ssh.ParsePrivateKey(priv.body)
	if err != nil {
		if strings.Contains(err.Error(), "cannot decode encrypted private keys") {
			return encryptedPrivKeyAuth(priv)
		}
		return nil, err
	}
	return ssh.PublicKeys(signer), nil
}

func encryptedPrivKeyAuth(priv privateKey) (ssh.AuthMethod, error) {
	fmt.Fprintf(os.Stderr, "This SSH key is encrypted. Please enter passphrase for key '%s':", priv.path)
	passphrase, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return nil, err
	}
	fmt.Fprintln(os.Stderr)

	signer, err := DecryptSSHKey(priv.body, passphrase)
	if err != nil {
		return nil, err
	}
	return ssh.PublicKeys(signer), nil
}
