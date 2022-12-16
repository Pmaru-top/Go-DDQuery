package main

import (
	"github.com/FishZe/Go-DDQuery"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strconv"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("/query", func(c *gin.Context) {
		if c.Query("uid") == "" && c.Query("name") == "" {
			c.JSON(http.StatusOK, gin.H{
				"error": "uid or name is required",
			})
			return
		}
		var uid int64
		var name string
		if c.Query("uid") != "" {
			uid, _ = strconv.ParseInt(c.Query("uid"), 10, 64)
		} else {
			name = c.Query("name")
		}
		route, err := Go_DDQuery.DDQuery(name, uid)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"error": err.Error(),
			})
			return
		}
		if c.Query("show") == "1" {
			file, _ := os.ReadFile(route)
			_, err2 := c.Writer.WriteString(string(file))
			if err2 != nil {
				c.JSON(http.StatusOK, gin.H{
					"error": err2.Error(),
				})
				return
			}
		} else {
			c.JSON(http.StatusOK, gin.H{
				"route": route,
				"error": nil,
			})
		}
	})
	err := r.Run(":20004")
	if err != nil {
		return
	}
}
