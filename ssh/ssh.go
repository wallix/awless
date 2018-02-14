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
	Proxy                   *Client
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

func (c *Client) DialWithUsers(usernames ...string) error {
	var err error
	var client *gossh.Client

	hostport := fmt.Sprintf("%s:%d", c.IP, c.Port)

	for _, user := range usernames {
		newConfig := *c.Config
		newConfig.User = user
		if !c.StrictHostKeyChecking {
			newConfig.HostKeyCallback = gossh.InsecureIgnoreHostKey()
		}
		client, err = gossh.Dial("tcp", hostport, &newConfig)
		if err != nil {
			continue
		} else {
			c.logger.ExtraVerbosef("dialed %s successfully with user %s", hostport, user)
			c.User = user
			c.Client = client
			return nil
		}
	}

	return fmt.Errorf("unable to authenticate to %s for users %q. Last error: %s", hostport, usernames, err)
}

func (c *Client) NewClientWithProxy(destinationHost string, destinationPort int, usernames ...string) (*Client, error) {
	hostport := fmt.Sprintf("%s:%d", destinationHost, destinationPort)
	for _, user := range usernames {
		netConn, err := c.Dial("tcp", hostport)
		if err != nil {
			return nil, fmt.Errorf("cannot dial from %s:%d to %s:%d - %s", c.IP, c.Port, destinationHost, destinationPort, err)
		}
		c.logger.ExtraVerbosef("successful tcp connection from %s:%d to %s:%d", c.IP, c.Port, destinationHost, destinationPort)
		newConfig := *c.Config
		newConfig.User = user
		if !c.StrictHostKeyChecking {
			newConfig.HostKeyCallback = gossh.InsecureIgnoreHostKey()
		}
		conn, chans, reqs, err := gossh.NewClientConn(netConn, hostport, &newConfig)
		if err != nil {
			netConn.Close()
			c.logger.ExtraVerbosef("cannot proxy with user %s (err: %s)", user, err)
			continue
		}
		c.logger.ExtraVerbosef("proxied successfully with user %s", user)

		return &Client{
			Client:  gossh.NewClient(conn, chans, reqs),
			Proxy:   c,
			IP:      destinationHost,
			User:    user,
			Keypath: c.Keypath,
			Port:    destinationPort,
			InteractiveTerminalFunc: func(*gossh.Client) error { return nil },
			StrictHostKeyChecking:   c.StrictHostKeyChecking,
			logger:                  logger.DiscardLogger,
		}, nil
	}

	return nil, fmt.Errorf("cannot proxy from %s:%d to %s:%d with users %q", c.IP, c.Port, destinationHost, destinationPort, usernames)
}

func (c *Client) CloseAll() error {
	if c != nil {
		if c.Client != nil {
			return c.Client.Close()
		}
		if c.Proxy != nil {
			return c.Proxy.Close()
		}
	}
	return nil
}

func (c *Client) Connect() (err error) {
	args, installed := c.localExec()
	if installed {
		c.logger.Infof("Login as '%s' on '%s'; client '%s'", c.User, c.IP, args[0])
		c.logger.ExtraVerbosef("running locally %s", args)
		if err := c.CloseAll(); err != nil {
			c.logger.Warning("could not close properly SSH awless client before delegating")
		}
		if c.Proxy != nil {
			return workaroundExeCVEThroughScript(args)
		}
		return syscall.Exec(args[0], args, os.Environ())
	}

	c.logger.Infof("No SSH. Fallback on builtin client. Login as '%s' on '%s'", c.User, c.IP)
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
	if c.Proxy != nil {
		var keyArg string
		if k := c.Proxy.Keypath; len(k) > 0 {
			keyArg = fmt.Sprintf("-i %s", k)
		}
		extraOpts["ProxyCommand"] = fmt.Sprintf("ssh %s %s@%s -p %d -W %%h:%%p", keyArg, c.Proxy.User, c.Proxy.IP, c.Proxy.Port)
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
	args := []string{bin}
	if len(c.Keypath) > 0 {
		args = append(args, "-i", c.Keypath)
	}
	if c.Port != 22 {
		args = append(args, "-p", strconv.Itoa(c.Port))
	}
	if !c.StrictHostKeyChecking {
		args = append(args, "-o", "StrictHostKeychecking=no")
	}

	args = append(args, fmt.Sprintf("%s@%s", c.User, c.IP))

	if c.Proxy != nil {
		var keyArg string
		if k := c.Proxy.Keypath; len(k) > 0 {
			keyArg = fmt.Sprintf("-i %s", k)
		}
		args = append(args, "-o", fmt.Sprintf("ProxyCommand='ssh %s %s@%s -p %d -W %%h:%%p'", keyArg, c.Proxy.User, c.Proxy.IP, c.Proxy.Port))
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

const tmpProxyCommandScriptFilename = "awless-ssh-proxycommand"

// This hack is used to circumvent a bug i cannot yet figure out
// Bug: when executing syscall.Exec(args[0], args, os.Environ()) and args contains
// the proxy command (typically args := []string{"/usr/bin/ssh", "ec2-user@172.31.78.138", "-o", "StrictHostKeychecking=no", "-o", "ProxyCommand='ssh ec2-user@52.26.181.76 -W [%h]:%p'"}
// we get an error like (in Go, Python):
//     /bin/bash: 1: exec: ssh ec2-user@52.26.181.76 -W [172.31.78.138]:22: not found
//     ssh_exchange_identification: Connection closed by remote host
//
// Since execve(2) can take as the first argument a filename, the workaround is to use
// a temporary script to execute this command.
//
// Note that the file cannot be removed since we syscall for another process. So the first time
// it is created and after that only truncated (reuse the same file)
func workaroundExeCVEThroughScript(args []string) error {
	fpath := filepath.Join(os.TempDir(), tmpProxyCommandScriptFilename)
	logger.ExtraVerbosef("using script %s", fpath)
	tmpExec, err := os.OpenFile(fpath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0700)
	if err != nil {
		return err
	}

	script := fmt.Sprintf("#! /bin/bash\n%s", strings.Join(args, " "))
	if _, err := tmpExec.Write([]byte(script)); err != nil {
		return err
	}
	if err := tmpExec.Close(); err != nil {
		return err
	}
	return syscall.Exec(tmpExec.Name(), []string{}, os.Environ())
}
