package controllers

import (
	"bitbucket/api_format/db"
	"bitbucket/api_format/domain"
	"bitbucket/api_format/models"
	"bitbucket/api_format/utils"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func Createtenant(c *fiber.Ctx) error {
	var tenant models.Tenants

	if err := c.BodyParser(&tenant); err != nil {
		return err
	}

	tenantchecking := domain.Gettenantbyname(tenant.Tenantname)
	tenant.Apptenantid = tenantchecking
	// tenant.Tenantname = tname

	fmt.Println("------tenantchecking-----", tenantchecking)
	// fmt.Println("------tenantchecking-----", tname)

	if tenantchecking == 0 {
		fmt.Println("came in")
		info1, err := domain.Createtenant(tenant)
		if info1 == 0 || err != nil {
			c.JSON(fiber.Map{
				"code":    http.StatusConflict,
				"message": "Tenant creation unsuccessful, check tenant's details!",
				"status":  false,
			})
		}
		tenant.Apptenantid = info1
	} else {
		return c.JSON(fiber.Map{
			"code":    http.StatusConflict,
			"message": "Tenant creation unsuccessful, tenant already exists!",
			"status":  false,
		})
	}
	fmt.Println("-----tenantid----:::", tenant.Tenantid)
	result := domain.Gettenantbyid(tenant.Apptenantid)
	return c.JSON(fiber.Map{
		"code":    http.StatusCreated,
		"message": "Success",
		"status":  true,
		"details": result,
	})
}

func Updatetenant(c *fiber.Ctx) error {
	// Parse request body into tenant struct
	var tenant models.Tenants
	if err := c.BodyParser(&tenant); err != nil {
		return err
	}

	// Fetch existing data from the database
	existingtenant := domain.Gettenantbyid(tenant.Apptenantid)

	// Check if the tenant exists
	if existingtenant.Apptenantid == 0 {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"code":    http.StatusNotFound,
			"message": "tenant not found!",
			"status":  false,
		})
	}
	// Update only the fields provided
	result := db.DB.Table("tenants").Where("apptenantid = ?", tenant.Apptenantid).Updates(&tenant)
	if result.Error != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"code":    http.StatusInternalServerError,
			"message": "Failed to update tenant",
			"status":  false,
		})
	}
	// Return the updated data
	return c.JSON(fiber.Map{
		"code":    http.StatusOK,
		"message": "tenant updated successfully",
		"status":  true,
	})
}

func Gettenantapps(c *fiber.Ctx) error {
	tid, _ := strconv.Atoi(c.Query("apptenantid"))

	res1 := domain.Gettenantidbytid(tid)

	res2 := domain.Gettenantappsbytid(res1)

	return c.JSON(fiber.Map{
		"code":    http.StatusCreated,
		"message": "Success",
		"status":  true,
		"details": res2,
	})
}

func Getavailableapps(c *fiber.Ctx) error {
	tid, _ := strconv.Atoi(c.Query("apptenantid"))

	var tenants []utils.App_application = domain.Getavailablepermissionsbyid(tid)

	return c.JSON(fiber.Map{
		"code":    http.StatusCreated,
		"message": "Success",
		"status":  true,
		"details": tenants,
	})
}

func Getallbytenant(c *fiber.Ctx) error {

	appid := c.Query("orgname")

	var data []models.Tenants

	q1 := `SELECT * FROM tenants where apptenantname = '` + appid + `'`

	println(q1)

	db.DB.Raw(q1).Find(&data)

	return c.JSON(fiber.Map{
		"code":    http.StatusOK,
		"message": "Success",
		"status":  true,
		"details": data,
	})
}

type Orgresponse struct {
	Apptenantid   int    `json:"apptenantid"`
	Apptenantname string `json:"apptenantname"`
}

func Getallorganizations(c *fiber.Ctx) error {

	var data []Orgresponse

	q1 := `SELECT MIN(apptenantid) AS apptenantid, apptenantname FROM tenants GROUP BY apptenantname`

	println(q1)

	db.DB.Raw(q1).Find(&data)

	return c.JSON(fiber.Map{
		"code":    http.StatusOK,
		"message": "Success",
		"status":  true,
		"details": data,
	})
}

func Gettenantbyid(c *fiber.Ctx) error {

	tid, _ := strconv.Atoi(c.Query("apptenantid"))

	var tenants models.Tenants = domain.Gettenantbyid(tid)

	return c.JSON(fiber.Map{
		"code":    http.StatusCreated,
		"message": "Success",
		"status":  true,
		"details": tenants,
	})
}

func Deletetenantbyid(c *fiber.Ctx) error {

	tid, _ := strconv.Atoi(c.Query("apptenantid"))

	res1, res2 := domain.Deletetenantbyid(tid)
	if res1 != false {
		return c.JSON(fiber.Map{
			"code":    http.StatusCreated,
			"message": "Successfully deleted.",
			"status":  true,
		})
	} else {
		return c.JSON(fiber.Map{
			"code":    http.StatusBadRequest,
			"message": res2,
			"status":  false,
		})
	}
}

func Getalltenants(c *fiber.Ctx) error {

	var user []models.Tenants = domain.Getalltenants()

	return c.JSON(fiber.Map{
		"code":    http.StatusCreated,
		"message": "Success",
		"status":  true,
		"details": user,
	})
}
