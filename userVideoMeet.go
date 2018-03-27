package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/wangtuanjie/ip17mon"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//视频上传
func uservideoup(conn *mgo.Collection, collectionmeet *mgo.Collection, mintime int32) {

	mapfunc := `function() { 
        uid = parseInt(this.userId)
        if(this.userId == this.speakerId && 
              uid > 9999999 &&
			((this.resourceId==101 && this.deviceType==1)||(this.resourceId==101 && this.deviceType==5)||(this.resourceId==101 && this.deviceType==6)||(this.resourceId==100))
			){
			var org,fin
			if(this.lossRateOriginal==0){
				org=0
			}else{
				org=this.lossRateOriginal/100.0
			}
			if(this.lossRateFinal==0){
				fin=0
			}else{
				fin=this.lossRateFinal/100.0
			}
            emit({userid:this.userId, speakerid:this.speakerId,resourceId:this.resourceId,relayIp: this.relayIp} , 
                  {userIp: this.userIp, 
                      relayIp: this.relayIp, 
                      meetingId: this.meetingId,
                      deviceType: this.deviceType,
                      networkType: this.networkType,
					  resourceId:this.resourceId,
                      codeRate: this.curCodeRate, 
                      frameRate: this.frameRate, 
                      lossRateOrg: org,
                      lossRateFin: fin, 
                      delay: this.delayTimeIntArray,
                      odbw: this.origDetectedBW,
                      deviceType: this.deviceType,
                      networkType: this.networkType}
                      );
        }
    }`

	reducefunc := `  function(userId, info) {
        reducedVal = {
            userIp: info[0].userIp, 
            relayIp: info[0].relayIp, 
            meetingId: info[0].meetingId,
            deviceType: info[0].deviceType, 
            networkType: info[0].networkType,
			resourceId: info[0].resourceId,
            crArr: new Array(), 
            frArr: new Array(), 
            lossOrgArr: new Array(),
            loss: 0, 
            delayArr: new Array(), 
            delayLoss: 0,
            odbw: 0,
            dropline: false,
            deviceType: info[0].deviceType,
            networkType: info[0].networkType,
            };
        for (var idx = 0; idx < info.length; idx++) {
            if (info[idx].codeRate === undefined){
                Array.prototype.push.apply(reducedVal.crArr, info[idx].crArr);
                Array.prototype.push.apply(reducedVal.frArr, info[idx].frArr);
                Array.prototype.push.apply(reducedVal.lossOrgArr, info[idx].lossOrgArr);
                reducedVal.loss += info[idx].loss;
                reducedVal.delayLoss += info[idx].delayLoss;
                Array.prototype.push.apply(reducedVal.delayArr, info[idx].delayArr);
                reducedVal.odbw = info[idx].odbw;
                reducedVal.offline = info[idx].offline;
            }else{
                reducedVal.crArr.push(info[idx].codeRate);
                reducedVal.frArr.push(info[idx].frameRate);
                if (info[idx].lossRateFin > 0) {
                    reducedVal.loss += 1;
                }
                if (info[idx].lossRateOrg > 0) {
                    reducedVal.lossOrgArr.push(info[idx].lossRateOrg);
                }
                var ds = info[idx].delay.split('|');
                ds.forEach(function(d) {
                    if (d === "-1" || d === "") {
                        reducedVal.delayLoss += 1;
                    } else {
                        reducedVal.delayArr.push(parseInt(d));
                    }
                });
                if (info[idx].odbw !== -1 &&
                    reducedVal.odbw === 0){
                    reducedVal.odbw = info[idx].odbw
                }
                if (reducedVal.odbw !== 0 &&
                    info[idx].odbw === -1){
                    if (reducedVal.meetingId === info[idx].meetingId){
                        reducedVal.dropline = true;
                    }else{
                        reducedVal.meetingId = info[idx].meetingId
                        reducedVal.odbw = 0;
                    }
                }
            }
        }
        return reducedVal
    }`
	finalizefunc := `function(key, val){
           
		 if(val.loss == undefined){
			return {
                userip: "**",
                relayip: "",
                meetingid: "",
                crarr: "",
                frarr: "",
                lossorgarr: "",
                loss: 0, 
                delayarr: "",
                delayloss: 0,
                dropline: false,
                devicetype: 0,
                networktype: 0,

                }
		 }
            return {
                userip: val.userIp,
                relayip: val.relayIp,
                meetingid: val.meetingId,
                crarr: Array.join(val.crArr),
                frarr: Array.join(val.frArr),
                lossorgarr: Array.join(val.lossOrgArr),
                loss: val.loss, 
                delayarr: Array.join(val.delayArr),
                delayloss: val.delayLoss,
                dropline: val.dropline,
                devicetype: val.deviceType,
                networktype: val.networkType,
				resourceid: val.resourceId,

                }
        }`

	mapreduceup := &mgo.MapReduce{
		Map:      mapfunc,
		Reduce:   reducefunc,
		Finalize: finalizefunc,
	}

	var uservideosup []uservideoinfo

	_, err := conn.Find(bson.M{"timeStamp": bson.M{"$gte": mintime}}).MapReduce(mapreduceup, &uservideosup)
	if err != nil {
		fmt.Println(err, "372")
		if len(uservideosup) == 0 {
			return
		}
	}

	relaymap := globeCfg.Output.Relaymap

	//视频上行
	for _, v := range uservideosup {

		if v.Id.Userid == "" {
			continue
		}
		if v.Value.Userip == "**" {
			continue
		}
		usernumber, err := strconv.ParseInt(v.Id.Userid, 10, 64)
		if err != nil {
			continue
		}

		var userent string
		var meetent string
		var usertype string
		userent = usermap[usernumber]
		meetent = userent
		if usermap[usernumber] == "非商业" || usermap[usernumber] == "" {
			userent = "非商业"
			isBusiness, usetid := isBusinessmeet(collectionmeet, int(v.Value.Meetingid))
			if !isBusiness {
				usertype = "非商业"
				userent = "非商业"
				meetent = "非商业"

			} else {
				usertype = "商业"
				meetent = usermap[usetid]
			}

		} else {
			usertype = "商业"
		}
		if strings.Index(usermap[usernumber], "(演示)") != -1 {
			usertype = "演示"
		}

		var relayDomIsp string
		var userDom string
		var userIsp string

		var meetid string
		meetid = strconv.Itoa(int(v.Value.Meetingid))
		relayid := strings.Split(v.Value.Relayip, ":")
		if relaymap[relayid[0]] != "" {
			relayDomIsp = relaymap[relayid[0]]
		} else {
			relayDomIsp = "unknown"
		}
		userid := strings.Split(v.Value.Userip, ":")
		loc, err := ip17mon.Find(userid[0])
		if err != nil {
			fmt.Println(err, "162")
			return
		}
		if loc.City == "N/A" {
			if loc.Region == "N/A" {
				if loc.Country == "N/A" {
					userDom = "unknown"
				} else {
					userDom = loc.Country
				}
			} else {
				userDom = loc.Region
			}
		} else {
			userDom = loc.City
		}
		userIsp = loc.Isp

		userkeys := userkey{
			Userid:     v.Id.Userid,
			Meetintid:  meetid,
			Speakerid:  v.Id.Speakerid,
			Resourceid: v.Value.Resourceid,
			Relay:      relayDomIsp,
		}

		if uservideoqualityup[userkeys].Id.Userid == "" {

			user := uservideoinfo{}
			user.Id = v.Id
			user.Value.Usertype = usertype
			user.Value.Meetingid = v.Value.Meetingid
			user.Value.Userent = userent
			user.Value.Meetent = meetent
			user.Value.Devicetype = v.Value.Devicetype
			user.Value.Networktype = v.Value.Networktype
			user.Value.Resourceid = v.Value.Resourceid
			user.Value.Userdom = userDom
			user.Value.Userisp = userIsp
			user.Value.Relaydomisp = relayDomIsp
			user.Value.Crarr = v.Value.Crarr
			user.Value.Frarr = v.Value.Frarr
			user.Value.Delayarr = v.Value.Delayarr
			user.Value.Delayloss = v.Value.Delayloss
			user.Value.Loss = v.Value.Loss
			user.Value.Lossorgarr = v.Value.Lossorgarr
			user.Value.Dropline = v.Value.Dropline
			v.Value.Users = 1
			uservideoqualityup[userkeys] = user
		} else {
			user := uservideoqualityup[userkeys]
			user.Value.Relaydomisp = relayDomIsp
			user.Value.Crarr += "," + v.Value.Crarr
			user.Value.Frarr += "," + v.Value.Frarr
			user.Value.Delayarr += "," + v.Value.Delayarr
			user.Value.Delayloss += v.Value.Delayloss
			user.Value.Loss += v.Value.Loss
			user.Value.Lossorgarr += "," + v.Value.Lossorgarr
			if v.Value.Dropline == true {
				user.Value.Dropline = true
			}
			v.Value.Users += 1
			uservideoqualityup[userkeys] = user
		}

	}

}

