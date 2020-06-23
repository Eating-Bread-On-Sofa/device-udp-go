// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2018 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

// This package provides a implementation of a ProtocolDriver interface.
//
package driver

import (
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"

	dsModels "github.com/edgexfoundry/device-sdk-go/pkg/models"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
	"github.com/edgexfoundry/go-mod-core-contracts/models"
)

var once sync.Once
var driver *Driver

type Driver struct {
	lc            logger.LoggingClient
	asyncCh       chan<- *dsModels.AsyncValues
}

func NewProtocolDriver() dsModels.ProtocolDriver {
	once.Do(func() {
		driver = new(Driver)
	})
	return driver
}

func (d *Driver) DisconnectDevice(deviceName string, protocols map[string]models.ProtocolProperties) error {
	d.lc.Info(fmt.Sprintf("Driver.DisconnectDevice: device-udp-go driver is disconnecting to %s", deviceName))
	return nil
}

func (d *Driver) Initialize(lc logger.LoggingClient, asyncCh chan<- *dsModels.AsyncValues) error {
	d.lc = lc
	d.asyncCh = asyncCh
	return nil
}

func (d *Driver) HandleReadCommands(deviceName string, protocols map[string]models.ProtocolProperties, reqs []dsModels.CommandRequest) (res []*dsModels.CommandValue, err error) {
	if len(reqs) != 1 {
		err = fmt.Errorf("Driver.HandleReadCommands; too many command requests; only one supported")
		return
	}

	d.lc.Debug(fmt.Sprintf("Driver.HandleReadCommands: protocols: %v resource: %v attributes: %v", protocols, reqs[0].DeviceResourceName, reqs[0].Attributes))

	res = make([]*dsModels.CommandValue, 1)
	now := time.Now().UnixNano()
	data := make([]byte, 4096)
	var UDP = protocols["udp"]
	addr := UDP["Address"]
	// 创建连接
	client, err := net.Dial("udp", addr)
	if err != nil {
		fmt.Println("连接失败!", err)
	}
	defer client.Close()

	if reqs[0].DeviceResourceName == "randomnumber" {
		// 发送数据
		senddata := []byte("rand")
		_, err = client.Write(senddata)
		if err != nil {
			fmt.Println("发送数据失败!", err)
		}
		// 接收数据
		n, err := client.Read(data[:])
		if err != nil {
			fmt.Println("读取数据失败!", err)
		}
		//将字节数组转化为string类型的
		var t=string(data[:n])
		num ,_ := strconv.Atoi(t)
		cv, _ := dsModels.NewInt32Value(reqs[0].DeviceResourceName, now, int32(num))
		res[0] = cv
	} else if reqs[0].DeviceResourceName == "ping" {
		// 发送数据
		senddata := []byte("ping")
		_, err = client.Write(senddata)
		if err != nil {
			fmt.Println("发送数据失败!", err)
		}
		// 接收数据
		n, err := client.Read(data[:])
		if err != nil {
			fmt.Println("读取数据失败!", err)
		}
		//将字节数组转化为string类型的
		var t=string(data[:n])
		cv := dsModels.NewStringValue(reqs[0].DeviceResourceName, now, t)
		res[0] = cv
	}else if reqs[0].DeviceResourceName == "message"{
		senddata := []byte("message")
		_, err = client.Write(senddata)
		if err != nil {
			fmt.Println("发送数据失败!", err)
		}
		// 接收数据
		n, err := client.Read(data[:])
		if err != nil {
			fmt.Println("读取数据失败!", err)
		}
		//将字节数组转化为string类型的
		var t=string(data[:n])
		cv := dsModels.NewStringValue(reqs[0].DeviceResourceName, now, t)
		res[0] = cv
	} else{
		err = fmt.Errorf("don't find deviceresource")
	}

	return
}

func (d *Driver) HandleWriteCommands(deviceName string, protocols map[string]models.ProtocolProperties, reqs []dsModels.CommandRequest,
	params []*dsModels.CommandValue) error {
	if len(reqs) != 1 {
		err := fmt.Errorf("Driver.HandleWriteCommands; too many command requests; only one supported")
		return err
	}
	if len(params) != 1 {
		err := fmt.Errorf("Driver.HandleWriteCommands; the number of parameter is not correct; only one supported")
		return err
	}

	d.lc.Debug(fmt.Sprintf("Driver.HandleWriteCommands: protocols: %v, resource: %v, parameters: %v", protocols, reqs[0].DeviceResourceName, params))
	var err error

	data := make([]byte, 4096)
	var UDP = protocols["udp"]
	addr := UDP["Address"]
	// 创建连接
	client, err := net.Dial("udp", addr)
	if err != nil {
		fmt.Println("连接失败!", err)
	}
	defer client.Close()
	if params[0].DeviceResourceName == "message" {
		senddata, err := params[0].StringValue();
		if err != nil{
			return fmt.Errorf("Driver.HandleWriteCommands: %v", err)
		}
		data = []byte(senddata)
		_, err = client.Write(data)
		if err != nil{
			return fmt.Errorf("set failed: %v", err)
		}
	}else{
		return fmt.Errorf("Driver.HandleWriteCommands: there is no matched device resource for %s", params[0].String())
	}

	return nil
}

func (d *Driver) Stop(force bool) error {
	d.lc.Info("Driver.Stop: device-udp-go driver is stopping...")
	return nil
}

func (d *Driver) AddDevice(deviceName string, protocols map[string]models.ProtocolProperties, adminState models.AdminState) error {
	d.lc.Debug(fmt.Sprintf("a new Device is added: %s", deviceName))
	return nil
}

func (d *Driver) UpdateDevice(deviceName string, protocols map[string]models.ProtocolProperties, adminState models.AdminState) error {
	d.lc.Debug(fmt.Sprintf("Device %s is updated", deviceName))
	return nil
}

func (d *Driver) RemoveDevice(deviceName string, protocols map[string]models.ProtocolProperties) error {
	d.lc.Debug(fmt.Sprintf("Device %s is removed", deviceName))
	return nil
}
