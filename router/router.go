package router

import (
	service "github.com/shravanasati/scopex-go-assignment/service"

	"github.com/gin-gonic/gin"
)

// NewRoutes router global
func NewRoutes() *gin.Engine {

	router := gin.Default()
	v1 := router.Group("/api")

	// register router from each controller service
	service.RoutesLoginLogout(v1)
	service.RoutesUser(v1)
	service.RoutesUserDetail(v1)

	service.RoutesPost(v1)

	return router
}
