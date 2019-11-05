package main

import (
    logger "github.com/shengkehua/xlog4go"
    "github.com/s900274/magneto/pkg/downstream"
)

func main() {

    flagSet := Flagset()

    if err := initConf(flagSet); err != nil {
        logger.Error("%s", err.Error())
    }

    if err := downstream.InitRedisClient(); err != nil {
        logger.Error("init downstreams redis fail:%s", err.Error())
        return
    }

    if err := initLogger(); err != nil {
        logger.Error("%s", err.Error())
    }
    defer logger.Close()

    if err := initKafkaProducer(); err != nil {
        logger.Error("initKafkaProducer failed;err:%s", err.Error())
        return
    }

    if err := initKafkaConsumer(); err != nil {
       logger.Error("initKafkaConsumer failed;err:%s", err.Error())
       return
    }

    RunHttpServer()

    logger.Info("Server exit")
}
