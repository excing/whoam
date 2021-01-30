package main

import (
	"whoam.xyz/ent/service"
)

// InitService initialize service related business
func InitService() {
	_, err := client.Service.Query().Where(service.IDEQ(MainServiceID)).First(ctx)
	if err != nil {
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

// PostService 提交服务注册
func PostService(c *Context) error {
	var form struct {
		ServiceID   string `json:"service_id" binding:"required"`
		ServiceName string `json:"service_name" binding:"required"`
		ServiceDesc string `json:"service_desc"`
		Domain      string `json:"domain" binding:"required,url"`
		CloneURI    string `json:"clone_uri" binding:"required"`
	}
	err := c.ShouldBindJSON(&form)
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

	if err != nil {
		return c.InternalServerError(err.Error())
	}

	return c.NoContent()
}
