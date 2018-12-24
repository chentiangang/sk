package model

type CMDB struct {
	Name     Conf    `yaml:"name"`
	Children []*CMDB `yaml:"children"`
}

type Conf struct {
	Host      []*Host   `yaml:"host"`
	HostGroup HostGroup `yaml:"hostGroup"`
	User      User      `yaml:"user"`
	UserGroup UserGroup `yaml:"userGroup"`
}

type Host struct {
	Os    Os    `yaml:"os"`
	Hosts Hosts `yaml:"hosts"`
}

type Os struct {
	MemSize  string `yaml:"memSize"`
	DiskSize string `yaml:"diskSize"`
	Cpu      string `yaml:"cpu"`
}

type Hosts struct {
	HostName   string `yaml:"hostName"`
	Ip         string `yaml:"ip"`
	SysVersion string `yaml:"sysVersion"`
}

type HostGroup struct {
	GroupName string `yaml:"groupName"`
	HostIp    string `yaml:"hostIp"`
}

type User struct {
	Name string `yaml:"name"`
	Id   string `yaml:"id"`
}

type UserGroup struct {
	GroupName string `yaml:"groupName"`
	UserId    string `yaml:"userId"`
}
