package main

import (
	"Current/Gee/gee"
	"Current/tools/logc"
	"net/http"
	"os"
)

func GetTodoList(ctx *gee.Context) {
	todoFileName := "./todo.txt"
	todoFileFd, err := os.OpenFile(todoFileName, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		logc.Error("[GetTodoList] Open file error, err=[%+v]", err)

	}
	var todo []byte
	todoFileFd.Read(todo)
	ctx.WriteString(http.StatusOK, "%v", string(todo))
}

func main() {
	engine := gee.Default()
	engine.GET("/todo", GetTodoList)
	engine.Run("127.0.0.1:529")
}
