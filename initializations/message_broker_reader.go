package initializations

import (
	"bmt_order_service/internal/injectors"
	"log"
)

func initMessageBrokerReader() {
	reader, err := injectors.InitMessageBroker()
	if err != nil {
		log.Fatalf("an error occur when initiallizating ORDER READERS: %v", err)
	}

	reader.InitReaders()
}
