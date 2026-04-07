package util

import (
	"encoding/json"

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


func MapToProtoStruct(input interface{}) *structpb.Struct {
	if input == nil {
		return nil
	}

	// 1. Dùng bson.MarshalExtJSON để convert sang JSON bytes
	// Hàm này sẽ xử lý được cả bson.D, bson.M và các kiểu primitive (ObjectID, Decimal128)
	data, err := bson.MarshalExtJSON(input, false, false)
	if err != nil {
		return nil
	}

	// 2. Unmarshal vào map[string]interface{} chuẩn của Go
	var cleanMap map[string]interface{}
	if err := json.Unmarshal(data, &cleanMap); err != nil {
		return nil
	}

	// 3. Tạo structpb. Struct này giờ đã nhận "cleanMap" là kiểu chuẩn
	s, err := structpb.NewStruct(cleanMap)
	if err != nil {
		return nil
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