//视频下行
func uservideodown(conn *mgo.Collection, collectionmeet *mgo.Collection, mintime int32) {

	mapfunc2 := `function() { 
        uid = parseInt(this.userId)
        if(this.userId != this.speakerId && 
              uid > 9999999 &&
			((this.resourceId==101 && this.deviceType==1)||(this.resourceId==101 && this.deviceType==5)||(this.resourceId==101 && this.deviceType==6)||(this.resourceId==100))
			){
				var org,fin
			if(this.lossRateOriginal==0){
				org=0
			}else{
				org=this.lossRateOriginal/100.0
			}
			if(this.lossRateFinal==0){
				fin=0
			}else{
				fin=this.lossRateFinal/100.0
			}
            emit({userid:this.userId, speakerid:this.speakerId,resourceId:this.resourceId,relayIp: this.relayIp} , 
                  {userIp: this.userIp, 
                      relayIp: this.relayIp, 
                      meetingId: this.meetingId,
                      deviceType: this.deviceType,
                      networkType: this.networkType,
					  resourceId:this.resourceId, 
                      codeRate: this.curCodeRate, 
                      frameRate: this.frameRate, 
                      lossRateOrg: org,
                      lossRateFin: fin, 
                      delay: this.delayTimeIntArray,
                      odbw: this.origDetectedBW,
                      deviceType: this.deviceType,
                      networkType: this.networkType}
                      );
        }
    }`
	reducefunc := `  function(userId, info) {
        reducedVal = {
            userIp: info[0].userIp, 
            relayIp: info[0].relayIp, 
            meetingId: info[0].meetingId,
            deviceType: info[0].deviceType, 
            networkType: info[0].networkType, 
			resourceId: info[0].resourceId, 
            crArr: new Array(), 
            frArr: new Array(), 
            lossOrgArr: new Array(),
            loss: 0, 
            delayArr: new Array(), 
            delayLoss: 0,
            odbw: 0,
            dropline: false,
            deviceType: info[0].deviceType,
            networkType: info[0].networkType,
            };
        for (var idx = 0; idx < info.length; idx++) {
            if (info[idx].codeRate === undefined){
                Array.prototype.push.apply(reducedVal.crArr, info[idx].crArr);
                Array.prototype.push.apply(reducedVal.frArr, info[idx].frArr);
                Array.prototype.push.apply(reducedVal.lossOrgArr, info[idx].lossOrgArr);
                reducedVal.loss += info[idx].loss;
                reducedVal.delayLoss += info[idx].delayLoss;
                Array.prototype.push.apply(reducedVal.delayArr, info[idx].delayArr);
                reducedVal.odbw = info[idx].odbw;
                reducedVal.offline = info[idx].offline;
            }else{
                reducedVal.crArr.push(info[idx].codeRate);
                reducedVal.frArr.push(info[idx].frameRate);
                if (info[idx].lossRateFin > 0) {
                    reducedVal.loss += 1;
                }
                if (info[idx].lossRateOrg > 0) {
                    reducedVal.lossOrgArr.push(info[idx].lossRateOrg);
                }
                var ds = info[idx].delay.split('|');
                ds.forEach(function(d) {
                    if (d === "-1" || d === "") {
                        reducedVal.delayLoss += 1;
                    } else {
                        reducedVal.delayArr.push(parseInt(d));
                    }
                });
                if (info[idx].odbw !== -1 &&
                    reducedVal.odbw === 0){
                    reducedVal.odbw = info[idx].odbw
                }
                if (reducedVal.odbw !== 0 &&
                    info[idx].odbw === -1){
                    if (reducedVal.meetingId === info[idx].meetingId){
                        reducedVal.dropline = true;
                    }else{
                        reducedVal.meetingId = info[idx].meetingId
                        reducedVal.odbw = 0;
                    }
                }
            }
        }
        return reducedVal
    }`
	finalizefunc := `function(key, val){
           
		 if(val.loss == undefined){
			return {
                userip: "**",
                relayip: "",
                meetingid: "",
                crarr: "",
                frarr: "",
                lossorgarr: "",
                loss: 0, 
                delayarr: "",
                delayloss: 0,
                dropline: false,
                devicetype: 0,
                networktype: 0,

                }
		 }
            return {
                userip: val.userIp,
                relayip: val.relayIp,
                meetingid: val.meetingId,
                crarr: Array.join(val.crArr),
                frarr: Array.join(val.frArr),
                lossorgarr: Array.join(val.lossOrgArr),
                loss: val.loss, 
                delayarr: Array.join(val.delayArr),
                delayloss: val.delayLoss,
                dropline: val.dropline,
                devicetype: val.deviceType,
                networktype: val.networkType,
				resourceid: val.resourceId,
                }
        }`

	mapreducedown := &mgo.MapReduce{
		Map:      mapfunc2,
		Reduce:   reducefunc,
		Finalize: finalizefunc,
	}

	var uservideosdown []uservideoinfo

	_, err := conn.Find(bson.M{"timeStamp": bson.M{"$gte": mintime}}).MapReduce(mapreducedown, &uservideosdown)
	if err != nil {
		fmt.Println(err, "379")
		if len(uservideosdown) == 0 {
			return
		}
	}
	relaymap := globeCfg.Output.Relaymap

	//视频下行
	for _, v := range uservideosdown {

		if v.Id.Userid == "" {
			continue
		}
		if v.Value.Userip == "**" {
			continue
		}
		usernumber, err := strconv.ParseInt(v.Id.Userid, 10, 64)
		if err != nil {
			continue
		}
		var userent string
		var meetent string
		var usertype string
		userent = usermap[usernumber]
		meetent = userent
		if usermap[usernumber] == "非商业" || usermap[usernumber] == "" {
			userent = "非商业"
			isBusiness, usetid := isBusinessmeet(collectionmeet, int(v.Value.Meetingid))
			if !isBusiness {
				usertype = "非商业"
				userent = "非商业"
				meetent = "非商业"

			} else {
				usertype = "商业"
				meetent = usermap[usetid]
			}

		} else {
			usertype = "商业"
		}
		if strings.Index(usermap[usernumber], "(演示)") != -1 {
			usertype = "演示"
		}
		var relayDomIsp string
		var userDom string
		var userIsp string

		var meetid string
		meetid = strconv.Itoa(int(v.Value.Meetingid))
		relayid := strings.Split(v.Value.Relayip, ":")
		if relaymap[relayid[0]] != "" {
			relayDomIsp = relaymap[relayid[0]]
		} else {
			relayDomIsp = "unknown"
		}
		userid := strings.Split(v.Value.Userip, ":")
		loc, err := ip17mon.Find(userid[0])
		if err != nil {
			fmt.Println(err, "162")
			return
		}
		if loc.City == "N/A" {
			if loc.Region == "N/A" {
				if loc.Country == "N/A" {
					userDom = "unknown"
				} else {
					userDom = loc.Country
				}
			} else {
				userDom = loc.Region
			}
		} else {
			userDom = loc.City
		}
		userIsp = loc.Isp

		if relaymap[relayid[0]] != "" {
			relayDomIsp = relaymap[relayid[0]]
		} else {
			relayDomIsp = "unknown"
		}

		userkeys := userkey{
			Userid:     v.Id.Userid,
			Meetintid:  meetid,
			Speakerid:  v.Id.Speakerid,
			Resourceid: v.Value.Resourceid,
			Relay:      relayDomIsp,
		}

		if uservideoqualitydown[userkeys].Id.Userid == "" {

			user := uservideoinfo{}
			user.Id = v.Id
			user.Value.Usertype = usertype
			user.Value.Meetingid = v.Value.Meetingid
			user.Value.Userent = userent
			user.Value.Meetent = meetent
			user.Value.Devicetype = v.Value.Devicetype
			user.Value.Networktype = v.Value.Networktype
			user.Value.Resourceid = v.Value.Resourceid
			user.Value.Userdom = userDom
			user.Value.Userisp = userIsp
			user.Value.Relaydomisp = relayDomIsp
			user.Value.Crarr = v.Value.Crarr
			user.Value.Frarr = v.Value.Frarr
			user.Value.Delayarr = v.Value.Delayarr
			user.Value.Delayloss = v.Value.Delayloss
			user.Value.Loss = v.Value.Loss
			user.Value.Lossorgarr = v.Value.Lossorgarr
			user.Value.Dropline = v.Value.Dropline
			v.Value.Users = 1
			uservideoqualitydown[userkeys] = user
		} else {
			user := uservideoqualitydown[userkeys]
			user.Value.Relaydomisp = relayDomIsp
			user.Value.Crarr += "," + v.Value.Crarr
			user.Value.Frarr += "," + v.Value.Frarr
			user.Value.Delayarr += "," + v.Value.Delayarr
			user.Value.Delayloss += v.Value.Delayloss
			user.Value.Loss += v.Value.Loss
			user.Value.Lossorgarr += "," + v.Value.Lossorgarr
			if v.Value.Dropline == true {
				user.Value.Dropline = true
			}
			v.Value.Users += 1
			uservideoqualitydown[userkeys] = user
		}

	}

}

