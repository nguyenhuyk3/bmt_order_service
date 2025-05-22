package readers

import (
	"bmt_order_service/dto/message"
	"bmt_order_service/dto/request"
	"bmt_order_service/global"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

func (m *MessageBrokerReader) startReader(topic string) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{
			global.Config.ServiceSetting.KafkaSetting.KafkaBroker_1,
			global.Config.ServiceSetting.KafkaSetting.KafkaBroker_2,
			global.Config.ServiceSetting.KafkaSetting.KafkaBroker_3,
		},
		GroupID:        global.PAYMENT_SERVICE_GROUP,
		Topic:          topic,
		CommitInterval: time.Second * 5,
	})
	defer reader.Close()

	for {
		message, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("error reading message: %v\n", err)
			continue
		}

		m.processMessage(topic, message.Value)
	}
}

func (m *MessageBrokerReader) processMessage(topic string, value []byte) {
	switch topic {
	case global.BMT_PAYMENT_PUBLIC_OUTBOXES:
		var messageData message.BMTPublicOutboxesMsg
		if err := json.Unmarshal(value, &messageData); err != nil {
			log.Printf("failed to unmarshal payment message: %v\n", err)
			return
		}

		m.handleCreateSubOrder(messageData)

	default:
		log.Printf("unknown topic received: %s", topic)
	}
}

func (m *MessageBrokerReader) handleCreateSubOrder(messageData message.BMTPublicOutboxesMsg) {
	var payload message.PayloadPaymentData
	if err := json.Unmarshal([]byte(messageData.After.Payload), &payload); err != nil {
		log.Printf("failed to parse payload (%s): %v", messageData.After.EventType, err)
		return
	}

	var orderRedisKey string = fmt.Sprintf("%s%d", global.ORDER, payload.OrderId)
	var subOrder request.SubOrder
	if err := m.RedisClient.Get(orderRedisKey, &subOrder); err != nil {
		log.Printf("failed to get data with key %s", orderRedisKey)
		return
	}

	switch messageData.After.EventType {
	case global.PAYMENT_SUCCESS:
		err := m.SqlQuery.CreateSubOrderTran(m.Context, subOrder, true)
		if err != nil {
			log.Printf("failed to create sub order tran (%s): %v", messageData.After.EventType, err)
			return
		}
		log.Printf("create sub order tran (%s) successfully", messageData.After.EventType)

	case global.PAYMENT_FAILED:
		err := m.SqlQuery.CreateSubOrderTran(m.Context, subOrder, false)
		if err != nil {
			log.Printf("failed to update sub order tran (%s): %v", messageData.After.EventType, err)
			return
		}
		log.Printf("update sub order tran (%s) successfully", messageData.After.EventType)

	default:
		log.Printf("unknown event type received: %s", messageData.After.EventType)
	}
}
