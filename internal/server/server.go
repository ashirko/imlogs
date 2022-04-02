package server

import (
	"github.com/ashirko/imlogs/internal/utils"
	"github.com/gin-gonic/gin"
	"log"
)

const WEB_SERVER_HOST = "localhost"
const WEB_SERVER_PORT = "8081"

func Start() error {
	router := initRouter()
	addHandlers(router)
	address := getAddress()
	log.Printf("Starting Web Server on address %s...\n", address)
	err := router.Run(address)
	return err
}

func initRouter() *gin.Engine {
	return gin.Default()
}

func getAddress() string {
	host := utils.GetEnvStr("WEB_SERVER_HOST", WEB_SERVER_HOST)
	port := utils.GetEnvStr("WEB_SERVER_PORT", WEB_SERVER_PORT)
	return host + ":" + port
}
