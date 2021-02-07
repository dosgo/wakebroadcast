package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sabhiram/go-wol/wol"
	"log"
	"net"
	"wakebroadcast/notify"
	"os"
	"runtime"
)
func main(){
	address := "0.0.0.0:6661"
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
		}else{
			fmt.Printf("端口被占用程序已经在运行了!")
		}
		os.Exit(1)
	}
	defer conn.Close()
	go udpRecv(conn);
	//判断windows
	if(runtime.GOOS=="windows"){
		notify.GuiInit();
	}else{
		select{

		}
	}
}

func udpRecv(conn *net.UDPConn){
	for {
		// Here must use make and give the lenth of buffer
		data := make([]byte,1024)
		len, rAddr, err := conn.ReadFromUDP(data)
		if err != nil {
			fmt.Println(err)
			continue
		}
		jsonData := make(map[string]interface{})
		err = json.Unmarshal(data[:len], &jsonData)
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
		fmt.Println("wakeUp ---  mac:"+mac+"ip:"+ip)
		err=wakeUp(mac,ip,lip);
		var  code=0;
		if(err!=nil){
			log.Printf("err:%v\r\n",err)
			code=1;
		}
		back:=fmt.Sprintf("{\"code\":%d,\"msg\":\"\"}",code);
		conn.WriteTo([]byte(back),rAddr)
	}
}


func wakeUp(macAddr string, ip string,lip string) error{
	if(macAddr==""){
		return errors.New("mac null")
	}
	if(ip==""){
		ip="192.168.6.255"
	}
	if(lip==""){
		lip="192.168.6.221"
	}

	//生成魔术包结构，
	mp, err := wol.New(macAddr)
	if err != nil {
		return err;
	}
	bs, err := mp.Marshal()
	if err != nil {
		return err;
	}
	//laddr, err := net.ResolveUDPAddr("udp4", lip+":7777");
	// 这里设置接收者的IP地址为广播地址
	raddr, err := net.ResolveUDPAddr("udp4", ip+":7777");
	if err != nil {
		return err;
	}
	conn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		return err;
	}
	conn.Write(bs)
	conn.Close()
	return nil;
}


