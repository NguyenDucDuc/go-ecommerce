package util

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/status"
)

const (
	Success               = "SUCCESS"
	ErrBadRequest         = "BAD_REQUEST"           // Dữ liệu gửi lên không hợp lệ
	ErrUnauthorized       = "UNAUTHORIZED"          // Chưa đăng nhập hoặc token hết hạn
	ErrForbidden          = "FORBIDDEN"             // Đã đăng nhập nhưng không có quyền truy cập
	ErrNotFound           = "NOT_FOUND"             // Không tìm thấy tài nguyên (User, Product...)
	ErrConflict           = "CONFLICT"              // Dữ liệu đã tồn tại (trùng Email, SKU...)
	ErrTooManyRequests    = "TOO_MANY_REQUESTS"     // Bị giới hạn rate limit
	ErrInternalServer     = "INTERNAL_SERVER_ERROR" // Lỗi code, lỗi DB không xác định
	ErrServiceUnavailable = "SERVICE_UNAVAILABLE"   // Service đang bảo trì hoặc quá tải
	ErrInvalidCredentials = "INVALID_CREDENTIALS"   // Sai tài khoản hoặc mật khẩu
	ErrTokenExpired       = "TOKEN_EXPIRED"         // Token hết hạn
	ErrUserDisabled       = "USER_DISABLED"         // Tài khoản bị khóa
	ErrOtpInvalid         = "OTP_INVALID"           // Mã OTP sai
)

type AppError struct {
	StatusCode int    `json:"status_code"`
	Code       string `json:"code"`
	Message    string `json:"message"`
}

type ResponseData struct {
	StatusCode int    `json:"status_code"`
	Code       string `json:"code"`
	Message    string `json:"message"`
	Data       any    `json:"data"`
}

type ResponseError struct {
	StatusCode int    `json:"status_code"`
	Code       string `json:"code"`
	Message    string `json:"message"`
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(status int, code, message string) *AppError {
	return &AppError{
		Code:       code,
		StatusCode: status,
		Message:    message,
	}
}

func NewResponseData(ctx *gin.Context, statusCode int, code string, message string, data any) {
	rsp := &ResponseData{
		StatusCode: statusCode,
		Code:       code,
		Message:    message,
		Data:       data,
	}
	ctx.JSON(statusCode, rsp)
}

func NewBindingError(ctx *gin.Context, err error) {
	rspErr := &ResponseError{
		StatusCode: http.StatusBadRequest,
		Code: ErrBadRequest,
		Message: err.Error(),
	}
	ctx.JSON(http.StatusBadRequest, rspErr)
}

func NewResponseError(ctx *gin.Context, err error) {
    if appErr, ok := err.(*AppError); ok {
        rspErr := &ResponseError{
            StatusCode: appErr.StatusCode,
            Code:       appErr.Code,
            Message:    appErr.Message,
        }
        ctx.JSON(appErr.StatusCode, rspErr)
        return
    }

    // Mặc định ban đầu
    statusCode := http.StatusInternalServerError
    message := err.Error()

    // Kiểm tra xem lỗi có phải là từ gRPC (RPC error) không
    if st, ok := status.FromError(err); ok {
        // st.Message() sẽ chỉ lấy phần nội dung sau "desc ="
        // Ví dụ: "Insert to database failed"
        message = st.Message()
        
        // Nếu bạn muốn map mã lỗi gRPC sang HTTP Status Code tương ứng (tùy chọn)
        statusCode = runtime.HTTPStatusFromCode(st.Code()) 
    }

    rspErr := &ResponseError{
        StatusCode: statusCode,
        Code:       ErrInternalServer,
        Message:    message,
    }
    ctx.JSON(statusCode, rspErr)
}
