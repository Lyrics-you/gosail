# gosail
A tool for batch execution of shh commands, programmed with go.

一个使用go编写的批量执行ssh命令的工具。

![image-20220613110917195](https://github.com/Lyrics-you/gosail/blob/main/image-20220613110917195.png)

## 说明

支持：

- 并发执行
- 超时控制
- 输入参数
- 输出格式
- 颜色支持
- 输出适应窗口大小

命令：**gosail位置** + **主机IP**（-hosts\\-hostsfile\ips\ipfiles) "对应IP填写方式" + **命令**（-cmdline\-cmdfile） "对应命令填写方式" + **用户**（-u,IP方式为user@hosts,**可省略**） + **密码**（-p,有密钥可以**省略**） + **密钥位置**（-k,有密码或者默认密钥位置,**可省略**） + **其他指令**（看需求）

## 使用

### 编译

```go
go get ./...
go build .
```

### 提示

这里使用的是windows环境去处理，linux使用的是==gosail==二进制包，对应命令的斜杠替换为`/`即可。

## 参数

### 通用

#### 帮助

通过 -h -help -? 可以查看参数含义

```shell
  -ciphers string
        ciphers
  -cmdfile string
        cmdfile path
  -cmdline string
        command line
  -config string
        config file Path
  -fpath string
        write file path
  -hostfile string
        hostfile path
  -hosts string
        host address list
  -ipfile string
        ipfile path
  -ips string
        ip address list
  -j    print output in json format
  -k string
        ssh private key
  -keyexchanges string
        keyexchanges
  -l    linux mode : multi command combine with && ,such as date&&cd /opt&&ls
  -nl int
        max execute number (default 20)
  -otxt
        write result into txt
  -p string
        password
  -port int
        ssh port (default 22)
  -s    select host to login
  -tl int
        max timeout (default 30)
  -u string
        username
  -v    show version
```

#### 版本

通过 -v 可以查看版本信息

```shell
ToolName : gosail
Version : x.x.x
Email : Leyuan.Jia@Outlook.com
```

### 主机IP

#### 方式一：

hosts方式为分别指定主机IP，通过 ; 号或者 , 号作为命令和主机的分隔符。

IP指定时，需要加上共同用户名，通过-u指定，

支持user@hosts方式，IP地址前可以加上username（root@192.168.245.131)，可以省略-u参数。

#### -hosts:指定主机地址

命令如下：

`.\gosail.exe -hosts "192.168.245.131;192.168.245.132"  -cmdline "ls" -u root -p qwerty`

`.\gosail.exe -hosts "root@192.168.245.131;root@192.168.245.132"  -cmdline "ls" -p qwerty`

#### -hostfile:指定主机地址文件

主机地址每行写入，

IP指定时，需要加上共同用户名，通过-u指定，

支持user@hosts方式，可以省略-u参数。

命令如下：

`.\gosail.exe -cmdline "date"  -hostfile ".\examples\host-list" -u root -p qwerty`

hostfile

```text
192.168.245.131
root@192.168.245.132
root@192.168.245.133
```

#### 方式二:

如果输入的是 IP （-ips 或 -ipfile），那么允许IP地址段方式的输入，例如 192.168.245.131-192.168.245.133。 

#### -ips:指定主机IP段

允许IP地址段方式的输入，

需要加上共同用户名，通过-u指定

命令如下：
`.\gosail.exe -ips "192.168.245.131-192.168.245.132"  -cmdline "ls" -u root -p qwerty`

#### -ipfile:指定主机IP段文件

`.\gosail.exe -ipfile ".\examples\ip-list"  -cmdline "ls" -u root -p qwerty`

ipfile

```text
192.168.245.131-192.168.245.133
101.132.145.243
```

### 命令行

#### -cmdline:指定命令行

可以通过;分隔多个命令

`.\gosail.exe -hostfile ".\examples\host-list" -cmdline "date;ls" -u root -p qwerty`

### -cmdfile:指定命令行文件

也可以通过文本来存放主机组和命令组，通过换行符分隔。

命令如下：
`.\gosail.exe -hostfile ".\examples\host-list" -cmdfile ".\examples\cmdfile" -u root -p qwerty`

cmdfile

```text
cd /etc/sysconfig/network-scripts/
ls
date
```

echodate

```go
#!/bin/bash
 
#echo time
for((i=0;i<3;i++))
do
    sleep 1
    echo $(date +"%Y-%m-%d %H:%M:%S")
done
```

### 其他参数

#### -k SSH密钥

支持使用 ssh 密钥认证，此时如果输入 password ，则为作为 key 的密码

命令如下：
`.\gosail.exe -hostfile ".\example\host-list" -cmdline "ls" -u root -k "C:\Users\Taragrade\.ssh\id_rsa"`

默认密钥位置在 UserHomeDir的.ssh下的id_rsa

keyPath := path.Join(homePath, ".ssh", "id_rsa")

`.\gosail.exe -hostfile ".\example\host-list" -cmdline "ls" -u root`

#### -l liunx模式

