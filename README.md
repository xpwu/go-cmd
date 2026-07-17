# go-cmd

命令注册及参数的解析，很方便在一个可执行包中实现多个相关的功能，

## 1、 基本使用（执行包以./exe为例）   
### ./exe  [command]  < args >   
command: 通过RegisterCmdErr等函数注册的命令   
args: 该命令对应的命令行参数   
```go

import (
  "github.com/xpwu/go-cmd/cmd"
  "github.com/xpwu/go-cmd/arg"
)

func main() {
  cmd.RegisterCmd("print", "print config file", func(args *arg.Arg) {
    file := "config.default"
    args.String(&file, "f", "config file name")
    args.ParseAndRunHook()

    // --- print handle ---
    // 
  })

  cmd.RegisterCmd("listen", "listen port", func(args *arg.Arg) {
    port := ":80"
    ip := "127.0.0.1"
    args.String(&port, "port", "listen :port")
    args.String(&ip, "ip", "listen ip addr")
    args.ParseAndRunHook()
  
    // --- listen handle ---
    //
  })
  
  cmd.Run()
}
```
如上代码注册了两个命令（print 与 listen）   
print 支持一个参数 -f, 默认值为"config.default",    
listen 支持两个参数 -port 与 -ip   
cmd.Run() 在所有注册之后调用

```shell
// 查看帮助
./exe -h

// 查看 print 命令的帮助
./exe print -h

// 查看 listen 命令的帮助
./exe listen -h

// 以默认参数"config.default"执行 print 命令
./exe print 

// 以参数"config.json.default"执行 print 命令
./exe print -f config.json.default

// 以port参数"8080", ip默认值执行 listen 命令,
./exe listen -port 8080
```

## 2、RegisterKeepAliveCmd
通过RegisterCmd注册的Cmd在运行完后，整个进程退出，通过RegisterKeepAliveCmd
注册的Cmd在运行完后，不会退出进程，进程一直运行下去

## 3、命令名
如果注册了重名的命令名，后注册的会被修改为另外的非重名命令名，通过-h可以看到具体修改后的
名字。
#### 默认命令 cmd.DefaultCmdName   
如果是注册的cmd.DefaultCmdName (run)命令，后面注册的会覆盖前面注册的，只能有一个默认
命令名存在。如果注册了一个默认命令，以下两句运行结果等效。
```shell
// 明确指定运行默认命令
./exe run 
// 没有指定命令名，则运行默认命令，与上一句完全等效
./exe
```
如果默认命令带有参数，以-ok为例，下面几条命令相当
```shell
// 明确指定以参数ok运行run命令
./exe run -ok

// 带参数运行默认命令
./exe -ok

// 以ok的默认值运行默认命令
./exe
```
## 4、此lib实现的其它功能
#### 1）此lib中实现了config的打印、校验命令。
```
import (
  _ "github.com/xpwu/go-cmd/cmd/printconf"
  _ "github.com/xpwu/go-cmd/cmd/validconf"
)
```
在import模块中用如上的方式引入，则自动注册了这些命令。

#### 2）此lib 实现了配置文件的读取及解析
如下代码实现了xpwu/go-config格式的配置读取及解析
```go
import (
  "github.com/xpwu/go-cmd/arg" 
)

func main() {
  cmd.RegisterCmd("xxx", "start ...", func(args *arg.Arg) {

    arg.HookReadConfigTo(args)
    args.ParseAndRunHook()
		
		// --- handler ---
		// 
  })

  cmd.Run()
}
```

## 5、client cli
启动服务后，可以通过 client 命令启动一个与服务通信的交互式客户端，此交互式客户端的启动与退出不影响服务的运行。
代码中通过 clientcli.RegisterCmd(xxx) 方法可以注册交互式命令及对应的参数，再通过 clientcli.Start() 
即可启动注册的交互式命令。（先完成所有的注册，最后再启动）
```go

import (
  "github.com/xpwu/go-cmd/cmd"
  "github.com/xpwu/go-cmd/arg"
  "github.com/xpwu/go-cmd/clientcli"
)

func clientcliStart() {
  clientcli.RegisterCmd("sync", "force sync", func(args *arg.Arg) clientcli.AckToClient {
    ctx, logger := log.WithCtx(context.TODO())

    d, coll := "<db>", "<coll>"
    args.String(&d, "d", "db")
    args.ParseAndRunHook()

    // --- sync handle ---
    // 
   
    return "OK!"
  })

  clientcli.Start()
}

func main() {
  cmd.RegisterCmd("print", "print config file", func(args *arg.Arg) {
    file := "config.default"
    args.String(&file, "f", "config file name")
    args.ParseAndRunHook()

    // --- print handle ---
    //
  })

  cmd.RegisterKeepAliveCmd("listen", "listen port", func(args *arg.Arg) {
    port := ":80"
    ip := "127.0.0.1"
    args.String(&port, "port", "listen :port")
    args.String(&ip, "ip", "listen ip addr")
    args.ParseAndRunHook()
  
    // --- listen handle ---
    //

		clientcliStart()
		
		// --- else handle ---
  })
  
  cmd.Run()
}
```
通过 ./exe client 即可启动一个交互式客户端，输入 sync -d "user" 即可用参数 user 执行注册的 sync 命令