//音频上行
func useraudioup(conn *mgo.Collection, collectionmeet *mgo.Collection, mintime int32) {

	mapfunc3 := `function() { 
        uid = parseInt(this.userId)
        if(this.userId == this.speakerId &&
		                   uid > 9999999 &&
				((this.resourceId == 301 && this.deviceType == 5) || this.resourceId == 200 || (this.resourceId == 201 && this.deviceType == 6))
				){
					var org,fin
			if(this.lossRateOriginal==0){
				org=0
			}else{
				org=this.lossRateOriginal/100.0
			}
			if(this.lossRateFinal==0){
				fin=0
			}else{
				fin=this.lossRateFinal/100.0
			}
            emit({userid:this.userId, speakerid:this.speakerId,resourceId:this.resourceId,relayIp: this.relayIp} , 
                  {userIp: this.userIp, 
                      relayIp: this.relayIp, 
                      meetingId: this.meetingId,
                      deviceType: this.deviceType,
                      networkType: this.networkType,             
                      deviceType: this.deviceType,
                      networkType: this.networkType,
					  lossRateOrg: org,
                      lossRateFin: fin, 
					  one: this.oneEmpty,
					  two: this.twoEmpty,
					  three: this.thrEmpty,
					  foure: this.fouEmpty,
					  ten: this.tenEmpty}
                      );
        }
    }`

	reducefunc2 := `  function(userId, info) {
        reducedVal = {
            userIp: info[0].userIp, 
            relayIp: info[0].relayIp, 
            meetingId: info[0].meetingId,
            deviceType: info[0].deviceType, 
            networkType: info[0].networkType,             
            deviceType: info[0].deviceType,
            networkType: info[0].networkType,
			lossOrgArr: new Array(),
            loss: 0, 
			one: 0,
			two: 0,
			three: 0,
			foure: 0,
			ten: 0,
            };
        for (var idx = 0; idx < info.length; idx++) {
            if (info[idx].userIp === undefined){
				reducedVal.one += info[idx].one;
				reducedVal.two += info[idx].two;
				reducedVal.three += info[idx].three;
				reducedVal.foure += info[idx].foure;
				reducedVal.ten += info[idx].ten;
				Array.prototype.push.apply(reducedVal.lossOrgArr, info[idx].lossOrgArr);
				reducedVal.loss += info[idx].loss;
            }else{
				 if (info[idx].lossRateFin > 0) {
                    reducedVal.loss += 1;
                }
                if (info[idx].lossRateOrg > 0) {
                    reducedVal.lossOrgArr.push(info[idx].lossRateOrg);
                }
				if (info[idx].one > 0){
					reducedVal.one+=info[idx].one;
				}
			    if (info[idx].two > 0){
					reducedVal.two+=info[idx].two;
				}
				if (info[idx].three > 0){
					reducedVal.three+=info[idx].three;
				}
				if (info[idx].foure > 0){
					reducedVal.foure+=info[idx].foure;
				}
				if (info[idx].ten > 0){
					reducedVal.ten+=info[idx].ten;
				}		                                                          
            }
        }
        return reducedVal
    }`

	finalizefunc2 := `function(key, val){
            return {
                userip: val.userIp,
                relayip: val.relayIp,
                meetingid: val.meetingId,     
                devicetype: val.deviceType,
                networktype: val.networkType,
				oneempty: val.one,
				twoempty: val.two,
				threeempty: val.three,
				fourempty: val.foure,
				tenempty: val.ten,
				lossorgarr: Array.join(val.lossOrgArr),
                loss: val.loss, 
                }
        }`
	mapreduceaudioup := &mgo.MapReduce{
		Map:      mapfunc3,
		Reduce:   reducefunc2,
		Finalize: finalizefunc2,
	}

	var uservideosaudioup []uservideoinfo

	_, err := conn.Find(bson.M{"timeStamp": bson.M{"$gte": mintime}}).MapReduce(mapreduceaudioup, &uservideosaudioup)
	if err != nil {
		fmt.Println(err, "386")
		if len(uservideosaudioup) == 0 {
			return
		}
	}
	relaymap := globeCfg.Output.Relaymap
	//音频上行
	for _, v := range uservideosaudioup {

		if v.Id.Userid == "" {
			continue
		}
		if v.Value.Meetingid == 0 {
			continue
		}
		usernumber, err := strconv.ParseInt(v.Id.Userid, 10, 64)
		if err != nil {
			continue
		}
		var userent string
		var meetent string
		var usertype string
		userent = usermap[usernumber]
		meetent = userent
		if usermap[usernumber] == "非商业" || usermap[usernumber] == "" {
			userent = "非商业"
			isBusiness, usetid := isBusinessmeet(collectionmeet, int(v.Value.Meetingid))
			if !isBusiness {
				usertype = "非商业"
				userent = "非商业"
				meetent = "非商业"

			} else {
				usertype = "商业"
				meetent = usermap[usetid]
			}

		} else {
			usertype = "商业"
		}
		if strings.Index(usermap[usernumber], "(演示)") != -1 {
			usertype = "演示"
		}
		var relayDomIsp string
		var userDom string
		var userIsp string

		var meetid string
		meetid = strconv.Itoa(int(v.Value.Meetingid))
		relayid := strings.Split(v.Value.Relayip, ":")
		if relaymap[relayid[0]] != "" {
			relayDomIsp = relaymap[relayid[0]]
		} else {
			relayDomIsp = "unknown"
		}
		userid := strings.Split(v.Value.Userip, ":")
		loc, err := ip17mon.Find(userid[0])
		if err != nil {
			fmt.Println(err, "162")
			return
		}
		if loc.City == "N/A" {
			if loc.Region == "N/A" {
				if loc.Country == "N/A" {
					userDom = "unknown"
				} else {
					userDom = loc.Country
				}
			} else {
				userDom = loc.Region
			}
		} else {
			userDom = loc.City
		}
		userIsp = loc.Isp

		if relaymap[relayid[0]] != "" {
			relayDomIsp = relaymap[relayid[0]]
		} else {
			relayDomIsp = "unknown"
		}

		userkeys := userkey{
			Userid:     v.Id.Userid,
			Meetintid:  meetid,
			Speakerid:  v.Id.Speakerid,
			Resourceid: v.Value.Resourceid,
			Relay:      relayDomIsp,
		}
		if uservideoqualityaudioup[userkeys].Id.Userid == "" {

			user := uservideoinfo{}
			user.Id = v.Id
			user.Value.Usertype = usertype
			user.Value.Meetingid = v.Value.Meetingid
			user.Value.Userent = userent
			user.Value.Meetent = meetent
			user.Value.Devicetype = v.Value.Devicetype
			user.Value.Networktype = v.Value.Networktype
			user.Value.Userdom = userDom
			user.Value.Userisp = userIsp
			user.Value.Relaydomisp = relayDomIsp
			user.Value.Loss = v.Value.Loss
			user.Value.Lossorgarr = v.Value.Lossorgarr
			user.Value.Oneempty = v.Value.Oneempty
			user.Value.Twoempty = v.Value.Twoempty
			user.Value.Threeempty = v.Value.Threeempty
			user.Value.Fourempty = v.Value.Fourempty
			user.Value.Tenempty = v.Value.Tenempty
			user.Value.Users = 1
			uservideoqualityaudioup[userkeys] = user
		} else {
			user := uservideoqualityaudioup[userkeys]
			user.Value.Relaydomisp = relayDomIsp
			user.Value.Loss += v.Value.Loss
			user.Value.Lossorgarr += "," + v.Value.Lossorgarr
			user.Value.Oneempty += v.Value.Oneempty
			user.Value.Twoempty += v.Value.Twoempty
			user.Value.Threeempty += v.Value.Threeempty
			user.Value.Fourempty += v.Value.Fourempty
			user.Value.Tenempty += v.Value.Tenempty

			user.Value.Users += 1
			uservideoqualityaudioup[userkeys] = user
		}

	}

}

