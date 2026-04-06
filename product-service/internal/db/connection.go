package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

func ConnectDB(uri string, dbName string) *mongo.Database {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    clientOptions := options.Client().ApplyURI(uri)

    client, err := mongo.Connect(clientOptions) // V2 cho phép truyền options trực tiếp ở đây
    if err != nil {
        log.Fatalf("Lỗi khởi tạo MongoDB Client: %v", err)
    }

    err = client.Ping(ctx, readpref.Primary())
    if err != nil {
        log.Fatalf("Không thể ping tới MongoDB: %v", err)
    }

    fmt.Println("✅ Đã kết nối thành công tới MongoDB (Driver v2)!")

    return client.Database(dbName)
}