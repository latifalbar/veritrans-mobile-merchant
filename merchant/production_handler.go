package merchant

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/urlfetch"
)

// ChargeProduction handler will do the charging by adding Production Server Key into header
func ChargeProduction(c *gin.Context) {
	// Encode server key using base 64 string
	authorization := base64.StdEncoding.EncodeToString([]byte(ProductionToken))

	// HTTP client is using app engine
	appEngine := appengine.NewContext(c.Request)
	client := urlfetch.Client(appEngine)

	request, err := http.NewRequest("POST", VTBaseURLProduction+"/charge", c.Request.Body)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status_code": "400", "status_message": "Bad Request"})
	} else {
		request.Header.Add("Content-Type", "application/json")
		request.Header.Add("Accept", "application/json")
		request.Header.Add("Authorization", "Basic "+authorization)
		response, _ := client.Do(request)

		responseBody, _ := ioutil.ReadAll(response.Body)
		var respObj interface{}
		json.Unmarshal(responseBody, &respObj)
		c.JSON(http.StatusOK, respObj)
	}
}

// GetPromotionsProduction will get list of available promos
func GetPromotionsProduction(c *gin.Context) {
	appEngine := appengine.NewContext(c.Request)
	discountQuery := datastore.NewQuery(DiscountKey).Ancestor(ProductionPromotionsKey(appEngine, DiscountKey))
	var discountList []Discount
	installmentQuery := datastore.NewQuery(InstallmentKey).Ancestor(ProductionPromotionsKey(appEngine, InstallmentKey))
	var installmentList []Installment

	discountQuery.GetAll(appEngine, &discountList)

	installmentQuery.GetAll(appEngine, &installmentList)

	var discountListFinal []Discount
	var installmentListFinal []Installment

	if discountList == nil {
		discountListFinal = []Discount{}
	} else {
		discountListFinal = discountList
	}
	if installmentList == nil {
		installmentListFinal = []Installment{}
	} else {
		installmentListFinal = installmentList
	}
	var promotion = Promotion{InstallmentList: installmentListFinal, DiscountList: discountListFinal}

	c.JSON(http.StatusOK, gin.H{
		"status_code":    http.StatusOK,
		"status_message": "success",
		"data":           promotion,
	})
}

// InsertDiscountProduction is a function to save discount into datastore
func InsertDiscountProduction(c *gin.Context) {
	requestBody, _ := ioutil.ReadAll(c.Request.Body)
	var requestObj Discount
	json.Unmarshal(requestBody, &requestObj)

	requestPoint := &requestObj

	appEngine := appengine.NewContext(c.Request)
	key := datastore.NewIncompleteKey(appEngine, DiscountKey, ProductionPromotionsKey(appEngine, DiscountKey))

	if _, err := datastore.Put(appEngine, key, requestPoint); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status_code":    http.StatusBadRequest,
			"status_message": err.Error(),
		})
	} else {
		c.JSON(http.StatusCreated, gin.H{
			"status_code":    http.StatusCreated,
			"status_message": "Discount created",
		})
	}

}

// InsertInstallmentProduction is a function to save discount into datastore
func InsertInstallmentProduction(c *gin.Context) {
	requestBody, _ := ioutil.ReadAll(c.Request.Body)
	var requestObj Installment
	json.Unmarshal(requestBody, &requestObj)

	requestPoint := &requestObj

	appEngine := appengine.NewContext(c.Request)
	key := datastore.NewIncompleteKey(appEngine, InstallmentKey, ProductionPromotionsKey(appEngine, InstallmentKey))
	if _, err := datastore.Put(appEngine, key, requestPoint); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status_code":    http.StatusBadRequest,
			"status_message": err.Error(),
		})
	} else {
		c.JSON(http.StatusCreated, gin.H{
			"status_code":    http.StatusCreated,
			"status_message": "Installment created",
		})
	}
}

// ProductionPromotionsKey will create Production Key for Promotion
func ProductionPromotionsKey(c context.Context, entity string) *datastore.Key {
	return datastore.NewKey(c, entity, "production", 0, nil)
}
