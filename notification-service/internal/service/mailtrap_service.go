package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type MailTrapService struct {
	MailTrapToken string
	MailTrapUri string
	FromEmail      string
	FromName       string
	Client         *http.Client
}

type Recipient struct {
	Email string `json:"email"`
}

type SandboxPayload struct {
	From    map[string]string `json:"from"`
	To      []Recipient       `json:"to"`
	Subject string            `json:"subject"`
	HTML    string            `json:"html"`
}

func NewMailTrapService(token, url ,fromEmail, fromName string) *MailTrapService {
	return &MailTrapService{
		MailTrapToken: token,
		FromEmail: fromEmail,
		FromName: fromName,
		MailTrapUri: url,
		Client: &http.Client{Timeout: 15 * time.Second},
	}
}

func (s *MailTrapService) Send(to, subject, htmlContent string) error {
	// Kiểm tra dữ liệu đầu vào
	if to == "" {
		return fmt.Errorf("địa chỉ email người nhận (to) bị trống")
	}

	// Tạo payload chuẩn xác bằng Struct
	payload := SandboxPayload{
		From: map[string]string{
			"email": "hello@demomailtrap.com",
			"name":  "Mailtrap Test",
		},
		To: []Recipient{
			{Email: to},
		},
		Subject: subject,
		HTML:    htmlContent,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("lỗi tạo JSON: %v", err)
	}

	// Debug: In ra JSON để kiểm tra cấu trúc
	log.Printf("DEBUG: Đang gửi JSON tới Mailtrap: %s", string(jsonData))

	req, err := http.NewRequest("POST", s.MailTrapUri, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+s.MailTrapToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("mailtrap sandbox error: %d - %s", resp.StatusCode, string(body))
	}

	log.Printf("[✔] Mail đã được gửi thành công tới: %s", to)
	return nil
}