package entities_test

import (
	"errors"
	"testing"
	"time"

	"pay-service/internal/entities"
)

func TestNewSubscription_OK(testLib *testing.T) {
	accessUntil := time.Now().Add(24 * time.Hour)
	subscription, err := entities.NewSubscription(
		"34234",
		"3245324",
		entities.PaymentMethodSnapshot{RebillID: "12", Active: true, Brand: "t-bank", Last4: "2121"},
		accessUntil,
	)

	if err != nil {
		testLib.Fatalf("New subscription unexpected error: %v", err)
	}

	if subscription.SubscriptionStatus != entities.SubscriptionStatusActive {
		testLib.Fatalf("Unexpected subscription status: %v", subscription.SubscriptionStatus)
	}
}

func TestNewSubscription_userIdError(testLib *testing.T) {
	_, err := entities.NewSubscription(
		"",
		"3245324",
		entities.PaymentMethodSnapshot{RebillID: "12", Active: true, Brand: "t-bank", Last4: "2121"},
		time.Now().Add(24*time.Hour),
	)

	if !errors.Is(err, entities.ErrInvalidUserId) {
		testLib.Fatalf("Have no error if empty user id: %v", err)
	}
}

func TestNewSubscription_SubscriptionIdError(testLib *testing.T) {
	_, err := entities.NewSubscription(
		"scasac",
		"",
		entities.PaymentMethodSnapshot{RebillID: "12", Active: true, Brand: "t-bank", Last4: "2121"},
		time.Now().Add(24*time.Hour),
	)

	if !errors.Is(err, entities.ErrInvalidSubscriptionId) {
		testLib.Fatalf("Have no error if empty rebill id: %v", err)
	}
}

func TestNewSubscription_RebillIdError(testLib *testing.T) {
	_, err := entities.NewSubscription(
		"scasac",
		"cssc",
		entities.PaymentMethodSnapshot{RebillID: "", Active: true, Brand: "t-bank", Last4: "2121"},
		time.Now().Add(24*time.Hour),
	)

	if !errors.Is(err, entities.ErrInvalidRebillId) {
		testLib.Fatalf("Have no error if empty rebill id: %v", err)
	}
}

func TestNewSubscription_PaymentMethodError(testLib *testing.T) {
	_, err := entities.NewSubscription(
		"scasac",
		"cssc",
		entities.PaymentMethodSnapshot{RebillID: "ascsac", Active: false, Brand: "t-bank", Last4: "2121"},
		time.Now().Add(24*time.Hour),
	)

	if !errors.Is(err, entities.ErrInvalidPayMethod) {
		testLib.Fatalf("Have no error if empty rebill id: %v", err)
	}
}