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
	"strconv"
	"strings"
	"syscall"
	"text/template"
	"time"

	"github.com/wallix/awless/logger"

	gossh "golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

type Client struct {
	*gossh.Client
	Config                  *gossh.ClientConfig
	IP, User, Keypath       string
	Port                    int
	HostKeyCallback         gossh.HostKeyCallback
	StrictHostKeyChecking   bool
	InteractiveTerminalFunc func(*gossh.Client) error
	logger                  *logger.Logger
}

func InitClient(keyname string, keyFolders ...string) (*Client, error) {
	var auths []gossh.AuthMethod

	privkey, ok := findPrivateKeyFromName(keyname, keyFolders...)
	if ok {
		if a, err := privateKeyAuth(privkey); err == nil {
			auths = append(auths, a)
		}
	}

	if a, err := agentAuth(); err == nil {
		auths = append(auths, a)
	}

	if len(auths) == 0 {
		return nil, fmt.Errorf("No key provided and no SSH_AUTH_SOCK env variable set, unable to resolve auth")
	}

	return &Client{
		Config: &gossh.ClientConfig{
			Auth:            auths,
			Timeout:         2 * time.Second,
			HostKeyCallback: checkHostKey,
		},
		Keypath:                 privkey.path,
		logger:                  logger.DiscardLogger,
		InteractiveTerminalFunc: func(*gossh.Client) error { return nil },
		StrictHostKeyChecking:   true,
	}, nil
}

func (c *Client) SetLogger(l *logger.Logger) {
	c.logger = l
}

func (c *Client) SetStrictHostKeyChecking(hostKeyChecking bool) {
	c.StrictHostKeyChecking = hostKeyChecking
}

func (c *Client) DialWithUsers(usernames ...string) (*Client, error) {
	var err error
	for _, user := range usernames {
		c.logger.ExtraVerbosef("trying with user %s", user)
		newConfig := *c.Config
		newConfig.User = user
		if !c.StrictHostKeyChecking {
			newConfig.HostKeyCallback = gossh.InsecureIgnoreHostKey()
		}
		client, err := gossh.Dial("tcp", fmt.Sprintf("%s:%d", c.IP, c.Port), &newConfig)
		if err != nil {
			return c, err
		} else {
			c.User = user
			c.Client = client
			return c, nil
		}
	}

	return nil, fmt.Errorf("unable to authenticate with any of those users %q. Last error: %s", usernames, err)
}

func (c *Client) Connect() error {
	defer func() {
		if c.Client != nil {
			c.Client.Close()
		}
	}()

	args, installed := c.localExec()
	if installed {
		c.logger.Infof("Login as '%s' on '%s'; client '%s'\n", c.User, c.IP, args[0])
		return syscall.Exec(args[0], args, os.Environ())
	}

	c.logger.Infof("No SSH. Fallback on builtin client. Login as '%s' on '%s', using %s for auth\n", c.User, c.IP)
	return c.InteractiveTerminalFunc(c.Client)
}

func (c *Client) SSHConfigString(hostname string) string {
	var buf bytes.Buffer

	extraOpts := map[string]string{}
	if len(c.Keypath) > 0 {
		extraOpts["IdentityFile"] = c.Keypath
	}
	if !c.StrictHostKeyChecking {
		extraOpts["StrictHostKeychecking"] = "no"
	}
	if c.Port != 22 {
		extraOpts["Port"] = strconv.Itoa(c.Port)
	}

	params := struct {
		IP, User, Name string
		Extra          map[string]string
	}{c.IP, c.User, hostname, extraOpts}

	template.Must(template.New("ssh_config").Parse(`
Host {{ .Name }}
	Hostname {{ .IP }}
	User {{ .User }}
{{- range $key, $value := .Extra }}
	{{ $key }} {{ $value -}}
{{ end -}}
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
	if c.Port != 22 {
		args = append(args, "-p", strconv.Itoa(c.Port))
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

type privateKey struct {
	path string
	body []byte
}

func findPrivateKeyFromName(keyname string, keyFolders ...string) (privateKey, bool) {
	var priv privateKey

	if len(keyname) == 0 {
		return priv, false
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
		if _, err := os.Stat(folder); err != nil {
			continue
		}
		keyPaths = append(keyPaths, filepath.Join(folder, keyname))
		if !strings.HasPrefix(keyname, ".pem") {
			keyPaths = append(keyPaths, filepath.Join(folder, fmt.Sprintf("%s.pem", keyname)))
		}
	}

	for _, path := range keyPaths {
		b, err := ioutil.ReadFile(path)
		if err == nil {
			priv.path = path
			priv.body = b
			return priv, true
		}
		if !os.IsNotExist(err) {
			return priv, false
		}
	}

	return priv, false
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
