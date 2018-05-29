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
	//fmt.Println("当前企业通话数", len(entToNumber1))
	for k, v := range entToNumber1 {
		nodenow.WithLabelValues(k, "Now").Set(float64(v))
		entToNumber1[k] = 0
	}
	//新增企业通话数
	//fmt.Println("新增企业通话数", len(entToNumber2))
	for k, v := range entToNumber2 {
		nodenow.WithLabelValues(k, "New").Set(float64(v))
		entToNumber2[k] = 0
	}
	//企业未接通数
	fmt.Println("企业未接通数", len(entToNumber3))
	for k, v := range entToNumber3 {
		nodenow.WithLabelValues(k, "CallFail").Set(float64(v))
		entToNumber3[k] = 0
	}
	//企业延迟插入的通话数（时间与实际时间相差5m及以上）
	//fmt.Println("企业延迟插入的通话数", len(entToNumber4))
	for k, v := range entToNumber4 {
		nodenow.WithLabelValues(k, "Glean").Set(float64(v))
		entToNumber4[k] = 0
	}
	//企业通话总时长
	//fmt.Println("企业通话总时长", len(entToNumber5))
	for k, v := range entToNumber5 {
		nodenowcalltime.WithLabelValues(k, "CallTime").Set(float64(v))
		entToNumber5[k] = 0
	}
	//按地域划分的通话数（p2p）
	//fmt.Println("按地域划分的通话数（p2p）", len(domToNumber))
	for k, v := range domToNumber {
		if dommap[k] == "" {
			nodecallip.WithLabelValues("unknown", "dom").Set(float64(v))
		} else {
			nodecallip.WithLabelValues(dommap[k], "dom").Set(float64(v))
		}
		domToNumber[k] = 0
	}
	//按运营商划分的通话数（p2p）
	//fmt.Println("按运营商划分的通话数（p2p）", len(ispToNumber))
	for k, v := range ispToNumber {
		if ispmap[k] == "" {
			nodecallip.WithLabelValues("unknown", "isp").Set(float64(v))
		} else {
			nodecallip.WithLabelValues(ispmap[k], "isp").Set(float64(v))
		}
		ispToNumber[k] = 0
	}
	//按设备类型划分的通话数（p2p）
	//fmt.Println("按设备类型划分的通话数（p2p）", len(devToNumber))
	for k, v := range devToNumber {
		if devp2pmap[k] == "" {
			nodecallip.WithLabelValues("unknown", "dev").Set(float64(v))
		} else {
			nodecallip.WithLabelValues(devp2pmap[k], "dev").Set(float64(v))
		}
		devToNumber[k] = 0
	}
	//按网络类型划分的通话数（p2p）
	//fmt.Println("按网络类型划分的通话数（p2p）", len(netToNumber))
	for k, v := range netToNumber {
		if netp2pmap[k] == "" {
			nodecallip.WithLabelValues("unknown", "net").Set(float64(v))
		} else {
			nodecallip.WithLabelValues(netp2pmap[k], "net").Set(float64(v))
		}

		netToNumber[k] = 0
	}
	//企业通话时长（meet）
	//fmt.Println("企业通话时长（meet）", len(entTocallTime))
	for k, v := range entTocallTime {
		entcalltimemeet.WithLabelValues(k).Set(float64(v))
		entTocallTime[k] = 0
	}
	//按地域划分的通话数（meet）
	//fmt.Println("按地域划分的通话数（meet）", len(domToNumbermeet))
	for k, v := range domToNumbermeet {
		entipmeet.WithLabelValues(k, "dom").Set(float64(v))
		domToNumbermeet[k] = 0
	}
	//按运营商划分的通话数（meet）
	//fmt.Println("按运营商划分的通话数（meet）", len(ispToNumbermeet))
	for k, v := range ispToNumbermeet {
		entipmeet.WithLabelValues(k, "isp").Set(float64(v))
		ispToNumbermeet[k] = 0
	}
	//按设备类型划分的通话数（meet）
	//fmt.Println("按设备类型划分的通话数（meet）", len(devToNumbermeet))
	for k, v := range devToNumbermeet {
		entipmeet.WithLabelValues(k, "dev").Set(float64(v))
		devToNumbermeet[k] = 0
	}
	//按网络类型划分的通话数（meet）
	//fmt.Println("按网络类型划分的通话数（meet）", len(netToNumbermeet))
	for k, v := range netToNumbermeet {
		entipmeet.WithLabelValues(k, "net").Set(float64(v))
		netToNumbermeet[k] = 0
	}
	//fmt.Println("meetrelayup", len(meetrelayup))
	for k, v := range meetrelayup {
		meetrelay.WithLabelValues(k, "up").Set(float64(v / 60))
		meetrelayup[k] = 0
	}
	//fmt.Println("meetrelaydown", len(meetrelaydown))
	for k, v := range meetrelaydown {
		meetrelay.WithLabelValues(k, "down").Set(float64(v / 60))
		meetrelaydown[k] = 0
	}
	//p2p用户活跃统计
	//fmt.Println("p2p用户活跃统计", len(newusersp2p))
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
	//fmt.Println("meet用户活跃统计", len(newusersmeet))
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

	//user <-> speaker 之间的码流统计
	//fmt.Println("user <-> speaker 之间的码流统计", len(userspeaktraffic))
	var b int = 0
	for _, v := range userspeaktraffic {
		b += len(v.flow)
	}
	//fmt.Println("user <-> speaker 之间的码流统计flow", b)
	for k, v := range userspeaktraffic {
		if v.Relaydomisp == "" {
			fmt.Println(v, "172")
			continue
		}
		var updown string
		if k.Speakerid == k.Userid {
			updown = "up"
		} else {
			updown = "down"
		}
		var dev string
		devs := strings.Split(v.Mid, "_")
		Devicetypestr := devs[0]
		if devmeet[Devicetypestr] != "" {
			dev = devmeet[Devicetypestr]
		} else {
			dev = "unknown"
		}
		if v.flow != "" {
			trfloatArr := strArrToFloatArr(strings.Split(v.flow, ","))
			trAvgValue, _ := stats.Mean(trfloatArr)
			trStdValue, _ := stats.StandardDeviationPopulation(trfloatArr)
			trMaxValue, _ := stats.Max(trfloatArr)
			trMinValue, _ := stats.Min(trfloatArr)

			trAvg.WithLabelValues(k.Userid, k.Speakerid, k.Meetintid, updown, v.UserDom, v.UserIsp, v.Relaydomisp, dev).Set(trAvgValue)
			trStd.WithLabelValues(k.Userid, k.Speakerid, k.Meetintid, updown, v.UserDom, v.UserIsp, v.Relaydomisp, dev).Set(trStdValue)
			trMax.WithLabelValues(k.Userid, k.Speakerid, k.Meetintid, updown, v.UserDom, v.UserIsp, v.Relaydomisp, dev).Set(trMaxValue)
			trMin.WithLabelValues(k.Userid, k.Speakerid, k.Meetintid, updown, v.UserDom, v.UserIsp, v.Relaydomisp, dev).Set(trMinValue)
			a := userspeaktraffic[k]
			a.flow = ""
			a.UserIp = ""
			userspeaktraffic[k] = a
		} else {
			trAvg.WithLabelValues(k.Userid, k.Speakerid, k.Meetintid, updown, v.UserDom, v.UserIsp, v.Relaydomisp, dev).Set(0)
			trStd.WithLabelValues(k.Userid, k.Speakerid, k.Meetintid, updown, v.UserDom, v.UserIsp, v.Relaydomisp, dev).Set(0)
			trMax.WithLabelValues(k.Userid, k.Speakerid, k.Meetintid, updown, v.UserDom, v.UserIsp, v.Relaydomisp, dev).Set(0)
			trMin.WithLabelValues(k.Userid, k.Speakerid, k.Meetintid, updown, v.UserDom, v.UserIsp, v.Relaydomisp, dev).Set(0)

			trAvg.DeleteLabelValues(k.Userid, k.Speakerid, k.Meetintid, updown, v.UserDom, v.UserIsp, v.Relaydomisp, dev)
			trStd.DeleteLabelValues(k.Userid, k.Speakerid, k.Meetintid, updown, v.UserDom, v.UserIsp, v.Relaydomisp, dev)
			trMax.DeleteLabelValues(k.Userid, k.Speakerid, k.Meetintid, updown, v.UserDom, v.UserIsp, v.Relaydomisp, dev)
			trMin.DeleteLabelValues(k.Userid, k.Speakerid, k.Meetintid, updown, v.UserDom, v.UserIsp, v.Relaydomisp, dev)
			delete(userspeaktraffic, k)

		}

	}

	//relay间流量
	//fmt.Println("relay间流量", len(relaybetweens))
	var a int = 0
	for _, v := range relaybetweens {
		a += len(v.flow)
	}
	//fmt.Println("relay间流量flow", a)
	for k, v := range relaybetweens {

		meetid := strconv.Itoa(int(k.Meetintid))
		if v.flow != "" {
			relayflowfloatArr := strArrToFloatArr(strings.Split(v.flow, ","))
			relauflowAvgValue, _ := stats.Mean(relayflowfloatArr)
			relaytorelayflow.WithLabelValues(k.Relayid, k.Userip, k.Speakerid, meetid, k.Usertype).Set(relauflowAvgValue)
			v.flow = ""
			v.isnull = 1
			relaybetweens[k] = v
		} else {
			relaytorelayflow.WithLabelValues(k.Relayid, k.Userip, k.Speakerid, meetid, k.Usertype).Set(0)

			relaytorelayflow.DeleteLabelValues(k.Relayid, k.Userip, k.Speakerid, meetid, k.Usertype)

			if v.isnull == 1 {
				delete(relaybetweens, k)
				continue
			}
		}

		//fmt.Println("relay间流量", k.Relayid, k.Userip, k.Speakerid, k.Meetintid, k.Usertype, relauflowAvgValue)
	}
	//relay间传输丢包

	relaymap := globeCfg.Output.Relaymap
	//fmt.Println("relay间传输丢包", len(relaybetweenloss))
	for k, v := range relaybetweenloss {

		meetid := strconv.Itoa(int(k.Meetintid))
		LossFinArr := strArrToFloatArr(strings.Split(v.LossFinArr, ","))
		LossOrgArr := strArrToFloatArr(strings.Split(v.LossOrgArr, ","))
		relayFinAvglossValue, _ := stats.Mean(LossFinArr)
		relayOrgAvglossValue, _ := stats.Mean(LossOrgArr)

		relayFinMaxlossValue, _ := stats.Max(LossFinArr)
		relayOrgMaxlossValue, _ := stats.Max(LossOrgArr)
		var relauFinMinlossValue float64
		var relauOrgMinlossValue float64
		relauFinMinlossValue = 0
		relauOrgMinlossValue = 0
		for _, v := range LossFinArr {
			if v == 0 {
				continue
			}
			if (v > relauFinMinlossValue && relauFinMinlossValue == 0) || (v < relauFinMinlossValue && relauFinMinlossValue != 0) {
				relauFinMinlossValue = v
			}
		}
		for _, v := range LossOrgArr {
			if v == 0 {
				continue
			}
			if (v > relauOrgMinlossValue && relauOrgMinlossValue == 0) || (v < relauOrgMinlossValue && relauOrgMinlossValue != 0) {
				relauOrgMinlossValue = v
			}
		}

		var relayFinlossValue float64
		for _, v := range LossFinArr {
			if v != 0 {
				relayFinlossValue++
			}
		}
		var relayOrglossValue float64
		for _, v := range LossOrgArr {
			if v != 0 {
				relayOrglossValue++
			}
		}
		var relayDomIsp string

		relayid := strings.Split(k.Mid, "_")
		dev := strconv.Itoa(k.Deviceid)

		if relaymap[dev+"|"+relayid[0]] != "" {
			relayDomIsp = relaymap[dev+"|"+relayid[0]]
		} else {
			relayDomIsp = "unknown"
		}

		if v.isnull == 0 {
			relaytorelayFinAvgloss.WithLabelValues(k.Relayid, k.Userip, k.Speakerid, meetid, k.Usertype, relayDomIsp).Set(relayFinAvglossValue)
			relaytorelayFinMaxloss.WithLabelValues(k.Relayid, k.Userip, k.Speakerid, meetid, k.Usertype, relayDomIsp).Set(relayFinMaxlossValue)
			relaytorelayFinMinloss.WithLabelValues(k.Relayid, k.Userip, k.Speakerid, meetid, k.Usertype, relayDomIsp).Set(relauFinMinlossValue)
			relaytorelayFinloss.WithLabelValues(k.Relayid, k.Userip, k.Speakerid, meetid, k.Usertype, relayDomIsp).Set(relayFinlossValue)

			relaytorelayOrgAvgloss.WithLabelValues(k.Relayid, k.Userip, k.Speakerid, meetid, k.Usertype, relayDomIsp).Set(relayOrgAvglossValue)
			relaytorelayOrgMaxloss.WithLabelValues(k.Relayid, k.Userip, k.Speakerid, meetid, k.Usertype, relayDomIsp).Set(relayOrgMaxlossValue)
			relaytorelayOrgMinloss.WithLabelValues(k.Relayid, k.Userip, k.Speakerid, meetid, k.Usertype, relayDomIsp).Set(relauOrgMinlossValue)
			relaytorelayOrgloss.WithLabelValues(k.Relayid, k.Userip, k.Speakerid, meetid, k.Usertype, relayDomIsp).Set(relayOrglossValue)

			userlossvalues := relaybetweenloss[k]
			userlossvalues.LossFinArr = ""
			userlossvalues.LossOrgArr = ""
			userlossvalues.isnull = 1
			relaybetweenloss[k] = userlossvalues
		} else {
			relaytorelayFinAvgloss.WithLabelValues(k.Relayid, k.Userip, k.Speakerid, meetid, k.Usertype, relayDomIsp).Set(0)
			relaytorelayFinMaxloss.WithLabelValues(k.Relayid, k.Userip, k.Speakerid, meetid, k.Usertype, relayDomIsp).Set(0)
			relaytorelayFinMinloss.WithLabelValues(k.Relayid, k.Userip, k.Speakerid, meetid, k.Usertype, relayDomIsp).Set(0)
			relaytorelayFinloss.WithLabelValues(k.Relayid, k.Userip, k.Speakerid, meetid, k.Usertype, relayDomIsp).Set(0)

			relaytorelayOrgAvgloss.WithLabelValues(k.Relayid, k.Userip, k.Speakerid, meetid, k.Usertype, relayDomIsp).Set(0)
			relaytorelayOrgMaxloss.WithLabelValues(k.Relayid, k.Userip, k.Speakerid, meetid, k.Usertype, relayDomIsp).Set(0)
			relaytorelayOrgMinloss.WithLabelValues(k.Relayid, k.Userip, k.Speakerid, meetid, k.Usertype, relayDomIsp).Set(0)
			relaytorelayOrgloss.WithLabelValues(k.Relayid, k.Userip, k.Speakerid, meetid, k.Usertype, relayDomIsp).Set(0)

			relaytorelayFinAvgloss.DeleteLabelValues(k.Relayid, k.Userip, k.Speakerid, meetid, k.Usertype, relayDomIsp)
			relaytorelayFinMaxloss.DeleteLabelValues(k.Relayid, k.Userip, k.Speakerid, meetid, k.Usertype, relayDomIsp)
			relaytorelayFinMinloss.DeleteLabelValues(k.Relayid, k.Userip, k.Speakerid, meetid, k.Usertype, relayDomIsp)
			relaytorelayFinloss.DeleteLabelValues(k.Relayid, k.Userip, k.Speakerid, meetid, k.Usertype, relayDomIsp)

			relaytorelayOrgAvgloss.DeleteLabelValues(k.Relayid, k.Userip, k.Speakerid, meetid, k.Usertype, relayDomIsp)
			relaytorelayOrgMaxloss.DeleteLabelValues(k.Relayid, k.Userip, k.Speakerid, meetid, k.Usertype, relayDomIsp)
			relaytorelayOrgMinloss.DeleteLabelValues(k.Relayid, k.Userip, k.Speakerid, meetid, k.Usertype, relayDomIsp)
			relaytorelayOrgloss.DeleteLabelValues(k.Relayid, k.Userip, k.Speakerid, meetid, k.Usertype, relayDomIsp)

			delete(relaybetweenloss, k)
			continue
		}

		//fmt.Println("relay间传输丢包", k.Relayid, k.Userip, k.Speakerid, k.Meetintid, k.Usertype, relayDomIsp, relayOrglossValue)
	}

	//文档上下行
	//fmt.Println("文档上下行", len(userfilequalityupdown))
	for k, v := range userfilequalityupdown {
		var updown string
		if v.Id.Speakerid == v.Id.Userid {
			updown = "up"
		} else {
			updown = "down"
		}
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
			//最大延迟
			maxdelayup, _ = stats.Max(delayfloatArr)

			meetidstr := strconv.FormatInt(int64(v.Value.Meetingid), 10)
			crAvgValue, _ := stats.Mean(CrfloatArr)
			crStdValue, _ := stats.StandardDeviationPopulation(CrfloatArr)
			crMaxValue, _ := stats.Max(CrfloatArr)
			crMinValue, _ := stats.Min(CrfloatArr)
			frAvgValue, _ := stats.Mean(FrfloatArr)
			frStdValue, _ := stats.StandardDeviationPopulation(FrfloatArr)
			frMaxValue, _ := stats.Max(FrfloatArr)
			frMinValue, _ := stats.Min(FrfloatArr)
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
			filecrAvg.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype).Set(float64(int64(crAvgValue)))
			filecrStd.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype).Set(float64(int64(crStdValue)))
			filefrAvg.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype).Set(float64(int64(frAvgValue)))
			filefrStd.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype).Set(float64(int64(frStdValue)))
			filecrMax.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype).Set(float64(int64(crMaxValue)))
			filecrMin.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype).Set(float64(int64(crMinValue)))
			filefrMax.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype).Set(float64(int64(frMaxValue)))
			filefrMin.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype).Set(float64(int64(frMinValue)))
			filedelayAvg.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype).Set(float64(int64(delayAvgValue)))
			filedelayStd.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype).Set(float64(int64(delayStdValue)))
			filedelayLoss.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype).Set(float64(int64(v.Value.Delayloss)))
			fileloss.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype).Set(float64(int64(v.Value.Loss)))
			filelossOrg.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype).Set(float64(int64(lossorg)))
			filelossOrgAvg.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype).Set(float64(int64(lossAvgValue)))
			filelossOrgStd.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype).Set(float64(int64(lossStdValue)))
			filevideomaxdelay.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype).Set(maxdelayup)
			if v.Value.Dropline == true {
				filedropline.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype).Set(float64(1))

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
			filecrAvg.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype).Set(float64(0))
			filecrStd.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype).Set(float64(0))
			filefrAvg.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype).Set(float64(0))
			filefrStd.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype).Set(float64(0))
			filecrMax.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype).Set(float64(0))
			filecrMin.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype).Set(float64(0))
			filefrMax.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype).Set(float64(0))
			filefrMin.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype).Set(float64(0))
			filedelayAvg.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype).Set(float64(0))
			filedelayStd.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype).Set(float64(0))
			filedelayLoss.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype).Set(float64(0))
			fileloss.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype).Set(float64(0))
			filelossOrg.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype).Set(float64(0))
			filelossOrgAvg.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype).Set(float64(0))
			filelossOrgStd.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype).Set(float64(0))
			filedropline.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype).Set(float64(0))
			filevideomaxdelay.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype).Set(float64(0))

			filecrAvg.DeleteLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype)
			filecrStd.DeleteLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype)
			filefrAvg.DeleteLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype)
			filefrStd.DeleteLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype)
			filecrMax.DeleteLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype)
			filecrMin.DeleteLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype)
			filefrMax.DeleteLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype)
			filefrMin.DeleteLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype)
			filedelayAvg.DeleteLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype)
			filedelayStd.DeleteLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype)
			filedelayLoss.DeleteLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype)
			fileloss.DeleteLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype)
			filelossOrg.DeleteLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype)
			filelossOrgAvg.DeleteLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype)
			filelossOrgStd.DeleteLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype)
			filedropline.DeleteLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype)
			filevideomaxdelay.DeleteLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype)

			delete(userfilequalityupdown, k)
			continue
		}

		user := userfilequalityupdown[k]
		user.Value.Crarr = ""
		user.Value.Frarr = ""
		user.Value.Delayarr = ""
		user.Value.Delayloss = 0
		user.Value.Loss = 0
		user.Value.Lossorgarr = ""
		user.Value.Dropline = false
		user.Value.Users = 0
		userfilequalityupdown[k] = user
	}

	//视频上下行
	//fmt.Println("视频上下行", len(uservideoqualityupdown))
	for k, v := range uservideoqualityupdown {
		var updown string
		if v.Id.Speakerid == v.Id.Userid {
			updown = "up"
		} else {
			updown = "down"
		}
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
			crMaxValue, _ := stats.Max(CrfloatArr)
			crMinValue, _ := stats.Min(CrfloatArr)
			frMaxValue, _ := stats.Max(FrfloatArr)
			frMinValue, _ := stats.Min(FrfloatArr)

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
			crAvg.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype).Set(float64(int64(crAvgValue)))
			crStd.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype).Set(float64(int64(crStdValue)))
			frAvg.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype).Set(float64(int64(frAvgValue)))
			frStd.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype).Set(float64(int64(frStdValue)))
			crMax.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype).Set(float64(int64(crMaxValue)))
			crMin.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype).Set(float64(int64(crMinValue)))
			frMax.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype).Set(float64(int64(frMaxValue)))
			frMin.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype).Set(float64(int64(frMinValue)))
			delayAvg.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype).Set(float64(int64(delayAvgValue)))
			delayStd.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype).Set(float64(int64(delayStdValue)))
			delayLoss.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype).Set(float64(int64(v.Value.Delayloss)))
			loss.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype).Set(float64(int64(v.Value.Loss)))
			lossOrg.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype).Set(float64(int64(lossorg)))
			lossOrgAvg.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype).Set(float64(int64(lossAvgValue)))
			lossOrgStd.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype).Set(float64(int64(lossStdValue)))
			videomaxdelay.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype).Set(maxdelayup)
			if v.Value.Dropline == true {
				dropline.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype).Set(float64(1))

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
			crAvg.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype).Set(float64(0))
			crStd.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype).Set(float64(0))
			frAvg.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype).Set(float64(0))
			frStd.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype).Set(float64(0))
			crMax.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype).Set(float64(0))
			crMin.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype).Set(float64(0))
			frMax.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype).Set(float64(0))
			frMin.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype).Set(float64(0))
			delayAvg.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype).Set(float64(0))
			delayStd.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype).Set(float64(0))
			delayLoss.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype).Set(float64(0))
			loss.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype).Set(float64(0))
			lossOrg.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype).Set(float64(0))
			lossOrgAvg.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype).Set(float64(0))
			lossOrgStd.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype).Set(float64(0))
			dropline.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype).Set(float64(0))
			videomaxdelay.WithLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype).Set(float64(0))

			crAvg.DeleteLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype)
			crStd.DeleteLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype)
			frAvg.DeleteLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype)
			frStd.DeleteLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype)
			crMax.DeleteLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype)
			crMin.DeleteLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype)
			frMax.DeleteLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype)
			frMin.DeleteLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype)
			delayAvg.DeleteLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype)
			delayStd.DeleteLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype)
			delayLoss.DeleteLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype)
			loss.DeleteLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype)
			lossOrg.DeleteLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype)
			lossOrgAvg.DeleteLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype)
			lossOrgStd.DeleteLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype)
			dropline.DeleteLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype)
			videomaxdelay.DeleteLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, updown, v.Value.Usertype)

			delete(uservideoqualityupdown, k)
			continue
		}

		user := uservideoqualityupdown[k]
		user.Value.Crarr = ""
		user.Value.Frarr = ""
		user.Value.Delayarr = ""
		user.Value.Delayloss = 0
		user.Value.Loss = 0
		user.Value.Lossorgarr = ""
		user.Value.Dropline = false
		user.Value.Users = 0
		uservideoqualityupdown[k] = user
	}

	//音频上行
	//fmt.Println("音频上行", len(uservideoqualityaudioup))
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

			audioloss.DeleteLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "up", v.Value.Usertype)
			audiolossOrg.DeleteLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "up", v.Value.Usertype)
			audiolossAvg.DeleteLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "up", v.Value.Usertype)
			audiolossStd.DeleteLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "up", v.Value.Usertype)

			Emptyaudiobag.DeleteLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "up", "one", v.Value.Usertype)
			Emptyaudiobag.DeleteLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "up", "two", v.Value.Usertype)
			Emptyaudiobag.DeleteLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "up", "three", v.Value.Usertype)
			Emptyaudiobag.DeleteLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "up", "foure", v.Value.Usertype)
			Emptyaudiobag.DeleteLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "up", "ten", v.Value.Usertype)

			delete(uservideoqualityaudioup, k)
			continue
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
	//fmt.Println("音频下行", len(uservideoqualityaudiodown))
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

			audioloss.DeleteLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "down", v.Value.Usertype)
			audiolossOrg.DeleteLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "down", v.Value.Usertype)
			audiolossAvg.DeleteLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "down", v.Value.Usertype)
			audiolossStd.DeleteLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "down", v.Value.Usertype)

			Emptyaudiobag.DeleteLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "down", "one", v.Value.Usertype)
			Emptyaudiobag.DeleteLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "down", "two", v.Value.Usertype)
			Emptyaudiobag.DeleteLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "down", "three", v.Value.Usertype)
			Emptyaudiobag.DeleteLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "down", "foure", v.Value.Usertype)
			Emptyaudiobag.DeleteLabelValues(v.Id.Userid, v.Id.Speakerid, meetidstr, dev, net, v.Value.Userent, v.Value.Meetent, v.Value.Userdom, v.Value.Userisp, v.Value.Relaydomisp, "down", "ten", v.Value.Usertype)
			delete(uservideoqualityaudiodown, k)
			continue
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
