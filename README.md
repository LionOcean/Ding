# Ding

> An application to transfer files instantly in wails(golang)

> version:  1.0.0

> Org:  LionOcean

> Technology stack:  Typescript + Vite + React + Golang + Wails v2

### Features
 * 采用wails v2包生成app、exe可执行文件，方便离线使用。
 * wails使用系统自带webview，而不是electron内置，软件打包后安装包2～3M，比内置chromium方案磁盘占用、内存占用率优秀很多
 * golang的go routine在读取文件和cpu计算会充分利用多核处理器，处理速度更快
 * 对go的GC阈值做了限制，文件传输过程中由于读取文件到buffer会申请很大的内存，虽然做了传输分片，但这个体积依旧不小。go是个带GC的语言，go的GC不会在变量未使用后及时回收（会过段时间自己执行GC），因为GC运行会让cpu占用率过高，但我们这个传输文件不涉及到服务波动，故可以让GC多次执行而避免内存常驻

### Usage Help

### Download
> 见[release](https://github.com/LionOcean/Ding/releases)
 
### Directory Tree
``` bash
    ├─build                 编译平台相关的配置文件、可执行程序
    ├─frontend              wails展示的前端静态资源，wails不强关联前端框架和构建工具
    │  ├─src
    │  ├─wailsjs           wails在加载静态资源时自动生成的方法bindings，见wails文档
    │  ├─index.html
    │  ├─package.json
    │  └─vite.config.ts     vite config配置
    ├─screenshot  
    ├─scripts               wails的常用命令组合的shell脚本，见wails文档
    ├─app.go                wails程序的app结构体，主要用于绑定go方法到js运行时，见wails文档
    ├─main.go               wails程序的入口，初始化，见wails文档
    └─wails.json            wails cli打包程序需要的配置，见wails文档
```

### Development Setup
#### Required dependencies
- Go 1.19+
- wails v2.4.1
- NPM (Node 16.17.1+)
> 详细见[wails文档](https://wails.io/zh-Hans/docs/gettingstarted/installation)

#### Dev liveload
项目逻辑分为前端和后端
- 前端可以使用任意框架，兼容性不用过于考虑，因为windows平台使用的webview2(和chromium一致)，前端使用go绑定方法也十分简单，都是挂载windows对象
- 后端的go方法可以使用任意go module，当需要绑定给js运行时，只需要在`main.go`里`bind`挂载struct即可，例如挂载`ding/transfer Peer结构体实例`
- 程序运行中使用前端项目的hot load热更新，打开一个终端运行`wails dev`即可，你也可以在浏览器中打开wails dev server，同样支持debug绑定的go方法

#### Build
```bash
# 使用upx压缩打包，upx是个命令行工具，需要提前安装
$ wails build --upx
```
虽然go支持交叉编译(任意本机系统架构都支持编译到其他系统架构)，但是由于wails使用了不同系统架构的webview程序绑定，所以wails目前只能跨平台运行，而不支持交叉编译

- 本机是mac平台
    - 可以编译darwin、windows和linux三个平台，对于arm64，必须芯片是arm64架构，可以统一编译amd64、arm64架构兼容的universal版本
- 本机是windows平台
    - 只能编译windows，可以单独编译amd64、arm64架构
- 本机是linux平台
    - 可以编译linux、windows和linux三个平台，对于arm64，必须芯片是arm64架构
- 所有平台均可使用upx来压缩，压缩比例非常强
- mac平台编译完的是app文件，windows默认编译完的是exe，所有平台均可使用nsis来打包exe
- wails使用的方案是系统自带的webview，目前macos/linux主流版本均内置，windows平台使用的webview2只在部分win10和正式win11内置，所以当你的windwos系统中不存在webview2时，程序启动后会引导你安装，大概118M左右
