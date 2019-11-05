package define

var Cfg ServiceConfig


const (
    EVENT_CHAT_MESSAGE          = "chat message"
    EVENT_CONNECT               = "connection"
    EVENT_DISCONNECT            = "disconnection"
    EVENT_ERROR                 = "error"

    ROOM_NAME = "%v:%v"
    PLAT_ROOM_NAME = "ROOM1"
    CONSUMER_GROUP  = "CHATROOM_GROUP_%v"

    KAFKA_TOPIC_CHATROOM = "CHATROOM"
)

const (
    BROADCAST_TYPE_UNICAST = 1
    BROADCAST_TYPE_ALL = 2
)