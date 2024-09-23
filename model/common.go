package model

type User struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Age       int    `json:"age"`
	MobNumber string    `json:"mob_number"`
	Email string `json:"email"`
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
	Password    string `json:"password"`
	Email string `json:"email"`
}

type Signup struct {
	EmployeeID string `json:"employee_id"`
	FirstName    string `json:"first_name"`
	LastName string `json:"last_name"`
	MobileNumber string `json:"mobile_number"`
	Email string `json:"email"`
	DateOfBirth string `json:"date_of_birth"`
	Gender string `json:"gender"`
	Password string `json:"password"`
}
