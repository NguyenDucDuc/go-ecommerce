package module

import (
	"go-ecommerce/notification-service/internal/config"
	"go-ecommerce/notification-service/internal/service"
)

type MailTrapModule struct {
	Service service.IMailTrapService
}

func NewMailTrapModule(cfg *config.NotificationServiceConfig) *MailTrapModule {
	service := service.NewMailTrapService(cfg.MailTrapToken,cfg.MailTrapUrl ,cfg.MailTrapMailSender, cfg.MailTrapNameSender)
	return &MailTrapModule{
		Service: service,
	}
}