package worker

import (
	"go-ecommerce/common/pkg/rabbitmq"
	"go-ecommerce/order-service/internal/service"
	"log"
)

type Worker struct {
	rabbitService rabbitmq.IRabbitMQService
	orderService  *service.OrderService // Hoặc ProductService tùy service
}

func NewWorker(rabbitService rabbitmq.IRabbitMQService, orderSvc *service.OrderService) *Worker {
	return &Worker{
		rabbitService: rabbitService,
		orderService:  orderSvc,
	}
}

// Start sẽ là nơi đăng ký tất cả các consumer
func (w *Worker) Start() {
	log.Println("Starting RabbitMQ Workers...")

	// Worker xử lý khi inventory thành công
	go w.rabbitService.Consume(
		"order_inventory_success_queue",
		"inventory.success",
		"topic_exchange",
		w.orderService.HandleInventorySuccess,
	)

	// Worker xử lý khi inventory thất bại
	go w.rabbitService.Consume(
		"order_inventory_failed_queue",
		"inventory.failed",
		"topic_exchange",
		w.orderService.HandleInventoryFailed,
	)
}