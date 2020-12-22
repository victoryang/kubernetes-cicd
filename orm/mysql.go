package orm

import (
	"fmt"
	
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // mysql driver

	"github.com/victoryang/kubernetes-cicd/config"
)

// MySQL orm
var MySQL *gorm.DB

// TODO: add mysql table index
func InitMysqlModule(c *config.DatabaseConfig) {
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return "sce_" + defaultTableName
	}

	connStr := fmt.Sprintf("%v:%v@tcp(127.0.0.1:3306/rolling?parseTime=true&charset=utf8&loc=Local)", c.Username, c.Password)
	var err error
	MySQL, err = gorm.Open("mysql", connStr)
	if err != nil {
		panic(fmt.Sprintf("failed to connect database: %s", err))
	}
}
