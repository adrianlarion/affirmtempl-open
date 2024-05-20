package main

import (
	"context"
	"errors"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"runtime/pprof"
	"time"

	"github.com/adrianlarion/affirmtempl-open/internal/model"
	"github.com/adrianlarion/affirmtempl-open/internal/model/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/patrickmn/go-cache"
	"github.com/rs/zerolog"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"

	"github.com/newrelic/go-agent/v3/newrelic"
)

type application struct {
	user        *mysql.UserModel
	log         echo.Logger
	errorLog    *zerolog.Logger
	grafana_app *newrelic.Application
}

func main() {

	esrv := echo.New()
	//from here turn on/off logging level
	esrv.Logger.SetLevel(log.DEBUG)

	//db
	dsn := flag.String("dsn", "affirmtempl:pass@/affirmtempl?parseTime=true", "mysql data source name")
	port := flag.String("port", "4000", "port for app")
	//profile
	var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
	//parse
	flag.Parse()
	sessionSecret := os.Getenv("AFFIRMTEMPL_SESSION_SECRET")
	if len(sessionSecret) <= 0 {
		esrv.Logger.Fatal(errors.New("no session secret key provided or available as os env"))
	}

	esrv.Logger.Info("Trying to open db")
	db, err := openDB(*dsn)
	if err != nil {
		esrv.Logger.Fatal(err)
	}
	defer db.Close()

	//user cache
	userCache := cache.New(10*1*time.Minute, 10*time.Minute)

	//profiling
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	//profiling

	//error logger
	errorLog, err := createErrLogger()
	if err != nil {
		esrv.Logger.Fatal(err)
	}

	//grafana
	grafana_app, err := newrelic.NewApplication(
		newrelic.ConfigAppName("affirmtempl"),
		newrelic.ConfigLicense(""),
		newrelic.ConfigAppLogForwardingEnabled(true),
	)

	//app
	app := &application{
		user:        &mysql.UserModel{DB: db, Cache: userCache},
		log:         esrv.Logger,
		errorLog:    errorLog,
		grafana_app: grafana_app,
	}

	//--------------------------
	//middleware provided by echo
	//--------------------------
	//panic recover
	esrv.Use(middleware.Recover())
	//body limit
	esrv.Use(middleware.BodyLimit("35K"))
	//secure header
	esrv.Use(middleware.Secure())
	//timeout
	esrv.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: 5 * time.Second,
	}))
	//logger
	esrv.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `{"time":"${time_rfc3339_nano}","remote_ip":"${remote_ip}",` +
			`"host":"${host}","method":"${method}","uri":"${uri}",` +
			`"status":${status},"error":"${error}","latency_human":"${latency_human}"` +
			`` + "\n",
		CustomTimeFormat: "2006-01-02 15:04:05.00000",
	}))
	//sessions
	esrv.Use(session.Middleware(sessions.NewCookieStore([]byte(sessionSecret))))

	//DEV stuff, not prod
	if os.Getenv(model.PROD_OS_ENV_KEY) == "no" {
		app._createTestUser()
	}

	//routes
	esrv.Static("/static", "static")
	esrv.GET("/", app.home)
	esrv.GET("/card-default-affirm-arr/:id", app.getDefaultAffirmArr)
	esrv.GET("/card-affirm-arr/:id", app.getAffirmArr)
	esrv.PUT("/card-affirm-arr/:id", app.putAffirmArr, checkAuth)
	esrv.PUT("/card-fav/:id", app.putFavStatusCard, checkAuth)
	esrv.PUT("/affirm-fav/", app.putFavStatusAffirm, checkAuth)
	esrv.GET("/settings", app.getSettings)
	esrv.PUT("/settings", app.putSettings, checkAuth)
	esrv.POST("/oauth", app.googleLogin)
	esrv.GET("/oauth-logout", app.googleLogout)
	esrv.GET("/user-auth-status", app.userAuthStatus)

	esrv.GET("/terms-of-use", app.tos)
	esrv.GET("/privacy-policy", app.privacyPolicy)

	//graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		if err := esrv.Start(":" + *port); err != nil && err != http.ErrServerClosed {
			esrv.Logger.Fatal("shutting down server")
		}

	}()

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := esrv.Shutdown(ctx); err != nil {
		esrv.Logger.Fatal(err)
	}

}
