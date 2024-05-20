package model

import (
	"math"
	"os"
	"strconv"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func (c *AffirmCard) FillAffirmArrValsFromJson(affirmArr []string) error{
	affirmations:=[]AffirmEntry{}
	for _, affirm := range affirmArr{
		affirmEntry := AffirmEntry{Content:affirm}
		affirmations=append(affirmations,affirmEntry)
	}
	c.Affirmations=affirmations


	return nil
}

func (c *AffirmCard) UpdateFavCardAffirmArr(json_map map[string]interface{})error{
	if json_map["fav"].(bool){
		c.Affirmations = appendAffirmToAffirmEntriesIfNotExists(json_map["affirm"].(string),c.Affirmations)
	}else{
		c.Affirmations=removeAffirmEntryIfExists(json_map["affirm"].(string),c.Affirmations)
	}
	return nil
}

func (c *AffirmCard) FillFavStatusFromJson(json_map map[string]interface{}) error{
	c.Fav=json_map["fav"].(bool)
	return nil
}

func (s *Settings) FillValuesFromJson(json_map map[string]interface{}) error{ 
	
	autoplay:=json_map["autoplay"].(bool)
	randomAffirm:=json_map["randomAffirm"].(bool)
	autoplayDuration, err := strconv.ParseInt(json_map["autoplayDuration"].(string),10,64)
	if err != nil{
		return err
	}
	s.Autoplay=autoplay
	s.RandomAffirm=randomAffirm
	s.AutoplayDuration=autoplayDuration
	return nil
}




func sortAffirmCardsByFav(affirmCards []*AffirmCard)[]*AffirmCard{

	// sortedCards := make([]*AffirmCard,len(affirmCards))
	sortedCards := make([]*AffirmCard,len(affirmCards))
	copy(sortedCards,affirmCards)


	favPtr:=-1
	for i:=0; i<len(sortedCards);i++{
		if sortedCards[i].Fav{
			favPtr+=1
			sortedCards[favPtr], sortedCards[i]=sortedCards[i],sortedCards[favPtr]
		}
	}
	return sortedCards
}

func SplitAffirmCards(affirmCards []*AffirmCard) []TwoCards {
	arrLen:=int64(math.Ceil(float64(len(affirmCards))/2))
	var arr = make([]TwoCards,arrLen)

	//sort by fav so that they are first
	affirmCards = sortAffirmCardsByFav(affirmCards)

	soFarIndexCount:=0
	arrIndex:=0
	for _,v := range affirmCards{

		if soFarIndexCount==0{
			arr[arrIndex].One=v
		}else if soFarIndexCount==1{
			arr[arrIndex].Two=v
		}

		soFarIndexCount++
		if soFarIndexCount >= 2{
			arrIndex++
			soFarIndexCount=0
		}

	}
	

	return arr
}


func appendAffirmToAffirmEntriesIfNotExists(affirm string, affirmations []AffirmEntry)[]AffirmEntry{
	exists:=false
	for _, v := range affirmations{
		if v.Content==affirm{
			exists=true
			break
		}
	}
	if !exists{
		affirmations=append(affirmations,AffirmEntry{Content:affirm})
	}
	return affirmations

}

func removeAffirmEntryIfExists(affirm string, affirmations []AffirmEntry) []AffirmEntry{
	existsIndex:=-1
	for i:=0; i<len(affirmations);i++{
		if affirmations[i].Content==affirm{
			existsIndex=i
			break
		}
	}
	if existsIndex>=0{
		affirmations=append(affirmations[:existsIndex],affirmations[existsIndex+1:]...)
	}
	return affirmations
}
//creates or appends to a session
func CreateAppendMainSessionWithKeyVal(c echo.Context, k, v string) {
	  sess, _ := session.Get(SESSION_NAME, c)
	  sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 30,
		// MaxAge:   60 * 20,
		HttpOnly: true,
	  }
	  sess.Values[k] = v
	  sess.Save(c.Request(), c.Response())
}

func ReadKeyFromMainSessKey(c echo.Context, k string) (string, error){
	  sess, err := session.Get(SESSION_NAME, c)
	  v:=""
	  if err!=nil || sess.Values[k]==nil{
		  return v,err
	  }
	  v=sess.Values[k].(string)
	  return v, nil
}

func EmptyKeyFromMainSess(c echo.Context, k string) error{
	  sess, err := session.Get(SESSION_NAME, c)
	  if err != nil{
		  return err
	  }
	  sess.Values[k] = nil
	  sess.Save(c.Request(), c.Response())
	  return nil
}

func UserIsAuthenticated(c echo.Context) bool{
	if os.Getenv(PROD_OS_ENV_KEY)=="no"{
		return true
	}
	v, err:=ReadKeyFromMainSessKey(c,SESS_UUID_KEY)
	if err!=nil || len(v)<=0{
		return false
	}
	return true
}
