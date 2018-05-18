package types

// StringActivityInfo stores json format the front-end wanted
type StringActivityInfo struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	StartTime       string `json:"startTime"`
	EndTime         string `json:"endTime"`
	Campus          int    `json:"campus"`
	Location        string `json:"location"`
	EnrollCondition string `json:"enrollCondition"`
	Sponsor         string `json:"sponsor"`
	Type            int    `json:"type"`
	PubStartTime    string `json:"pubStartTime"`
	PubEndTime      string `json:"pubEndTime"`
	Detail          string `json:"detail"`
	Reward          string `json:"reward"`
	Introduction    string `json:"introduction"`
	Requirement     string `json:"requirement"`
	Poster          string `json:"poster"`
	Qrcode          string `json:"qrcode"`
	Email           string `json:"email"`
}

// IntActivityInfo stores json format the front-end wanted
type IntActivityInfo struct {
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

// ActivityIntroduction include required information in activity list page
type ActivityIntroduction struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	StartTime       int64  `json:"startTime"`
	EndTime         int64  `json:"endTime"`
	Campus          int    `json:"campus"`
	EnrollCondition string `json:"enrollCondition"`
	Sponsor         string `json:"sponsor"`
	PubStartTime    int64  `json:"pubStartTime"`
	PubEndTime      int64  `json:"pubEndTime"`
	Verified        int    `json:"verified"`
}

// ActivityList defines the return format
type ActivityList struct {
	Content []ActivityIntroduction `json:"content"`
}

// PCUserInfo include required login response
type PCUserInfo struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Logo  string `json:"logo"`
	Token string `json:"token"`
}

// PCUserRequest include user login message
type PCUserRequest struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

// PCUserSignInfo include sign up message
type PCUserSignInfo struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Logo     string `json:"logo"`
	Evidence string `json:"evidence"`
	Info     string `json:"info"`
}

// FileInfo include filename of the file
type FileInfo struct {
	Filename string `json:"filename"`
}
