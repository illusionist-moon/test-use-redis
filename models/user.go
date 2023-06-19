package models

import (
	"gorm.io/gorm"
)

const (
	RankSize = 10
)

type User struct {
	ID       int    `json:"id" gorm:"column:id;primary_key;auto_increment"`
	UserName string `json:"username" gorm:"column:user_name;unique"`
	Password string `json:"password" gorm:"column:password"`
	Points   int    `json:"points" gorm:"column:points"`
}

type Rank struct {
	UserName string `json:"username" gorm:"column:user_name"`
	Points   int    `json:"points" gorm:"column:points"`
}

func Exists(db *gorm.DB, username string) (int, string, bool) {
	var user User
	// 按照 用户名 进行查找
	err := db.Select("id", "password").Where("user_name = ?", username).First(&user).Error
	if err != nil {
		return 0, "", false
	} else {
		return user.ID, user.Password, true
	}
}

func CreateUser(db *gorm.DB, username, password string) error {
	user := &User{
		UserName: username,
		Password: password,
		Points:   0,
	}
	err := db.Create(user).Error
	if err != nil {
		return err
	}
	return nil
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
	err := db.Model(&User{}).Select("points, user_name").Limit(RankSize).Order("points desc, user_name").Find(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}
