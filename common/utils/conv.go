package util

import (
	"go.mongodb.org/mongo-driver/v2/bson"
	"google.golang.org/protobuf/types/known/structpb"
)

// ToDecimal128 convert string sang bson.Decimal128 (Dùng khi nhận data từ API/GRPC)
func ToDecimal128(s string) bson.Decimal128 {
	d, err := bson.ParseDecimal128(s)
	if err != nil {
		// Trả về 0 nếu lỗi, hoặc bạn có thể xử lý panic tùy logic
		d, _ = bson.ParseDecimal128("0")
	}
	return d
}

// DecimalToString convert ra string để hiển thị chính xác tuyệt đối
func DecimalToString(d bson.Decimal128) string {
	return d.String()
}


func MapToProtoStruct(m bson.M) *structpb.Struct {
	s, err := structpb.NewStruct(m)
	if err != nil {
		// Nếu lỗi, trả về một struct rỗng thay vì nil để tránh panic
		s, _ = structpb.NewStruct(map[string]interface{}{})
	}
	return s
}

// ProtoStructToMap: Chuyển ngược lại từ gRPC sang bson.M để lưu vào DB
func ProtoStructToMap(s *structpb.Struct) bson.M {
	if s == nil {
		return bson.M{}
	}
	return bson.M(s.AsMap())
}