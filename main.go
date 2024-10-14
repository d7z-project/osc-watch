package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/hypebeast/go-osc/osc"
	"tinygo.org/x/bluetooth"
)

var (
	mac     = ""
	host    = "localhost"
	port    = 9000
	adapter = bluetooth.DefaultAdapter
)

func init() {
	flag.StringVar(&mac, "mac", "", "过滤的MAC地址")
	flag.StringVar(&host, "addr", "10.0.3.115", "VRChat OSC 地址")
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
	atomicBpm := ""

	go func() {
		for {
			if atomicBpm == "" {
				continue
			}

			message := osc.NewMessage("/chatbox/input")
			message.Append(strings.TrimSpace(atomicBpm))
			message.Append(true)
			message.Append(true)
			_ = oscClient.Send(message)
			time.Sleep(2 * time.Second)
		}
	}()

	go func() {
		<-ctx.Done()
		_ = adapter.StopScan()
	}()

	log.Printf("开始监听来自 %s 的广播数据", mac)
	maxBPM := 0
	if err := adapter.Scan(func(adapter *bluetooth.Adapter, device bluetooth.ScanResult) {
		if device.Address.MAC.String() == mac {
			for _, data := range device.ManufacturerData() {
				if data.CompanyID == 0x0157 {
					bpm := int(data.Data[3])

					if bpm == 255 {
						log.Printf("当前心率不正常 ( == 255 )，可能未开启小米手环运动模式")
						atomicBpm = "小米手环离线了"
						continue
					} else {
						log.Printf("[ %d dBm] 当前手环心率为 %d BPM", device.RSSI, bpm)
					}
					if maxBPM < bpm {
						maxBPM = bpm
					}
					note := "不想动 orz"
					if 130 >= bpm && bpm > 100 {
						note = "动起来了"
					}
					if 150 >= bpm && bpm > 130 {
						note = "火力全开"
					}
					if bpm > 150 {
						note = "感觉快寄了 orz"
					}
					atomicBpm = fmt.Sprintf("BLE RSSI: %d dBm\nHeart Rate: %03d / %03d BPM\n\n%s", device.RSSI, bpm, maxBPM, note)
				}
			}
		}
	}); err != nil {
		log.Printf("%v", err)
	}
	log.Printf("任务结束")
}
