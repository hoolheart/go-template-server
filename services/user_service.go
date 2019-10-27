package services

import (
	"log"
	"time"
	"net/http"
	"git.dev.tencent.com/petit_kayak/go-template-server/common"
	"git.dev.tencent.com/petit_kayak/go-template-server/database"
	"git.dev.tencent.com/petit_kayak/go-template-server/user"
)

type RegisterPara struct {
	Email 		string
	Name			string
	Password	string
}

type LoginPara struct {
	email 		string
	password	string
}

type StartSessionPara struct {
	UserID		string
}

type CheckSessionPara struct {
	SessionID	string
}

type FinishSessionPara struct {
	SessionID	string
}

type LogoutPara struct {
	SessionID	string
}

type UserInfo struct {
	UserID	string
	Email 	string 
	Name		string
}

type SessionInfo struct {
	SessionID 	string
	UserID			string
	CreatedAt	time.Time
}

type GeneralEcho struct {
	Result	bool
	Error		string
}

type SuccessEcho struct {
	Result			bool
	UserInfo		UserInfo
	SessionInfo	SessionInfo
}

func register(writer http.ResponseWriter, request *http.Request) {
	para := RegisterPara{}//prepare parameter
	echo := GeneralEcho{}//prepare echo
	err := common.LoadJSON(request.Body,&para)//parse
	if err!= nil || len(para.Name)==0 || len(para.Email)==0 || len(para.Password)==0 {
		if err!=nil {
			log.Printf("Failed to parse register parameter: %s", err.Error())//log
		} else {
			log.Printf("Failed to fetch valid parameter: name %s, email %s, len(password) %d", para.Name, para.Email, len(para.Password))//log
		}
		echo = GeneralEcho{false,"Parse Parameter Error"}
	} else {
		_, err = user.New(database.TemplateDb,para.Name,para.Email,para.Password)//try to create user
		if err!=nil {
			log.Printf("Failed to create new user with name %s, email %s: %s",para.Name,para.Email,err.Error())//log
			echo = GeneralEcho{false,"Create User Error"}
		} else {
			log.Printf("Successd to create user with name %s and email %s",para.Name,para.Email)//log
			echo = GeneralEcho{true,""}
		}
	}

	err = common.ExportJSON(&echo,writer)//export echo
	if err!= nil {
		log.Printf("Failed to return echo: %s", err.Error())
	}
}
