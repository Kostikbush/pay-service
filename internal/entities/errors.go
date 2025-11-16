package entities

import "errors"

var (
	ErrInvalidParam = errors.New("param in not valid")
	ErrInvalidUserId = errors.New("user id not valid")
	ErrInvalidSubscriptionId = errors.New("subscription id not valid")
	ErrInvalidPayMethod = errors.New("payment method not valid")
	ErrInvalidRebillId = errors.New("rebill id not valid")
)