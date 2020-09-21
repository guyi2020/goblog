package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func main()  {
	router := gin.Default()
	router.POST("/post", func(c *gin.Context) {
		ids := c.QueryMap("ids")
		names := c.QueryMap("names")
		fmt.Printf("idssss: %v; namessss: %v", ids, names)
	})
	router.Run(":8080")
}
