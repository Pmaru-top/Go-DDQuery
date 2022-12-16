package main

import (
	"fmt"
	"github.com/FishZe/Go-DDQuery"
	"os/exec"
)

func main() {
	route, err := Go_DDQuery.DDQuery("夜然z", 0)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(route)
	cmd := exec.Command("cmd", "/k", "start", route)
	err = cmd.Start()
	if err != nil {
		return
	}
}
