package service

import (
	"context"
	order "go-ecommerce/common/gen-proto/orders"
	product "go-ecommerce/common/gen-proto/products"
	pkg_redis "go-ecommerce/common/pkg/redis"
	util "go-ecommerce/common/utils"
	"go-ecommerce/order-service/internal/model"
	"go-ecommerce/order-service/internal/repository"
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
	order.UnimplementedOrderServiceServer
}

func NewOrderService(repo repository.IOrderRepository, redis pkg_redis.IRedisService, productClient product.ProductServiceClient) *OrderService {
	return &OrderService{
		repo: repo,
		redis: redis,
		productClient: productClient,
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
		util.PrettyPrint(productInfo)
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
		Status: order.OrderStatus_PENDING,
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

	return orderRsp, nil
}