package repository

import (
	"context"
	"fmt"
	util "go-ecommerce/common/utils"
	"go-ecommerce/product-service/internal/model"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type InventoryRepository struct {
	collection *mongo.Collection
}

func NewInventoryRepository(db *mongo.Database) IInventoryRepository {
	return &InventoryRepository{
		collection: db.Collection("inventories"),
	}
}

func (inventoryRepo *InventoryRepository) Create(ctx context.Context, inventory *model.Inventory) (*model.Inventory, error) {
	res, err := inventoryRepo.collection.InsertOne(ctx, inventory)
	if err != nil {
		return &model.Inventory{}, util.NewAppError(http.StatusInternalServerError, util.ErrInternalServer, "Failed insert into database")
	}

	inventory.ID = res.InsertedID.(bson.ObjectID)
	return inventory, nil
}

func (inventoryRepo *InventoryRepository) UpdateOne(ctx context.Context, id string, updateData bson.M) error {
    pID, err := bson.ObjectIDFromHex(id)
    if err != nil {
        return fmt.Errorf("invalid object id: %w", err)
    }

    filter := bson.M{"_id": pID}
    
    // 1. Luôn cập nhật thời gian vào bản ghi
    updateData["updated_at"] = time.Now()

    // 2. Bọc toàn bộ data vào $set để MongoDB hiểu đây là lệnh cập nhật field
    // Cách này giúp bạn truyền bson.M{"available_stock": 10} từ Service cực kỳ gọn
    finalUpdate := bson.M{
        "$set": updateData,
    }

    // 3. Thực thi với context (có thể là sessCtx từ transaction)
    _, err = inventoryRepo.collection.UpdateOne(ctx, filter, finalUpdate)
    if err != nil {
        return fmt.Errorf("failed to update inventory: %w", err)
    }

    return nil
}

func (r *InventoryRepository) UpdateStock(ctx context.Context, productId string, quantity int) error {
	pID, err := bson.ObjectIDFromHex(productId)
	if err != nil {
		return fmt.Errorf("invalid product id format: %w", err)
	}

	// Filter tìm theo product_id (vì bảng inventory lưu product_id làm ngoại khóa)
	// Điều kiện phụ: available_stock phải đủ để trừ
	filter := bson.M{
		"product_id":      pID,
		"available_stock": bson.M{"$gte": quantity}, 
	}

	// Update logic:
	// 1. Giảm available_stock đi một lượng 'quantity'
	// 2. Tăng reserved_stock lên một lượng 'quantity' (để giữ chỗ)
	update := bson.M{
		"$inc": bson.M{
			"available_stock": -quantity,
			"reserved_stock":  quantity,
		},
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("lỗi database khi cập nhật inventory: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("sản phẩm %s không tồn tại trong kho hoặc không đủ số lượng", productId)
	}

	return nil
}