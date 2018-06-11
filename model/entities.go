package model

import "time"

// ActivityInfo store activity information
type ActivityInfo struct {
	ID              int    `xorm:"pk autoincr 'id'"`
	Name            string `xorm:"varchar(30) notnull"`
	StartTime       *time.Time
	EndTime         *time.Time
	Campus          int
	Location        string `xorm:"varchar(100)"`
	EnrollCondition string `xorm:"varchar(50)"`
	Sponsor         string `xorm:"varchar(50)"`
	Type            int
	PubStartTime    *time.Time
	PubEndTime      *time.Time
	Detail          string `xorm:"varchar(150)" `
	EnrollWay       string
	EnrollEndTime   *time.Time
	Reward          string `xorm:"varchar(30)"`
	Introduction    string `xorm:"varchar(50)"`
	Requirement     string `xorm:"varchar(50)"`
	Poster          string `xorm:"varchar(64)"`
	Qrcode          string `xorm:"varchar(64)"`
	CanEnrolled     int
	Email           string `xorm:"varchar(255)"`
	Verified        int    `xorm:"int 'verified'"`
	PCUserID        int    `xorm:"int 'pcuser_id'"`
}

// PCUser stores pc user message
type PCUser struct {
	ID           int    `xorm:"pk autoincr 'id'"`
	Name         string `xorm:"varchar(45) notnull"`
	Email        string `xorm:"varchar(255) notnull"`
	Logo         string `xorm:"varchar(70) notnull"`
	Evidence     string `xorm:"varchar(70) notnull"`
	Info         string `xorm:"varchar(150)"`
	Verified     int    `xorm:"int notnull"`
	Account      string `xorm:"varchar(64)"`
	Password     string `xorm:"varchar(64)"`
	RegisterTime *time.Time
}

// Message struct
type Message struct {
	ID      int    `xorm:"pk autoincr 'id'"`
	Subject string `xorm:"varchar(60) notnull"`
	Body    string `xorm:"varchar(150) notnull"`
	PubTime *time.Time
}

// Message_Pcuser struct
type MessagePCUser struct {
	PCUserId  int `xorm:"int 'pcuser_id'"`
	MessageId int `xorm:"int 'message_id'"`
}

type ActApplyInfo struct {
	ID        int    `xorm:"pk autoincr 'id'"`
	ActId     int    `xorm:"int notnull pk 'actid'"`
	UserId    string `xorm:"varchar(64) notnull pk 'userid'"`
	UserName  string `xorm:"varchar(64) username"`
	StudentId string `xorm:"varchar(64) studentid"`
	Phone     string `xorm:"varchar(20)"`
	School    string `xorm:"varchar(100)"`
}

// SortablePCUserList implement sort interface
type SortablePCUserList []PCUser

// TableName defines table name
func (u ActivityInfo) TableName() string {
	return "activity"
}

func (u PCUser) TableName() string {
	return "pcuser"
}

func (u Message) TableName() string {
	return "message"
}

func (u MessagePCUser) TableName() string {
	return "message_pcuser"
}

func (u ActApplyInfo) TableName() string {
	return "actapply"
}

// Sort interface for list
func (s SortablePCUserList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s SortablePCUserList) Len() int {
	return len(s)
}

func (s SortablePCUserList) Less(i, j int) bool {
	if s[i].RegisterTime != nil && s[j].RegisterTime != nil {
		return s[i].RegisterTime.Before(*s[j].RegisterTime)
	}
	return s[i].ID <= s[j].ID
}
