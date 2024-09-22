package postgresql

import (
	"context"
	"crud/model"
	"encoding/json"
	// "fmt"
	"github.com/go-redis/redis/v8"
	"database/sql"
)

func New() SqlManager {
	return &sqlService{}
}

type sqlService struct {
}

func (r sqlService) SetUserData(ctx context.Context, db *sql.DB, users model.User) error {
	var err error

		sqlStatement := `INSERT INTO users (name, age, mob_number, email, ) VALUES ($1, $2, $3, $4) RETURNING id`
        id := 0
		err = db.QueryRow(sqlStatement, users.Name, users.Age, users.MobNumber, users.Email).Scan(&id)
        if err != nil {
            return err
        }
    
	return err
}

// func (r sqlService) GetUserData(ctx context.Context, client *redis.Client, userData string) (model.User, error) {
// 	var users model.StudentInfo
// 	// Get data from Redis by user ID
// 	userJSON, err := client.Get(ctx, "students_information").Result()
// 	if err != nil {
// 		return model.User{}, err
// 	}
// 	// Parse JSON data into User struct
// 	err = json.Unmarshal([]byte(userJSON), &users)
// 	if err != nil {
// 		return model.User{}, err
// 	}
// 	for _, user := range users.Students {
// 		if user.ID == userData {
// 			return user, err

// 		}
// 	}

// 	return model.User{}, fmt.Errorf("User with ID %s does not exist", userData)
// }

// func (r sqlService) DeleteUserData(ctx context.Context, client *redis.Client, userID string) error {
// 	var users model.StudentInfo
// 	var data []model.User

// 	// Get data from Redis by user ID
// 	userJSON, err := client.Get(ctx, "students_information").Result()
// 	if err != nil {
// 		return err
// 	}
// 	// Parse JSON data into User struct
// 	err = json.Unmarshal([]byte(userJSON), &users)
// 	if err != nil {
// 		return err
// 	}

// 	for index, user := range users.Students {
// 		if user.ID != userID {
// 			data = append(data, users.Students[index])
// 		}
// 	}
// 	users.Students = data
// 	usersJson, err := json.Marshal(users)

// 	// Set data in Redis with key as user ID
// 	err = client.Set(ctx, "students_information", usersJson, 0).Err()
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (r sqlService) UpdateUserData(ctx context.Context, client *redis.Client, data model.User) error {

// 	var users model.StudentInfo
// 	// Get data from Redis by user ID
// 	userJSON, err := client.Get(ctx, "students_information").Result()
// 	if err != nil {
// 		return err
// 	}
// 	// Parse JSON data into User struct
// 	err = json.Unmarshal([]byte(userJSON), &users)
// 	if err != nil {
// 		return err
// 	}

// 	for index, user := range users.Students {
// 		if user.ID == data.ID {
// 			users.Students[index] = data
// 			userJSON, err := json.Marshal(users)

// 			// Set data in Redis with key as user ID
// 			err = client.Set(ctx, "students_information", userJSON, 0).Err()
// 			if err != nil {
// 				return err
// 			}

// 		}
// 	}
// 	return nil
// }

func (r sqlService) GetAll(ctx context.Context, client *redis.Client) (model.StudentInfo, error) {
	var users model.StudentInfo

	studentData, err := client.Get(ctx, "students_information").Result()
	if err != nil {
		return users, err
	}
	err = json.Unmarshal([]byte(studentData), &users)
	if err != nil {
		return model.StudentInfo{}, err
	}
	return users, nil
}
