package oredis

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/vmihailenco/msgpack"
	"time"
)

type Client struct {
	New *redis.Client
}

var (
	Ser *Client
)

func New(option redis.Options) *Client {
	Ser = &Client{New: redis.NewClient(&option)}
	return Ser
}
func (r *Client) SetType(key string, val interface{}, expiration time.Duration) (err error) {
	b, err := msgpack.Marshal(val)
	if err != nil {
		return
	}
	_, err = r.New.Set(key, b, expiration).Result()
	if err != nil {
		return
	}
	return
}

func (r *Client) LPushExpire(key string, expiration time.Duration, values ...interface{}) (err error) {
	_, err = r.New.LPush(key, values...).Result()
	if err != nil {
		return
	}
	r.New.Expire(key, expiration)
	return
}
func (r *Client) LPush(key string, values ...interface{}) (err error) {
	_, err = r.New.LPush(key, values...).Result()
	if err != nil {
		return
	}
	return
}
func (r *Client) LLen(key string) (int64, error) {
	return r.New.LLen(key).Result()
}
func (r *Client) LRange(key string, start, stop int64) ([]string, error) {
	return r.New.LRange(key, start, stop).Result()
}

func (r *Client) HSETExpire(key string, expiration time.Duration, values ...interface{}) (err error) {
	_, err = r.New.LPush(key, values...).Result()
	if err != nil {
		return
	}
	r.New.Expire(key, expiration)
	return
}
func (r *Client) HSET(name, key string, value interface{}) (err error) {
	_, err = r.New.HSet(name, key, value).Result()
	if err != nil {
		return
	}
	return
}
func (r *Client) HGET(name, key string) (string, error) {
	return r.New.HGet(name, key).Result()
}
func (r *Client) HDel(name, key string) (int64, error) {
	return r.New.HDel(name, key).Result()
}
func (r *Client) HSETType(name, key string, value interface{}) (err error) {
	b, err := msgpack.Marshal(value)
	if err != nil {
		return
	}
	r.HSET(name, key, b)
	return
}
func (r *Client) HGETType(name, key string, value interface{}) error {
	b, err := r.HGET(name, key)
	if err != nil {
		return err
	}
	_ = msgpack.Unmarshal([]byte(b), &value)
	return nil
}
func (r *Client) HLEN(name string) (int64, error) {
	return r.New.HLen(name).Result()
}
func (r *Client) HSCAN(name string, offset, limit int64) (val map[string]string) {
	val = make(map[string]string)
	iter := r.New.HScan(name, uint64(offset), "", limit).Iterator()
	var i = 0
	var k string
	for iter.Next() {
		if i == 0 {
			k = iter.Val()
			val[k] = ""
			i++
		} else {
			val[k] = iter.Val()
			i = 0
		}
	}
	return val
}
func (r *Client) HSCANAll(name string) (val map[string]string, err error) {
	val = make(map[string]string)
	limit, err := r.HLEN(name)
	if err != nil {
		return
	}
	iter := r.New.HScan(name, 0, "", limit*2).Iterator()
	var i = 0
	var k string
	for iter.Next() {
		if i == 0 {
			k = iter.Val()
			val[k] = ""
			i++
		} else {
			val[k] = iter.Val()
			i = 0
		}
	}
	return val, nil
}

func (r *Client) Set(key string, val string, expiration time.Duration) (err error) {
	_, err = r.New.Set(key, val, expiration).Result()
	if err != nil {
		return
	}
	return
}

func (r *Client) SetKVType(key string, val interface{}) error {
	return r.SetType(key, val, 0)
}
func (r *Client) SetKV(key, val string) error {
	return r.Set(key, val, 0)
}

func (r *Client) Get(key string) (string, error) {
	return r.New.Get(key).Result()
}
func (r *Client) GetType(key string, v interface{}) (err error) {
	b, err := r.New.Get(key).Bytes()
	if err != nil {
		return
	}
	_ = msgpack.Unmarshal(b, &v)
	return
}
func (r *Client) Del(keys ...string) (int64, error) {
	return r.New.Del(keys...).Result()
}
func (r *Client) Keys(keys string) ([]string, error) {
	return r.New.Keys(keys).Result()
}

func (r *Client) SetNX(key string, val interface{}, expiration time.Duration) (bool, error) {
	return r.New.SetNX(key, val, expiration).Result()
}
func (r *Client) TTL(keys string) (time.Duration, error) {
	return r.New.TTL(keys).Result()
} // Redis 剩下有效时间 单位秒
func (r *Client) Close() error {
	return r.New.Close()
}

func (r *Client) GetTimeDataTypeByTime(key string, v interface{}, t int64) (err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
			err = errors.New("error")
			return
		}
	}()
	var data TimeData
	err = r.GetType(key, &data)
	if err != nil {
		return
	}
	if t < data.StartTime.Unix() {
		return errors.New(fmt.Sprintf("未达到开始时间,还剩余%d", data.StartTime.Unix()-t))
	}
	if t > data.EndTime.Unix() {
		return errors.New(fmt.Sprintf("已超时"))
	}
	_ = msgpack.Unmarshal(data.Data, &v)
	return
} // redis 保存有效时间
func (r *Client) GetTimeDataType(key string, v interface{}) (err error) {
	return r.GetTimeDataTypeByTime(key, v, time.Now().Unix())
} // redis 保存有效时间

func (r *Client) SetTimeData(key string, val interface{}, start, endTime time.Time) (err error) {
	b, err := msgpack.Marshal(val)
	if err != nil {
		return
	}
	var td = TimeData{
		EndTime:   endTime,
		StartTime: start,
		Data:      b,
	}
	t := time.Now().Unix()
	t = endTime.Unix() - t + 60
	err = r.SetType(key, td, time.Duration(t)*time.Second)
	if err != nil {
		return
	}
	return
}

type TimeData struct {
	EndTime   time.Time //到期时间
	StartTime time.Time // 开始时间
	Data      []byte    // 数据
} // 时间对象存入redis，用来判断诗句超时现象
