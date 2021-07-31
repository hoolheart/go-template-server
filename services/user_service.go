package services

import (
	"log"
	"time"
	"net/http"
	"gitee.com/alpha-bootes/go-template-server/common"
	"gitee.com/alpha-bootes/go-template-server/database"
	"gitee.com/alpha-bootes/go-template-server/user"
)

// RegisterPara contains the parameters for registration
type RegisterPara struct {
	Email 		string	`json:"email"`
	Name			string	`json:"name"`
	Password	string	`json:"password"`
}

// LoginPara contains the parameters for login
type LoginPara struct {
	Email 		string	`json:"email"`
	Password	string	`json:"password"`
}

// UserPara contains the parameters for existing user
type UserPara struct {
	UserID		string	`json:"user_id"`
}

// SessionPara contains the parameters for existing session
type SessionPara struct {
	SessionID	string	`json:"session_id"`
}

// ChangePasswordPara contains the paramters to change user password
type ChangePasswordPara struct {
	UserID			string	`json:"user_id"`
	OldPassword	string	`json:"old_password"`
	NewPassword	string	`json:"new_password"`
}

// UserInfo contains information of a user
type UserInfo struct {
	UserID	string	`json:"user_id"`
	Email 	string 	`json:"email"`
	Name		string	`json:"name"`
}

// SessionInfo contains information of a session
type SessionInfo struct {
	SessionID 	string	`json:"session_id"`
	UserID			string	`json:"user_id"`
	CreatedAt	time.Time	`json:"created_at"`
}

// GeneralEcho is a general echo
type GeneralEcho struct {
	Result	bool		`json:"result"`
	Error		string	`json:"error"`
}

// SuccessEcho is an echo with information of user and session when succeeded
type SuccessEcho struct {
	Result			bool				`json:"result"`
	UserInfo		UserInfo		`json:"user"`
	SessionInfo	SessionInfo	`json:"session"`
}

func getUsers(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {//check method
		writer.WriteHeader(404)//service not found
		return
	}
	writer.Header().Add("Content-Type","application/json")//content type
	echo := GeneralEcho{false,""}//preapre echo
	users, err := user.Users(database.TemplateDb)//get user list
	if err != nil {
		log.Printf("Failed to fetch all users: %s", err.Error())//log
		echo.Error = "Get users failed"//fill echo
	} else {
		log.Printf("Get %d users", len(users))//log
		usersEcho := make([]UserInfo,0)//prepare echo
		for _,u := range users {
			usersEcho = append(usersEcho, UserInfo{
				UserID 	: u.UUID,
				Email 	: u.Email,
				Name 		: u.Name,
			})//append every user info
		}
		err = common.ExportJSON(usersEcho,writer)//return users info
		echo.Result = true//mark result
	}
	if !echo.Result {
		err = common.ExportJSON(echo,writer)//return error
	}
	if err!=nil {
		log.Printf("Failed to return echo: %s", err.Error())
	}
}

func userHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("Content-Type","application/json")//content type
	switch request.Method {
	case http.MethodGet://get user
		getUser(writer,request)
	case http.MethodPost://register
		register(writer,request)
	case http.MethodPut://update
		update(writer,request)
	case http.MethodDelete://unregister
		unregister(writer,request)
	default:
		writer.WriteHeader(404)//service not found
	}
}

func getUser(writer http.ResponseWriter, request *http.Request) {
	para := UserPara{}//prepare parameter
	echo := GeneralEcho{false,""}//prepare general echo
	err := common.LoadJSON(request.Body,&para)//parse parameter
	if err!=nil {
		log.Printf("Failed to parse parameter: %s", err.Error())//log
		echo.Error = "Parse parameter error"//fill echo error
		request.ParseForm()
		para.UserID = request.FormValue("user_id")
	}
	if len(para.UserID)==0 {
		log.Println("Failed to parse user ID from the request")//log
		echo.Error = "No user ID specified"//fill echo error
	} else {
		u, err := user.Get(database.TemplateDb,para.UserID)//try login
		if err!=nil || u.ID==0 {
			log.Printf("Failed to login with ID %s", para.UserID)//log
			echo.Error = "Get user error"//fill echo error
		} else {
			log.Printf("Succeed to get info of user %s", para.UserID)
			echo.Result = true//mark success
			err = common.ExportJSON(&UserInfo{u.UUID,u.Email,u.Name},writer)//export info
		}
	}
	if !echo.Result {
		err = common.ExportJSON(&echo,writer)//export echo
	}
	if err!= nil {
		log.Printf("Failed to return echo: %s", err.Error())
	}
}

