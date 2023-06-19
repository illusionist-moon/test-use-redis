package models

import (
	"ChildrenMath/pkg/util"
	"encoding/json"
	"fmt"
	"time"
)

func GenerateUserPointsKeyForRead(userID int) string {
	return fmt.Sprintf("user:%d:points:r", userID)
}

func GenerateUserPointsKeyForWrite(userID int) string {
	return fmt.Sprintf("user:%d:points:w", userID)
}

func GetPointsRankFromRedis() ([]Rank, error) {
	var (
		rankJson string
		rank     []Rank
		err      error
	)
	rankJson, err = Rdb.Get("rank").Result()
	if err != nil {
		rank, err = GetPointsRank(DB)
		if err != nil {
			return nil, err
		}
		jsonData, err := json.Marshal(rank)
		if err != nil {
			return nil, err
		}
		err = Rdb.Set("rank", jsonData, time.Minute*1).Err()
		if err != nil {
			return nil, err
		}
	} else {
		err = json.Unmarshal(util.Str2Byte(rankJson), &rank)
		if err != nil {
			return nil, err
		}
	}
	fmt.Println(rank)
	return rank, nil
}

func GetOwnPointsFromRedisWithSave(userID int) (int, error) {
	var (
		ownPoints int
		err       error
	)
	ownPoints, err = Rdb.Get(GenerateUserPointsKeyForRead(userID)).Int()
	if err != nil {
		ownPoints, err = Rdb.Get(GenerateUserPointsKeyForWrite(userID)).Int()
		if err != nil {
			ownPoints, err = GetUserPoints(DB, userID)
			if err != nil {
				return 0, err
			}
			err = Rdb.Set(GenerateUserPointsKeyForWrite(userID), ownPoints, -1).Err()
		}
		// 写回
		err = UpdatePoints(DB, userID, ownPoints)
		if err != nil {
			return 0, err
		}
		err = Rdb.Set(GenerateUserPointsKeyForRead(userID), ownPoints, time.Minute*1).Err()
		if err != nil {
			return 0, err
		}
	}
	return ownPoints, nil
}

func IncreaseOwnPointsInRedis(userID, points int) error {
	readKey := GenerateUserPointsKeyForRead(userID)
	writeKey := GenerateUserPointsKeyForWrite(userID)
	tx := Rdb.TxPipeline()

	var err error
	err = tx.IncrBy(readKey, int64(points)).Err()
	if err != nil {
		return err
	}
	err = tx.Expire(readKey, time.Minute*1).Err()
	if err != nil {
		return err
	}

	err = tx.IncrBy(writeKey, int64(points)).Err()
	if err != nil {
		return err
	}

	_, err = tx.Exec()
	return err
}
