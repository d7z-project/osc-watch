package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"strings"
	"syscall"

	"github.com/hypebeast/go-osc/osc"
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

	oscClient := osc.NewClient(host, port)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM)
	defer stop()

	if err := adapter.Enable(); err != nil {
		log.Fatal("蓝牙开启失败", err)
	}

	go func() {
		<-ctx.Done()
		adapter.StopScan()
	}()

	log.Printf("开始监听来自 %s 的广播数据", mac)

	if err := adapter.Scan(func(adapter *bluetooth.Adapter, device bluetooth.ScanResult) {
		if device.Address.MAC.String() == mac {
			for _, data := range device.ManufacturerData() {
				if data.CompanyID == 0x0157 {
					bpm := data.Data[3]
					if bpm == 255 {
						log.Printf("当前心率不正常 ( == 255 )，可能未开启小米手环运动模式")
						continue
					} else {
						log.Printf("[ %d dBm] 当前手环心率为 %d BPM", device.RSSI, bpm)
					}
					msg1 := osc.NewMessage("/avatar/parameters/Heartrate")
					msg1.Append(float32(float32(bpm)/127.0 - 1.0))
					_ = oscClient.Send(msg1)
					msg2 := osc.NewMessage("/avatar/parameters/Heartrate2")
					msg2.Append(float32(float32(bpm) / 127.0))
					_ = oscClient.Send(msg2)
					msg3 := osc.NewMessage("/avatar/parameters/Heartrate3")
					msg3.Append(int32(bpm))
					_ = oscClient.Send(msg3)
				}
			}
		}
	}); err != nil {
		log.Printf("%v", err)
	}
	log.Printf("任务结束")
}

func must(action string, err error) {
	if err != nil {
		panic("failed to " + action + ": " + err.Error())
	}
}
