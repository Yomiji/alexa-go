package alexa

// constants
type IntentName string
func (b IntentName) String() string {
	return string(b)
}
// built in intents
const (
	//HelpIntent is the Alexa built-in Help Intent
	HelpIntent IntentName = "AMAZON.HelpIntent"

	//CancelIntent is the Alexa built-in Cancel Intent
	CancelIntent IntentName = "AMAZON.CancelIntent"

	//StopIntent is the Alexa built-in Stop Intent
	StopIntent IntentName = "AMAZON.StopIntent"
)

type RequestType string
func (b RequestType) String() string {
	return string(b)
}
const (
	LaunchRequest           RequestType = "LaunchRequest"
	CanFulfillIntentRequest RequestType = "CanFulfillIntentRequest"
	IntentRequest           RequestType = "IntentRequest"
	SessionEndedRequest     RequestType = "SessionEndedRequest"
)

// locales
type Locale string
func (b Locale) String() string {
	return string(b)
}
const (
	// LocaleItalian is the locale for Italian
	LocaleItalian Locale = "it-IT"

	// LocaleGerman is the locale for standard dialect German
	LocaleGerman Locale = "de-DE"

	// LocaleAustralianEnglish is the locale for Australian English
	LocaleAustralianEnglish Locale = "en-AU"

	//LocaleCanadianEnglish is the locale for Canadian English
	LocaleCanadianEnglish Locale = "en-CA"

	//LocaleBritishEnglish is the locale for UK English
	LocaleBritishEnglish Locale = "en-GB"

	//LocaleIndianEnglish is the locale for Indian English
	LocaleIndianEnglish Locale = "en-IN"

	//LocaleAmericanEnglish is the locale for American English
	LocaleAmericanEnglish Locale = "en-US"

	// LocaleJapanese is the locale for Japanese
	LocaleJapanese = "ja-JP"
)

func IsEnglish(locale Locale) bool {
	switch locale {
	case LocaleAmericanEnglish:
		return true
	case LocaleIndianEnglish:
		return true
	case LocaleBritishEnglish:
		return true
	case LocaleCanadianEnglish:
		return true
	case LocaleAustralianEnglish:
		return true
	default:
		return false
	}
}

// request

// Request is an Alexa skill request
// see https://developer.amazon.com/docs/custom-skills/request-and-response-json-reference.html#request-format
type Request struct {
	Version string  `json:"version"`
	Session Session `json:"session"`
	Body    ReqBody `json:"request"`
	Context Context `json:"context"`
}

type User struct {
	UserID      string `json:"userId"`
	AccessToken string `json:"accessToken,omitempty"`
}

// Session represents the Alexa skill session
type Session struct {
	New         bool                   `json:"new"`
	SessionID   string                 `json:"sessionId"`
	Application Application            `json:"application"`
	Attributes  map[string]interface{} `json:"attributes"`
	User        User                   `json:"user"`
}

// Video viewport shape constants
type Shape string

const (
	ROUND     Shape = "ROUND"
	RECTANGLE Shape = "RECTANGLE"
)
func (b Shape) String() string {
	return string(b)
}
// New: Video support
type Viewport struct {
	Experiences []struct {
		ArcMinuteWidth  int  `json:"arcMinuteWidth,omitempty"`
		ArcMinuteHeight int  `json:"arcMinuteHeight,omitempty"`
		CanRotate       bool `json:"canRotate"`
		CanResize       bool `json:"canResize"`
	} `json:"experiences"`
	Shape              Shape    `json:"shape"`
	PixelWidth         int      `json:"pixelWidth"`
	PixelHeight        int      `json:"pixelHeight"`
	CurrentPixelWidth  int      `json:"currentPixelWidth"`
	CurrentPixelHeight int      `json:"currentPixelHeight"`
	Dpi                int      `json:"dpi"`
	Touch              []string `json:"touch"`
	Keyboard           []string `json:"keyboard,omitempty"`
	Video              struct {
		Codecs []string `json:"codecs,omitempty"`
	} `json:"video,omitempty"`
}

type Device struct {
	DeviceID string `json:"deviceId,omitempty"`
}

type Application struct {
	ApplicationID string `json:"applicationId,omitempty"`
}

type System struct {
	APIAccessToken string      `json:"apiAccessToken"`
	Device         Device      `json:"device,omitempty"`
	Application    Application `json:"application,omitempty"`
	APIEndpoint    string      `json:"apiEndpoint"`
}

// Context represents the Alexa skill request context
type Context struct {
	System   System   `json:"System,omitempty"`
	Viewport Viewport `json:"Viewport,omitempty"`
}

// ReqBody is the actual request information
type ReqBody struct {
	Type        RequestType `json:"type"`
	RequestID   string `json:"requestId"`
	Timestamp   string `json:"timestamp"`
	Locale      Locale `json:"locale"`
	Intent      Intent `json:"intent,omitempty"`
	Reason      string `json:"reason,omitempty"`
	DialogState string `json:"dialogState,omitempty"`
}

// Intent is the Alexa skill intent
type Intent struct {
	Name  IntentName      `json:"name"`
	Slots map[string]Slot `json:"slots"`
}

// Slot is an Alexa skill slot
type Slot struct {
	Name        string      `json:"name"`
	Value       string      `json:"value"`
	Resolutions Resolutions `json:"resolutions"`
}

type ResolutionPerAuthority []struct {
	Values []struct {
		Value struct {
			Name string `json:"name"`
			Id   string `json:"id"`
		} `json:"value"`
	} `json:"values"`
}

type Resolutions struct {
	ResolutionPerAuthority ResolutionPerAuthority `json:"resolutionsPerAuthority"`
}
