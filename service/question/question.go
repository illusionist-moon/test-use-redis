package question

import (
	"math/rand"
	"time"
)

const (
	MultiMax = 10
	OtherMax = 100

	Count = 10
)

type Questions struct {
	Num1 uint32
	Num2 uint32
}

// GenerateQuestions api保证了传入该函数的operator一定是正确且合法的
func GenerateQuestions(operator string) []*Questions {
	var max uint32
	if operator == "*" {
		max = MultiMax
	} else {
		max = OtherMax
	}
	ans := make([]*Questions, 0)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < Count; i++ {
		n1 := r.Uint32() % max
		var n2 uint32
		if operator == "/" || operator == "-" {
			for {
				if n1 != 0 {
					break
				}
				n1 = r.Uint32() % max
			}
			n2 = r.Uint32() % n1
		} else {
			n2 = r.Uint32() % max
		}
		if operator == "/" {
			for {
				if n2 != 0 {
					break
				}
				n2 = r.Uint32() % n1
			}
		}
		ans = append(ans, &Questions{
			Num1: n1,
			Num2: n2,
		})
	}
	return ans
}

func Judge(nums []int, op string) bool {
	switch op {
	case "+":
		return nums[0]+nums[1] == nums[2]
	case "-":
		return nums[0]-nums[1] == nums[2]
	case "*":
		return nums[0]*nums[1] == nums[2]
	case "/":
		return nums[0]/nums[1] == nums[2]
	default:
		// 事实上，这一步永远不会走到
		return false
	}
}
