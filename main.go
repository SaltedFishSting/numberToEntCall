package main

import (
	"container/list"
	"database/sql"
	"fmt"
	"gproto"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/axgle/mahonia"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/client_golang/prometheus/push"
	"github.com/tealeg/xlsx"
	"github.com/wangtuanjie/ip17mon"
	"gopkg.in/mgo.v2"
	"gopkg.in/yaml.v2"
)

//配置文件yaml
type RConfig struct {
	Gw struct {
		Addr           string `yaml:"addr"`
		HttpListenPort int    `yaml:"httpListenPort"`
	}
	Output struct {
		Prometheus      bool              `yaml:"prometheus"`
		PushGateway     bool              `yaml:"pushGateway"`
		PushGatewayAddr string            `yaml:"pushGatewayAddr"`
		MonitorID       string            `yaml:"monitorID"`
		Period          int               `yaml:"period"`
		Ip              string            `yaml:"ip"`
		Db              string            `yaml:"db"`
		Table1          string            `yaml:"table1"`
		Table2          string            `yaml:"table2"`
		URL             string            `yaml:"url"`
		Dom             map[string]string `yaml:"dom"`
		Isp             map[string]string `yaml:"isp"`
		Devp2p          map[string]string `yaml:"devp2p"`
		Netp2p          map[string]string `yaml:"netp2p"`
		Ipmeet          string            `yaml:"ipmeet"`
		Dbmeet          string            `yaml:"dbmeet"`
		Table1meet      string            `yaml:"table1meet"`
		Devmeet         map[string]string `yaml:"devmeet"`
		Netmeet         map[string]string `yaml:"netmeet"`
		CnlogS          string            `yaml:"cnlogS"`
		Relaymap        map[string]string `yaml:"relaymap"`
	}
}

//mongodb 模型

type sss struct {
	TimeStamp int64 `bson:"timeStamp"`
}

type meetinfo struct {
	EndTime      int32  `bson:"endTime"`
	MeetingId    int64  `bson:"meetingId"`
	StartTime    int32  `bson:"startTime"`
	Duration     int64  `bson:"duration"`
	UserIdList   string `bson:"userIdList"`
	UserCount    int    `bson:"userCount"`
	QosTableName string `bson:"qosTableName"`
}

type p2pinfo struct {
	InsertTime        int64  `bson:"insertTime"`
	Endtime           int64  `bson:"endtime"`
	Starttime         int64  `bson:"starttime"`
	EventType         int64  `bson:"eventType"`
	Called            int64  `bson:"called"`
	Caller            int64  `bson:"caller"`
	SidReporter       string `bson:"sidReporter"`
	DissconnectedTime int64  `bson:"dissconnectedTime"`
}
type p2pispcity struct {
	Id    string "_id"
	Value struct {
		Ldom   string
		Lsip   string
		Ldev   string
		Lnet   string
		Rdom   string
		Rsip   string
		Rdev   string
		Rnet   string
		Number string
	}
}
type uservideoinfo struct {
	Id    string "_id"
	Value struct {
		users       int
		userIp      string
		relayIp     string
		meetingId   int
		deviceType  int
		networkType int
		crAvg       int
		crStd       int
		frAvg       int
		frStd       int
		lossOrg     int
		lossOrgAvg  int
		lossOrgStd  int
		loss        int
		delayAvg    int
		delayStd    int
		delayLoss   int
		dropline    bool
	}
}
type relayflow struct {
	Id    string "_id"
	Value struct {
		Usernum int
		Upbw    int
		Downbw  int
	}
}
type meet2MapReduce struct {
	Id    string "_id"
	Value struct {
		Addr    string
		Dev     int
		Net     int
		Isnew   bool
		Timelen int
		Meetid  int
	}
}

