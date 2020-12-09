package main

// InitRAS init ras service
func InitRAS() {
}

type newRAS struct {
	Subject     string
	PostURL     string
	RedirectURI string
}

// NewArticle can create a new article
func NewArticle(c *Context) error {
	return c.NoContent()
}

// NewAccord can create a new accord
func NewAccord(c *Context) error {
	return c.NoContent()
}

// NewRAS create a new random anonymous space
func NewRAS(c *Context) error {
	var form newRAS
	if err := c.ParseForm(&form); err != nil {
		return c.BadRequest(err.Error())
	}

	return c.NoContent()
}
