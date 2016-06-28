package merchant

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"golang.org/x/net/context"

	"github.com/dchest/uniuri"
	"github.com/gin-gonic/gin"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/urlfetch"
)

// InstallmentKey is entity name for Installment model
var InstallmentKey = "Installment"

// DiscountKey is entity name for Discount model
var DiscountKey = "Discount"

// TokenKey is entity name for AuthenticatedModel
var TokenKey = "Token"

// CardsKey is entity name for Cards
var CardsKey = "Cards"

// GetVersion handler will return current version of the app.
func GetVersion(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"version":       Version,
		"merchant_name": MerchantName,
	})
}

// Charge handler will do the charging by adding Server Key into header
func Charge(c *gin.Context) {
	// Encode server key using base 64 string
	authorization := base64.StdEncoding.EncodeToString([]byte(Token))

	// HTTP client is using app engine
	appEngine := appengine.NewContext(c.Request)
	client := urlfetch.Client(appEngine)

	request, err := http.NewRequest("POST", VTBaseURL+"/charge", c.Request.Body)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status_code": http.StatusBadRequest, "status_message": err.Error()})
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

// GetPromotions will get list of available promos
func GetPromotions(c *gin.Context) {
	appEngine := appengine.NewContext(c.Request)
	discountQuery := datastore.NewQuery(DiscountKey).Ancestor(SandboxPromotionsKey(appEngine, DiscountKey))
	var discountList []Discount
	installmentQuery := datastore.NewQuery(InstallmentKey).Ancestor(SandboxPromotionsKey(appEngine, InstallmentKey))
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

// InsertDiscount is a function to save discount into datastore
func InsertDiscount(c *gin.Context) {
	requestBody, _ := ioutil.ReadAll(c.Request.Body)
	var requestObj Discount
	json.Unmarshal(requestBody, &requestObj)

	requestPoint := &requestObj

	appEngine := appengine.NewContext(c.Request)
	key := datastore.NewIncompleteKey(appEngine, DiscountKey, SandboxPromotionsKey(appEngine, DiscountKey))

	if _, err := datastore.Put(appEngine, key, requestPoint); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status_code":    http.StatusBadRequest,
			"status_message": err.Error(),
		})
	} else {
		c.JSON(http.StatusCreated, gin.H{
			"status_code":    http.StatusCreated,
			"status_message": "Discount created.",
		})
	}

}

// InsertInstallment is a function to save discount into datastore
func InsertInstallment(c *gin.Context) {
	requestBody, _ := ioutil.ReadAll(c.Request.Body)
	var requestObj Installment
	json.Unmarshal(requestBody, &requestObj)

	requestPoint := &requestObj

	appEngine := appengine.NewContext(c.Request)
	key := datastore.NewIncompleteKey(appEngine, InstallmentKey, SandboxPromotionsKey(appEngine, InstallmentKey))
	if _, err := datastore.Put(appEngine, key, requestPoint); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status_code":    http.StatusBadRequest,
			"status_message": err.Error(),
		})
	} else {
		c.JSON(http.StatusCreated, gin.H{
			"status_code":    http.StatusCreated,
			"status_message": "Installment created.",
		})
	}
}

// SandboxPromotionsKey will create Sandbox Key for Promotion
func SandboxPromotionsKey(c context.Context, entity string) *datastore.Key {
	return datastore.NewKey(c, entity, "sandbox", 0, nil)
}

// GenerateAuth will generate custom authentication token
func GenerateAuth(c *gin.Context) {
	appEngine := appengine.NewContext(c.Request)
	randomString := uniuri.NewLen(32)

	requestObj := AuthenticatedModel{Token: randomString, Cards: []Card{}}
	key := datastore.NewIncompleteKey(appEngine, TokenKey, SandboxPromotionsKey(appEngine, TokenKey))

	if _, err := datastore.Put(appEngine, key, &requestObj); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"status_code": http.StatusBadGateway, "status_message": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"X-Auth": randomString})
	}
}

// GetCards will return card list saved by specific token
func GetCards(c *gin.Context) {
	if c.Request.Header.Get("x-auth") != "" {
		appEngine := appengine.NewContext(c.Request)
		tokens := GetTokenList(appEngine)
		if tokens == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"status_code": http.StatusUnauthorized, "status_message": "Authentication Token is invalid."})
		} else {
			if CheckTokenValidity(tokens, c.Request.Header.Get("x-auth")) {
				cardQuery := datastore.NewQuery(CardsKey).Ancestor(SandboxPromotionsKey(appEngine, CardsKey))
				var cards []Card
				cardQuery.GetAll(appEngine, &cards)
				if cards != nil {
					c.JSON(http.StatusOK, gin.H{"status_code": http.StatusOK, "status_message": "Success", "data": cards})
				} else {
					c.JSON(http.StatusOK, gin.H{"status_code": http.StatusOK, "status_message": "Success", "data": []Card{}})
				}
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"status_code": http.StatusUnauthorized, "status_message": "Authentication Token is invalid."})
			}
		}
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"status_code": http.StatusUnauthorized, "status_message": "Authentication Token is invalid."})
	}
}

// RegisterCard will save a card
func RegisterCard(c *gin.Context) {
	appEngine := appengine.NewContext(c.Request)
	requestBody, _ := ioutil.ReadAll(c.Request.Body)
	var requestObj CardRequest
	json.Unmarshal(requestBody, &requestObj)
	token := c.Request.Header.Get("x-auth")

	if requestObj.StatusCode == "200" && requestObj.MaskedCard != "" && requestObj.SavedTokenID != "" {
		tokenList := GetTokenList(appEngine)
		if CheckTokenValidity(tokenList, token) {
			cardList := GetCardList(appEngine)
			card := Card{SavedTokenID: requestObj.SavedTokenID, MaskedCard: requestObj.MaskedCard}
			if !CheckCard(cardList, card) {
				// Save the card
				key := datastore.NewIncompleteKey(appEngine, CardsKey, SandboxPromotionsKey(appEngine, CardsKey))

				if _, err := datastore.Put(appEngine, key, &card); err != nil {
					c.JSON(http.StatusBadGateway, gin.H{"status_code": http.StatusBadGateway, "status_message": err.Error()})
				} else {
					c.JSON(http.StatusCreated, gin.H{"status_code": http.StatusCreated, "status_message": "Card is saved."})
				}
			} else {
				c.JSON(http.StatusConflict, gin.H{"status_code": http.StatusConflict, "status_message": "The card with same token ID is already present."})
			}
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"status_code": http.StatusUnauthorized, "status_message": "Authentication token is not valid."})
		}
	} else {
		if token != "" {
			c.JSON(http.StatusBadRequest, gin.H{"status_code": http.StatusBadRequest, "status_message": "Status code from PAPI must be 200. Masked card and saved token ID cannot be null."})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"status_code": http.StatusUnauthorized, "status_message": "Authentication token is not valid."})
		}
	}
}
