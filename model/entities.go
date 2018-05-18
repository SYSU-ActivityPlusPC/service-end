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
	Reward          string `xorm:"varchar(30)"`
	Introduction    string `xorm:"varchar(50)"`
	Requirement     string `xorm:"varchar(50)"`
	Poster          string `xorm:"varchar(64)"`
	Qrcode          string `xorm:"varchar(64)"`
	Email           string `xorm:"varchar(255)"`
	Verified        int    `xorm:"int 'verified'"`
}

// PCUser stores pc user message
type PCUser struct {
	ID       int    `xorm:"pk autoincr 'id'"`
	Name     string `xorm:"varchar(45) notnull"`
	Email    string `xorm:"varchar(255) notnull"`
	Logo     string `xorm:"varchar(70) notnull"`
	Evidence string `xorm:"varchar(70) notnull"`
	Info     string `xorm:"varchar(150)"`
	Verified int    `xorm:"int notnull"`
	Account  string `xorm:"varchar(64)"`
	Password string `xorm:"varchar(64)"`
}

// TableName defines table name
func (u ActivityInfo) TableName() string {
	return "activity"
}

func (u PCUser) TableName() string {
	return "pcuser"
}
