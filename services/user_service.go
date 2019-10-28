package services

import (
	"log"
	"time"
	"net/http"
	"git.dev.tencent.com/petit_kayak/go-template-server/common"
	"git.dev.tencent.com/petit_kayak/go-template-server/database"
	"git.dev.tencent.com/petit_kayak/go-template-server/user"
)

// RegisterPara contains the parameters for registration
type RegisterPara struct {
	Email 		string
	Name			string
	Password	string
}

// LoginPara contains the parameters for login
type LoginPara struct {
	Email 		string
	Password	string
}

// StartSessionPara contains the parameters for starting new session
type StartSessionPara struct {
	UserID		string
}

// CheckSessionPara contains the parameters for checking existing session
type CheckSessionPara struct {
	SessionID	string
}

// FinishSessionPara contains the parameters for finishing existing session
type FinishSessionPara struct {
	SessionID	string
}

// LogoutPara contains the parameters for logout
type LogoutPara struct {
	SessionID	string
}

// UserInfo contains information of a user
type UserInfo struct {
	UserID	string
	Email 	string 
	Name		string
}

// SessionInfo contains information of a session
type SessionInfo struct {
	SessionID 	string
	UserID			string
	CreatedAt	time.Time
}

// GeneralEcho is a general echo
type GeneralEcho struct {
	Result	bool
	Error		string
}

// SuccessEcho is an echo with information of user and session when succeeded
type SuccessEcho struct {
	Result			bool
	UserInfo		UserInfo
	SessionInfo	SessionInfo
}

func register(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("Content-Type","application/json")//content type
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

func login(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("Content-Type","application/json")//content type
	para := LoginPara{}//prepare parameter
	echo := GeneralEcho{false,""}//prepare general echo
	err := common.LoadJSON(request.Body,&para)//parse parameter
	if err!=nil {
		log.Printf("Failed to parse register parameter: %s", err.Error())//log
		echo.Error = "Parse Parameter Error"//fill echo error
	} else {
		u, err := user.GetUserByEmail(database.TemplateDb,para.Email,para.Password)//try login
		if err!=nil || u.ID==0 {
			log.Printf("Failed to login with email %s", para.Email)//log
			echo.Error = "Login Error"//fill echo error
		} else {
			session, err := u.StartSession(database.TemplateDb)//try to start session
			if err!=nil || session.ID==0 {
				log.Printf("Failed to start session for user %s", u.UUID)//log
				echo.Error = "Start Session Error"//fill echo error
			} else {
				echo.Result = true//mark success
				infoEcho := SuccessEcho{true,
					UserInfo{u.UUID,u.Email,u.Name},
					SessionInfo{session.UUID,session.UserID,session.CreatedTs}}//fill info echo
				err = common.ExportJSON(&infoEcho,writer)//export info
			}
		}
	}
	if !echo.Result {
		err = common.ExportJSON(&echo,writer)//export echo
	}
	if err!= nil {
		log.Printf("Failed to return echo: %s", err.Error())
	}
}
