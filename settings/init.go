package settings

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

var DataSettings Settings

type Settings struct {
	Port string
	DB   MySql `yaml:"primary_db"`
}

var LoginExpirationDuration = time.Duration(1) * time.Hour
var JwtSigningMethod = jwt.SigningMethodHS256
var JwtSignatureKey = []byte("FO00XG0P0ndas1CrOkBos55JAYUDA44e4ateByVNEU@20200218Gara2BCA")

func init() {

	// load file config
	file, err := os.Open("./config.yaml")
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	// close file config
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(file)

	// load config to config variable
	configDecoded := yaml.NewDecoder(file)
	err = configDecoded.Decode(&DataSettings)
	if err != nil {
		fmt.Println("File config is not valid")
		os.Exit(0)
	}
}
