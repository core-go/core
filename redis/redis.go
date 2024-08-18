package redis

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	"log"
	"net/url"
	"os"
	"time"
)

func NewRedisPool(uri string) (*redis.Pool, error) {
	rUrl, _ := url.Parse(uri)
	redisPassword, _ := rUrl.User.Password()
	pool := &redis.Pool{
		Wait: true,
		Dial: func() (redis.Conn, error) {
			if redisPassword != "" {
				c, err := redis.Dial("tcp", rUrl.Host, redis.DialPassword(redisPassword))
				if err != nil {
					return nil, err
				}
				return c, nil
			}
			c, err := redis.Dial("tcp", rUrl.Host)
			if err != nil {
				return nil, err
			}
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
	return pool, nil
}

type Config struct {
	Url             string     `yaml:"url" mapstructure:"url" json:"url,omitempty" gorm:"column:url" bson:"url,omitempty" dynamodbav:"url,omitempty" firestore:"url,omitempty"`
	IdleTimeout     int64      `yaml:"idle_timeout" mapstructure:"idle_timeout" json:"idleTimeout,omitempty" gorm:"column:idletimeout" bson:"idleTimeout,omitempty" dynamodbav:"idleTimeout,omitempty" firestore:"idleTimeout,omitempty"`
	MaxConnLifetime int64      `yaml:"max_conn_lifetime" mapstructure:"max_conn_lifetime" json:"maxConnLifetime,omitempty" gorm:"column:maxconnlifetime" bson:"maxConnLifetime,omitempty" dynamodbav:"maxConnLifetime,omitempty" firestore:"maxConnLifetime,omitempty"`
	MaxActive       int        `yaml:"max_active" mapstructure:"max_active" json:"maxActive,omitempty" gorm:"column:maxactive" bson:"maxActive,omitempty" dynamodbav:"maxActive,omitempty" firestore:"maxActive,omitempty"`
	MaxIdle         int        `yaml:"max_idle" mapstructure:"max_idle" json:"maxIdle,omitempty" gorm:"column:maxIdle" bson:"maxIdle,omitempty" dynamodbav:"maxIdle,omitempty" firestore:"maxIdle,omitempty"`
	DB              int        `yaml:"db" mapstructure:"db" json:"db,omitempty" gorm:"column:db" bson:"db,omitempty" dynamodbav:"db,omitempty" firestore:"db,omitempty"`
	ConnectTimeout  int64      `yaml:"connect_timeout" mapstructure:"connect_timeout" json:"connectTimeout,omitempty" gorm:"column:connecttimeout" bson:"connectTimeout,omitempty" dynamodbav:"connectTimeout,omitempty" firestore:"connectTimeout,omitempty"`
	KeepAlive       int64      `yaml:"keep_alive" mapstructure:"keep_alive" json:"keepAlive,omitempty" gorm:"column:keepalive" bson:"keepAlive,omitempty" dynamodbav:"keepAlive,omitempty" firestore:"keepAlive,omitempty"`
	ReadTimeout     int64      `yaml:"read_timeout" mapstructure:"read_timeout" json:"readTimeout,omitempty" gorm:"column:readtimeout" bson:"readTimeout,omitempty" dynamodbav:"readTimeout,omitempty" firestore:"readTimeout,omitempty"`
	WriteTimeout    int64      `yaml:"write_timeout" mapstructure:"write_timeout" json:"writeTimeout,omitempty" gorm:"column:writetimeout" bson:"writeTimeout,omitempty" dynamodbav:"writeTimeout,omitempty" firestore:"writeTimeout,omitempty"`
	Wait            *bool      `yaml:"wait" mapstructure:"wait" json:"wait,omitempty" gorm:"column:wait" bson:"wait,omitempty" dynamodbav:"wait,omitempty" firestore:"wait,omitempty"`
	TLSEnable       *bool      `yaml:"tls_enable" mapstructure:"tls_enable" json:"tlsEnable,omitempty" gorm:"column:tlsenable" bson:"tlsEnable,omitempty" dynamodbav:"tlsEnable,omitempty" firestore:"tlsEnable,omitempty"`
	TLSSkipVerify   *bool      `yaml:"tls_skip_verify" mapstructure:"tls_skip_verify" json:"tlsSkipVerify,omitempty" gorm:"column:tlsskipverify" bson:"tlsSkipVerify,omitempty" dynamodbav:"tlsSkipVerify,omitempty" firestore:"tlsSkipVerify,omitempty"`
	TLS             *TLSConfig `yaml:"tls" mapstructure:"tls" json:"tls,omitempty" gorm:"column:tls" bson:"tls,omitempty" dynamodbav:"tls,omitempty" firestore:"tls,omitempty"`
}
type TLSConfig struct {
	InsecureSkipVerify *bool  `yaml:"insecure_skip_verify" mapstructure:"insecure_skip_verify" json:"insecureSkipVerify,omitempty" gorm:"column:insecureskipverify" bson:"insecureSkipVerify,omitempty" dynamodbav:"insecureSkipVerify,omitempty" firestore:"insecureSkipVerify,omitempty"`
	CertFile           string `yaml:"cert_file" mapstructure:"cert_file" json:"certFile,omitempty" gorm:"column:certfile" bson:"certFile,omitempty" dynamodbav:"certFile,omitempty" firestore:"certFile,omitempty"`
	KeyFile            string `yaml:"key_file" mapstructure:"key_file" json:"keyFile,omitempty" gorm:"column:keyfile" bson:"keyFile,omitempty" dynamodbav:"keyFile,omitempty" firestore:"keyFile,omitempty"`
	CaFile             string `yaml:"ca_file" mapstructure:"ca_file" json:"caFile,omitempty" gorm:"column:cafile" bson:"caFile,omitempty" dynamodbav:"caFile,omitempty" firestore:"caFile,omitempty"`
}

func NewDialOptions(c Config, pass ...string) []redis.DialOption {
	os := make([]redis.DialOption, 0)
	if len(pass) > 0 {
		o := redis.DialPassword(pass[0])
		os = append(os, o)
	}
	if c.DB > 0 {
		o := redis.DialDatabase(c.DB)
		os = append(os, o)
	}
	if c.ConnectTimeout > 0 {
		o := redis.DialConnectTimeout(time.Duration(c.ConnectTimeout) * time.Millisecond)
		os = append(os, o)
	}
	if c.KeepAlive > 0 {
		o := redis.DialKeepAlive(time.Duration(c.KeepAlive) * time.Millisecond)
		os = append(os, o)
	}
	if c.ReadTimeout > 0 {
		o := redis.DialReadTimeout(time.Duration(c.ReadTimeout) * time.Millisecond)
		os = append(os, o)
	}
	if c.WriteTimeout > 0 {
		o := redis.DialWriteTimeout(time.Duration(c.WriteTimeout) * time.Millisecond)
		os = append(os, o)
	}
	if c.TLSEnable != nil {
		o0 := redis.DialUseTLS(*c.TLSEnable)
		os = append(os, o0)
		if *c.TLSEnable == true {
			if c.TLSSkipVerify != nil {
				o := redis.DialTLSSkipVerify(*c.TLSSkipVerify)
				os = append(os, o)
			}
			if c.TLS != nil {
				tls := CreateTLSConfig(*c.TLS)
				o := redis.DialTLSConfig(tls)
				os = append(os, o)
			}
		}
	}
	return os
}
func CreateTLSConfig(c TLSConfig) (t *tls.Config) {
	t = &tls.Config{}
	if c.CertFile != "" && c.KeyFile != "" && c.CaFile != "" {
		cert, err := tls.LoadX509KeyPair(c.CertFile, c.KeyFile)
		if err != nil {
			log.Fatalf("%v", err)
		}

		caCert, err := os.ReadFile(c.CaFile)
		if err != nil {
			log.Fatalf("%v", err)
		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		t = &tls.Config{
			Certificates: []tls.Certificate{cert},
			RootCAs:      caCertPool,
		}
	}
	if c.InsecureSkipVerify != nil {
		t.InsecureSkipVerify = *c.InsecureSkipVerify
	}
	return t
}
func NewRedisPoolByConfig(c Config) (*redis.Pool, error) {
	rUrl, er0 := url.Parse(c.Url)
	if er0 != nil {
		return nil, er0
	}
	p := redis.Pool{}
	if c.MaxIdle > 0 {
		p.MaxIdle = c.MaxIdle
	}
	if c.MaxActive > 0 {
		p.MaxActive = c.MaxActive
	}
	if c.IdleTimeout > 0 {
		p.IdleTimeout = time.Duration(c.IdleTimeout) * time.Millisecond
	}
	if c.MaxConnLifetime > 0 {
		p.MaxConnLifetime = time.Duration(c.MaxConnLifetime) * time.Millisecond
	}
	var options []redis.DialOption
	redisPassword, ok := rUrl.User.Password()
	if ok {
		options = NewDialOptions(c, redisPassword)
	} else {
		options = NewDialOptions(c)
	}
	wait := true
	if c.Wait != nil {
		wait = *c.Wait
	}
	redisPool := &redis.Pool{
		MaxIdle:         p.MaxIdle,
		MaxActive:       p.MaxActive,
		Wait:            wait,
		IdleTimeout:     p.IdleTimeout,
		MaxConnLifetime: p.MaxConnLifetime,
		Dial: func() (redis.Conn, error) {
			client, err := redis.Dial("tcp", rUrl.Host, options...)
			if err != nil {
				return nil, err
			}
			return client, nil
		},
		TestOnBorrow: func(con redis.Conn, t time.Time) error {
			_, err := con.Do("PING")
			return err
		},
	}
	return redisPool, nil
}

func Set(pool *redis.Pool, key string, value interface{}, timeToLive time.Duration) error {
	conn := pool.Get()
	defer conn.Close()
	s, ok := value.(string)
	if ok {
		_, err := conn.Do("SET", key, s, "EX", int(timeToLive))
		return err
	} else {
		s2, ok2 := value.(*string)
		if ok2 {
			_, err := conn.Do("SET", key, s2, "EX", int(timeToLive))
			return err
		} else {
			valueJson, err := json.Marshal(value)
			if err != nil {
				return err
			}
			_, err = conn.Do("SET", key, valueJson, "EX", int(timeToLive))
			return err
		}
	}
}

func Expire(pool *redis.Pool, key string, timeToLive time.Duration) (bool, error) {
	conn := pool.Get()
	defer conn.Close()
	result, err := conn.Do("EXPIRE", key, int(timeToLive))
	if err != nil {
		return false, err
	}
	if rs, ok := result.(int64); ok && rs == 0 {
		return false, nil
	}
	return true, nil
}

func Get(pool *redis.Pool, key string) (string, error) {
	conn := pool.Get()
	defer conn.Close()
	result, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return "", err
	}

	return string(result), nil
}

func GetAndDecode(pool *redis.Pool, key string, obj interface{}) (string, error) {
	conn := pool.Get()
	defer conn.Close()
	result, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return "", err
	}
	er1 := json.Unmarshal(result, &obj)
	return string(result), er1
}

func Exists(pool *redis.Pool, key string) (bool, error) {
	conn := pool.Get()
	defer conn.Close()
	result, err := conn.Do("EXISTS", key)
	if err != nil {
		return false, err
	}
	if rs, ok := result.(int64); ok && rs == 0 {
		return false, nil
	}
	return true, nil
}

func Delete(pool *redis.Pool, key string) (bool, error) {
	conn := pool.Get()
	defer conn.Close()
	count, err := conn.Do("DEL", key)
	if err != nil {
		return false, err
	}
	if rs, ok := count.(int64); ok && rs == 0 {
		return false, nil
	}
	return true, nil
}

func Clear(pool *redis.Pool) error {
	conn := pool.Get()
	defer conn.Close()
	_, err := conn.Do("FLUSHDB")
	return err
}

func GetMany(pool *redis.Pool, keys []string) (map[string]string, []string, error) {
	result := make(map[string]string)
	keyNil := make([]string, 0)
	keysSlice := make([]interface{}, 0)
	for _, value := range keys {
		keysSlice = append(keysSlice, value)
	}
	conn := pool.Get()
	defer conn.Close()
	count, err := conn.Do("MGET", keysSlice...)
	if err != nil {
		return nil, nil, err
	}
	if mapp, ok := count.([]interface{}); ok {
		for i, key := range keys {
			if mapp[i] != nil {
				result[key] = string(mapp[i].([]uint8))
			} else {
				keyNil = append(keyNil, key)
			}
		}
	}
	return result, keyNil, nil
}

func Keys(pool *redis.Pool) ([]string, error) {
	conn := pool.Get()
	defer conn.Close()
	list, err := conn.Do("KEYS", "*")
	if err != nil {
		return nil, err
	}
	rs := make([]string, 0)
	if s, ok := list.([]interface{}); ok {
		for i := range s {
			if k, yes := s[i].([]uint8); yes {
				rs = append(rs, string(k))
			}
		}
	}

	return rs, err
}

func Count(pool *redis.Pool) (int64, error) {
	conn := pool.Get()
	defer conn.Close()
	list, err := conn.Do("KEYS", "*")
	if err != nil {
		return 0, err
	}
	if s, ok := list.([]interface{}); ok {
		count := len(s)
		return int64(count), nil
	}
	return 0, nil

}

func Size(pool *redis.Pool) (int64, error) {
	conn := pool.Get()
	defer conn.Close()
	size, err := conn.Do("DBSIZE")
	if err != nil {
		return 0, err
	}
	if s, ok := size.(int64); ok {
		return s, err
	}
	return 0, nil
}
