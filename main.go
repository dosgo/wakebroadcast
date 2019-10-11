package main

import (
	"encoding/json"
	"fmt"
	"github.com/sabhiram/go-wol"
	"log"
	"net"
	"wakebroadcast/notify"
	"os"
	"runtime"
)
func main(){
	address := "0.0.0.0:666"
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println(err)
		if(runtime.GOOS=="windows") {
			notify.MsgBox("信息","端口被占用程序已经在运行了!");
		}
		os.Exit(1)
	}
	//判断windows
	if(runtime.GOOS=="windows"){
		notify.GuiInit();
	}
	defer conn.Close()
	for {
		// Here must use make and give the lenth of buffer
		data := make([]byte,1024)
		_, _, err := conn.ReadFromUDP(data)
		if err != nil {
			fmt.Println(err)
			continue
		}
		jsonData := make(map[string]interface{})
		err = json.Unmarshal(data, &jsonData)
		if err != nil {
			fmt.Println(err)
			continue
		}
		var mac  =""
		if value, ok := jsonData["mac"]; ok {
			_mac, ok1 := value.(string)
			if ok1 {
				mac=_mac;
			}
		}

		var ip  =""
		if value, ok := jsonData["ip"]; ok {
			_ip, ok1 := value.(string)
			if ok1 {
				ip=_ip;
			}
		}

		var lip  =""
		if value, ok := jsonData["lip"]; ok {
			_lip, ok1 := value.(string)
			if ok1 {
				lip=_lip;
			}
		}
		wakeUp(mac,ip,lip);
	}
}


func wakeUp(macAddr string, ip string,lip string) bool{
	if(macAddr==""){
		log.Println("mac null")
		return false
	}
	if(ip==""){
		ip="192.168.6.255"
	}
	if(lip==""){
		ip="192.168.6.221"
	}

	//生成魔术包结构，
	mp, err := wol.New(macAddr)
	if err != nil {
		return false;
	}
	bs, err := mp.Marshal()
	if err != nil {
		return false;
	}
	//laddr, err := net.ResolveUDPAddr("udp4", lip+":7777");
	// 这里设置接收者的IP地址为广播地址
	raddr, err := net.ResolveUDPAddr("udp4", ip+":9");
	conn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		println(err.Error())
		return false;
	}
	conn.Write(bs)
	conn.Close()
	return true;
}


