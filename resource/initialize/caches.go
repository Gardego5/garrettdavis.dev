package initialize

import (
	"github.com/Gardego5/garrettdavis.dev/model"
	"github.com/Gardego5/garrettdavis.dev/utils/bimarshal"
	"github.com/google/go-github/v66/github"
	"github.com/redis/go-redis/v9"
)

// registry for all the caches to add to the context
func Caches(rdb *redis.Client) bimarshal.RegisteredCaches {
	return bimarshal.Caches{
		"user":         bimarshal.Register[github.User](bimarshal.JSON),
		"access-token": bimarshal.Register[model.GHAccessToken](bimarshal.MessagePack),
		"subject":      bimarshal.Register[model.Subject](bimarshal.MessagePack),
	}.Build(rdb)
}
