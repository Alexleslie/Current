package main

import (
	"gee/gee"
	"time"

	"net/http"
)

func main() {

	r := gee.Default()
	r.Use(gee.Logger(), gee.Recovery())

	//r := gee.New()
	//r.Use(gee.Logger())
	//r.SetFuncMap(template.FuncMap{
	//	"FormatAsDate": FormatAsDate,
	//})
	//r.LoadHTMLGlob("templates/*")
	//r.Static("/assets", "./static")
	//
	//stu1 := &student{Name: "Yu", Age: 18}
	//stu2 := &student{Name: "Geektutu", Age: 20}
	//r.GET("/", func(c *gee.Context) {
	//	c.HTML(http.StatusOK, "css.tmpl", nil)
	//})
	//r.GET("/students", func(c *gee.Context) {
	//	c.HTML(http.StatusOK, "arr.tmpl", gee.H{
	//		"title":  "yu",
	//		"stuARR": [2]*student{stu1, stu2},
	//	})
	//})

	r.GET("/date", func(c *gee.Context) {
		c.WriteHTML(http.StatusOK, "custom_func.tmpl", gee.H{
			"title": "yu",
			"now":   time.Date(2022, 8, 2, 0, 0, 0, 0, time.UTC),
		})
	})

	r.GET("/", func(c *gee.Context) {
		c.WriteString(http.StatusOK, "Hello Yu\n")
	})

	r.GET("/panic", func(c *gee.Context) {
		names := []string{"Yu", "Geektutu"}
		c.WriteString(http.StatusOK, names[100])
	})
	err := r.Run(":9999")
	if err != nil {
		return
	}
}
