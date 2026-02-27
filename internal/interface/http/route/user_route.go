package route

import (
	"go-split/internal/interface/http/handler"
	"go-split/internal/interface/middleware"

	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(
	router *gin.Engine,
	userHandler *handler.UserHandler,
) {
	userRouter := router.Group("api/v1/users")

	// auth
	authRouter := userRouter.Group("auth")
	{
		authRouter.POST("/login", userHandler.Login)
		authRouter.POST("/register", userHandler.Register)
		authRouter.POST("/refresh-token", userHandler.RefreshToken)
		authRouter.POST("/logout", middleware.AuthMiddleware(), userHandler.Logout)
		authRouter.POST("/update-profile", middleware.AuthMiddleware(), userHandler.UpdateProfile)
	}

	// admin
	adminRouter := userRouter.Group("admin")
	adminRouter.Use(middleware.AuthMiddleware())
	adminRouter.Use(middleware.AdminMiddleware())
	{
		adminRouter.POST("/", userHandler.CreateUser)
		adminRouter.GET("/list", userHandler.GetUsers)
		adminRouter.GET("/:id", userHandler.GetUserById)
		adminRouter.PUT("/:id", userHandler.UpdateUser)
		adminRouter.DELETE("/:id", userHandler.SoftDeleteUser)
	}

}
