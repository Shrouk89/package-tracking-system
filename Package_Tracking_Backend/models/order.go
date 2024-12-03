package models

type Order struct {
	ID              int64  `db:"id" json:"id"`
	UserID          int64  `db:"user_id" json:"user_id"`
	CourierID       int64  `db:"courier_id" json:"courier_id"`
	PickupLocation  string `db:"pickup_location" json:"pickupLocation"`
	DropoffLocation string `db:"dropoff_location" json:"dropoffLocation"`
	PackageDetails  string `db:"package_details" json:"packageDetails"`
	DeliveryTime    string `db:"delivery_time" json:"deliveryTime"`
	Status          string `db:"status" json:"status"`
}
