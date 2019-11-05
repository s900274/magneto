package server

import (
    "encoding/json"
    "net/http"
    "os"
    "path/filepath"
    "regexp"
    "strings"
    "fmt"
    "github.com/gin-gonic/gin"
    "github.com/googollee/go-socket.io"
    logger "github.com/shengkehua/xlog4go"
    "github.com/swaggo/gin-swagger"
    "github.com/swaggo/gin-swagger/swaggerFiles"
    "github.com/s900274/magneto/internal/define"
    _ "github.com/s900274/magneto/internal/server/docs"
    "github.com/s900274/magneto/internal/server/middleware"
    "github.com/s900274/magneto/internal/server/services"
    cv "github.com/s900274/magneto/pkg/chat-violation"
    "github.com/s900274/magneto/pkg/helpers/common"
    "github.com/s900274/magneto/pkg/helpers/utils"
    "github.com/s900274/magneto/internal/server/apimodel"
    "github.com/s900274/magneto/pkg/helpers/kafkaproducer"
    "github.com/pkg/errors"
)

type HServer struct {
    RunDirPath string
}

//type HandlerFunc func(w http.ResponseWriter, r *http.Request)

func NewHTTPServer() *HServer {
    s := &HServer{}

    return s
}

var SocketioServer, _ = socketio.NewServer(nil)

var SocketioSocket socketio.Socket

// @title Swagger test_swag API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host
// @BasePath /
func (s *HServer) InitHttpServer() error {

    defer func() {
        err := recover()
        if err != nil {
            logger.Error("magneto panic err: %s", err)
            stackInfo := utils.GetStackInfo()
            //utils.CallSlack(stackInfo, define.SLACK_PANIC_CHANNEL, define.SLACK_PANIC_SENDFROM)
            logger.Error("magneto panic stackinfo: %s", stackInfo)
        }
    }()

    serverAddr := fmt.Sprintf("%s:%d", define.Cfg.Host, define.Cfg.Http_server_port)

    runDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
    txtPath := fmt.Sprintf("%s/../config/list.txt", runDir)
    cv.MessageFilter.InitChatViolation(txtPath)

    router := s.Router()
    gin.SetMode(gin.DebugMode)
    router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
    return router.Run(serverAddr)
}

func (s *HServer) Router() *gin.Engine {
    r := gin.Default()
    s.RunDirPath, _ = filepath.Abs(filepath.Dir(os.Args[0]))
    templeteDir := fmt.Sprintf("%s/../web/templates/*", s.RunDirPath)
    r.LoadHTMLGlob(templeteDir)
    staticDir := fmt.Sprintf("%s/../web/static", s.RunDirPath)

    magneto := r.Group("/magneto")
    {
        magneto.Static("/static", staticDir)
        wsClient := magneto.Group("/wscli")
        {
            wsClient.GET("", services.WebSocketClient)
        }
        fmt.Println("web socket client url: http://127.0.0.1:<:port>/wscli")

        v1 := magneto.Group("/v1")
        {
            handshakeController(v1)
        }
    }

    return r
}

func getFrontUserInfoFromCtx(ctx *gin.Context) (*bool, error) {
    _, exists := ctx.Get("UserAuth")
    if !exists {
        return nil, errors.New("not found user")
    }

    r := true
    return &r, nil

}

func socketHandler(ctx *gin.Context) {

    SocketioServer.On(define.EVENT_CONNECT, func(so socketio.Socket) {
        // Assign the socket to a global variable
        SocketioSocket = so

        isAuth, _ := getFrontUserInfoFromCtx(ctx)
        if !*isAuth{
            errMsg := "token fail"
            so.Emit(define.EVENT_CHAT_MESSAGE, apimodel.WSRespFmt(
                nil,
                define.ERR_CHECKTOKEN_ERROR,
                errMsg,
                "",
                "",
                0,
            ))
            so.Disconnect()
            return
        }

        roomName := define.PLAT_ROOM_NAME
        logger.Debug("%v join room : %v", so.Id(), roomName)

        // Join room
        so.Join(roomName)

        // listen the chat message event
        so.On(define.EVENT_CHAT_MESSAGE, func(msg string) {
            // Log the message body
            logger.Debug("Request body : %v", msg)

            // Decode request json body
            msgObj := apimodel.ChatMessageRequest{}
            common.Json2Struct(msg, &msgObj)

            // Swearing Filter
            //oldMsg := msgObj.Msg
            msgObj.Msg = cv.MessageFilter.WordsFilter(msgObj.Msg)

            // If the message is empty then do nothing
            if msgObj.Msg == "" {
                return
            }

            respMsg := apimodel.ChatMessageResponse{}
            respMsg.Msg = msgObj.Msg
            respMsg.SetAccount(msgObj.UUID)

            msgString := apimodel.WSRespFmt(
                respMsg,
                define.ERR_OK,
                "",
                msgObj.UUID,
                define.EVENT_CHAT_MESSAGE,
                define.BROADCAST_TYPE_UNICAST,
            )

            // If the sending message is different between original message (be filtered) then do not broadcast
            //if oldMsg != msgObj.Msg {
            //    so.Emit(define.EVENT_CHAT_MESSAGE, msgString)
            //    return
            //}

            UnicastMessage(roomName, msgString)
        })

        // listen the disconnect event
        so.On(define.EVENT_DISCONNECT, func() {
            logger.Debug("%v on disconnect", so.Id())
        })
    })

    // Socket io error handler
    SocketioServer.On(define.EVENT_ERROR, func(so socketio.Socket, err error) {
        logger.Error("Failed : %v", err)
    })

    SocketioServer.ServeHTTP( ctx.Writer, ctx.Request)

    //[GIN-debug] [WARNING] Headers were already written. Wanted to override status code 200 with 400
    ctx.Abort()
}


