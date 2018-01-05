package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func adduser(c *gin.Context) {
	userid := c.Query("userid")
	entname := c.Query("entname")
	starttime := c.Query("starttime")
	endtime := c.Query("endtime")
	id := strconv.Itoa(keyid)
	userbyte := []byte(id + "|" + userid + "|" + entname + "|" + starttime + "|" + endtime + "\r\n")
	f, _ := os.OpenFile("user.txt", os.O_WRONLY|os.O_APPEND, 0666)
	f.Write(userbyte)
	timestart, _ := time.Parse("2006/01/02 15:04:05", starttime)
	timeend, _ := time.Parse("2006/01/02 15:04:05", endtime)

	for _, v := range temporaryUser {
		if v.userid == userid {
			if (timestart.Unix() < v.endtime.Unix() && timestart.Unix() > v.starttime.Unix()) || (timeend.Unix() < v.endtime.Unix() && timeend.Unix() > v.starttime.Unix()) {
				fmt.Println("添加记录有重复时间")
				c.String(200, "添加记录有重复时间")
				return
			}
		}
	}
	temporaryUser[keyid] = usertime{
		userid:    userid,
		entname:   entname,
		starttime: timestart,
		endtime:   timeend,
	}
	keyid++
	fmt.Println("user.txt文件写入:", string(userbyte))
	c.String(200, "添加记录"+string(userbyte))
}
func showuser(c *gin.Context) {
	fmt.Println("show temporaryUser:", temporaryUser)
	var returnstring string
	for k, v := range temporaryUser {
		kstr := strconv.Itoa(k)
		starttimestr := v.starttime.Format("2006/01/02 15:04:05")
		endtimestr := v.endtime.Format("2006/01/02 15:04:05")
		returnstring = returnstring + "{id:" + kstr + ",userid:" + v.userid + ",entname:" + v.entname + ",starttime:" + starttimestr + ",endtime:" + endtimestr + "}\r\n"
	}
	c.String(200, returnstring)
}

func deleteuser(c *gin.Context) {
	id := c.Query("id")
	fmt.Println("delete id:", id)
	idint, err := strconv.Atoi(id)
	if err != nil {
		c.String(200, "输入有误")
	}

	userid := temporaryUser[idint].userid
	useridiint, _ := strconv.Atoi(userid)
	if usermap[int64(useridiint)] != "" {
		delete(usermap, int64(useridiint))
	}
	userbyte := []byte("delete" + "|" + id + "\r\n")
	f, _ := os.OpenFile("user.txt", os.O_WRONLY|os.O_APPEND, 0666)
	f.Write(userbyte)
	delete(temporaryUser, idint)
	c.String(200, "delete id:"+id)
}

func updateuser(c *gin.Context) {
	id := c.Query("id")
	userid := c.Query("userid")
	entname := c.Query("entname")
	starttime := c.Query("starttime")
	endtime := c.Query("endtime")
	fmt.Println("update id:", id)
	idint, err := strconv.Atoi(id)
	if err != nil {
		c.String(200, "输入有误")
	}
	timestart, _ := time.Parse("2006/01/02 15:04:05", starttime)
	timeend, _ := time.Parse("2006/01/02 15:04:05", endtime)
	temporaryUser[idint] = usertime{
		userid:    userid,
		entname:   entname,
		starttime: timestart,
		endtime:   timeend,
	}
	userbyte := []byte("update" + id + "|" + userid + "|" + entname + "|" + starttime + "|" + endtime + "\r\n")
	f, _ := os.OpenFile("user.txt", os.O_WRONLY|os.O_APPEND, 0666)
	f.Write(userbyte)
	c.String(200, "update"+id+"|"+userid+"|"+entname+"|"+starttime+"|"+endtime)
}
