package main

import (
	"fmt"
	"net/http"
	"os/exec"
	"strconv"

	Go_DDQuery "github.com/Pmaru-top/Go-DDQuery"
	"github.com/gin-gonic/gin"
)

func main() {
	query()

	//startHttpServer("127.0.0.1:8964")
}

func query() {
	var input string
	fmt.Println("请输入要查询的昵称或UID (如: 陈睿 / UID: 1):")
	_, err := fmt.Scanln(&input)
	if err != nil {
		return
	}
	var uid int64
	var name string
	if len(input) > 4 && (input[:4] == "UID:" || input[:4] == "uid:") {
		input = input[4:]
		uid, err = strconv.ParseInt(input, 10, 64)
		if err != nil {
			fmt.Println("输入的UID格式错误")
			return
		}
	} else {
		name = input
	}
	bytes, err := Go_DDQuery.DDQuery(name, uid)
	if err != nil {
		fmt.Println("查詢失敗:\n", err)
		return
	}

	directory, err := Go_DDQuery.SavePic("out", bytes)

	if err != nil {
		fmt.Println("保存失敗:\n", err)
		return
	}

	fmt.Println(directory)
	cmd := exec.Command("cmd", "/k", "start", directory)
	err = cmd.Start()
	if err != nil {
		fmt.Println("打開失敗，請手動打開\n", directory)
	}
}

func startHttpServer(bind string) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("/query", httpHandler)

	fmt.Println("HttpApi服務已在", bind, "上開啓")
	err := r.Run(bind)

	if err != nil {
		fmt.Println("HttpServer出現異常:\n", err)
		return
	}
}

func httpHandler(c *gin.Context) {
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
	bytes, err := Go_DDQuery.DDQuery(name, uid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}
	_, err = c.Writer.Write(bytes)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}
}