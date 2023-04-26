# go-cmd

配置、解析命令行参数，匹配一个具体的Command。   

此模块中实现了config的打印、校验、交互模式命令。
```
import (
  _ "github.com/xpwu/go-cmd/cmd/interactive"
  _ "github.com/xpwu/go-cmd/cmd/printconf"
  _ "github.com/xpwu/go-cmd/cmd/validconf"
)
```
在import模块中用如上的方式引入，即可使用需要的默认命令。   

## Usage
### RegisterCmd
在init()中使用 RegisterCmd 注册命名，如果需要命令行参数，可以通过 args 解析。  
参数的解析方式为：(顺序不可改变)   
1、定义接收数据的参数及默认值；  
2、使用Arg相应类型的函数，设置参数；  
3、调用Arg.Parse()即可获取到参数。  
#### arg.ReadConfig(xxx) 
实现了读取配置文件相关的arg解析及json格式的配置文件读取，    
使用方式：
```
// 可以添加其他的参数解析命令
arg.ReadConfig(arg)
// 可以添加其他的参数解析命令
arg.Parse()

// 配置读取成功
```


### DefaultCommand  
使用 DefaultCmdName 即可注册默认命令 


### main 
在main()的最后执行 cmd.Run() 方法，此方法需要在main()中执行，
而不能在init()中执行，1、init()主要负责收集；2、如果用到config，
init()时还没有读取到config值  


## Note
1、因为命令名不能重名，所以代码中指定的命令名有可能被修改，以 -h 输出的命令名为准；   
2、DefaultCmdName 的名字不会修改


## Else
### exe
可以获取执行文件的 绝对路径 、 名字 


### client-cli
可以添加服务程序的交互式访问的客户端，   
调用clientcli.Listen(xxx)函数即可。 

* 使用：   
1、启动服务；2、使用 './xxx client' 命令启动此服务的客户端模式，即可通信。


