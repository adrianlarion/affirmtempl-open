package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/adrianlarion/affirmtempl-open/internal/model"
	"github.com/adrianlarion/affirmtempl-open/internal/model/csv"
	"github.com/adrianlarion/affirmtempl-open/internal/view/page"

	"github.com/labstack/echo/v4"
	"google.golang.org/api/idtoken"
)

func (app *application) home(c echo.Context) error {

	txn := app.grafana_app.StartTransaction("home - GET - /")
	defer txn.End()

	app.log.Debug("home()")

	allCards, err := app.userCards(c)

	if err != nil {
		app.log.Error(err)
		app.writeToErrLog(err, c)
		return echo.NewHTTPError(http.StatusInternalServerError, model.ErrGenericServerError)
	}
	twoCardsArr := model.SplitAffirmCards(allCards)

	if err != nil {
		app.log.Error(err)
		app.writeToErrLog(err, c)
		return echo.NewHTTPError(http.StatusInternalServerError, model.ErrGenericServerError)
	}

	isAuth := model.UserIsAuthenticated(c)
	userName, err := app.userName(c)
	if err != nil {
		app.log.Error(err)
		app.writeToErrLog(err, c)
		return echo.NewHTTPError(http.StatusInternalServerError, model.ErrGenericServerError)
	}

	return render(c, http.StatusOK, page.Home(twoCardsArr, isAuth, userName))
}

func (app *application) tos(c echo.Context) error {
	return render(c, http.StatusOK, page.Tos())
}

func (app *application) privacyPolicy(c echo.Context) error {
	return render(c, http.StatusOK, page.PrivacyPolicy())
}

func (app *application) userAuthStatus(c echo.Context) error {

	txn := app.grafana_app.StartTransaction("userAuthStatus - GET - /user-auth-status")
	defer txn.End()

	if model.UserIsAuthenticated(c) {
		return c.NoContent(http.StatusOK)
	} else {
		return c.NoContent(http.StatusUnauthorized)
	}
}

func (app *application) googleLogout(c echo.Context) error {

	txn := app.grafana_app.StartTransaction("googleLogout - GET - /oauth-logout")
	defer txn.End()

	app.log.Debug("googleLogout()")
	err := model.EmptyKeyFromMainSess(c, model.SESS_UUID_KEY)
	if err != nil {
		app.log.Error(err)
		app.writeToErrLog(err, c)
		return echo.NewHTTPError(http.StatusInternalServerError, model.ErrGenericServerError)
	}

	// app.log.Debug("redirecting user to home page()")
	// return c.Redirect(http.StatusMovedPermanently,"/")
	return nil

}

func (app *application) googleLogin(c echo.Context) error {

	txn := app.grafana_app.StartTransaction("googleLogin - POST - /oauth")
	defer txn.End()

	token := c.FormValue("credential")
	googleClientId := model.GOOGLE_CLIENT_ID_OAUTH

	payload, err := idtoken.Validate(c.Request().Context(), token, googleClientId)
	if err != nil {
		app.log.Error(err)
		app.writeToErrLog(err, c)
		return echo.NewHTTPError(http.StatusInternalServerError, model.ErrGenericServerError)
	}

	fmt.Println(payload.Claims)
	claims := payload.Claims

	givenName := ""
	email := ""
	sub := ""

	_, ok := claims["given_name"]
	if ok {
		givenName = claims["given_name"].(string)
	}

	_, ok = claims["email"]
	if ok {
		email = claims["email"].(string)
	}

	_, ok = claims["sub"]
	if ok {
		sub = claims["sub"].(string)
	}

	//if empty email return error
	if len(email) <= 0 {
		app.log.Error(errors.New("empty email"))
		app.writeToErrLog(err, c)
		return echo.NewHTTPError(http.StatusInternalServerError, model.ErrGenericServerError)
	}

	//check if user exists
	_, err = app.user.GetByEmail(email)

	if err != nil {
		// use doent exist
		if errors.Is(err, model.ErrNoRecord) {
			//create new user
			allCards, err := csv.DefaultCardsFromCsv()

			err = app.user.Insert(givenName,
				email,
				allCards,
				&model.Settings{true, 25000, true},
				&model.Auth{Sub: sub},
				&model.Meta{GoogleLogin: true},
				sub,
			)
			//create a new user got us an error
			if err != nil {
				app.log.Error(err)
				app.writeToErrLog(err, c)
				return echo.NewHTTPError(http.StatusInternalServerError, model.ErrGenericServerError)
			}

		} else {
			//something else went wrong
			app.log.Error(err)
			app.writeToErrLog(err, c)
			return echo.NewHTTPError(http.StatusInternalServerError, model.ErrGenericServerError)

		}
	}

	//create session
	model.CreateAppendMainSessionWithKeyVal(c, model.SESS_UUID_KEY, email)

	//redirect to home page
	return c.Redirect(http.StatusMovedPermanently, "/")

}

