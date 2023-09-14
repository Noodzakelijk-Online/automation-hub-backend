package config

import (
	"log"
	"os"
	"strconv"
)

const (
	serverPort     string = "SERVER_PORT"
	baseUrl        string = "BASE_URL"
	nginxContainer string = "NGINX_CONTAINER"
	configDir      string = "CONFIG_DIR"
)

type Configuration struct {
	ConfigDir      string
	BaseUrl        string
	ServerPort     string
	NginxContainer string
}

var Config Configuration

func Init() {
	numPort := getEnvInt(serverPort, 80)
	if numPort < 0 || numPort > 65535 {
		log.Fatalf("error: Port %d is not valid, please check the environment variable: %s", numPort, serverPort)
	}
	Config = Configuration{
		ConfigDir:      getEnvString(configDir, "/app/sites-enabled"),
		BaseUrl:        getEnvString(baseUrl, "/config-manager"),
		ServerPort:     ":" + strconv.Itoa(numPort),
		NginxContainer: getEnvString(nginxContainer, "gateway"),
	}
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

func getEnvString(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	log.Printf("Using default value for %s: %s", key, defaultValue)
	return defaultValue
}
