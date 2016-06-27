package merchant

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

// CheckTokenValidity check validity of authentication token
func CheckTokenValidity(tokenList []AuthenticatedModel, token string) bool {
	var isValid = false
	for _, a := range tokenList {
		if a.Token == token {
			return true
		}
	}
	return isValid
}

// GetTokenList will return token list
func GetTokenList(context context.Context) []AuthenticatedModel {
	tokenQuery := datastore.NewQuery(TokenKey).Ancestor(SandboxPromotionsKey(context, TokenKey))
	var tokens []AuthenticatedModel
	tokenQuery.GetAll(context, &tokens)
	return tokens
}

// GetTokenListProduction will return token list
func GetTokenListProduction(context context.Context) []AuthenticatedModel {
	tokenQuery := datastore.NewQuery(TokenKey).Ancestor(ProductionPromotionsKey(context, TokenKey))
	var tokens []AuthenticatedModel
	tokenQuery.GetAll(context, &tokens)
	return tokens
}

// GetCardList will return card list from a token
func GetCardList(context context.Context) []Card {
	cardQuery := datastore.NewQuery(CardsKey).Ancestor(SandboxPromotionsKey(context, CardsKey))
	var cards []Card
	cardQuery.GetAll(context, &cards)
	return cards
}

// GetCardListProduction will return card list from a token
func GetCardListProduction(context context.Context) []Card {
	cardQuery := datastore.NewQuery(CardsKey).Ancestor(ProductionPromotionsKey(context, CardsKey))
	var cards []Card
	cardQuery.GetAll(context, &cards)
	return cards
}

// CheckCard to check if card already saved
func CheckCard(cardList []Card, card Card) bool {
	var isCardAlreadySaved = false

	for _, a := range cardList {
		if a.SavedTokenID == card.SavedTokenID {
			isCardAlreadySaved = true
		}
	}

	return isCardAlreadySaved
}
