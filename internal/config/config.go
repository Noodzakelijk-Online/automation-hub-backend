package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	serverPort       string = "SERVER_PORT"
	baseUrl          string = "BASE_URL"
	nginxContainer   string = "NGINX_CONTAINER"
	configDir        string = "CONFIG_DIR"
	dbHost           string = "DB_HOST"
	dbPort           string = "DB_PORT"
	dbName           string = "DB_NAME"
	dbUser           string = "DB_USER"
	dbPassword       string = "DB_PASSWORD"
	imageMaxSizeInMb string = "IMAGE_MAX_SIZE_IN_MB"
	imageExtensions  string = "IMAGE_EXTENSIONS"
	imageSaveDir     string = "IMAGE_SAVE_DIR"
)

type Configuration struct {
	ConfigDir       string
	BaseUrl         string
	ServerPort      string
	NginxContainer  string
	DbHost          string
	DbPort          int
	DbName          string
	DbUser          string
	DbPassword      string
	ImageMaxSize    int64
	ImageExtensions []string
	ImageSaveDir    string
}

var AppConfig Configuration

func Init() {
	servNumPort := getEnvInt(serverPort, 80)
	if err := validatePort(servNumPort); err != nil {
		panic(err)
	}
	dbNumPort := getEnvInt(dbPort, 5432)
	if err := validatePort(dbNumPort); err != nil {
		panic(err)
	}

	AppConfig = Configuration{
		ConfigDir:       getEnvString(configDir, "/app/sites-enabled"),
		BaseUrl:         getEnvString(baseUrl, "/api"),
		ServerPort:      ":" + strconv.Itoa(servNumPort),
		NginxContainer:  getEnvString(nginxContainer, "gateway"),
		DbHost:          getEnvString(dbHost, "postgres-automation"),
		DbPort:          dbNumPort,
		DbName:          getEnvString(dbName, "automation"),
		DbUser:          getEnvString(dbUser, "postgres"),
		DbPassword:      getEnvString(dbPassword, "postgres"),
		ImageMaxSize:    getEnvInt64(imageMaxSizeInMb, 5*1024*1024),
		ImageExtensions: getImageExtensions(),
		ImageSaveDir:    getEnvString(imageSaveDir, "images"),
	}
	ensureImageDirExists()
}

func ensureImageDirExists() {
	if _, err := os.Stat(AppConfig.ImageSaveDir); os.IsNotExist(err) {
		err := os.MkdirAll(AppConfig.ImageSaveDir, 0755)
		if err != nil {
			log.Fatalf("Failed to create directory %s: %v", AppConfig.ImageSaveDir, err)
		}
	}
}

func getImageExtensions() []string {
	extStr := getEnvString(imageExtensions, ".jpg,.png,.jpeg")
	return strings.Split(extStr, ",")
}

func validatePort(port int) error {
	if port < 0 || port > 65535 {
		return fmt.Errorf("error: Port %d is not valid", port)
	}
	return nil
}

func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		intVal, err := strconv.Atoi(value)
		if err == nil {
			return intVal
		}
	}
	log.Printf("Using default value for %s: %v", key, defaultValue)
	return defaultValue
}

func getEnvInt64(key string, defaultValue int64) int64 {
	if value, exists := os.LookupEnv(key); exists {
		intVal, err := strconv.ParseInt(value, 10, 64)
		if err == nil {
			return intVal
		}
	}
	log.Printf("Using default value for %s: %v", key, defaultValue)
	return defaultValue
}

func getEnvString(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	log.Printf("Using default value for %s: %s", key, defaultValue)
	return defaultValue
}
