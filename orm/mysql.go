package orm

import (
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"gorm.io/driver/mysql"

	"github.com/victoryang/kubernetes-cicd/config"
)

// MySQL orm
var MySQL *gorm.DB

// TODO: add mysql table index
func InitMysqlModule(c *config.DatabaseConfig) {
	var err error
	dsn := fmt.Sprintf("%v:%v@tcp(127.0.0.1:3306)/rolling?parseTime=true&charset=utf8&loc=Local", c.Username, c.Password)

	config := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: "sce_",
			SingularTable: false,
		},
	}
	MySQL, err = gorm.Open(mysql.Open(dsn), config)
	if err != nil {
		panic(fmt.Sprintf("failed to connect database: %s", err))
	}
}