func ConvertToString(src string, srcCode string, tagCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(tagCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	return result
}

//活跃用户判断
var ruc *gproto.RecUserClient

//配置文件
var globeCfg *RConfig

//企业to企业id
var entToNumber map[string]int

//企业id to企业
var numberToEnt map[int]string

//企业当前通话数
var entToNumber1 map[string]int

//企业当前新增通话数
var entToNumber2 map[string]int

//企业当前未接通数
var entToNumber3 map[string]int

//延迟插入的通话数
var entToNumber4 map[string]int

//企业通话时长统计
var entToNumber5 map[string]int64

//按地域的通话数
var domToNumber map[string]int
var dommap map[string]string
var domToNumbermeet map[string]int

//按运营商的通话数
var ispToNumber map[string]int
var ispmap map[string]string
var ispToNumbermeet map[string]int

//按设备的通话数
var devToNumber map[string]int
var devp2pmap map[string]string
var devmeetmap map[string]string
var devToNumbermeet map[string]int

//按网络类型的通话数
var netToNumber map[string]int
var netp2pmap map[string]string
var netmeetmap map[string]string
var netToNumbermeet map[string]int

//企业通话时长meet
var entTocallTime map[string]int

//meet的relay上行流量
var meetrelayup map[string]int

//meet的relay下行流量
var meetrelaydown map[string]int

//用户视讯号to企业id
var numberArray [100000000]int

//商业用户在线和新增
var syuser map[string]int

//上一次mongodb最后插入的时间
var lastTime int64 = 0 //p2p
var lastTime2 int64 = 0
var lasttime int32 = 0 //meet
var timelist list.List

//上一次查询的meetTable
var lasttablemeet string = "meetTable"

//活跃用户统计(p2p)
var newusersp2p map[string]int
var dayOldusersp2p map[string]int
var monOldusersp2p map[string]int
var weekOldusersp2p map[string]int

//活跃用户统计(meet)
var newusersmeet map[string]int
var dayOldusersmeet map[string]int
var weekOldusersmeet map[string]int
var monOldusersmeet map[string]int

//视频体验统计
var uservideoquality map[string]uservideoinfo
var relayvideoquality map[string]uservideoinfo

//prometheus var
var (
	nodenow         *(prometheus.GaugeVec)
	nodenowcalltime *(prometheus.GaugeVec)
	nodecallip      *(prometheus.GaugeVec)
	node            *(prometheus.GaugeVec)
	synode          *(prometheus.GaugeVec)
	entcalltimemeet *(prometheus.GaugeVec)
	entipmeet       *(prometheus.GaugeVec)
	meetrelay       *(prometheus.GaugeVec)
	p2pactive       *(prometheus.GaugeVec)
	meetactive      *(prometheus.GaugeVec)
	crAvg           *(prometheus.GaugeVec)
	crStd           *(prometheus.GaugeVec)
	frAvg           *(prometheus.GaugeVec)
	frStd           *(prometheus.GaugeVec)
	delayAvg        *(prometheus.GaugeVec)
	delayStd        *(prometheus.GaugeVec)
	delayLoss       *(prometheus.GaugeVec)
	dropline        *(prometheus.GaugeVec)
	loss            *(prometheus.GaugeVec)
	lossOrg         *(prometheus.GaugeVec)
	lossOrgAvg      *(prometheus.GaugeVec)
	lossOrgStd      *(prometheus.GaugeVec)
	//解析excal后生成的map
	usermap map[int64]string
)

//获取当前各个企业视频号号段范围
func getNumber(db *sql.DB) {
	//初始化map
	entToNumber = make(map[string]int)
	numberToEnt = make(map[int]string)
	entToNumber1 = make(map[string]int)
	entToNumber2 = make(map[string]int)
	entToNumber3 = make(map[string]int)
	entToNumber4 = make(map[string]int)
	entToNumber5 = make(map[string]int64)
	domToNumber = make(map[string]int)
	ispToNumber = make(map[string]int)
	devToNumber = make(map[string]int)
	netToNumber = make(map[string]int)
	entTocallTime = make(map[string]int)
	domToNumbermeet = make(map[string]int)
	ispToNumbermeet = make(map[string]int)
	devToNumbermeet = make(map[string]int)
	netToNumbermeet = make(map[string]int)
	meetrelayup = make(map[string]int)
	meetrelaydown = make(map[string]int)
	syuser = make(map[string]int)
	newusersp2p = make(map[string]int)
	dayOldusersp2p = make(map[string]int)
	weekOldusersp2p = make(map[string]int)
	monOldusersp2p = make(map[string]int)
	newusersmeet = make(map[string]int)
	dayOldusersmeet = make(map[string]int)
	weekOldusersmeet = make(map[string]int)
	monOldusersmeet = make(map[string]int)
	uservideoquality = make(map[string]uservideoinfo)
	relayvideoquality = make(map[string]uservideoinfo)
	dommap = globeCfg.Output.Dom
	ispmap = globeCfg.Output.Isp
	devp2pmap = globeCfg.Output.Devp2p
	netp2pmap = globeCfg.Output.Netp2p
	devmeetmap = globeCfg.Output.Devmeet
	netmeetmap = globeCfg.Output.Netmeet
	fmt.Println("字典列表如下")
	fmt.Println(dommap)
	fmt.Println(ispmap)
	fmt.Println(devp2pmap)
	fmt.Println(netp2pmap)
	fmt.Println(devmeetmap)
	fmt.Println(netmeetmap)
	fmt.Println("relaymap: ", globeCfg.Output.Relaymap)
	rows, err := db.Query("select entName, t4.nubeEnd,t4.nubeStart, t1.entType from t_ent as t1 left join t_ent_appkey_conf as t2 on  t1.entId=t2.entId left join t_ent_appkey_nube_config as t3 on t2.appKey=t3.appKey left join t_ent_appkey_nubeinte_conf as t4 on  t3.intervalId=t4.intervalId  group by t4.intervalId")
	if err != nil {
		fmt.Println(err, "257")
		return
	}
	i := 1
	numberToEnt[0] = "Unknown"
	for rows.Next() {
		var entName string
		var numberEnd int
		var numberStart int
		var entType int
		err = rows.Scan(&entName, &numberEnd, &numberStart, &entType)
		if err != nil {
			fmt.Println("企业号段空值(有号码没号段)", entName, numberEnd, numberStart, entType)
			fmt.Println(err, "270")
			continue
		}
		if entToNumber[entName] == 0 {
			entToNumber[entName] = i
			numberToEnt[i] = entName
			i++
		}
		for i := numberStart; i <= numberEnd; i++ {
			if i > 99999999 {
				continue
			}
			numberArray[i] = entToNumber[entName]
		}
	}
}

//推送数据
func Observe() {

	synode.WithLabelValues("nowuser").Set(float64(syuser["now"]))
	synode.WithLabelValues("newuser").Set(float64(syuser["new"]))
	synode.WithLabelValues("realuser").Set(float64(syuser["realuser"]))

	syuser["now"] = 0
	syuser["new"] = 0
	syuser["realuser"] = 0
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
	fmt.Println("meetrelayup: ", meetrelayup)
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

	fmt.Println("puseOK")
}

//推送现在正在会议的数量
func Observemeet(meetumber int) {

	node.WithLabelValues("meet").Set(float64(meetumber))
}

//推送商业用户的会议数量
func syObserve(meetumber int, nowuser int, newuser int) {
	fmt.Println("商业会议数量", meetumber, "在线用户", nowuser, "新增用户", newuser)
	synode.WithLabelValues("meet").Set(float64(meetumber))

}

//加载配置文件
func loadConfig() {
	if err := ip17mon.Init("mydata4vipweek2.dat"); err != nil {
		panic(err)
	}
	cfgbuf, err := ioutil.ReadFile("cfg.yaml")
	if err != nil {
		panic("not found cfg.yaml")
	}
	rfig := RConfig{}
	err = yaml.Unmarshal(cfgbuf, &rfig)
	if err != nil {
		panic("invalid cfg.yaml")
	}
	globeCfg = &rfig
	fmt.Println("Load config -'cfg.yaml'- ok...")
	//加载excal文件
	usermap = make(map[int64]string)
	excelFileName := "user.xlsx"
	xlFile, err := xlsx.OpenFile(excelFileName)
	if err != nil {
		fmt.Println(err, "463")
	}
	for _, sheet := range xlFile.Sheets {
		for _, row := range sheet.Rows {
			meetid := 0
			for i, cell := range row.Cells {
				if i == 0 {
					meetid, _ = strconv.Atoi(cell.String())
				}
				if i == 1 {
					if cell.String() == "商业用户" {
						usermap[int64(meetid)] = "商业"
					} else {
						usermap[int64(meetid)] = "非商业"
					}
				}
				if i == 2 {
					if usermap[int64(meetid)] == "商业" {
						usermap[int64(meetid)] = cell.String()
					}
				}
			}
		}
	}
	for k, v := range usermap {
		fmt.Println(k, v)
	}

}
func init() {
	loadConfig() //加载配置文件

	nodenow = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "p2p",
		Subsystem: "entToNumber",
		Name:      "callNumberStatistics",
		Help:      "state",
	}, []string{
		"entName",
		"callType"})

	nodenowcalltime = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "p2p",
		Subsystem: "entToNumber",
		Name:      "callTime",
		Help:      "state",
	}, []string{
		"entName",
		"callType"})
	nodecallip = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "p2p",
		Subsystem: "entToNumber",
		Name:      "domIspDevNet",
		Help:      "state",
	}, []string{
		"domIspDevNetContent",
		"type"})
	node = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "meet",
		Subsystem: "info",
		Name:      "number",
		Help:      "rcstate"}, []string{"meet"})
	synode = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "symeet",
		Subsystem: "info",
		Name:      "number",
		Help:      "rcstate"}, []string{"meet"})

	entcalltimemeet = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "meet",
		Subsystem: "ent",
		Name:      "calltime",
		Help:      "state",
	}, []string{
		"entconfName",
	})

	entipmeet = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "meet",
		Subsystem: "ent",
		Name:      "ip",
		Help:      "state",
	}, []string{
		"domIspDevNetContent",
		"type"})
	meetrelay = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "meet",
		Subsystem: "relay",
		Name:      "updown",
		Help:      "state",
	}, []string{
		"domIsp",
		"updown"})
	p2pactive = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "p2p",
		Subsystem: "user",
		Name:      "active",
		Help:      "state",
	}, []string{
		"entName",
		"Type"})
	meetactive = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "meet",
		Subsystem: "user",
		Name:      "active",
		Help:      "state",
	}, []string{
		"entconfName",
		"Type"})
	crAvg = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "user",
		Subsystem: "quality",
		Name:      "crAvg",
		Help:      "state",
	}, []string{
		"userId",
		"meetingId",
		"userEnt",
		"meetEnt",
		"userDom",
		"userIsp",
		"relayDomIsp"})
	crStd = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "user",
		Subsystem: "quality",
		Name:      "crStd ",
		Help:      "state",
	}, []string{
		"userId",
		"meetingId",
		"userEnt",
		"meetEnt",
		"userDom",
		"userIsp",
		"relayDomIsp"})
	frAvg = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "user",
		Subsystem: "quality",
		Name:      "frAvg",
		Help:      "state",
	}, []string{
		"userId",
		"meetingId",
		"userEnt",
		"meetEnt",
		"userDom",
		"userIsp",
		"relayDomIsp"})
	frStd = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "user",
		Subsystem: "quality",
		Name:      "frStd",
		Help:      "state",
	}, []string{
		"userId",
		"meetingId",
		"userEnt",
		"meetEnt",
		"userDom",
		"userIsp",
		"relayDomIsp"})
	delayAvg = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "user",
		Subsystem: "quality",
		Name:      "delayAvg",
		Help:      "state",
	}, []string{
		"userId",
		"meetingId",
		"userEnt",
		"meetEnt",
		"userDom",
		"userIsp",
		"relayDomIsp"})
	delayStd = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "user",
		Subsystem: "quality",
		Name:      "delayStd",
		Help:      "state",
	}, []string{
		"userId",
		"meetingId",
		"userEnt",
		"meetEnt",
		"userDom",
		"userIsp",
		"relayDomIsp"})
	delayLoss = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "user",
		Subsystem: "quality",
		Name:      "delayLoss",
		Help:      "state",
	}, []string{
		"userId",
		"meetingId",
		"userEnt",
		"meetEnt",
		"userDom",
		"userIsp",
		"relayDomIsp"})
	dropline = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "user",
		Subsystem: "quality",
		Name:      "dropline",
		Help:      "state",
	}, []string{
		"userId",
		"meetingId",
		"userEnt",
		"meetEnt",
		"userDom",
		"userIsp",
		"relayDomIsp"})
	loss = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "user",
		Subsystem: "quality",
		Name:      "loss",
		Help:      "state",
	}, []string{
		"userId",
		"meetingId",
		"userEnt",
		"meetEnt",
		"userDom",
		"userIsp",
		"relayDomIsp"})
	lossOrg = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "user",
		Subsystem: "quality",
		Name:      "lossOrg",
		Help:      "state",
	}, []string{
		"userId",
		"meetingId",
		"userEnt",
		"meetEnt",
		"userDom",
		"userIsp",
		"relayDomIsp"})
	lossOrgAvg = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "user",
		Subsystem: "quality",
		Name:      "lossOrgAvg",
		Help:      "state",
	}, []string{
		"userId",
		"meetingId",
		"userEnt",
		"meetEnt",
		"userDom",
		"userIsp",
		"relayDomIsp"})
	lossOrgStd = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "user",
		Subsystem: "quality",
		Name:      "lossOrgStd",
		Help:      "state",
	}, []string{
		"userId",
		"meetingId",
		"userEnt",
		"meetEnt",
		"userDom",
		"userIsp",
		"relayDomIsp"})
	prometheus.MustRegister(node)
	prometheus.MustRegister(synode)
	prometheus.MustRegister(nodenow)
	prometheus.MustRegister(nodenowcalltime)
	prometheus.MustRegister(nodecallip)
	prometheus.MustRegister(entcalltimemeet)
	prometheus.MustRegister(entipmeet)
	prometheus.MustRegister(meetrelay)
	prometheus.MustRegister(p2pactive)
	prometheus.MustRegister(meetactive)
	prometheus.MustRegister(crAvg)
	prometheus.MustRegister(crStd)
	prometheus.MustRegister(frAvg)
	prometheus.MustRegister(frStd)
	prometheus.MustRegister(delayAvg)
	prometheus.MustRegister(delayStd)
	prometheus.MustRegister(delayLoss)
	prometheus.MustRegister(dropline)
	prometheus.MustRegister(loss)
	prometheus.MustRegister(lossOrg)
	prometheus.MustRegister(lossOrgAvg)
	prometheus.MustRegister(lossOrgStd)

}

