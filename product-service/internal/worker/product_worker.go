package worker

import (
	"go-ecommerce/common/pkg/rabbitmq"
	"go-ecommerce/product-service/internal/service"
	"log"
)

type ProductWorker struct {
	rabbitService rabbitmq.IRabbitMQService
	productService *service.ProductService
}

func NewProductWorker(rabbit rabbitmq.IRabbitMQService, productSvc *service.ProductService) *ProductWorker {
	return &ProductWorker{
		rabbitService:  rabbit,
		productService: productSvc,
	}
}

func (w *ProductWorker) Start() {
	log.Println("[*] Product Worker is starting...")

	// Worker lắng nghe sự kiện tạo đơn hàng để kiểm tra kho
	go w.rabbitService.Consume(
		"product_inventory_queue",
		"order.created",
		"topic_exchange",
		w.productService.OrderCreated, // Đây là hàm xử lý logic trừ kho
	)
}