package utils

import (
    "runtime"
    "fmt"
)

func stack() []byte {
    buf := make([]byte, 1024)
    n := runtime.Stack(buf, false)
    return buf[:n]
}

func GetStackInfo() string {

    stackInfo := stack()

    return fmt.Sprintf("%s", stackInfo)
}

type SpreadCmd struct {
    Uid       string
    Pwd       string
    SendFrom  string
    SendTo    []string
    Channel   string
    SendType  []int
    MsgSource string
}
