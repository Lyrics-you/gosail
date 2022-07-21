# gosail
gosail is a free and open source batch and concurrent command execution system,designed to execute commands on multiple servers or k8s pods and get results with speed and efficiency.You can also copy(pull or push) files by it.

gosail 是一个免费开源的批处理并发命令执行系统，旨在在多个服务器或 k8s pod 上执行命令并快速高效地获取结果。还可以通过它复制（拉取或推送）文件。

<img src=".\gosail-exec.png" alt="gosail-exec" style="zoom: 50%;" />

<img src=".\gosail-k8s-cli.png" alt="gosail-k8s-cli" style="zoom:50%;" />

## 说明

支持：

- 并发执行
- 超时控制
- 输入参数
- 输出格式
- 颜色支持
- 输出适应窗口大小
- 支持交互终端

## 使用

### 编译

```go
go get ./...
go build .
```

### 提示

这里使用的是windows环境去处理，linux使用的是==gosail==二进制包，对应命令的斜杠替换为`/`即可。

## 参数

使用cobra框架

### help

通过 --help -? 可以查看参数含义

```shell
gosail is a free and open source batch and concurrent command execution system,
designed to execute commands on multiple servers or k8s pods and get results with speed and efficiency.
You can also copy(pull or push) files by it.

Usage:
  gosail [flags]
  gosail [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  login       Login host to do something
  version     Version subcommand show gosail version info

Flags:
      --ciphers string        ssh ciphers
      --config string         host execute config
  -?, --help                  help for this command
      --hostfile string       for gosail cli loginhost
  -K, --key string            id_rsa.pub key filepath
      --keyexchanges string   ssh keyexchanges
  -N, --numlimit int          max execute number (default 20)
  -p, --password string       for gosail cli password
  -T, --timelimit int         max timeout (default 30)
  -u, --username string       for gosail cli username
  -v, --version               gosail version

Use "gosail [command] --help" for more information about a command.
```

### version

通过 `version` 可以查看版本信息, `-d`可以查看版本信息描述

```shell
Name          : gosail⛵
Version       : x.x.x
Email         : Leyuan.Jia@Outlook.com
```

### color

通过上下键选择颜色

共有black/red/green/yellow/blue/magenta/white几种

<img src=".\gosail-color.png" alt="gosail-color" style="zoom: 50%;" />

### login

`gosail login --help`

```shell
eg. : gosail login -h <hostfile> [-u "<username>"] [-p "<password>"] [--prot "<port>"]

If the ssh port is 22, can omit port arg
eg. : gosail login -h <hostfile> [-u "<username>"] [-p "<password>"]

If the hostfile or hostline contain hosts in the format username@host, can omit u arg 
eg. : gosail login -h <hostfile> [-p "<password>"]

If specified the K arg or has default id_rsa.pub key, can omit p arg
eg. : gosail login -h <hostfile>

Usage:
  gosail login [flags]
  gosail login [command]

Available Commands:
  exec        Exec can execute commands concurrently and in batches on all hosts
  k8s         K8s master to do something
  pull        Pull can copy file from hosts concurrently, and create folders of host to distinguish
  push        Pull can copy file to hosts concurrently, and create folders that do not exist       

Flags:
  -h, --hostfile string   hostfile
      --hostline string   hostline
  -i, --ipfile string     ipfile
      --ipline string     ipline
  -p, --password string   host password
      --port int          ssh port (default 22)
  -u, --username string   host username

Global Flags:
      --ciphers string        ssh ciphers
      --config string         config
  -?, --help                  help for this command
  -K, --key string            id_rsa.pub key filepath
      --keyexchanges string   ssh keyexchanges
  -N, --numlimit int          max execute number (default 20)
  -T, --timelimit int         max timeout (default 30)

Use "gosail login [command] --help" for more information about a command.
```

#### 方式一：

hosts方式为分别指定主机IP，通过 ; 号或者 , 号作为命令和主机的分隔符。

IP指定时，需要加上共同用户名，通过-u指定，

支持user@hosts方式，IP地址前可以加上username（root@192.168.245.131)，可以省略-u参数。

#### --hostline 指定主机地址

#### --hostfile 指定主机地址文件

主机地址每行写入，

IP指定时，需要加上共同用户名，通过-u指定，

支持user@hosts方式，可以省略-u参数。

#### 方式二:

如果输入的是 IP ，那么允许IP地址段方式的输入，例如 192.168.245.131-192.168.245.133。 

#### --ipline 指定主机IP段

