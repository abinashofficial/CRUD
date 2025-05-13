package fields

import (
	"crud/model"
	"crud/store/postgresql"
	"crud/store/redismanager"
	"crud/utils"
	"database/sql"
	"github.com/go-redis/redis/v8"
	"strings"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
	"github.com/gorilla/websocket"
	"log"
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

// var otpStore = struct {
// 	mu   sync.Mutex
// 	data map[string]string // Stores OTPs with email as the key
// }{data: make(map[string]string)}

// OTPData stores the OTP and its expiration time
type OTPData struct {
	OTP        string
	ExpiryTime time.Time
}

var otpStore = struct {
	mu   sync.Mutex
	data map[string]OTPData // Stores OTPs and their expiration times, keyed by email
}{data: make(map[string]OTPData)}


// Upgrade HTTP connection to WebSocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins (adjust for production)
	},
}

// Store user connections
var clients = make(map[string]*websocket.Conn)
var mu sync.Mutex


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
		return
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




// Function to send a notification to a specific user
func sendNotification(userID string, message string) {
	mu.Lock()
	defer mu.Unlock()
	conn, exists := clients[userID]
	if exists {
		err := conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			fmt.Printf("Error sending message to %s: %v\n", userID, err)
		} else {
			fmt.Printf("Message sent to %s: %s\n", userID, message)
		}
	} else {
		fmt.Printf("User %s is not connected\n", userID)
	}
}




func (h fieldHandler) Login(w http.ResponseWriter, r *http.Request) {
	req := model.Signup{}
	// ctx := r.Context()
	password:= ""
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if req.Email !=""{
		query := "SELECT employee_id, first_name, last_name,mobile_number,  email, date_of_birth,gender,  password, country_code, photo_url, coins FROM employees WHERE email = $1"
		err = h.sqlDB.QueryRow(query, req.Email).Scan(&req.EmployeeID,&req.FirstName,&req.LastName,&req.MobileNumber, &req.Email,&req.DateOfBirth,&req.Gender, &password, &req.CountryCode, &req.PhotoUrl, &req.Coins)
		if err != nil {
			utils.ErrorResponse(w, "Invalid Email - "+ err.Error(), http.StatusBadRequest)
			return
		}
	}else if req.MobileNumber !=""{
		query := "SELECT employee_id, first_name, last_name,mobile_number,  email, date_of_birth,gender,  password, country_code, photo_url, coins FROM employees WHERE mobile_number = $1"
		err = h.sqlDB.QueryRow(query, req.MobileNumber).Scan(&req.EmployeeID,&req.FirstName,&req.LastName,&req.MobileNumber, &req.Email,&req.DateOfBirth,&req.Gender, &password, &req.CountryCode, &req.PhotoUrl, &req.Coins)
		if err != nil {
			utils.ErrorResponse(w, "Invalid Mobile Number - "+ err.Error(), http.StatusBadRequest)
			return
		}
	}
if password != req.Password{
		utils.ErrorResponse(w, "Password Invalid", http.StatusUnauthorized)
		return
	}

	req.Token, err = utils.GenerateJWT(req.Email)
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


			query = `INSERT INTO employees (first_name, last_name, mobile_number, email, date_of_birth, gender, password, country_code, photo_url ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING employee_id`
		err = h.sqlDB.QueryRow(query, req.FirstName, req.LastName,req.MobileNumber, req.Email, req.DateOfBirth, req.Gender, req.Password, req.CountryCode, req.PhotoUrl).Scan(&req.EmployeeID)
		if err != nil {
			utils.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}
		req.Password = ""

	utils.ReturnResponse(w, http.StatusOK, req)
}

