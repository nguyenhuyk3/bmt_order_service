package message

type PayloadPaymentData struct {
	OrderId int32 `json:"order_id" binding:"required"`
}