对于 linux ，支持 linuxMode 模式，也就是将命令组合通过 && 连接后，使用 session.Run() 运行。

不显示主机名，只有返回结果。

命令如下：
`.\gosail.exe -hostfile ".\example\host-list" -cmdline "cd /opt;ls" -u root -p qwerty -l`

#### -c json配置加载

也可以为每个主机定义不同的配置参数，以 json 格式加载配置。

`.\gosail.exe -c ".\example\ssh.json"`

json配置参数：参考

```go
type SSHHost struct {
	Host      string    `json:"Host"`
	Port      int       `json:"Port"`
	Username  string    `json:"Username"`
	Password  string    `json:"Password"`
	CmdFile   string    `json:"CmdFile"`
	CmdLine   string    `json:"CmdLine"`
	CmdList   []string  `json:"CmdList"`
	Key       string    `json:"Key"`
	LinuxMode bool      `json:"LinuxMode"`
	Result    SSHResult `json:"-"`
}
```

ssh.json

```json
{
    "Hosts": [{
            "Host": "192.168.245.131",
            "Port": 22,
            "Username": "root",
            "Password": "qwerty",
            "CmdLine": "ls"
        },
        {
            "Host": "192.168.80.132",
            "Port": 22,
            "Username": "root",
            "Password": "",
            "key": "",
            "linuxMode": true,
            "CmdFile": "cmdfile"
        }
    ],
    "Global": {
        "Ciphers": "aes128-ctr,aes192-ctr,aes256-ctr,aes128-cbc,3des-cbc",
        "KeyExchanges": "diffie-hellman-group1-sha1,curve25519-sha256@libssh.org,ecdh-sha2-nistp256,ecdh-sha2-nistp384,ecdh-sha2-nistp521,diffie-hellman-group-exchange-sha256,diffie-hellman-group14-sha1"
    }
}
```

#### -j json格式输出

输出可以打成 json 格式，方便程序处理。

命令如下：

`.\gosail.exe -c ".\examples\ssh.json" -j`

#### -otxt 输出txt文件

#### -path 输出文件位置

也可以把输出结果存到以主机名命名的文本中，比如用来做配置备份

命令如下：

`.\gosail.exe -c ".\examples\ssh.json" -path "./" -otxt`

#### -s select 选择主机登录

可以通过输入id登录主机，并且显示主机是否可以登录的状态。

命令如下：

`.\gosail.exe -hostfile ".\examples\host-list" -cmdline "cd /etc && ls" -s`

```shell
✋Server List:
Enter the 0~2 to select the host, other input will exit!
0   : 192.168.245.13  [x]
1   : 192.168.245.132 [√]
2   : 192.168.245.133 [√]
Input id :
```



# gocy

依赖于gosail的一个并发复制文件（pull\push)的工具。

试想两种需求，往多台主机上传递文件（push）,或者从多台主机上拉取文件（pull）。

![image-20220613122328741](https://github.com/Lyrics-you/gosail/blob/main/image-20220613122328741.png)

## 使用

### 编译

```go
go get ./...
cd gocy/ && go build .
```

## 参数

### 帮助

通过 -h -help -? 可以查看参数含义

```shell
  -ciphers string
        ciphers
  -config string
        config file Path
  -hostfile string
        hostfile path
  -hosts string
        host address list
  -ipfile string
        ipfile path
  -ips string
        ip address list
  -k string
        ssh private key
  -keyexchanges string
        keyexchanges
  -nl int
        max execute number (default 20)
  -p string
        password
  -path string
        pull or push's destination path
  -port int
        ssh port (default 22)
  -pull string
        pull's source path
  -push string
        push's source path
  -s    select host to login
  -tl int
        max timeout (default 30)
  -u string
        username
  -v    show version
```

### 版本

通过 -v 可以查看版本信息

```shell
ToolName : gocy
Version : x.x.x
Email : Leyuan.Jia@Outlook.com
```

### pull/push

其他参数可以参考gosail使用，不多描述。

pull/push，底层通过gosail执行scp命令，所以，最好运行主机与目标主机之间已经建立免密。

```shell
ssh-keygen
ssh-copy-id -i ~/.ssh/id_rsa.pub username@hostname
```

#### pull

从主机批量并发拉取文件到本地，本地文件支持相对路径以及其他主机目录（username@host:/path)，文件下每个主机的文件以各自主机名作为区分。

`./gocy -hostfile "./examples/host-list" -pull "/root/demo/" -path "../demo/" `

`./gocy -hostfile "./examples/host-list" -pull "/root/demo/" -path "root@192.168.245.131:/root/demo"`

```shell
.
├── 192.168.245.132
│   └── examples
│       ├── cmdfile
│       ├── echodate
│       ├── example-cmdfile
│       ├── host-list
│       ├── host-list-name
│       ├── ip-list
│       └── ssh.json
└── 192.168.245.133
    └── examples
        ├── cmdfile
        ├── echodate
        ├── example-cmdfile
        ├── host-list
        ├── host-list-name
        ├── ip-list
        └── ssh.json
```

#### push

从本地批量并发推送文件到主机，本地文件支持相对路径以及其他主机目录（username@host:/path)

