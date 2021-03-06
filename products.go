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
	ReferenceName          string           `json:"referenceName,omitempty"`
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
	InSkillProducts []InSkillProduct `json:"inSkillProducts"`
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

var loggingEnabled = false

func ToggleDebugLogging(value bool) {
	loggingEnabled = value
}

// Assuming a slot's ID corresponds to an ISP reference name, get the ISP from the slice that
// matches the slot
func SimpleSlotToProducts(slot Slot, products []InSkillProduct) (product *InSkillProduct) {
	for _,resolution := range slot.Resolutions.ResolutionPerAuthority {
		for _,value := range resolution.Values {
			for _,product := range products {
				if value.Value.Id == product.ReferenceName {
					return &product
				}
			}
		}
	}
	return nil
}

// Gets In-Skill Products for the user, based on Alexa Skill API
// see https://developer.amazon.com/docs/in-skill-purchase/in-skill-product-service.html
func GetInSkillProducts(request Request) (products []InSkillProduct, err error) {
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
	var inSkillResponse = &InSkillProductResponse{}
	err = json.Unmarshal(body, inSkillResponse)
	checkErr(err)

	// set return value
	products = inSkillResponse.InSkillProducts

	if loggingEnabled {
		slog.Info("Retrieved %v products", len(products))
	}

	// return the product list
	return products, nil
}

// get all products purchased
func GetPurchasedProducts(products []InSkillProduct) []InSkillProduct {
	result := make([]InSkillProduct, 0, len(products))

	for _, v := range products {
		if v.Entitled == ENTITLED {
			result = append(result, v)
		}
	}

	return result
}

// get all products available to purchase
func GetAvailableProducts(products []InSkillProduct) []InSkillProduct {
	result := make([]InSkillProduct, 0, len(products))

	for _, v := range products {
		if v.Entitled == NOT_ENTITLED {
			result = append(result, v)
		}
	}

	return result
}