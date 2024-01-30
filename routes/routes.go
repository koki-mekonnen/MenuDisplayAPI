package routes

import (
	"foodorderapi/internals/controllers"
	"foodorderapi/middleware"

	"github.com/labstack/echo/v4"
)

func Foodorderroutes(e *echo.Echo) {
	// ADMIN ROUTES
	adminroute := e.Group("/admin")
	adminroute.POST("/", controllers.Signupadmin)
	adminroute.POST("/login", controllers.Signinadmin)

	admin := e.Group("/admins")
	admin.Use(middleware.ValidateToken)
	admin.PUT("/update/:id", controllers.UpdateAdmin)
	admin.DELETE("/delete/:id", controllers.DeleteAdmin)
	admin.POST("/logout", controllers.Logout)

	// controllers routes to merchants
	admin.POST("/signupmerchant", controllers.Signupmerchant)
	admin.GET("/allmerchant/", controllers.Getallmerchant)
	admin.GET("/singlemerchant/:id", controllers.Singlemerchant)
	admin.PUT("/updatemerchant/:id", controllers.UpdateMerchantbyAdmin)
	admin.DELETE("/deletemerchant/:id", controllers.DeleteMerchant)
	admin.PUT("/deactivatemerchant/:id", controllers.DeactivateMerchant)
	admin.PUT("/activatemerchant/:id", controllers.ActivateMerchant)

	// MERCHANT ROUTES

	merchantauth := e.Group("/merchant")
	merchantauth.POST("/", controllers.Signin)
	merchantauth.POST("/forgetpassword", controllers.Forgetpassword)
	merchantauth.POST("/displaymenu", controllers.DisplayMenu)

	merchantroute := e.Group("/merchants")

	// ROUTES WHICH NEED MERCHANT TOKEN
	merchantroute.Use(middleware.ValidateToken)
	merchantroute.POST("/logout", controllers.Logout)
	merchantroute.PUT("/updateprofile/:id", controllers.UpdateMerchant)

	menuroute := e.Group("/merchantmenu")
	menuroute.Use(middleware.ValidateToken)
	menuroute.POST("/addnewmenu", controllers.CreateMenu)
	menuroute.GET("/getallmenus", controllers.ShowAllMenus)
	menuroute.GET("/getsinglemenu/:id", controllers.GetFood)
	menuroute.PUT("/updatemenu/:id", controllers.UpdateMenu)
	menuroute.DELETE("/deletemenu/:id", controllers.DeleteMenu)
	menuroute.POST("/orderfood/:id", controllers.OrderFood)

}
