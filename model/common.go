package model

type User struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Age       int    `json:"age"`
	MobNumber int    `json:"mob_number"`
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
