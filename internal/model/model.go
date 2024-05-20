package model

import (
	"errors"
	"time"
)

var (
	ErrNoRecord           = errors.New("models: no matching record found")
	ErrDuplicateEmail     = errors.New("models: user with email already exists")
	ErrGenericServerError = errors.New("The server has encountered an internal error")
)

const DEFAULT_USER_EMAIL = "test@test.com"
const DEFAULT_USER_NAME = "default"
const GOOGLE_CLIENT_ID_OAUTH = ""
const SESSION_NAME = "easyaffirm_session"
const SESS_UUID_KEY = "sub"
const TEST_USER_EMAIL = "test@test.com"
const PROD_OS_ENV_KEY = "AFFIRMTEMPL_PROD"

type User struct {
	ID       int64         `json:"id"`
	Cards    []*AffirmCard `json:"cards"`
	Settings *Settings     `json:"settings"`
	Name     string        `json:"name"`
	Email    string        `json:"email"`
	Created  time.Time     `json:"-"`
	Active   bool          `json:"active"`
	Auth     *Auth         `json:"auth"`
	Meta     *Meta         `json:"meta"`
	Sub      string        `json:"-"`
}

type Auth struct {
	//dummy mostly
	Sub string `json:"sub"`
}

type Meta struct {
	//dummy mostly
	GoogleLogin bool `json:"google_login"`
}

type Settings struct {
	Autoplay         bool  `json:"autoplay"`
	AutoplayDuration int64 `json:"autoplayDuration"`
	RandomAffirm     bool  `json:"randomAffirm"`
}

type AffirmCard struct {
	Fav                 bool          `json:"fav"`
	Title               string        `json:"title"`
	ImgPath             string        `json:"imgPath"`
	ID                  int64         `json:"id"`
	Affirmations        []AffirmEntry `json:"affirmations"`
	DefaultAffirmations []AffirmEntry `json:"defaultAffirmations"`
}

type AffirmEntry struct {
	Content string `json:"content"`
}

type TwoCards struct {
	One *AffirmCard
	Two *AffirmCard
}
