package apimodel

import (
    "fmt"
)

type TokenVerification struct {
    UserToken string `json:"userToken"`
    UserAccount string `json:"userAccount"`
    Origin string `json:"origin"`
    From string `json:"from"` //來源: ex: front(前端) / owner(業主) / master(主控)
}

type Account struct {
    UserAccount string `json:"userAccount"`
}

func (s *Account)SetAccount(account string)  {

    s.UserAccount = fmt.Sprintf("%s***", account[0:len(account)-3])
}

type ChatMessageRequest struct {
    Msg string `json:"msg"`
    UUID string `json:"uuid"`
    Type string `json:"type"`
}

type BetOrderRequest struct {
    Data interface{} `json:"data"`
    UUID string `json:"uuid"`
    Type string `json:"type"`
}


type BetOrderResponse struct {
    Account
    Msg interface{} `json:"msg"`
}

type ChatMessageResponse struct {
    Account
    Msg string `json:"msg"`
}
