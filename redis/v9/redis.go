package redis

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"net/url"
	"time"
)

type FailoverConfig struct {
	MasterName              *string        `yaml:"master_name" mapstructure:"master_name" json:"masterName,omitempty" gorm:"column:mastername" bson:"masterName,omitempty" dynamodbav:"masterName,omitempty" firestore:"masterName,omitempty"`
	SentinelAddrs           []string       `yaml:"sentinel_addrs" mapstructure:"sentinel_addrs" json:"sentinelAddrs,omitempty" gorm:"column:sentineladdrs" bson:"sentinelAddrs,omitempty" dynamodbav:"sentinelAddrs,omitempty" firestore:"sentinelAddrs,omitempty"`
	Password                *string        `yaml:"password" mapstructure:"password" json:"password,omitempty" gorm:"column:password" bson:"password,omitempty" dynamodbav:"password,omitempty" firestore:"password,omitempty"`
	SentinelPassword        *string        `yaml:"sentinel_password" mapstructure:"sentinel_password" json:"sentinelPassword,omitempty" gorm:"column:sentinelpassword" bson:"sentinelPassword,omitempty" dynamodbav:"sentinelPassword,omitempty" firestore:"sentinelPassword,omitempty"`
	MaxRetries              *int           `yaml:"max_retries" mapstructure:"max_retries" json:"maxRetries,omitempty" gorm:"column:maxretries" bson:"maxRetries,omitempty" dynamodbav:"maxRetries,omitempty" firestore:"maxRetries,omitempty"`
	DialTimeout             *time.Duration `yaml:"dial_timeout" mapstructure:"dial_timeout" json:"dialTimeout,omitempty" gorm:"column:dialtimeout" bson:"dialTimeout,omitempty" dynamodbav:"dialTimeout,omitempty" firestore:"dialTimeout,omitempty"`
	SentinelUsername        *string        `yaml:"sentinel_username" mapstructure:"sentinel_username" json:"sentinelUsername,omitempty" gorm:"column:sentinelusername" bson:"sentinelUsername,omitempty" dynamodbav:"sentinelUsername,omitempty" firestore:"sentinelUsername,omitempty"`
	Username                *string        `yaml:"username" mapstructure:"username" json:"username,omitempty" gorm:"column:username" bson:"username,omitempty" dynamodbav:"username,omitempty" firestore:"username,omitempty"`
	RouteByLatency          *bool          `yaml:"route_by_latency" mapstructure:"route_by_latency" json:"routeByLatency,omitempty" gorm:"column:routebylatency" bson:"routeByLatency,omitempty" dynamodbav:"routeByLatency,omitempty" firestore:"routeByLatency,omitempty"`
	RouteRandomly           *bool          `yaml:"route_randomly" mapstructure:"route_randomly" json:"routeRandomly,omitempty" gorm:"column:routerandomly" bson:"routeRandomly,omitempty" dynamodbav:"routeRandomly,omitempty" firestore:"routeRandomly,omitempty"`
	ReplicaOnly             *bool          `yaml:"replica_only" mapstructure:"replica_only" json:"replicaOnly,omitempty" gorm:"column:replicaonly" bson:"replicaOnly,omitempty" dynamodbav:"replicaOnly,omitempty" firestore:"replicaOnly,omitempty"`
	UseDisconnectedReplicas *bool          `yaml:"use_disconnected_replicas" mapstructure:"use_disconnected_replicas" json:"useDisconnectedReplicas,omitempty" gorm:"column:usedisconnectedreplicas" bson:"useDisconnectedReplicas,omitempty" dynamodbav:"useDisconnectedReplicas,omitempty" firestore:"useDisconnectedReplicas,omitempty"`
	DB                      *int           `yaml:"db" mapstructure:"db" json:"db,omitempty" gorm:"column:db" bson:"db,omitempty" dynamodbav:"db,omitempty" firestore:"db,omitempty"`
	ReadTimeout             *time.Duration `yaml:"read_timeout" mapstructure:"read_timeout" json:"readTimeout,omitempty" gorm:"column:readtimeout" bson:"readTimeout,omitempty" dynamodbav:"readTimeout,omitempty" firestore:"readTimeout,omitempty"`
	WriteTimeout            *time.Duration `yaml:"write_timeout" mapstructure:"write_timeout" json:"writeTimeout,omitempty" gorm:"column:writetimeout" bson:"writeTimeout,omitempty" dynamodbav:"writeTimeout,omitempty" firestore:"writeTimeout,omitempty"`
	MaxRetryBackoff         *time.Duration `yaml:"max_retry_backoff" mapstructure:"max_retry_backoff" json:"maxRetryBackoff,omitempty" gorm:"column:maxretrybackoff" bson:"maxRetryBackoff,omitempty" dynamodbav:"maxRetryBackoff,omitempty" firestore:"maxRetryBackoff,omitempty"`
	MinRetryBackoff         *time.Duration `yaml:"min_retry_backoff" mapstructure:"min_retry_backoff" json:"minRetryBackoff,omitempty" gorm:"column:minretrybackoff" bson:"minRetryBackoff,omitempty" dynamodbav:"minRetryBackoff,omitempty" firestore:"minRetryBackoff,omitempty"`
	ContextTimeoutEnabled   *bool          `yaml:"context_timeout_enabled" mapstructure:"context_timeout_enabled" json:"contextTimeoutEnabled,omitempty" gorm:"column:contexttimeoutenabled" bson:"contextTimeoutEnabled,omitempty" dynamodbav:"contextTimeoutEnabled,omitempty" firestore:"contextTimeoutEnabled,omitempty"`
	PoolSize                *int           `yaml:"pool_size" mapstructure:"pool_size" json:"poolSize,omitempty" gorm:"column:poolsize" bson:"poolSize,omitempty" dynamodbav:"poolSize,omitempty" firestore:"poolSize,omitempty"`
	PoolTimeout             *time.Duration `yaml:"pool_timeout" mapstructure:"pool_timeout" json:"poolTimeout,omitempty" gorm:"column:pooltimeout" bson:"poolTimeout,omitempty" dynamodbav:"poolTimeout,omitempty" firestore:"poolTimeout,omitempty"`
	MinIdleConns            *int           `yaml:"min_idle_conns" mapstructure:"min_idle_conns" json:"minIdleConns,omitempty" gorm:"column:minidleconns" bson:"minIdleConns,omitempty" dynamodbav:"minIdleConns,omitempty" firestore:"minIdleConns,omitempty"`
	MaxIdleConns            *int           `yaml:"max_idle_conns" mapstructure:"max_idle_conns" json:"maxIdleConns,omitempty" gorm:"column:maxidleconns" bson:"maxIdleConns,omitempty" dynamodbav:"maxIdleConns,omitempty" firestore:"maxIdleConns,omitempty"`
	ConnMaxIdleTime         *time.Duration `yaml:"conn_max_idle_time" mapstructure:"conn_max_idle_time" json:"connMaxIdleTime,omitempty" gorm:"column:connmaxidletime" bson:"connMaxIdleTime,omitempty" dynamodbav:"connMaxIdleTime,omitempty" firestore:"connMaxIdleTime,omitempty"`
	ConnMaxLifetime         *time.Duration `yaml:"conn_max_lifetime" mapstructure:"conn_max_lifetime" json:"connMaxLifetime,omitempty" gorm:"column:connmaxlifetime" bson:"connMaxLifetime,omitempty" dynamodbav:"connMaxLifetime,omitempty" firestore:"connMaxLifetime,omitempty"`
}

