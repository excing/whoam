package main

import (
	"whoam.xyz/ent"
	"whoam.xyz/ent/method"
	"whoam.xyz/ent/service"
)

// InitService initialize service related business
func InitService() {
	_service, err := client.Service.Query().Where(service.IDEQ(MainServiceID)).First(ctx)
	if err != nil {
		_service, err = client.Service.Create().
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

	methods := []*serviceMethodBody{
		{
			Name:    "GetUser",
			Subject: "获取你的基本信息",
			Route:   "/user",
			Scope:   method.ScopePrivate,
		},
		{
			Name:    "GetOAuthState",
			Subject: "获取你的授权状态",
			Route:   "/user/oauth/state",
			Scope:   method.ScopePrivate,
		},
		{
			Name:    "PostService",
			Subject: "注册服务到 WHOAM",
			Route:   "/service",
			Scope:   method.ScopePublic,
		},
		{
			Name:    "PostServiceMethod",
			Subject: "为服务添加方法",
			Route:   "/service/:id/method",
			Scope:   method.ScopePublic,
		},
		{
			Name:    "PostServicePermission",
			Subject: "为服务添加所需的权限",
			Route:   "/service/:id/permission",
			Scope:   method.ScopePublic,
		},
	}

	var _methodCreates []*ent.MethodCreate
	for _, _method := range methods {
		_, err := client.Method.Update().
			Where(
				method.And(
					method.NameEQ(_method.Name),
					method.HasOwnerWith(service.IDEQ(MainServiceID)))).
			SetSubject(_method.Subject).
			SetRoute(_method.Route).
			SetScope(_method.Scope).
			Save(ctx)

		if err != nil {
			_methodCreate := client.Method.Create().
				SetName(_method.Name).
				SetSubject(_method.Subject).
				SetRoute(_method.Route).
				SetScope(_method.Scope).
				SetOwner(_service)

			_methodCreates = append(_methodCreates, _methodCreate)
		}
	}

	if _, err := client.Method.CreateBulk(_methodCreates...).Save(ctx); err != nil {
		panic(err)
	}
}

type postServiceForm struct {
	ServiceID   string `json:"service_id" binding:"required"`
	ServiceName string `json:"service_name" binding:"required"`
	ServiceDesc string `json:"service_desc"`
	Domain      string `json:"domain" binding:"required,url"`
	CloneURI    string `json:"clone_uri" binding:"required"`
}

// PostService 提交服务注册
func PostService(c *Context) error {
	var form postServiceForm
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

type serviceMethodBody struct {
	Name    string       `json:"name" binding:"required"`
	Route   string       `json:"route" binding:"required"`
	Subject string       `json:"subject" binding:"required"`
	Scope   method.Scope `json:"scope" binding:"required"`
}

// PostServiceMethod receive service method submission
func PostServiceMethod(c *Context) error {
	id, err := c.GetQueryString("id")
	if id == "" {
		return c.BadRequest(err.Error())
	}

	service, err := client.Service.Query().
		Where(service.IDEQ(id)).
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

type servicePermissionBody struct {
	ServiceID   string   `json:"serviceId" binding:"required"`
	Permissions []string `json:"permissions" binding:"required"`
}

// PostServicePermission receives the list of permissions required to add the specified service
func PostServicePermission(c *Context) error {
	id, err := c.GetQueryString("id")
	if id == "" {
		return c.BadRequest(err.Error())
	}

	src, err := client.Service.Query().
		Where(service.IDEQ(id)).
		Only(ctx)

	if err != nil {
		return c.BadRequest(err.Error())
	}

	var body []*servicePermissionBody

	err = c.ShouldBindJSON(&body)
	if err != nil {
		return c.BadRequest(err.Error())
	}

	srcUpdate := src.Update()

	for _, v := range body {
		_methods, err := client.Service.
			Query().
			Where(service.IDEQ(v.ServiceID)).
			QueryMethods().
			Where(method.NameIn(v.Permissions...)).
			All(ctx)

		if err != nil {
			return c.BadRequest(err.Error())
		}

		srcUpdate.AddPermissions(_methods...)
	}

	_, err = srcUpdate.Save(ctx)
	if err != nil {
		return c.InternalServerError(err.Error())
	}

	return c.NoContent()
}
