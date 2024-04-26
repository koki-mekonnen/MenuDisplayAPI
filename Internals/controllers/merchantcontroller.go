package controllers

import (
	"errors"
	"foodorderapi/internals/config"
	"foodorderapi/internals/models"

	"foodorderapi/utils"
	"net/http"

	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func Signin(c echo.Context) error {

	db := config.DB()
	var payload *models.Merchantsignin
	var general *models.General
	var merchant *models.Merchant

	if err := c.Bind(&payload); err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}

		return c.JSON(http.StatusInternalServerError, data)
	}

	if res := db.Where("merchant_shortcode = ?", payload.MerchantShortcode).Find(&merchant); res.Error != nil {
		data := map[string]interface{}{
			"message": "you are not allowed to use this system! Please Register first",
		}

		return c.JSON(http.StatusForbidden, data)
	}

	if res := db.Where("merchant_shortcode = ?", payload.MerchantShortcode).Find(&general); res.Error != nil {

		data := map[string]interface{}{
			"message": "you are not allowed to use this system! Please Register first",
		}

		return c.JSON(http.StatusForbidden, data)

	}

	if !general.IsActive {
		return c.JSON(http.StatusForbidden, "you are currently inactive Please contact the support team")
	}

	password := payload.Password

	// Verify password
	passCheck := utils.VerifyPassword(general.Password, password)
	passCheck2 := utils.VerifyPassword(merchant.Password, password)

	if !passCheck || !passCheck2 {
		return c.JSON(http.StatusConflict, "incorrect password! Please try again")
	}

	ttl := 24 * time.Hour

	access_token, err := utils.Createtoken(ttl, general.Id, general.Role, []byte(general.PrivateKey))
	if err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}
		return c.JSON(http.StatusBadRequest, data)
	}

	response := map[string]interface{}{
		"access_token": access_token,
		"user":         merchant,
	}

	return c.JSON(http.StatusOK, response)
}

func UpdateMerchant(c echo.Context) error {

	db := config.DB()
	// merchantID := c.Param("id")

	merchantID := c.Get("merchantID").(string)

	// Retrieve the existing merchant from the database
	var existingMerchant *models.Merchant
	var exitinggeneraluser *models.General

	if res := db.Where("id = ?", merchantID).Find(&existingMerchant); res.Error != nil {
		data := map[string]interface{}{
			"message": "Merchant not found",
		}
		return c.JSON(http.StatusInternalServerError, data)
	}

	if res := db.Where("merchant_id=?", merchantID).Find(&exitinggeneraluser); res.Error != nil {
		data := map[string]interface{}{
			"message": "merchant not found",
		}

		return c.JSON(http.StatusInternalServerError, data)
	}

	var payload *models.Merchant

	if err := c.Bind(&payload); err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, data)
	}

	hashedpassword, err := utils.Hashpassword(payload.Password)

	if err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}

		return c.JSON(http.StatusBadGateway, data)
	}

	// Update the merchant properties

	existingMerchant.BusinessName = payload.BusinessName
	existingMerchant.OwnerName = payload.OwnerName
	existingMerchant.ContactPerson = payload.ContactPerson
	existingMerchant.Email = payload.Email
	existingMerchant.Phonenumber = payload.Phonenumber
	existingMerchant.IsUpdated = true
	existingMerchant.Password = hashedpassword

	existingMerchant.UpdatedAt = time.Now()

	// Save the updated merchant to the database
	if err := db.Save(&existingMerchant).Error; err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, data)
	}

	// Update the general table properties
	exitinggeneraluser.Name = payload.BusinessName
	exitinggeneraluser.Email = payload.Email
	exitinggeneraluser.Phonenumber = payload.Phonenumber
	exitinggeneraluser.Password = hashedpassword
	exitinggeneraluser.UpdatedAt = time.Now()

	// Save the updated general table to the database
	if err := db.Save(&exitinggeneraluser).Error; err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, data)
	}

	merchantResponse := &models.MerchantResponse{
		Id:            existingMerchant.Id,
		BusinessName:  existingMerchant.BusinessName,
		OwnerName:     existingMerchant.OwnerName,
		ContactPerson: existingMerchant.ContactPerson,
		Email:         existingMerchant.Email,
		Phonenumber:   existingMerchant.Phonenumber,
		PublicKey:     existingMerchant.PublicKey,

		IsUpdated: existingMerchant.IsUpdated,
		CreatedAt: existingMerchant.CreatedAt,
		UpdatedAt: existingMerchant.UpdatedAt,
	}

	return c.JSON(http.StatusOK, merchantResponse)
}

// for forget password

func Forgetpassword(c echo.Context) error {
	db := config.DB()

	var payload *models.Signininputs
	var merchant *models.Merchant

	if err := c.Bind(&payload); err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}

		return c.JSON(http.StatusInternalServerError, data)
	}

	MerchantShortcode := payload.MerchantShortcode

	if MerchantShortcode == 0 {
		data := map[string]interface{}{
			"message": "Please fill out your phone number.",
		}
		return c.JSON(http.StatusUnauthorized, data)
	}

	

	if err := db.Where("merchant_shortcode = ?", MerchantShortcode).First(&merchant).Error; err != nil {
		// Handle query execution error
		if errors.Is(err, gorm.ErrRecordNotFound) {

			return c.JSON(http.StatusUnauthorized, "You are not allowed to use this system! Please register first.")
		}

		// Handle other query errors
		data := map[string]interface{}{
			"message": "An error occurred while querying the database.",
		}
		return c.JSON(http.StatusInternalServerError, data)
	}

	
    response := map[string]interface{}{
		
		"user":         merchant,
	}

	return c.JSON(http.StatusOK, response)

}




func GetMerchantByShortCode(c echo.Context)error{
	db:=config.DB()
	var merchant models.Merchant


	var reqBody struct {
		MerchantShortcode int64 `json:"merchantshortcode"`
	}

	if err := c.Bind(&reqBody); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if res:=db.Where("merchant_shortcode=?",reqBody.MerchantShortcode).Find(&merchant);res.Error!=nil{
		data := map[string]interface{}{
			"message": "Merchant not found",
		}
		return c.JSON(http.StatusInternalServerError, data)

	}


	return c.JSON(http.StatusOK,merchant)
}



