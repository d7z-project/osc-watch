package main

import (
	"flag"
	"log"
	"strings"
	"tinygo.org/x/bluetooth"
)

var (
	mac     string = ""
	host    string = "localhost"
	port    int    = 9000
	adapter        = bluetooth.DefaultAdapter
)

func init() {
	flag.StringVar(&mac, "mac", "", "过滤的MAC地址")
	flag.StringVar(&host, "addr", "localhost", "VRChat OSC 地址")
	flag.IntVar(&port, "port", 9000, "VRChat OSC 端口")
}

func main() {
	flag.Parse()
	mac = strings.ToUpper(mac)
	if mac == "" {
		log.Fatal("请指定蓝牙 MAC 地址")
	}

	_, err := bluetooth.ParseMAC(mac)
	if err != nil {
		log.Fatal("MAC 地址格式错误")
	}

	if err := adapter.Enable(); err != nil {
		log.Fatal("蓝牙开启失败", err)
	}
	log.Printf("开始监听来自 %s 的广播数据", mac)

	var r bluetooth.ScanResult
	_ = adapter.Scan(func(a *bluetooth.Adapter, result bluetooth.ScanResult) {
		if result.Address.String() == mac {
			log.Printf("find %s", result.LocalName())
			_ = adapter.StopScan()
			r = result
		}
	})

	connect, err := adapter.Connect(r.Address, bluetooth.ConnectionParams{})
	if err != nil {
		log.Fatalf("连接失败 ！ %v", err)
	}
	services, err := connect.DiscoverServices(nil)
	if err != nil {
		log.Fatal(err)
	}
	for _, service := range services {
		println(service.UUID().String())
	}

}
