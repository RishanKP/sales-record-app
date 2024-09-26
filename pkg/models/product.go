package models

import "time"

type SaleRecord struct {
	OrderID         string    `json:"order_id" bson:"order_id"`
	ProductID       string    `json:"product_id" bson:"product_id"`
	CustomerID      string    `json:"customer_id" bson:"customer_id"`
	ProductName     string    `json:"product_name" bson:"product_name"`
	Category        string    `json:"category" bson:"category"`
	Region          string    `json:"region" bson:"region"`
	DateOfSale      time.Time `json:"date_of_sale" bson:"date_of_sale"`
	QuantitySold    int       `json:"quantity_sold" bson:"quantity_sold"`
	UnitPrice       float64   `json:"unit_price" bson:"unit_price"`
	Discount        float64   `json:"discount" bson:"discount"`
	ShippingCost    float64   `json:"shipping_cost" bson:"shipping_cost"`
	PaymentMethod   string    `json:"payment_method" bson:"payment_method"`
	CustomerName    string    `json:"customer_name" bson:"customer_name"`
	CustomerEmail   string    `json:"customer_email" bson:"customer_email"`
	CustomerAddress string    `json:"customer_address" bson:"customer_address"`
}
