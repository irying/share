package routers

import (
	"backend/api/controllers/user"
	"backend/api/middleware"
	"backend/conf"
	"backend/system/exception"
	"github.com/DeanThompson/ginpprof"
	"github.com/gin-gonic/contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitRun() {

	engine := gin.Default()
	ginpprof.Wrap(engine)
	engine.Use(middleware.Recovery(recoveryHandler))
	engine.Use(cors.Default())

	// backend route
	backend := engine.Group("/backend")
	{
		engine.LoadHTMLGlob("../views/*")
		backend.Static("/static/js", "../static/js")
		backend.Static("/assets", "../static/assets")
		backend.Static("/images", "../static/images")
		backend.GET("/login", func(c *gin.Context) {
			c.HTML(http.StatusOK, "login.tmpl", gin.H{})
		})
		backend.Use(middleware.JWT()).GET("/index", user.Index)
	}

	// api route
	engine.POST("/api/login", user.Login)
	api := engine.Group("/api")
	{
		api.Use(middleware.JWT())
		users := api.Group("/users")
		{
			users.GET("/profiles", user.Profile)
			users.POST("/profile_update", user.UpdateProfile)
		}
	}

	// service run
	if conf.Conf.App.Api.ApiTls == true {
		engine.RunTLS(conf.Conf.App.Api.ApiTlsAddr, "keys/server.crt", "keys/server.key")
	} else {
		engine.Run(conf.Conf.App.Api.ApiAddr)
	}
}

// 捕捉 panic
func recoveryHandler(c *gin.Context, err interface{}) {
	c.JSON(http.StatusInternalServerError, gin.H{"code": exception.SUCCESS, "message": err})
}