func (h fieldHandler) PasswordChange(w http.ResponseWriter, r *http.Request) {
	req := model.Signup{}
	// ctx := r.Context()

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if req.Email ==""{
		query := "SELECT employee_id, first_name, last_name,mobile_number,  email, date_of_birth,gender, country_code, photo_url FROM employees WHERE mobile_number = $1"
		// password:= ""
		err = h.sqlDB.QueryRow(query, req.MobileNumber).Scan(&req.EmployeeID,&req.FirstName,&req.LastName,&req.MobileNumber, &req.Email,&req.DateOfBirth,&req.Gender, &req.CountryCode, &req.PhotoUrl)
		if err != nil {
			utils.ErrorResponse(w, "Invalid Mobile Number", http.StatusBadRequest)
			return
		}

		sqlStatement := `UPDATE employees
        SET password = $1
        WHERE mobile_number = $2`
	_,err = h.sqlDB.Exec(sqlStatement, req.Password, req.MobileNumber)
	if err != nil {
		utils.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return		
	}
	}else{
		query := "SELECT employee_id, first_name, last_name,mobile_number,  email, date_of_birth,gender, country_code, photo_url FROM employees WHERE email = $1"
		// password:= ""
		err = h.sqlDB.QueryRow(query, req.Email).Scan(&req.EmployeeID,&req.FirstName,&req.LastName,&req.MobileNumber, &req.Email,&req.DateOfBirth,&req.Gender, &req.CountryCode, &req.PhotoUrl)
		if err != nil {
			utils.ErrorResponse(w, "Invalid Email", http.StatusBadRequest)
			return
		}

		sqlStatement := `UPDATE employees
        SET password = $1
        WHERE email = $2`
	_,err = h.sqlDB.Exec(sqlStatement, req.Password, req.Email)
	if err != nil {
		utils.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return		
	}
	}
	query := "SELECT email FROM employees WHERE email = $1"
    err = h.sqlDB.QueryRow(query, req.Email).Scan(&req.Email)
	if err != nil {
		utils.ErrorResponse(w,"Invalid Email", http.StatusBadRequest)
		return
	}
	sqlStatement := `UPDATE employees
        SET password = $1
        WHERE email = $2`
	_,err = h.sqlDB.Exec(sqlStatement, req.Password, req.Email)
	if err != nil {
		utils.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return		
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
		             SET first_name = $1, last_name = $2, mobile_number = $3, email = $4, date_of_birth = $5, gender = $6, country_code =$7, photo_url =$8, coins = $9
        			WHERE employee_id = $10`
	_,err = h.sqlDB.Exec(sqlStatement, req.FirstName, req.LastName, req.MobileNumber, req.Email, req.DateOfBirth, req.Gender,req.CountryCode, req.PhotoUrl,req.Coins, req.EmployeeID)
	if err != nil {
		utils.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	notifyUser(req.Email , req.Coins )
	utils.ReturnResponse(w, http.StatusOK, req)
}





// SendOTPHandler handles OTP generation and sending
func (h fieldHandler)SendOTPEmailHandler(w http.ResponseWriter, r *http.Request) {

	type Response struct {
		Message string `json:"message"`
	}

	var req model.Signup
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Generate OTP
	otp, err := utils.GenerateOTP(6)
	if err != nil {
		http.Error(w, "Failed to generate OTP", http.StatusInternalServerError)
		return
	}

	// Store OTP with expiration time
	otpStore.mu.Lock()
	otpStore.data[req.Email] = OTPData{
		OTP:        otp,
		ExpiryTime: time.Now().Add(5 * time.Minute), // Set expiry to 1 minute
	}
	otpStore.mu.Unlock()
	// Send OTP via email
	subject := "Your OTP Code"
	body := fmt.Sprintf("Your OTP code is: %s", otp)
	if err := utils.SendEmail(req.Email, subject, body); err != nil {
		http.Error(w, "Failed to send OTP", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(Response{Message: "OTP sent successfully"})
	go CleanupExpiredOTPs()

}

// SendOTPHandler handles OTP generation and sending
func (h fieldHandler)SendOTPMobHandler(w http.ResponseWriter, r *http.Request) {

	type Response struct {
		Message string `json:"message"`
	}

	var req model.Signup
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Generate OTP
	otp, err := utils.GenerateOTP(6)
	if err != nil {
		http.Error(w, "Failed to generate OTP", http.StatusInternalServerError)
		return
	}

	// Store OTP with expiration time
	otpStore.mu.Lock()
	otpStore.data[req.MobileNumber] = OTPData{
		OTP:        otp,
		ExpiryTime: time.Now().Add(5 * time.Minute), // Set expiry to 1 minute
	}
	otpStore.mu.Unlock()

	// Send OTP via mobile no
	// Define the external API URL
	url := "http://login4.spearuc.com/MOBILE_APPS_API/sms_api.php?type=smsquicksend&user=iitmadras&pass=welcome&sender=EVOLGN&t_id=1707166841244742343&to_mobileno="+req.MobileNumber+"&sms_text=Dear%20Applicant,%20Your%20OTP%20for%20Mobile%20No.%20Verification%20is%20"+otp+"%20-%20Prison%20Birds%20Tech%20.%20MJPTBCWREIS%20-%20EVOLGN%20" 

	// Perform the GET request
	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, "Failed to Send OTP", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Check if the status code is 200 (OK)
	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Failed to Send OTP in Status Code", http.StatusInternalServerError)
		return
		}

	json.NewEncoder(w).Encode(Response{Message: "OTP sent successfully"})
	go CleanupExpiredOTPs()

}

// VerifyOTPHandler handles OTP verification
func (h fieldHandler)VerifyOTPHandler(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		Email string `json:"email"`
		Mobile string `json:"mobile_number"`
		OTP   string `json:"otp"`
	}
	type Response struct {
		Valid bool   `json:"valid"`
		Error string `json:"error,omitempty"`
	}

	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Check OTP
	otpStore.mu.Lock()

	if req.Email ==""{
		otpData, exists := otpStore.data[req.Mobile]
		if !exists {
			otpStore.mu.Unlock()
			json.NewEncoder(w).Encode(Response{Valid: false, Error: "No OTP found"})
			return
		}
	
		// Check expiration time
		if time.Now().After(otpData.ExpiryTime) {
			delete(otpStore.data, req.Email) // Remove expired OTP
			otpStore.mu.Unlock()
			json.NewEncoder(w).Encode(Response{Valid: false, Error: "OTP expired"})
			return
		}

			// Check OTP value
	if otpData.OTP != req.OTP {
		otpStore.mu.Unlock()
		json.NewEncoder(w).Encode(Response{Valid: false, Error: "Invalid OTP"})
		return
	}
	}else{
	otpData, exists := otpStore.data[req.Email]
	if !exists {
		otpStore.mu.Unlock()
		json.NewEncoder(w).Encode(Response{Valid: false, Error: "No OTP found"})
		return
	}

	// Check expiration time
	if time.Now().After(otpData.ExpiryTime) {
		delete(otpStore.data, req.Email) // Remove expired OTP
		otpStore.mu.Unlock()
		json.NewEncoder(w).Encode(Response{Valid: false, Error: "OTP expired"})
		return
	}
		// Check OTP value
		if otpData.OTP != req.OTP {	
			otpStore.mu.Unlock()
			json.NewEncoder(w).Encode(Response{Valid: false, Error: "Invalid OTP"})
			return
		}

	}
	


	// OTP is valid, remove it from the store
	delete(otpStore.data, req.Email)
	otpStore.mu.Unlock()
	json.NewEncoder(w).Encode(Response{Valid: true})
	go CleanupExpiredOTPs()

}

// CleanupExpiredOTPs periodically removes expired OTPs from the store
func CleanupExpiredOTPs() {
		otpStore.mu.Lock()
		now := time.Now()
		for email, otpData := range otpStore.data {
			if now.After(otpData.ExpiryTime) {
				delete(otpStore.data, email)
			}
		}
		otpStore.mu.Unlock()
}


func (h fieldHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	req := model.Signup{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if req.Email ==""{
		query := "SELECT employee_id, first_name, last_name,mobile_number,  email, date_of_birth,gender, country_code, photo_url, coins FROM employees WHERE mobile_number = $1"
		// password:= ""
		err = h.sqlDB.QueryRow(query, req.MobileNumber).Scan(&req.EmployeeID,&req.FirstName,&req.LastName,&req.MobileNumber, &req.Email,&req.DateOfBirth,&req.Gender, &req.CountryCode, &req.PhotoUrl, &req.Coins)
		if err != nil {
			utils.ErrorResponse(w, "Invalid Mobile Number", http.StatusBadRequest)
			return
		}
	}else{
		query := "SELECT employee_id, first_name, last_name,mobile_number,  email, date_of_birth,gender, country_code, photo_url, coins FROM employees WHERE email = $1"
		// password:= ""
		err = h.sqlDB.QueryRow(query, req.Email).Scan(&req.EmployeeID,&req.FirstName,&req.LastName,&req.MobileNumber, &req.Email,&req.DateOfBirth,&req.Gender, &req.CountryCode, &req.PhotoUrl, &req.Coins)
		if err != nil {
			utils.ErrorResponse(w, "Invalid Email", http.StatusBadRequest)
			return
		}
	}

	req.Token, err = utils.GenerateJWT(req.Email)
	if err != nil {
		utils.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.ReturnResponse(w, http.StatusOK, req)
}



// Handle WebSocket connections
func (h fieldHandler)HandleConnections(w http.ResponseWriter, r *http.Request) {
	// Parse user ID from query parameters
	userID := r.URL.Query().Get("userId")
	if userID == "" {
		http.Error(w, "Missing userId", http.StatusBadRequest)
		return
	}

	// Upgrade the connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading connection:", err)
		return
	}
	defer conn.Close()

	// Add user to clients map
	mu.Lock()
	clients[userID] = conn
	mu.Unlock()
	fmt.Printf("User %s connected\n", userID)
	if userID != "abinash1411999@gmail.com" {
		go sendNotification("abinash1411999@gmail.com", userID+ "  Online")
	}
	// Listen for messages from the client
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf("User %s disconnected\n", userID)
			if userID != "abinash1411999@gmail.com" {
				go sendNotification("abinash1411999@gmail.com", userID+ "  Offline")
			}

			break
		}
		fmt.Printf("Received message from %s: %s\n", userID, msg)
	}
	// Remove user from clients map on disconnect
	mu.Lock()
	delete(clients, userID)
	mu.Unlock()

}

type Client struct {
	userID  string
	writer  http.ResponseWriter
	flusher http.Flusher
	done    chan struct{}
}

var sseClients = make(map[string][]Client)// userID → message channel





func (h fieldHandler)SSEHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userId")
	if userID == "" {
		http.Error(w, "Missing userId", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	done := r.Context().Done() // This will be closed when the client disconnects

	client := Client{
		userID:  userID,
		writer:  w,
		flusher: flusher,
		done:    make(chan struct{}),
	}

	sseClients[userID] = append(sseClients[userID], client)

	log.Printf("Client connected: %s", userID)

	notify := w.(http.CloseNotifier).CloseNotify()

	go func() {
		<-notify
		log.Printf("Client disconnected: %s", userID)
		removeClient(userID, client)
	}()

	// Keep the connection alive
	for {
		select {
		case <-done:
			return
		default:
			// heartbeat (optional)
		}
	}
}



func removeClient(userID string, client Client) {
	activeClients := sseClients[userID]
	updatedClients := make([]Client, 0)

	for _, c := range activeClients {
		if c.writer != client.writer {
			updatedClients = append(updatedClients, c)
		}
	}

	if len(updatedClients) == 0 {
		delete(sseClients, userID)
		log.Printf("No more connections for user: %s. Cache cleared.", userID)
		// 🔥 Clear user-specific cache here if needed
	} else {
		sseClients[userID] = updatedClients
	}
}


func  notifyUser(userID string, coins int) {
	for _, client := range sseClients[userID] {
		msg := fmt.Sprintf("data: {\"coins\": %d}\n\n", coins)
		fmt.Println(msg)
		client.writer.Write([]byte(msg))
		client.flusher.Flush()
	}
	
}