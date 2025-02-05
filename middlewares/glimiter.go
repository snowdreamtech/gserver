package middlewares

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	libredis "github.com/redis/go-redis/v9"
	"github.com/snowdreamtech/gserver/pkg/configs"
	ghttp "github.com/snowdreamtech/gserver/pkg/net/http"
	"github.com/snowdreamtech/gserver/pkg/tools"

	limiter "github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	memory "github.com/ulule/limiter/v3/drivers/store/memory"
	sredis "github.com/ulule/limiter/v3/drivers/store/redis"
)

// CustomLimitReachedHandler is the Custom LimitReachedHandler used by a new Middleware.
func CustomLimitReachedHandler(c *gin.Context) {
	ghttp.NegotiateResponse(c, http.StatusTooManyRequests, ghttp.NewResponse(ghttp.TooManyRequests, "Too many requests, Please try again later.", nil))

	c.Abort()
}

// CustomErrorHandler is the Custom ErrorHandler used by a new Middleware.
func CustomErrorHandler(c *gin.Context, err error) {
	panic(err)
}

// RateLimiter RateLimiter
func RateLimiter() gin.HandlerFunc {
	tools.DebugPrintF("[INFO] Starting Middleware %s", "RateLimiter")

	app := configs.GetAppConfig()

	if app.RateLimiter == "" {
		return Empty()
	}

	var store limiter.Store

	// Define a limit rate to 4 requests per hour.
	// You can also use the simplified format "<limit>-<period>"", with the given
	// periods:
	//
	// * "S": second
	// * "M": minute
	// * "H": hour
	// * "D": day
	//
	// Examples:
	//
	// * 5 reqs/second: "5-S"
	// * 10 reqs/minute: "10-M"
	// * 1000 reqs/hour: "1000-H"
	// * 2000 reqs/day: "2000-D"
	//
	rate, err := limiter.NewRateFromFormatted(app.RateLimiter)

	if err != nil {
		tools.DebugPrintF(err.Error())
		return Empty()
	}

	store = limiterRedisStore()

	if store == nil {
		store = limiterInmemoryStore()
	}

	// Create a new middleware with the limiter instance.
	middleware := mgin.NewMiddleware(limiter.New(store, rate), mgin.WithErrorHandler(CustomErrorHandler), mgin.WithLimitReachedHandler(CustomLimitReachedHandler))

	return middleware
}

func limiterRedisStore() limiter.Store {
	r := configs.GetRedisConfig()

	// Create a redis client.
	host := r.Host + ":" + strconv.Itoa(r.Port)

	// Create a redis option.
	option := &libredis.Options{
		Addr:     host,
		Password: r.Password,
		DB:       1,
	}

	// Create a redis client.
	client := libredis.NewClient(option)

	// Create a store with the redis client.
	store, err := sredis.NewStoreWithOptions(client, limiter.StoreOptions{
		Prefix:   "RateLimiter",
		MaxRetry: 3,
	})

	if err != nil {
		return nil
	}

	return store
}

func limiterInmemoryStore() limiter.Store {
	store := memory.NewStore()

	return store
}
