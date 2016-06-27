package merchant

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Token is server key
var Token = "VT-server-MFoCwh-MkpoSlOMdqiWCx9WB"

// ProductionToken is production server key
var ProductionToken = "VT-server-nt-LofCmn8UhgUcK8T8Za2ep"

// VTBaseURL is base URL of VT PAPI
var VTBaseURL = "https://api.sandbox.veritrans.co.id/v2"

// VTBaseURLProduction is base URL of production VT PAPI
var VTBaseURLProduction = "https://api.veritrans.co.id/v2"

// Version number of the app
var Version = "0.2.0"

// MerchantName is name of the merchant
var MerchantName = "rakawm"

func init() {

	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	router := gin.Default()

	router.GET("/version", GetVersion)

	sandboxRouter := router.Group("/api")
	{
		sandboxRouter.POST("/charge", Charge)
		sandboxRouter.GET("/promotions", GetPromotions)
		writeSandboxRouter := sandboxRouter.Group("/promotions").Use(CheckHeaders())
		{
			writeSandboxRouter.POST("/discount", InsertDiscount)
			writeSandboxRouter.POST("/installment", InsertInstallment)
		}
		sandboxRouter.POST("/auth", GenerateAuth)
		sandboxRouter.GET("/card", GetCards)
		sandboxRouter.POST("/card/register", RegisterCard)
	}

	productionRouter := router.Group("/api-prod")
	{
		productionRouter.POST("/charge", ChargeProduction)
		productionRouter.GET("/promotions", GetPromotionsProduction)
		writeProductionRouter := productionRouter.Group("/promotions").Use(CheckHeaders())
		{
			writeProductionRouter.POST("/discount", InsertDiscountProduction)
			writeProductionRouter.POST("/installment", InsertInstallmentProduction)
		}
		productionRouter.POST("/auth", GenerateAuthProduction)
		productionRouter.GET("/card", GetCardsProduction)
		productionRouter.POST("/card/register", RegisterCardProduction)
	}

	http.Handle("/", router)
}