//音频下行
func useraudiodown(conn *mgo.Collection, collectionmeet *mgo.Collection, mintime int32) {

	mapfunc4 := `function() { 
        uid = parseInt(this.userId)
        if(this.userId != this.speakerId &&
		                   uid > 9999999 &&
				((this.resourceId == 301 && this.deviceType == 5) || this.resourceId == 200 || (this.resourceId == 201 && this.deviceType == 6))
				){
					var org,fin
			if(this.lossRateOriginal==0){
				org=0
			}else{
				org=this.lossRateOriginal/100
			}
			if(this.lossRateFinal==0){
				fin=0
			}else{
				fin=this.lossRateFinal/100
			}
            emit({userid:this.userId, speakerid:this.speakerId,resourceId:this.resourceId,relayIp: this.relayIp}, 
                  {userIp: this.userIp, 
                      relayIp: this.relayIp, 
                      meetingId: this.meetingId,
                      deviceType: this.deviceType,
                      networkType: this.networkType,                      
                      deviceType: this.deviceType,
                      networkType: this.networkType,
					  lossRateOrg: org,
                      lossRateFin: fin, 
					  one: this.oneEmpty,
					  two: this.twoEmpty,
					  three: this.thrEmpty,
					  foure: this.fouEmpty,
					  ten: this.tenEmpty}
                      );
        }
    }`

	reducefunc2 := `  function(userId, info) {
        reducedVal = {
            userIp: info[0].userIp, 
            relayIp: info[0].relayIp, 
            meetingId: info[0].meetingId,
            deviceType: info[0].deviceType, 
            networkType: info[0].networkType,            
            deviceType: info[0].deviceType,
            networkType: info[0].networkType,
			lossOrgArr: new Array(),
            loss: 0, 
			one: 0,
			two: 0,
			three: 0,
			foure: 0,
			ten: 0,
            };
        for (var idx = 0; idx < info.length; idx++) {
            if (info[idx].userIp === undefined){
				reducedVal.one += info[idx].one;
				reducedVal.two += info[idx].two;
				reducedVal.three += info[idx].three;
				reducedVal.foure += info[idx].foure;
				reducedVal.ten += info[idx].ten;
				Array.prototype.push.apply(reducedVal.lossOrgArr, info[idx].lossOrgArr);
				reducedVal.loss += info[idx].loss;
            }else{
				 if (info[idx].lossRateFin > 0) {
                    reducedVal.loss += 1;
                }
                if (info[idx].lossRateOrg > 0) {
                    reducedVal.lossOrgArr.push(info[idx].lossRateOrg);
                }
				if (info[idx].one > 0){
					reducedVal.one+=info[idx].one;
				}
			    if (info[idx].two > 0){
					reducedVal.two+=info[idx].two;
				}
				if (info[idx].three > 0){
					reducedVal.three+=info[idx].three;
				}
				if (info[idx].foure > 0){
					reducedVal.foure+=info[idx].foure;
				}
				if (info[idx].ten > 0){
					reducedVal.ten+=info[idx].ten;
				}		                                          
            }
        }
        return reducedVal
    }`

	finalizefunc2 := `function(key, val){
           
		
		 if(val.loss == undefined){
			return {
                userip: "**",
                relayip: "",
                meetingid: "",                          
                lossorgarr: "",
                loss: 0,                               
                devicetype: 0,
                networktype: 0,

                }
		 }
            return {
                userip: val.userIp,
                relayip: val.relayIp,
                meetingid: val.meetingId,     
                devicetype: val.deviceType,
                networktype: val.networkType,
				oneempty: val.one,
				twoempty: val.two,
				threeempty: val.three,
				fourempty: val.foure,
				tenempty: val.ten,
				lossorgarr: Array.join(val.lossOrgArr),
                loss: val.loss, 
                }
        }`

	mapreduceaudiodown := &mgo.MapReduce{
		Map:      mapfunc4,
		Reduce:   reducefunc2,
		Finalize: finalizefunc2,
	}

	var uservideosaudiodown []uservideoinfo
	_, err := conn.Find(bson.M{"timeStamp": bson.M{"$gte": mintime}}).MapReduce(mapreduceaudiodown, &uservideosaudiodown)
	if err != nil {
		fmt.Println(err, "393")
		if len(uservideosaudiodown) == 0 {
			return
		}
	}

	relaymap := globeCfg.Output.Relaymap
	//音频下行
	for _, v := range uservideosaudiodown {

		if v.Id.Userid == "" {
			continue
		}
		if v.Value.Meetingid == 0 {
			continue
		}
		usernumber, err := strconv.ParseInt(v.Id.Userid, 10, 64)
		if err != nil {
			continue
		}
		var userent string
		var meetent string
		var usertype string
		userent = usermap[usernumber]
		meetent = userent
		if usermap[usernumber] == "非商业" || usermap[usernumber] == "" {
			userent = "非商业"
			isBusiness, usetid := isBusinessmeet(collectionmeet, int(v.Value.Meetingid))
			if !isBusiness {
				usertype = "非商业"
				userent = "非商业"
				meetent = "非商业"

			} else {
				usertype = "商业"
				meetent = usermap[usetid]
			}

		} else {
			usertype = "商业"
		}
		if strings.Index(usermap[usernumber], "(演示)") != -1 {
			usertype = "演示"
		}
		var relayDomIsp string
		var userDom string
		var userIsp string

		var meetid string
		meetid = strconv.Itoa(int(v.Value.Meetingid))
		relayid := strings.Split(v.Value.Relayip, ":")
		if relaymap[relayid[0]] != "" {
			relayDomIsp = relaymap[relayid[0]]
		} else {
			relayDomIsp = "unknown"
		}
		userid := strings.Split(v.Value.Userip, ":")
		loc, err := ip17mon.Find(userid[0])
		if err != nil {
			fmt.Println(err, "162")
			return
		}
		if loc.City == "N/A" {
			if loc.Region == "N/A" {
				if loc.Country == "N/A" {
					userDom = "unknown"
				} else {
					userDom = loc.Country
				}
			} else {
				userDom = loc.Region
			}
		} else {
			userDom = loc.City
		}
		userIsp = loc.Isp

		if relaymap[relayid[0]] != "" {
			relayDomIsp = relaymap[relayid[0]]
		} else {
			relayDomIsp = "unknown"
		}

		userkeys := userkey{
			Userid:     v.Id.Userid,
			Meetintid:  meetid,
			Speakerid:  v.Id.Speakerid,
			Resourceid: v.Value.Resourceid,
			Relay:      relayDomIsp,
		}
		if uservideoqualityaudiodown[userkeys].Id.Userid == "**" {
			continue
		}
		if uservideoqualityaudiodown[userkeys].Id.Userid == "" {

			user := uservideoinfo{}
			user.Id = v.Id
			user.Value.Usertype = usertype
			user.Value.Meetingid = v.Value.Meetingid
			user.Value.Userent = userent
			user.Value.Meetent = meetent
			user.Value.Devicetype = v.Value.Devicetype
			user.Value.Networktype = v.Value.Networktype
			user.Value.Userdom = userDom
			user.Value.Userisp = userIsp
			user.Value.Relaydomisp = relayDomIsp
			user.Value.Loss = v.Value.Loss
			user.Value.Lossorgarr = v.Value.Lossorgarr
			user.Value.Oneempty = v.Value.Oneempty
			user.Value.Twoempty = v.Value.Twoempty
			user.Value.Threeempty = v.Value.Threeempty
			user.Value.Fourempty = v.Value.Fourempty
			user.Value.Tenempty = v.Value.Tenempty
			user.Value.Users = 1
			uservideoqualityaudiodown[userkeys] = user
		} else {
			user := uservideoqualityaudiodown[userkeys]
			user.Value.Relaydomisp = relayDomIsp
			user.Value.Loss += v.Value.Loss
			user.Value.Lossorgarr += "," + v.Value.Lossorgarr
			user.Value.Oneempty += v.Value.Oneempty
			user.Value.Twoempty += v.Value.Twoempty
			user.Value.Threeempty += v.Value.Threeempty
			user.Value.Fourempty += v.Value.Fourempty
			user.Value.Tenempty += v.Value.Tenempty
			user.Value.Users += 1
			uservideoqualityaudiodown[userkeys] = user
		}

	}

}

