package models

type Issue struct {
	IssueTitle        string `json:"issue_title"`
	URL               string `json:"url"`
	AdditionalContext string `json:"additional_context"`
}

type Entry struct {
	Matchers Match  `json:"matchers"`
	Category string `json:"category"`
	Tag      string `json:"tag"`
}

type Match struct {
	Headers      map[string]interface{} `json:"headers"`
	Strings      []string               `json:"strings"`
	ResponseCode int                    `json:"response_code"`
	Condition    string                 `json:"condition"`
}

type GenzaiResult struct {
	Target        string  `json:"target"`
	IoTidentified string  `json:"iot_identified"`
	Category      string  `json:"category"`
	Issues        []Issue `json:"issues"`
}


type ScanRequest struct {
	Targets []string `json:"targets"`
}

type ScanResponse struct {
	Results      []GenzaiResult `json:"results"`
	Targets      []string       `json:"targets"`
	TotalScanned int            `json:"total_scanned"`
	TimeElapsed  string         `json:"time_elapsed"`
	Errors       []string       `json:"errors,omitempty"`
}

type Response struct {
	Results []GenzaiResult `json:"Results"`
	Targets []string       `json:"Targets"`
}

type GenzaiDB map[string]Entry
type VendorDB map[string]interface{}
type VendorVulnsDB map[string]interface{}