func register(writer http.ResponseWriter, request *http.Request) {
	para := RegisterPara{}//prepare parameter
	echo := GeneralEcho{}//prepare echo
	err := common.LoadJSON(request.Body,&para)//parse
	if err!= nil || len(para.Name)==0 || len(para.Email)==0 || len(para.Password)==0 {
		if err!=nil {
			log.Printf("Failed to parse parameter: %s", err.Error())//log
		} else {
			log.Printf("Failed to fetch valid parameter: name %s, email %s, len(password) %d", para.Name, para.Email, len(para.Password))//log
		}
		echo = GeneralEcho{false,"Parse parameter error"}
	} else {
		_, err = user.New(database.TemplateDb,para.Name,para.Email,para.Password)//try to create user
		if err!=nil {
			log.Printf("Failed to create new user with name %s, email %s: %s",para.Name,para.Email,err.Error())//log
			echo = GeneralEcho{false,"Create user error"}
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

func update(writer http.ResponseWriter, request *http.Request) {
	para := UserInfo{}//prepare parameter
	echo := GeneralEcho{false,""}//prepare echo
	err := common.LoadJSON(request.Body,&para)//parse parameter
	if err!=nil {
		log.Printf("Failed to parse parameter: %s", err.Error())//log
		echo.Error = "Parse parameter error"//fill echo error
	} else if len(para.UserID)==0 {
		log.Println("Failed to parse user ID from the request")//log
		echo.Error = "No user ID specified"//fill echo error
	} else if len(para.Email)==0 || len(para.Name)==0 {
		log.Printf("New parameters of user %s is invalid", para.UserID)//log
		echo.Error = "Email and name can't be empty"//fill error
	} else {
		u := user.User{
			UUID:para.UserID,
			Email:para.Email,
			Name:para.Name,
		}//prepare user
		err = u.Update(database.TemplateDb)//update user
		if err!=nil {
			log.Printf("Failed to update user %s: %s",para.UserID,err.Error())//log
			echo.Error = "Update user failed"//fill error
		} else {
			log.Printf("User %s updated", para.UserID)//log
			echo.Result = true//mark success
		}
	}
	err = common.ExportJSON(&echo,writer)//export echo
	if err!= nil {
		log.Printf("Failed to return echo: %s", err.Error())
	}
}

func unregister(writer http.ResponseWriter, request *http.Request) {
	para := UserPara{}//prepare parameter
	echo := GeneralEcho{false,""}//prepare echo
	err := common.LoadJSON(request.Body,&para)//parse parameter
	if err!=nil {
		log.Printf("Failed to parse parameter: %s", err.Error())
		echo.Error = "Parse parameter error"
	} else if len(para.UserID)==0 {
		log.Println("Failed to parse user ID from the request")//log
		echo.Error = "No user ID specified"//fill echo error
	} else {
		u := user.User{
			UUID : para.UserID,
		}//prepare user
		err = u.Delete(database.TemplateDb)//try to delete user
		if err!=nil {
			log.Printf("Failed to unregister user %s: %s", para.UserID, err.Error())
			echo.Error = "Unregister failed"
		} else {
			log.Printf("Delete user %s", para.UserID)
			echo.Result = true
		}
	}
	err = common.ExportJSON(&echo,writer)//export echo
	if err!= nil {
		log.Printf("Failed to return echo: %s", err.Error())
	}
}

func changePassword(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPut {//check method
		writer.WriteHeader(404)//service not found
		return
	}
	writer.Header().Add("Content-Type","application/json")//return json
	para := ChangePasswordPara{}//prepare parameter
	echo := GeneralEcho{false,""}//prepare echo
	err := common.LoadJSON(request.Body,&para)//parse parameter
	if err!=nil {
		log.Printf("Failed to parse parameter: %s", err.Error())//log
		echo.Error = "Parse parameter error"//fill echo error
	} else if len(para.UserID)==0 {
		log.Println("Failed to parse user ID from the request")//log
		echo.Error = "No user ID specified"//fill echo error
	} else if len(para.OldPassword)==0 || len(para.NewPassword)==0 {
		log.Printf("Input password of user %s is empty", para.UserID)//log
		echo.Error = "Password is empty"//fill error
	} else {
		u := user.User{
			UUID:para.UserID,
		}//prepare user
		err = u.ChangePassword(database.TemplateDb,para.OldPassword,para.NewPassword)//change password
		if err!=nil {
			log.Printf("Failed to change password of user %s: %s", para.UserID, err.Error())//log
			echo.Error = "Change password failed"//fill error
		} else {
			log.Printf("Password of user %s has been changed", para.UserID)//log
			echo.Result = true//mark success
		}
	}
	err = common.ExportJSON(&echo,writer)//export echo
	if err!= nil {
		log.Printf("Failed to return echo: %s", err.Error())
	}
}

func login(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {//check method
		writer.WriteHeader(404)//service not found
		return
	}
	writer.Header().Add("Content-Type","application/json")//content type
	para := LoginPara{}//prepare parameter
	echo := GeneralEcho{false,""}//prepare general echo
	err := common.LoadJSON(request.Body,&para)//parse parameter
	if err!=nil {
		log.Printf("Failed to parse parameter: %s", err.Error())//log
		echo.Error = "Parse parameter error"//fill echo error
	} else if len(para.Email)==0 || len(para.Password)==0 {
		log.Println("Failed to parse email or/and password from the request")//log
		echo.Error = "No Email or/and password specified"//fill echo error
	} else {
		u, err := user.GetUserByEmail(database.TemplateDb,para.Email,para.Password)//try login
		if err!=nil || u.ID==0 {
			log.Printf("Failed to login with email %s", para.Email)//log
			echo.Error = "Login Error"//fill echo error
		} else {
			session, err := u.StartSession(database.TemplateDb)//try to start session
			if err!=nil || session.ID==0 {
				log.Printf("Failed to start session for user %s", u.UUID)//log
				echo.Error = "Start session error"//fill echo error
			} else {
				log.Printf("Succeed to login with the email %s",para.Email)
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

func logout(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodDelete {//check method
		writer.WriteHeader(404)//service not found
		return
	}
	writer.Header().Add("Content-Type","application/json")//content type
	finishSession(writer,request)//call session finish
}

func sessionHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("Content-Type","application/json")//content type
	switch request.Method {
	case http.MethodGet://check session
		checkSession(writer,request)
	case http.MethodPost://start session
		startSession(writer,request)
	case http.MethodDelete://finish
		finishSession(writer,request)
	default:
		writer.WriteHeader(404)//service not found
	}
}

func startSession(writer http.ResponseWriter, request *http.Request) {
	para := UserPara{}//prepare parameter
	echo := GeneralEcho{false,""}//prepare general echo
	err := common.LoadJSON(request.Body,&para)//parse parameter
	if err!=nil {
		log.Printf("Failed to parse parameter: %s", err.Error())//log
		echo.Error = "Parse parameter error"//fill echo error
	} else if len(para.UserID)==0 {
		log.Println("Failed to parse user ID from the request")//log
		echo.Error = "No user ID specified"//fill echo error
	} else {
		u, err := user.Get(database.TemplateDb,para.UserID)//try login
		if err!=nil || u.ID==0 {
			log.Printf("Failed to login with ID %s", para.UserID)//log
			echo.Error = "Get user error"//fill echo error
		} else {
			session, err := u.StartSession(database.TemplateDb)//try to start session
			if err!=nil || session.ID==0 {
				log.Printf("Failed to start session for user %s", u.UUID)//log
				echo.Error = "Start session error"//fill echo error
			} else {
				log.Printf("Succeed to start session for user %s", para.UserID)
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

func checkSession(writer http.ResponseWriter, request *http.Request) {
	para := SessionPara{}//prepare parameter
	echo := GeneralEcho{false,""}//prepare echo
	err := common.LoadJSON(request.Body,&para)//parse parameter
	if err!=nil {
		//log.Printf("Failed to parse parameter: %s", err.Error())//log
		//echo.Error = "Parse parameter error"//fill echo error
		request.ParseForm()
		para.SessionID = request.FormValue("session_id")//get session id from request form
		err = nil
	}
	if len(para.SessionID)==0 {
		log.Println("Failed to parse session ID from the request")//log
		echo.Error = "No session ID specified"//fill echo error
	} else {
		session := user.Session{UUID:para.SessionID}//prepare session
		valid, err := session.Check(database.TemplateDb)//check session
		if err!=nil || !valid || session.ID==0 {
			log.Printf("Session %s is invalid",para.SessionID)//log
			echo.Error = "Session is invalid"//fill error
		} else {
			u, err := session.User(database.TemplateDb)//get user
			if err!=nil || u.ID==0 {
				log.Printf("Failed to get corresponding user of the session %s: %s",
					para.SessionID, common.If(err!=nil,err.Error(),"ID==0").(string))//log
				echo.Error = "Get user error"//fill error
			} else {
				log.Printf("Session %s is valid", para.SessionID)//log
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

func finishSession(writer http.ResponseWriter, request *http.Request) {
	para := SessionPara{}//prepare parameter
	echo := GeneralEcho{false,""}//prepare echo
	err := common.LoadJSON(request.Body,&para)//parse parameter
	if err!=nil {
		log.Printf("Failed to parse parameter: %s", err.Error())//log
		echo.Error = "Parse parameter error"//fill echo error
	} else if len(para.SessionID)==0 {
		log.Println("Failed to parse session ID from the request")//log
		echo.Error = "No session ID specified"//fill echo error
	} else {
		session := user.Session{UUID:para.SessionID}//prepare session
		valid, err := session.Check(database.TemplateDb)//check session
		if err!=nil || !valid || session.ID==0 {
			log.Printf("Session %s is invalid",para.SessionID)//log
			echo.Error = "Session is invalid"//fill error
		} else {
			err = session.Delete(database.TemplateDb)//delete session
			if err!=nil {
				log.Printf("Failed to finish session %s: %s", para.SessionID, err.Error())//log
				echo.Error = "Finish session error"//fill error
			} else {
				log.Printf("Session %s is finished", para.SessionID)//log
				echo.Result = true//mark success
			}
		}
	}
	err = common.ExportJSON(&echo,writer)//export echo
	if err!= nil {
		log.Printf("Failed to return echo: %s", err.Error())
	}
}
