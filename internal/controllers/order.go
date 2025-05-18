package controllers

import (
	"bmt_order_service/dto/request"
	"bmt_order_service/global"
	"bmt_order_service/internal/responses"
	"bmt_order_service/internal/services"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type OrderController struct {
	OrderService services.IOrder
}

func NewOrderController(orderService services.IOrder) *OrderController {
	return &OrderController{
		OrderService: orderService,
	}
}

func (o *OrderController) CreateOrder(c *gin.Context) {
	var req request.Order
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.FailureResponse(c, http.StatusBadRequest, fmt.Sprintf("invalid request: %v", err))
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	orderedBy := c.GetString(global.X_USER_EMAIL)
	req.OrderedBy = orderedBy

	status, err := o.OrderService.CreateOrder(ctx, req)
	if err != nil {
		responses.FailureResponse(c, status, err.Error())
		return
	}

	responses.SuccessResponse(c, status, "add new order perform successfully", nil)
}
