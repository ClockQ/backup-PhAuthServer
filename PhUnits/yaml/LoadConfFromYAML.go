package yaml

import (
	"io/ioutil"
	"fmt"
	"gopkg.in/yaml.v2"
	"github.com/alfredyang1986/BmServiceDef/BmPodsDefine"
)

func LoadConfFromYAML(path string) (conf *BmPodsDefine.Conf) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
		fmt.Println("error")
	}

	conf = &BmPodsDefine.Conf{}
	err = yaml.Unmarshal(data, conf)
	if err != nil {
		panic(err)
	}
	return
}
