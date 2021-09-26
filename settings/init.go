package settings

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var DataSettings Settings

type Settings struct {
	Port                string
	ProductPhotosFolder string `yaml:"product_photos_folder"`
	SecretKey           string `yaml:"secret_key"`
	JwtExpiredTime      string `yaml:"jwt_expired_time"`
	DB                  MySql  `yaml:"primary_db"`
}

var LoginExpirationDuration = time.Duration(1) * time.Hour
var JwtSigningMethod = jwt.SigningMethodHS256
var JwtSignatureKey = []byte("")
var StaticFolder = "static/"
var ProductPhotosPrefixUrl = "/product-photos/"

func init() {

	// load file config
	file, err := os.Open("./config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	// close file config
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Println(err)
		}
	}(file)

	// load config to config variable
	configDecoded := yaml.NewDecoder(file)
	err = configDecoded.Decode(&DataSettings)

	if DataSettings.ProductPhotosFolder != "" {
		if !strings.Contains(DataSettings.ProductPhotosFolder, "/") {
			DataSettings.ProductPhotosFolder = DataSettings.ProductPhotosFolder + "/"
		}
		DataSettings.ProductPhotosFolder = StaticFolder + DataSettings.ProductPhotosFolder
		_, err = os.Stat(DataSettings.ProductPhotosFolder)
		if os.IsNotExist(err) {
			err = os.Mkdir(DataSettings.ProductPhotosFolder, 0755)
			if err != nil {
				log.Fatal(fmt.Sprintf("cannot create folder %s.", DataSettings.ProductPhotosFolder))
			}
		}
	} else {
		log.Fatal("you have to set `product_photos_folder` value in config.yaml ")
	}

	if DataSettings.JwtExpiredTime != "" {
		expiredTime, err := strconv.Atoi(DataSettings.JwtExpiredTime)
		if err != nil {
			log.Fatal("error to load 'jwt_expired_time' config.")
		}
		LoginExpirationDuration = time.Duration(expiredTime) * time.Hour
	}

	JwtSignatureKey = []byte(DataSettings.SecretKey)

	if err != nil {
		log.Fatal("File config is not valid")
	}
}
