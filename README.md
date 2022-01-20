# go-cmd

配置、解析命令行参数，匹配一个具体的Command。
此模块中实现了config的打印、校验、读取命令。

## Usage
### Package
在init()中使用 RegisterCmd 注册命名，如果需要命令行参数，可以通过 args 解析。  
参数的解析方式为：(顺序不可改变)   
1、定义接收数据的参数及默认值；  
2、使用Arg相应类型的函数，设置参数；  
3、调用Arg.Parse()即可获取到参数。  


### DefaultCommand  
使用 DefaultCmdName 即可注册默认命令 


### main 
在main()的最后执行 cmd.Run() 方法，此方法需要在main()中执行，
而不能在init()中执行，1、init()主要负责收集；2、如果用到config，
init()时还没有读取到config值  


## Note
1、因为命令名不能重名，所以代码中指定的命令名有可能被修改，以 -h 输出的命令名为准；   
2、DefaultCmdName 的名字不会修改
