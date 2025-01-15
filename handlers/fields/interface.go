package fields

import "net/http"

type Handler interface {
	CreateAll(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	GetAll(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	Signup(w http.ResponseWriter, r *http.Request)
	PasswordChange(w http.ResponseWriter, r *http.Request)
	UpdateProfile(w http.ResponseWriter, r *http.Request)
	SendOTPEmailHandler(w http.ResponseWriter, r *http.Request)
	SendOTPMobHandler(w http.ResponseWriter, r *http.Request)
	VerifyOTPHandler(w http.ResponseWriter, r *http.Request)
	GetUser(w http.ResponseWriter, r *http.Request)
	GetImage(w http.ResponseWriter, r *http.Request)


}
