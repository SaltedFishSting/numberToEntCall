package main

import (
	"fmt"
	"strconv"
	"strings"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Meetqualified struct {
	Ctocquality string `bson:"c2cquality"`
	MeetingId   int64  `bson:"meetingId"`
	Userlist    string `bson:"userIdList`
}

type userc2c struct {
	Caller    string
	Called    string
	Meetingid int64
	Qualified string
	Cpu       string
	Usertype  string
}
type usermeet struct {
	Qualified float64
	Usertype  string
}

func getmeetqualifieds(DB *mgo.Collection, time int32) ([]Meetqualified, error) {
	var qualifieds []Meetqualified
	err := DB.Find(bson.M{"endTime": bson.M{"$gte": time}}).Select(bson.M{"c2cquality": 1, "meetingId": 1, "endTime": 1, "userIdList": 1}).All(&qualifieds)
	if err != nil {
		fmt.Println(err, "18")
	}
	var c2cs, c2c []string

	for _, v := range qualifieds {
		if v.Ctocquality == "" {
			continue
		}
		var userType string
		userType = "非商业"
		userlist := strings.Split(v.Userlist, "_")
		for _, v := range userlist {
			intuserid, _ := strconv.ParseInt(v, 10, 64)
			if usermap[intuserid] != "非商业" && usermap[intuserid] != "" {
				userType = "商业"
			}
		}

		c2cs = strings.Split(v.Ctocquality, "|")
		var a, b float64

		for _, v2 := range c2cs {
			c2c = strings.Split(v2, "/")

			a++
			if c2c[2] == "1" {
				b++
			}
			userc2cone := userc2c{
				Caller:    c2c[0],
				Called:    c2c[1],
				Meetingid: v.MeetingId,
				Qualified: c2c[2],
				Cpu:       c2c[4],
				Usertype:  userType,
			}
			q, _ := strconv.Atoi(c2c[3])
			meetc2cqualifiedmap[userc2cone] = float64(q)
		}
		um := usermeet{
			Qualified: b / a,
			Usertype:  userType,
		}
		meetqualified[int(v.MeetingId)] = um
	}
	return qualifieds, err
}
