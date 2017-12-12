package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/wangtuanjie/ip17mon"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func uservideo(conn *mgo.Collection, mintime int32) {

	mapfunc := `function() { 
        uid = parseInt(this.userId)
        if(this.userId === this.speakerId && 
              uid > 9999999 &&
			((this.resourceId==101 && this.deviceType==5)||(this.resourceId==101 && this.deviceType==6)||(this.resourceId==100))
			){
            emit(this.userId, 
                  {userIp: this.userIp, 
                      relayIp: this.relayIp, 
                      meetingId: this.meetingId,
                      deviceType: this.deviceType,
                      networkType: this.networkType, 
                      codeRate: this.curCodeRate, 
                      frameRate: this.frameRate, 
                      lossRateOrg: this.lossRateOriginal,
                      lossRateFin: this.lossRateFinal, 
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
            lo = 0;
            loa = 0;
            los = 0;
            if (val.lossOrgArr.length > 0) {
                lo = val.lossOrgArr.length;
                loa = Array.avg(val.lossOrgArr);
                los = Array.stdDev(val.lossOrgArr);
            }
            return {
                userIp: val.userIp, 
                relayIp: val.relayIp, 
                meetingId: val.meetingId,
                deviceType: val.deviceType,
                networkType: val.networkType,  
                crAvg: Array.avg(val.crArr), 
                crStd: Array.stdDev(val.crArr), 
                frAvg: Array.avg(val.frArr), 
                frStd: Array.stdDev(val.frArr), 
                lossOrg: lo,
                lossOrgAvg: loa,
                lossOrgStd: los,
                loss: val.loss, 
                delayAvg: Array.avg(val.delayArr), 
                delayStd: Array.stdDev(val.delayArr), 
                delayLoss: val.delayLoss,
                dropline: val.dropline,
                deviceType: val.deviceType,
                networkType: val.networkType,
                }
        }`
	mapreduce := &mgo.MapReduce{
		Map:      mapfunc,
		Reduce:   reducefunc,
		Finalize: finalizefunc,
	}
	var uservideos []uservideoinfo
	_, err := conn.Find(bson.M{"timeStamp": bson.M{"$gte": mintime}}).MapReduce(mapreduce, &uservideos)
	if err != nil {
		fmt.Println(err, "137")
		return
	}
	relaymap := globeCfg.Output.Relaymap

	for _, v := range uservideos {
		usernumber, err := strconv.ParseInt(v.Id, 10, 64)
		if err != nil {
			continue
		}
		if usermap[usernumber] == "非商业" || usermap[usernumber] == "" {
			continue
		}
		var relayDomIsp string
		var userDom string
		var userIsp string
		var ent string
		var meetid string
		meetid = strconv.Itoa(v.Value.meetingId)
		relayid := strings.Split(v.Value.relayIp, ":")
		userid := strings.Split(v.Value.userIp, ":")
		loc, err := ip17mon.Find(userid[0])
		if err != nil {
			fmt.Println(err, "144")
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
		userIsp += loc.Isp
		ent = usermap[usernumber]
		relayDomIsp = relaymap[relayid[0]]
		crAvg.WithLabelValues(v.Id, meetid, ent, ent, userDom, userIsp, relayDomIsp).Set(float64(v.Value.crAvg))
		crStd.WithLabelValues(v.Id, meetid, ent, ent, userDom, userIsp, relayDomIsp).Set(float64(v.Value.crStd))
		frAvg.WithLabelValues(v.Id, meetid, ent, ent, userDom, userIsp, relayDomIsp).Set(float64(v.Value.frAvg))
		frStd.WithLabelValues(v.Id, meetid, ent, ent, userDom, userIsp, relayDomIsp).Set(float64(v.Value.frStd))
		delayAvg.WithLabelValues(v.Id, meetid, ent, ent, userDom, userIsp, relayDomIsp).Set(float64(v.Value.delayAvg))
		delayStd.WithLabelValues(v.Id, meetid, ent, ent, userDom, userIsp, relayDomIsp).Set(float64(v.Value.delayStd))
		delayLoss.WithLabelValues(v.Id, meetid, ent, ent, userDom, userIsp, relayDomIsp).Set(float64(v.Value.delayLoss))
		loss.WithLabelValues(v.Id, meetid, ent, ent, userDom, userIsp, relayDomIsp).Set(float64(v.Value.loss))
		lossOrg.WithLabelValues(v.Id, meetid, ent, ent, userDom, userIsp, relayDomIsp).Set(float64(v.Value.lossOrg))
		lossOrgAvg.WithLabelValues(v.Id, meetid, ent, ent, userDom, userIsp, relayDomIsp).Set(float64(v.Value.lossOrgAvg))
		lossOrgStd.WithLabelValues(v.Id, meetid, ent, ent, userDom, userIsp, relayDomIsp).Set(float64(v.Value.lossOrgStd))
		if v.Value.dropline == true {
			dropline.WithLabelValues(v.Id, meetid, ent, ent, userDom, userIsp, relayDomIsp).Set(float64(1))

		}
	}

}
