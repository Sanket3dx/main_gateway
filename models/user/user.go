package userModel

import (
	"database/sql"
	"fmt"
	mysql_configer "main_gateway/database"
	"main_gateway/utils"
	"time"

	"github.com/golang-jwt/jwt"
)

type LoginDetails struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserType struct {
	UserTypeID int    `db:"user_type_id" json:"user_type_id"`
	UserType   string `db:"user_type" json:"user_type"`
	UpdatedOn  string `db:"updated_on" json:"updated_on"`
	UpdatedBy  int    `db:"updated_by" json:"updated_by"`
	Deleted    int    `db:"deleted" json:"deleted"`
}

type User struct {
	UserID     int    `db:"user_id" json:"user_id"`
	UserTypeID int    `db:"user_type_id" json:"user_type_id"`
	ClinicID   int    `db:"clinic_id" json:"clinic_id"`
	Name       string `db:"name" json:"name"`
	Gender     string `db:"gender" json:"gender"`
	DOB        string `db:"dob" json:"dob"`
}

type UserLogin struct {
	UserLoginID int    `db:"user_login_id" json:"user_login_id"`
	Username    string `db:"username" json:"username"`
	Password    string `db:"password" json:"password"`
	UpdatedOn   string `db:"updated_on" json:"updated_on"`
	UpdatedBy   int    `db:"updated_by" json:"updated_by"`
	Status      int    `db:"status" json:"status"`
	Deleted     int    `db:"deleted" json:"deleted"`
}

func AuthenticateUser(username, password string) (int, error) {
	db := mysql_configer.InitDB() // Initialize your database connection
	defer db.Close()
	secKey := "6fcf6c4d0ebd8b0d5e96839d8a240740"
	password, _ = utils.EncryptToHex([]byte(password), []byte(secKey))
	fmt.Println(password)
	query := "SELECT user_login_id FROM user_logins WHERE username = ? AND password = ? AND status = 1"
	var userID int
	err := db.QueryRow(query, username, password).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("authentication failed")
		}
		// Database error
		return 0, err
	}
	return userID, nil
}

func GetUser(userID int) (*User, error) {
	db := mysql_configer.InitDB()
	defer db.Close()

	stmt, err := db.Prepare("SELECT user_type_id, clinic_id, name, gender, dob FROM Users WHERE user_id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(userID)
	var user User

	err = row.Scan(
		&user.UserTypeID,
		&user.ClinicID,
		&user.Name,
		&user.Gender,
		&user.DOB,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("User with ID %d not found", userID)
		}
		return nil, err
	}

	user.UserID = userID

	return &user, nil
}

func GenrateJwtWithClaims(user User) (string, error) {
	configs := utils.GetProjectConfig()
	secretKey := []byte(configs.SecretKey)
	claims := jwt.MapClaims{
		"user_id":      user.UserID,
		"user_type_id": user.UserTypeID,
		"exp":          time.Now().Add(time.Hour * 1).Unix(), // Token will expire in 1 hour
		"nbf":          time.Now().Unix(),                    // Token not valid before this time
		"iat":          time.Now().Unix(),                    // Token issued at this time
	}

	// Create a new JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with your secret key
	tokenString, err := token.SignedString(secretKey)
	fmt.Println(tokenString)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