func (app *application) getDefaultAffirmArr(c echo.Context) error {
	txn := app.grafana_app.StartTransaction("getDefaultAffirmArr - GET - /card-default-affirm-arr")
	defer txn.End()

	idParam := c.Param("id")
	app.log.Debug("getDefaultAffirmArr() id ", idParam)

	id, err := strconv.ParseInt(idParam, 10, 64)

	if err != nil {
		app.log.Error(err)
		app.writeToErrLog(err, c)
		return echo.NewHTTPError(http.StatusInternalServerError, model.ErrGenericServerError)

	}

	allCards, err := app.userCards(c)
	if err != nil {
		app.log.Error(err)
		app.writeToErrLog(err, c)
		return echo.NewHTTPError(http.StatusInternalServerError, model.ErrGenericServerError)
	}
	card := cardFromCardsById(allCards, id)

	return c.JSON(http.StatusOK, card.DefaultAffirmations)

}

func (app *application) getAffirmArr(c echo.Context) error {

	txn := app.grafana_app.StartTransaction("getAffirmArr - GET - /card-affirm-arr")
	defer txn.End()

	idParam := c.Param("id")
	app.log.Debug("getAffirmArr() id ", idParam)

	id, err := strconv.ParseInt(idParam, 10, 64)

	if err != nil {
		app.log.Error(err)
		app.writeToErrLog(err, c)
		return echo.NewHTTPError(http.StatusInternalServerError, model.ErrGenericServerError)

	}

	allCards, err := app.userCards(c)
	if err != nil {
		app.log.Error(err)
		app.writeToErrLog(err, c)
		return echo.NewHTTPError(http.StatusInternalServerError, model.ErrGenericServerError)
	}
	card := cardFromCardsById(allCards, id)

	return c.JSON(http.StatusOK, card.Affirmations)

}

func (app *application) userCards(c echo.Context) ([]*model.AffirmCard, error) {
	//returns current user's cards
	user, err := app.user.CurrentUser(c)
	if err != nil {
		return nil, err
	}
	return user.Cards, nil
}

func (app *application) userName(c echo.Context) (string, error) {
	user, err := app.user.CurrentUser(c)
	if err != nil {
		return "", err
	}
	return user.Name, nil
}

func (app *application) userSingleCard(cardId int64) (*model.AffirmCard, error) {

	return nil, nil
}

func cardFromCardsById(cards []*model.AffirmCard, id int64) *model.AffirmCard {

	for _, v := range cards {
		if v.ID == id {
			return v
		}
	}
	return nil
}