func NewFailoverClient(c FailoverConfig) *redis.Client {
	ops := &redis.FailoverOptions{
		SentinelAddrs: c.SentinelAddrs,
	}
	if c.MasterName != nil {
		ops.MasterName = *c.MasterName
	}
	if c.MaxRetries != nil {
		ops.MaxRetries = *c.MaxRetries
	}
	if c.DialTimeout != nil {
		ops.DialTimeout = *c.DialTimeout
	}
	if c.Password != nil {
		ops.Password = *c.Password
	}
	if c.SentinelPassword != nil {
		ops.SentinelPassword = *c.SentinelPassword
	}
	if c.SentinelUsername != nil {
		ops.SentinelUsername = *c.SentinelUsername
	}
	if c.Username != nil {
		ops.Username = *c.Username
	}
	if c.RouteByLatency != nil {
		ops.RouteByLatency = *c.RouteByLatency
	}
	if c.RouteRandomly != nil {
		ops.RouteRandomly = *c.RouteRandomly
	}
	if c.ReplicaOnly != nil {
		ops.ReplicaOnly = *c.ReplicaOnly
	}
	if c.UseDisconnectedReplicas != nil {
		ops.UseDisconnectedReplicas = *c.UseDisconnectedReplicas
	}
	if c.DB != nil {
		ops.DB = *c.DB
	}
	if c.ReadTimeout != nil {
		ops.ReadTimeout = *c.ReadTimeout
	}
	if c.WriteTimeout != nil {
		ops.WriteTimeout = *c.WriteTimeout
	}
	if c.MaxRetryBackoff != nil {
		ops.MaxRetryBackoff = *c.MaxRetryBackoff
	}
	if c.MinRetryBackoff != nil {
		ops.MinRetryBackoff = *c.MinRetryBackoff
	}
	if c.ContextTimeoutEnabled != nil {
		ops.ContextTimeoutEnabled = *c.ContextTimeoutEnabled
	}
	if c.PoolSize != nil {
		ops.PoolSize = *c.PoolSize
	}
	if c.PoolTimeout != nil {
		ops.PoolTimeout = *c.PoolTimeout
	}
	if c.MinIdleConns != nil {
		ops.MinIdleConns = *c.MinIdleConns
	}
	if c.MaxIdleConns != nil {
		ops.MaxIdleConns = *c.MaxIdleConns
	}
	if c.ConnMaxIdleTime != nil {
		ops.ConnMaxIdleTime = *c.ConnMaxIdleTime
	}
	if c.ConnMaxLifetime != nil {
		ops.ConnMaxLifetime = *c.ConnMaxLifetime
	}
	rc := redis.NewFailoverClient(ops)
	return rc
}

