package main

import (
	"strconv"
	"strings"

	"whoam.xyz/ent/article"
)

// InitRAS init ras service
func InitRAS() {
}

type newArticleForm struct {
	Subject string `schema:"subject,required"`
	Note    string
}

// NewArticle can create a new article
func NewArticle(c *Context) error {
	var form newArticleForm
	if err := c.ParseForm(&form); err != nil {
		return c.BadRequest(err.Error())
	}

	article, err := client.Article.Create().SetSubject(form.Subject).SetNote(form.Note).Save(ctx)

	if err != nil {
		return c.InternalServerError(err.Error())
	}

	return c.Ok(strconv.Itoa(article.ID))
}

type newAccordForm struct {
	Name     string `schema:"name,required"`
	About    string
	Articles string `schema:"articles,required"`
}

// NewAccord can create a new accord
func NewAccord(c *Context) error {
	var form newAccordForm
	if err := c.ParseForm(&form); err != nil {
		return c.BadRequest(err.Error())
	}

	articleIDStr := strings.Split(form.Articles, ",")
	articleIDs := make([]int, len(articleIDStr))
	for i, v := range articleIDStr {
		id, err := strconv.Atoi(strings.Trim(v, " "))
		if err != nil {
			return c.BadRequest(err.Error())
		}

		articleIDs[i] = id
	}

	ids, err := client.Article.Query().Where(article.IDIn(articleIDs...)).IDs(ctx)
	if len(ids) < len(articleIDs) {
		missIDs := make([]int, len(articleIDs)-len(ids))
		i := 0
		for _, n := range articleIDs {
			miss := true
			for _, m := range ids {
				if n == m {
					miss = false
					continue
				}
			}
			if miss {
				missIDs[i] = n
				i++
			}
		}
		return c.BadRequest("Invalid article ids: %v", missIDs)
	}
	if err != nil {
		return c.BadRequest(err.Error())
	}

	accord, err := client.Accord.Create().SetName(form.Name).SetAbout(form.About).AddArticleIDs(articleIDs...).Save(ctx)
	if err != nil {
		return c.InternalServerError(err.Error())
	}

	return c.Ok(accord.ID)
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
