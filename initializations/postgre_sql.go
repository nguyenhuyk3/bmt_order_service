package initializations

import (
	"bmt_order_service/global"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func initPostgreSql() {
	config := global.Config.ServiceSetting.PostgreSql
	dbName := global.Config.ServiceSetting.PostgreSql.DbName
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.Username, config.Password, dbName)
	ctx := context.Background()
	db, err := pgxpool.New(ctx, connStr)
	if err != nil {
		log.Println("error connecting to the database:", err)
		return
	}

	if err := db.Ping(ctx); err != nil {
		log.Println("error pinging the database:", err)

		os.Exit(1)
	}

	log.Println("=============== successfully connected to the database ===============")

	global.Postgresql = db
}
