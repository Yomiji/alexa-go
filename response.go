package alexa

//NewSimpleResponse builds a session response
func NewSimpleResponse(title string, text string) Response {
	r := Response{
		Version: "1.0",
		Body: ResBody{
			OutputSpeech: &Payload{
				Type: "PlainText",
				Text: text,
			},
			Card: &Payload{
				Type:    "Simple",
				Title:   title,
				Content: text,
			},
			ShouldEndSession: true,
		},
	}
	return r
}

// Response Types

// Response is the response back to the Alexa speech service
type Response struct {
	Version           string                 `json:"version"`
	SessionAttributes map[string]interface{} `json:"sessionAttributes,omitempty"`
	Body              ResBody                `json:"response"`
}

//TODO: Ensure that directives is passable to ISP API
// ResBody is the actual body of the response
type ResBody struct {
	OutputSpeech     *Payload         `json:"outputSpeech,omitempty"`
	Card             *Payload         `json:"card,omitempty"`
	Reprompt         *Reprompt        `json:"reprompt,omitempty"`
	Directives       []Directives     `json:"directives,omitempty"`
	ShouldEndSession bool             `json:"shouldEndSession"`
	CanFulfillIntent CanFulfillIntent `json:"canFulFillIntent,omitempty"`
}

type CanFulfillSlotDefinition struct {
	CanUnderstand string `json:"canUnderstand,omitempty"`
	CanFulfill    string `json:"canFulfill,omitempty"`
}

// Body structure of the CanFulfillIntentResponse
type CanFulfillIntent struct {
	CanFulfill string                              `json:"canFulfill,omitempty"`
	Slots      map[string]CanFulfillSlotDefinition `json:"slots,omitempty"`
}

// Reprompt is imformation
type Reprompt struct {
	OutputSpeech *Payload `json:"outputSpeech,omitempty"`
}
type ISPPayload struct {
	InSkillProduct InSkillProduct `json:"InSkillProduct,omitempty"`
	UpsellMessage  string         `json:"upsellMessage,omitempty"`
}

// Directives is imformation
type Directives struct {
	Type          string         `json:"type,omitempty"`
	Name          string         `json:"name,omitempty"`
	Payload       ISPPayload     `json:"payload,omitempty"`
	Token         string         `json:"token,omitempty"`
	SlotToElicit  string         `json:"slotToElicit,omitempty"`
	UpdatedIntent *UpdatedIntent `json:"UpdatedIntent,omitempty"`
	PlayBehavior  string         `json:"playBehavior,omitempty"`
	AudioItem     AudioItem      `json:"audioItem,omitempty"`
}
type AudioItem struct {
	Stream Stream `json:"stream,omitempty"`
}

type Stream struct {
	Token                string `json:"token,omitempty"`
	URL                  string `json:"url,omitempty"`
	OffsetInMilliseconds int    `json:"offsetInMilliseconds,omitempty"`
}

// UpdatedIntent is to update the Intent
type UpdatedIntent struct {
	Name               string                 `json:"name,omitempty"`
	ConfirmationStatus string                 `json:"confirmationStatus,omitempty"`
	Slots              map[string]interface{} `json:"slots,omitempty"`
}

// Image ...
type Image struct {
	SmallImageURL string `json:"smallImageUrl,omitempty"`
	LargeImageURL string `json:"largeImageUrl,omitempty"`
}

// Payload ...
type Payload struct {
	Type    string `json:"type,omitempty"`
	Title   string `json:"title,omitempty"`
	Text    string `json:"text,omitempty"`
	SSML    string `json:"ssml,omitempty"`
	Content string `json:"content,omitempty"`
	Image   Image  `json:"image,omitempty"`
}
