package models

type Order struct {
	ID              int64  `db:"id" json:"id"`
	UserID          int64  `db:"user_id" json:"user_id"`
	PickupLocation  string `db:"pickup_location" json:"pickup_location"`
	DropoffLocation string `db:"dropoff_location" json:"dropoff_location"`
	PackageDetails  string `db:"package_details" json:"package_details"`
	DeliveryTime    string `db:"delivery_time" json:"delivery_time"`
	Status          string `db:"status" json:"status"`
}
