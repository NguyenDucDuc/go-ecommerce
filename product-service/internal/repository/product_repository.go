package repository

import (
	"context"
	util "go-ecommerce/common/utils"
	"go-ecommerce/product-service/internal/model"
	"net/http"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ProductRepository struct {
	collection *mongo.Collection
}

func NewProductRepository(db *mongo.Database) IProductRepository {
	return &ProductRepository{
		collection: db.Collection("products"),
	}
}

func (productRepo *ProductRepository) Create(ctx context.Context, product *model.Product) (*model.Product, error) {
	result, err := productRepo.collection.InsertOne(ctx, product)
	if err != nil {
		return &model.Product{}, util.NewAppError(http.StatusInternalServerError, util.ErrInternalServer, "Failed insert to database")
	}

	product.ID = result.InsertedID.(bson.ObjectID)

	return product, nil
}

func (productRepo *ProductRepository) FindAll(ctx context.Context, skip, limit int, orderBy, sort string)([]*model.ProductWithInventory, int, error) {
	// 1. Xác định chiều sắp xếp
	sortOrder := -1
	if sort == "" || sort == "desc" {
		sortOrder = -1
	}
	if sort == "asc" {
		sortOrder = 1
	}
	if orderBy == "" {
		orderBy = "created_at"
	}

	// 2. Xây dựng Pipeline với Facet
	pipeline := mongo.Pipeline{
    {{Key: "$facet", Value: bson.D{
        {Key: "metadata", Value: mongo.Pipeline{
            {{Key: "$count", Value: "total"}},
        }},
        {Key: "data", Value: mongo.Pipeline{
				{{Key: "$sort", Value: bson.D{{Key: orderBy, Value: sortOrder}}}},
				{{Key: "$skip", Value: skip}},
				{{Key: "$limit", Value: limit}},
				// BẮT ĐẦU LOOKUP TẠI ĐÂY
				{{Key: "$lookup", Value: bson.D{
					{Key: "from", Value: "inventories"},         // Tên collection inventory
					{Key: "localField", Value: "_id"},         // Trường _id của product
					{Key: "foreignField", Value: "product_id"},// Trường product_id của inventory
					{Key: "as", Value: "inventory_info"},      // Tên field kết quả (mảng)
				}}},
				// Nếu bạn muốn inventory là 1 object thay vì mảng (quan hệ 1-1)
				{{Key: "$unwind", Value: bson.D{
					{Key: "path", Value: "$inventory_info"},
					{Key: "preserveNullAndEmptyArrays", Value: true}, // Giữ lại product nếu chưa có inventory
				}}},
			}},
		}}},
	}

	cursor, err := productRepo.collection.Aggregate(ctx, pipeline)

	if err != nil {
		return nil, 0, util.NewAppError(http.StatusInternalServerError, util.ErrInternalServer, "Failed aggregate database")
	}
	defer cursor.Close(ctx)

	// 3. Định nghĩa cấu trúc để nhận kết quả từ Facet
	var results []struct {
		Metadata []struct {
			Total int `bson:"total"`
		} `bson:"metadata"`
		Data []*model.ProductWithInventory `bson:"data"`
	}



	if err := cursor.All(ctx, &results); err != nil {
		return nil, 0, util.NewAppError(http.StatusInternalServerError, util.ErrInternalServer, "Failed aggregate database")
	}

	// 4. Xử lý kết quả trả về
	if len(results) == 0 || len(results[0].Metadata) == 0 {
		return []*model.ProductWithInventory{}, 0, nil
	}


	return results[0].Data, results[0].Metadata[0].Total, nil
}

func (productRepo *ProductRepository) FindById(ctx context.Context, productId bson.ObjectID)(model.Product) {
	var product model.Product

	filter := bson.D{{Key: "_id", Value: productId}}
	err := productRepo.collection.FindOne(ctx, filter).Decode(&product)
	if err != nil {
		return model.Product{}
	} 

	return product
}


