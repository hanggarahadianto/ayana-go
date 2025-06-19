package dto

import (
	"time"
)

type CustomerResponse struct {
	ID            string        `json:"id"`
	Name          string        `json:"name"`
	Address       string        `json:"address"`
	Phone         string        `json:"phone"`
	Status        string        `json:"status"`
	Marketer      string        `json:"marketer"`
	Amount        int64         `json:"amount"`
	PaymentMethod string        `json:"payment_method"`
	DateInputed   *time.Time    `json:"date_inputed"`
	HomeID        string        `json:"home_id"`
	ProductUnit   string        `json:"product_unit"`
	BankName      string        `json:"bank_name"`
	Home          *HomeResponse `json:"home"` // gunakan pointer biar null-safe
}

type HomeResponse struct {
	ID         string `json:"id"`
	ClusterID  string `json:"cluster_id"`
	Type       string `json:"type"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	Bathroom   int    `json:"bathroom"`
	Bedroom    int    `json:"bedroom"`
	Square     int    `json:"square"`
	Price      int64  `json:"price"`
	Quantity   int    `json:"quantity"`
	Status     string `json:"status"`
	Sequence   int    `json:"sequence"`
	StartPrice int64  `json:"start_price"`
}
