package main

import (
	"crud/app"
)

// // OTPData stores the OTP and its expiration time
// type OTPData struct {
// 	OTP        string
// 	ExpiryTime time.Time
// }

// var otpStore = struct {
// 	mu   sync.Mutex
// 	data map[string]OTPData // Stores OTPs and their expiration times, keyed by email
// }{data: make(map[string]OTPData)}

// // CleanupExpiredOTPs periodically removes expired OTPs from the store
// func CleanupExpiredOTPs() {
// 	ticker := time.NewTicker(2 * time.Minute) // Run every 2 minutes
// 	defer ticker.Stop()

// 	for range ticker.C {
// 		otpStore.mu.Lock()
// 		now := time.Now()
// 		for email, otpData := range otpStore.data {
// 			if now.After(otpData.ExpiryTime) {
// 				delete(otpStore.data, email)
// 			}
// 		}
// 		otpStore.mu.Unlock()
// 	}
// }

func main() {
	app.Start()
}
