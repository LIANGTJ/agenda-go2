## 简介

 agenda是会议管理的命令行程序，由go语言实现。

## 安装

```shell
go get -u github.com/Binly42/agenda-go
```



## 用法

```shell
Usage:
  agenda [flags]
  agenda [command]

Available Commands and local flags:
  add        -s startTime //create Meeting 
  			-e endTime 
  			-t title 
  			-p participator 
  			
  delete     -u           // delete user or meeting
  			-m 			// note that -u & -m can't appear at the same time
  			-t title         
            
  edit       -i           //edit meeting, invite a participator or delete participator
  			-d 
			-p participator 
  help        Help about any command
  
  login      -u username  //login 
  			-p password
  
  logout     			 //logout
  
  register   -u username  //register for further use
  			-p password
  			[-e] email
  			[-t] phone
  			
  search     -u           //search users or meetings
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

register 

```shell
$ ./agenda-go.exe register -u lyb -p 12345 -e lyb@gmail.com -t 2080290
register called
register sucessfully!

```

login

```shell
$ ./agenda-go.exe login -u lyb -p12345
login called
login called by lyb
login with info password: 12345
login sucessfully!

```



search user

```shell
$ ./agenda-go.exe search -u
search called
+-------+---------------+---------+
| NAME  |     EMAIL     |  PHONE  |
+-------+---------------+---------+
| sky   |               |         |
| lrd   |               |         |
| lyb   | lyb@gmail.com | 2080290 |
| binly | binly@git.com | 1234567 |
+-------+---------------+---------+

```

create meeting

```shell
$ ./agenda-go.exe add -s "2011-01-01 10:00:34" -e "2011-01-02 08:00:34" -t MatrixShareMeeting -p lrd
create Meeting called
start: 2011-01-01 10:00:34 +0000 UTC end: 2011-01-02 08:00:34 +0000 UTC
sucessfully create meeting
```

search meeting

```shell
$ ./agenda-go.exe search -s "2011-01-01 10:00:34" -e "2011-01-02 08:00:34"
search called
+--------------------+---------+----------------------+----------------------+---------------+
|       TITLE        | SPONSOR |      STARTTIME       |       ENDTIME        | PARTICIPATORS |
+--------------------+---------+----------------------+----------------------+---------------+
| MatrixShareMeeting | lyb     | 2011-01-01T10:00:34Z | 2011-01-02T08:00:34Z |  lyb          |
+--------------------+---------+----------------------+----------------------+---------------+

```

add participator

```shell
$ ./agenda-go.exe edit -i -u sky -t MatrixShareMeeting
edit called
$ ./agenda-go.exe search -s "2011-01-01 10:00:34" -e "2011-01-02 08:00:34"
search called
+--------------------+---------+----------------------+----------------------+---------------+
|       TITLE        | SPONSOR |      STARTTIME       |       ENDTIME        | PARTICIPATORS |
+--------------------+---------+----------------------+----------------------+---------------+
| MatrixShareMeeting | lyb     | 2011-01-01T10:00:34Z | 2011-01-02T08:00:34Z |  sky lyb      |
+--------------------+---------+----------------------+----------------------+---------------+

```

 logout

```shell
$ ./agenda-go.exe logout
logout called
logout sucessfully!

```

delete meeting

```shell
$ ./agenda-go.exe delete -m -t MatrixShareMeeting
delete called
$ ./agenda-go.exe search -s "2011-01-01 10:00:34" -e "2011-01-02 08:00:34"
search called

```

delete current user

```shell
$ ./agenda-go.exe delete -u
delete called
CancelAccount successfully
$ ./agenda-go.exe login -u lyb -p12345
login called
login called by lyb
login with info password: 12345
Error[login]： a nil user/*user is to be used

```






