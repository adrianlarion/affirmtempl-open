package mysql

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/adrianlarion/affirmtempl-open/internal/model"
	"github.com/adrianlarion/affirmtempl-open/internal/model/csv"
	"github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/patrickmn/go-cache"
)

// global
var defaultUser *model.User

type UserModel struct {
	DB    *sql.DB
	Cache *cache.Cache
}

func (u *UserModel) CurrentUser(c echo.Context) (*model.User, error) {

	defaultUser, err := u.DefaultUser()
	if err != nil {
		return nil, err
	}

	if model.UserIsAuthenticated(c) {

		email, err := model.ReadKeyFromMainSessKey(c, model.SESS_UUID_KEY)
		if os.Getenv(model.PROD_OS_ENV_KEY) == "no" {
			email = model.TEST_USER_EMAIL
		}

		if err != nil {
			return nil, err
		}
		user, err := u.GetByEmail(email)
		if err != nil {
			return nil, err
		}
		return user, nil
	}

	return defaultUser, nil
}

func (u *UserModel) _testUser() (*model.User, error) {
	user, err := u.GetByEmail(model.TEST_USER_EMAIL)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return user, nil
}

// returns in memory user
func (u *UserModel) DefaultUser() (*model.User, error) {
	cards, err := csv.DefaultCardsFromCsv()

	if err != nil {
		return nil, err
	}

	if defaultUser == nil {
		defaultUser = &model.User{
			ID:       -1,
			Cards:    cards,
			Settings: &model.Settings{true, 25000, true},
			Name:     "default",
			Email:    "default@default.com",
			Auth:     &model.Auth{},
			Meta:     &model.Meta{},
			Sub:      "",
		}
	}
	return defaultUser, nil

}

func (u *UserModel) GetByEmail(email string) (*model.User, error) {

	//check if user is present in cache
	res, found := u.Cache.Get(email)

	//user is in cache, return it
	if found {
		// fmt.Println("User with email found in cache ",email)
		var userFromCache = res.(*model.User)
		return userFromCache, nil
	}

	//not in cache, do the usual stuff
	// fmt.Println("Accessing database for user with email ",email)

	user := &model.User{}
	var cardsString string
	var settingsString string
	var settings *model.Settings = &model.Settings{}
	var cards []*model.AffirmCard

	var authString string
	var metaString string
	var auth *model.Auth = &model.Auth{}
	var meta *model.Meta = &model.Meta{}
	// var cards []interface{}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// Defer cancel to make sure that we cancel the context before the Get() method returns
	defer cancel()

	stmt := "SELECT id, name, email, cards, settings, auth, meta, sub, created FROM user WHERE email = ?"
	err := u.DB.QueryRowContext(ctx, stmt, email).Scan(&user.ID, &user.Name, &user.Email, &cardsString, &settingsString, &authString, &metaString, &user.Sub, &user.Created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrNoRecord
		} else {
			return nil, err
		}
	}

	err = json.Unmarshal([]byte(authString), auth)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(metaString), meta)
	if err != nil {
		return nil, err
	}

	user.Auth = auth
	user.Meta = meta

	err = json.Unmarshal([]byte(settingsString), settings)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(cardsString), &cards)
	if err != nil {
		return nil, err
	}

	user.Settings = settings
	user.Cards = cards

	//before returning, add it to cache
	// fmt.Println("Addin guser with email to cache ",email)
	u.Cache.Set(email, user, cache.DefaultExpiration)

	return user, nil
}

func (u *UserModel) Insert(name, email string, cards []*model.AffirmCard, settings *model.Settings, auth *model.Auth, meta *model.Meta, sub string) error {

	stmt := `INSERT INTO user (name,email,cards,settings,auth,meta,sub,created) VALUES(
		?,?,?,?,?,?,?,UTC_TIMESTAMP()
	)`

	jsonAuth, err := json.Marshal(auth)
	if err != nil {
		return err
	}
	jsonMeta, err := json.Marshal(meta)
	if err != nil {
		return err
	}

	jsonCards, err := json.Marshal(cards)
	if err != nil {
		return err
	}
	settingsJson, err := json.Marshal(settings)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// Defer cancel to make sure that we cancel the context before the Get() method returns
	defer cancel()

	_, err = u.DB.ExecContext(ctx, stmt, name, email, jsonCards, settingsJson, jsonAuth, jsonMeta, sub)
	if err != nil {
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
				return model.ErrDuplicateEmail
			}
		}
		return err
	}
	return nil

}

func (u *UserModel) Delete(email string) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// Defer cancel to make sure that we cancel the context before the Get() method returns
	defer cancel()

	stmt := "DELETE FROM user WHERE email = ?"
	_, err := u.DB.ExecContext(ctx, stmt, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.ErrNoRecord
		} else {
			return err
		}
	}
	return nil
}

func (u *UserModel) Update(name, email string, cards []*model.AffirmCard, settings *model.Settings, auth *model.Auth, meta *model.Meta, sub string) error {
	//remove old entry from cache
	// fmt.Println("removing user with email from cache ",email)
	u.Cache.Delete(email)

	// stmt :=` UPDATE user
	// SET name = $1, email = $2, cards = $3, settings = $4
	// WHERE email = $2
	// `
	stmt := "UPDATE user SET name = ?, email = ?, cards = ?, settings = ?, auth = ?, meta = ?, sub = ?  WHERE email = ?"

	jsonAuth, err := json.Marshal(auth)
	if err != nil {
		return err
	}
	jsonMeta, err := json.Marshal(meta)
	if err != nil {
		return err
	}

	jsonCards, err := json.Marshal(cards)
	if err != nil {
		return err
	}
	settingsJson, err := json.Marshal(settings)
	if err != nil {
		return err
	}

	args := []interface{}{
		name,
		email,
		jsonCards,
		settingsJson,
		jsonAuth,
		jsonMeta,
		sub,
		email,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// Defer cancel to make sure that we cancel the context before the Get() method returns
	defer cancel()

	result, err := u.DB.ExecContext(ctx, stmt, args...)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		// return model.ErrNoRecord
	}

	return nil

}
