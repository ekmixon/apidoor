package gateway

import (
	"context"
	"os"
	"regexp"

	"github.com/go-redis/redis/v8"
)

var rdb = redis.NewClient(&redis.Options{
	Addr:     os.Getenv("REDIS_HOST"),
	Password: "",
	DB:       0,
})

type Field struct {
	Re   *regexp.Regexp
	Path string
}

type KeyData map[string][]Field

var Data = make(KeyData)

func GetAPIURL(ctx context.Context, key, path string) (string, error) {
	fields, ok := Data[key]
	if !ok {
		return "", &MyError{Message: "unauthorized request"}
	}

	for _, v := range fields {
		if v.Re.Match([]byte(path)) {
			return v.Path[1:], nil
		}
	}

	return "", &MyError{Message: "unauthorized request"}
}

func init() {
	ctx := context.Background()
	for _, k := range rdb.Keys(ctx, "*").Val() {
		for _, hk := range rdb.HKeys(ctx, k).Val() {
			re, err := regexp.Compile(hk)
			if err != nil {
				panic(err)
			}
			Data[k] = append(Data[k], Field{
				Re:   re,
				Path: rdb.HGet(ctx, k, hk).Val(),
			})
		}
	}
}
