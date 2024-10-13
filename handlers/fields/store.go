package fields

import (
	"crud/model"
	"crud/store/postgresql"
	"crud/store/redismanager"
	"crud/utils"
	"database/sql"
	"encoding/json"
	// "log"
	"net/http"
	"github.com/go-redis/redis/v8"
	"strings"
)


func New(cacheRepo redismanager.CacheManager, client *redis.Client, sqlRepo postgresql.SqlManager, 	sqlDB *sql.DB) Handler {
	return &fieldHandler{
		cacheRepo: cacheRepo,
		client:    client,
		sqlRepo:sqlRepo,
		sqlDB:sqlDB,
	}
}

type fieldHandler struct {
	cacheRepo redismanager.CacheManager
	client    *redis.Client
	sqlDB        *sql.DB
	sqlRepo postgresql.SqlManager 
}

func (h fieldHandler) CreateAll(w http.ResponseWriter, r *http.Request) {
	var req model.StudentInfo
	ctx := r.Context()
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.cacheRepo.SetUserData(ctx, h.client, req)
	if err != nil {
		utils.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.ReturnResponse(w, http.StatusOK, req)
}

func (h fieldHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req := model.User{}
	var err error

	req.ID, err = utils.GetURLParam(r, "student-info")

	user, err := h.cacheRepo.GetUserData(ctx, h.client, req.ID)
	if err != nil {
		utils.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.ReturnResponse(w, http.StatusOK, user)
}

func (h fieldHandler) Update(w http.ResponseWriter, r *http.Request) {
	req := model.User{}
	ctx := r.Context()

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = h.cacheRepo.UpdateUserData(ctx, h.client, req)
	if err != nil {
		utils.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.ReturnResponse(w, http.StatusOK, req)
}

func (h fieldHandler) Create(w http.ResponseWriter, r *http.Request) {
	req := model.User{}
	// ctx := r.Context()

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// err = h.cacheRepo.CreateUserData(ctx, h.client, req)
	// if err != nil {
	// 	utils.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	// if h.sqlDB == nil {
    //     http.Error(w, "Internal server error", http.StatusInternalServerError)
    //     log.Println("Database connection is nil")
    //     return
    // }

	// err = h.sqlRepo.SetUserData(ctx, h.sqlDB, req)
	// if err != nil {
	// 	utils.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	sqlStatement := `INSERT INTO users (name, age, mob_number, email, ) VALUES ($1, $2, $3, $4) RETURNING id`
	id := 0
	err = h.sqlDB.QueryRow(sqlStatement, req.Name, req.Age, req.MobNumber, req.Email).Scan(&id)
	if err != nil {
		utils.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
	}
	utils.ReturnResponse(w, http.StatusOK, req)
}



func (h fieldHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// req := model.User{}
	ctx := r.Context()

	// err := json.NewDecoder(r.Body).Decode(&req)
	// if err != nil {
	// 	utils.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	req, err := utils.GetURLParam(r, "student-info")


	err = h.cacheRepo.DeleteUserData(ctx, h.client, req)
	if err != nil {
		utils.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.ReturnResponse(w, http.StatusOK, "Deleted Succesful")
}

func (h fieldHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var users model.StudentInfo
	var err error
	users, err = h.cacheRepo.GetAll(ctx, h.client)
	if err != nil {
		utils.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.ReturnResponse(w, http.StatusOK, users)
}


func (h fieldHandler) Login(w http.ResponseWriter, r *http.Request) {
	req := model.Signup{}
	// ctx := r.Context()

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	query := "SELECT employee_id, first_name, last_name,mobile_number,  email, date_of_birth,gender,  password, access_token FROM employees WHERE email = $1"
	password:= ""
    err = h.sqlDB.QueryRow(query, req.Email).Scan(&req.EmployeeID,&req.FirstName,&req.LastName,&req.MobileNumber, &req.Email,&req.DateOfBirth,&req.Gender, &password, &req.Token)
	if err != nil {
		utils.ErrorResponse(w, "Invalid Email", http.StatusBadRequest)
	}else if password != req.Password{
		utils.ErrorResponse(w, "Password Invalid", http.StatusUnauthorized)
	}

	req.Token, err = utils.GenerateJWT(req.Email)
	if err != nil {
		utils.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sqlStatement := `UPDATE employees
        SET access_token = $1
        WHERE email = $2`
	_,err = h.sqlDB.Exec(sqlStatement, req.Token, req.Email)
	if err != nil {
		utils.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	req.Password = ""
	utils.ReturnResponse(w, http.StatusOK, req)
}

func (h fieldHandler) Signup(w http.ResponseWriter, r *http.Request) {
	req := model.Signup{}
	// ctx := r.Context()

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	query := "SELECT email FROM employees WHERE email = $1"
    err = h.sqlDB.QueryRow(query, req.Email).Scan(&req.Email)
		if err == nil {
			utils.ErrorResponse(w, "Email ID Already Exist", http.StatusBadRequest)
			return
		}
		query = "SELECT mobile_number FROM employees WHERE mobile_number = $1"
		err = h.sqlDB.QueryRow(query, req.MobileNumber).Scan(&req.MobileNumber)
			if err == nil {
				utils.ErrorResponse(w, "Mobile Number Already Exist", http.StatusUnauthorized)
				return
			}

		    req.Token, err = utils.GenerateJWT(req.Email)
			if err != nil {
				utils.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
				return
			}


			query = `INSERT INTO employees (first_name, last_name, mobile_number, email, date_of_birth, gender, password, access_token ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING employee_id`
		err = h.sqlDB.QueryRow(query, req.FirstName, req.LastName,req.MobileNumber, req.Email, req.DateOfBirth, req.Gender, req.Password, req.Token).Scan(&req.EmployeeID)
		if err != nil {
			utils.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}
		req.Password = ""

	utils.ReturnResponse(w, http.StatusOK, req)
}

func (h fieldHandler) PasswordChange(w http.ResponseWriter, r *http.Request) {
	req := model.Login{}
	// ctx := r.Context()

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	query := "SELECT email FROM employees WHERE email = $1"
    err = h.sqlDB.QueryRow(query, req.Email).Scan(&req.Email)
	if err != nil {
		utils.ErrorResponse(w,"Invalid Email", http.StatusBadRequest)
	}
	sqlStatement := `UPDATE employees
        SET password = $1
        WHERE email = $2`
	_,err = h.sqlDB.Exec(sqlStatement, req.Password, req.Email)
	if err != nil {
		utils.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
	}
	req.Password = ""
	utils.ReturnResponse(w, http.StatusOK, req)
}


func (h fieldHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")


	req := model.Signup{}
	// ctx := r.Context()

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	    // Simulate validation of received token
		_, err = utils.ValidateJWT(tokenString)
		if err != nil {
			utils.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
			return		}

	sqlStatement := `UPDATE employees
		             SET first_name = $1, last_name = $2, mobile_number = $3, email = $4, date_of_birth = $5, gender = $6
        			WHERE employee_id = $7`
	_,err = h.sqlDB.Exec(sqlStatement, req.FirstName, req.LastName, req.MobileNumber, req.Email, req.DateOfBirth, req.Gender, req.EmployeeID)
	if err != nil {
		utils.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
	}
	utils.ReturnResponse(w, http.StatusOK, req)
}