type Config struct {
	Url             string         `yaml:"url" mapstructure:"url" json:"url,omitempty" gorm:"column:url" bson:"url,omitempty" dynamodbav:"url,omitempty" firestore:"url,omitempty"`
	Addr            string         `yaml:"addr" mapstructure:"addr" json:"addr,omitempty" gorm:"column:addr" bson:"addr,omitempty" dynamodbav:"addr,omitempty" firestore:"addr,omitempty"`
	Password        *string        `yaml:"password" mapstructure:"password" json:"password,omitempty" gorm:"column:password" bson:"password,omitempty" dynamodbav:"password,omitempty" firestore:"password,omitempty"`
	TLSConfig       *bool          `yaml:"tls_config" mapstructure:"tls_config" json:"tlsConfig,omitempty" gorm:"column:tls_config" bson:"tlsConfig,omitempty" dynamodbav:"tlsConfig,omitempty" firestore:"tlsConfig,omitempty"`
	MaxRetries      int            `yaml:"max_retries" mapstructure:"max_retries" json:"maxRetries,omitempty" gorm:"column:maxretries" bson:"maxRetries,omitempty" dynamodbav:"maxRetries,omitempty" firestore:"maxRetries,omitempty"`
	PoolSize        int            `yaml:"pool_size" mapstructure:"pool_size" json:"poolSize,omitempty" gorm:"column:poolsize" bson:"poolSize,omitempty" dynamodbav:"poolSize,omitempty" firestore:"poolSize,omitempty"`
	DialTimeout     *time.Duration `yaml:"dial_timeout" mapstructure:"dial_timeout" json:"dialTimeout,omitempty" gorm:"column:dialtimeout" bson:"dialTimeout,omitempty" dynamodbav:"dialTimeout,omitempty" firestore:"dialTimeout,omitempty"`
	PoolTimeout     *time.Duration `yaml:"pool_timeout" mapstructure:"pool_timeout" json:"poolTimeout,omitempty" gorm:"column:pooltimeout" bson:"poolTimeout,omitempty" dynamodbav:"poolTimeout,omitempty" firestore:"poolTimeout,omitempty"`
	ReadTimeout     *time.Duration `yaml:"read_timeout" mapstructure:"read_timeout" json:"readTimeout,omitempty" gorm:"column:readtimeout" bson:"readTimeout,omitempty" dynamodbav:"readTimeout,omitempty" firestore:"readTimeout,omitempty"`
	WriteTimeout    *time.Duration `yaml:"write_timeout" mapstructure:"write_timeout" json:"writeTimeout,omitempty" gorm:"column:writetimeout" bson:"writeTimeout,omitempty" dynamodbav:"writeTimeout,omitempty" firestore:"writeTimeout,omitempty"`
	MaxRetryBackoff *time.Duration `yaml:"max_retry_backoff" mapstructure:"max_retry_backoff" json:"maxRetryBackoff,omitempty" gorm:"column:maxretrybackoff" bson:"maxRetryBackoff,omitempty" dynamodbav:"maxRetryBackoff,omitempty" firestore:"maxRetryBackoff,omitempty"`
	MinRetryBackoff *time.Duration `yaml:"min_retry_backoff" mapstructure:"min_retry_backoff" json:"minRetryBackoff,omitempty" gorm:"column:minretrybackoff" bson:"minRetryBackoff,omitempty" dynamodbav:"minRetryBackoff,omitempty" firestore:"minRetryBackoff,omitempty"`
	MinIdleConns    int            `yaml:"min_idle_conns" mapstructure:"min_idle_conns" json:"minIdleConns,omitempty" gorm:"column:minidleconns" bson:"minIdleConns,omitempty" dynamodbav:"minIdleConns,omitempty" firestore:"minIdleConns,omitempty"`
	DB              int            `yaml:"db" mapstructure:"db" json:"db,omitempty" gorm:"column:db" bson:"db,omitempty" dynamodbav:"db,omitempty" firestore:"db,omitempty"`
}

