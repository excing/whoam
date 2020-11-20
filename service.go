package main

import (
	"golang.org/x/crypto/bcrypt"
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

type serviceForm struct {
	ServiceID   string `schema:"serviceId,required"`
	ServiceName string `schema:"serviceName,required"`
	ServiceDesc string
}

// PostServicer 提交服务注册
func PostServicer(c *Context) error {
	var form serviceForm
	err := c.ParseForm(&form)
	if err != nil {
		return c.BadRequest(err.Error())
	}

	var service Service
	err = db.Where("serviceId=?", form.ServiceID).Find(&service).Error

	if err == nil && 0 != service.ID {
		return c.Created("This service is already registered")
	}

	token, _ := New64BitUUID()

	service.ServiceID = form.ServiceID
	service.ServiceName = form.ServiceName
	service.ServiceDesc = form.ServiceDesc
	service.ServiceToken = token

	err = db.Create(&service).Error

	if err != nil {
		return c.InternalServerError(err.Error())
	}

	return c.Ok("Post successed")
}

// DeleteServicer 注销服务
func DeleteServicer(c *Context) error {
	serviceID := c.PostForm("serviceId")
	serviceToken := c.PostForm("serviceToken")

	if "" == serviceID {
		return c.BadRequest("serviceId is empty")
	} else if "" == serviceToken {
		return c.BadRequest("serviceToken is empty")
	}

	serviceProvider, providerOk := serviceMap[serviceID]

	if providerOk {
		err := bcrypt.CompareHashAndPassword([]byte(serviceProvider.ServiceTokenEncode), []byte(serviceToken))

		if err == nil {
			delete(serviceMap, serviceID)
			return c.NoContent()
		}
	}

	return c.Unauthorized("serviceToken is an invalid value")
}
