package alexa

import (
	"encoding/json"
	"fmt"
	"github.com/Yomiji/slog"
	"io/ioutil"
	"net/http"
	"time"
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
	ProductId              string           `json:"productId,omitempty"`
	ReferenceName          string           `json:"reference_name,omitempty"`
	Type                   ProductType      `json:"type,omitempty"`
	Name                   string           `json:"name,omitempty"`
	Summary                string           `json:"summary,omitempty"`
	Entitled               EntitlementState `json:"entitled,omitempty"`
	Purchasable            PurchaseState    `json:"purchasable,omitempty"`
	EntitlementReason      string           `json:"entitlementReason,omitempty"`
	ActiveEntitlementCount int              `json:"activeEntitlementCount,omitempty"`
	PurchaseMode           PurchaseMode     `json:"purchaseMode,omitempty"`
}

type InSkillProductResponse struct {
	InSkillProducts []InSkillProduct `json:"inSkillProductsResponse"`
	IsTruncated     string           `json:"isTruncated"`
	NextToken       string           `json:"nextToken"`
}

func checkErr(err error) {
	if err != nil {
		slog.Fail("error occurred: %v", err)
		panic(err)
	}
}

var ispDefaultRequestTimeout time.Duration = 30
var ISPRequestTimeout time.Duration = 30
var client = &http.Client{}

func init() {

}
// Gets In-Skill Products for the user, based on Alexa Skill API
// see https://developer.amazon.com/docs/in-skill-purchase/in-skill-product-service.html
func GetInSkillProducts(request Request, loggingEnabled bool) (products []InSkillProduct, err error) {
	defer func() {
		if v := recover(); v != nil {
			products = nil
			if e, ok := v.(error); ok {
				err = e
			}
		}
	}()

	if loggingEnabled  {
		slog.Debug("Entering GetInSkillProducts")
	}

	// set client timeout
	//noinspection GoBoolExpressions
	if ISPRequestTimeout > 0 {
		client.Timeout = ISPRequestTimeout * time.Second
	} else {
		client.Timeout = ispDefaultRequestTimeout * time.Second
	}

	if loggingEnabled {
		slog.Debug("Constructing client for ISP deployment")
	}


	// get api host
	apiHost := request.Context.System.APIEndpoint + "/v1/users/~current/skills/~current/inSkillProducts"


	if loggingEnabled {
		slog.Debug("Generating request for endpoint %s", apiHost)
	}

	// begin building request to ISP API
	getRequest, err := http.NewRequest(http.MethodGet, apiHost, http.NoBody)
	checkErr(err)

	// establish required headers for ISP api
	getRequest.Header.Set("Accept-Language", string(request.Body.Locale))
	getRequest.Header.Set("Authorization", fmt.Sprintf("Bearer %s", request.Context.System.APIAccessToken))
	if loggingEnabled {
		slog.Debug("Performing request: %v", getRequest)
	}
	resp, err := client.Do(getRequest)

	if loggingEnabled {
		slog.Debug("Request completed.")
	}
	checkErr(err)

	// defer the close, ensuring the panic happens to recover later
	defer func() {
		err := resp.Body.Close()
		checkErr(err)
	}()

	// get the body bytes
	body, err := ioutil.ReadAll(resp.Body)
	checkErr(err)

	if loggingEnabled {
		slog.Debug("Body retrieved: %v", string(body))
	}
	// convert to a product list
	err = json.Unmarshal(body, &products)
	checkErr(err)

	if loggingEnabled {
		slog.Info("Retrieved %v products", len(products))
	}

	// return the product list
	return products, nil
}
