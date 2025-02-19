package middleware

import "github.com/gin-gonic/gin"

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
      clientToken := c.Request.Header.get("token");
      if clientToken == "" {
        
      }
  }
}