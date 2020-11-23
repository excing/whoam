package main

import (
	"gorm.io/gorm"
)

// Service basic information of restful api service
type Service struct {
	gorm.Model
	ServiceID    string
	ServiceName  string
	ServiceDesc  string
	ServiceToken string
}

type postServiceForm struct {
	ServiceID   string `schema:"serviceId,required"`
	ServiceName string `schema:"serviceName,required"`
	ServiceDesc string
}

// PostServicer 提交服务注册
func PostServicer(c *Context) error {
	var form postServiceForm
	err := c.ParseForm(&form)
	if err != nil {
		return c.BadRequest(err.Error())
	}

	var service Service
	if err = db.Where("serviceId=?", form.ServiceID).Find(&service).Error; err == nil && 0 != service.ID {
		return c.Created("This service is already registered")
	}

	token := New64BitID()

	service.ServiceID = form.ServiceID
	service.ServiceName = form.ServiceName
	service.ServiceDesc = form.ServiceDesc
	service.ServiceToken = token

	err = db.Create(&service).Error
	if err != nil {
		return c.InternalServerError(err.Error())
	}

	return c.Ok(token)
}

type deleteServiceForm struct {
	ServiceID    string `schema:"serviceId,required"`
	ServiceToken string `schema:"serviceToken,required"`
}

// DeleteServicer 注销服务
func DeleteServicer(c *Context) error {
	var form deleteServiceForm
	err := c.ParseForm(&form)
	if err != nil {
		return c.BadRequest(err.Error())
	}

	var service Service
	err = db.Where("serviceId=?", form.ServiceID).Find(&service).Error
	if err != nil {
		return c.Forbidden("Service not registered")
	}

	err = db.Delete(&service).Error
	if err != nil {
		return c.InternalServerError(err.Error())
	}

	return c.NoContent()
}