//文档上行
func userfileup(conn *mgo.Collection, collectionmeet *mgo.Collection, mintime int32) {

	mapfunc := `function() { 
        uid = parseInt(this.userId)
        if(this.userId == this.speakerId && 
              uid > 9999999 &&
			this.mediaType==2
			){
			var org,fin
			if(this.lossRateOriginal==0){
				org=0
			}else{
				org=this.lossRateOriginal/100.0
			}
			if(this.lossRateFinal==0){
				fin=0
			}else{
				fin=this.lossRateFinal/100.0
			}
            emit({userid:this.userId, speakerid:this.speakerId,resourceId:this.resourceId,relayIp: this.relayIp} , 
                  {userIp: this.userIp, 
                      relayIp: this.relayIp, 
                      meetingId: this.meetingId,
                      deviceType: this.deviceType,
                      networkType: this.networkType,
					  resourceId:this.resourceId,
                      codeRate: this.curCodeRate, 
                      frameRate: this.frameRate, 
                      lossRateOrg: org,
                      lossRateFin: fin, 
                      delay: this.delayTimeIntArray,
                      odbw: this.origDetectedBW,
                      deviceType: this.deviceType,
                      networkType: this.networkType}
                      );
        }
    }`

	reducefunc := `  function(userId, info) {
        reducedVal = {
            userIp: info[0].userIp, 
            relayIp: info[0].relayIp, 
            meetingId: info[0].meetingId,
            deviceType: info[0].deviceType, 
            networkType: info[0].networkType,
			resourceId: info[0].resourceId,
            crArr: new Array(), 
            frArr: new Array(), 
            lossOrgArr: new Array(),
            loss: 0, 
            delayArr: new Array(), 
            delayLoss: 0,
            odbw: 0,
            dropline: false,
            deviceType: info[0].deviceType,
            networkType: info[0].networkType,
            };
        for (var idx = 0; idx < info.length; idx++) {
            if (info[idx].codeRate === undefined){
                Array.prototype.push.apply(reducedVal.crArr, info[idx].crArr);
                Array.prototype.push.apply(reducedVal.frArr, info[idx].frArr);
                Array.prototype.push.apply(reducedVal.lossOrgArr, info[idx].lossOrgArr);
                reducedVal.loss += info[idx].loss;
                reducedVal.delayLoss += info[idx].delayLoss;
                Array.prototype.push.apply(reducedVal.delayArr, info[idx].delayArr);
                reducedVal.odbw = info[idx].odbw;
                reducedVal.offline = info[idx].offline;
            }else{
                reducedVal.crArr.push(info[idx].codeRate);
                reducedVal.frArr.push(info[idx].frameRate);
                if (info[idx].lossRateFin > 0) {
                    reducedVal.loss += 1;
                }
                if (info[idx].lossRateOrg > 0) {
                    reducedVal.lossOrgArr.push(info[idx].lossRateOrg);
                }
                var ds = info[idx].delay.split('|');
                ds.forEach(function(d) {
                    if (d === "-1" || d === "") {
                        reducedVal.delayLoss += 1;
                    } else {
                        reducedVal.delayArr.push(parseInt(d));
                    }
                });
                if (info[idx].odbw !== -1 &&
                    reducedVal.odbw === 0){
                    reducedVal.odbw = info[idx].odbw
                }
                if (reducedVal.odbw !== 0 &&
                    info[idx].odbw === -1){
                    if (reducedVal.meetingId === info[idx].meetingId){
                        reducedVal.dropline = true;
                    }else{
                        reducedVal.meetingId = info[idx].meetingId
                        reducedVal.odbw = 0;
                    }
                }
            }
        }
        return reducedVal
    }`
	finalizefunc := `function(key, val){
           
		 if(val.loss == undefined){
			return {
                userip: "**",
                relayip: "",
                meetingid: "",
                crarr: "",
                frarr: "",
                lossorgarr: "",
                loss: 0, 
                delayarr: "",
                delayloss: 0,
                dropline: false,
                devicetype: 0,
                networktype: 0,

                }
		 }
            return {
                userip: val.userIp,
                relayip: val.relayIp,
                meetingid: val.meetingId,
                crarr: Array.join(val.crArr),
                frarr: Array.join(val.frArr),
                lossorgarr: Array.join(val.lossOrgArr),
                loss: val.loss, 
                delayarr: Array.join(val.delayArr),
                delayloss: val.delayLoss,
                dropline: val.dropline,
                devicetype: val.deviceType,
                networktype: val.networkType,
				resourceid: val.resourceId,

                }
        }`

	mapreduceup := &mgo.MapReduce{
		Map:      mapfunc,
		Reduce:   reducefunc,
		Finalize: finalizefunc,
	}

	var uservideosup []uservideoinfo

	_, err := conn.Find(bson.M{"timeStamp": bson.M{"$gte": mintime}}).MapReduce(mapreduceup, &uservideosup)
	if err != nil {
		fmt.Println(err, "1220")
		if len(uservideosup) == 0 {
			return
		}
	}

	relaymap := globeCfg.Output.Relaymap

	//文档上行
	for _, v := range uservideosup {

		if v.Id.Userid == "" {
			continue
		}
		if v.Value.Userip == "**" {
			continue
		}
		usernumber, err := strconv.ParseInt(v.Id.Userid, 10, 64)
		if err != nil {
			continue
		}

		var userent string
		var meetent string
		var usertype string
		userent = usermap[usernumber]
		meetent = userent
		if usermap[usernumber] == "非商业" || usermap[usernumber] == "" {
			userent = "非商业"
			isBusiness, usetid := isBusinessmeet(collectionmeet, int(v.Value.Meetingid))
			if !isBusiness {
				usertype = "非商业"
				userent = "非商业"
				meetent = "非商业"

			} else {
				usertype = "商业"
				meetent = usermap[usetid]
			}

		} else {
			usertype = "商业"
		}
		if strings.Index(usermap[usernumber], "(演示)") != -1 {
			usertype = "演示"
		}

		var relayDomIsp string
		var userDom string
		var userIsp string

		var meetid string
		meetid = strconv.Itoa(int(v.Value.Meetingid))
		relayid := strings.Split(v.Value.Relayip, ":")
		if relaymap[relayid[0]] != "" {
			relayDomIsp = relaymap[relayid[0]]
		} else {
			relayDomIsp = "unknown"
		}
		userid := strings.Split(v.Value.Userip, ":")
		loc, err := ip17mon.Find(userid[0])
		if err != nil {
			fmt.Println(err, "1282")
			return
		}
		if loc.City == "N/A" {
			if loc.Region == "N/A" {
				if loc.Country == "N/A" {
					userDom = "unknown"
				} else {
					userDom = loc.Country
				}
			} else {
				userDom = loc.Region
			}
		} else {
			userDom = loc.City
		}
		userIsp = loc.Isp

		userkeys := userkey{
			Userid:     v.Id.Userid,
			Meetintid:  meetid,
			Speakerid:  v.Id.Speakerid,
			Resourceid: v.Value.Resourceid,
			Relay:      relayDomIsp,
		}

		if userfilequalityup[userkeys].Id.Userid == "" {

			user := uservideoinfo{}
			user.Id = v.Id
			user.Value.Usertype = usertype
			user.Value.Meetingid = v.Value.Meetingid
			user.Value.Userent = userent
			user.Value.Meetent = meetent
			user.Value.Devicetype = v.Value.Devicetype
			user.Value.Networktype = v.Value.Networktype
			user.Value.Resourceid = v.Value.Resourceid
			user.Value.Userdom = userDom
			user.Value.Userisp = userIsp
			user.Value.Relaydomisp = relayDomIsp
			user.Value.Crarr = v.Value.Crarr
			user.Value.Frarr = v.Value.Frarr
			user.Value.Delayarr = v.Value.Delayarr
			user.Value.Delayloss = v.Value.Delayloss
			user.Value.Loss = v.Value.Loss
			user.Value.Lossorgarr = v.Value.Lossorgarr
			user.Value.Dropline = v.Value.Dropline
			v.Value.Users = 1
			userfilequalityup[userkeys] = user
		} else {
			user := userfilequalityup[userkeys]
			user.Value.Relaydomisp = relayDomIsp
			user.Value.Crarr += "," + v.Value.Crarr
			user.Value.Frarr += "," + v.Value.Frarr
			user.Value.Delayarr += "," + v.Value.Delayarr
			user.Value.Delayloss += v.Value.Delayloss
			user.Value.Loss += v.Value.Loss
			user.Value.Lossorgarr += "," + v.Value.Lossorgarr
			if v.Value.Dropline == true {
				user.Value.Dropline = true
			}
			v.Value.Users += 1
			userfilequalityup[userkeys] = user
		}

	}

}

