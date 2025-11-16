package services

import "fmt"

type Service struct {}

func NewService() *Service {return &Service{}}

func (service *Service) InitPay(userId string) (string, error) {

	return fmt.Sprintf("User pay, %s", userId), nil
}