`./gocy -hostfile "./examples/host-list" -push "../demo" -path "/root/demo/"`

`./gocy -hostfile "./examples/host-list" -push "root@192.168.245.131:/root/demo" -path "/root/demo/"`

### tar

pull时，可添加`-tar`参数将远端的文件压缩后，进行拉取。

`./gocy -hostfile "./examples/host-list" -pull "/root/demo/" -path "../demo/"· -tar`

### 登录主机

最后入`-s`可通过id登录主机，详情见gosail



# gobars

依赖于gosail的一个批量在k8s容器中执行命令工具

![image-20220628194009292](https://github.com/Lyrics-you/gosail/blob/main/image-20220628194009292.png)

![image-20220628193758224](https://github.com/Lyrics-you/gosail/blob/main/image-20220628193758224.png)

## 使用

### 编译

```go
go get ./...
cd gobars/ && go build .
```

## 参数

### 帮助

通过 -h -help -? 可以查看参数含义

```shell
  -app string
        k8s app name    
  -c string
        k8s container   
  -ciphers string       
        ciphers
  -cmdline string       
        command line    
  -config string        
        config file Path
  -fpath string
        write file path
  -hostfile string
        hostfile path
  -hosts string
        host address list
  -ipfile string
        ipfile path
  -ips string
        ip address list
  -j    print output in json format
  -k string
        ssh private key
  -keyexchanges string
        keyexchanges
  -l string
        k8s label
  -n string
        k8s namespace
  -nl int
        max execute number (default 20)
  -otxt
        write result into txt
  -p string
        password
  -path string
        pull's destination path
  -port int
        ssh port (default 22)
  -pull string
        pull's source path
  -s    select host to login
  -scp
        k8s cp function
  -tl int
        max timeout (default 30)
  -u string
        username
  -v    show version
```

### 版本

通过 -v 可以查看版本信息

```shell
ToolName : gobars
Version : x.x.x
Email : Leyuan.Jia@Outlook.com
```

### cmdline

指定k8s的master集群地址、pod的namespace、容器名称以及执行命令，即可批量并发在pod中执行命令

`./gobars -hostfile "../examples/host-list-k8s" -n ingress-nginx -c nginx-ingress-controller -cmdline "date" `

```shell
👇===============> nginx-ingress-controller-5bb8fb4bb6-2ndh7 <===============[0  ]
👉 ------------> date
Tue Jun 28 10:53:05 UTC 2022

👇===============> nginx-ingress-controller-5bb8fb4bb6-rgm4w <===============[1  ]
👉 ------------> date
Tue Jun 28 10:53:05 UTC 2022

👇===============> nginx-ingress-controller-5bb8fb4bb6-twmzv <===============[2  ]
👉 ------------> date
Tue Jun 28 10:53:05 UTC 2022

👌Finshed!
```

### pull

指定k8s的master集群地址、pod的namespace、容器名称，需要加入参数-copy

pull功能底层使用gosail执行scp命令，所以，最好运行主机与目标主机之间已经建立免密，详情见gocy。

从pod批量并发拉取文件到本地，文件下每个pod的文件以各自名称作为区分。

`./gobars -hostfile "../examples/host-list-k8s" -n ingress-nginx -c nginx-ingress-controller -copy -pull "/etc/nginx" -path "./demo/"`

```shell
👇===============> nginx-ingress-controller-5bb8fb4bb6-2ndh7 (nginx-ingress-controller) <===============[0  ]
/usr/bin/scp -r root@192.168.245.133:nginx-ingress-controller-5bb8fb4bb6-2ndh7 ./demo/192.168.245.133/ Done!

👇===============> nginx-ingress-controller-5bb8fb4bb6-rgm4w (nginx-ingress-controller) <===============[1  ]
/usr/bin/scp -r root@192.168.245.133:nginx-ingress-controller-5bb8fb4bb6-rgm4w ./demo/192.168.245.133/ Done!

👇===============> nginx-ingress-controller-5bb8fb4bb6-twmzv (nginx-ingress-controller) <===============[2  ]
/usr/bin/scp -r root@192.168.245.133:nginx-ingress-controller-5bb8fb4bb6-twmzv ./demo/192.168.245.133/ Done!

👌Finshed!
```

```shell
[root@centos-7-01 gobars]# cd demo/
[root@centos-7-01 demo]# ls
192.168.245.133
[root@centos-7-01 demo]# cd 192.168.245.133/
[root@centos-7-01 192.168.245.133]# ls
nginx-ingress-controller-5bb8fb4bb6-2ndh7  nginx-ingress-controller-5bb8fb4bb6-rgm4w  nginx-ingress-controller-5bb8fb4bb6-twmzv
```

### tar

pull时，可添加`-tar`参数将远端的文件压缩后，进行拉取。

`./gobars -hostfile "../examples/host-list-k8s" -n ingress-nginx -c nginx-ingress-controller -copy -pull "/etc/nginx" -path "./demo/" -tar`

### 登录主机

最后入 `-s` 可通过id登录主机，详情见gosail
