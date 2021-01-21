package cloudinit

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"testing"
)

func Test(t *testing.T) {
	tmpCloudInit := CloudInit{
		DirPath:  "",
		MetaData: MetaData{},
		UserData: UserData{},
		NetworkConfig: NetworkCon{
			Version: 0,
			Config: []NetworkConfig{
				{
					Type:       NetworkConfigTypePhysical,
					Name:       "test",
					MacAddress: "00:11:22:33:44:55",
					//Subnets: []NetworkConfigSubnet{
					//	{
					//
					//	},
					//},
				},
			},
		},
	}

	yaml, _ := yaml.Marshal(tmpCloudInit)
	fmt.Println(string(yaml))
}
