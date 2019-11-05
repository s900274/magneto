package main

import (
    "errors"
    "flag"
    "github.com/BurntSushi/toml"
    logger "github.com/shengkehua/xlog4go"
    "github.com/s900274/magneto/internal/define"
    "github.com/s900274/magneto/internal/server"
    "github.com/s900274/magneto/internal/kafkaconsumer"
    "github.com/s900274/magneto/pkg/helpers/kafkaproducer"
    "log"
    "os"
    "sync"
)


func initConf(flagSet *flag.FlagSet) error {

    flagSet.Parse(os.Args[1:])
    if len(flagSet.Lookup("config").Value.(flag.Getter).Get().(string)) > 0 {
        configFile := flagSet.Lookup("config").Value.String()
        if configFile != "" {
            _, err := toml.DecodeFile(configFile, &define.Cfg)
            if err != nil {
                log.Fatalf("ERROR: failed to load config file %s - %s\n", configFile, err.Error())
                return err
            }

        } else {
            log.Fatalln("ERROR: config file is nil")
            err := errors.New("ERROR: config file is nil")
            return err
        }
    } else {
        log.Fatalln("ERROR: no config param given")
        err := errors.New("ERROR: no config param given")
        return err
    }

    return nil
}

func initLogger() error {
    err := logger.SetupLogWithConf(define.Cfg.LogFile)
    return err
}

func initKafkaProducer() error {
    for tag, _ := range define.Cfg.KafkaConsumerCfg {
        _ = kafkaproducer.AddNewTopic(define.Cfg.KafkaProducerCfg.BrokerList, define.Cfg.KafkaProducerCfg.PartitionNum, tag)
        logger.Debug("Add topic %v for %v partitions", tag, define.Cfg.KafkaProducerCfg.PartitionNum)
    }

    for i := 0; i < define.Cfg.KafkaProducerCfg.ProduceNum; i++ {
        producer, err := kafkaproducer.NewSimAsyncProducer(define.Cfg.KafkaProducerCfg.BrokerList, define.Cfg.KafkaProducerCfg.PartitionNum)
        if err != nil {
            logger.Error("NewSimAsyncProducer failed;err:%s", err.Error())
            return err
        }
        producer.Start()
        kafkaproducer.ProducerList = append(kafkaproducer.ProducerList, producer)
    }
    return nil
}

func initKafkaConsumer() error {
    logger.Debug("init kafka consumer")
    logger.Debug("init kafka topic cfg : %v", define.Cfg.KafkaConsumerCfg)
    for tag, _ := range define.Cfg.KafkaConsumerCfg {

        logger.Debug("init kafka topic : %v", tag)
        consumerGroup, err := kafkaconsumer.NewKafkaConsumer(tag)
        if err != nil {
            logger.Error("New Kafka Consumer failed;err:%s", err.Error())
            return err
        }
        go consumerGroup.Start()
    }
    return nil
}

func RunHttpServer()  {
    var wg sync.WaitGroup

    wg.Add(1)
    go func() {
        defer wg.Done()
        HServer := server.NewHTTPServer()
        err := HServer.InitHttpServer()
        if nil != err {
            logger.Error("HTTPServerStart failed, err :%v", err)
            return
        }
    }()
    wg.Wait()
}