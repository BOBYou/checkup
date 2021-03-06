package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/gaochao1/sw"
	"github.com/hel2o/checkup/g"
)

type Ping struct {
	Ip  string
	Rtt float64
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	g.RunPid()
	cfg := flag.String("c", "cfg.json", "configuration file")
	version := flag.Bool("v", false, "show version")
	flag.Parse()

	if *version {
		fmt.Println(g.VERSION)
		os.Exit(0)
	}

	g.MyLog()
	g.ParseConfig(*cfg)
	CheckupIps := AllCheckupIp()
	var FailsInterval int
	var HostName string
	if g.Config().Checkup.HostName == "" {
		HostName, _ = os.Hostname()
	} else {
		HostName = g.Config().Checkup.HostName
	}

	for {
		go Interval(CheckupIps, HostName, &FailsInterval)
		time.Sleep(time.Duration(g.Config().Checkup.Interval) * time.Second)
	}
}

func Interval(ips []string, hostName string, failsInterval *int) {
	t1 := time.Now()
	fails := 0
	iptables := CheckIptable()

	chs := make([]chan Ping, len(ips))

	for i, v := range ips {
		chs[i] = make(chan Ping)
		go GoPing(chs[i], v)
	}

	for _, ch := range chs {
		t2 := time.Now()
		Result := <-ch
		if Result.Rtt == -1 {
			fails++
		}
		log.Println("IP:", Result.Ip, "	Rtt:", Result.Rtt, "	Cost:", time.Since(t2))
	}
	f := float64(fails) / float64(len(ips))
	log.Println("iptables: ", iptables)

	if f >= g.Config().Checkup.FailureRate {
		*failsInterval++
		log.Println("FailsInterval:", *failsInterval)
		if *failsInterval >= g.Config().Checkup.FailsInterval && iptables == false {
			Iptables("sh/input")
			urlData_input := make(map[string][]string, 2)
			content_input := hostName + "\nPing次数：" + strconv.Itoa(len(ips)) + "\n失败次数：" + strconv.Itoa(fails) + "\n执行操作：Add Reject\n" + "时间：" + time.Now().Format("2006-01-02 15:04:05") + "\n以上内容通过SendWeChat_Api发送"
			urlData_input["content"] = []string{content_input}
			urlData_input["to"] = []string{g.Config().Checkup.To}

			_, err := g.Post(g.Config().Checkup.PostUrl, urlData_input)
			if err != nil {
				log.Println("error: %v", err)
			}
			log.Println("intput")
		}
	} else if f < g.Config().Checkup.FailureRate && iptables == true {

		Iptables("sh/remove")
		*failsInterval = 0
		urlData_remove := make(map[string][]string, 2)
		content_remove := hostName + "\nPing次数：" + strconv.Itoa(len(ips)) + "\n失败次数：" + strconv.Itoa(fails) + "\n执行操作：Remove Reject\n" + "时间：" + time.Now().Format("2006-01-02 15:04:05") + "\n以上内容通过SendWeChat_Api发送"
		urlData_remove["content"] = []string{content_remove}
		urlData_remove["to"] = []string{g.Config().Checkup.To}

		_, err := g.Post(g.Config().Checkup.PostUrl, urlData_remove)
		if err != nil {
			log.Println("error: %v", err)
		}
		log.Println("remove")
	} else {
		*failsInterval = 0
		log.Println("nothing to do!")
	}

	elapsed := time.Since(t1)
	log.Println("Runing time: ", elapsed, " Fails:", fails, " FailsInterval:", *failsInterval, "\n")
}

func GoPing(ch chan Ping, ip string) {
	var ping Ping
	fastPingMode := g.Config().Checkup.FastPingMode
	timeOut := g.Config().Checkup.PingTimeout * g.Config().Checkup.PingRetry
	rtt, err := sw.PingRtt(ip, timeOut, fastPingMode)

	if err != nil {
		ping.Rtt = -1
		ping.Ip = ip
		ch <- ping
		return
	}

	ping.Rtt = rtt
	ping.Ip = ip
	ch <- ping
	return
}

func AllCheckupIp() (allIp []string) {
	Checkup := g.Config().Checkup.IpRange

	if len(Checkup) > 0 {
		for _, sip := range Checkup {
			aip := sw.ParseIp(sip)
			for _, ip := range aip {
				allIp = append(allIp, ip)
			}
		}
	}
	return allIp
}

func Iptables(run string) {
	cmd := exec.Command("/bin/sh", run)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Println("error", err)
	}
	log.Println(string(out))
}

func CheckIptable() bool {
	cmd := exec.Command("iptables", "-L", "-n")
	out, _ := cmd.CombinedOutput()
	check := strings.Contains(string(out), "reject keepalived connect")
	return check
}
