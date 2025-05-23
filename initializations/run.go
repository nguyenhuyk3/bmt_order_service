package initializations

import (
	"bmt_order_service/global"
	"fmt"
)

func Run() {
	loadConfigs()
	initPostgreSql()
	initRedis()
	initMessageBrokerReader()

	r := initRouter()

	r.Run(fmt.Sprintf("0.0.0.0:%s", global.Config.Server.ServerPort))
}
