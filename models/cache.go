package models

import (
	"ChildrenMath/pkg/util"
	"ChildrenMath/service/emailverify"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"strconv"
	"time"
)

const (
	userPointsZset       = "user:points"
	cacheExpireTime      = time.Minute * 10
	emailVCodeExpireTime = time.Minute * emailverify.ExpireTime
)

type UserInfo struct {
	UserId   int    `json:"userid"`
	UserName string `json:"username"`
}

func generateUserPointsKeyForRead(userID int) string {
	return fmt.Sprintf("user:%d:points:r", userID)
}

func GetOwnPointsFromRedisWithSave(userID int, userName string) (int, error) {
	var (
		ownPoints int
		err       error
	)
	ownPoints, err = Rdb.Get(generateUserPointsKeyForRead(userID)).Int()
	if err != nil {
		ownPoints, err = GetPointsFromZsetInRedis(userID, userName)
		if err != nil {
			ownPoints, err = InitPointsKeysInRedis(userID, userName)
			return ownPoints, err
		}
		// 写回
		err = UpdatePoints(DB, userID, ownPoints)
		if err != nil {
			return 0, err
		}
		err = Rdb.Set(generateUserPointsKeyForRead(userID), ownPoints, cacheExpireTime).Err()
		if err != nil {
			return 0, err
		}
	}
	return ownPoints, nil
}

// InitPointsKeysInRedis
// 该方法在Zset中无该用户数据时调用，会在Redis中新建一个带超时的key，并在一个Zset类型的key中添加一个元素
func InitPointsKeysInRedis(userID int, userName string) (int, error) {
	points, err := GetUserPoints(DB, userID)
	if err != nil {
		return 0, nil
	}

	keyByte, _ := json.Marshal(&UserInfo{
		UserId:   userID,
		UserName: userName,
	})

	tx := Rdb.TxPipeline()
	err = tx.ZAdd(userPointsZset, redis.Z{
		Score:  float64(points),
		Member: util.Byte2Str(keyByte),
	}).Err()
	if err != nil {
		return 0, err
	}
	err = tx.Set(generateUserPointsKeyForRead(userID), points, cacheExpireTime).Err()
	if err != nil {
		return 0, err
	}
	_, err = tx.Exec()
	return points, err
}

func GetPointsFromZsetInRedis(userID int, userName string) (int, error) {
	keyByte, _ := json.Marshal(&UserInfo{
		UserId:   userID,
		UserName: userName,
	})

	points, err := Rdb.ZScore(userPointsZset, util.Byte2Str(keyByte)).Result()
	if err != nil {
		return 0, err
	}
	return int(points), nil
}

func IncreaseOwnPointsInRedis(userID, addPoints int, userName string) error {
	readKey := generateUserPointsKeyForRead(userID)
	tx := Rdb.TxPipeline()

	var (
		err        error
		currPoints int
	)

	currPoints, err = GetPointsFromZsetInRedis(userID, userName)
	if err != nil {
		return err
	}
	err = tx.Set(readKey, currPoints+addPoints, cacheExpireTime).Err()
	if err != nil {
		return err
	}

	err = tx.ZIncrBy(userPointsZset, float64(addPoints), strconv.Itoa(userID)).Err()
	if err != nil {
		_, err = InitPointsKeysInRedis(userID, userName)
		if err != nil {
			return err
		}
	}

	_, err = tx.Exec()
	return err
}

func GetPointsRankByIDFromRedis() ([]Rank, error) {
	z, err := Rdb.ZRevRangeWithScores(userPointsZset, 0, 9).Result()
	if err != nil {
		return nil, err
	}
	rank := make([]Rank, len(z))
	for i := 0; i < len(z); i++ {
		var temp UserInfo
		json.Unmarshal(util.Str2Byte(z[i].Member.(string)), &temp)
		rank[i].UserID = temp.UserId
		rank[i].UserName = temp.UserName
		rank[i].Points = int(z[i].Score)
	}
	return rank, nil
}

func generateEmailVCodeKey(keyEmail string) string {
	return fmt.Sprintf("vcode:%s", keyEmail)
}
func SaveVCodeInRedis(keyEmail, vCode string) error {
	var err error
	_, err = GetVCodeFromRedis(keyEmail)
	if err == nil {
		return emailverify.ErrSend
	}
	err = Rdb.Set(generateEmailVCodeKey(keyEmail), vCode, emailVCodeExpireTime).Err()
	return err
}

func GetVCodeFromRedis(keyEmail string) (string, error) {
	vCode, err := Rdb.Get(generateEmailVCodeKey(keyEmail)).Result()
	if err != nil {
		return "", err
	}
	return vCode, nil
}
