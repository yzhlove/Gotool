# 允许外部访问
bind 0.0.0.0  
# 通信端口
port {{ .Port }}
# 后台运行
daemonize yes
# 关闭保护模式
protected-mode no
# 进程 PID 文件存放路径
pidfile {{ .PidPath }}
# 日志文件存放路径
logfile {{ .LogPath }}
# 是否以集群的形式启动
cluster-enabled yes
# 集群配置文件名
cluster-config-file {{ .ClusterCfgName }}
# 各个节点相互发送消息的频率，单位为毫秒。某节点发现与其他节点最后通信时间超过cluster-node-timeout/2时会发送ping命令，同时带上slots槽数组（2KB）和整个集群1/10的状态数据（10个节点的状态数据约1KB），该参数也会影响故障转移时间
cluster-node-timeout 15000
# RDB 文件名
dbfilename "dump.rdb"
# 指定数据文件(AOF,RDB)存放位置，必须要指定不同的目录位置，不然会丢失数据
dir {{ .DataDir }}
cluster-announce-ip redis-cluster
cluster-announce-port {{ .Port }}
cluster-announce-bus-port {{ .BusPort }}
