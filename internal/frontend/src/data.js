import {log} from './logger.js'
import {showErrorModal} from './utils.js'
import {MAX_AFFIRMATIONS} from './const.js'

export async function affirmTextIsFav(affirmText){
    let arr = await getAffirmArr(0)
    for (let i=0; i<arr.length;i++){
        if (arr[i]["content"]==affirmText){
            return true
        }
    }
    return false
}

export async function getAffirmArr(cardId){
    return await _getCardData(cardId)
}

export async function getDefaultAffirmArr(cardId){
    // return [{"content":"def1"},{"content":"def2"}]
    return await _getDefaultCardData(cardId)
}

export function affirmArrNeedsBackendData(cardId){
    if (Object.keys(_cardDataCache).length  === 0){
        return true
    }
    return _cardDataCache[cardId].dirty
}

export function settingsNeedsBackendData(){
    return _settingsCache.dirty
}



export async function setAffirmArr(affirmArr,cardId){
    _cardDataCache[cardId].dirty=true
    log.info("getting fresh card data from server for card with id %s",cardId)
    let cardUrl = "/card-affirm-arr/"+cardId
    await $.ajax({
        url:cardUrl,
        type:"PUT",
        data:affirmArr,
        success: function(result){
        },
        error: function(jqXHR,textStatus,errorThrown){
            log.error({error:errorThrown})
            showErrorModal()
            return null
        }
    })

}

export async function setCardFav(fav,cardId){
    log.info({fav:fav},"putting fresh fav status for card data on server for card "+ cardId)
    
    let url = "/card-fav/"+cardId
    await $.ajax({
        url:url,
        type:"PUT",
        data:fav,
    })
}

export async function setAffirmFav(fav){
    //by setting fav for affirm, we need to get fresh data for card 0 (fav card) wmich contains our fav affirm arr
    _cardDataCache[0].dirty=true

    log.info({fav:fav},"putting fresh fav status for affirm data on server for card ")
    let url = "/affirm-fav/"
    await $.ajax({
        url:url,
        type:"PUT",
        data:fav,
        success:function(result){
            // location.reload();
        },
        error: function(jqXHR,textStatus,errorThrown){
            log.error({error:errorThrown})
            showErrorModal()
            return null
        }

    })
    
}

export async function setSettings(settings){
    log.info({settings:settings},"putting fresh setttings data on server()")
    let settingsUrl="/settings"
    _settingsCache.dirty=true
    await $.ajax({
        url:settingsUrl,
        type:"PUT",
        data:settings,
        success:function(result){
        },
        error: function(jqXHR,textStatus,errorThrown){
            log.error({error:errorThrown})
            showErrorModal()
            return null
        }

    })
}

export async function logoutUser(){
    let url="/oauth-logout"
        await $.ajax({
            url:url,
        })
}

export async function getUserAuthStatus(){

    let url="/user-auth-status"
        await $.ajax({
            url:url,
            data:null,
            success:function(result){
            },
            error: function(jqXHR,textStatus,errorThrown){
            }
        })
}

export async function getSettings(){
    if (_settingsCache.dirty){
        log.info("getting fresh setttings data from server()")
        let settingsUrl="/settings"
        await $.ajax({
            url:settingsUrl,
            data:null,
            success:function(result){
                _settingsCache.data=result
                _settingsCache.dirty=false
                log.info({result:result},"settings data from server is:")
            },
            error: function(jqXHR,textStatus,errorThrown){
                log.error({error:errorThrown})
                showErrorModal()
                return null
            }

        })
    }
    return _settingsCache.data
}

let _cardDataCache = {}
let _defaultCardDataCache = {}
let _settingsCache = {data:null,dirty:true}

 
function _genCacheData(){
    //card affirm arr from user
    if (Object.keys(_cardDataCache).length  === 0){
        log.info("generating card data cache")
        for (let i=0; i<MAX_AFFIRMATIONS;i++){
            _cardDataCache[i.toString()]={data:null,dirty:true}
        }
        log.info({_cardDataCache:_cardDataCache},"generated  data cache")
    }
    //default card affirm arr set by developer
    if (Object.keys(_defaultCardDataCache).length  === 0){
        log.info("generating default card data cache")
        for (let i=0; i<MAX_AFFIRMATIONS;i++){
            _defaultCardDataCache[i.toString()]={data:null,dirty:true}
        }
        log.info({_cardDataCache:_defaultCardDataCache},"generated  default data cache")
    }
}

async function _getDefaultCardData(cardId){
    _genCacheData()

    if (_defaultCardDataCache[cardId].dirty){
        log.info("getting fresh card data from server for card with id %s",cardId)
        let cardUrl = "/card-default-affirm-arr/"+cardId
        await $.ajax({
            url:cardUrl,
            data:null,
            success: function(result){
                _defaultCardDataCache[cardId].dirty=false
                _defaultCardDataCache[cardId].data=result
                log.info({result:result},"card data from server is:")
            },
            error: function(jqXHR,textStatus,errorThrown){
                log.error({error:errorThrown})
                showErrorModal()
                return null
            }
        })
    }
    return _defaultCardDataCache[cardId].data
}

async function _getCardData(cardId){
    _genCacheData()

    if (_cardDataCache[cardId].dirty){
        log.info("getting fresh card data from server for card with id %s",cardId)
        let cardUrl = "/card-affirm-arr/"+cardId
        await $.ajax({
            url:cardUrl,
            data:null,
            success: function(result){
                _cardDataCache[cardId].dirty=false
                _cardDataCache[cardId].data=result
                log.info({result:result},"card data from server is:")
            },
            error: function(jqXHR,textStatus,errorThrown){
                log.error({error:errorThrown})
                showErrorModal()
                return null
            }
        })
    }
    return _cardDataCache[cardId].data
}
