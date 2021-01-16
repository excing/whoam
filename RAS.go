package main

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"whoam.xyz/ent"
	"whoam.xyz/ent/accord"
	"whoam.xyz/ent/article"
	"whoam.xyz/ent/ras"
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
	Subject     string `json:"subject"`
	PostURI     string `json:"post_uri"`
	RedirectURI string `json:"redirect_uri"`
	Accord      int    `json:"accord"`
	Voters      []int  `json:"voters"`
}

type rasVote struct {
	UserID int
	State  string
	Note   string
}

// NewRAS can create a new RAS
// The voter list is determined by the creator of RAS
func NewRAS(c *Context) error {
	var body newRASForm
	err := c.ShouldBindJSON(&body)
	if err != nil {
		return c.BadRequest(err.Error())
	}

	if len(body.Voters) != 10 {
		return c.BadRequest("Please fix the number of voters to 10")
	}

	err = WithTx(ctx, client, func(tx *ent.Tx) error {
		ras, err := tx.RAS.
			Create().
			SetSubject(body.Subject).
			SetPostURI(body.PostURI).
			SetRedirectURI(body.RedirectURI).
			SetAccordID(body.Accord).
			Save(ctx)
		if err != nil {
			return err
		}

		// Issues: Voters are not stored, if the service restarts, there is a loss
		var votes []*rasVote
		for _, v := range body.Voters {
			votes = append(votes, &rasVote{UserID: v})
			var _rass []int
			err = rasBox.ValI(v, &_rass)
			if err != nil {
				return err
			}

			_rass = append(_rass, ras.ID)
			err = rasBox.SetValI(v, &_rass)
			if err != nil {
				return err
			}
		}
		err = rasBox.Val(strconv.Itoa(ras.ID), &votes)
		return err
	})

	if err != nil {
		return c.InternalServerError(err.Error())
	}

	// `RAS.id` will be sent with the voting result when calling back `RedirectURI`
	// A `post` may produce multiple votes, when the voting result is abstention
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

// GetRasVotes get all vote that specified RAS
func GetRasVotes(c *Context) error {
	rasID, err := c.ParamInt("rasId")
	if err != nil {
		c.BadRequest(err.Error())
	}

	return c.Ok(rasID)
}

// GetUserVotes get all vote that specified User
func GetUserVotes(c *Context) error {
	rasID, err := c.ParamInt("userId")
	if err != nil {
		c.BadRequest(err.Error())
	}

	return c.Ok(rasID)
}

// GetPostVotes get all vote that specified Post
func GetPostVotes(c *Context) error {
	rasID, err := c.ParamInt("postUri")
	if err != nil {
		c.BadRequest(err.Error())
	}

	return c.Ok(rasID)
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

	var _votes []*rasVote
	err = rasBox.Val(strconv.Itoa(form.RASID), &_votes)
	if err != nil {
		return c.BadRequest(err.Error())
	}

	// Has the vote been completed
	completed := true
	for _, _vote := range _votes {
		if _vote.UserID == userID {
			if vote.StateValidator(vote.State(_vote.State)) != nil {
				return c.Conflict("Voted on this post")
			}
			_vote.State = form.State
			_vote.Note = form.Note
		} else if vote.StateValidator(vote.State(_vote.State)) != nil {
			completed = false
		}
	}

	if !completed {
		// Refresh voting list
		err = rasBox.SetVal(strconv.Itoa(form.RASID), &_votes)
		if err != nil {
			return c.InternalServerError(err.Error())
		}

		return c.NoContent()
	}

	// If the voting is completed, process and save the voting result

	err = WithTx(ctx, client, func(tx *ent.Tx) error {
		voteCreates := make([]*ent.VoteCreate, len(_votes))

		_subject := ""
		allowedCount := 0
		rejectedCount := 0
		abstainedCount := 0
		for i, _vote := range _votes {
			_subject += _vote.Note + ";"
			state := vote.State(_vote.State)
			switch state {
			case vote.StateAllowed:
				allowedCount++
			case vote.StateRejected:
				rejectedCount++
			case vote.StateAbstained:
				abstainedCount++
			}

			voteCreates[i] = tx.Vote.Create().SetDstID(form.RASID).SetState(state).SetSubject(form.Note)
		}

		_state := ras.StateAbstained
		if 4 <= allowedCount && 4 <= rejectedCount {
			if allowedCount == rejectedCount || allowedCount <= rejectedCount {
				_state = ras.StateAllowed
			} else {
				_state = ras.StateRejected
			}
		} else if 4 <= allowedCount {
			_state = ras.StateAllowed
		} else if 4 <= rejectedCount {
			_state = ras.StateRejected
		} else {
			_state = ras.StateAbstained
		}

		_, err := tx.RAS.Update().Where(ras.IDEQ(form.RASID)).SetState(_state).SetSubject(_subject).Save(ctx)
		if err != nil {
			return err
		}

		rasRedirectURI, err := tx.RAS.Query().Where(ras.IDEQ(form.RASID)).Select(ras.FieldRedirectURI).String(ctx)
		if err != nil {
			return err
		}

		_resp := url.Values{
			"state":   {_state.String()},
			"subject": {_subject},
			"rasId":   {strconv.Itoa(form.RASID)},
		}
		_, err = http.PostForm(rasRedirectURI, _resp)
		if err != nil {
			return err
		}

		_, err = tx.Vote.CreateBulk(voteCreates...).Save(ctx)
		if err != nil {
			return err
		}

		ok := rasBox.DelString(strconv.Itoa(form.RASID))
		print("Remove RAS(", form.RASID, "): ", ok)

		// todo Unhandled cache exception
		for _, _vote := range _votes {
			var _ids []int
			rasBox.ValI(_vote.UserID, &_ids)
			if 1 == len(_ids) {
				rasBox.DelInt(int64(_vote.UserID))
			} else {
				news := make([]int, len(_ids)-1)
				for i, _id := range _ids {
					if form.RASID == _id {
						i--
						continue
					}
					news[i] = _id
				}
				rasBox.SetValI(_vote.UserID, &news)
			}
		}

		return err
	})

	if err != nil {
		return c.InternalServerError(err.Error())
	}

	return c.NoContent()
}
