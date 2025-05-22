package readers

import (
	"bmt_order_service/db/sqlc"
	"bmt_order_service/global"
	"bmt_order_service/internal/services"
	"context"
	"log"
)

type MessageBrokerReader struct {
	SqlQuery    sqlc.IStore
	RedisClient services.IRedis
	Context     context.Context
}

var topics = []string{
	global.BMT_PAYMENT_PUBLIC_OUTBOXES,
}

func NewMessageBrokerReader(
	sqlQuery sqlc.IStore,
	redisClient services.IRedis,
) *MessageBrokerReader {
	return &MessageBrokerReader{
		SqlQuery:    sqlQuery,
		RedisClient: redisClient,
		Context:     context.Background(),
	}
}

func (m *MessageBrokerReader) InitReaders() {
	log.Printf("=============== Showtime Service is listening to messages ... ===============\n\n\n")

	for _, topic := range topics {
		go m.startReader(topic)
	}
}
