package main

import (
	"go-ecommerce/common/pkg/rabbitmq"
	util "go-ecommerce/common/utils"
	"go-ecommerce/notification-service/internal/config"
	"go-ecommerce/notification-service/internal/module"
	"go-ecommerce/notification-service/internal/worker"
	"log"
)

func main() {
	util.LoadEnv()
	cfg := config.NewNotificationServiceConfig()

	// rabbit mq
	rabbitMQService, err := rabbitmq.NewRabbitMQ(cfg.RabbitMQUri)
	if err != nil {
		log.Fatal(err)
	}

	// load module
	mailTrapModule := module.NewMailTrapModule(cfg)

	// notification worker
	mailTrapWorker := worker.NewNotificationWorker(rabbitMQService, mailTrapModule.Service)
	mailTrapWorker.Start()

	select {}
}