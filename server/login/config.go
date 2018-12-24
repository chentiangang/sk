package main

import (
	"io/ioutil"
	"os/user"
	"path"

	"golang.org/x/crypto/ssh"
	yaml "gopkg.in/yaml.v2"
)

type Node struct {
	Name       string  `json:"name"`
	Host       string  `json:"host"`
	User       string  `json:"user"`
	Port       int     `json:"port"`
	KeyPath    string  `json:"keypath"`
	Passphrase string  `json:"passphrase"`
	Password   string  `json:"password"`
	Children   []*Node `json:"children"`
}

func (n *Node) String() string {
	return n.Name
}

func (n *Node) user() string {
	if n.User == "" {
		return "root"
	}
	return n.User
}

func (n *Node) port() int {
	if n.Port <= 0 {
		return 22
	}
	return n.Port
}

func (n *Node) password() ssh.AuthMethod {
	if n.Password == "" {
		return nil
	}
	return ssh.Password(n.Password)
}

var cfg []*Node

func GetConfig() []*Node {
	return cfg
}

func LoadConfig() error {
	u, err := user.Current()
	if err != nil {
		return err
	}

	b, err := ioutil.ReadFile(path.Join(u.HomeDir, ".sk.yml"))
	if err != nil {
		return err
	}

	var c []*Node
	err = yaml.Unmarshal(b, &c)
	if err != nil {
		return err
	}
	cfg = c
	return nil
}
