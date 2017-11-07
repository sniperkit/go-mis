package services

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"bitbucket.org/go-mis/config"
	"github.com/go-redis/redis"
)

const (
	prefixRecLoan = "GOMIS_RECOMENDEDLOAN"
)

// RedisClient - store redis client driver and additional information
type RedisClient struct {
	*redis.Client
}

// NewClientRedis - instance new client redis
func NewClientRedis() (*RedisClient, error) {
	redisClient := new(RedisClient)
	port := strconv.Itoa(config.Configuration.Redis.Port)
	fmt.Println("Redis port: ", port)
	client := redis.NewClient(&redis.Options{
		Addr:     config.Configuration.Redis.Address + ":" + port,
		Password: config.Configuration.Redis.Password,
		DB:       config.Configuration.Redis.Db,
	})
	if _, err := client.Ping().Result(); err != nil {
		log.Println("[ERROR] Unable to connect to redis", err)
		panic(err)
	}
	redisClient.Client = client
	return redisClient, nil
}

// GetPRecomendedLoanKey - get prefix recomended loan that will be stored in redis
func (r *RedisClient) GetPRecomendedLoanKey(investorID string) (string, error) {
	if len(strings.TrimSpace(investorID)) == 0 {
		return "", errors.New("Investor ID can not be empty")
	}
	prefix := prefixRecLoan + "_" + investorID
	return prefix, nil
}

// GetRecomendedLoan - get recomended loan from redis
func (r *RedisClient) GetRecomendedLoan(key string) ([]byte, error) {
	if len(strings.TrimSpace(key)) == 0 {
		log.Println("[ERROR] Key can not be empty")
		return nil, errors.New("Key can not be empty")
	}
	b, err := r.Get(key).Bytes()
	if err != nil {
		log.Println("[ERROR] ", err)
		return nil, err
	}
	return b, nil
}

func (r *RedisClient) GetKeys(key string) ([]string, error) {
	if len(strings.TrimSpace(key)) == 0 {
		key = "*"
	}
	fmt.Println("Key: ", key)
	cmd := r.Keys("*" + key + "*")
	if cmd.Err() != nil {
		log.Println("[ERROR]", cmd.Err())
		return nil, cmd.Err()
	}
	return cmd.Val(), nil
}

// GetAllRecomendedLoan - get all data recomended loan from redis
func (r *RedisClient) GetAllRecomendedLoan() ([]string, error) {
	data := make([]string, 0)
	key := r.GetPrefixKeyRecomendedLoan()
	keyList, err := r.GetKeys(key)
	if err != nil {
		return nil, err
	}
	for i := range keyList {
		b, err := r.GetRecomendedLoan(keyList[i])
		if err != nil {
			return nil, err
		}
		data = append(data, string(b))
	}
	fmt.Println(data)
	return data, nil
}

// SaveRecomendedLoan - save recomended loan into redis
func (r *RedisClient) SaveRecomendedLoan(investorID string, loanByte []byte) error {
	if len(loanByte) == 0 {
		return errors.New("Loan can not be empty")
	}
	key, err := r.GetPRecomendedLoanKey(investorID)
	if err != nil {
		log.Println("[ERROR] ", err)
		return err
	}
	cmd := r.Set(key, string(loanByte), 5*time.Minute)
	if err := cmd.Err(); err != nil {
		return err
	}
	return nil
}

func (r *RedisClient) GetPrefixKeyRecomendedLoan() string {
	return prefixRecLoan
}
