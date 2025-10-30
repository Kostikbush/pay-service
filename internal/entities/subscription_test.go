package subscription

import (
	"testing"
)

func TestNewSubscription_OK(testLib *testing.T) {
	subscription, err := CreateActiveSubscription("34234", "3245324", PaymentMethodSnapshot{RebillID: "12", Active: true, Brand: "t-bank", Last4: "2121"})

	if err != nil {
		testLib.Fatalf("New subscription unexpected error: %v", err)
	}

	if subscription.SubscriptionStatus != SubscriptionStatusActive {
		testLib.Fatalf("Unexpected subscription status: %v", subscription.SubscriptionStatus)
	}
}

func TestNewSubscription_userIdError(testLib *testing.T) {
	_, err := CreateActiveSubscription("", "3245324", PaymentMethodSnapshot{RebillID: "12", Active: true, Brand: "t-bank", Last4: "2121"})

	if(err != ErrUserIdNotValid) {
		testLib.Fatalf("Have no error if empty user id: %v", err)
	}
}