func main() {
	//判断活跃用户的模块

	ruc = &gproto.RecUserClient{}
	cnlogS := globeCfg.Output.CnlogS
	err := ruc.Init(cnlogS)
	fmt.Println("cnlogS: ", cnlogS)
	if err != nil {
		fmt.Println(err, "593")
	}
	defer ruc.Destroy()
	mip := globeCfg.Output.Ipmeet
	mdb := globeCfg.Output.Dbmeet
	mtable1 := globeCfg.Output.Table1meet

	url := globeCfg.Output.URL //readeuc:Hzl@20170920@tcp(192.168.101.17:3306)/enterpriseuc
	db, err := sql.Open("mysql", url+"?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	ip := globeCfg.Output.Ip
	mgodb := globeCfg.Output.Db
	session, err := mgo.Dial(ip)
	table1 := globeCfg.Output.Table1
	table2 := globeCfg.Output.Table2
	fmt.Println("p2pmongoDb配置", ip, mgodb, table1, table2)
	if err != nil {
		fmt.Println(err, "614")
		return
	}
	defer session.Close()

	collection := session.DB(mgodb).C(table1)
	collection2 := session.DB(mgodb).C(table2)

	//loop
	go func() {
		fmt.Println("Program startup ok...")
		//从mysql中获取企业号段初始化
		getNumber(db)

		fmt.Println("dbOK")
		for {
			//从mongodb中获取数据
			toPromtheus(mip, mdb, mtable1)
			getCall(collection)
			ispCityCall(collection2)
			Observe()
			//是否推送数据给PushGatway
			if globeCfg.Output.PushGateway {
				var info = make(map[string]string)
				info["monitorID"] = globeCfg.Output.MonitorID
				if err := push.FromGatherer("rt", info, globeCfg.Output.PushGatewayAddr, prometheus.DefaultGatherer); err != nil {
					fmt.Println("FromGatherer:", err)
				}
			}
			fmt.Println("****** ", time.Now(), " ******")
			fmt.Println()
			fmt.Println()
			time.Sleep(time.Duration(globeCfg.Output.Period) * time.Second)
		}
	}()
	//设置prometheus监听的ip和端口
	if globeCfg.Output.Prometheus {
		go func() {
			fmt.Println("ip", globeCfg.Gw.Addr)
			fmt.Println("port", globeCfg.Gw.HttpListenPort)
			http.Handle("/metrics", promhttp.Handler())
			http.ListenAndServe(fmt.Sprintf("%s:%d", globeCfg.Gw.Addr, globeCfg.Gw.HttpListenPort), nil)
		}()
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c)
	//	signal.Notify(c, os.Interrupt, os.Kill)
	s := <-c
	fmt.Println("exitss", s)
}
