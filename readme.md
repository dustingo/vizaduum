### VIZADUUM
LINUX 进程数量、端口监控报警工具,可发送至报警中转服务器(由rabbitmq接收)，以及alertmanager；
正常存在外网的情况下，可以直接发送至mq，内网服务器可发送至存在外网ip的管理机中的alertmanager

> 依赖

Rabbitmq
Alertmanager
Consul

> 设计思路

1. 使用Consul作为工具的远程配置管理器,便于部署以及后续的配置文件修改
2. 定时读取Consul配置,实现配置文件的热更
3. 充分利用Consul的KV，使用key作为各服务器唯一的配置文件
4. 当前报警会写入到consul定义的program.file文件中，当业务进程处于异常或宕机时，不会重复报警
<!-- 5. 当program.file中字段为 "pause"时(文件内容，不是此字段值)，会处于维护状态，任何报警都不会发送，以实现维护时间断的静默 -->
5. 新增查询状态、设置维护、解除维护接口，删除了以program.file为依据的状态维护；接口信息参考router

> Consul KV说明
Key与服务器标识(ip/hostname等)对应，作为各服务器的唯一配置文件名称;
总共分为两部分:
1. 游戏业务配置key:
2. mq以及alertmanager配置key
如下所示:
key: game1
Value:  
```yaml
program:
  app: vizaduum #项目名称 
  cluster: vizaduum-国服 #项目逻辑服
  name:  game1 #功能服名称(如 game,gate等等)
  instance: 192.168.116.130 # hostname或IP
  altype: tt #报警方式
  team: tt-Mobile_Game_Alert #报警组
  call: 0 #是否电话通知
  channel: rabbitmq # 报警途径: rabbitmq|alertmanager
  file: record.last # 报警状态记录文件,会自动创建到工具所在目录
  interval: 3 # 检擦进程频率，单位 分钟
  port: ":7783" #工具启动的端口
service: #业务进程/端口检测说明
  - method: port  # 检测方式,port:端口检测; default: “ps”检测
    uniqname: mobilealert # 服务名
    ports:    #只有method=port 时使用
      - 8300
      - 8301
  - method: default
    uniqname: consul
    count: 2  #只有method=default 时使用，期望的consul进程数量
  - method: ura
    uniqname: zack
    count: 1
```
Key: sysconfig  
```yaml
rabbitmq:
  user: 
  pass: 
  url: 
  queue: 
alertmanager:
  url: "http://192.168.116.130:9093"
```

> 启动

```shell
./vizaduum --consul-address="192.168.116.130:8500" --consul-token="7d8dc827-02ea-28dc-4685-1a0ddde03434" --game-key="game1" --sys-key="sysconfig"
```
> 附录 Consul安装配置

```shell
sudo yum install -y yum-utils
sudo yum-config-manager --add-repo https://rpm.releases.hashicorp.com/RHEL/hashicorp.repo
sudo yum -y install consul
```
