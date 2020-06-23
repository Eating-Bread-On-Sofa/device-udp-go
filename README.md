# Device-Udp-Go
## 使用前须安装
在使用设备服务前需要安装以下依赖：
- Ubuntu的Linux系统上必须安装Go Lang（1.13或者更高版本）。
- 安装“make”程序，之后会用make指令生成可执行二进制文件。通过如下指令：
```
sudo apt install build-essential
```
- 推荐安装一个Go语言编辑器，如GoLand。（可自选，不安装也行，用于开发Go语言代码）
## 相关代码说明
在本设备服务的的cmd/res目录下，有一个configuration.toml文件，cmd/res/docker下也有一个configuration.toml文件。
docker目录下的configuration.toml主要用于生成镜像，开启容器时使用，可以不做处理。目前，主要通过生成可执行文件运行设备服务代码，
因此可以根据需要修改cmd/res目录下的configuration.toml文件。
- 其中最上面的[service]部分主要是说明本设备服务的IP地址，可根据需要修改“Host”和“Port”里面的内容，别的暂时不需要修改：
```cassandraql
[Service]
Host = "0.0.0.0"   //这部分改成运行设备服务的IP地址
Port = 49990       //端口号可以根据需要修改
ConnectRetries = 3
Labels = []
OpenMsg = "udp device started"
MaxResultCount = 50000
Timeout = 5000
EnableAsyncReadings = true
AsyncBufferSize = 16
```
- 下面的[Registry]和[client]是将设备服务注册到控制台和edgexfoundry。（设备服务这个微服务可以不和edgexfoundry跑在同一台服务器上，这就需要这部分代码将DeviceService和edgexfoundry联系起来）
```cassandraql
[Registry]
Host = "0.0.0.0"   \\这部分为edgexfoundry运行的IP地址，注册后，可以通过8500端口到控制台查看本微服务。（下面端口不可更改）
Port = 8500
CheckInterval = "10s"
FailLimit = 3
FailWaitTime = 10
Type = "consul"   \\本条语句不能缺少 

[Clients]
  [Clients.Data]
  Name = "edgex-core-data"
  Protocol = "http"
  Host = "0.0.0.0"   \\这部分为edgexfoundry运行的IP地址，注册后，DeviceService接受的数据将会推送到coredata模块保存（下面端口不可更改）
  Port = 48080
  Timeout = 5000

  [Clients.Metadata]
  Name = "edgex-core-metadata"
  Protocol = "http"
  Host = "0.0.0.0"  \\这部分为edgexfoundry运行的IP地址，注册后，DeviceService将会连接到metadata模块（下面端口不可更改）
  Port = 48081
  Timeout = 5000

  [Clients.Logging]
  Name = "edgex-support-logging"
  Protocol = "http"
  Host = "0.0.0.0"  \\这部分为edgexfoundry运行的IP地址，注册后，DeviceService将会连接到logging模块（下面端口不可更改）
  Port = 48061
```
- 接下来的[device]部分，暂时还没弄懂，官网给的每个微服务都有这部分代码，可以不用管他，这部分代码不变。[logging]部分是将日志保存的位置，以及保存什么级别及以上的日志（debug、info、warn等），根据需要自己调整。
- 最后，[[DeviceList]]这部分对于设备服务添加所属设备非常重要。一个设备服务可以添加多个按此协议通讯的设备，其中一种方式是在以下部分进行注册：
```cassandraql
[[DeviceList]]   //一个设备要想注册到本设备服务，需要添加如下内容，若想注册多个设备，可以写多个[[DeviceList]]
  Name = "Udp-device01"   //设备的名字，可以随意起，但不同设备不要重复
  Profile = "Udp-Device"  //说明设备遵从哪个.yml文件。（如本设备服务里面cmd/res下面提供的udp-test-device.yaml）。名字必须为.yml文件最开头“name”部分的名字，绝对不能写错
  Description = "Example of Udp Device"  //描述，可以随便写，没有要求。这部分对应edgexfoundry中UI界面设备的相关备注信息
  Labels = [ "industrial" ]   //标签，可以随便写，没有要求。这部分对应edgexfoundry中UI界面设备的相关备注信息
  [DeviceList.Protocols]
    [DeviceList.Protocols.udp]  //本协议这部分最后写的udp，当开发其它协议时，这里也可以写别的，写driver函数时根据这部分来定位协议和设备
      Address = "192.168.0.141:8888"  //这个根据协议具体需要，自定义要编辑的内容，如mqtt还要指定订阅主体topic等。
  [[DeviceList.AutoEvents]]  //这部分及以下可选写，也可以写多个。用于定时向设备读取数据，并将数据推送到edgexfoundry的coredata里面。
    Frequency = "30s"   //向设备读取数据的频率
    OnChange = false    //这个值等于“false”，表明每隔上述间隔，必须要向coredata上传数据，等于“true”时，表明只有接受的数据发生改变时，才会向coredata推送数据
    Resource = "randomnumber"  //这个值表明到底想读设备的哪个资源，必须是.yaml文件里面deviceResources下面“name”的值，不能随便写
```
上述提到的.yaml文件为设备配置文件，里面指明了edgexfoundy可下达的coreCommands有哪些，里面的deviceCommands是将coreCommands转换成设备可以理解的命令。deviceResources里面写出了设备可以提供哪些数据。
```cassandraql
name: "Udp-Device"     //必须与configuration.toml里面[[DeviceList]]下的Name名称对应 
manufacturer: "Dell Technologies"  //不太清楚，不用管
model: "1"   //不太清楚，不用管
labels:
 - "test"   //标签，可以随便写，没有要求。这部分对应edgexfoundry中UI界面设备的相关备注信息
description: "simulate a device"  ////描述，可以随便写，没有要求。这部分对应edgexfoundry中UI界面设备的相关备注信息

deviceResources:       //定义设备具有哪些可调用的资源
    -   
        name: "randomnumber"
        description: "get random number"
        attributes:
            { type: "random" }
        properties:     //关于传输数据的相关性质，根据需要定义
            value:
                { type: "INT32", readWrite: "R", defaultValue: "0.00", minimum: "0.00", maximum: "100.00"  }
            units:
                { type: "String", readWrite: "R", defaultValue: "" }
    -
        name: "ping"
        description: "device awake"
        properties:
            value:
                { type: "String", size: "0", readWrite: "R", defaultValue: "oops" }
            units:
                { type: "String", readWrite: "R", defaultValue: "" }
    -
        name: "message"
        description: "device notification message"
        properties:
            value:
                { type: "String", size: "0", readWrite: "W" ,scale: "", offset: "", base: ""  }
            units:
                { type: "String", readWrite: "R", defaultValue: "" }

deviceCommands:    //用于将coreCommands里面的命令转换成设备理解的命令
    -
        name: "Random"   //必须与coreCommands里面的name值保持一致
        get:
            -
                { operation: "get", object: "randomnumber", property: "value", parameter: "Random" }      //根据需要自定义
    -
        name: "testping"
        get:
            -
                { index: "1", operation: "get", deviceResource: "ping"}
    -   name: "testmessage"
        get:
            - 
                { index: "1", operation: "get", deviceResource: "message"}
        set:
            -
                { index: "1", operation: "set", deviceResource: "message"}

coreCommands:   //表明edgexfoundry中的coredata模块可以向DeviceService下达哪些命令，下面的命令可根据需要自定义
  -
    name: "Random"    
    get:               //与读相关的命令
        path: "/api/v1/device/{deviceId}/Random"       //api接口
        responses:
          -
            code: "200"     //表明读取数据成功
            description: ""   //描述信可随意写
            expectedValues: ["randomnumber"]  //期待响应的值，里面的内容必须和上面deviceResources里面的name对应
          -
            code: "503"     //表明读取数据失败
            description: "service unavailable"
            expectedValues: []
  -
    name: "testping"
    get:
        path: "/api/v1/device/{deviceId}/testping"
        responses:
          -
            code: "200"
            description: "ping the device"
            expectedValues: ["ping"]
          -
            code: "503"
            description: "service unavailable"
            expectedValues: []
  -
    name: "testmessage"
    get:
      path: "/api/v1/device/{deviceId}/testmessage"
      responses:
        - code: "200"
          description: "get the message"
          expectedValues: ["message"]
        - code: "500"
          description: "internal server error"
          expectedValues: []
    put:   //与向设备发送指令相关的命令
      path: "/api/v1/device/{deviceId}/testmessage"
      parameterNames: ["message"]    //里面的内容必须和上面deviceResources里面的name对应
      responses:
        -
          code: "204"  //向设备发送命令成功
          description: "set the message."  //可随意写
          expectedValues: []   
        -
          code: "500"  //向设备发送命令失败
          description: "service unavailable"
          expectedValues: []
```
相关协议的核心代码在driver目录下的driver.go文件中，可具体看里面的HandleReadCommands()和HandleWriteCommands()函数。
## 设备服务运行
- 运行edgexfoundry，如通过docker将程序跑起来。访问https://github.com/edgexfoundry/developer-scripts/tree/master/releases/fuji/compose-files，下载里面的docker-compose-fuji.yml到本地，并更名为docker-compose.yml
（注意本DeviceService是基于fuji版本开发的，因此建议edgexfoundry也使用fuji版本）。然后打开终端，输入如下指令，运行edgexfoundry：
```cassandraql
sudo docker-compose -f docker-compose.yml up -d
```
输入如下指令，关闭edgexfoundry运行：
```cassandraql
sudo docker-compose -f docker-compose.yml down
```
- 若收集真实设备的数据，则运行真实设备；若用于测试，本服务mock目录下的device.go是一个模拟设备，可运行代码进行测试(注意，需要将device.go的IP地址稍作修改)。移动到mock文件下，打开终端，输入：
```cassandraql
go run device.go
```
- 将cmd/res目录下的configuration.toml文件里面的IP地址改好后，返回到device-udp-go主文件夹下，打开终端，输入如下命令构建设备服务，生成二进制可执行文件：
```cassandraql
make build
```
-  移动到cmd目录下，可以看见可执行文件device-udp-go，输入如下命令启动设备服务：
```cassandraql
./device-udp-go
```
终端界面将会显示日志。
## 读取数据方式
- 方式一：
打开edgexfoundry的UI界面：在浏览器输入：http://(IP地址):4000,登录后，注册网关，在DeviceService下面可以看到edgex-device-udp，点击就可以看到注册的设备，可以进行get或set操作。
- 方式二： 通过postman测试，在浏览器输入：http://(IP地址):48080/api/v1/event/device/(设备名称)/10，可以读取从设备发送到coredata里面的最近10条数据
## License
[Apache-2.0](LICENSE)
