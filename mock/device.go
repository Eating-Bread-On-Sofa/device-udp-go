package main

import (
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"time"
)

func main() {
	// 创建监听
	listen,err := net.ListenUDP("udp",&net.UDPAddr{
		IP:net.IPv4(0,0,0,0),
		Port:8888,
	})
	if err != nil{
		fmt.Println("监听失败，错误：",err)
		return
	}
	fmt.Println("listen Start...:")

	defer listen.Close()
	for {
		// 读取数据
		var data [1024]byte
		n,addr,err := listen.ReadFromUDP(data[:])
		if err != nil{
			fmt.Println("接收udp数据失败，错误：",err)
			continue
		}
		fmt.Printf("data:%v addr:%v count:%v\n", string(data[:n]), addr, n)

		// 返回数据
		senddata := string(data[:n])
		if senddata == "rand"{
			rand.Seed(time.Now().UnixNano())
			i := rand.Intn(128) + 1
			fmt.Println(i)
			str := strconv.Itoa(i)
			backdata := []byte(str)
			_, err = listen.WriteToUDP(backdata,addr)
			if err != nil{
				fmt.Println("发送数据失败，错误：",err)
				continue
			}
		}else if senddata == "ping"{
			backdata := []byte("pong")
			_, err = listen.WriteToUDP(backdata,addr)
			if err != nil{
				fmt.Println("发送数据失败，错误：",err)
				continue
			}
		}else{
			backdata := []byte("ok!")
			_, err = listen.WriteToUDP(backdata,addr)
			if err != nil{
				fmt.Println("发送数据失败，错误：",err)
				continue
			}
		}
	}
}
