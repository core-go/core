package v8

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"net/url"
	"time"
)

type Config struct {
	Url                string `mapstructure:"url" json:"url,omitempty" gorm:"column:url" bson:"url,omitempty" dynamodbav:"url,omitempty" firestore:"url,omitempty"`
	MaxRetries         int    `mapstructure:"max_retries" json:"maxRetries,omitempty" gorm:"column:maxretries" bson:"maxRetries,omitempty" dynamodbav:"maxRetries,omitempty" firestore:"maxRetries,omitempty"`
	PoolSize           int    `mapstructure:"pool_size" json:"poolSize,omitempty" gorm:"column:poolsize" bson:"poolSize,omitempty" dynamodbav:"poolSize,omitempty" firestore:"poolSize,omitempty"`
	IdleTimeout        int64  `mapstructure:"idle_timeout" json:"idleTimeout,omitempty" gorm:"column:idletimeout" bson:"idleTimeout,omitempty" dynamodbav:"idleTimeout,omitempty" firestore:"idleTimeout,omitempty"`
	DialTimeout        int64  `mapstructure:"dial_timeout" json:"dialTimeout,omitempty" gorm:"column:dialtimeout" bson:"dialTimeout,omitempty" dynamodbav:"dialTimeout,omitempty" firestore:"dialTimeout,omitempty"`
	PoolTimeout        int64  `mapstructure:"pool_timeout" json:"poolTimeout,omitempty" gorm:"column:pooltimeout" bson:"poolTimeout,omitempty" dynamodbav:"poolTimeout,omitempty" firestore:"poolTimeout,omitempty"`
	ReadTimeout        int64  `mapstructure:"read_timeout" json:"readTimeout,omitempty" gorm:"column:readtimeout" bson:"readTimeout,omitempty" dynamodbav:"readTimeout,omitempty" firestore:"readTimeout,omitempty"`
	WriteTimeout       int64  `mapstructure:"write_timeout" json:"writeTimeout,omitempty" gorm:"column:writetimeout" bson:"writeTimeout,omitempty" dynamodbav:"writeTimeout,omitempty" firestore:"writeTimeout,omitempty"`
	MaxConnAge         int64  `mapstructure:"max_conn_age" json:"maxConnAge,omitempty" gorm:"column:maxconnage" bson:"maxConnAge,omitempty" dynamodbav:"maxConnAge,omitempty" firestore:"maxConnAge,omitempty"`
	IdleCheckFrequency int64  `mapstructure:"idle_check_frequency" json:"idleCheckFrequency,omitempty" gorm:"column:idlecheckfrequency" bson:"idleCheckFrequency,omitempty" dynamodbav:"idleCheckFrequency,omitempty" firestore:"idleCheckFrequency,omitempty"`
	MaxRetryBackoff    int64  `mapstructure:"max_retry_backoff" json:"maxRetryBackoff,omitempty" gorm:"column:maxretrybackoff" bson:"maxRetryBackoff,omitempty" dynamodbav:"maxRetryBackoff,omitempty" firestore:"maxRetryBackoff,omitempty"`
	MinRetryBackoff    int64  `mapstructure:"min_retry_backoff" json:"minRetryBackoff,omitempty" gorm:"column:minretrybackoff" bson:"minRetryBackoff,omitempty" dynamodbav:"minRetryBackoff,omitempty" firestore:"minRetryBackoff,omitempty"`
	MinIdleConns       int    `mapstructure:"min_idle_conns" json:"minIdleConns,omitempty" gorm:"column:minidleconns" bson:"minIdleConns,omitempty" dynamodbav:"minIdleConns,omitempty" firestore:"minIdleConns,omitempty"`
	DB                 int    `mapstructure:"db" json:"db,omitempty" gorm:"column:db" bson:"db,omitempty" dynamodbav:"db,omitempty" firestore:"db,omitempty"`
}

