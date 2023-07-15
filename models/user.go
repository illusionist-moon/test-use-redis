package models

import (
	"gorm.io/gorm"
)

const (
	RankSize = 10
)

type User struct {
	ID       int    `json:"id" gorm:"column:id;primary_key;auto_increment"`
	UserName string `json:"username" gorm:"column:user_name"`
	Email    string `json:"email" gorm:"column:email;unique"`
	Password string `json:"password" gorm:"column:password"`
	Points   int    `json:"points" gorm:"column:points"`
}

type Rank struct {
	UserID   int    `json:"userid" gorm:"column:id"`
	UserName string `json:"username" gorm:"column:user_name"`
	Points   int    `json:"points" gorm:"column:points"`
}

func Exists(db *gorm.DB, email string) bool {
	var user User
	// 按照 用户名 进行查找
	err := db.Omit("id", "user_name", "email", "password", "points").Where("email = ?", email).First(&user).Error
	if err != nil {
		return false
	} else {
		return true
	}
}

func GetUserInfoByEmail(db *gorm.DB, email string) (int, string, string, error) {
	var user User
	// 按照 用户名 进行查找
	err := db.Select("id", "user_name", "password").Where("email = ?", email).First(&user).Error
	if err != nil {
		return 0, "", "", err
	}
	return user.ID, user.UserName, user.Password, nil
}

func GetUserNameByID(db *gorm.DB, id int) (string, error) {
	var user User
	// 按照 用户名 进行查找
	err := db.Select("user_name").Where("id = ?", id).First(&user).Error
	if err != nil {
		return "", err
	}
	return user.UserName, nil
}

func CreateUser(db *gorm.DB, userName, email, password string) (int, error) {
	user := &User{
		UserName: userName,
		Email:    email,
		Password: password,
		Points:   0,
	}
	err := db.Create(user).Error
	if err != nil {
		return 0, err
	}
	return user.ID, nil
}

func UpdatePassword(db *gorm.DB, email, newPassword string) error {
	err := db.Model(&User{}).Where("email = ?", email).Update("password", newPassword).Error
	return err
}

func UpdatePoints(db *gorm.DB, userID int, points int) error {
	err := db.Model(&User{}).Where("id = ?", userID).Update("points", points).Error
	return err
}

func GetUserPoints(db *gorm.DB, userID int) (int, error) {
	var res int
	err := db.Model(&User{}).Select("points").Where("id = ?", userID).Take(&res).Error
	if err != nil {
		return 0, err
	}
	return res, nil
}

func GetPointsRank(db *gorm.DB) ([]Rank, error) {
	var res []Rank
	err := db.Model(&User{}).Select("points, id, user_name").Limit(RankSize).Order("points desc, id").Find(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

func UpdateUserName(db *gorm.DB, id int, newUserName string) error {
	err := db.Model(&User{}).Where("id = ?", id).Update("user_name", newUserName).Error
	return err
}
