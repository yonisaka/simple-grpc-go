package config

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
)

var DB *gorm.DB

func InitializeConnDB() {
	dbDriver := viper.GetString(`database.driver`)
	dbHost := viper.GetString(`database.host`)
	dbPort := viper.GetString(`database.port`)
	dbUser := viper.GetString(`database.user`)
	dbPass := viper.GetString(`database.pass`)
	dbName := viper.GetString(`database.name`)

	var err error
	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", dbUser, dbPass, dbHost, dbPort, dbName)
	DB, err = gorm.Open(dbDriver, connection)

	if err != nil {
		log.Println("Cannot connect to database", dbDriver)
		log.Fatal("connection error:", err)
	} else {
		log.Println("We are connected to the database", dbDriver)
	}
}