func (app *application) putFavStatusAffirm(c echo.Context) error {
	txn := app.grafana_app.StartTransaction("putFavStatusAffirm - PUT - /affirm-fav")
	defer txn.End()

	json_map := make(map[string]interface{})
	err := json.NewDecoder(c.Request().Body).Decode(&json_map)

	if err != nil {
		app.log.Error(err)
		app.writeToErrLog(err, c)
		return echo.NewHTTPError(http.StatusInternalServerError, model.ErrGenericServerError)
	}

	allCards, err := app.userCards(c)
	if err != nil {
		app.log.Error(err)
		app.writeToErrLog(err, c)
		return echo.NewHTTPError(http.StatusInternalServerError, model.ErrGenericServerError)
	}
	card := cardFromCardsById(allCards, 0)

	err = card.UpdateFavCardAffirmArr(json_map)

	if err != nil {
		app.log.Error(err)
		app.writeToErrLog(err, c)
		return echo.NewHTTPError(http.StatusInternalServerError, model.ErrGenericServerError)
	}

	////add cards to user, write to db
	user, err := app.user.CurrentUser(c)
	if err != nil {
		app.log.Error(err)
		app.writeToErrLog(err, c)
		return echo.NewHTTPError(http.StatusInternalServerError, model.ErrGenericServerError)
	}
	replaceCardInUser(card, user)

	err = app.user.Update(user.Name, user.Email, user.Cards, user.Settings, user.Auth, user.Meta, user.Sub)
	if err != nil {
		app.log.Error(err)
		app.writeToErrLog(err, c)
		return echo.NewHTTPError(http.StatusInternalServerError, model.ErrGenericServerError)
	}
	return nil
}

func (app *application) putFavStatusCard(c echo.Context) error {

	txn := app.grafana_app.StartTransaction("putFavStatusCard - PUT - /card-fav")
	defer txn.End()

	idParam := c.Param("id")
	app.log.Debug("putFavStatusCard() id ", idParam)

	//get json
	json_map := make(map[string]interface{})
	err := json.NewDecoder(c.Request().Body).Decode(&json_map)

	if err != nil {
		app.log.Error(err)
		app.writeToErrLog(err, c)
		return echo.NewHTTPError(http.StatusInternalServerError, model.ErrGenericServerError)
	}

	//get single card
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		app.log.Error(err)
		app.writeToErrLog(err, c)
		return echo.NewHTTPError(http.StatusInternalServerError, model.ErrGenericServerError)

	}
	allCards, err := app.userCards(c)
	if err != nil {
		app.log.Error(err)
		app.writeToErrLog(err, c)
		return echo.NewHTTPError(http.StatusInternalServerError, model.ErrGenericServerError)
	}
	card := cardFromCardsById(allCards, id)

	err = card.FillFavStatusFromJson(json_map)

	if err != nil {
		app.log.Error(err)
		app.writeToErrLog(err, c)
		return echo.NewHTTPError(http.StatusInternalServerError, model.ErrGenericServerError)
	}

	////add cards to user, write to db
	user, err := app.user.CurrentUser(c)
	if err != nil {
		app.log.Error(err)
		app.writeToErrLog(err, c)
		return echo.NewHTTPError(http.StatusInternalServerError, model.ErrGenericServerError)
	}
	replaceCardInUser(card, user)

	err = app.user.Update(user.Name, user.Email, user.Cards, user.Settings, user.Auth, user.Meta, user.Sub)
	if err != nil {
		app.log.Error(err)
		app.writeToErrLog(err, c)
		return echo.NewHTTPError(http.StatusInternalServerError, model.ErrGenericServerError)
	}

	return nil
}

