package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/thoas/go-funk"
)

// Note: This is just an example for a tutorial

var (
	VALID_AUTHENTICATIONS  = []string{"user", "admin", "subscriber"}
)

type User struct {
	Username 	string `form:"username" json:"username" xml:"username" binding:"required"`
	AuthType 	string `form:"authType" json:"authType" xml:"authType" binding:"required"`
}

func main() {

	router := gin.Default()

	store := sessions.NewCookieStore([]byte("sessionSuperSecret"))
	router.Use(sessions.Sessions("sessionName", store))

	api := router.Group("/api/v1")
	// no authentication endpoints
	{
		api.POST("/login", loginHandler)
		api.GET("/message/:msg", noAuthMessageHandler)
	}
	// basic authentication endpoints
	{
		basicAuth := api.Group("/")
		basicAuth.Use(AuthenticationRequired())
		{
			basicAuth.GET("/logout", logoutHandler)
		}
	}
	// admin authentication endpoints
	{
		adminAuth := api.Group("/admin")
		adminAuth.Use(AuthenticationRequired("admin"))
		{
			adminAuth.GET("/message/:msg", adminMessageHandler)
		}
	}
	// subscriber authentication endpoints
	{
		subscriberAuth := api.Group("/")
		subscriberAuth.Use(AuthenticationRequired("subscriber"))
		{
			subscriberAuth.GET("/subscriber/message/:msg", subscriberMessageHandler)
		}
	}

	apiV2 := router.Group("/api/v2")
	router.Use(AuthenticationRequired("admin", "subscriber"))
	// admin and subscriber authentication endpoints
	{
		apiV2.POST("/post/message/:msg", postMessageHandler)
	}

	router.Run(":8080")
}

func loginHandler(c *gin.Context) {
	var user User
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	session := sessions.Default(c)

	if strings.Trim(user.Username, " ") == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username can't be empty"})
	}
	if ! funk.ContainsString(VALID_AUTHENTICATIONS, user.AuthType) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid auth type"})
	}

	// Note: This is just an example, in real world AuthType would be set by business logic and not the user
	session.Set("user", user.Username)
	session.Set("authType", user.AuthType)

	err := session.Save()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate session token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "authentication successful"})
}

func logoutHandler(c *gin.Context) {
	session := sessions.Default(c)

	// this would only be hit if the user was authenticated
	session.Delete("user")
	session.Delete("authType")

	err := session.Save()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate session token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "successfully logged out"})

}

func noAuthMessageHandler(c *gin.Context) {
	msg := c.Param("msg")
	c.JSON(http.StatusOK, gin.H{"message": msg})
}

func subscriberMessageHandler(c *gin.Context) {
	msg := c.Param("msg")

	session := sessions.Default(c)
	user := session.Get("user")

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Hello Subscriber %s, here's your message: %s", user, msg)})
}

func adminMessageHandler(c *gin.Context) {
	msg := c.Param("msg")

	session := sessions.Default(c)
	user := session.Get("user")

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Hello Admin %s, here's your message: %s", user, msg)})
}

func postMessageHandler(c *gin.Context) {
	msg := c.Param("msg")

	session := sessions.Default(c)
	user := session.Get("user")

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Hello Admin/Subscriber %s, your message: %s will be posted", user, msg)})
}