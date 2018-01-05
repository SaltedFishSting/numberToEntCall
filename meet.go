package main

import (
	"container/list"
	"ebase"
	"fmt"
	"strconv"
	"strings"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/wangtuanjie/ip17mon"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//统计各企业的通话时长meet 和 统计在线的用户数量，新增的用户数量
func meetEntToTime(meetsession *mgo.Session, nowtime meetinfo, min10time int32) {
	collection3 := meetsession.DB("meetingDB").C(globeCfg.Output.Table1meet)
	tablesmeet := strings.Split(nowtime.QosTableName, "/")
	tablemeet := "MeetingQosStatic" + tablesmeet[len(tablesmeet)-1]
	fmt.Println("本次查询时间:", min10time)
	fmt.Println("会议查询新表 ", tablemeet)
	//查询meet企业用户通话时间的mapreduce
	mapfunc := ` function() { 
        if (this.userId == this.speakerId &&((this.resourceId == 301 && this.deviceType == 5) || this.resourceId == 200 || (this.resourceId == 201 && this.deviceType == 6))){
            emit(this.userId, {addr: this.userIp, dev: this.deviceType, net: this.networkType, isnew:(this.origDetectedBW===-1), timelen: 5,meetid: this.meetingId}); 
        }
    }`
	reducefunc := `function(userId, info) { 
        reducedVal = { addr: info[0].addr, dev: info[0].dev, net: info[0].net, isnew: info[0].isnew, timelen: 0, meetid: info[0].meetid};
		   for (var idx = 0; idx < info.length; idx++) {
            reducedVal.timelen += info[idx].timelen;
        }

        return reducedVal;
    }`
	finalizefunc := `function(key, val){return val}`
	mapreduce := &mgo.MapReduce{
		Map:      mapfunc,
		Reduce:   reducefunc,
		Finalize: finalizefunc,
	}
	//查询meet的relay上下行流量的mapreduce
	mapfunc2 := ` function() { 
	 if ((this.resourceId == 301 && this.deviceType == 5) || this.resourceId == 200){
       var bw = this.bandwidth;
        if (bw < 0 || bw>100000){
            bw = 0;
        }
        var upbw = 0;
        var downbw = 0;
        if(this.userId === this.speakerId){
            upbw = bw;
        }else{
            downbw = bw;
        }
        emit(this.relayIp, {count:1, user: this.userId, upbw:upbw, downbw:downbw}); 
		}
    }`
	reducefunc2 := `function(relayIp, info) { 
        reducedVal = { usernum: 0, upbw: 0, downbw: 0, obj: {} };
        for (var idx = 0; idx < info.length; idx++) {
             reducedVal.obj[info[idx].user] = 1;
             reducedVal.upbw += info[idx].upbw;
             reducedVal.downbw += info[idx].downbw;
        }
        return reducedVal;
    }`
	finalizefunc2 := `function(key, val){
		if (val.obj != undefined){
			 delete val.obj["undefined"];
			   }
            return val
        }`
	mapreduce2 := &mgo.MapReduce{
		Map:      mapfunc2,
		Reduce:   reducefunc2,
		Finalize: finalizefunc2,
	}
	fmt.Println("lasttablname: ", lasttablemeet)
	if lasttablemeet == "meetTable" || lasttablemeet == tablemeet {
		var meet2MapReduces []meet2MapReduce
		var relayflows []relayflow
		collection := meetsession.DB("meetingDB").C(tablemeet)
		//用户体验统计
		userVideoAudioUpAndDown(collection, collection3, min10time)
		_, err := collection.Find(bson.M{"timeStamp": bson.M{"$gte": min10time}}).MapReduce(mapreduce, &meet2MapReduces)
		if err != nil {
			fmt.Println(err, "74")
			return
		}
		_, err = collection.Find(bson.M{"timeStamp": bson.M{"$gte": min10time}}).MapReduce(mapreduce2, &relayflows)
		if err != nil {
			fmt.Println(err, "92")
			return
		}
		for _, v := range relayflows {
			var ID string
			relaymap := globeCfg.Output.Relaymap
			ids := strings.Split(v.Id, ":")
			for k, v1 := range relaymap {
				if k == ids[0] {
					ID = v1
					break
				}
				ID = v.Id
			}
			meetrelayup[ID] += v.Value.Upbw
			meetrelaydown[ID] += v.Value.Downbw
		}
		fmt.Println("1 ", "查询meet总总数（5m）: ", len(meet2MapReduces), "表名: ", tablemeet)
		//fmt.Println("meet内容（meet2MapReduces）: ", meet2MapReduces)
		//fmt.Println("查询会议relay内容（relayflows）：", relayflows)

		for _, v := range meet2MapReduces {
			id, _ := strconv.Atoi(v.Id)
			//当前用户在商业会议中的数量
			isBusiness, _ := isBusinessmeet(collection3, v.Value.Meetid)
			if isBusiness {
				fmt.Println("1在商业会议中user", id)
				syuser["realuser"]++
			}

			if usermap[int64(id)] == "非商业" || usermap[int64(id)] == "" {
				continue
			}
			syuser["now"]++
			if v.Value.Isnew == true {
				intv2, _ := strconv.ParseInt(v.Id, 10, 64)
				meetidstring := strconv.Itoa(v.Value.Meetid)
				err, newuser, monOlduser, weekOlduser, dayOlduser, _ := isActivemeet(intv2, meetidstring)
				if err != nil {
					fmt.Println(err, "128")
				} else {
					if newuser == true {
						newusersmeet[usermap[intv2]]++
					}
					if monOlduser == true {
						monOldusersmeet[usermap[intv2]]++
					}
					if weekOlduser == true {
						weekOldusersmeet[usermap[intv2]]++
					}
					if dayOlduser == true {
						dayOldusersmeet[usermap[intv2]]++
					}
				}
				syuser["new"]++
			}

			entstring := usermap[int64(id)]
			fmt.Println("1商业用户", entstring, int64(id))
			entTocallTime[entstring] = entTocallTime[entstring] + v.Value.Timelen
			addr := strings.Split(v.Value.Addr, ":")
			ipToNumber(v.Id, addr[0], v.Value.Dev, v.Value.Net)
		}
		lasttablemeet = tablemeet
	} else {

		table1number := strings.Split(tablemeet, "_")
		table1lastnumber := strings.Split(lasttablemeet, "_")
		if table1number[0] == table1lastnumber[0] {
			lnumber, _ := strconv.Atoi(table1lastnumber[1])
			snumber, _ := strconv.Atoi(table1number[1])
			userlist := make(map[int]int)
			for lnumber <= snumber {
				lnumberstring := strconv.Itoa(lnumber)
				meettable := table1number[0] + "_" + lnumberstring
				var meet2MapReduces []meet2MapReduce
				var relayflows []relayflow

				collection2 := meetsession.DB("meetingDB").C(meettable)
				//用户体验统计
				userVideoAudioUpAndDown(collection2, collection3, min10time)
				_, err := collection2.Find(bson.M{"timeStamp": bson.M{"$gte": min10time}}).MapReduce(mapreduce, &meet2MapReduces)
				if err != nil {
					fmt.Println(err, "135")
					return
				}
				_, err = collection2.Find(bson.M{"timeStamp": bson.M{"$gte": min10time}}).MapReduce(mapreduce2, &relayflows)
				if err != nil {
					fmt.Println(err, "140")
					return
				}
				for _, v := range relayflows {
					var ID string
					relaymap := globeCfg.Output.Relaymap
					ids := strings.Split(v.Id, ":")
					for k, v1 := range relaymap {
						if k == ids[0] {
							ID = v1
							break
						}
						ID = v.Id
					}
					meetrelayup[ID] += v.Value.Upbw
					meetrelaydown[ID] += v.Value.Downbw
				}
				fmt.Println("2 ", "查询meet总总数（5m）: ", len(meet2MapReduces), "表名: ", meettable)
				//fmt.Println("meet内容（meet2MapReduces）: ", meet2MapReduces)
				//fmt.Println("查询会议relay内容（relayflows）：", relayflows)
				for _, v := range meet2MapReduces {
					id, _ := strconv.Atoi(v.Id)
					if userlist[id] == 0 {
						//用户是否在商业会议中
						isBusiness, _ := isBusinessmeet(collection3, v.Value.Meetid)

						if isBusiness {
							fmt.Println("2在商业会议中user", id)
							userlist[id] = 1
							syuser["realuser"]++
						}
						if usermap[int64(id)] == "非商业" || usermap[int64(id)] == "" {
							continue
						}
						//用户在线
						syuser["now"]++
						//新增用户
						if v.Value.Isnew == true {
							intv2, _ := strconv.ParseInt(v.Id, 10, 64)
							meetidstring := strconv.Itoa(v.Value.Meetid)
							err, newuser, monOlduser, weekOlduser, dayOlduser, _ := isActivemeet(intv2, meetidstring)
							if err != nil {
								fmt.Println(err, "213")
							} else {
								if newuser == true {
									newusersmeet[usermap[intv2]]++
								}
								if monOlduser == true {
									monOldusersmeet[usermap[intv2]]++
								}
								if weekOlduser == true {
									weekOldusersmeet[usermap[intv2]]++
								}
								if dayOlduser == true {
									dayOldusersmeet[usermap[intv2]]++
								}
							}
							syuser["new"]++
						}
						addr := strings.Split(v.Value.Addr, ":")
						//用户的 地域，运营商，网络类型，设备类型统计
						ipToNumber(v.Id, addr[0], v.Value.Dev, v.Value.Net)
						userlist[id] = 1
					}
					if usermap[int64(id)] == "非商业" || usermap[int64(id)] == "" {
						continue
					}
					//entstring := numberToEnt[numberArray[id]]
					entstring := usermap[int64(id)]
					fmt.Println("2商业用户", entstring, int64(id))
					entTocallTime[entstring] = entTocallTime[entstring] + v.Value.Timelen
				}
				lnumber++
			}
		} else {
			var meet2MapReduces []meet2MapReduce
			var relayflows []relayflow
			collection2 := meetsession.DB("meetingDB").C(tablemeet)
			//用户体验统计
			userVideoAudioUpAndDown(collection2, collection3, min10time)

			_, err := collection2.Find(bson.M{"timeStamp": bson.M{"$gte": min10time}}).MapReduce(mapreduce, &meet2MapReduces)
			if err != nil {
				fmt.Println(err, "178")
				return
			}
			_, err = collection2.Find(bson.M{"timeStamp": bson.M{"$gte": min10time}}).MapReduce(mapreduce2, &relayflows)
			if err != nil {
				fmt.Println(err, "183")
				return
			}
			for _, v := range relayflows {
				var ID string
				relaymap := globeCfg.Output.Relaymap
				ids := strings.Split(v.Id, ":")
				for k, v1 := range relaymap {
					if k == ids[0] {
						ID = v1
						break
					}
					ID = v.Id
				}
				meetrelayup[ID] += v.Value.Upbw
				meetrelaydown[ID] += v.Value.Downbw
			}
			fmt.Println("3 ", "查询meet总总数（5m）: ", len(meet2MapReduces), "表名: ", tablemeet)
			//	fmt.Println("meet内容（meet2MapReduces）: ", meet2MapReduces)
			//fmt.Println("查询会议relay内容（relayflows）：", relayflows)
			for _, v := range meet2MapReduces {

				//用户是否在商业会议中
				isBusiness, _ := isBusinessmeet(collection3, v.Value.Meetid)
				if isBusiness {
					syuser["realuser"]++
				}
				id, _ := strconv.Atoi(v.Id)
				if usermap[int64(id)] == "非商业" || usermap[int64(id)] == "" {
					continue
				}
				syuser["now"]++
				if v.Value.Isnew == true {
					intv2, _ := strconv.ParseInt(v.Id, 10, 64)
					meetidstring := strconv.Itoa(v.Value.Meetid)
					err, newuser, monOlduser, weekOlduser, dayOlduser, _ := isActivemeet(intv2, meetidstring)
					if err != nil {
						fmt.Println(err, "288")
					} else {
						if newuser == true {
							newusersmeet[usermap[intv2]]++
						}
						if monOlduser == true {
							monOldusersmeet[usermap[intv2]]++
						}
						if weekOlduser == true {
							weekOldusersmeet[usermap[intv2]]++
						}
						if dayOlduser == true {
							dayOldusersmeet[usermap[intv2]]++
						}
					}
					syuser["new"]++
				}
				//entstring := numberToEnt[numberArray[id]]
				entstring := usermap[int64(id)]
				fmt.Println("3商业用户", entstring, int64(id))
				entTocallTime[entstring] = entTocallTime[entstring] + v.Value.Timelen
				addr := strings.Split(v.Value.Addr, ":")

				ipToNumber(v.Id, addr[0], v.Value.Dev, v.Value.Net)
			}
		}
		lasttablemeet = tablemeet
	}

}

//从mongodb获取数据meet数据统计 统计当前的会议数量
func toPromtheus(ip string, db string, table1 string) {

	looptime := int32(globeCfg.Output.Period)
	session, err := mgo.Dial(ip)
	if err != nil {
		fmt.Println(err, "222")
		return
	}
	defer session.Close()
	collection := session.DB(db).C(table1)
	var nowtime meetinfo
	err = collection.Find(bson.M{}).Sort("-endTime").Limit(1).Select(bson.M{"endTime": 1, "qosTableName": 1, "meetingId": 1}).One(&nowtime)
	fmt.Println("meet表1的endtime", nowtime.EndTime)
	var min10time int32
	if lasttime == 0 {
		min10time = nowtime.EndTime - looptime
		lasttime = nowtime.EndTime
	} else {
		min10time = lasttime
		lasttime = nowtime.EndTime
	}

	meetEntToTime(session, nowtime, min10time)
	var meetinfos []meetinfo
	err = collection.Find(bson.M{"endTime": bson.M{"$gte": min10time}}).Select(bson.M{"meetingId": 1, "duration": 1, "startTime": 1, "endTime": 1, "userIdList": 1, "userCount": 1, "qosTableName": 1}).All(&meetinfos)
	if err != nil {
		fmt.Println(err, "243")
	}
	Observemeet(len(meetinfos))
	sumber := 0
	i := list.New()
	for _, v := range meetinfos {
		//fmt.Println(v.EndTime)
		strings := strings.Split(v.UserIdList, "_")
		for _, v2 := range strings {
			intv2, _ := strconv.ParseInt(v2, 10, 64)
			if intv2 < 999999 {
				continue
			}
			if usermap[intv2] != "非商业" && usermap[intv2] != "" {
				fmt.Println("商业会议id", v.MeetingId, intv2)
				i.PushBack(v.MeetingId)
				sumber++
				break
			}
		}
	}
	timelist = *i
	syObserve(sumber)

}

//统计地域 设备 运营商 网络类型 的在线人数(meet)
func ipToNumber(id string, ip string, dev int, net int) {
	loc, err := ip17mon.Find(ip)
	if err != nil {
		fmt.Println(err, "307")
		return
	}
	if loc.City == "N/A" {
		if loc.Region == "N/A" {
			if loc.Country == "N/A" {
				domToNumbermeet["unknown"]++
			} else {
				domToNumbermeet[loc.Country]++
			}
		} else {
			domToNumbermeet[loc.Region]++
		}
	} else {
		domToNumbermeet[loc.City]++
	}
	ispToNumbermeet[loc.Isp]++
	devstring := strconv.Itoa(dev)
	netstring := strconv.Itoa(net)
	if devmeetmap[devstring] == "" {
		devToNumbermeet["unknown"]++
	} else {
		devToNumbermeet[devmeetmap[devstring]]++
	}
	if netmeetmap[netstring] == "" {
		netToNumbermeet["unkeown"]++
	} else {
		netToNumbermeet[netmeetmap[netstring]]++
	}

}

//判断会议用户是否活跃
func isActivemeet(user int64, callId string) (err error, newuser, monOlduser, weekOlduser, dayOlduser bool, ent string) {
	callId = ConvertToString(callId, "gbk", "utf-8")
	err, newuserbool, monOlduserbool, weekOlduserbool, dayOlduserbool, ent := ruc.AddAct("conf", user, ebase.GetTimestampS(), 0, true, callId)
	return err, newuserbool, monOlduserbool, weekOlduserbool, dayOlduserbool, ent

}

func isBusinessmeet(collection *mgo.Collection, meetid int) (bool, int64) {
	var meetinfoOne meetinfo
	err := collection.Find(bson.M{"meetingId": meetid}).Select(bson.M{"userIdList": 1}).One(&meetinfoOne)
	if err != nil {
		fmt.Println(meetid)
		fmt.Println(err, "426")
		return false, 0
	}
	userlist := strings.Split(meetinfoOne.UserIdList, "_")
	for _, v := range userlist {
		intuserid, _ := strconv.ParseInt(v, 10, 64)
		if usermap[intuserid] != "非商业" && usermap[intuserid] != "" {
			return true, intuserid
		}
	}
	return false, 0

}
