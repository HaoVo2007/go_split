package route

import (
	"go-split/internal/interface/http/handler"
	"go-split/internal/interface/middleware"

	"github.com/gin-gonic/gin"
)

func SetupGroupRoutes(
	router *gin.Engine,
	groupHandler *handler.GroupHandler,
) {
	groupRouter := router.Group("api/v1/groups")
	groupRouter.Use(middleware.AuthMiddleware())
	{
		groupRouter.POST("/", groupHandler.CreateGroup)
		groupRouter.GET("/", groupHandler.GetGroups)
		groupRouter.GET("/:id", groupHandler.GetGroupById)
		groupRouter.PUT("/:id", groupHandler.UpdateGroup)
		groupRouter.DELETE("/:id", groupHandler.DeleteGroup)
		
		groupRouter.POST("/:id/members", groupHandler.AddGroupMember)
		groupRouter.DELETE("/:id/members/:memberId", groupHandler.RemoveGroupMember)
	}
}
