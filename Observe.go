package main

import (
	"fmt"
	"strconv"
	"strings"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/montanaflynn/stats"
)

//推送数据
func Observe() {
	//当前商业会议在线人数
	synode.WithLabelValues("nowuser").Set(float64(syuser["now"]))
	//当前商业会议新增人数
	synode.WithLabelValues("newuser").Set(float64(syuser["new"]))
	//当前参加商业会议的人数
	synode.WithLabelValues("realuser").Set(float64(syuser["realuser"]))
	//当前参加会议得人数
	synode.WithLabelValues("alluser").Set(float64(syuser["alluser"]))
	syuser["now"] = 0
	syuser["new"] = 0
	syuser["realuser"] = 0
	syuser["alluser"] = 0
	//当前企业通话数
	for k, v := range entToNumber1 {
		nodenow.WithLabelValues(k, "Now").Set(float64(v))
		entToNumber1[k] = 0
	}
	//新增企业通话数
	for k, v := range entToNumber2 {
		nodenow.WithLabelValues(k, "New").Set(float64(v))
		entToNumber2[k] = 0
	}
	//企业未接通数
	for k, v := range entToNumber3 {
		nodenow.WithLabelValues(k, "CallFail").Set(float64(v))
		entToNumber3[k] = 0
	}
	//企业延迟插入的通话数（时间与实际时间相差5m及以上）

	for k, v := range entToNumber4 {
		nodenow.WithLabelValues(k, "Glean").Set(float64(v))
		entToNumber4[k] = 0
	}
	//企业通话总时长
	for k, v := range entToNumber5 {
		nodenowcalltime.WithLabelValues(k, "CallTime").Set(float64(v))
		entToNumber5[k] = 0
	}
	//按地域划分的通话数（p2p）
	for k, v := range domToNumber {
		if dommap[k] == "" {
			nodecallip.WithLabelValues("unknown", "dom").Set(float64(v))
		} else {
			nodecallip.WithLabelValues(dommap[k], "dom").Set(float64(v))
		}
		domToNumber[k] = 0
	}
	//按运营商划分的通话数（p2p）
	for k, v := range ispToNumber {
		if ispmap[k] == "" {
			nodecallip.WithLabelValues("unknown", "isp").Set(float64(v))
		} else {
			nodecallip.WithLabelValues(ispmap[k], "isp").Set(float64(v))
		}
		ispToNumber[k] = 0
	}
	fmt.Println("devp2p", devToNumber)
	//按设备类型划分的通话数（p2p）
	for k, v := range devToNumber {
		if devp2pmap[k] == "" {
			nodecallip.WithLabelValues("unknown", "dev").Set(float64(v))
		} else {
			nodecallip.WithLabelValues(devp2pmap[k], "dev").Set(float64(v))
		}
		devToNumber[k] = 0
	}
	//按网络类型划分的通话数（p2p）
	for k, v := range netToNumber {
		if netp2pmap[k] == "" {
			nodecallip.WithLabelValues("unknown", "net").Set(float64(v))
		} else {
			nodecallip.WithLabelValues(netp2pmap[k], "net").Set(float64(v))
		}

		netToNumber[k] = 0
	}
	//企业通话时长（meet）
	fmt.Println("entTocallTime: ", entTocallTime)
	for k, v := range entTocallTime {
		entcalltimemeet.WithLabelValues(k).Set(float64(v))
		entTocallTime[k] = 0
	}
	//按地域划分的通话数（meet）
	fmt.Println("domToNumbermeet: ", domToNumbermeet)
	for k, v := range domToNumbermeet {
		entipmeet.WithLabelValues(k, "dom").Set(float64(v))
		domToNumbermeet[k] = 0
	}
	//按运营商划分的通话数（meet）
	fmt.Println("ispToNumbermeet: ", ispToNumbermeet)
	for k, v := range ispToNumbermeet {
		entipmeet.WithLabelValues(k, "isp").Set(float64(v))
		ispToNumbermeet[k] = 0
	}
	//按设备类型划分的通话数（meet）
	fmt.Println("devToNumbermeet: ", devToNumbermeet)
	for k, v := range devToNumbermeet {
		entipmeet.WithLabelValues(k, "dev").Set(float64(v))
		devToNumbermeet[k] = 0
	}
	//按网络类型划分的通话数（meet）
	fmt.Println("netToNumbermeet: ", netToNumbermeet)
	for k, v := range netToNumbermeet {
		entipmeet.WithLabelValues(k, "net").Set(float64(v))
		netToNumbermeet[k] = 0
	}

	for k, v := range meetrelayup {
		meetrelay.WithLabelValues(k, "up").Set(float64(v))
		meetrelayup[k] = 0
	}
	for k, v := range meetrelaydown {
		meetrelay.WithLabelValues(k, "down").Set(float64(v))
		meetrelaydown[k] = 0
	}
	//p2p用户活跃统计
	fmt.Println("newusersp2p: ", newusersp2p)
	for k, v := range newusersp2p {
		p2pactive.WithLabelValues(k, "new").Set(float64(v))
		newusersp2p[k] = 0
	}
	for k, v := range monOldusersp2p {
		p2pactive.WithLabelValues(k, "mon").Set(float64(v))
		monOldusersp2p[k] = 0
	}

	for k, v := range weekOldusersp2p {
		p2pactive.WithLabelValues(k, "week").Set(float64(v))
		weekOldusersp2p[k] = 0
	}
	for k, v := range dayOldusersp2p {
		p2pactive.WithLabelValues(k, "day").Set(float64(v))
		dayOldusersp2p[k] = 0
	}
	//meet用户活跃统计
	fmt.Println("newusersmeet: ", newusersmeet)
	for k, v := range newusersmeet {
		meetactive.WithLabelValues(k, "new").Set(float64(v))
		newusersmeet[k] = 0
	}
	for k, v := range monOldusersmeet {
		meetactive.WithLabelValues(k, "mon").Set(float64(v))
		monOldusersmeet[k] = 0
	}

	for k, v := range weekOldusersmeet {
		meetactive.WithLabelValues(k, "week").Set(float64(v))
		weekOldusersmeet[k] = 0
	}
	for k, v := range dayOldusersmeet {
		meetactive.WithLabelValues(k, "day").Set(float64(v))
		dayOldusersmeet[k] = 0
	}
	//文档上行
	for k, v := range userfilequalityup {

		if v.Value.Crarr != "" && v.Value.Frarr != "" && v.Value.Delayarr != "" {
			var dev, net string
			CrfloatArr := strArrToFloatArr(strings.Split(v.Value.Crarr, ","))
			FrfloatArr := strArrToFloatArr(strings.Split(v.Value.Frarr, ","))
			delayfloatArr := strArrToFloatArr(strings.Split(v.Value.Delayarr, ","))
			lossfloatArr := strArrToFloatArr(strings.Split(v.Value.Lossorgarr, ","))
			Devicetypestr := strconv.Itoa(v.Value.Devicetype)
			Networktypestr := strconv.Itoa(v.Value.Networktype)
			if devmeet[Devicetypestr] != "" {
				dev = devmeet[Devicetypestr]
				devstreamkey := Devicetypestr + "|" + strconv.Itoa(v.Value.Resourceid)
				if devstream[devstreamkey] != "" {
					dev = devstream[devstreamkey]
				}

			} else {
				dev = "unknown"
			}
			if netmeet[Networktypestr] != "" {
				net = netmeet[Networktypestr]
			} else {
				net = "unknown"
			}
			//上行最大延迟
			maxdelayup, _ = stats.Max(delayfloatArr)

			meetidstr := strconv.FormatInt(int64(v.Value.Meetingid), 10)
			crAvgValue, _ := stats.Mean(CrfloatArr)
			crStdValue, _ := stats.StandardDeviationPopulation(CrfloatArr)
			frAvgValue, _ := stats.Mean(FrfloatArr)
			frStdValue, _ := stats.StandardDeviationPopulation(FrfloatArr)
			delayAvgValue, _ := stats.Mean(delayfloatArr)
			delayStdValue, _ := stats.StandardDeviationPopulation(delayfloatArr)
			lossAvgValue, _ := stats.Mean(lossfloatArr)
			lossStdValue, _ := stats.StandardDeviationPopulation(lossfloatArr)

			var lossorg float64
			for _, v := range lossfloatArr {
				if v != 0 {
					lossorg++
				}
			}
			//			fmt.Println("crAvgValue", crAvgValue, "|", v.Value.Crarr)
			//			fmt.Println("frAvgValue", frAvgValue, "|", v.Value.Frarr)
			//			fmt.Println("delayAvgValue", delayAvgValue, "|", v.Value.Delayarr)
			//			fmt.Println("lossAvgValue", lossAvgValue, "|", v.Value.Lossorgarr, "|", lossfloatArr)
			//			fmt.Println("lossorg", lossorg)
			filecrAvg.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "up", v.Value.Usertype).Set(float64(int64(crAvgValue)))
			filecrStd.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "up", v.Value.Usertype).Set(float64(int64(crStdValue)))
			filefrAvg.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "up", v.Value.Usertype).Set(float64(int64(frAvgValue)))
			filefrStd.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "up", v.Value.Usertype).Set(float64(int64(frStdValue)))
			filedelayAvg.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "up", v.Value.Usertype).Set(float64(int64(delayAvgValue)))
			filedelayStd.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "up", v.Value.Usertype).Set(float64(int64(delayStdValue)))
			filedelayLoss.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "up", v.Value.Usertype).Set(float64(int64(v.Value.Delayloss)))
			fileloss.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "up", v.Value.Usertype).Set(float64(int64(v.Value.Loss)))
			filelossOrg.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "up", v.Value.Usertype).Set(float64(int64(lossorg)))
			filelossOrgAvg.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "up", v.Value.Usertype).Set(float64(int64(lossAvgValue)))
			filelossOrgStd.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "up", v.Value.Usertype).Set(float64(int64(lossStdValue)))
			filevideomaxdelay.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "up", v.Value.Usertype).Set(maxdelayup)
			if v.Value.Dropline == true {
				filedropline.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "up", v.Value.Usertype).Set(float64(1))

			}
		} else {
			var dev, net string
			Devicetypestr := strconv.Itoa(v.Value.Devicetype)
			Networktypestr := strconv.Itoa(v.Value.Networktype)
			if devmeet[Devicetypestr] != "" {
				dev = devmeet[Devicetypestr]
				devstreamkey := Devicetypestr + "|" + strconv.Itoa(v.Value.Resourceid)
				if devstream[devstreamkey] != "" {
					dev = devstream[devstreamkey]
				}
			} else {
				dev = "unknown"
			}
			if netmeet[Networktypestr] != "" {
				net = netmeet[Networktypestr]
			} else {
				net = "unknown"
			}
			meetidstr := strconv.FormatInt(int64(v.Value.Meetingid), 10)
			filecrAvg.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "up", v.Value.Usertype).Set(float64(0))
			filecrStd.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "up", v.Value.Usertype).Set(float64(0))
			filefrAvg.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "up", v.Value.Usertype).Set(float64(0))
			filefrStd.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "up", v.Value.Usertype).Set(float64(0))
			filedelayAvg.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "up", v.Value.Usertype).Set(float64(0))
			filedelayStd.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "up", v.Value.Usertype).Set(float64(0))
			filedelayLoss.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "up", v.Value.Usertype).Set(float64(0))
			fileloss.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "up", v.Value.Usertype).Set(float64(0))
			filelossOrg.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "up", v.Value.Usertype).Set(float64(0))
			filelossOrgAvg.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "up", v.Value.Usertype).Set(float64(0))
			filelossOrgStd.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "up", v.Value.Usertype).Set(float64(0))
			filedropline.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "up", v.Value.Usertype).Set(float64(0))
			filevideomaxdelay.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "up", v.Value.Usertype).Set(float64(0))
			delete(userfilequalityup, k)
		}

		user := userfilequalityup[k]
		user.Value.Crarr = ""
		user.Value.Frarr = ""
		user.Value.Delayarr = ""
		user.Value.Delayloss = 0
		user.Value.Loss = 0
		user.Value.Lossorgarr = ""
		user.Value.Dropline = false
		user.Value.Users = 0
		userfilequalityup[k] = user
	}
	//文档下行
	for k, v := range userfilequalitydown {

		if v.Value.Crarr != "" && v.Value.Frarr != "" && v.Value.Delayarr != "" {
			var dev, net string
			CrfloatArr := strArrToFloatArr(strings.Split(v.Value.Crarr, ","))
			FrfloatArr := strArrToFloatArr(strings.Split(v.Value.Frarr, ","))
			delayfloatArr := strArrToFloatArr(strings.Split(v.Value.Delayarr, ","))
			lossfloatArr := strArrToFloatArr(strings.Split(v.Value.Lossorgarr, ","))
			Devicetypestr := strconv.Itoa(v.Value.Devicetype)
			Networktypestr := strconv.Itoa(v.Value.Networktype)
			if devmeet[Devicetypestr] != "" {
				dev = devmeet[Devicetypestr]
				devstreamkey := Devicetypestr + "|" + strconv.Itoa(v.Value.Resourceid)
				if devstream[devstreamkey] != "" {
					dev = devstream[devstreamkey]
				}

			} else {
				dev = "unknown"
			}
			if netmeet[Networktypestr] != "" {
				net = netmeet[Networktypestr]
			} else {
				net = "unknown"
			}
			maxdelaydown, _ = stats.Max(delayfloatArr)
			meetidstr := strconv.FormatInt(int64(v.Value.Meetingid), 10)
			crAvgValue, _ := stats.Mean(CrfloatArr)
			crStdValue, _ := stats.StandardDeviationPopulation(CrfloatArr)
			frAvgValue, _ := stats.Mean(FrfloatArr)
			frStdValue, _ := stats.StandardDeviationPopulation(FrfloatArr)
			delayAvgValue, _ := stats.Mean(delayfloatArr)
			delayStdValue, _ := stats.StandardDeviationPopulation(delayfloatArr)
			lossAvgValue, _ := stats.Mean(lossfloatArr)
			lossStdValue, _ := stats.StandardDeviationPopulation(lossfloatArr)

			var lossorg float64

			for _, v := range lossfloatArr {
				if v != 0 {
					lossorg++
				}
			}
			//			fmt.Println("crAvgValue", crAvgValue, "|", v.Value.Crarr)
			//			fmt.Println("frAvgValue", frAvgValue, "|", v.Value.Frarr)
			//			fmt.Println("delayAvgValue", delayAvgValue, "|", v.Value.Delayarr)
			//			fmt.Println("lossAvgValue", lossAvgValue, "|", v.Value.Lossorgarr, "|", lossfloatArr)
			//			fmt.Println("lossorg", lossorg)
			filecrAvg.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "down", v.Value.Usertype).Set(float64(int64(crAvgValue)))
			filecrStd.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "down", v.Value.Usertype).Set(float64(int64(crStdValue)))
			filefrAvg.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "down", v.Value.Usertype).Set(float64(int64(frAvgValue)))
			filefrStd.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "down", v.Value.Usertype).Set(float64(int64(frStdValue)))
			filedelayAvg.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "down", v.Value.Usertype).Set(float64(int64(delayAvgValue)))
			filedelayStd.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "down", v.Value.Usertype).Set(float64(int64(delayStdValue)))
			filedelayLoss.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "down", v.Value.Usertype).Set(float64(int64(v.Value.Delayloss)))
			fileloss.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "down", v.Value.Usertype).Set(float64(int64(v.Value.Loss)))
			filelossOrg.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "down", v.Value.Usertype).Set(float64(int64(lossorg)))
			filelossOrgAvg.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "down", v.Value.Usertype).Set(float64(int64(lossAvgValue)))
			filelossOrgStd.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "down", v.Value.Usertype).Set(float64(int64(lossStdValue)))
			filevideomaxdelay.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "down", v.Value.Usertype).Set(float64(maxdelaydown))
			if v.Value.Dropline == true {
				filedropline.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "down", v.Value.Usertype).Set(float64(1))

			}
		} else {
			var dev, net string
			Devicetypestr := strconv.Itoa(v.Value.Devicetype)
			Networktypestr := strconv.Itoa(v.Value.Networktype)
			if devmeet[Devicetypestr] != "" {
				dev = devmeet[Devicetypestr]
				devstreamkey := Devicetypestr + "|" + strconv.Itoa(v.Value.Resourceid)
				if devstream[devstreamkey] != "" {
					dev = devstream[devstreamkey]
				}
			} else {
				dev = "unknown"
			}
			if netmeet[Networktypestr] != "" {
				net = netmeet[Networktypestr]
			} else {
				net = "unknown"
			}
			meetidstr := strconv.FormatInt(int64(v.Value.Meetingid), 10)
			filecrAvg.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "down", v.Value.Usertype).Set(float64(0))
			filecrStd.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "down", v.Value.Usertype).Set(float64(0))
			filefrAvg.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "down", v.Value.Usertype).Set(float64(0))
			filefrStd.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "down", v.Value.Usertype).Set(float64(0))
			filedelayAvg.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "down", v.Value.Usertype).Set(float64(0))
			filedelayStd.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "down", v.Value.Usertype).Set(float64(0))
			filedelayLoss.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "down", v.Value.Usertype).Set(float64(0))
			fileloss.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "down", v.Value.Usertype).Set(float64(0))
			filelossOrg.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "down", v.Value.Usertype).Set(float64(0))
			filelossOrgAvg.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "down", v.Value.Usertype).Set(float64(0))
			filelossOrgStd.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "down", v.Value.Usertype).Set(float64(0))
			filedropline.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "down", v.Value.Usertype).Set(float64(0))
			filevideomaxdelay.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "down", v.Value.Usertype).Set(float64(0))
			delete(userfilequalitydown, k)
		}

		user := userfilequalitydown[k]
		user.Value.Crarr = ""
		user.Value.Frarr = ""
		user.Value.Delayarr = ""
		user.Value.Delayloss = 0
		user.Value.Loss = 0
		user.Value.Lossorgarr = ""
		user.Value.Dropline = false
		user.Value.Users = 0
		userfilequalitydown[k] = user
	}
	//视频上行
	for k, v := range uservideoqualityup {

		if v.Value.Crarr != "" && v.Value.Frarr != "" && v.Value.Delayarr != "" {
			var dev, net string
			CrfloatArr := strArrToFloatArr(strings.Split(v.Value.Crarr, ","))
			FrfloatArr := strArrToFloatArr(strings.Split(v.Value.Frarr, ","))
			delayfloatArr := strArrToFloatArr(strings.Split(v.Value.Delayarr, ","))
			lossfloatArr := strArrToFloatArr(strings.Split(v.Value.Lossorgarr, ","))
			Devicetypestr := strconv.Itoa(v.Value.Devicetype)
			Networktypestr := strconv.Itoa(v.Value.Networktype)
			if devmeet[Devicetypestr] != "" {
				dev = devmeet[Devicetypestr]
				devstreamkey := Devicetypestr + "|" + strconv.Itoa(v.Value.Resourceid)
				if devstream[devstreamkey] != "" {
					dev = devstream[devstreamkey]
				}

			} else {
				dev = "unknown"
			}
			if netmeet[Networktypestr] != "" {
				net = netmeet[Networktypestr]
			} else {
				net = "unknown"
			}
			//上行最大延迟
			maxdelayup, _ = stats.Max(delayfloatArr)

			meetidstr := strconv.FormatInt(int64(v.Value.Meetingid), 10)
			crAvgValue, _ := stats.Mean(CrfloatArr)
			crStdValue, _ := stats.StandardDeviationPopulation(CrfloatArr)
			frAvgValue, _ := stats.Mean(FrfloatArr)
			frStdValue, _ := stats.StandardDeviationPopulation(FrfloatArr)
			delayAvgValue, _ := stats.Mean(delayfloatArr)
			delayStdValue, _ := stats.StandardDeviationPopulation(delayfloatArr)
			lossAvgValue, _ := stats.Mean(lossfloatArr)
			lossStdValue, _ := stats.StandardDeviationPopulation(lossfloatArr)

			var lossorg float64
			for _, v := range lossfloatArr {
				if v != 0 {
					lossorg++
				}
			}
			//			fmt.Println("crAvgValue", crAvgValue, "|", v.Value.Crarr)
			//			fmt.Println("frAvgValue", frAvgValue, "|", v.Value.Frarr)
			//			fmt.Println("delayAvgValue", delayAvgValue, "|", v.Value.Delayarr)
			//			fmt.Println("lossAvgValue", lossAvgValue, "|", v.Value.Lossorgarr, "|", lossfloatArr)
			//			fmt.Println("lossorg", lossorg)
			crAvg.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "up", v.Value.Usertype).Set(float64(int64(crAvgValue)))
			crStd.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "up", v.Value.Usertype).Set(float64(int64(crStdValue)))
			frAvg.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "up", v.Value.Usertype).Set(float64(int64(frAvgValue)))
			frStd.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "up", v.Value.Usertype).Set(float64(int64(frStdValue)))
			delayAvg.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "up", v.Value.Usertype).Set(float64(int64(delayAvgValue)))
			delayStd.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "up", v.Value.Usertype).Set(float64(int64(delayStdValue)))
			delayLoss.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "up", v.Value.Usertype).Set(float64(int64(v.Value.Delayloss)))
			loss.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "up", v.Value.Usertype).Set(float64(int64(v.Value.Loss)))
			lossOrg.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "up", v.Value.Usertype).Set(float64(int64(lossorg)))
			lossOrgAvg.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "up", v.Value.Usertype).Set(float64(int64(lossAvgValue)))
			lossOrgStd.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "up", v.Value.Usertype).Set(float64(int64(lossStdValue)))
			videomaxdelay.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "up", v.Value.Usertype).Set(maxdelayup)
			if v.Value.Dropline == true {
				dropline.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "up", v.Value.Usertype).Set(float64(1))

			}
		} else {
			var dev, net string
			Devicetypestr := strconv.Itoa(v.Value.Devicetype)
			Networktypestr := strconv.Itoa(v.Value.Networktype)
			if devmeet[Devicetypestr] != "" {
				dev = devmeet[Devicetypestr]
				devstreamkey := Devicetypestr + "|" + strconv.Itoa(v.Value.Resourceid)
				if devstream[devstreamkey] != "" {
					dev = devstream[devstreamkey]
				}
			} else {
				dev = "unknown"
			}
			if netmeet[Networktypestr] != "" {
				net = netmeet[Networktypestr]
			} else {
				net = "unknown"
			}
			meetidstr := strconv.FormatInt(int64(v.Value.Meetingid), 10)
			crAvg.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "up", v.Value.Usertype).Set(float64(0))
			crStd.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "up", v.Value.Usertype).Set(float64(0))
			frAvg.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "up", v.Value.Usertype).Set(float64(0))
			frStd.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "up", v.Value.Usertype).Set(float64(0))
			delayAvg.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "up", v.Value.Usertype).Set(float64(0))
			delayStd.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "up", v.Value.Usertype).Set(float64(0))
			delayLoss.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "up", v.Value.Usertype).Set(float64(0))
			loss.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "up", v.Value.Usertype).Set(float64(0))
			lossOrg.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "up", v.Value.Usertype).Set(float64(0))
			lossOrgAvg.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "up", v.Value.Usertype).Set(float64(0))
			lossOrgStd.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "up", v.Value.Usertype).Set(float64(0))
			dropline.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "up", v.Value.Usertype).Set(float64(0))
			videomaxdelay.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "up", v.Value.Usertype).Set(float64(0))
			delete(uservideoqualityup, k)
		}

		user := uservideoqualityup[k]
		user.Value.Crarr = ""
		user.Value.Frarr = ""
		user.Value.Delayarr = ""
		user.Value.Delayloss = 0
		user.Value.Loss = 0
		user.Value.Lossorgarr = ""
		user.Value.Dropline = false
		user.Value.Users = 0
		uservideoqualityup[k] = user
	}
	//视频下行
	for k, v := range uservideoqualitydown {

		if v.Value.Crarr != "" && v.Value.Frarr != "" && v.Value.Delayarr != "" {
			var dev, net string
			CrfloatArr := strArrToFloatArr(strings.Split(v.Value.Crarr, ","))
			FrfloatArr := strArrToFloatArr(strings.Split(v.Value.Frarr, ","))
			delayfloatArr := strArrToFloatArr(strings.Split(v.Value.Delayarr, ","))
			lossfloatArr := strArrToFloatArr(strings.Split(v.Value.Lossorgarr, ","))
			Devicetypestr := strconv.Itoa(v.Value.Devicetype)
			Networktypestr := strconv.Itoa(v.Value.Networktype)
			if devmeet[Devicetypestr] != "" {
				dev = devmeet[Devicetypestr]
				devstreamkey := Devicetypestr + "|" + strconv.Itoa(v.Value.Resourceid)
				if devstream[devstreamkey] != "" {
					dev = devstream[devstreamkey]
				}

			} else {
				dev = "unknown"
			}
			if netmeet[Networktypestr] != "" {
				net = netmeet[Networktypestr]
			} else {
				net = "unknown"
			}
			maxdelaydown, _ = stats.Max(delayfloatArr)
			meetidstr := strconv.FormatInt(int64(v.Value.Meetingid), 10)
			crAvgValue, _ := stats.Mean(CrfloatArr)
			crStdValue, _ := stats.StandardDeviationPopulation(CrfloatArr)
			frAvgValue, _ := stats.Mean(FrfloatArr)
			frStdValue, _ := stats.StandardDeviationPopulation(FrfloatArr)
			delayAvgValue, _ := stats.Mean(delayfloatArr)
			delayStdValue, _ := stats.StandardDeviationPopulation(delayfloatArr)
			lossAvgValue, _ := stats.Mean(lossfloatArr)
			lossStdValue, _ := stats.StandardDeviationPopulation(lossfloatArr)

			var lossorg float64

			for _, v := range lossfloatArr {
				if v != 0 {
					lossorg++
				}
			}
			//			fmt.Println("crAvgValue", crAvgValue, "|", v.Value.Crarr)
			//			fmt.Println("frAvgValue", frAvgValue, "|", v.Value.Frarr)
			//			fmt.Println("delayAvgValue", delayAvgValue, "|", v.Value.Delayarr)
			//			fmt.Println("lossAvgValue", lossAvgValue, "|", v.Value.Lossorgarr, "|", lossfloatArr)
			//			fmt.Println("lossorg", lossorg)
			crAvg.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "down", v.Value.Usertype).Set(float64(int64(crAvgValue)))
			crStd.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "down", v.Value.Usertype).Set(float64(int64(crStdValue)))
			frAvg.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "down", v.Value.Usertype).Set(float64(int64(frAvgValue)))
			frStd.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "down", v.Value.Usertype).Set(float64(int64(frStdValue)))
			delayAvg.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "down", v.Value.Usertype).Set(float64(int64(delayAvgValue)))
			delayStd.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "down", v.Value.Usertype).Set(float64(int64(delayStdValue)))
			delayLoss.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "down", v.Value.Usertype).Set(float64(int64(v.Value.Delayloss)))
			loss.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "down", v.Value.Usertype).Set(float64(int64(v.Value.Loss)))
			lossOrg.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "down", v.Value.Usertype).Set(float64(int64(lossorg)))
			lossOrgAvg.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "down", v.Value.Usertype).Set(float64(int64(lossAvgValue)))
			lossOrgStd.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "down", v.Value.Usertype).Set(float64(int64(lossStdValue)))
			videomaxdelay.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "down", v.Value.Usertype).Set(float64(maxdelaydown))
			if v.Value.Dropline == true {
				dropline.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "down", v.Value.Usertype).Set(float64(1))

			}
		} else {
			var dev, net string
			Devicetypestr := strconv.Itoa(v.Value.Devicetype)
			Networktypestr := strconv.Itoa(v.Value.Networktype)
			if devmeet[Devicetypestr] != "" {
				dev = devmeet[Devicetypestr]
				devstreamkey := Devicetypestr + "|" + strconv.Itoa(v.Value.Resourceid)
				if devstream[devstreamkey] != "" {
					dev = devstream[devstreamkey]
				}
			} else {
				dev = "unknown"
			}
			if netmeet[Networktypestr] != "" {
				net = netmeet[Networktypestr]
			} else {
				net = "unknown"
			}
			meetidstr := strconv.FormatInt(int64(v.Value.Meetingid), 10)
			crAvg.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "down", v.Value.Usertype).Set(float64(0))
			crStd.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "down", v.Value.Usertype).Set(float64(0))
			frAvg.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "down", v.Value.Usertype).Set(float64(0))
			frStd.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "down", v.Value.Usertype).Set(float64(0))
			delayAvg.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "down", v.Value.Usertype).Set(float64(0))
			delayStd.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "down", v.Value.Usertype).Set(float64(0))
			delayLoss.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "down", v.Value.Usertype).Set(float64(0))
			loss.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "down", v.Value.Usertype).Set(float64(0))
			lossOrg.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "down", v.Value.Usertype).Set(float64(0))
			lossOrgAvg.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "down", v.Value.Usertype).Set(float64(0))
			lossOrgStd.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "down", v.Value.Usertype).Set(float64(0))
			dropline.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "down", v.Value.Usertype).Set(float64(0))
			videomaxdelay.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "down", v.Value.Usertype).Set(float64(0))
			delete(uservideoqualitydown, k)
		}

		user := uservideoqualitydown[k]
		user.Value.Crarr = ""
		user.Value.Frarr = ""
		user.Value.Delayarr = ""
		user.Value.Delayloss = 0
		user.Value.Loss = 0
		user.Value.Lossorgarr = ""
		user.Value.Dropline = false
		user.Value.Users = 0
		uservideoqualitydown[k] = user
	}
	//音频上行
	for k, v := range uservideoqualityaudioup {
		if v.Value.Meetingid == 0 {
			continue
		}
		//fmt.Println("up音频输出", v.Id, v.Value.Meetingid, "1,2,3,4,10", v.Value.Oneempty, v.Value.Twoempty, v.Value.Threeempty, v.Value.Fourempty, v.Value.Tenempty, "|", v.Value.Users)
		if v.Value.Users != 0 {
			var dev, net string
			Devicetypestr := strconv.Itoa(v.Value.Devicetype)
			Networktypestr := strconv.Itoa(v.Value.Networktype)
			if devmeet[Devicetypestr] != "" {
				dev = devmeet[Devicetypestr]
			} else {
				dev = "unknown"
			}
			if netmeet[Networktypestr] != "" {
				net = netmeet[Networktypestr]
			} else {
				net = "unknown"
			}

			meetidstr := strconv.FormatInt(int64(v.Value.Meetingid), 10)

			lossfloatArr := strArrToFloatArr(strings.Split(v.Value.Lossorgarr, ","))
			lossAvgValue, _ := stats.Mean(lossfloatArr)
			lossStdValue, _ := stats.StandardDeviationPopulation(lossfloatArr)

			var lossorg float64
			for _, v := range lossfloatArr {
				if v != 0 {
					lossorg++
				}
			}
			audioloss.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "up", v.Value.Usertype).Set(float64(v.Value.Loss))
			audiolossOrg.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "up", v.Value.Usertype).Set(float64(lossorg))
			audiolossAvg.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "up", v.Value.Usertype).Set(float64(lossAvgValue))
			audiolossStd.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "up", v.Value.Usertype).Set(float64(lossStdValue))

			Emptyaudiobag.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "up", "one", v.Value.Usertype).Set(float64(v.Value.Oneempty))
			Emptyaudiobag.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "up", "two", v.Value.Usertype).Set(float64(v.Value.Twoempty))
			Emptyaudiobag.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "up", "three", v.Value.Usertype).Set(float64(v.Value.Threeempty))
			Emptyaudiobag.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "up", "foure", v.Value.Usertype).Set(float64(v.Value.Fourempty))
			Emptyaudiobag.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "up", "ten", v.Value.Usertype).Set(float64(v.Value.Tenempty))

		} else {
			var dev, net string
			Devicetypestr := strconv.Itoa(v.Value.Devicetype)
			Networktypestr := strconv.Itoa(v.Value.Networktype)
			if devmeet[Devicetypestr] != "" {
				dev = devmeet[Devicetypestr]
			} else {
				dev = "unknown"
			}
			if netmeet[Networktypestr] != "" {
				net = netmeet[Networktypestr]
			} else {
				net = "unknown"
			}
			meetidstr := strconv.FormatInt(int64(v.Value.Meetingid), 10)

			audioloss.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "up", v.Value.Usertype).Set(float64(0))
			audiolossOrg.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "up", v.Value.Usertype).Set(float64(0))
			audiolossAvg.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "up", v.Value.Usertype).Set(float64(0))
			audiolossStd.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "up", v.Value.Usertype).Set(float64(0))

			Emptyaudiobag.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "up", "one", v.Value.Usertype).Set(float64(0))
			Emptyaudiobag.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "up", "two", v.Value.Usertype).Set(float64(0))
			Emptyaudiobag.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "up", "three", v.Value.Usertype).Set(float64(0))
			Emptyaudiobag.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "up", "foure", v.Value.Usertype).Set(float64(0))
			Emptyaudiobag.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "up", "ten", v.Value.Usertype).Set(float64(0))
			delete(uservideoqualityaudioup, k)
		}

		user := uservideoqualityaudioup[k]

		user.Value.Users = 0
		user.Value.Oneempty = 0
		user.Value.Twoempty = 0
		user.Value.Threeempty = 0
		user.Value.Fourempty = 0
		user.Value.Tenempty = 0
		uservideoqualityaudioup[k] = user
	}
	//音频下行
	for k, v := range uservideoqualityaudiodown {
		if v.Value.Meetingid == 0 {
			continue
		}
		//fmt.Println("down音频输出", v.Id, v.Value.Meetingid, "1,2,3,4,10", v.Value.Oneempty, v.Value.Twoempty, v.Value.Threeempty, v.Value.Fourempty, v.Value.Tenempty, "|", v.Value.Users)
		if v.Value.Users != 0 {
			var dev, net string

			Devicetypestr := strconv.Itoa(v.Value.Devicetype)
			Networktypestr := strconv.Itoa(v.Value.Networktype)
			if devmeet[Devicetypestr] != "" {
				dev = devmeet[Devicetypestr]
			} else {
				dev = "unknown"
			}
			if netmeet[Networktypestr] != "" {
				net = netmeet[Networktypestr]
			} else {
				net = "unknown"
			}

			meetidstr := strconv.FormatInt(int64(v.Value.Meetingid), 10)

			lossfloatArr := strArrToFloatArr(strings.Split(v.Value.Lossorgarr, ","))
			lossAvgValue, _ := stats.Mean(lossfloatArr)
			lossStdValue, _ := stats.StandardDeviationPopulation(lossfloatArr)

			var lossorg float64
			for _, v := range lossfloatArr {
				if v != 0 {
					lossorg++
				}
			}
			audioloss.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "down", v.Value.Usertype).Set(float64(v.Value.Loss))
			audiolossOrg.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "down", v.Value.Usertype).Set(float64(lossorg))
			audiolossAvg.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "down", v.Value.Usertype).Set(float64(lossAvgValue))
			audiolossStd.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "down", v.Value.Usertype).Set(float64(lossStdValue))

			Emptyaudiobag.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "down", "one", v.Value.Usertype).Set(float64(v.Value.Oneempty))
			Emptyaudiobag.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "down", "two", v.Value.Usertype).Set(float64(v.Value.Twoempty))
			Emptyaudiobag.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "down", "three", v.Value.Usertype).Set(float64(v.Value.Threeempty))
			Emptyaudiobag.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "down", "foure", v.Value.Usertype).Set(float64(v.Value.Fourempty))
			Emptyaudiobag.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "down", "ten", v.Value.Usertype).Set(float64(v.Value.Tenempty))

		} else {
			var dev, net string
			Devicetypestr := strconv.Itoa(v.Value.Devicetype)
			Networktypestr := strconv.Itoa(v.Value.Networktype)
			if devmeet[Devicetypestr] != "" {
				dev = devmeet[Devicetypestr]
			} else {
				dev = "unknown"
			}
			if netmeet[Networktypestr] != "" {
				net = netmeet[Networktypestr]
			} else {
				net = "unknown"
			}
			meetidstr := strconv.FormatInt(int64(v.Value.Meetingid), 10)

			audioloss.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "down", v.Value.Usertype).Set(float64(0))
			audiolossOrg.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "down", v.Value.Usertype).Set(float64(0))
			audiolossAvg.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "down", v.Value.Usertype).Set(float64(0))
			audiolossStd.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "down", v.Value.Usertype).Set(float64(0))

			Emptyaudiobag.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "down", "one", v.Value.Usertype).Set(float64(0))
			Emptyaudiobag.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "down", "two", v.Value.Usertype).Set(float64(0))
			Emptyaudiobag.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "down", "three", v.Value.Usertype).Set(float64(0))
			Emptyaudiobag.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "down", "foure", v.Value.Usertype).Set(float64(0))
			Emptyaudiobag.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "down", "ten", v.Value.Usertype).Set(float64(0))
			delete(uservideoqualityaudiodown, k)
		}

		user := uservideoqualityaudiodown[k]
		//内容制零
		user.Value.Users = 0
		user.Value.Oneempty = 0
		user.Value.Twoempty = 0
		user.Value.Threeempty = 0
		user.Value.Fourempty = 0
		user.Value.Tenempty = 0
		uservideoqualityaudiodown[k] = user
	}
	for k, v := range meetqualified {
		if v.Qualified == 0 {
			kstr := strconv.Itoa(k)
			meetqualifiedrate.WithLabelValues(kstr, v.Usertype).Set(0)
			delete(meetqualified, k)
			continue
		}
		kstr := strconv.Itoa(k)
		meetqualifiedrate.WithLabelValues(kstr, v.Usertype).Set(v.Qualified)
		meetqualified[k] = usermeet{
			Qualified: 0,
			Usertype:  v.Usertype,
		}
	}
	for k, v := range meetc2cqualifiedmap {
		if v == 0 {
			meetingidstr := strconv.Itoa(int(k.Meetingid))
			meetc2cqualified.WithLabelValues(k.Called, k.Caller, meetingidstr, k.Qualified, k.Cpu, k.Usertype).Set(0)
			delete(meetc2cqualifiedmap, k)
			continue
		}
		meetingidstr := strconv.Itoa(int(k.Meetingid))
		meetc2cqualified.WithLabelValues(k.Called, k.Caller, meetingidstr, k.Qualified, k.Cpu, k.Usertype).Set(v)
		meetc2cqualifiedmap[k] = 0
	}
	fmt.Println("puseOK")
}

//推送现在正在会议的数量
func Observemeet(meetumber int) {

	node.WithLabelValues("meet").Set(float64(meetumber))
}

//推送商业用户的会议数量
func syObserve(meetumber int) {
	fmt.Println("商业会议数量", meetumber)
	synode.WithLabelValues("meet").Set(float64(meetumber))

}

//字符串数组转换为float数组
func strArrToFloatArr(strArr []string) []float64 {
	i := len(strArr)
	floatArr := make([]float64, i)
	for i, v := range strArr {
		f, _ := strconv.ParseFloat(v, 64)
		floatArr[i] = f
	}
	return floatArr
}
