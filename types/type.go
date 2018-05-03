package types

// ActivityInfo stores json format the front-end wanted
type ActivityInfo struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	StartTime       string  `json:"startTime"`
	EndTime         string  `json:"endTime"`
	Campus          int    `json:"campus"`
	Location        string `json:"location"`
	EnrollCondition string `json:"enrollCondition"`
	Sponsor         string `json:"sponsor"`
	Type            int    `json:"type"`
	PubStartTime    string  `json:"pubStartTime"`
	PubEndTime      string  `json:"pubEndTime"`
	Detail          string `json:"detail"`
	Reward          string `json:"reward"`
	Introduction    string `json:"introduction"`
	Requirement     string `json:"requirement"`
	Poster          string `json:"poster"`
	Qrcode          string `json:"qrcode"`
	Email           string `json:"email"`
	Verified        int    `json:"verified"`
}
