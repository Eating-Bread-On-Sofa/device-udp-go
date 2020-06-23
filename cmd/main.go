// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2018-2019 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"github.com/edgexfoundry/device-sdk-go/pkg/startup"
	"github.com/edgexfoundry/device-udp-go"
	"github.com/edgexfoundry/device-udp-go/driver"
)

const (
	serviceName string = "edgex-device-udp"
)

func main() {
	d := driver.NewProtocolDriver()
	startup.Bootstrap(serviceName, device_udp.Version, d)
}