//文档下行
func userfiledown(conn *mgo.Collection, collectionmeet *mgo.Collection, mintime int32) {

	mapfunc2 := `function() { 
        uid = parseInt(this.userId)
        if(this.userId != this.speakerId && 
              uid > 9999999 &&
			this.mediaType==2
			){
				var org,fin
			if(this.lossRateOriginal==0){
				org=0
			}else{
				org=this.lossRateOriginal/100.0
			}
			if(this.lossRateFinal==0){
				fin=0
			}else{
				fin=this.lossRateFinal/100.0
			}
            emit({userid:this.userId, speakerid:this.speakerId,resourceId:this.resourceId,relayIp: this.relayIp} , 
                  {userIp: this.userIp, 
                      relayIp: this.relayIp, 
                      meetingId: this.meetingId,
                      deviceType: this.deviceType,
                      networkType: this.networkType,
					  resourceId:this.resourceId, 
                      codeRate: this.curCodeRate, 
                      frameRate: this.frameRate, 
                      lossRateOrg: org,
                      lossRateFin: fin, 
                      delay: this.delayTimeIntArray,
                      odbw: this.origDetectedBW,
                      deviceType: this.deviceType,
                      networkType: this.networkType}
                      );
        }
    }`
	reducefunc := `  function(userId, info) {
        reducedVal = {
            userIp: info[0].userIp, 
            relayIp: info[0].relayIp, 
            meetingId: info[0].meetingId,
            deviceType: info[0].deviceType, 
            networkType: info[0].networkType, 
			resourceId: info[0].resourceId, 
            crArr: new Array(), 
            frArr: new Array(), 
            lossOrgArr: new Array(),
            loss: 0, 
            delayArr: new Array(), 
            delayLoss: 0,
            odbw: 0,
            dropline: false,
            deviceType: info[0].deviceType,
            networkType: info[0].networkType,
            };
        for (var idx = 0; idx < info.length; idx++) {
            if (info[idx].codeRate === undefined){
                Array.prototype.push.apply(reducedVal.crArr, info[idx].crArr);
                Array.prototype.push.apply(reducedVal.frArr, info[idx].frArr);
                Array.prototype.push.apply(reducedVal.lossOrgArr, info[idx].lossOrgArr);
                reducedVal.loss += info[idx].loss;
                reducedVal.delayLoss += info[idx].delayLoss;
                Array.prototype.push.apply(reducedVal.delayArr, info[idx].delayArr);
                reducedVal.odbw = info[idx].odbw;
                reducedVal.offline = info[idx].offline;
            }else{
                reducedVal.crArr.push(info[idx].codeRate);
                reducedVal.frArr.push(info[idx].frameRate);
                if (info[idx].lossRateFin > 0) {
                    reducedVal.loss += 1;
                }
                if (info[idx].lossRateOrg > 0) {
                    reducedVal.lossOrgArr.push(info[idx].lossRateOrg);
                }
                var ds = info[idx].delay.split('|');
                ds.forEach(function(d) {
                    if (d === "-1" || d === "") {
                        reducedVal.delayLoss += 1;
                    } else {
                        reducedVal.delayArr.push(parseInt(d));
                    }
                });
                if (info[idx].odbw !== -1 &&
                    reducedVal.odbw === 0){
                    reducedVal.odbw = info[idx].odbw
                }
                if (reducedVal.odbw !== 0 &&
                    info[idx].odbw === -1){
                    if (reducedVal.meetingId === info[idx].meetingId){
                        reducedVal.dropline = true;
                    }else{
                        reducedVal.meetingId = info[idx].meetingId
                        reducedVal.odbw = 0;
                    }
                }
            }
        }
        return reducedVal
    }`
	finalizefunc := `function(key, val){
           
		 if(val.loss == undefined){
			return {
                userip: "**",
                relayip: "",
                meetingid: "",
                crarr: "",
                frarr: "",
                lossorgarr: "",
                loss: 0, 
                delayarr: "",
                delayloss: 0,
                dropline: false,
                devicetype: 0,
                networktype: 0,

                }
		 }
            return {
                userip: val.userIp,
                relayip: val.relayIp,
                meetingid: val.meetingId,
                crarr: Array.join(val.crArr),
                frarr: Array.join(val.frArr),
                lossorgarr: Array.join(val.lossOrgArr),
                loss: val.loss, 
                delayarr: Array.join(val.delayArr),
                delayloss: val.delayLoss,
                dropline: val.dropline,
                devicetype: val.deviceType,
                networktype: val.networkType,
				resourceid: val.resourceId,
                }
        }`

	mapreducedown := &mgo.MapReduce{
		Map:      mapfunc2,
		Reduce:   reducefunc,
		Finalize: finalizefunc,
	}

	var uservideosdown []uservideoinfo

	_, err := conn.Find(bson.M{"timeStamp": bson.M{"$gte": mintime}}).MapReduce(mapreducedown, &uservideosdown)
	if err != nil {
		fmt.Println(err, "1494")
		if len(uservideosdown) == 0 {
			return
		}
	}
	relaymap := globeCfg.Output.Relaymap

	//文档下行
	for _, v := range uservideosdown {

		if v.Id.Userid == "" {
			continue
		}
		if v.Value.Userip == "**" {
			continue
		}
		usernumber, err := strconv.ParseInt(v.Id.Userid, 10, 64)
		if err != nil {
			continue
		}
		var userent string
		var meetent string
		var usertype string
		userent = usermap[usernumber]
		meetent = userent
		if usermap[usernumber] == "非商业" || usermap[usernumber] == "" {
			userent = "非商业"
			isBusiness, usetid := isBusinessmeet(collectionmeet, int(v.Value.Meetingid))
			if !isBusiness {
				usertype = "非商业"
				userent = "非商业"
				meetent = "非商业"

			} else {
				usertype = "商业"
				meetent = usermap[usetid]
			}

		} else {
			usertype = "商业"
		}
		if strings.Index(usermap[usernumber], "(演示)") != -1 {
			usertype = "演示"
		}
		var relayDomIsp string
		var userDom string
		var userIsp string

		var meetid string
		meetid = strconv.Itoa(int(v.Value.Meetingid))
		relayid := strings.Split(v.Value.Relayip, ":")
		if relaymap[relayid[0]] != "" {
			relayDomIsp = relaymap[relayid[0]]
		} else {
			relayDomIsp = "unknown"
		}
		userid := strings.Split(v.Value.Userip, ":")
		loc, err := ip17mon.Find(userid[0])
		if err != nil {
			fmt.Println(err, "162")
			return
		}
		if loc.City == "N/A" {
			if loc.Region == "N/A" {
				if loc.Country == "N/A" {
					userDom = "unknown"
				} else {
					userDom = loc.Country
				}
			} else {
				userDom = loc.Region
			}
		} else {
			userDom = loc.City
		}
		userIsp = loc.Isp

		if relaymap[relayid[0]] != "" {
			relayDomIsp = relaymap[relayid[0]]
		} else {
			relayDomIsp = "unknown"
		}

		userkeys := userkey{
			Userid:     v.Id.Userid,
			Meetintid:  meetid,
			Speakerid:  v.Id.Speakerid,
			Resourceid: v.Value.Resourceid,
			Relay:      relayDomIsp,
		}

		if userfilequalitydown[userkeys].Id.Userid == "" {

			user := uservideoinfo{}
			user.Id = v.Id
			user.Value.Usertype = usertype
			user.Value.Meetingid = v.Value.Meetingid
			user.Value.Userent = userent
			user.Value.Meetent = meetent
			user.Value.Devicetype = v.Value.Devicetype
			user.Value.Networktype = v.Value.Networktype
			user.Value.Resourceid = v.Value.Resourceid
			user.Value.Userdom = userDom
			user.Value.Userisp = userIsp
			user.Value.Relaydomisp = relayDomIsp
			user.Value.Crarr = v.Value.Crarr
			user.Value.Frarr = v.Value.Frarr
			user.Value.Delayarr = v.Value.Delayarr
			user.Value.Delayloss = v.Value.Delayloss
			user.Value.Loss = v.Value.Loss
			user.Value.Lossorgarr = v.Value.Lossorgarr
			user.Value.Dropline = v.Value.Dropline
			v.Value.Users = 1
			userfilequalitydown[userkeys] = user
		} else {
			user := userfilequalitydown[userkeys]
			user.Value.Relaydomisp = relayDomIsp
			user.Value.Crarr += "," + v.Value.Crarr
			user.Value.Frarr += "," + v.Value.Frarr
			user.Value.Delayarr += "," + v.Value.Delayarr
			user.Value.Delayloss += v.Value.Delayloss
			user.Value.Loss += v.Value.Loss
			user.Value.Lossorgarr += "," + v.Value.Lossorgarr
			if v.Value.Dropline == true {
				user.Value.Dropline = true
			}
			v.Value.Users += 1
			userfilequalitydown[userkeys] = user
		}

	}

}
func userVideoAudioUpAndDown(conn *mgo.Collection, collectionmeet *mgo.Collection, mintime int32) {
	uservideoup(conn, collectionmeet, mintime)
	uservideodown(conn, collectionmeet, mintime)
	useraudioup(conn, collectionmeet, mintime)
	useraudiodown(conn, collectionmeet, mintime)
	userfileup(conn, collectionmeet, mintime)
	userfiledown(conn, collectionmeet, mintime)
}
