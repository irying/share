package conf

import (
	"backend/utils"
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"os"
)

var (
	// Conf holds the global app config.
	Conf config
)

type config struct {
	// 应用配置
	App app

	// MySQL
	DB database `toml:"database"`

	// 静态资源
	Static static

	// Redis
	Redis redis `toml:"redis"`
}

// 项目配置
type app struct {
	Api struct {
		ApiName        string   `toml:"apiName"`
		ApiVersion     string   `toml:"apiVersion"`
		ApiAddr        string   `toml:"apiAddr"`
		ApiTls         bool     `toml:"apiTls"`
		ApiTlsAddr     string   `toml:"apiTlsAddr"`
		ApiSecretKey   string   `toml:"apiSecretKey"`
		ApiUsername    string   `toml:"apiUsername"`
		ApiPassword    string   `toml:"apiPassword"`
		ApiLogFile     string   `toml:"apiLogFile"`
		ImageSavePath    string   `toml:"imageSavePath"`
		ImageAllowExts []string `toml:"imageAllowExts"`
		ImageMaxSize   int      `toml:"imageMaxSize"`
	}

	Rpc struct {
		RpcName      string `toml:"rpcName"`
		RpcVersion   string `toml:"rpcVersion"`
		RpcAddr      string `toml:"rpcAddr"`
		RpcTraceAddr string `toml:"rpcTraceAddr"`
		RpcLogFile   string `toml:"rpcLogFile"`
	}
}

type static struct {
}

type database map[string]map[string]string

type redis struct {
	Server string `toml:"server"`
	Pwd    string `toml:"pwd"`
}

func InitConfig(env string, app string) {

	if env == "" {
		env = "dev"
	}

	var configFile string
	if app == "api" {
		configFile = "../conf/conf_" + env + ".toml"
	}
	if app == "rpc" {
		configFile = "conf/conf_" + env + ".toml"
	}

	if app == "test" {
		configFile = "../../conf/conf_" + env + ".toml"
	}

	if _, err := os.Stat(configFile); err != nil {
		utils.LogFatalfError(err)
	} else {

		configBytes, err := ioutil.ReadFile(configFile)
		if err != nil {
			utils.LogFatalfError(err)
		}

		_, err = toml.Decode(string(configBytes), &Conf)
		if err != nil {
			utils.LogFatalfError(err)
		}

		utils.LogPrint("Load config from file : %s", configFile)
	}
}
