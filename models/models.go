package models

// Person ...Defines the structure of our data
type Person struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}

// Config ...Config layout
type Config struct {
	DBURL string `json:"DBURL"`
}
