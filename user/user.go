package user

import (
	"database/sql"
	"time"
	"git.dev.tencent.com/petit_kayak/go-template-server/common"
)

// User structure
type User struct {
	ID 				int64 // record ID in database
	UUID 			string // user UUID
	Email 		string // user e-mail
	Name 			string // user name
	Password 	string // user password
	CreatedTs time.Time // created timstamp of user
}

// Session structure
type Session struct {
	ID 				int64 // record ID of session
	UUID 			string // session UUID
	UserID 		string // relevant user UUID
	CreatedTs time.Time // created timstamp of session
}

// New creates a new user with given name, email and password
func New(db *sql.DB, name string, email string, password string) (user User, err error) {
	user = User{}//prepare user
	statement := "insert into t_user (uuid, name, email, password, created_at) values ($1, $2, $3, $4, $5) returning id, uuid, name, email, password, created_at"
	stmt, err := db.Prepare(statement)//prepare statement
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(
		common.CreateUUID(), 
		name, 
		email, 
		common.Encrypt(password), 
		time.Now()).Scan(
			&user.ID,
			&user.UUID,
			&user.Name,
			&user.Email,
			&user.Password,
			&user.CreatedTs)//apply statement
	return
}

// Fetch user info by using uuid
func (user *User)Fetch(db *sql.DB) (err error) {
	statement := "select id, email, name, password, created_at from t_user where uuid = $1"//prepare statement
	err = db.QueryRow(statement,user.UUID).Scan(&user.ID,&user.Email,&user.Name,&user.Password,&user.CreatedTs)//query from db
	return
}

// GetUserByEmail --- get a user by using email and password
func GetUserByEmail(db *sql.DB, email string, pw string) (user User, err error) {
	user = User{}//prepare user
	statement := "select id, uuid, email, name, password, created_at from t_user where email = $1 and password = $2"//prepare statement
	err = db.QueryRow(statement,email,common.Encrypt(pw)).Scan(&user.ID,&user.UUID,&user.Email,&user.Name,&user.Password,&user.CreatedTs)//query from db
	return
}

// Users --- get all users
func Users(db *sql.DB) (users []User, err error) {
	users = make([]User,0)//initial user list
	rows, err := db.Query("select id, uuid, name, email, password, created_at from t_user")//query from db
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {//check every record
		user := User{}
		if err = rows.Scan(&user.ID, &user.UUID, &user.Name, &user.Email, &user.Password, &user.CreatedTs); err != nil {//fetch user info from record
			return
		}
		users = append(users, user)//add user into list
	}
	return
}

// Update user name and email
func (user *User) Update(db *sql.DB) (err error) {
	statement := "update t_user set name = $2, email = $3 where id = $1"
	stmt, err := db.Prepare(statement)//prepare statement
	if err != nil {
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(user.ID, user.Name, user.Email)//execute statement
	return
}

// ChangePassword --- change user password
func (user *User) ChangePassword(db *sql.DB, curPw string, newPw string) (err error) {
	if len(newPw) == 0 {//check new password
		return common.NewError("go-template-server","user","ChangePassword","new password is empty")
	}
	err = user.Fetch(db)// fetch user info fist
	if err != nil {
		return
	}
	if common.Encrypt(curPw) != user.Password {//check current password
		return common.NewError("go-template-server","user","ChangePassword","current password is incorrect")
	}
	statement := "update t_user set password = $2 where id = $1"
	stmt, err := db.Prepare(statement)//prepare statement
	if err != nil {
		return
	}
	defer stmt.Close()
	user.Password = common.Encrypt(curPw)//update password
	_, err = stmt.Exec(user.ID, user.Password)//execute statement
	return
}

// Delete an existing user
func (user *User) Delete(db *sql.DB) (err error) {
	statement := "delete from t_user where id = $1"
	stmt, err := db.Prepare(statement)//prepare statement
	if err != nil {
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(user.ID)//execute statement
	return
}

// StartSession creates a new session for the user
func (user *User) StartSession(db *sql.DB) (session Session, err error) {
	session, _ = user.Session(db)//get existing session
	if session.ID != 0 {
		if err = session.Delete(db); err!=nil {//delete last session
			return
		}
		session = Session{}// prepare a blank session
	}
	statement := "insert into t_session (uuid, user_id, created_at) values ($1, $2, $3) returning id, uuid, user_id, created_at"
	stmt, err := db.Prepare(statement)//prepare insert statement
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(
		common.CreateUUID(), 
		user.UUID,
		time.Now()).Scan(
			&session.ID, 
			&session.UUID, 
			&session.UserID,
			&session.CreatedTs)//apply statement
	return
}

// Check whether the session is valid
func (session *Session) Check(db *sql.DB) (valid bool, err error) {
	err = db.QueryRow(
		"select id, uuid, user_id, created_at from t_session where uuid = $1",
		session.UUID).Scan(
			&session.ID,
			&session.UUID,
			&session.UserID,
			&session.CreatedTs)// try to query sesstion by its UUID
	if err != nil {
		valid = false
		return
	}
	valid = (session.ID!=0)
	return
}

// User queries the user correspondint to the session
func (session *Session) User(db *sql.DB) (user User, err error) {
	user = User{}//prepare user
	err = db.QueryRow("select id, uuid, name, email, created_at from t_user where uuid = $1", session.UserID).
		Scan(&user.ID, &user.UUID, &user.Name, &user.Email, &user.CreatedTs)
	return
}

// Session queries the last session of the user
func (user *User) Session(db *sql.DB) (session Session, err error) {
	session = Session{}//prepare session
	err = db.QueryRow(
		"select id, uuid, user_id, created_at from t_session where user_id = $1",
		user.UUID).Scan(
			&session.ID,
			&session.UUID,
			&session.UserID,
			&session.CreatedTs)//try to get last session
	return
}

// Delete time-out or finished session
func (session *Session) Delete(db *sql.DB) (err error) {
	statement := "delete from t_session where uuid = $1"
	stmt, err := db.Prepare(statement)//prepare statement
	if err != nil {
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(session.UUID)
	return
}
