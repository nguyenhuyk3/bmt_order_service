package provider

import (
	"bmt_order_service/db/sqlc"
	"bmt_order_service/global"
	"fmt"
	"log"
	"sync"

	"product"
	"showtime"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func ProvidePgxPool() *pgxpool.Pool {
	return global.Postgresql
}

func ProvideQueries() *sqlc.Queries {
	return sqlc.New(global.Postgresql)
}

var (
	productClient  product.ProductClient
	showtimeClient showtime.ShowtimeClient

	productClientOnce  sync.Once
	showtimeClientOnce sync.Once
)

func ProvideProductClient() product.ProductClient {
	productClientOnce.Do(func() {
		conn, err := grpc.Dial(
			fmt.Sprintf("localhost:%s", global.Config.Server.ProductRPCServerPort),
			grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("cannot connect to product service on port %s: %v", global.Config.Server.ProductRPCServerPort, err)
		}

		productClient = product.NewProductClient(conn)
	})

	return productClient
}

func ProvideShowtimeClient() showtime.ShowtimeClient {
	showtimeClientOnce.Do(func() {
		conn, err := grpc.Dial(
			fmt.Sprintf("localhost:%s", global.Config.Server.ShowtimeRPCServerPort),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			log.Fatalf("cannot connect to showtime service on port %s: %v", global.Config.Server.ShowtimeRPCServerPort, err)
		}

		showtimeClient = showtime.NewShowtimeClient(conn)
	})

	return showtimeClient
}
