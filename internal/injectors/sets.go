package injectors

import (
	"bmt_order_service/db/sqlc"
	"bmt_order_service/internal/implementations/redis"
	"bmt_order_service/internal/injectors/provider"

	"github.com/google/wire"
)

var dbSet = wire.NewSet(
	provider.ProvidePgxPool,
	sqlc.NewStore,
)

var redisSet = wire.NewSet(
	redis.NewRedisClient,
)

var productClientSet = wire.NewSet(
	provider.ProvideProductClient,
)