允许IP地址段方式的输入，

需要加上共同用户名，通过-u指定

#### --ipfile 指定主机IP段文件

主机地址每行写入，

IP指定时，需要加上共同用户名，通过-u指定，

### k8s

`gosail login k8s -?`

```shell
eg. : gosail login k8s -n "<namespace>" -a "<deployment.app>" [-c "<container>"]

Usage:
  gosail login k8s [flags]
  gosail login k8s [command]

Available Commands:
  exec        Exec can execute commands concurrently and in batches on all specified pods
  pull        Pull

Flags:
  -a, --app string         k8s deployment app
  -c, --container string   deployment container
  -l, --label string       deployment label
  -n, --namespace string   k8s namespace
      --shell string       container shell (default "sh")

Global Flags:
      --ciphers string        ssh ciphers
      --config string         config
  -?, --help                  help for this command
  -h, --hostfile string       hostfile
      --hostline string       hostline
  -i, --ipfile string         ipfile
      --ipline string         ipline
  -K, --key string            id_rsa.pub key filepath
      --keyexchanges string   ssh keyexchanges
  -N, --numlimit int          max execute number (default 20)
  -p, --password string       host password
      --port int              ssh port (default 22)
  -T, --timelimit int         max timeout (default 30)
  -u, --username string       host username

Use "gosail login k8s [command] --help" for more information about a command.
```

指定k8s的相关信息

后续exec\pull支持在k8s下使用

### exec

`gosail login exec -?`

```shell
eg. : gosail login exec [-b "<highlight>"] [-e] "<cmdline>" 
eg. : gosail login exec -e "<cmdline>" [-b "<highlight>"] mode [flags]
eg. : gosail login exec --cmdfile "<cmdfile>"

Usage:
  gosail login exec [flags]
  gosail login exec [command]

Available Commands:
  mode        Mode offers choices of exec output formats

Flags:
      --cmdfile string     exec cmdfile
  -e, --cmdline string     exec cmdline
  -b, --highlight string   bold highlight for cmdline and linuxmode

Global Flags:
      --ciphers string        ssh ciphers
      --config string         config
  -?, --help                  help for this command
      --host string           hostline
  -h, --hostfile string       hostfile
      --ip string             ipline
  -i, --ipfile string         ipfile
  -K, --key string            id_rsa.pub key filepath
      --keyexchanges string   ssh keyexchanges
  -N, --numlimit int          max execute number (default 20)
  -p, --password string       host password
      --port int              ssh port (default 22)
  -T, --timelimit int         max timeout (default 30)
  -u, --username string       host username

Use "gosail login exec [command] --help" for more information about a command.
```

#### --cmdline 指定命令行

可以通过;分隔多个命令

#### --cmdfile 指定命令行文件

也可以通过文本来存放主机组和命令组，通过换行符分隔。

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

#### -b:高亮

只有在cmdline和linuxmode中才其作用，可以实现高亮

#### k8s

不支持`cmdfile`,添加`-b`参数可以实现高亮

### mode

`gosail login exec mode -?` 

```shell
-j : use jsonmode to make the outpout with json format    
-l : use linuxmode to make the output without the hostname
-s : use selection to login hosts by their id

Usage:
  gosail login exec mode [flags]

Flags:
  -j, --jsonmode    json mode
  -l, --linuxmode   linux mode
  -s, --selection   select host to login

Global Flags:
      --ciphers string        ssh ciphers
      --cmdfile string        exec cmdfile
  -e, --cmdline string        exec cmdline
      --config string         config
  -?, --help                  help for this command
  -h, --hostfile string       hostfile
      --hostline string       hostline
  -i, --ipfile string         ipfile
      --ipline string         ipline
  -K, --key string            id_rsa.pub key filepath        
      --keyexchanges string   ssh keyexchanges
  -N, --numlimit int          max execute number (default 20)
  -p, --password string       host password
      --port int              ssh port (default 22)
  -T, --timelimit int         max timeout (default 30)       
  -u, --username string       host username
```

#### -l liunx模式

对于 linux ，支持 linuxMode 模式，也就是将命令组合通过 && 连接后，使用 session.Run() 运行。

不显示主机名，只有返回结果。

#### -j json格式输出

输出可以打成 json 格式，方便程序处理。

#### -s select 选择主机登录

可以通过输入id登录主机，并且显示主机是否可以登录的状态。

```shell
✋Server List:
Enter the 0~2 to select the host, other input will exit!
0   : 192.168.245.13  [x]
1   : 192.168.245.132 [√]
2   : 192.168.245.133 [√]
Input id :
```

