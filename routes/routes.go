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
	// merchantauth.POST("/displaymenu", controllers.DisplayMenu)
    // merchantauth.POST("/menubycategory/:categoryid", controllers.GetFoodByCategory)
  merchantauth.POST("/numberofmenubycategory", controllers.FoodNumberByCategory)
//   merchantauth.POST("/displaycategory", controllers.DisplayCategory)
	
	
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
	menuroute.PATCH("/updatemenu/:id", controllers.UpdateMenu)
	menuroute.DELETE("/deletemenu/:id", controllers.DeleteMenu)

    categoryroute :=e.Group("/category")
	categoryroute.Use(middleware.ValidateToken)
	categoryroute.POST("/new",controllers.CreateCategory)
	categoryroute.GET("/all",controllers.GetCategory)
	categoryroute.PATCH("/:id",controllers.EditCategory)
	categoryroute.DELETE("/:id",controllers.DeleteCategory)
	categoryroute.GET("/foods/:categoryid",controllers.MerchantGetFoodByCategory)
	


	//user routes

	userroutes:=e.Group("/user")
	userroutes.POST("/displayallmenu",controllers.DisplayMenu)
	userroutes.POST("/displayallcategory",controllers.DisplayCategory)
	userroutes.POST("/menubycategory/:categoryid", controllers.GetFoodByCategory)
    userroutes.POST("/numberofmenubycategory", controllers.FoodNumberByCategory)
userroutes.POST("/fetchmenusbyfastingstatus",controllers.FetchMenusByFastingStatus)




	// menuroute.POST("/orderfood/:id", controllers.OrderFood)


}