func (app *application) putAffirmArr(c echo.Context) error {
	txn := app.grafana_app.StartTransaction("putAffirmArr - PUT - /card-affirm-arr")
	defer txn.End()

	app.log.Debug("putAffirmArr()")

	idParam := c.Param("id")
	app.log.Debug("putAffirmArr() id ", idParam)

	json_map := make(map[string][]string)
	err := json.NewDecoder(c.Request().Body).Decode(&json_map)
	// err := json.Unmarshal([]byte(c.Request().Body),json_map)

	if err != nil {
		app.log.Error(err)
		app.writeToErrLog(err, c)
		return echo.NewHTTPError(http.StatusInternalServerError, model.ErrGenericServerError)
	}

	//get single card
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		app.log.Error(err)
		app.writeToErrLog(err, c)
		return echo.NewHTTPError(http.StatusInternalServerError, model.ErrGenericServerError)

	}
	allCards, err := app.userCards(c)
	if err != nil {
		app.log.Error(err)
		app.writeToErrLog(err, c)
		return echo.NewHTTPError(http.StatusInternalServerError, model.ErrGenericServerError)
	}
	card := cardFromCardsById(allCards, id)

	err = card.FillAffirmArrValsFromJson(json_map["affirmations"])

	if err != nil {
		app.log.Error(err)
		app.writeToErrLog(err, c)
		return echo.NewHTTPError(http.StatusInternalServerError, model.ErrGenericServerError)
	}

	////add cards to user, write to db
	user, err := app.user.CurrentUser(c)
	if err != nil {
		app.log.Error(err)
		app.writeToErrLog(err, c)
		return echo.NewHTTPError(http.StatusInternalServerError, model.ErrGenericServerError)
	}
	replaceCardInUser(card, user)

	err = app.user.Update(user.Name, user.Email, user.Cards, user.Settings, user.Auth, user.Meta, user.Sub)
	if err != nil {
		app.log.Error(err)
		app.writeToErrLog(err, c)
		return echo.NewHTTPError(http.StatusInternalServerError, model.ErrGenericServerError)
	}

	return nil
}

func replaceCardInUser(newCard *model.AffirmCard, user *model.User) {
	for i := 0; i < len(user.Cards); i++ {
		if user.Cards[i].ID == newCard.ID {
			user.Cards[i] = newCard
			break
		}
	}

}

func (app *application) putSettings(c echo.Context) error {

	txn := app.grafana_app.StartTransaction("putSettings - PUT - /settings")
	defer txn.End()

	app.log.Debug("putSettings()")

	json_map := make(map[string]interface{})
	err := json.NewDecoder(c.Request().Body).Decode(&json_map)

	if err != nil {
		app.log.Error(err)
		app.writeToErrLog(err, c)
		return echo.NewHTTPError(http.StatusInternalServerError, model.ErrGenericServerError)
	}

	s := model.Settings{}
	err = s.FillValuesFromJson(json_map)
	if err != nil {
		app.log.Error(err)
		app.writeToErrLog(err, c)
		return echo.NewHTTPError(http.StatusInternalServerError, model.ErrGenericServerError)
	}
	//add settings to user, write to db
	user, err := app.user.CurrentUser(c)
	if err != nil {
		app.log.Error(err)
		app.writeToErrLog(err, c)
		return echo.NewHTTPError(http.StatusInternalServerError, model.ErrGenericServerError)
	}
	user.Settings = &s

	err = app.user.Update(user.Name, user.Email, user.Cards, user.Settings, user.Auth, user.Meta, user.Sub)
	if err != nil {
		app.log.Error(err)
		app.writeToErrLog(err, c)
		return echo.NewHTTPError(http.StatusInternalServerError, model.ErrGenericServerError)
	}

	return nil
}

func (app *application) getSettings(c echo.Context) error {

	txn := app.grafana_app.StartTransaction("getSettings - GET - /settings")
	defer txn.End()

	settings, err := app.userSettings(c)
	if err != nil {
		app.log.Error(err)
		app.writeToErrLog(err, c)
		return echo.NewHTTPError(http.StatusInternalServerError, model.ErrGenericServerError)
	}
	return c.JSON(http.StatusOK, settings)
}

func (app *application) userSettings(c echo.Context) (*model.Settings, error) {
	user, err := app.user.CurrentUser(c)
	if err != nil {
		return nil, err
	}
	return user.Settings, nil
}
