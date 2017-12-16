[![Build Status](https://travis-ci.org/Binly42/agenda-go2.svg?branch=master)](https://travis-ci.org/Binly42/agenda-go2)
[![codecov](https://codecov.io/gh/LIANGTJ/agenda-go2/branch/master/graph/badge.svg)](https://codecov.io/gh/LIANGTJ/agenda-go2)

## 简介

 agenda-go2 是在 [agenda-go](https://github.com/Binly42/agenda-go) 的基础上继续开发作业的, 对应 pml老师 的 [这篇博客](http://blog.csdn.net/pmlpml/article/details/78727210) 及 `ex-service-agenda.html`(*服务程序开发实战 - Agenda*) 。

## 下载

```shell
go get -u github.com/Binly42/agenda-go2
```



## 用法

```shell
Usage:
  agenda [flags]
  agenda [command]

Available Commands and local flags:
  createM    -s startTime //create Meeting 
  			-e endTime 
  			-t title 
  			-p participator 
  			      
            
 
  help        Help about any command
  
  login      -u username  //login 
  			-p password
  
  logout     			 //logout
  
  register   -u username  //register for further use
  			-p password
  			[-e] email
  			[-t] phone
  			
  query      -u           //search users or meetings
  			-m			// note that -u & -m can't appear at the same time
  			-s startTime
  			-e endTime

Root Flags:
  -a, --author string         Author name for copyright attribution (default "YOUR NAME")
      --config string         config file (default is $HOME/.cobra.yaml) (default "./.cobra.yaml")
  -h, --help                  help for agenda
  -l, --license licensetext   Name of license for the project (can provide licensetext in config)
  -b, --projectbase string    base project directory eg. github.com/spf13/
      --viper                 Use Viper for configuration (default true)

Use "agenda [command] --help" for more information about a command.
```
## 实现原理

 大致上:

> + entity 包中实现基本数据结构 User, Meeting, UserList, MeetingList 等, 同时也实现了 agenda 系统中需要它们具备的功能, 基本是根据作业要求的 html 上的 "附件: Agenda 业务需求" 来划分的; 业务操作中只要是在语义上足够合理的, 都会实现成 一个 User 作为 actor 调用其对应的方法完成该事物 的模式, 比如: `user.CancelAccount()` 这样, 但是与此同时, 与 agenda 有关的具体逻辑, 则不在 entity 包中实现 ;

> + 与 agenda 有关的具体逻辑, 在 agenda 包中实现, 其中的业务操作(只要合理)都假设 当前登录用户(通过 `LoginedUser()` 得到) 作为执行者, 从而, 由执行者调用其对应方法 ;

> *  entity 包中已实现各个对象的 序列化/反序列化 和 输入/输出 操作, 但是还是要由 model 包中的具体实现传入 (绑定好文件的) encoder/decoder 才能完成事实上的文件读写(比如 将一个 UserList 保存到文件中) ;

> * 理想情况下, 面向用户端的 UI 部分应该只直接导入 agenda 包中暴露的接口 ;

> *  其他细节, 基本按照作业要求的 html 上的内容进行 ;

> + 具体的 CLI 接口和命令的解析等, 由 [LIANGTJ]( https://github.com/LIANGTJ) 完成, 其针对不同命令调用 agenda 包中的不同接口 ;



## 样例

启动服务器：

```shell
liangtj@ubuntu:~/Desktop/GoWorkSpace/src/github.com/LIANGTJ/agenda-go2/service$ go run main.go 
[info]2017/12/16 07:28:59 server.go:29: Listtening addr: :8080
Listtening addr: :8080
[error]2017/12/16 07:30:16 action.go:21: UNIQUE constraint failed: user_infos.name

(UNIQUE constraint failed: user_infos.name) 
[2017-12-16 07:30:16]  


```



register 

```shell
mock测试：
PS E:\GoWorkSpace\src\github.com\LIANGTJ\agenda-go2\cli> .\main.exe register -u root -p 123
register called
register successfully

服务端：
【200 response】
liangtj@ubuntu:~/Desktop/GoWorkSpace/src/github.com/LIANGTJ/agenda-go2/cli$ ./main register -u ltj  -p 123 -e ltj@163.com -t 12345
[info]2017/12/16 07:30:27 root.go:67: Can't read config: open ./.cobra.yaml: no such file or directory
register called
[Register] Response:  
register successfully ltj


【非200 response】:

liangtj@ubuntu:~/Desktop/GoWorkSpace/src/github.com/LIANGTJ/agenda-go2/cli$ ./main register
[info]2017/12/16 07:29:30 root.go:67: Can't read config: open ./.cobra.yaml: no such file or directory
register called
Error[registerd]： user regiestered invalid

liangtj@ubuntu:~/Desktop/GoWorkSpace/src/github.com/LIANGTJ/agenda-go2/cli$ ./main register -u root -p 123 -e ltj@163.com -t 12345
[info]2017/12/16 07:30:16 root.go:67: Can't read config: open ./.cobra.yaml: no such file or directory
register called
Error[registerd]： the user has been existed


```

login

```shell
mock测试：
PS E:\GoWorkSpace\src\github.com\LIANGTJ\agenda-go2\cli> .\main.exe login -u root -p 123
login called by root
login with info password: 123
Login Sucessfully root

服务端测试：
liangtj@ubuntu:~/Desktop/GoWorkSpace/src/github.com/LIANGTJ/agenda-go2/cli$ ./main login  -u ltj  -p 123
[info]2017/12/16 07:35:32 root.go:67: Can't read config: open ./.cobra.yaml: no such file or directory
login called by ltj
login with info password: 123
[Login] Response:  
Login Sucessfully ltj
```

因为没多大意义，所以以下都不再展示mock测试

query user

```shell
liangtj@ubuntu:~/Desktop/GoWorkSpace/src/github.com/LIANGTJ/agenda-go2/cli$ ./main query -u
[info]2017/12/16 07:36:05 root.go:67: Can't read config: open ./.cobra.yaml: no such file or directory
query called
[QueryAccountAll] Response:  
+--------+-------+-------+
|  NAME  | EMAIL | PHONE |
+--------+-------+-------+
| root   |       |       |
| matrix |       |       |
| ltj    |       | 12345 |
+--------+-------+-------+

```

create meeting

```shell
liangtj@ubuntu:~/Desktop/GoWorkSpace/src/github.com/LIANGTJ/agenda-go2/cli$ ./main createM -s "2011-01-01 10:00:34" -e "2011-01-02 08:00:34" -t MatrixShareMeeting -p lrd
create Meeting called
start: 2011-01-01 10:00:34 +0000 UTC end: 2011-01-02 08:00:34 +0000 UTC
sucessfully create meeting
```

 logout

```shell
liangtj@ubuntu:~/Desktop/GoWorkSpace/src/github.com/LIANGTJ/agenda-go2/cli$ ./main logout
[info]2017/12/16 07:38:39 root.go:67: Can't read config: open ./.cobra.yaml: no such file or directory
logout called
logout sucessfully
```