### pull 

k8s pull 相同

pull时，可添加`-tar`参数将远端的文件压缩后，进行拉取。

从主机批量并发拉取文件到本地，本地文件支持相对路径以及其他主机目录（username@host:/path)，文件下每个主机的文件以各自主机名作为区分。

```shell
👇===============> nginx-ingress-controller-5bb8fb4bb6-2ndh7 (nginx-ingress-controller) <===============[0  ]
/usr/bin/scp -r root@192.168.245.133:nginx-ingress-controller-5bb8fb4bb6-2ndh7 ./demo/192.168.245.133/ Done!

👇===============> nginx-ingress-controller-5bb8fb4bb6-rgm4w (nginx-ingress-controller) <===============[1  ]
/usr/bin/scp -r root@192.168.245.133:nginx-ingress-controller-5bb8fb4bb6-rgm4w ./demo/192.168.245.133/ Done!

👇===============> nginx-ingress-controller-5bb8fb4bb6-twmzv (nginx-ingress-controller) <===============[2  ]
/usr/bin/scp -r root@192.168.245.133:nginx-ingress-controller-5bb8fb4bb6-twmzv ./demo/192.168.245.133/ Done!

👌Finshed!
```

### push

k8s没有该命令

从本地批量并发推送文件到主机，本地文件支持相对路径以及其他主机目录（username@host:/path)

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

### 通用参数

#### --ciphers/--keyexchanges

ssh登录参数

#### -K SSH密钥

支持使用 ssh 密钥认证，此时如果输入 password ，则为作为 key 的密码

默认密钥位置在 UserHomeDir的.ssh下的id_rsa

keyPath := path.Join(homePath, ".ssh", "id_rsa")

#### -N/-T 数量和超时控制

-N 控制批量最大协程数

-T 控制ssh结果超时时间，默认为秒

#### --config json配置加载

也可以为每个主机定义不同的配置参数，执行不同的命令，以 json 格式加载配置。

`.\gosail.exe --config ".\example\ssh.json"`

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
            "CmdFile": "./examples/cmdfile"
        }
    ],
    "Global": {
        "NumLimit": 20,
        "TimeLimit": 30,
        "Ciphers": "aes128-ctr,aes192-ctr,aes256-ctr,aes128-cbc,3des-cbc",
        "KeyExchanges": "diffie-hellman-group1-sha1,curve25519-sha256@libssh.org,ecdh-sha2-nistp256,ecdh-sha2-nistp384,ecdh-sha2-nistp521,diffie-hellman-group-exchange-sha256,diffie-hellman-group14-sha1"
    }
}
```

## 交互终端

使用grumble框架

```shell
                        _ _  
    __ _  ___  ___  __ _(_) | 
   / _  |/ _ \/ _ |/ _  | | | 
  |  g  |  o  \ s \  a  i   l 
   \__, |\___/|___/\__,_|_|__|
   |___/   

⛵ x.x.x


gosail is a free and open source batch and concurrent command execution system,
designed to execute commands on multiple servers or k8s pods and get results with speed and efficiency.
You can also copy(pull or push) files by it.

Commands:
=========
  clear          clear the screen
  exec           Exec can execute commands concurrently and in batches on all hosts and k8s pods
  exit           exit the shell
  help           use 'help [command]' for command help
  login, select  Login host to do something
  mode           Mode offers choices of exec output formats
  pull           Pull can copy file from hosts or pods, and create folders to distinguish
  push           Pull can copy file to hosts concurrently, and create folders that do not exist
  set            Set the gosail config
  show           Show the hosts

Sub Commands:
=============

login:
  k8s  K8s master to do something
```

通过`gosail`或者`gosail --hostfile "<hostfile>" [-u "<username>"] [-p "<password">]`进入交互界面

### login

进入后通过`login`以及其参数加载宿主机的信息

### set

设置 key\ciphers\keyexchanges\numlimit\timelimit等参数

### show

显示hosts

### exec

批量在主机和k8s pod中执行命令，k8s支持通过`-b`高亮

直接通过`exec`可进入循环执行命令行模式

<img src=".\gosail-interact.png" alt="gosail-interact" style="zoom: 50%;" />

`clear`可以清理多余屏幕

`Ctrl+C`、`exit`或者`quit`可以退出

### mode

设置linuxmode\jsonmode\selection等参数

k8s的linuxmode为false

### pull和push

同命令行中的参数，k8s不支持push
