//package lbapi
//
//import (
//	"github.com/gin-gonic/gin"
//	"github.com/go-resty/resty/v2"
//	lbutils "github.com/scarpart/distributed-task-scheduler/load-balancer/lb-utils"
//)
//
//func getAllNodes(ctx *gin.Context) {
//	client := resty.New()
//	request := client.R().
//		SetHeaders(lbutils.HeaderToMap(ctx.Request.Header)).
//		SetBody(ctx.Request.Body).
//		SetContext(ctx.Request.Context())
//
//	server := loadBalancer.PickServer()
//
//	// Forward the request to another server, based on the choice made by the load balancer
//	response, err := request.Execute(ctx.Request.Method, server.BaseUrl+ctx.Request.URL.Path)
//	if err != nil {	
//		ctx.JSON(response.StatusCode(), gin.H{"error": err.Error()})
//		return
//	}
//
//	// Return the response headers and body back to the client 
//	for key, values := range response.Header() {
//		for _, val := range values {
//			ctx.Writer.Header().Set(key, val)
//		}
//	}
//	ctx.Data(response.StatusCode(), response.Header().Get("Content-Type"), response.Body())
//}
