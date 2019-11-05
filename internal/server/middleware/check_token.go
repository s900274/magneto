package middleware

import (
    "github.com/gin-gonic/gin"
)


func CheckToken() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Set("UserAuth", 1)
        c.Next()
        return
    }
}
