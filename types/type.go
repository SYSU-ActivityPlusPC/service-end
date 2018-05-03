package types

// ActivityInfo stores json format the front-end wanted
type ActivityInfo struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	StartTime       int64  `json:"startTime"`
	EndTime         int64  `json:"endTime"`
	Campus          int    `json:"campus"`
	Location        string `json:"location"`
	EnrollCondition string `json:"enrollCondition"`
	Sponsor         string `json:"sponsor"`
	Type            int    `json:"type"`
	PubStartTime    int64  `json:"pubStartTime"`
	PubEndTime      int64  `json:"pubEndTime"`
	Detail          string `json:"detail"`
	Reward          string `json:"reward"`
	Introduction    string `json:"introduction"`
	Requirement     string `json:"requirement"`
	Poster          string `json:"poster"`
	Qrcode          string `json:"qrcode"`
	Email           string `json:"email"`
	Verified        int    `json:"verified"`
}
