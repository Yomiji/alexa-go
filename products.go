package alexa

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const InSkillProductListRequestAPI string = "/v1/users/~current/skills/~current/inSkillProducts"

type ProductType string

const (
	SUBSCRIPTION ProductType = "SUBSCRIPTION"
	ENTITLEMENT  ProductType = "ENTITLEMENT"
	CONSUMABLE   ProductType = "CONSUMABLE"
)

type PurchaseState string

const (
	PURCHASABLE     PurchaseState = "PURCHASABLE"
	NOT_PURCHASABLE PurchaseState = "NOT_PURCHASABLE"
)

type EntitlementState string

const (
	ENTITLED     EntitlementState = "ENTITLED"
	NOT_ENTITLED EntitlementState = "NOT_ENTITLED"
)

type EntitlementReason string

const (
	AUTO_ENTITLED EntitlementReason = "AUTO_ENTITLED"
)

type PurchaseMode string

const (
	TEST PurchaseMode = "TEST"
	LIVE PurchaseMode = "LIVE"
)

type InSkillProduct struct {
	ProductId              string            `json:"product_id"`
	ReferenceName          string            `json:"reference_name"`
	Type                   ProductType       `json:"type"`
	Name                   string            `json:"name"`
	Summary                string            `json:"summary"`
	Entitled               EntitlementState  `json:"entitled"`
	Purchasable            PurchaseState     `json:"purchasable"`
	EntitlementReason      EntitlementReason `json:"entitlementReason,omitempty"`
	ActiveEntitlementCount int               `json:"activeEntitlementCount"`
	PurchaseMode           PurchaseMode      `json:"purchaseMode"`
}

type InSkillProductResponse struct {
	InSkillProducts []InSkillProduct `json:"inSkillProductsResponse"`
	IsTruncated     string           `json:"isTruncated"`
	NextToken       string           `json:"nextToken"`
}
func checkErr(err error) {
	if err != nil {
		//TODO: logging mechanism here
		panic(err)
	}
}

// Gets In-Skill Products for the user, based on Alexa Skill API
// see https://developer.amazon.com/docs/in-skill-purchase/in-skill-product-service.html
func GetInSkillProducts(request Request) (products []InSkillProduct, err error) {
	defer func() {
		if v := recover(); v != nil {
			products = nil
			if e,ok := v.(error); ok {
				err = e
			}
		}
	}()

	// establish http client
	client := &http.Client{}

	// get api host
	apiHost := request.Context.System.APIEndpoint

	// begin building request to ISP API
	getRequest, err := http.NewRequest(http.MethodGet, apiHost, nil)
	checkErr(err)

	// establish required headers for ISP api
	getRequest.Header.Add("Accept-Language", string(request.Body.Locale))
	getRequest.Header.Add("Authorization", "Bearer " + string(request.Context.System.APIAccessToken))

	resp, err := client.Do(getRequest)
	checkErr(err)

	// defer the close, ensuring the panic happens to recover later
	defer func() {
		err := resp.Body.Close()
		checkErr(err)
	}()

	// get the body bytes
	body, err := ioutil.ReadAll(resp.Body)
	checkErr(err)

	// convert to a product list
	err = json.Unmarshal(body, products)
	checkErr(err)

	// return the product list
	return products, nil
}
