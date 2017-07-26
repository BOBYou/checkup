package g

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/toolkits/file"
)



type CheckupConfig struct {
	IpRange []string `json:"ipRange"`
	PingTimeout int `json:"pingTimeout"`
	PingRetry   int `json:"pingRetry"`
	FastPingMode	bool	`json:"fastPingMode"`
	Interval int    `json:"interval"`
	PostUrl string    `json:"postUrl"`
	HostName string    `json:"hostName"`
	FailureRate float64	`json:"failureRate"`
	To	string `json:"to"`
	FailsInterval int `json:"failsInterval"`
}




type GlobalConfig struct {
	Checkup      *CheckupConfig      `json:"checkup"`
}

var (
	ConfigFile string
	config     *GlobalConfig
	lock       = new(sync.RWMutex)
)

func Config() *GlobalConfig {
	lock.RLock()
	defer lock.RUnlock()
	return config
}


func ParseConfig(cfg string) {
	if cfg == "" {
		fmt.Println("use -c to specify configuration file")
	}

	if !file.IsExist(cfg) {
		fmt.Println("config file:", cfg, "is not existent. maybe you need `mv cfg.example.json cfg.json`")
	}

	ConfigFile = cfg

	configContent, err := file.ToTrimString(cfg)
	if err != nil {
		fmt.Println("read config file:", cfg, "fail:", err)
	}

	var c GlobalConfig
	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		fmt.Println("parse config file:", cfg, "fail:", err)
	}

	lock.Lock()
	defer lock.Unlock()

	config = &c

	fmt.Println("read config file:", cfg, "successfully")

}
