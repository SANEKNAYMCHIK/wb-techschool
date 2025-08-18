package models

import (
	"math/rand"
	"time"
)

type TestOrderGenerator struct {
	items []Item
}

func createItem() *Item {
	chrtId := rand.Intn(100000) + 1
	price := rand.Intn(100000) + 1
	sample := Item{
		ChrtID:      chrtId,
		TrackNumber: "WBILMTESTTRACK1",
		Price:       price,
		RID:         "ab4219087a764ae0btest",
		Name:        "Test Product",
		Sale:        0,
		Size:        "0",
		TotalPrice:  price,
		NmID:        1,
		Brand:       "Test Brand",
		Status:      202,
	}
	return &sample
}

func NewTestOrderGenerator() *TestOrderGenerator {
	// Create items for testing
	numItems := rand.Intn(5) + 1
	orderItems := make([]Item, numItems)
	for i := range numItems {
		orderItems[i] = *createItem()
	}
	return &TestOrderGenerator{
		items: orderItems,
	}
}

func (g *TestOrderGenerator) GenerateOrder() *Order {
	// Generate random UID
	orderUID := generateRandomString(19)

	// count total price
	var totalPrice int
	for _, item := range g.items {
		totalPrice += item.TotalPrice
	}

	return &Order{
		OrderUID:    orderUID,
		TrackNumber: "WBILMTESTTRACK",
		Entry:       "WBIL",
		Delivery: Delivery{
			Name:    "Test Testov",
			Phone:   "+9720000000",
			Zip:     "2639809",
			City:    "Kiryat Mozkin",
			Address: "Ploshad Mira 15",
			Region:  "Kraiot",
			Email:   "test@gmail.com",
		},
		Payment: Payment{
			Transaction:  orderUID,
			RequestID:    "",
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       totalPrice,
			PaymentDt:    time.Now().Unix(),
			Bank:         "alpha",
			DeliveryCost: 0,
			GoodsTotal:   totalPrice,
			CustomFee:    0,
		},
		Items:             g.items,
		Locale:            "en",
		InternalSignature: "",
		CustomerID:        "test",
		DeliveryService:   "meest",
		Shardkey:          "9",
		SmID:              99,
		DateCreated:       time.Now(),
		OofShard:          "1",
	}
}

func generateRandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
