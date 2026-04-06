package util

import (
	"encoding/json"
	"fmt"
)

// PrettyPrint nhận vào interface{} để log mọi loại dữ liệu
func PrettyPrint(v interface{}) {
	// MarshalIndent giúp tạo chuỗi JSON có thụt lề (2 dấu cách)
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Printf("❌ PrettyPrint Error: %v\n", err)
		return
	}
	
	fmt.Println("--- DEBUG OBJECT ---")
	fmt.Println(string(b))
	fmt.Println("--------------------")
}