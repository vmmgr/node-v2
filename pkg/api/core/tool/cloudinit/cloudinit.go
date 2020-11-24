package cloudinit

import (
	"fmt"
	"github.com/vmmgr/node/pkg/api/core/tool/file"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

type CloudInit struct {
	DirPath       string
	MetaData      MetaData   `json:"meta"`
	UserData      UserData   `json:"user"`
	NetworkConfig NetworkCon `json:"network"`
}

func NewCloudInitHandler(input CloudInit) *CloudInit {
	return &CloudInit{
		DirPath:       input.DirPath,
		MetaData:      input.MetaData,
		UserData:      input.UserData,
		NetworkConfig: input.NetworkConfig,
	}
}

// humstackのcloud-init部分を引用
// https://github.com/ophum/humstack/tree/master/pkg/utils/cloudinit

func (c *CloudInit) Generate() error {
	metaDataYAML, err := yaml.Marshal(c.MetaData)
	if err != nil {
		return err
	}

	userDataYAML, err := yaml.Marshal(c.UserData)
	if err != nil {
		return err
	}

	networkConfigYAML, err := yaml.Marshal(c.NetworkConfig)
	if err != nil {
		return err
	}

	metaDataPath := filepath.Join(c.DirPath, "meta-data")
	err = ioutil.WriteFile(metaDataPath, metaDataYAML, 0666)
	if err != nil {
		return err
	}

	userDataPath := filepath.Join(c.DirPath, "user-data")
	userDataYAML = []byte(fmt.Sprintf("#cloud-config\n%s", userDataYAML))
	err = ioutil.WriteFile(userDataPath, userDataYAML, 0666)
	if err != nil {
		return err
	}

	networkConfigPath := filepath.Join(c.DirPath, "network-config")
	err = ioutil.WriteFile(networkConfigPath, networkConfigYAML, 0666)
	if err != nil {
		return err
	}

	command := "cloud-localds"
	args := []string{
		"-N",
		networkConfigPath,
		filepath.Join(c.DirPath, "cloudinit.img"),
		userDataPath,
		metaDataPath,
	}

	//ファイルがすでに存在する場合は削除する
	if file.ExistsCheck(filepath.Join(c.DirPath, "cloudinit.img")) {
		if err := os.Remove(filepath.Join(c.DirPath, "cloudinit.img")); err != nil {
			log.Println(err)
			return err
		}
	}
	cmd := exec.Command(command, args...)
	if _, err := cmd.CombinedOutput(); err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}
