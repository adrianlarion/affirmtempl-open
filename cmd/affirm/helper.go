package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/a-h/templ"
	"github.com/adrianlarion/affirmtempl-open/internal/model"
	"github.com/adrianlarion/affirmtempl-open/internal/model/csv"
	"github.com/labstack/echo/v4"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rs/zerolog"
)

func render(ctx echo.Context, statusCode int, t templ.Component) error {
	ctx.Response().Writer.WriteHeader(statusCode)
	ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	return t.Render(ctx.Request().Context(), ctx.Response().Writer)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func createErrLogger() (*zerolog.Logger, error) {
	p := filepath.Join(".", "log")
	err := os.MkdirAll(p, os.ModePerm)

	errFilePath := "./log/err.txt"
	file, err := os.OpenFile(errFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	if err != nil {
		return nil, err
	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	logger := zerolog.New(file)

	return &logger, nil
}

// func (app *application)createDefaultUser(){

// 	app.user.Delete(model.DEFAULT_USER_EMAIL)

// 	allCards, err := csv.DefaultCardsFromCsv()
// 	err = app.user.Insert(model.DEFAULT_USER_NAME,model.DEFAULT_USER_EMAIL, allCards,&model.Settings{},&model.Auth{},&model.Meta{})
// 	if err != nil{
// 		app.log.Error(err)
// 	}
// }

func (app *application) writeToErrLog(err error, c echo.Context) {

	r := c.Request()
	app.errorLog.Debug().
		Str("remote addr", r.RemoteAddr).
		Str("proto", r.Proto).
		Str("method", r.Method).
		Str("uri", r.URL.RequestURI()).
		Msg(err.Error())

	//grafana error log
	now := time.Now()
	msg := fmt.Sprintf("remote addr: %s, proto: %s, method: %s, uri: %s, err: %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI(), err.Error())
	app.grafana_app.RecordLog(newrelic.LogData{Timestamp: now.Unix(), Severity: "error", Message: msg})
}

func (app *application) _createTestUser() error {

	testUser, _ := app.user.GetByEmail(model.TEST_USER_EMAIL)
	if testUser != nil {
		return nil
	}
	app.log.Info("creating test user")

	allCards, err := csv.DefaultCardsFromCsv()

	err = app.user.Insert("Test",
		model.TEST_USER_EMAIL,
		allCards,
		&model.Settings{true, 25000, true},
		&model.Auth{Sub: "testsub"},
		&model.Meta{GoogleLogin: true},
		"testsub",
	)
	//create a new user got us an error
	if err != nil {
		app.log.Error(err)
		return err
	}
	return nil

}
