package main

import (
	"whoam.xyz/ent/method"
	"whoam.xyz/ent/service"
)

// InitService initialize service related business
func InitService() {
	if ok, err := client.Service.Query().Where(service.IDEQ(MainServiceID)).Exist(ctx); !ok || err != nil {
		_, err = client.Service.Create().
			SetID(MainServiceID).
			SetName("whoam service").
			SetSubject("Support OAuth authorization, support service registration, support RAS.").
			SetDomain("https://www.whoam.xyz").
			SetCloneURI("https://github.com/excing/whoam.git").
			Save(ctx)

		if err != nil {
			panic(err)
		}
	}
}

type postServiceForm struct {
	ServiceID   string `schema:"serviceId,required"`
	ServiceName string `schema:"serviceName,required"`
	ServiceDesc string
	Domain      string `schema:"domain,required"`
	CloneURI    string `schema:"cloneUri,required"`
}

// PostServicer 提交服务注册
func PostServicer(c *Context) error {
	var form postServiceForm
	err := c.ParseForm(&form)
	if err != nil {
		return c.BadRequest(err.Error())
	}

	_, err = client.Service.Create().
		SetID(form.ServiceID).
		SetName(form.ServiceName).
		SetSubject(form.ServiceDesc).
		SetDomain(form.Domain).
		SetCloneURI(form.CloneURI).
		Save(ctx)

	return c.NoContent()
}

type serviceMethodBody struct {
	Name    string `json:"name"`
	Route   string `json:"route"`
	Subject string `json:"subject"`
	Scope   string `json:"scope"`
}

// PostServiceMethod receive service method submission
func PostServiceMethod(c *Context) error {
	service, err := client.Service.Query().
		Where(service.IDEQ(c.Param("id"))).
		Only(ctx)

	if err != nil {
		return c.BadRequest(err.Error())
	}

	var body []*serviceMethodBody

	err = c.ShouldBindJSON(&body)
	if err != nil {
		return c.BadRequest(err.Error())
	}

	tx, err := client.Tx(ctx)
	if err != nil {
		return c.InternalServerError(err.Error())
	}

	for _, m := range body {
		_, err = tx.Method.Create().
			SetName(m.Name).
			SetRoute(m.Route).
			SetSubject(m.Subject).
			SetScope(method.Scope(m.Scope)).
			SetOwner(service).
			Save(ctx)

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
