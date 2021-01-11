package main

import (
	"api/app/controllers"
	"api/utils"
)

func main() {
	utils.LoggingSettings("tamarock.log")
	controllers.StartWebServer()
}
