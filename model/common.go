package model

import "time"

type User struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Age       int    `json:"age"`
	MobNumber string `json:"mobile_number"`
	Email     string `json:"email"`
}

type StudentInfo struct {
	Students []User `json:"students_information"`
}

type GetStudentInfo struct {
	StudentID string `json:"id"`
}

type UpdateStudentInfo struct {
	Students User `json:"students_information"`
}

type Login struct {
	Password     string `json:"password"`
	Email        string `json:"email"`
	MobileNumber string `json:"mobile_number"`
}

type Signup struct {
	EmployeeID   string `json:"employee_id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	MobileNumber string `json:"mobile_number"`
	Email        string `json:"email"`
	DateOfBirth  string `json:"date_of_birth"`
	Gender       string `json:"gender"`
	Password     string `json:"password"`
	Token        string `json:"access_token"`
	CountryCode  string `json:"country_code"`
	PhotoUrl     string `json:"photo_url"`
	Coins        int    `json:"coins"`
}

type Internship struct {
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	MobileNumber string `json:"mobile_number"`
	Email        string `json:"email"`
	DateOfBirth  string `json:"date_of_birth"`
	Gender       string `json:"gender"`
	CountryCode  string `json:"country_code"`
	PhotoUrl     string `json:"photo_url"`
	Role         string `json:"role"`
	Duration     string `json:"duration"`
	Status string `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	CollegeName string `json:"college_name"`

}

type Visitor struct {
	FullName     string    `json:"full_name"`
	MobileNumber string    `json:"mobile_number"`
	Email        string    `json:"email"`
	Gender       string    `json:"gender"`
	CountryCode  string    `json:"country_code"`
	PhotoUrl     string    `json:"photo_url"`
	Branch       string    `json:"branch"`
	Purpose      string    `json:"purpose"`
	CreatedAt    time.Time `json:"created_at"`
}
