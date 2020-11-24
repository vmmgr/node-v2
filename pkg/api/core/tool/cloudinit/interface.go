package cloudinit

// humstackのcloud-init部分を参考

type MetaData struct {
	InstanceID    string `yaml:"instance-id"`
	LocalHostName string `yaml:"local-hostname"`
}

type UsersData struct {
	Name              string   `yaml:"name"`
	Password          string   `yaml:"password"`
	Groups            string   `yaml:"groups"`
	Shell             string   `yaml:"shell"`
	Sudo              []string `yaml:"sudo"`
	SSHAuthorizedKeys []string `yaml:"ssh-authorized-keys"`
	SSHPWAuth         bool     `yaml:"ssh_pwauth"`
	LockPasswd        bool     `yaml:"lock_passwd"`
}

type UserData struct {
	Password  string `yaml:"password"`
	ChPasswd  string `yaml:"chpasswd"`
	SshPwAuth bool   `yaml:"ssh_pwauth"`
}

type NetworkConfigSubnetType string

const (
	NetworkConfigSubnetTypeStatic NetworkConfigSubnetType = "static"
)

type NetworkConfigSubnet struct {
	Type    NetworkConfigSubnetType `yaml:"type"`
	Address string                  `yaml:"address"`
	Netmask string                  `yaml:"netmask"`
	Gateway string                  `yaml:"gateway"`
	DNS     []string                `yaml:"dns_nameservers"`
}

type NetworkConfigType string

const (
	NetworkConfigTypePhysical NetworkConfigType = "physical"
)

type NetworkCon struct {
	Version int32           `yaml:"version"`
	Config  []NetworkConfig `yaml:"config"`
}

type NetworkConfig struct {
	Type       NetworkConfigType     `yaml:"type"`
	Name       string                `yaml:"name"`
	MacAddress string                `yaml:"mac_address"`
	Subnets    []NetworkConfigSubnet `yaml:"subnets"`
}
