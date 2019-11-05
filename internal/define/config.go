package define

type ServiceConfig struct {
    LogFile          string `flag:"logfile" cfg:"logfile" toml:"logfile"`
    Http_server_port int    `flag:"http-server-port" cfg:"http_server_port" toml:"http_server_port"`
    Host             string
    RedisCfg         RedisConfig
    KafkaProducerCfg KafkaProducer
    KafkaConsumerCfg map[string]KafkaConsumer
}

type RedisConfig struct {
    Redis_svr           []string // redis server
    Redis_conn_timeout  int      // redis连接超时 毫秒
    Redis_read_timeout  int      // redis读超时 毫秒
    Redis_write_timeout int      // redis写超时 毫秒
    Redis_max_idle      int      // 最大空闲连接
    Redis_max_active    int      // 最大活动连接
    Redis_expire_second int      // redis数据过期时间 秒 线上配置10分钟 600
}

// Kafka Producer Config
type KafkaProducer struct {
    BrokerList        string
    BatchNum          int
    PartitionNum      int
    ProduceNum        int
    DialTimeout       int
    WriteTimeout      int
    ReadTimeout       int
    ReturnError       bool
    ReturnSuccess     bool
    FlushFrequency    int
    ChannelBufferSize int
}

type KafkaConsumer struct {
    Topic                string
    Group                string
    ProcessTimeout       int
    CommitInterval       int
    RetryTimes           int
    MetaMaxRetry         int
    ChannelSize          int
    HttpTimeout          int
    ZookeeperTimeout     int
    MetaRefreshFrequency int
    ZookeeperChroot      string
    RetryInterval        int
    ZookeeperAddresses   []string
}
