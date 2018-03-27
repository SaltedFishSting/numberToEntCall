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
	"github.com/gin-gonic/gin"
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
		Devstream       map[string]string `yaml:"devstream"`
	}
}

var devstream map[string]string

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
	Id struct {
		Userid     string
		Speakerid  string
		Resourceid int
	} "_id"
	Value struct {
		Userip      string
		Relayip     string
		Userent     string
		Meetent     string
		Userdom     string
		Userisp     string
		Relaydomisp string
		Meetingid   float64
		Devicetype  int
		Networktype int
		Resourceid  int
		Crarr       string
		Frarr       string
		Lossorgarr  string
		Loss        float64
		Delayarr    string
		Delayloss   int
		Dropline    bool
		Users       float64
		Oneempty    int
		Twoempty    int
		Threeempty  int
		Fourempty   int
		Tenempty    int
		Usertype    string
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
type userkey struct {
	Userid     string
	Meetintid  string
	Speakerid  string
	Relay      string
	Resourceid int
}

var uservideoqualityup map[userkey]uservideoinfo
var uservideoqualitydown map[userkey]uservideoinfo
var uservideoqualityaudioup map[userkey]uservideoinfo
var uservideoqualityaudiodown map[userkey]uservideoinfo
var userfilequalityup map[userkey]uservideoinfo
var userfilequalitydown map[userkey]uservideoinfo
var devmeet map[string]string
var netmeet map[string]string
var maxdelayup float64
var maxdelaydown float64

type usertime struct {
	userid    string
	entname   string
	starttime time.Time
	endtime   time.Time
}

//会议合格率
var meetqualified map[int]usermeet

var meetc2cqualifiedmap map[userc2c]float64

//prometheus var
var (
	nodenow           *(prometheus.GaugeVec)
	nodenowcalltime   *(prometheus.GaugeVec)
	nodecallip        *(prometheus.GaugeVec)
	node              *(prometheus.GaugeVec)
	synode            *(prometheus.GaugeVec)
	entcalltimemeet   *(prometheus.GaugeVec)
	entipmeet         *(prometheus.GaugeVec)
	meetrelay         *(prometheus.GaugeVec)
	p2pactive         *(prometheus.GaugeVec)
	meetactive        *(prometheus.GaugeVec)
	filecrAvg         *(prometheus.GaugeVec)
	filecrStd         *(prometheus.GaugeVec)
	filefrAvg         *(prometheus.GaugeVec)
	filefrStd         *(prometheus.GaugeVec)
	filedelayAvg      *(prometheus.GaugeVec)
	filedelayStd      *(prometheus.GaugeVec)
	filedelayLoss     *(prometheus.GaugeVec)
	filedropline      *(prometheus.GaugeVec)
	fileloss          *(prometheus.GaugeVec)
	filelossOrg       *(prometheus.GaugeVec)
	filelossOrgAvg    *(prometheus.GaugeVec)
	filelossOrgStd    *(prometheus.GaugeVec)
	crAvg             *(prometheus.GaugeVec)
	crStd             *(prometheus.GaugeVec)
	frAvg             *(prometheus.GaugeVec)
	frStd             *(prometheus.GaugeVec)
	delayAvg          *(prometheus.GaugeVec)
	delayStd          *(prometheus.GaugeVec)
	delayLoss         *(prometheus.GaugeVec)
	dropline          *(prometheus.GaugeVec)
	loss              *(prometheus.GaugeVec)
	lossOrg           *(prometheus.GaugeVec)
	lossOrgAvg        *(prometheus.GaugeVec)
	lossOrgStd        *(prometheus.GaugeVec)
	Emptyaudiobag     *(prometheus.GaugeVec)
	audioloss         *(prometheus.GaugeVec)
	audiolossAvg      *(prometheus.GaugeVec)
	audiolossStd      *(prometheus.GaugeVec)
	audiolossOrg      *(prometheus.GaugeVec)
	videomaxdelay     *(prometheus.GaugeVec)
	filevideomaxdelay *(prometheus.GaugeVec)

	//会议合格率
	meetqualifiedrate *(prometheus.GaugeVec)
	//会议c2c合格率和cpu使用情况
	meetc2cqualified *(prometheus.GaugeVec)
	//解析excal后生成的map
	usermap map[int64]string
	//临时usermap
	temporaryUser map[int]usertime
	//临时usermap的key
	keyid int
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
	uservideoqualityup = make(map[userkey]uservideoinfo)
	uservideoqualitydown = make(map[userkey]uservideoinfo)
	uservideoqualityaudioup = make(map[userkey]uservideoinfo)
	uservideoqualityaudiodown = make(map[userkey]uservideoinfo)
	userfilequalityup = make(map[userkey]uservideoinfo)
	userfilequalitydown = make(map[userkey]uservideoinfo)
	temporaryUser = make(map[int]usertime)
	meetqualified = make(map[int]usermeet)
	meetc2cqualifiedmap = make(map[userc2c]float64)
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

//加载用户excel文件
func loadexcel() {
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
	timenow := time.Now().Unix()
	for _, v := range temporaryUser {
		if timenow > v.starttime.Unix() && timenow < v.endtime.Unix() {
			uid, _ := strconv.Atoi(v.userid)
			usermap[int64(uid)] = v.entname
		}
		if timenow > v.endtime.Unix() {
			uid, _ := strconv.Atoi(v.userid)
			delete(usermap, int64(uid))
		}
	}

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

}
func init() {
	loadConfig() //加载配置文件
	devstream = globeCfg.Output.Devstream
	devmeet = globeCfg.Output.Devmeet
	netmeet = globeCfg.Output.Netmeet
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
		"Speakerid",
		"meetingId",
		"deviceType",
		"networkType",
		"userEnt",
		"meetEnt",
		"userDom",
		"userIsp",
		"relayDomIsp",
		"updown",
		"userType"})
	filecrAvg = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "user",
		Subsystem: "quality",
		Name:      "filecrAvg",
		Help:      "state",
	}, []string{
		"userId",
		"Speakerid",
		"meetingId",
		"deviceType",
		"networkType",
		"userEnt",
		"meetEnt",
		"userDom",
		"userIsp",
		"relayDomIsp",
		"updown",
		"userType"})
	crStd = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "user",
		Subsystem: "quality",
		Name:      "crStd",
		Help:      "state",
	}, []string{
		"userId",
		"Speakerid",
		"meetingId",
		"deviceType",
		"networkType",
		"userEnt",
		"meetEnt",
		"userDom",
		"userIsp",
		"relayDomIsp",
		"updown",
		"userType"})
	filecrStd = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "user",
		Subsystem: "quality",
		Name:      "filecrStd",
		Help:      "state",
	}, []string{
		"userId",
		"Speakerid",
		"meetingId",
		"deviceType",
		"networkType",
		"userEnt",
		"meetEnt",
		"userDom",
		"userIsp",
		"relayDomIsp",
		"updown",
		"userType"})
	frAvg = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "user",
		Subsystem: "quality",
		Name:      "frAvg",
		Help:      "state",
	}, []string{
		"userId",
		"Speakerid",
		"meetingId",
		"deviceType",
		"networkType",
		"userEnt",
		"meetEnt",
		"userDom",
		"userIsp",
		"relayDomIsp",
		"updown",
		"userType"})
	filefrAvg = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "user",
		Subsystem: "quality",
		Name:      "filefrAvg",
		Help:      "state",
	}, []string{
		"userId",
		"Speakerid",
		"meetingId",
		"deviceType",
		"networkType",
		"userEnt",
		"meetEnt",
		"userDom",
		"userIsp",
		"relayDomIsp",
		"updown",
		"userType"})
	frStd = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "user",
		Subsystem: "quality",
		Name:      "frStd",
		Help:      "state",
	}, []string{
		"userId",
		"Speakerid",
		"meetingId",
		"deviceType",
		"networkType",
		"userEnt",
		"meetEnt",
		"userDom",
		"userIsp",
		"relayDomIsp",
		"updown",
		"userType"})
	filefrStd = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "user",
		Subsystem: "quality",
		Name:      "filefrStd",
		Help:      "state",
	}, []string{
		"userId",
		"Speakerid",
		"meetingId",
		"deviceType",
		"networkType",
		"userEnt",
		"meetEnt",
		"userDom",
		"userIsp",
		"relayDomIsp",
		"updown",
		"userType"})
	delayAvg = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "user",
		Subsystem: "quality",
		Name:      "delayAvg",
		Help:      "state",
	}, []string{
		"userId",
		"Speakerid",
		"meetingId",
		"deviceType",
		"networkType",
		"userEnt",
		"meetEnt",
		"userDom",
		"userIsp",
		"relayDomIsp",
		"updown",
		"userType"})
	filedelayAvg = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "user",
		Subsystem: "quality",
		Name:      "filedelayAvg",
		Help:      "state",
	}, []string{
		"userId",
		"Speakerid",
		"meetingId",
		"deviceType",
		"networkType",
		"userEnt",
		"meetEnt",
		"userDom",
		"userIsp",
		"relayDomIsp",
		"updown",
		"userType"})
	delayStd = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "user",
		Subsystem: "quality",
		Name:      "delayStd",
		Help:      "state",
	}, []string{
		"userId",
		"Speakerid",
		"meetingId",
		"deviceType",
		"networkType",
		"userEnt",
		"meetEnt",
		"userDom",
		"userIsp",
		"relayDomIsp",
		"updown",
		"userType"})
	filedelayStd = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "user",
		Subsystem: "quality",
		Name:      "filedelayStd",
		Help:      "state",
	}, []string{
		"userId",
		"Speakerid",
		"meetingId",
		"deviceType",
		"networkType",
		"userEnt",
		"meetEnt",
		"userDom",
		"userIsp",
		"relayDomIsp",
		"updown",
		"userType"})
	delayLoss = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "user",
		Subsystem: "quality",
		Name:      "delayLoss",
		Help:      "state",
	}, []string{
		"userId",
		"Speakerid",
		"meetingId",
		"deviceType",
		"networkType",
		"userEnt",
		"meetEnt",
		"userDom",
		"userIsp",
		"relayDomIsp",
		"updown",
		"userType"})
	filedelayLoss = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "user",
		Subsystem: "quality",
		Name:      "filedelayLoss",
		Help:      "state",
	}, []string{
		"userId",
		"Speakerid",
		"meetingId",
		"deviceType",
		"networkType",
		"userEnt",
		"meetEnt",
		"userDom",
		"userIsp",
		"relayDomIsp",
		"updown",
		"userType"})
	dropline = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "user",
		Subsystem: "quality",
		Name:      "dropline",
		Help:      "state",
	}, []string{
		"userId",
		"Speakerid",
		"meetingId",
		"deviceType",
		"networkType",
		"userEnt",
		"meetEnt",
		"userDom",
		"userIsp",
		"relayDomIsp",
		"updown",
		"userType"})
	filedropline = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "user",
		Subsystem: "quality",
		Name:      "filedropline",
		Help:      "state",
	}, []string{
		"userId",
		"Speakerid",
		"meetingId",
		"deviceType",
		"networkType",
		"userEnt",
		"meetEnt",
		"userDom",
		"userIsp",
		"relayDomIsp",
		"updown",
		"userType"})
	loss = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "user",
		Subsystem: "quality",
		Name:      "loss",
		Help:      "state",
	}, []string{
		"userId",
		"Speakerid",
		"meetingId",
		"deviceType",
		"networkType",
		"userEnt",
		"meetEnt",
		"userDom",
		"userIsp",
		"relayDomIsp",
		"updown",
		"userType"})
	fileloss = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "user",
		Subsystem: "quality",
		Name:      "fileloss",
		Help:      "state",
	}, []string{
		"userId",
		"Speakerid",
		"meetingId",
		"deviceType",
		"networkType",
		"userEnt",
		"meetEnt",
		"userDom",
		"userIsp",
		"relayDomIsp",
		"updown",
		"userType"})
	lossOrg = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "user",
		Subsystem: "quality",
		Name:      "lossOrg",
		Help:      "state",
	}, []string{
		"userId",
		"Speakerid",
		"meetingId",
		"deviceType",
		"networkType",
		"userEnt",
		"meetEnt",
		"userDom",
		"userIsp",
		"relayDomIsp",
		"updown",
		"userType"})
	filelossOrg = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "user",
		Subsystem: "quality",
		Name:      "filelossOrg",
		Help:      "state",
	}, []string{
		"userId",
		"Speakerid",
		"meetingId",
		"deviceType",
		"networkType",
		"userEnt",
		"meetEnt",
		"userDom",
		"userIsp",
		"relayDomIsp",
		"updown",
		"userType"})
	lossOrgAvg = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "user",
		Subsystem: "quality",
		Name:      "lossOrgAvg",
		Help:      "state",
	}, []string{
		"userId",
		"Speakerid",
		"meetingId",
		"deviceType",
		"networkType",
		"userEnt",
		"meetEnt",
		"userDom",
		"userIsp",
		"relayDomIsp",
		"updown",
		"userType"})
	filelossOrgAvg = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "user",
		Subsystem: "quality",
		Name:      "filelossOrgAvg",
		Help:      "state",
	}, []string{
		"userId",
		"Speakerid",
		"meetingId",
		"deviceType",
		"networkType",
		"userEnt",
		"meetEnt",
		"userDom",
		"userIsp",
		"relayDomIsp",
		"updown",
		"userType"})
	lossOrgStd = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "user",
		Subsystem: "quality",
		Name:      "lossOrgStd",
		Help:      "state",
	}, []string{
		"userId",
		"Speakerid",
		"meetingId",
		"deviceType",
		"networkType",
		"userEnt",
		"meetEnt",
		"userDom",
		"userIsp",
		"relayDomIsp",
		"updown",
		"userType"})
	filelossOrgStd = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "user",
		Subsystem: "quality",
		Name:      "filelossOrgStd",
		Help:      "state",
	}, []string{
		"userId",
		"Speakerid",
		"meetingId",
		"deviceType",
		"networkType",
		"userEnt",
		"meetEnt",
		"userDom",
		"userIsp",
		"relayDomIsp",
		"updown",
		"userType"})
	Emptyaudiobag = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "user",
		Subsystem: "quality",
		Name:      "Emptyaudiobag",
		Help:      "state",
	}, []string{
		"userId",
		"Speakerid",
		"meetingId",
		"deviceType",
		"networkType",
		"userEnt",
		"meetEnt",
		"userDom",
		"userIsp",
		"relayDomIsp",
		"updown",
		"emptygrades",
		"userType"})
	audiolossOrg = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "user",
		Subsystem: "quality",
		Name:      "audiolossOrg",
		Help:      "state",
	}, []string{
		"userId",
		"Speakerid",
		"meetingId",
		"deviceType",
		"networkType",
		"userEnt",
		"meetEnt",
		"userDom",
		"userIsp",
		"relayDomIsp",
		"updown",
		"userType"})
	audiolossAvg = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "user",
		Subsystem: "quality",
		Name:      "audiolossAvg",
		Help:      "state",
	}, []string{
		"userId",
		"Speakerid",
		"meetingId",
		"deviceType",
		"networkType",
		"userEnt",
		"meetEnt",
		"userDom",
		"userIsp",
		"relayDomIsp",
		"updown",
		"userType"})
	audiolossStd = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "user",
		Subsystem: "quality",
		Name:      "audiolossStd",
		Help:      "state",
	}, []string{
		"userId",
		"Speakerid",
		"meetingId",
		"deviceType",
		"networkType",
		"userEnt",
		"meetEnt",
		"userDom",
		"userIsp",
		"relayDomIsp",
		"updown",
		"userType"})
	audioloss = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "user",
		Subsystem: "quality",
		Name:      "audioloss",
		Help:      "state",
	}, []string{
		"userId",
		"Speakerid",
		"meetingId",
		"deviceType",
		"networkType",
		"userEnt",
		"meetEnt",
		"userDom",
		"userIsp",
		"relayDomIsp",
		"updown",
		"userType"})
	meetqualifiedrate = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "meet",
		Subsystem: "quality",
		Name:      "qualified",
		Help:      "state",
	}, []string{
		"meetingId",
		"userType"})
	meetc2cqualified = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "meet",
		Subsystem: "quality",
		Name:      "c2c",
		Help:      "state",
	}, []string{
		"caller",
		"called",
		"meetingId",
		"qualified",
		"cpu",
		"userType"})
	videomaxdelay = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "user",
		Subsystem: "quality",
		Name:      "videomaxdelay",
		Help:      "state",
	}, []string{
		"userId",
		"Speakerid",
		"meetingId",
		"deviceType",
		"networkType",
		"userEnt",
		"meetEnt",
		"userDom",
		"userIsp",
		"relayDomIsp",
		"updown",
		"userType"})
	filevideomaxdelay = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "user",
		Subsystem: "quality",
		Name:      "filevideomaxdelay",
		Help:      "state",
	}, []string{
		"userId",
		"Speakerid",
		"meetingId",
		"deviceType",
		"networkType",
		"userEnt",
		"meetEnt",
		"userDom",
		"userIsp",
		"relayDomIsp",
		"updown",
		"userType"})

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
	prometheus.MustRegister(filecrAvg)
	prometheus.MustRegister(filecrStd)
	prometheus.MustRegister(filefrAvg)
	prometheus.MustRegister(filefrStd)
	prometheus.MustRegister(filedelayAvg)
	prometheus.MustRegister(filedelayStd)
	prometheus.MustRegister(filedelayLoss)
	prometheus.MustRegister(filedropline)
	prometheus.MustRegister(fileloss)
	prometheus.MustRegister(filelossOrg)
	prometheus.MustRegister(filelossOrgAvg)
	prometheus.MustRegister(filelossOrgStd)
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
	prometheus.MustRegister(Emptyaudiobag)
	prometheus.MustRegister(audiolossOrg)
	prometheus.MustRegister(audiolossAvg)
	prometheus.MustRegister(audiolossStd)
	prometheus.MustRegister(audioloss)
	prometheus.MustRegister(meetqualifiedrate)
	prometheus.MustRegister(meetc2cqualified)
	prometheus.MustRegister(videomaxdelay)
	prometheus.MustRegister(filevideomaxdelay)

}

func main() {

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

			//加载用户excel文件
			loadexcel()
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
	//开启aduser接口端口10088
	go func() {
		router := gin.Default()
		router.GET("adduser", adduser)
		router.GET("showuser", showuser)
		router.GET("deleteuser", deleteuser)
		router.GET("updateuser", updateuser)
		router.Run(":10088")
	}()
	c := make(chan os.Signal, 1)
	signal.Notify(c)
	//	signal.Notify(c, os.Interrupt, os.Kill)
	s := <-c
	fmt.Println("exitss", s)
}
