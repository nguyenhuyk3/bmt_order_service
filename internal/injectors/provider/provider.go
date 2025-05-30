package provider

import (
	"bmt_order_service/db/sqlc"
	"bmt_order_service/global"
	"log"
	"sync"

	"product"

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
	productClient     product.ProductClient
	productClientOnce sync.Once
)

func ProvideProductClient() product.ProductClient {
	productClientOnce.Do(func() {
		conn, err := grpc.Dial("localhost:50033", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("cannot connect to product service on port 50033: %v", err)
		}
		productClient = product.NewProductClient(conn)
	})

	return productClient
}