func NewRedisClientByConfig(c Config) (*redis.Client, error) {
	rUrl, er1 := url.Parse(c.Url)
	if er1 != nil {
		return nil, er1
	}
	redisPassword, _ := rUrl.User.Password()
	options := redis.Options{
		Addr:     rUrl.Host,
		Password: redisPassword,
	}
	if c.DB > 0 {
		options.DB = c.DB
	}
	if c.MaxRetries > 0 {
		options.MaxRetries = c.MaxRetries
	}
	if c.PoolSize > 0 {
		options.PoolSize = c.PoolSize
	}
	if c.IdleTimeout > 0 {
		options.IdleTimeout = time.Duration(c.IdleTimeout) * time.Millisecond
	}
	if c.DialTimeout > 0 {
		options.DialTimeout = time.Duration(c.DialTimeout) * time.Millisecond
	}
	if c.PoolTimeout > 0 {
		options.PoolTimeout = time.Duration(c.PoolTimeout) * time.Millisecond
	}
	if c.ReadTimeout > 0 {
		options.ReadTimeout = time.Duration(c.ReadTimeout) * time.Millisecond
	}
	if c.WriteTimeout > 0 {
		options.WriteTimeout = time.Duration(c.WriteTimeout) * time.Millisecond
	}
	if c.MaxConnAge > 0 {
		options.MaxConnAge = time.Duration(c.MaxConnAge) * time.Millisecond
	}
	if c.IdleCheckFrequency > 0 {
		options.IdleCheckFrequency = time.Duration(c.IdleCheckFrequency) * time.Millisecond
	}
	if c.MaxRetryBackoff > 0 {
		options.MaxRetryBackoff = time.Duration(c.MaxRetryBackoff) * time.Millisecond
	}
	if c.MinRetryBackoff > 0 {
		options.MinRetryBackoff = time.Duration(c.MinRetryBackoff) * time.Millisecond
	}
	if c.MinIdleConns > 0 {
		options.MinIdleConns = c.MinIdleConns
	}
	client := redis.NewClient(&options)
	return client, nil
}
func NewRedisClient(uri string) (*redis.Client, error) {
	rUrl, er1 := url.Parse(uri)
	if er1 != nil {
		return nil, er1
	}
	redisPassword, _ := rUrl.User.Password()
	redisDB := 0
	redisOptions := redis.Options{
		Addr:     rUrl.Host,
		Password: redisPassword,
		DB:       redisDB,
	}
	client := redis.NewClient(&redisOptions)
	return client, nil
}

func Set(ctx context.Context, client *redis.Client, key string, value interface{}, timeToLive time.Duration) error {
	status := client.Set(ctx, key, value, timeToLive)
	return status.Err()
}
func SetNX(ctx context.Context, client *redis.Client, key string, value interface{}, timeToLive time.Duration) error {
	status := client.SetNX(ctx, key, value, timeToLive)
	return status.Err()
}

func Expire(ctx context.Context, client *redis.Client, key string, timeToLive time.Duration) (bool, error) {
	return client.Expire(ctx, key, timeToLive).Result()
}

func Get(ctx context.Context, client *redis.Client, key string) (string, error) {
	res, err := client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}

	return res, nil
}
func GetAndDecode(ctx context.Context, client *redis.Client, key string, obj interface{}) (string, error) {
	res, er0 := client.Get(ctx, key).Result()
	if er0 != nil {
		return "", er0
	}
	er1 := json.Unmarshal([]byte(res), &obj)
	return res, er1
}
func Exists(ctx context.Context, client *redis.Client, key string) (bool, error) {
	result, err := client.Do(ctx, "EXISTS", key).Int()
	if err != nil {
		return false, err
	}
	if result == 0 {
		return false, nil
	}
	return true, nil
}
func GetMany(ctx context.Context, client *redis.Client, keys []string) (map[string]string, []string, error) {
	result := make(map[string]string)
	keyNil := make([]string, 0)
	res, err := client.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, nil, err
	}
	for i, key := range keys {
		if res[i] != nil {
			result[key] = res[i].(string)
		} else {
			keyNil = append(keyNil, key)
		}
	}
	return result, keyNil, nil
}
func Random(ctx context.Context, client *redis.Client) (key string, value string,err error) {
	key, err = client.RandomKey(ctx).Result()
	if err != nil {
		return
	}
	value, err = Get(ctx, client, key)
	return
}
func Delete(ctx context.Context, client *redis.Client, key string) (bool, error) {
	count, err := client.Do(ctx, "DEL", key).Int()
	if err != nil {
		return false, err
	}
	if count == 0 {
		return false, nil
	}
	return true, nil
}

func Clear(ctx context.Context, client *redis.Client) error {
	status := client.Do(ctx, "flushdb")
	return status.Err()
}
func Keys(ctx context.Context, client *redis.Client) ([]string, error) {
	cmd := client.Do(ctx, "KEYS", "*")
	err := cmd.Err()
	if err != nil {
		return nil, err
	}

	args := cmd.Args()
	keys := make([]string, len(args))
	for index, key := range args {
		keys[index] = fmt.Sprint(key)
	}

	return keys, nil
}

func Count(ctx context.Context, client *redis.Client) (int64, error) {
	cmd := client.Do(ctx, "KEYS", "*")
	err := cmd.Err()
	if err != nil {
		return 0, err
	}
	return int64(len(cmd.Args())), nil
}

func Size(ctx context.Context, client *redis.Client) (int64, error) {
	cmd := client.DBSize(ctx)
	return cmd.Result()
}
