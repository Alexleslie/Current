package main

import (
	"Current/Gorm/geeORM"
	"Current/tools/logc"
	_ "github.com/mattn/go-sqlite3"
)

type Student struct {
	Name string
	Age  int
}

type Teacher struct {
	Name  string
	Age   int
	Level int
}

func main() {
	engine, err := geeORM.NewEngine("sqlite3", "gee.db")
	if err != nil {
		logc.Error("err=[%+v]", err)
	}
	session := engine.NewSession()
	session.SetSchemaByInstance(&Teacher{}).CreateTable()
	session.Insert(&Teacher{"Yu", 18, 1})
	var teacher []Teacher
	session.Limit(1).Where("Level", 2).FindAll(&teacher)
	logc.Info("[%+v]", teacher)
	session.DropTable()

}
