package main

import (
	"gorm.io/gorm"
)

// Service basic information of restful api service
type Service struct {
	gorm.Model
	ServiceID string
	Name      string
	Token     string
	Describe  string
}

// ServiceMethod service method
type ServiceMethod struct {
	Name  string
	Scope string
	// Fields []ServiceMethodField
}

// ServiceMethodField service method field
// type ServiceMethodField struct {
// 	Name     string
// 	Kind     string
// 	Default  string
// 	Describe string
// 	Required string
// }

// InitService initialize service related business
func InitService() {
	db.AutoMigrate(&Service{})
}

// ServiceAuthorize service authorize
func ServiceAuthorize(c *Context, service *Service) bool {
	token := c.GetHeader("Authorization")
	if db.Where("token=?", token).Find(service).Error != nil || 0 == service.ID {
		return false
	}

	return true
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
	if err = db.Where("service_id=?", form.ServiceID).Find(&service).Error; err == nil && 0 != service.ID {
		return c.Created("This service is already registered")
	}

	token := New64BitID()

	service.ServiceID = form.ServiceID
	service.Name = form.ServiceName
	service.Describe = form.ServiceDesc
	service.Token = token

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
	var service Service
	if !ServiceAuthorize(c, &service) {
		return c.Unauthorized("Invalid token, please refresh the access token with Referh_token")
	}

	err := db.Delete(&service).Error
	if err != nil {
		return c.InternalServerError(err.Error())
	}

	return c.NoContent()
}

// PostServiceMethod receive service method submission
func PostServiceMethod(c *Context) error {
	var service Service
	if !ServiceAuthorize(c, &service) {
		return c.Unauthorized("Invalid token, please refresh the access token with Referh_token")
	}

	methods, err := c.FormArray("methods")
	if err != nil {
		return c.BadRequest(err.Error())
	}

	return c.Ok(methods)
}
