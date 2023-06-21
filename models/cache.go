package models

import (
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

const (
	userPointsZset = "user:points"
	expireTime     = time.Minute * 10
)

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
		ownPoints, err = GetPointsFromZsetInRedis(userName)
		if err != nil {
			ownPoints, err = InitPointsKeysInRedis(userID, userName)
			return ownPoints, err
		}
		// 写回
		err = UpdatePoints(DB, userID, ownPoints)
		if err != nil {
			return 0, err
		}
		err = Rdb.Set(generateUserPointsKeyForRead(userID), ownPoints, expireTime).Err()
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
	tx := Rdb.TxPipeline()
	err = tx.ZAdd(userPointsZset, redis.Z{
		Score:  float64(points),
		Member: userName,
	}).Err()
	if err != nil {
		return 0, err
	}
	err = tx.Set(generateUserPointsKeyForRead(userID), points, expireTime).Err()
	if err != nil {
		return 0, err
	}
	_, err = tx.Exec()
	return points, err
}

func GetPointsFromZsetInRedis(userName string) (int, error) {
	points, err := Rdb.ZScore(userPointsZset, userName).Result()
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

	currPoints, err = GetPointsFromZsetInRedis(userName)
	if err != nil {
		return err
	}
	err = tx.Set(readKey, currPoints+addPoints, expireTime).Err()
	if err != nil {
		return err
	}

	err = tx.ZIncrBy(userPointsZset, float64(addPoints), userName).Err()
	if err != nil {
		_, err = InitPointsKeysInRedis(userID, userName)
		if err != nil {
			return err
		}
	}

	_, err = tx.Exec()
	return err
}

func GetPointsRankFromRedis() ([]Rank, error) {
	z, err := Rdb.ZRevRangeWithScores(userPointsZset, 0, 9).Result()
	if err != nil {
		return nil, err
	}
	rank := make([]Rank, len(z))
	for i := 0; i < len(z); i++ {
		rank[i].UserName = z[i].Member.(string)
		rank[i].Points = int(z[i].Score)
	}
	return rank, nil
}
