package service

import (
	"context"
	"encoding/json"
	order "go-ecommerce/common/gen-proto/orders"
	product "go-ecommerce/common/gen-proto/products"
	"go-ecommerce/common/pkg/rabbitmq"
	pkg_redis "go-ecommerce/common/pkg/redis"
	util "go-ecommerce/common/utils"
	"go-ecommerce/order-service/internal/model"
	"go-ecommerce/order-service/internal/repository"
	"log"
	"net/http"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type OrderService struct {
	repo repository.IOrderRepository
	redis pkg_redis.IRedisService
	productClient product.ProductServiceClient
	rabbitMQService rabbitmq.IRabbitMQService
	order.UnimplementedOrderServiceServer
}

func NewOrderService(repo repository.IOrderRepository, redis pkg_redis.IRedisService, productClient product.ProductServiceClient, rabbitMQService rabbitmq.IRabbitMQService) *OrderService {
	return &OrderService{
		repo: repo,
		redis: redis,
		productClient: productClient,
		rabbitMQService: rabbitMQService,
	}
}

func (orderService *OrderService) CreateOrder(ctx context.Context, input *order.CreateOrderDto) (*order.Order, error) {
	userId, err := bson.ObjectIDFromHex(input.UserId)
	if err != nil {
		return &order.Order{}, util.NewAppError(http.StatusInternalServerError, util.ErrInternalServer, "Invalid user_id")
	}
	
	var orderItems []model.OrderItem
	var total float64
	for _, item := range input.Items {
		pId, _ := bson.ObjectIDFromHex(item.ProductId)
		findProdDto := &product.FindByIdDto{
			ProductId: item.ProductId,
		}
		productInfo, _ := orderService.productClient.FindById(ctx, findProdDto)
		orderItem := model.OrderItem{
			ProductID: pId,
			ProductName: productInfo.Name,
			Price: util.ToDecimal128(productInfo.Price),
			Quantity: int(item.Quantity),
		}

		orderItems = append(orderItems, orderItem)
		priceFloat, err := strconv.ParseFloat(productInfo.Price, 64)
		if err != nil {
			priceFloat = 0
		}
		total += priceFloat * float64(item.Quantity)
	}

	userId, _ = bson.ObjectIDFromHex(input.UserId)
	orderModel := &model.Order{
		OrderCode: "ORD" + time.Now().Format("20060102150405"),
		UserID: userId,
		Items: orderItems,
		TotalAmount: util.Float64ToDecimal128(total),
		Status: "PENDING",
		ShippingAddress: input.ShippingAddress,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	res, err := orderService.repo.Create(ctx, orderModel)
	if err != nil {
		return &order.Order{}, err
	}

	orderItemRsp := make([]*order.OrderItem, len(res.Items))
	for i, item := range res.Items {
		itemRsp := &order.OrderItem{
			ProductId: item.ProductID.Hex(),
			ProductName: item.ProductName,
			Price: util.DecimalToString(item.Price),
		}
		orderItemRsp[i] = itemRsp
	}

	orderRsp := &order.Order{
		Id: res.ID.Hex(),
		OrderCode: res.OrderCode,
		UserId: res.UserID.Hex(),
		TotalAmount: res.TotalAmount.String(),
		Status: res.Status,
		Items: orderItemRsp,
		ShippingAddress: res.ShippingAddress,
		CreatedAt: timestamppb.New(res.CreatedAt),
		UpdatedAt: timestamppb.New(res.UpdatedAt),
	}

	// message queue to product
	msg := map[string]interface{}{
		"order_id": res.ID.Hex(),
		"items": input.Items,
	}

	err = orderService.rabbitMQService.Publish("topic_exchange", "order.created", msg)
	if err != nil {
		log.Printf("Failed to publish message: %v", err)
	}

	return orderRsp, nil
}

// HandleInventorySuccess xử lý khi kho đã được giữ chỗ thành công
func (s *OrderService) HandleInventorySuccess(body []byte) {
    var msg map[string]interface{}
    json.Unmarshal(body, &msg)
    
    orderID := msg["order_id"].(string)
    log.Printf("Order %s: Kho đã sẵn sàng. Cập nhật status -> CONFIRMED", orderID)

    // Update status trong MongoDB thành CONFIRMED (hoặc DONE)
	updateData := bson.M{
		"status": "CONFIRMED",
	}
	oId, _ := bson.ObjectIDFromHex(orderID)
    s.repo.UpdateOne(context.Background(), oId, updateData)
}

// HandleInventoryFailed xử lý khi kho không đủ
func (s *OrderService) HandleInventoryFailed(body []byte) {
    var msg map[string]interface{}
    json.Unmarshal(body, &msg)
    
    orderID := msg["order_id"].(string)
    reason := msg["reason"].(string)
    log.Printf("Order %s thất bại do: %s. Cập nhật status -> FAILED", orderID, reason)

    // Update status trong MongoDB thành FAILED
	updateData := bson.M{
		"status": "FAILED",
	}

	oId, _ := bson.ObjectIDFromHex(orderID)
    s.repo.UpdateOne(context.Background(), oId, updateData)
}