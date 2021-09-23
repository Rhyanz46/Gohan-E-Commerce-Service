package settings

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

type MySql struct {
	Host     string `yaml:"host"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DBName   string `yaml:"db_name"`
	Port     string `yaml:"port"`
}

func (m MySql) CreateConnection() (db *gorm.DB) {
	fmt.Println("connecting to mysql \t\t: ", m.DBName, " ("+m.Host+")")
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       "" + m.Username + ":" + m.Password + "@tcp(" + m.Host + ":" + m.Port + ")/" + m.DBName + "?charset=utf8&parseTime=True&loc=Local", // data source name
		DefaultStringSize:         256,                                                                                                                               // default size for string fields
		DisableDatetimePrecision:  true,                                                                                                                              // disable datetime precision, which not supported before MySQL 5.6
		DontSupportRenameIndex:    true,                                                                                                                              // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
		DontSupportRenameColumn:   true,                                                                                                                              // `change` when rename column, rename column not supported before MySQL 8, MariaDB
		SkipInitializeWithVersion: false,                                                                                                                             // auto configure based on currently MySQL version
	}), &gorm.Config{})
	if err != nil {
		log.Fatal("cannot \t\t: connect to mysql database : ", m.DBName, " ("+m.Host+")")
	} else {
		fmt.Println("connected to mysql \t\t: ", m.DBName, " ("+m.Host+")")
	}
	return db
}
