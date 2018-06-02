package types

// ActivityInfo stores json format the front-end wanted
type ActivityInfo struct {
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
	EnrollWay       string `json:"enrollWay"`
	EnrollEndTime   string `json:"enrollEndTime"`
	CanEnrolled     int    `json:"can_enrolled"`
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
	Type            int    `json:"type"`
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

// PCUserDetailedInfo include sign up message
type PCUserDetailedInfo struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	Logo         string `json:"logo"`
	Evidence     string `json:"evidence"`
	Info         string `json:"info"`
	Account      string `json:"account"`
	RegisterTime string `json:"registerTime"`
}

// FileInfo include filename of the file
type FileInfo struct {
	Filename string `json:"filename"`
}

// EmailContent contains email message
type EmailContent struct {
	From    string
	To      string
	Subject string
	Content string
}

// PCUserListInfo in the list
type PCUserListInfo struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Logo         string `json:"logo"`
	Verified     int    `json:"verified"`
	RegisterTime string `json:"register_time"`
}

// Message info passed from front-end
type MessageInfo struct {
	ID 		 int    `json:"id"`
	Subject  string `json:"subject"`
	Body 	 string `json:"body"`
	PCUserId []int  `json:"pcuserId"`
}

// MessageIntroduction
type MessageIntroduction struct {
	ID 		 int    `json:"id"`
	Subject  string `json:"subject"`
	Body 	 string `json:"body"`
	PubTime  string  `json:"pubTime"`
	SendTo []string  `json:"sendTo"`
}

// MessageList defines the return format
type MessageList struct {
	Content []MessageIntroduction `json:"content"`
}