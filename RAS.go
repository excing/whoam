package main

import "time"

// RAS random anonymous space
type RAS struct {
	ID        uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT:10000"`
	Subject   string    `json:"subject"`
	Post      string    `json:"post"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"-"`
}

// InitRAS init ras service
func InitRAS() {
	db.AutoMigrate(&RAS{})
}

// NewRAS create a new random anonymous space
func NewRAS(c *Context) error {
	return c.NoContent()
}
