package engine

import (
	"database/sql"

	helmet "github.com/danielkov/gin-helmet"
	"github.com/dusansimic/yaas/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

type options struct {
	DB          *sql.DB
	WAPI        api.WriteAPIBlocking
	CookieStore cookie.Store
	ServerStore memstore.Store
	Params      map[string]interface{}
}

// Option sets an option
type Option func(*options)

// New returns a new engine
func New(opts ...Option) *gin.Engine {
	o := options{}
	o.Params = make(map[string]interface{})
	for _, opt := range opts {
		opt(&o)
	}
	engine := gin.Default()

	cfg := cors.DefaultConfig()
	cfg.AllowAllOrigins = true
	cfg.AddAllowHeaders("Authorization", "authorization")
	cfg.AllowCredentials = true
	engine.Use(cors.New(cfg))
	engine.Use(helmet.Default())

	engine.Use(sessions.SessionsMap(map[string]sessions.Store{
		"cookie": o.CookieStore,
		"server": o.ServerStore,
	}))

	engine.Use(handlers.HandleErrors())

	engine.POST("/event", handlers.ParseEvent(o.DB), handlers.CreateRecord(o.DB), handlers.AddRecord(o.DB))
	engine.GET("/error", handlers.Error())

	auth := engine.Group("/auth")
	{
		auth.GET("/verify", handlers.AuthenticateJWT(), handlers.Ok)
		auth.POST("/login", handlers.Login(o.DB, o.Params["frontBase"].(string)))
		auth.POST("/logout", handlers.AuthenticateJWT(), handlers.Logout(o.DB))
		auth.POST("/register", handlers.Register(o.DB))
		auth.GET("/available", handlers.AvailableUser(o.DB))
	}

	user := engine.Group("/user")
	user.Use(handlers.AuthenticateJWT())
	{
		domain := user.Group("/domain")
		{
			// Gets one or multiple domains and their information
			domain.GET("", handlers.GetDomain(o.DB))
			// Adds a new domain to the database
			domain.POST("", handlers.AddDomain(o.DB))
			// Updates a domain in the database
			domain.PUT("", handlers.UpdateDomain(o.DB))
			// Deletes a domain and all it's events from the database
			domain.DELETE("")

			// domain.OPTIONS("", handlers.Preflight)
		}

		stats := user.Group("/stats")
		{
			stats.GET("/:id", handlers.Stats(o.DB))
		}
	}

	return engine
}

// WithDB sets a new database connection as an option
func WithDB(db *sql.DB) Option {
	return func(o *options) {
		o.DB = db
	}
}

func WithWriteAPI(w api.WriteAPIBlocking) Option {
	return func(o *options) {
		o.WAPI = w
	}
}

// WithCookieStore sets a cookie store as an option
func WithCookieStore(store sessions.Store) Option {
	return func(o *options) {
		o.CookieStore = store
	}
}

// WithServerStore sets a new server store as an option
func WithServerStore(store sessions.Store) Option {
	return func(o *options) {
		o.ServerStore = store
	}
}

// WithParam sets a new param as an option
func WithParam(k string, v interface{}) Option {
	return func(o *options) {
		o.Params[k] = v
	}
}
