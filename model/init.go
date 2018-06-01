package model

import (
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
)

// Engine is mysql engine
var Engine *xorm.Engine

func init() {
	DBADDRESS := os.Getenv("DATABASE_ADDRESS")
	if len(DBADDRESS) == 0 {
		DBADDRESS = "localhost"
	}
	DBPORT := os.Getenv("DATABASE_PORT")
	if len(DBPORT) != 0 && DBPORT[0] != ':' {
		DBPORT = ":" + DBPORT
	}
	// url := fmt.Sprintf("root:root@tcp(%s%s)/activityplus?charset=utf8", DBADDRESS, DBPORT)
	// url := fmt.Sprintf("root:root@tcp(119.29.155.194:5000)/activityplus?charset=utf8")
	url := fmt.Sprintf("root:root@tcp(127.0.0.1:3306)/activityplus?charset=utf8")
	var err error
	engine, err := xorm.NewEngine("mysql", url)
	if err != nil {
		panic(err)
	}
	Engine = engine
	if os.Getenv("DEVELOP") == "TRUE" {
		Engine.Ping()
		Engine.ShowSQL(true)
		Engine.Logger().SetLevel(core.LOG_DEBUG)
	}
}