func NewRedisClientByConfig(c Config) (*redis.Client, error) {
	options := redis.Options{}
	if len(c.Url) > 0 {
		rUrl, er1 := url.Parse(c.Url)
		if er1 != nil {
			return nil, er1
		}
		redisPassword, _ := rUrl.User.Password()
		options.Addr = rUrl.Host
		options.Password = redisPassword
	}
	if len(c.Addr) > 0 {
		options.Addr = c.Addr
	}
	if c.Password != nil && len(*c.Password) > 0 {
		options.Password = *c.Password
	}
	if c.TLSConfig != nil && *c.TLSConfig == true {
		options.TLSConfig = &tls.Config{MinVersion: tls.VersionTLS12}
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
	if c.DialTimeout != nil {
		options.DialTimeout = *c.DialTimeout
	}
	if c.PoolTimeout != nil {
		options.PoolTimeout = *c.PoolTimeout
	}
	if c.ReadTimeout != nil {
		options.ReadTimeout = *c.ReadTimeout
	}
	if c.WriteTimeout != nil {
		options.WriteTimeout = *c.WriteTimeout
	}
	if c.MaxRetryBackoff != nil {
		options.MaxRetryBackoff = *c.MaxRetryBackoff
	}
	if c.MinRetryBackoff != nil {
		options.MinRetryBackoff = *c.MinRetryBackoff
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
	var v string
	v, ok := value.(string)
	if ok == false {
		json, err := json.Marshal(value)
		if err != nil {
			return err
		}
		v = string(json)
	}
	status := client.Set(ctx, key, v, timeToLive)
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
func Random(ctx context.Context, client *redis.Client) (key string, value string, err error) {
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
func Keys(ctx context.Context, client *redis.Client) (keys []string, err error) {
	cmd := client.Do(ctx, "KEYS", "*")
	err = cmd.Err()
	if err != nil {
		return nil, err
	}
	args := cmd.Val()
	if argsStr, ok := args.([]interface{}); ok {
		keys = make([]string, len(argsStr))
		for index, key := range argsStr {
			keys[index] = fmt.Sprint(key)
		}
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
