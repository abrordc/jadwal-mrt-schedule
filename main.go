package main

import (
	"github.com/abrordc/jadwal-mrt-schedule/modules/station"
	"github.com/gin-gonic/gin"
)

func main(){
	InintiateRouter()
}

func InintiateRouter() {
	var (
		router = gin.Default()
		api = router.Group("/v1/api")
	)

	station.Initiate(api)

	router.Run(":8080")
}