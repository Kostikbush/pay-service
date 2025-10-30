package subscription

import (
	"errors"
	"time"

	utils "pay-service/internal/utils"
)

// Статус подписки
type SubscriptionStatus string

var (
	ErrUserIdNotValid         = errors.New("userId: user id not valid")
	ErrSubscriptionIdNotValid = errors.New("id: subscription id not valid")
	ErrPayMethodMustBeActive  = errors.New("active: payment method must be active")
	ErrRebillIDNotValid       = errors.New("rebillID: rebillID not valid")
)

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
	// Последняя дата оплаты
	PayDate time.Time
	// id пользователя из MongoDB
	UserID string
	// информация об оплате
	PaymentMethod PaymentMethodSnapshot
	// Дата создания подписки
	CreatedAt time.Time
	// количество списаний по подписке (просто информационное поле)
	Count int64
	// Доступ до
	AccessUntil time.Time
}

type PaymentMethodSnapshot struct {
	RebillID string // ключ для рекуррентного списания у Т-Банка
	Active   bool   // можно по ней списать сейчас
	Brand    string // "VISA", "MIR", и т.п. — можно отдать во фронт
	Last4    string // "1234" — последние 4 цифры карты
}

func CreateActiveSubscription(userID string, ID string, paymentMethod PaymentMethodSnapshot) (*Subscription, error) {
	if userID == "" {
		return nil, ErrUserIdNotValid
	}

	if ID == "" {
		return nil, ErrSubscriptionIdNotValid
	}

	if !paymentMethod.Active {
		return nil, ErrPayMethodMustBeActive
	}

	if paymentMethod.RebillID == "" {
		return nil, ErrRebillIDNotValid
	}

	return &Subscription{
		ID:                 ID,
		SubscriptionStatus: SubscriptionStatusActive,
		PayDate:            time.Now(),
		UserID:             userID,
		PaymentMethod:      paymentMethod,
		CreatedAt:          time.Now(),
		Count:              1,
		AccessUntil:				time.Now().AddDate(0, 1, 0),
	}, nil
}

func (subscription *Subscription) MarkPaid() {
	subscription.Count += 1
	subscription.PayDate = time.Now()
	subscription.SubscriptionStatus = SubscriptionStatusActive
	subscription.AccessUntil = time.Now().AddDate(0, 1, 0)
}

func (subscription *Subscription) CancelSubscription() {
	today := time.Now()

	if utils.AtLeastOneMonthPassed(subscription.PayDate, today) {
		subscription.SubscriptionStatus = SubscriptionStatusDisable
	}else {
		subscription.SubscriptionStatus = SubscriptionStatusCancel
	}
}

func (subscription *Subscription) ExpireIfPeriodEnded() {
	subscription.SubscriptionStatus = SubscriptionStatusDisable
}

func (subscription *Subscription) FailedPaySubscription() {
	subscription.SubscriptionStatus = SubscriptionStatusWaiting
}

func (subscription *Subscription) SetPayInfoSubscription(RebillID string, Active bool, Brand string, Last4 string) error {
	if !Active {
		return ErrPayMethodMustBeActive
	}

	if RebillID == "" {
		return ErrRebillIDNotValid
	}

	subscription.PaymentMethod.Active = Active
	subscription.PaymentMethod.Brand = Brand
	subscription.PaymentMethod.Last4 = Last4
	subscription.PaymentMethod.RebillID = RebillID

	return nil
}