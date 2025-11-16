package entities

import (
	"fmt"
	"time"
)

type SubscriptionStatus string

const (
	// Подписка активна и можно списывать - доступ есть
	SubscriptionStatusActive SubscriptionStatus = "active"
	// Подписка неактивна, списывать нельзя - доступа нет
	SubscriptionStatusDisable SubscriptionStatus = "disable"
	// Подписка отключена, сохранить функционал на оплаченный период, но списывать нельзя - доступ есть до оплаченного периода
	SubscriptionStatusCancel SubscriptionStatus = "cancel"
	// Подписка включена, но не получилось списать деньги с карты - попытка списывать да, но доступа нет
	SubscriptionStatusWaiting SubscriptionStatus = "waiting"
)

type Subscription struct {
	ID                 string
	SubscriptionStatus SubscriptionStatus
	LastPayDate time.Time
	UserID string
	PaymentMethod PaymentMethodSnapshot
	CreatedAt time.Time
	PaymentsCount int64
	AccessUntil time.Time
}

type PaymentMethodSnapshot struct {
	RebillID string // ключ для рекуррентного списания у Т-Банка
	Active   bool   // можно по ней списать сейчас
	Brand    string // "VISA", "MIR", и т.п. — можно отдать во фронт
	Last4    string // "1234" — последние 4 цифры карты
}

func NewSubscription(userID string, ID string, paymentMethod PaymentMethodSnapshot, accessUntil time.Time) (*Subscription, error) {
	if userID == "" {
		return nil, fmt.Errorf("invalid param: %w", ErrInvalidUserId)
	}

	if ID == "" {
		return nil, fmt.Errorf("invalid param: %w", ErrInvalidSubscriptionId)
	}

	if !paymentMethod.Active {
		return nil, fmt.Errorf("invalid param: %w", ErrInvalidPayMethod)
	}

	if paymentMethod.RebillID == "" {
		return nil, fmt.Errorf("invalid param: %w", ErrInvalidRebillId)
	}

	return &Subscription{
		ID:                 ID,
		SubscriptionStatus: SubscriptionStatusActive,
		LastPayDate:        time.Now(),
		UserID:             userID,
		PaymentMethod:      paymentMethod,
		CreatedAt:          time.Now(),
		PaymentsCount:              1,
		AccessUntil:				accessUntil,
	}, nil
}

func (subscription *Subscription) MarkPaid(accessUntil time.Time) {
	subscription.PaymentsCount += 1
	subscription.LastPayDate = time.Now()
	subscription.SubscriptionStatus = SubscriptionStatusActive
	subscription.AccessUntil = accessUntil
}

func (subscription *Subscription) SetStatus(status SubscriptionStatus) {
	subscription.SubscriptionStatus = status
}

func (subscription *Subscription) SetPayInfoSubscription(RebillID string, Active bool, Brand string, Last4 string) error {
	if !Active {
		return fmt.Errorf("pay method not active: %w", ErrInvalidParam)
	}

	if RebillID == "" {
		return fmt.Errorf("rebill id is empty: %w", ErrInvalidParam)
	}

	subscription.PaymentMethod.Active = Active
	subscription.PaymentMethod.Brand = Brand
	subscription.PaymentMethod.Last4 = Last4
	subscription.PaymentMethod.RebillID = RebillID

	return nil
}

// func (subscription *Subscription) SubscriptionIsActive() bool { // На уровне сервиса - это бизнес логика?
// 	if subscription.SubscriptionStatus == SubscriptionStatusWaiting ||
// 		 subscription.SubscriptionStatus == SubscriptionStatusDisable {
// 			return false
// 	}

// 	today := time.Now()

// 	isExpired := subscription.AccessUntil.Before(today)

// 	if isExpired {
// 		return false
// 	}

// 	if subscription.SubscriptionStatus == SubscriptionStatusActive || subscription.SubscriptionStatus == SubscriptionStatusCancel  {
// 		return true
// 	}

// 	return false
// }