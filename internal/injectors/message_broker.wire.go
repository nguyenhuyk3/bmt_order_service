//go:build wireinject

package injectors

import (
	"bmt_order_service/internal/implementations/message_broker/readers"

	"github.com/google/wire"
)

func InitMessageBroker() (*readers.MessageBrokerReader, error) {
	wire.Build(
		dbSet,
		redisSet,

		readers.NewMessageBrokerReader,
	)

	return nil, nil
}
