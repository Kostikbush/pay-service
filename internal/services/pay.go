package pay

import "fmt"

type Service struct {}

func NewService() *Service {return &Service{}}

func (service *Service) InitPay(userId string) string {
	return fmt.Sprintf("User pay, %s", userId)
}