package main

import (
	"strconv"

	"whoam.xyz/ent/accord"
)

// InitRAS init ras service
func InitRAS() {
}

type newAccordForm struct {
	Name  string `schema:"name,required"`
	About string
}

// NewAccord can create a new accord
func NewAccord(c *Context) error {
	var form newAccordForm
	if err := c.ParseForm(&form); err != nil {
		return c.BadRequest(err.Error())
	}

	return c.NoContent()
}

type accordArticleBody struct {
	Subject string `json:"subject"`
	Note    string `json:"note"`
}

// PostAccordArticle adds Articles to Accord
func PostAccordArticle(c *Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return c.BadRequest("id is invalid")
	}

	_accord, err := client.Accord.Query().
		Where(accord.IDEQ(id)).
		Only(ctx)

	if err != nil {
		return c.BadRequest(err.Error())
	}

	var body []*accordArticleBody

	err = c.ShouldBindJSON(&body)
	if err != nil {
		return c.BadRequest(err.Error())
	}

	return c.Ok(&_accord)
}

type newRASForm struct {
	Subject     string
	PostURL     string
	RedirectURI string
}

// NewRAS create a new random anonymous space
func NewRAS(c *Context) error {
	var form newRASForm
	if err := c.ParseForm(&form); err != nil {
		return c.BadRequest(err.Error())
	}

	return c.NoContent()
}
