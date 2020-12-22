package app

import (
	"fmt"
	"net/http"

	"github.com/victoryang/kubernetes-cicd/auth"

	"github.com/gin-gonic/gin"
)

// TODO: 重构思路: 固化流程，将参数获取和逻辑部分作为参数传进去

func login(c *gin.Context) {
	// Get param
	param := &auth.LoginParam{}
	err := c.BindJSON(param)
	if err != nil {
		err = fmt.Errorf("Wrong Login Param: %s", err)
		HttpResponseWithBadRequest(c, err)
		return
	}
	c.Set("username", param.Username)

	// Logic
	var userinfo *auth.User
	userinfo, err = auth.Login(param)
	if err != nil {
		err = fmt.Errorf("Authorized Failed: %v", err)
		HttpResponseWithBadRequest(c, err)
		return
	}
	HttpResponseWithSuccess(c, userinfo)

}

func authorize() gin.HandlerFunc {

	return func (c *gin.Context) {

		token, err := c.Cookie("token")
		if err != nil {
			err = fmt.Errorf("Get Token failed: %v", err)
			HttpResponseWithForbidden(c, err)
			return
		}
		var user *auth.User
		user, err = auth.GetUserByToken(token)
		if err != nil {
			err = fmt.Errorf("Authorized failed: %v", err)
			HttpResponseWithForbidden(c, err)
			return
		}
		c.Set("username", user.Name)
		return
	}
}