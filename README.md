# go-commandline

配置、解析命令行参数。实现了config打印、校验、读取的命令。

## usage
### Package
1、定义接收数据的参数；  
2、使用tinyFlag package中的相应类型的函数解析参数；  
3、使用RegisterCommand()注册一个与配置相应的命令。  

以上操作常常在 package的init()中配置  

### DefaultCommand  
1、没有其他配置项匹配成功时，就会执行默认配置  
2、使用RegisterDefaultCommand() function 或者 DefaultCommand struct
即可注册默认命令

### main 
在main()的最后执行 command.Run() 方法，此方法需要在main()中执行，
而不能在init()中执行，1、init()主要负责收集；2、如果用到config，
init()时还没有读取到config值  

