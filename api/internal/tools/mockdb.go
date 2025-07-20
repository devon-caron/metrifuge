package tools

import "time"

type mockDB struct{}

var mockLoginDetails = map[string]LoginDetails{
	"alex": {
		AuthToken: "123ABC",
		Username:  "alex",
	},
	"jason": {
		AuthToken: "456DEF",
		Username:  "jason",
	},
	"marie": {
		AuthToken: "789GHI",
		Username:  "marie",
	},
}

var mockCoinDetails = map[string]CoinDetails{
	"alex": {
		Coins:    100,
		Username: "alex",
	},
	"jason": {
		Coins:    200,
		Username: "jason",
	},
	"marie": {
		Coins:    300,
		Username: "marie",
	},
}

// The (d *mockDB) part is called the method receiver in Go. Let's break it down in more detail:

// Method Receiver:
// In Go, (d *mockDB) declares that this function is a method associated with the mockDB type.
// Pointer Receiver:
// The * before mockDB makes this a pointer receiver. This means the method receives a pointer to a mockDB instance, not a copy of it.
// Receiver Variable Name:
// d is the name given to the receiver within the method. It's similar to this or self in other languages, but in Go, you explicitly name it.
// Usage:
// Inside the method, you can use d to access fields or call other methods of the mockDB instance.
// Why Use a Pointer Receiver:

// Efficiency: Avoids copying the entire struct when the method is called.
// Ability to Modify: Allows the method to modify the original mockDB instance.
// Consistency: Often used when other methods of the type use pointer receivers.

// Syntax Variations:

// (d mockDB) without the * would be a value receiver.
// The receiver variable name d could be any valid identifier, but short names are conventional.

func (this *mockDB) GetUserLoginDetails(username string) *LoginDetails {
	// Simulate DB Call
	time.Sleep(time.Second * 1)

	var clientData = LoginDetails{}
	clientData, ok := mockLoginDetails[username]
	if !ok {
		return nil
	}

	return &clientData
}

func (this *mockDB) GetUserCoins(username string) *CoinDetails {
	// Simulate DB Call
	time.Sleep(time.Second * 1)

	var clientData = CoinDetails{}
	clientData, ok := mockCoinDetails[username]
	if !ok {
		return nil
	}

	return &clientData
}

func (this *mockDB) SetupDatabase() error {
	return nil
}