func handshakeController(v1 *gin.RouterGroup) *gin.RouterGroup {

    hsGroup := v1.Group("/socket.io")
    {
        hsGroup.GET("/", middleware.CheckToken(), socketHandler)
        hsGroup.POST("/", middleware.CheckToken(), socketHandler)
        hsGroup.Handle("WS", "/", middleware.CheckToken(), socketHandler)
        hsGroup.Handle("WSS", "/", middleware.CheckToken(), socketHandler)
        hsGroup.GET("/unicast", unicastMessageHandler)
        hsGroup.GET("/multicast", multicastPlatMessageHandler)
        hsGroup.GET("/broadcast", broadcastMessageHandler)
        hsGroup.GET("/sessions", sessionCounter)
    }

    return v1
}

func sessionCounter(ctx *gin.Context) {

    ctx.JSON(http.StatusOK, gin.H{
        "code" : http.StatusOK,
        "sessions": SocketioServer.Count(),// session count
    })
    return
}

func messageProducer(roomName string, msgString string) {

    userInfo := &kafkaproducer.KfkJobData{
        Topic: define.KAFKA_TOPIC_CHATROOM,
        Key:    roomName,
        Value:  msgString,
    }
    kafkaproducer.Jobchan <- userInfo
}

func MessageConsumer(roomName string, msgString string) {

    var msgData = &apimodel.WSResponse{}
    err := json.Unmarshal([]byte(msgString), msgData)

    if err != nil {
        logger.Error("Decode message failed: %v", err)
    } else {
        SocketioServer.BroadcastTo(roomName, msgData.Event, msgString)
    }
}


func unicastMessageHandler(c *gin.Context) {
    roomName := fmt.Sprintf(define.ROOM_NAME, c.Request.FormValue("platId") ,c.Request.FormValue("memberId"))

    UnicastMessage(roomName, apimodel.WSRespFmt(
        "send to user MSG",
        define.ERR_OK,
        "",
        "",
        "",
        define.BROADCAST_TYPE_UNICAST,
    ))
}

//send all client
func UnicastMessage(roomName string, msg string) {

    logger.Debug("Room Name : %v", roomName)
    messageProducer(roomName, msg)
}

func multicastPlatMessageHandler(c *gin.Context) {

    var platSlice []string

    platSlice = append(platSlice, c.Request.FormValue("platId"))

    MulticastPlatMessage(platSlice,
        apimodel.WSRespFmt(
            "send to user MSG",
            define.ERR_OK,
            "",
            "",
            "",
            define.BROADCAST_TYPE_UNICAST,
        ))
}

func MulticastPlatMessage(platMap []string, msg string) {

    platConcat := strings.Join(platMap,"|")

    for _, v := range SocketioSocket.Rooms() {
        logger.Debug("Room Name : %v", v)
        if m, _ := regexp.MatchString("^("+platConcat+")$", v); m {
            messageProducer(v, msg)
        }
    }
}

func broadcastMessageHandler(c *gin.Context) {
    BroadcastMessage(apimodel.WSRespFmt(
        "send to user MSG",
        define.ERR_OK,
        "",
        "",
        "",
        define.BROADCAST_TYPE_ALL,
    ))
}

func BroadcastMessage(msg string) {
    for _, v := range SocketioSocket.Rooms() {
        logger.Debug("Room Name : %v", v)
        messageProducer(v, msg)
    }
}
