package ssh

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"text/template"
	"time"

	"github.com/wallix/awless/logger"

	gossh "golang.org/x/crypto/ssh"
	agent "golang.org/x/crypto/ssh/agent"
	"golang.org/x/crypto/ssh/knownhosts"
	"golang.org/x/crypto/ssh/terminal"
)

type Client struct {
	*gossh.Client
	authmethod              gossh.AuthMethod
	Keypath, IP, User       string
	HostKeyCallback         gossh.HostKeyCallback
	StrictHostKeyChecking   bool
	InteractiveTerminalFunc func(*gossh.Client) error
	logger                  *logger.Logger
}

func InitClient(keyname string, keyFolders ...string) (*Client, error) {
	privkey, err := resolvePrivateKey(keyname, keyFolders...)
	if err != nil {
		return nil, err
	}

	authmethod, err := resolveAuthMethod(privkey)
	if err != nil {
		return nil, err
	}

	cli := &Client{
		logger:                  logger.DiscardLogger,
		InteractiveTerminalFunc: func(*gossh.Client) error { return nil },
		Keypath:                 privkey.path,
		authmethod:              authmethod,
		StrictHostKeyChecking:   true,
	}

	return cli, nil
}

func (c *Client) SetLogger(l *logger.Logger) {
	c.logger = l
}

func (c *Client) SetStrictHostKeyChecking(hostKeyChecking bool) {
	c.StrictHostKeyChecking = hostKeyChecking
}

func (c *Client) DialWithUsers(usernames ...string) (*Client, error) {
	var failures int
	for _, user := range usernames {
		c.logger.Verbosef("trying with user %s", user)
		client, err := gossh.Dial("tcp", c.IP+":22", c.buildClientConfig(user))
		if err == nil {
			c.User = user
			c.Client = client
			break
		}
		if err != nil && strings.Contains(err.Error(), "unable to authenticate") {
			failures++
			if len(usernames) == failures {
				return c, fmt.Errorf("with users %q: %s", usernames, err)
			}
			continue
		} else if err != nil {
			return c, err
		}
	}

	return c, nil
}

func (c *Client) Connect() error {
	defer func() {
		if c.Client != nil {
			c.Client.Close()
		}
	}()

	args, installed := c.localExec()
	if installed {
		c.logger.Infof("Login as '%s' on '%s', using keypair '%s' with ssh client '%s'\n", c.User, c.IP, c.Keypath, args[0])
		return syscall.Exec(args[0], args, os.Environ())
	}

	c.logger.Infof("No SSH. Fallback on builtin client. Login as '%s' on '%s', using keypair '%s'\n", c.User, c.IP, c.Keypath)
	return c.InteractiveTerminalFunc(c.Client)
}

func (c *Client) buildClientConfig(username string) *gossh.ClientConfig {
	config := &gossh.ClientConfig{
		User:            username,
		Auth:            []gossh.AuthMethod{c.authmethod},
		Timeout:         2 * time.Second,
		HostKeyCallback: gossh.InsecureIgnoreHostKey(),
	}

	if c.StrictHostKeyChecking {
		config.HostKeyCallback = checkHostKey
	}

	return config
}

func (c *Client) SSHConfigString(hostname string) string {
	var buf bytes.Buffer

	params := struct {
		IP, User, Keypath, Name string
	}{c.IP, c.User, c.Keypath, hostname}

	template.Must(template.New("ssh_config").Parse(`
Host {{ .Name }}
	Hostname {{ .IP }}
	User {{ .User }}
	IdentityFile {{ .Keypath }}
`)).Execute(&buf, params)

	return buf.String()
}

func (c *Client) ConnectString() string {
	args, _ := c.localExec()
	return strings.Join(args, " ")
}

func (c *Client) localExec() ([]string, bool) {
	exists := true
	bin, err := exec.LookPath("ssh")
	if err != nil {
		exists = false
		bin = "ssh"
	}
	args := []string{bin, fmt.Sprintf("%s@%s", c.User, c.IP)}
	if len(c.Keypath) > 0 {
		args = append(args, "-i", c.Keypath)
	}
	if !c.StrictHostKeyChecking {
		args = append(args, "-o", "StrictHostKeychecking=no")
	}

	return args, exists
}

func DecryptSSHKey(key []byte, password []byte) (gossh.Signer, error) {
	block, _ := pem.Decode(key)
	pem, err := x509.DecryptPEMBlock(block, password)
	if err != nil {
		return nil, err
	}
	sshkey, err := x509.ParsePKCS1PrivateKey(pem)
	if err != nil {
		return nil, err
	}
	return gossh.NewSignerFromKey(sshkey)
}

