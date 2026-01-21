package main

import (
	"fmt"
	"encoding/json"
	"github.com/go-webauthn/webauthn/webauthn"
	"os"
)

const (
 user_file string = "./users.json"
)

type Session struct {
	Data *webauthn.SessionData
}

func NewSession() *Session {
	return new(Session)
}

func (s *Session) Save(sess *webauthn.SessionData) {
	s.Data = sess
}

func (s *Session) Load() (*webauthn.SessionData, error) {
	if s.Data == nil { return nil, fmt.Errorf("Did nor find session") }
	return s.Data, nil
}

type Users map[string]*User

func (u Users) Load() error {
	file, err := os.Open(user_file)
	if(err!=nil){ return err }
	return json.NewDecoder(file).Decode(&u)
}

func (u Users) Save() error {
	file, err := os.Open(user_file)
	if(err!=nil){ return err }
	return json.NewEncoder(file).Encode(&u)
}

func (u Users) AddUser(uo *User) (error) {
	user := u[uo.Email]
	if(user != nil){
		return fmt.Errorf("User %s exists!",uo.Email)
	}
	u[uo.Email] = uo
	return u.Save()
}

func (u Users) DeleteUser(uo *User) (error) {
	user := u[uo.Email]
	if(user != nil){
		return fmt.Errorf("User %s exists!",uo.Email)
	}
	delete(u,uo.Email)
	return nil
}

func (u Users) SetUser(uo *User) (error) {
	user := u[uo.Email]
	if(user != nil){
		return fmt.Errorf("User %s exists!",uo.Email)
	}
	u[uo.Email] = uo
	return nil
}

func (u Users) GetUser(email string) (*User, error) {
	user := u[email]

	if(user == nil){
		return nil,fmt.Errorf("User %s not found!",email)
	}

	return user, nil
}

type User struct {
	Email string `json:"email"`
	Name string `json:"name"`
	Session *Session `json:"-"`
}

func NewUser(email string,name string) *User {
	return &User{
		Email:email,
		Name: name,
	}	
}

func (u User) WebAuthnID() []byte {
	return []byte(u.Email)
}

func (u User) WebAuthnName() string {
	return u.Name
}

func (u User) WebAuthnDisplayName() string {
	return u.Name
}

func (u User)Encode() ([]byte, error) {
	return json.Marshal(&u)
}

func (u User)Decode(data []byte) (error) {
	return json.Unmarshal(data,&u)
}

