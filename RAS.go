package main

import (
	"strconv"
	"strings"

	"whoam.xyz/ent"
	"whoam.xyz/ent/accord"
	"whoam.xyz/ent/article"
	"whoam.xyz/ent/schema"
	"whoam.xyz/ent/user"
	"whoam.xyz/ent/vote"
)

var rasBox *Box

// InitRAS init ras service
func InitRAS() {
	// size: 3M
	// default timeout: 30day
	rasBox = NewBox(3*1024*1024, 30*24*60*60)
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

// GetArticle returns 200 and the specified Article,
// if the specified Article cannot be found, it returns a non-2XX status code
func GetArticle(c *Context) error {
	id, err := c.ParamInt("id")
	if err != nil {
		return c.BadRequest(err.Error())
	}

	article, err := client.Article.Query().Where(article.IDEQ(id)).Only(ctx)

	if err != nil {
		return c.NotFound(err.Error())
	}

	return c.Ok(&article)
}

// GetArticles returns 200 and the Article List of the specified page number,
// if not, it returns a non-2XX status code
func GetArticles(c *Context) error {
	start, ok, err := c.QueryInt("start")
	if err != nil {
		if ok {
			return c.BadRequest(err.Error())
		}
	}
	count, ok, err := c.QueryInt("count")
	if err != nil {
		if ok {
			return c.BadRequest(err.Error())
		}
		count = 10
	}

	articles, err := client.Article.Query().Order(ent.Desc(article.FieldID)).Offset(start).Limit(count).All(ctx)

	if err != nil {
		return c.NotFound(err.Error())
	}

	return c.Ok(&articles)
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
		return c.NotFound("Invalid article ids: %v", missIDs)
	}
	if err != nil {
		return c.NotFound(err.Error())
	}

	accord, err := client.Accord.Create().SetName(form.Name).SetAbout(form.About).AddArticleIDs(articleIDs...).Save(ctx)
	if err != nil {
		return c.InternalServerError(err.Error())
	}

	return c.Ok(accord.ID)
}

// GetAccord returns 200 and the specified accord,
// if cannot be found, it returns a non-2XX status code
func GetAccord(c *Context) error {
	id, err := c.ParamInt("id")
	if err != nil {
		return c.BadRequest(err.Error())
	}

	_accord, err := client.Accord.Query().Where(accord.IDEQ(id)).Only(ctx)

	if err != nil {
		return c.NotFound(err.Error())
	}

	return c.Ok(&_accord)
}

// GetAccords returns 200 and the accords of the specified page number,
// if cannot be found, it returns a non-2XX status code
func GetAccords(c *Context) error {
	start, ok, err := c.QueryInt("start")
	if err != nil {
		if ok {
			return c.BadRequest(err.Error())
		}
	}
	count, ok, err := c.QueryInt("count")
	if err != nil {
		if ok {
			return c.BadRequest(err.Error())
		}
		count = 10
	}

	accords, err := client.Accord.Query().Order(ent.Desc(accord.FieldID)).Offset(start).Limit(count).All(ctx)

	if err != nil {
		return c.NotFound(err.Error())
	}

	return c.Ok(&accords)
}

// GetAccordArticles returns 200 and the articles of the specified accord,
// if cannot be found, it returns a non-2XX status code
func GetAccordArticles(c *Context) error {
	id, err := c.ParamInt("id")
	if err != nil {
		return c.BadRequest(err.Error())
	}

	_accord, err := client.Accord.Query().Where(accord.IDEQ(id)).Only(ctx)

	if err != nil {
		return c.NotFound(err.Error())
	}

	articles, err := _accord.QueryArticles().All(ctx)

	if err != nil {
		return c.NotFound(err.Error())
	}

	return c.Ok(&articles)
}

type newRASForm struct {
	Subject     string `schema:"subject,required"`
	PostURI     string `schema:"post_uri,required"`
	RedirectURI string `schema:"redirect_uri,required"`
	Accord      int    `schema:"accord,required"`
}

// NewRAS can create a new RAS
func NewRAS(c *Context) error {
	var form newRASForm
	err := c.ParseForm(&form)
	if err != nil {
		return c.BadRequest(err.Error())
	}

	tx, err := client.Tx(ctx)

	ras, err := tx.RAS.
		Create().
		SetSubject(form.Subject).
		SetPostURI(form.PostURI).
		SetRedirectURI(form.RedirectURI).
		SetAccordID(form.Accord).
		Save(ctx)

	if err != nil {
		c.InternalServerError(err.Error())
	}

	users, err := tx.User.Query().
		Order(schema.Rand()).
		Limit(10).
		Select(user.FieldID).
		Ints(ctx)
	if err != nil {
		c.InternalServerError(err.Error())
	}

	for _, v := range users {
		var m []*ent.RAS
		err = rasBox.ValI(v, &m)
		if err != nil {
			return c.InternalServerError(err.Error())
		}

		m = append(m, ras)
		err = rasBox.SetValI(v, &m)
		if err != nil {
			return c.InternalServerError(err.Error())
		}
	}

	err = tx.Commit()
	if err != nil {
		return c.InternalServerError(err.Error())
	}

	return c.NoContent()
}

// GetRAS get RAS
func GetRAS(c *Context) error {
	id, err := c.ParamInt("userId")

	if err != nil {
		c.BadRequest(err.Error())
	}

	var m []*ent.RAS
	err = rasBox.ValI(id, &m)
	if err != nil {
		return c.InternalServerError(err.Error())
	}

	if 0 == len(m) {
		return c.NoContent()
	}

	return c.Ok(&m[0])
}

// GetVotes get all vote that specified RAS post
func GetVotes(c *Context) error {
	return c.NoContent()
}

type postVoteForm struct {
	RASID int
	State string
	Note  string
}

// VoteRAS vote RAS
func VoteRAS(c *Context) error {
	token, err := FilterJWTToken(c.GetHeader("Authorization"), signingKey)
	if err != nil {
		return c.Unauthorized("Invalid token, please login again")
	}

	userID := int(token.OtherID)

	var form postVoteForm
	err = c.ParseForm(&form)
	if err != nil {
		return c.BadRequest(err.Error())
	}

	m := make(map[int]*ent.RAS)
	err = rasBox.ValI(userID, &m)
	if err != nil {
		return c.BadRequest(err.Error())
	}

	ras, ok := m[form.RASID]
	if !ok {
		return c.BadRequest(err.Error())
	}

	_, err = client.Vote.Create().
		SetDstID(ras.ID).
		SetState(vote.State(form.State)).
		SetSubject(form.Note).
		Save(ctx)
	if err != nil {
		c.InternalServerError(err.Error())
	}

	return c.NoContent()
}
