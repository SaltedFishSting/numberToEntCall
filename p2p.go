package main

import (
	"ebase"
	"fmt"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//从mongodb获取数据p2p通话统计 统计企业p2p的当前通话数，未接通通话数，新增通话数，通话时长
func getCall(collection *mgo.Collection) {

	looptime := int64(globeCfg.Output.Period)
	var p2pinfos []p2pinfo
	var nowtime p2pinfo

	err := collection.Find(bson.M{"isCaller": true}).Sort("-insertTime").Limit(1).One(&nowtime)
	if err != nil {
		fmt.Println(err, "21")
		return
	}
	var min10time int64
	if lastTime == 0 {
		lastTime = nowtime.InsertTime
		min10time = lastTime - looptime*1000
	} else {
		if nowtime.InsertTime == lastTime {
			min10time = lastTime
			lastTime = nowtime.InsertTime
		} else {
			min10time = lastTime + 1000
			lastTime = nowtime.InsertTime
		}

	}
	err = collection.Find(bson.M{"insertTime": bson.M{"$gte": min10time}, "isCaller": true}).Select(bson.M{"insertTime": 1, "starttime": 1, "endtime": 1, "called": 1, "caller": 1, "eventType": 1, "sidReporter": 1, "dissconnectedTime": 1}).All(&p2pinfos)
	if err != nil {
		fmt.Println(err, "34")
	}
	fmt.Println("p2p min10time:", min10time)
	fmt.Println("p2p通话总数(5m):", len(p2pinfos))
	fmt.Println("p2p通话内容(p2pinfos):", p2pinfos)
	for _, v := range p2pinfos {
		if v.Called > 99999999 || v.Caller > 99999999 {
			fmt.Println("出现不合理用户视讯号", v.Called, v.Caller)
			continue
		}
		if v.Starttime >= min10time {
			calledent := numberToEnt[numberArray[v.Called]]
			callerent := numberToEnt[numberArray[v.Caller]]
			entToNumber2[calledent]++
			entToNumber2[callerent]++
			if v.Called >= 999999 {
				err, newuser, monOlduser, weekOlduser, dayOlduser, _ := isActivep2p(v.Called, v.SidReporter)

				if err != nil {
					fmt.Println(err, "50")
				} else {
					if newuser == true {
						newusersp2p[calledent]++
					}
					if monOlduser == true {
						monOldusersp2p[calledent]++
					}
					if weekOlduser == true {
						weekOldusersp2p[calledent]++
					}
					if dayOlduser == true {
						dayOldusersp2p[calledent]++
					}
				}
			}
			if v.Caller >= 999999 {
				err, newuser, monOlduser, weekOlduser, dayOlduser, _ := isActivep2p(v.Caller, v.SidReporter)

				if err != nil {
					fmt.Println(err, "64")
				} else {
					if newuser == true {
						newusersp2p[callerent]++
					}
					if monOlduser == true {
						monOldusersp2p[callerent]++
					}
					if weekOlduser == true {
						weekOldusersp2p[calledent]++
					}
					if dayOlduser == true {
						dayOldusersp2p[callerent]++
					}
				}
			}

		}
		if v.InsertTime-v.Endtime >= 5*60*1000 {
			entToNumber4[numberToEnt[numberArray[v.Called]]]++
			entToNumber4[numberToEnt[numberArray[v.Caller]]]++
		}
		if v.DissconnectedTime != 0 {
			ent1 := numberToEnt[numberArray[v.Called]]
			ent2 := numberToEnt[numberArray[v.Caller]]
			time := v.DissconnectedTime - v.Starttime
			if time >= 0 {

				entToNumber5[ent1] = entToNumber5[ent1] + time
				entToNumber5[ent2] = entToNumber5[ent2] + time
			} else {
				fmt.Println("通话时间为负", v.SidReporter)
			}
		}
		entToNumber1[numberToEnt[numberArray[v.Called]]]++
		entToNumber1[numberToEnt[numberArray[v.Caller]]]++
		if v.EventType != 0 && v.EventType != -10000 {
			entToNumber3[numberToEnt[numberArray[v.Called]]]++
			entToNumber3[numberToEnt[numberArray[v.Caller]]]++
		}
	}
}

//统计企业通话用户的 地域，通话设备，使用的运营商网络类型，网络类型（wifi/有线/4G。。。）
func ispCityCall(collection *mgo.Collection) {
	fmt.Println("start")
	looptime := int64(globeCfg.Output.Period)
	var p2pispcitys []p2pispcity
	var nowtime p2pinfo
	err := collection.Find(bson.M{}).Select(bson.M{"insertTime": 1}).Sort("-insertTime").Limit(1).One(&nowtime)
	if err != nil {
		fmt.Println(err, "110")
		return
	}
	var time int64
	if lastTime2 == 0 {

		lastTime2 = nowtime.InsertTime
		time = lastTime2 - looptime*1000
	} else {
		time = lastTime2
		lastTime2 = nowtime.InsertTime
	}
	mapfunc := `   function(){
        for (var key in this.callBaseLog) {
            onelog = this.callBaseLog[key];
            if( (onelog.indexOf("eventtype=ext_dev_info")>=0)  
                ){
                idx0 = onelog.indexOf("loc_dom=");
                idx1 = onelog.indexOf("loc_isp=");
                idx2 = onelog.indexOf("loc_dev=");
                idx3 = onelog.indexOf("loc_dev_id=");
				idx4 = onelog.indexOf("rem_dom=");
				idx5 = onelog.indexOf("rem_isp=");
				idx6 = onelog.indexOf("rem_dev=");
				idx7 = onelog.indexOf("rem_dev_name=");
				idx8 = onelog.indexOf("loc_net=");
				idx9 = onelog.indexOf("rem_net=");
				
                ldom = onelog.substr(idx0+8, idx1-(idx0+9))
                lsip = onelog.substr(idx1+8, idx2-(idx1+9))
                ldev = onelog.substr(idx2+8, idx3-(idx2+9)) 
				
				rdom = onelog.substr(idx4+8, idx5-(idx4+9))
				rsip = onelog.substr(idx5+8, idx6-(idx5+9))               
				rdev = onelog.substr(idx6+8, idx7-(idx6+9))
				lnet = onelog.substr(idx8+8, 1) 
				rnet = onelog.substr(idx9+8, 1) 
				var cns = this.sid.split('_');
				number = cns[cns.length-1]
                emit(this.sid, {ldom:ldom,lsip:lsip,ldev:ldev,lnet:lnet,rdom:rdom,rsip:rsip,rdev:rdev,rnet:rnet,number:number});
            }
        }
    }`

	reducefunc := ` function(sid, info){
        return info[0]
    }`

	finalizefunc := ` function(sid, info){
            return info
        }`
	mapreduce := &mgo.MapReduce{
		Map:      mapfunc,
		Reduce:   reducefunc,
		Finalize: finalizefunc,
	}
	_, err = collection.Find(bson.M{"insertTime": bson.M{"$gte": time}}).MapReduce(mapreduce, &p2pispcitys)
	if err != nil {
		fmt.Println(err, "168")
		return
	}
	//	for _, v := range p2pispcitys {
	//		fmt.Println(v.Value.Ldev, "||", v.Value.Ldom, "||", v.Value.Lnet, "||", v.Value.Lsip, "||", v.Value.Number, "||", v.Value.Rdev, "||", v.Value.Rdom, "||", v.Value.Rnet, "||", v.Value.Rsip)
	//	}

	fmt.Println("p2p的ip统计总数(5m)", len(p2pispcitys))
	fmt.Println("p2p的ip统计内容(p2pispcitys)", p2pispcitys)
	for _, v := range p2pispcitys {
		domToNumber[v.Value.Ldom]++
		domToNumber[v.Value.Rdom]++
		ispToNumber[v.Value.Lsip]++
		ispToNumber[v.Value.Rsip]++
		devToNumber[v.Value.Ldev]++
		devToNumber[v.Value.Rdev]++
		netToNumber[v.Value.Lnet]++
		netToNumber[v.Value.Rnet]++
	}

}

//判断p2p用户是否活跃
func isActivep2p(user int64, callId string) (err error, newuser, monOlduser, weekOlduser, dayOlduser bool, ent string) {
	callId = ConvertToString(callId, "gbk", "utf-8")
	err, newuserbool, monOlduserbool, weekOlduserbool, dayOlduserbool, ent := ruc.AddAct("p2p", user, ebase.GetTimestampS(), 0, true, callId)
	return err, newuserbool, monOlduserbool, weekOlduserbool, dayOlduserbool, ent

}
