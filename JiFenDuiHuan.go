package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_"github.com/go-sql-driver/mysql"
)
func main(){
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/test?charset=utf8")
	if err!=nil{
		fmt.Println("Your open is wrong:",err)
	}
	router:=gin.Default()
	router.POST("/login", func(c *gin.Context) {
		cookie, err := c.Cookie("cookie")
		if err != nil {
			c.SetCookie("cookie", "v_cookie", 100, "/", "localhost", false, true)
			username := c.PostForm("username")
			Sl, err := db.Query("select point from Point where username=?", username)
			if err != nil {
				fmt.Println("Select error:", err)
			}
			defer Sl.Close()
			var point int
			for Sl.Next() {
				err := Sl.Scan(&point)
				if err != nil {
					fmt.Println("Scan point error:", err)
				}
			}
			point += 20
			stmt, err := db.Prepare("update point set point=? where username=?")
			if err != nil {
				fmt.Println("Update error:", err)
			}
			defer stmt.Close()
			stmt.Exec(point, username)
			c.JSON(200, gin.H{username: "签到成功，已获得20积分。"})
			}else {
			c.JSON(500, gin.H{"false": "您今天已经签到了，请明天再来吧。"})
		}
		fmt.Println(cookie)
	})
	router.POST("/point", func(c *gin.Context) {
		username:=c.PostForm("username")
		prize:=c.PostForm("prize")
		Sl1,err:=db.Query("select point from point where username=?")
		if err!=nil{
			fmt.Println("Select error:",err)
		}
		defer Sl1.Close()
		var point1 int
		for Sl1.Next(){
			err:=Sl1.Scan(&point1)
			if err!=nil{
				fmt.Println("Scan error:",err)
			}
		}
		fmt.Println(prize)
		Sl2,err:=db.Prepare("select point from prize where username=? ")
		if err!=nil{
			fmt.Println("Selecte errror:",err)
		}
		defer Sl2.Close()
		var point2 int
		for Sl1.Next(){
			err:=Sl1.Scan(&point2)
			if err!=nil{
				fmt.Println("Scan error:",err)
			}
		}
		if point1>=point2{
			point1-=point2
			stmt,_:=db.Prepare("update point set point=? where username=?")
			defer stmt.Close()
			stmt.Exec(point1,username)
			c.Writer.WriteString("兑换成功。")
		}else{
			c.Writer.WriteString("您的积分不足。")
		}
	})
	router.Run()
}