func resolveAuthMethod(priv privateKey) (gossh.AuthMethod, error) {
	var authmethod gossh.AuthMethod

	if len(priv.body) > 0 {
		signer, err := gossh.ParsePrivateKey(priv.body)
		if err != nil && strings.Contains(err.Error(), "cannot decode encrypted private keys") {
			fmt.Fprintf(os.Stderr, "This SSH key is encrypted. Please enter passphrase for key '%s':", priv.path)
			var passphrase []byte
			passphrase, err = terminal.ReadPassword(int(syscall.Stdin))
			if err != nil {
				return nil, err
			}
			fmt.Fprintln(os.Stderr)
			signer, err = DecryptSSHKey(priv.body, passphrase)

			authmethod = gossh.PublicKeys(signer)
		}

		return authmethod, err
	} else {
		sshAuthSock := os.Getenv("SSH_AUTH_SOCK")
		if len(sshAuthSock) == 0 {
			return nil, fmt.Errorf("No key provided and no SSH_AUTH_SOCK env variable set, unable to resolve auth")
		}

		agentUnixSock, err := net.Dial("unix", sshAuthSock)
		if err != nil {
			return nil, err
		}
		authmethod = gossh.PublicKeysCallback(agent.NewClient(agentUnixSock).Signers)
	}

	return authmethod, nil
}

type privateKey struct {
	path string
	body []byte
}

func resolvePrivateKey(keyname string, keyFolders ...string) (priv privateKey, err error) {
	// if keyname is zero, assume that the agent will be used
	if len(keyname) == 0 {
		return
	}

	keyPaths := []string{
		keyname,
	}
	if !strings.HasPrefix(keyname, ".pem") {
		keyPaths = append(keyPaths, fmt.Sprintf("%s.pem", keyname))
	}
	for _, folder := range keyFolders {
		if filepath.IsAbs(keyname) {
			break
		}
		if _, err = os.Stat(folder); err != nil {
			continue
		}
		keyPaths = append(keyPaths, filepath.Join(folder, keyname))
		if !strings.HasPrefix(keyname, ".pem") {
			keyPaths = append(keyPaths, filepath.Join(folder, fmt.Sprintf("%s.pem", keyname)))
		}
	}

	for _, path := range keyPaths {
		priv.body, err = ioutil.ReadFile(path)
		if err == nil {
			priv.path = path
			return
		}
		if !os.IsNotExist(err) {
			return
		}
	}

	err = fmt.Errorf("cannot find SSH key '%s'. Searched at paths '%s'", keyname, strings.Join(keyPaths, "','"))

	return
}

func checkHostKey(hostname string, remote net.Addr, key gossh.PublicKey) error {
	var knownHostsFiles []string
	var fileToAddKnownKey string

	opensshFile := filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts")
	if _, err := os.Stat(opensshFile); err == nil {
		knownHostsFiles = append(knownHostsFiles, opensshFile)
		fileToAddKnownKey = opensshFile
	}

	awlessFile := filepath.Join(os.Getenv("__AWLESS_HOME"), "known_hosts")
	if _, err := os.Stat(awlessFile); err == nil {
		knownHostsFiles = append(knownHostsFiles, awlessFile)
	}
	if fileToAddKnownKey == "" {
		fileToAddKnownKey = awlessFile
	}

	checkKnownHostFunc, err := knownhosts.New(knownHostsFiles...)
	if err != nil {
		return err
	}
	knownhostsErr := checkKnownHostFunc(hostname, remote, key)
	keyError, ok := knownhostsErr.(*knownhosts.KeyError)
	if !ok {
		return knownhostsErr
	}
	if len(keyError.Want) == 0 {
		if trustKeyFunc(hostname, remote, key, fileToAddKnownKey) {
			f, err := os.OpenFile(fileToAddKnownKey, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
			if err != nil {
				return err
			}
			defer f.Close()
			_, err = f.WriteString(knownhosts.Line([]string{hostname}, key) + "\n")
			return err
		} else {
			return errors.New("Host public key verification failed.")
		}
	}

	var knownKeyInfos string
	var knownKeyFiles []string
	for _, knownKey := range keyError.Want {
		knownKeyInfos += fmt.Sprintf("\n-> %s (%s key in %s:%d)", gossh.FingerprintSHA256(knownKey.Key), knownKey.Key.Type(), knownKey.Filename, knownKey.Line)
		knownKeyFiles = append(knownKeyFiles, fmt.Sprintf("'%s:%d'", knownKey.Filename, knownKey.Line))
	}

	return fmt.Errorf(`
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
AWLESS DETECTED THAT THE REMOTE HOST PUBLIC KEY HAS CHANGED
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@

Host key for '%s' has changed and you did not disable strict host key checking.
Someone may be trying to intercept your connection (man-in-the-middle attack). Otherwise, the host key may have been changed.

The fingerprint for the %s key sent by the remote host is %s.
You persisted:%s

To get rid of this message, update %s`, hostname, key.Type(), gossh.FingerprintSHA256(key), knownKeyInfos, strings.Join(knownKeyFiles, ","))
}

var trustKeyFunc func(hostname string, remote net.Addr, key gossh.PublicKey, keyFileName string) bool = func(hostname string, remote net.Addr, key gossh.PublicKey, keyFileName string) bool {
	fmt.Printf("awless could not validate the authenticity of '%s' (unknown host)\n", hostname)
	fmt.Printf("%s public key fingerprint is %s.\n", key.Type(), gossh.FingerprintSHA256(key))
	fmt.Printf("Do you want to continue connecting and persist this key to '%s' (yes/no)? ", keyFileName)
	var yesorno string
	_, err := fmt.Scanln(&yesorno)
	if err != nil {
		return false
	}
	return strings.ToLower(yesorno) == "yes"
}
