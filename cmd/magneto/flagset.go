package main

import (
    //"time"
    "flag"
)

func Flagset() *flag.FlagSet {
    flagSet := flag.NewFlagSet("app", flag.ExitOnError)
    flagSet.String("config", "", "path to config file")

    //global flag
    flagSet.String("logfile", "", "log output file")
    flagSet.Int("http_server_port", 5000, "http port")

    return flagSet
}
