package main

import (
	"time"
)

// RAS random anonymous space
type RAS struct {
	ID          uint `gorm:"primary_key;AUTO_INCREMENT:10000"`
	Subject     string
	PostURL     string
	RedirectURI string
	Status      int
	CreatedAt   *time.Time
}

// RASVote RAS vote
// type RASVote struct {
// 	ID        uint
// 	UserID    uint
// 	RASID     uint
// 	RAS       RAS
// 	Status    int
// 	Remark    string
// 	CreatedAt *time.Time
// }

// InitRAS init ras service
func InitRAS() {
	db.AutoMigrate(&RAS{})
}

type newRAS struct {
	Subject     string
	PostURL     string
	RedirectURI string
}

// NewRAS create a new random anonymous space
func NewRAS(c *Context) error {
	var form newRAS
	if err := c.ParseForm(&form); err != nil {
		return c.BadRequest(err.Error())
	}

	ras := RAS{
		Subject:     form.Subject,
		PostURL:     form.PostURL,
		RedirectURI: form.RedirectURI,
	}

	if err := db.Create(&ras).Error; err != nil {
		return c.InternalServerError(err.Error())
	}

	return c.NoContent()
}
