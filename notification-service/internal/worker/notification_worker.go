package worker

import (
	"encoding/json"
	"fmt"
	"go-ecommerce/common/pkg/rabbitmq"
	"go-ecommerce/notification-service/internal/service"
	"log"
)

type NotificationWorker struct {
	rabbitMQService rabbitmq.IRabbitMQService
	mailTrapService service.IMailTrapService
}

func NewNotificationWorker(rabbitService rabbitmq.IRabbitMQService, mailTrapService service.IMailTrapService) *NotificationWorker {
	return &NotificationWorker{
		rabbitMQService: rabbitService,
		mailTrapService: mailTrapService,
	}
}

func (w *NotificationWorker) Start() {
	log.Println("[*] Notification Worker is starting...")

	go w.rabbitMQService.Consume(
		"notification_queue",
		"user.created",
		"topic_exchange",
		func(body []byte) { // 1. Mở hàm ẩn danh ở đây
			var msg struct {
				Email string `json:"email"`
				Otp string `json:"otp"`
				}
				if err := json.Unmarshal(body, &msg); err != nil {
				log.Printf("Lỗi giải mã message: %v", err)
				return
   			}

			htmlContent := fmt.Sprintf(`
				<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; border: 1px solid #e0e0e0; border-radius: 8px; overflow: hidden;">
				<div style="background-color: #5b42f3; color: white; padding: 20px; text-align: center; font-size: 24px; font-weight: bold;">
					Mã OTP Của Bạn
				</div>
				
				<div style="padding: 30px;">
					<p style="font-size: 16px; color: #333;">Xin chào,</p>
					<p style="font-size: 16px; color: #333;">Mã xác thực (OTP) cho tài khoản của bạn là:</p>
					
					<div style="background-color: #f4f4f4; padding: 20px; text-align: center; font-size: 32px; font-weight: bold; color: #5b42f3; margin: 20px 0; border-radius: 4px; letter-spacing: 5px;">
						%s
					</div>
					
					<p style="font-size: 14px; color: #666;">Mã OTP này có hiệu lực trong <b>2 phút</b>. Vui lòng không chia sẻ mã này với bất kỳ ai.</p>
					<p style="font-size: 14px; color: #666;">Nếu bạn không yêu cầu mã này, hãy bỏ qua email này.</p>
					<p style="font-size: 14px; color: #666;">Cảm ơn bạn đã sử dụng dịch vụ của chúng tôi!</p>
				</div>
				
				<div style="background-color: #f9f9f9; padding: 15px; text-align: center; font-size: 12px; color: #999;">
					© 2026 Tên Công Ty Của Bạn. Tất cả các quyền được bảo lưu.
				</div>
			</div>`, msg.Otp)
			
        	err := w.mailTrapService.Send(msg.Email, "Xác thực tài khoản", htmlContent)
			if err != nil {
				log.Printf("Lỗi gửi mail: %v", err)
			}
		}, 
	)

}