package settings

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"strconv"
	"time"
)

var DataSettings Settings

type Settings struct {
	Port           string
	StaticFolder   string `yaml:"static_folder"`
	SecretKey      string `yaml:"secret_key"`
	JwtExpiredTime string `yaml:"jwt_expired_time"`
	DB             MySql  `yaml:"primary_db"`
}

var LoginExpirationDuration = time.Duration(1) * time.Hour
var JwtSigningMethod = jwt.SigningMethodHS256
var JwtSignatureKey = []byte("")

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

	_, err = os.Stat(DataSettings.StaticFolder)
	if os.IsNotExist(err) {
		log.Fatal(fmt.Sprintf("Folder %s does not exist.", DataSettings.StaticFolder))
	}

	if DataSettings.JwtExpiredTime != "" {
		expiredTime, err := strconv.Atoi(DataSettings.JwtExpiredTime)
		if err != nil {
			log.Fatal("error to load 'jwt_expired_time' config.")
			return
		}
		LoginExpirationDuration = time.Duration(expiredTime) * time.Hour
	}

	JwtSignatureKey = []byte(DataSettings.SecretKey)

	if err != nil {
		fmt.Println("File config is not valid")
		os.Exit(0)
	}
}
