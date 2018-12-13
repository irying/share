package serviceTests

import "backend/conf"

// dev/pro
const Env = "dev"

func init()  {
	// init config
	conf.InitConfig(Env,"test")
}
