package snowflake

import (
	"math/rand"
	"strconv"
	"sync"
	"time"
)

const (
	epoch         = int64(1400000000000)              // 设置起始时间
	timestampBits = uint(41)                          // 时间戳占用位数
	workerIdBits  = uint(10)                          // 机器id所占位数
	sequenceBits  = uint(12)                          // 序列所占的位数
	timestampMax  = int64(-1 ^ (-1 << timestampBits)) // 时间戳最大值
	//workerIdMax    = int64(-1 ^ (-1 << workerIdBits))  // 支持的最大机器id数量
	sequenceMask   = int64(-1 ^ (-1 << sequenceBits)) // 支持的最大序列id数量
	workerIdShift  = sequenceBits                     // 机器id左移位数
	timestampShift = sequenceBits + workerIdBits      // 时间戳左移位数
)

var s *Snowflake

type Snowflake struct {
	sync.Mutex          //锁
	lastTimeStamp int64 //时间戳
	workerId      int64 //机器id
	sequence      int64 //序列
}

func Generate() string {
	if s == nil {
		s = &Snowflake{
			workerId: rand.Int63n(1024),
		}
	}
	return s.New()
}

func (s *Snowflake) New() string {
	s.Lock()
	now := time.Now().UnixNano() / 1000000
	if s.lastTimeStamp == now {
		s.sequence = (s.sequence + 1) & sequenceMask
		if s.sequence == 0 {
			for now <= s.lastTimeStamp {
				now = time.Now().UnixNano() / 1000000
			}
		}
	} else {
		s.sequence = 0
	}
	t := now - epoch
	if t > timestampMax {
		s.Unlock()
		return ""
	}
	s.lastTimeStamp = now
	s.Unlock()
	shift := (t)<<timestampShift | (s.workerId << workerIdShift) | (s.sequence)
	res := strconv.FormatInt(shift, 10)
	return res
}
