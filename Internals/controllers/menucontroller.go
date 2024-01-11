package controllers

import (
	"foodorderapi/internals/config"
	"foodorderapi/internals/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

// CreateMenu creates a new menu item

func CreateMenu(c echo.Context) error {

	role := c.Get("role").(string)

	// Check if the role is merchant
	if role != "merchant" {
		data := map[string]interface{}{
			"message": "Access denied. Only merchants can perform this operation.",
		}
		return c.JSON(http.StatusForbidden, data)
	}

	db := config.DB()

	var merchants *models.Merchant
	merchantID := c.Get("merchantID").(string)

	if res := db.Where("id = ?", merchantID).Find(&merchants); res.Error != nil {
		data := map[string]interface{}{
			"message": "Merchant not found",
		}
		return c.JSON(http.StatusInternalServerError, data)
	}

	var menu *models.Menu

	if err := c.Bind(&menu); err != nil {
		return c.String(http.StatusBadRequest, "Invalid request payload")
	}

	newmenu := &models.Menu{
		FoodName:          menu.FoodName,
		Description:       menu.Description,
		Price:             menu.Price,
		Image:menu.Image,
		MerchantID:        merchants.Id,
		MerchantShortCode: merchants.MerchantShortcode,
	}

	if err := db.Create(&newmenu).Error; err != nil {
		return c.String(http.StatusInternalServerError, "Failed to create menu")
	}

	return c.JSON(http.StatusCreated, newmenu)
}

func ShowAllMenus(c echo.Context) error {

	role := c.Get("role").(string)

	// Check if the role is merchant
	if role != "merchant" {
		data := map[string]interface{}{
			"message": "Access denied. Only merchants can perform this operation.",
		}
		return c.JSON(http.StatusForbidden, data)
	}

	db := config.DB()
	var merchant *models.Merchant
	var menus []models.Menu

	merchantID := c.Get("merchantID").(string)

	if res := db.Where("id = ?", merchantID).Find(&merchant); res.Error != nil {
		data := map[string]interface{}{
			"message": "Merchant not found",
		}
		return c.JSON(http.StatusInternalServerError, data)
	}

	if res := db.Where("merchant_id = ?", merchantID).Find(&menus); res.Error != nil {
		data := map[string]interface{}{
			"message": "Merchant not found",
		}
		return c.JSON(http.StatusInternalServerError, data)
	}




	return c.JSON(http.StatusOK, menus)

}

func GetFood(c echo.Context) error {
	db := config.DB()
	id := c.Param("id")

	var food models.Menu

	if res := db.Where("id = ?", id).First(&food); res.Error != nil {
		data := map[string]interface{}{
			"message": res.Error.Error(),
		}

		return c.JSON(http.StatusNotFound, data)
	}

	return c.JSON(http.StatusOK, food)
}

func UpdateMenu(c echo.Context) error {
	role := c.Get("role").(string)

	// Check if the role is merchant
	if role != "merchant" {
		data := map[string]interface{}{
			"message": "Access denied. Only merchants can perform this operation.",
		}
		return c.JSON(http.StatusForbidden, data)
	}

	db := config.DB()

	menuID := c.Param("menuID")

	var menu models.Menu
	var payload models.UpdateMenu
	var merchant *models.Merchant

	merchantID := c.Get("merchantID").(string)

	if res := db.Where("id = ?", merchantID).Find(&merchant); res.Error != nil {
		data := map[string]interface{}{
			"message": "Merchant not found",
		}
		return c.JSON(http.StatusInternalServerError, data)
	}

	if res := db.Where("merchant_id = ?", merchantID).Find(&menu); res.Error != nil {
		data := map[string]interface{}{
			"message": "Merchant not found",
		}
		return c.JSON(http.StatusInternalServerError, data)
	}

	if res := db.Where("id=?", menuID).Find(&menu); res.Error != nil {
		data := map[string]interface{}{
			"message": "Menu not found",
		}
		return c.JSON(http.StatusInternalServerError, data)
	}

	if err := c.Bind(&payload); err != nil {
		return c.String(http.StatusBadRequest, "Invalid request payload")
	}

	menu.FoodName = payload.FoodName
	menu.Description = payload.Description
	menu.Price = payload.Price
	menu.Image=payload.Image

	if err := db.Save(&menu).Error; err != nil {
		return c.String(http.StatusInternalServerError, "Failed to update menu")
	}

	return c.JSON(http.StatusOK, menu)
}

func DeleteMenu(c echo.Context) error {

	role := c.Get("role").(string)

	// Check if the role is merchant
	if role != "merchant" {
		data := map[string]interface{}{
			"message": "Access denied. Only merchants can perform this operation.",
		}
		return c.JSON(http.StatusForbidden, data)
	}
	db := config.DB()

	id := c.Param("id")

	var menu models.Menu
	var merchant models.Merchant
	merchantID := c.Get("merchantID").(string)

	if res := db.Where("id = ?", merchantID).Find(&merchant); res.Error != nil {
		data := map[string]interface{}{
			"message": "Merchant not found",
		}
		return c.JSON(http.StatusInternalServerError, data)
	}

	if res := db.Where("merchant_id = ?", merchantID).Find(&menu); res.Error != nil {
		data := map[string]interface{}{
			"message": "Merchant not found",
		}
		return c.JSON(http.StatusInternalServerError, data)
	}

	if res := db.Where("id = ?", id).First(&menu); res.Error != nil {
		data := map[string]interface{}{
			"message": res.Error.Error(),
		}

		return c.JSON(http.StatusNotFound, data)
	}

	if res := db.Delete(&menu); res.Error != nil {
		data := map[string]interface{}{
			"message": res.Error.Error(),
		}

		return c.JSON(http.StatusInternalServerError, data)
	}

	data := map[string]interface{}{
		"message": "Food item deleted successfully",
	}

	return c.JSON(http.StatusOK, data)
}

// Food Order Routes

func OrderFood(c echo.Context) error {
	db := config.DB()
	id := c.Param("id")

	var food models.Menu
	var payload models.Order

	if res := db.Where("id = ?", id).First(&food); res.Error != nil {
		data := map[string]interface{}{
			"message": res.Error.Error(),
		}

		return c.JSON(http.StatusNotFound, data)
	}
	if err := c.Bind(&payload); err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}

		return c.JSON(http.StatusInternalServerError, data)
	}

	var cost = food.Price * float64(payload.Quantity)

	neworder := &models.Order{

		MenuID:    food.ID,
		Quantity:  payload.Quantity,
		TotalCost: cost,
	}

	if err := db.Create(&neworder).Error; err != nil {

		data := map[string]interface{}{
			"message": err.Error(),
		}

		return c.JSON(http.StatusInternalServerError, data)
	}

	return c.JSON(http.StatusOK, neworder)
}



func DisplayMenu(c echo.Context) error {
	db := config.DB()
	
var reqBody struct {
		MerchantShortcode int64 `json:"merchantshortcode"`
	}

	if err := c.Bind(&reqBody); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}


	var food []models.Menu

	
	

	if res := db.Where("merchant_short_code = ?", reqBody.MerchantShortcode).Find(&food); res.Error != nil {
		data := map[string]interface{}{
			"message": res.Error.Error(),
		}

		return c.JSON(http.StatusNotFound, data)
	}

	return c.JSON(http.StatusOK, food)
}