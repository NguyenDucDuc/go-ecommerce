package config

import util "go-ecommerce/common/utils"

type NotificationServiceConfig struct {
	MailTrapToken    string
	MailTrapUrl        string
	MailTrapMailSender string
	MailTrapNameSender string
	RabbitMQUri string
}

func NewNotificationServiceConfig() *NotificationServiceConfig {
	return &NotificationServiceConfig{
		MailTrapToken: util.GetEnv("MAILTRAP_API_KEY", "123"),
		MailTrapUrl: util.GetEnv("MAILTRAP_URL", "123"),
		MailTrapMailSender: util.GetEnv("MAILTRAP_MAIL_SENDER", "nguyenducduc26797@gmail.com"),
		MailTrapNameSender: util.GetEnv("MAILTRAP_NAME_SENDER", "go-ecommerce"),
		RabbitMQUri: util.GetEnv("RABBIT_MQ_URI", "amqp://root:admin123@localhost:5672/"),
	}
}