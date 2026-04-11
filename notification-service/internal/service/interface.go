package service

type IMailTrapService interface {
	Send(to, subject, htmlContent string) error
}