package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"time"

	"golang.org/x/crypto/ssh/terminal"

	"golang.org/x/crypto/ssh"
)

var (
	DefaultCiphers = []string{
		"aes128-ctr",
	}
)

type Client interface {
	Login()
}

type defaultClient struct {
	clientConfig *ssh.ClientConfig
	node         *Node
}

func NewClient(node *Node) Client {
	u, err := user.Current()
	if err != nil {
		l.Error(err)
		return nil
	}

	var authMethods []ssh.AuthMethod
	var pemBytes []byte
	if node.KeyPath == "" {
		pemBytes, err = ioutil.ReadFile(path.Join(u.HomeDir, ".ssh/id_dsa"))
	} else {
		pemBytes, err = ioutil.ReadFile(node.KeyPath)
	}

	if err != nil {
		l.Error(err)
	} else {
		var signer ssh.Signer
		if node.Passphrase != "" {
			signer, err = ssh.ParsePrivateKeyWithPassphrase(pemBytes, []byte(node.Passphrase))
		} else {
			signer, err = ssh.ParsePrivateKey(pemBytes)
		}
		if err != nil {
			l.Error(err)
		} else {
			authMethods = append(authMethods, ssh.PublicKeys(signer))
		}
	}
	password := node.password()
	if password != nil {
		interactive := func(user, instruction string, questions []string, echos []bool) (answers []string, err error) {
			answers = make([]string, len(questions))
			for n := range questions {
				answers[n] = node.Password
			}
			return answers, nil
		}
		authMethods = append(authMethods, ssh.KeyboardInteractive(interactive))
		authMethods = append(authMethods, password)
	}
	scfg := &ssh.ClientConfig{
		User:            node.user(),
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Second * 10,
	}
	scfg.SetDefaults()
	scfg.Ciphers = append(scfg.Ciphers, DefaultCiphers...)

	return &defaultClient{
		clientConfig: scfg,
		node:         node,
	}
}

func (c *defaultClient) Login() {
	host := c.node.Host
	port := c.node.port()
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", host, port), c.clientConfig)
	if err != nil {
		l.Error(err)
		return
	}
	defer client.Close()
	l.Infof("Connect server ssh -p %d %s@%s version: %s\n", port, c.node.user(), host, string(client.ServerVersion()))
	session, err := client.NewSession()
	if err != nil {
		l.Error(err)
		return
	}
	defer session.Close()
	fd := int(os.Stdin.Fd())
	state, err := terminal.MakeRaw(fd)
	if err != nil {
		l.Error(err)
		return
	}
	defer terminal.Restore(fd, state)

	w, h, err := terminal.GetSize(fd)
	if err != nil {
		l.Error(err)
		return
	}
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	err = session.RequestPty("xterm", h, w, modes)
	if err != nil {
		l.Error(err)
		return
	}

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	err = session.Shell()

	if err != nil {
		l.Error(err)
		return
	}

	// interval get terminal size
	// fix resize issue
	go func() {
		var (
			ow = w
			oh = h
		)
		for {
			cw, ch, err := terminal.GetSize(fd)
			if err != nil {
				break
			}
			if cw != ow || ch != oh {
				err = session.WindowChange(ch, cw)
				if err != nil {
					break
				}
				ow = cw
				oh = ch
			}
			time.Sleep(time.Second)
		}
	}()
	session.Wait()
}
