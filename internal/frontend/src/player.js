import {fxFlash, getCurrentAffirmTextFromPlayer} from './utils.js'
import {getUserAuthStatus,getAffirmArr, getSettings, affirmArrNeedsBackendData, settingsNeedsBackendData,setAffirmFav,affirmTextIsFav} from './data.js'
import {showLoader, hideLoader,getSelectedCardId,showLoginModal,showErrorModal} from './utils.js'
import {NO_AFFIRM_INFO_TEXT} from './const.js'
import {log} from './logger.js'

//PRIVATE
//----------------------------------------------------------------
//----------------------------------------------------------------



let _settings
let _affirmArr
let _curIndex
let _autoplayTimeout
let _autoplayDuration

let favSelector = "#app-player-btn-fav"
let favUnfav = 'favorite_border'
let favFav = 'favorite'

function _playerIsVisible(){
    let c=$("#app-home-player").is(":visible")
    return c
}


function _changeAffirmText(affirmText){
    $(".app-affirm-text").fadeOut(200,function(){
    $(this).text(affirmText).fadeIn(200);
    })
}

function _changeFavStatus(isFav){
    if (isFav){
        $(favSelector).html(favFav)
        $(favSelector).removeClass("text-type-dark-transp")
        $(favSelector).addClass("text-logo-main-transp")
    }else{
        $(favSelector).html(favUnfav)
        $(favSelector).removeClass("text-logo-main-transp")
        $(favSelector).addClass("text-type-dark-transp")
    }

}

function _incrementCurIndex(){
    let isResetting = false
    let nextIndex=_curIndex+1
    if (nextIndex>=_affirmArr.length) {
        nextIndex=0
        isResetting=true
    }
    _curIndex=nextIndex
    return isResetting
}

async function _getNextAffirmText(){
    let isResetting = _incrementCurIndex()
    if (isResetting === true){
    
        //on favorite card, when we get to the end of affirms, get fresh data from the server since it might have been changed
        let cardId = getSelectedCardId()
        if (cardId==="0"){

            let needBackendDataAffirmArr = affirmArrNeedsBackendData(cardId)

            //get data from server
            if (needBackendDataAffirmArr){
                showLoader()
            }
            _affirmArr = await getAffirmArr(cardId)

            if ( needBackendDataAffirmArr){
                hideLoader()
            }

        }
        fxFlash()
        _shuffleArrayIfNeeded()
    }

    if (_affirmArrIsEmpty()){
        return NO_AFFIRM_INFO_TEXT
    }
     
    if (!_affirmArr){
        return NO_AFFIRM_INFO_TEXT
    }

    if (_curIndex >= _affirmArr.length){
        return ""
    }
    return _affirmArr[_curIndex].content
}

async function _getAffirmFavStatus(){
    if (_affirmArrIsEmpty()){
        return false
    }
    return await affirmTextIsFav(_affirmArr[_curIndex]["content"])
}

 
async function _playNextAffirm(){
    if (_autoplayTimeout != null){
        window.clearTimeout(_autoplayTimeout)
    }
    _changeAffirmText(await _getNextAffirmText())
    _changeFavStatus(await _getAffirmFavStatus())
    if (_settings["autoplay"]===true){
        _setTimeoutAutoplay()
    }
}

//autoplay
function _setTimeoutAutoplay(){
    _autoplayTimeout=window.setTimeout(function(){
            _playNextAffirm()
    },_autoplayDuration)

}
function _stopAutoplayTimeout(){
    window.clearTimeout(_autoplayTimeout)
}




function _pauseAutoplay(){
    _stopAutoplayTimeout()
}
function _resumeAutoplay(){
    _setTimeoutAutoplay()
}

function _shuffle(array) {
  let currentIndex = array.length,  randomIndex;

  // While there remain elements to shuffle.
  while (currentIndex > 0) {

    // Pick a remaining element.
    randomIndex = Math.floor(Math.random() * currentIndex);
    currentIndex--;

    // And swap it with the current element.
    [array[currentIndex], array[randomIndex]] = [
      array[randomIndex], array[currentIndex]];
  }

  return array;
}

function _shuffleArrayIfNeeded(){
    if (_settings["randomAffirm"]){
        _affirmArr=_shuffle(_affirmArr)
    }
}


function _affirmArrIsEmpty(){
    if (_affirmArr.length<=0){
        return true
    }
    return false
}

function _showNoAffirmTextFav(){
    $(".app-affirm-text").text(NO_AFFIRM_INFO_TEXT)
    //set fav to false if no affirms
    _changeFavStatus(false)
}

function _getAutoplayDuration(){
    return _settings["autoplayDuration"]*1000
}


//EXPORTED
//----------------------------------------------------------------
//----------------------------------------------------------------
export async function resetPlayer(cardId){
    _curIndex=0;

    //affirmation text to empty
    $(".app-affirm-text").text("")

    let needBackendDataSettings = settingsNeedsBackendData()
    let needBackendDataAffirmArr = affirmArrNeedsBackendData(cardId)

    //get data from server
    if (needBackendDataSettings || needBackendDataAffirmArr){
        showLoader()
    }
    _settings = await getSettings()
    _affirmArr = await getAffirmArr(cardId)
    _shuffleArrayIfNeeded()

    if (needBackendDataSettings || needBackendDataAffirmArr){
        hideLoader()
    }

    //init player with data from server
    _autoplayDuration=_getAutoplayDuration()

    if (!_affirmArrIsEmpty()){
        //jquery
        $(".app-affirm-text").text(_affirmArr[_curIndex].content)
        //start with correct fav status
        _changeFavStatus(await _getAffirmFavStatus())
        //EXEC
        _stopAutoplayTimeout()
        if (_settings["autoplay"]===true){
            _setTimeoutAutoplay()
        }
    }else{
        _showNoAffirmTextFav()
        //EXEC
        _stopAutoplayTimeout()

    }

}

export function playerFunctionality(){

    ////back button was pressed, go home
    //  if (window.history && window.history.pushState) {

    //    window.history.pushState('forward', null, './#forward');
    //    $(window).on('popstate', function() {

    //    $("#app-player-btn-home").trigger('click')

    //    });
        //
            //
    //change cursor to hand on hover over home and fav affirm buttons
    $("#app-player-btn-home").hover(function(){
        $(this).css('cursor','pointer')
    })
    $("#app-player-btn-fav").hover(function(){
        $(this).css('cursor','pointer')
    })


    //back navigation
    $(window).on('popstate', function() {
      $("#app-player-btn-home").trigger("click")
    });   //  }

    //EVENTS
    //on click we display next affirmation
    $("#app-home-player").on("click",function(e){
        if (e.target.tagName === 'I'){
            e.stopPropagation();
        }else{
            //if no affirms, do nothing on click
            if (_affirmArrIsEmpty()){
                return
            }else{
                _playNextAffirm()
            }
        }
    })
    //on lose focuse we pause
    window.addEventListener('focus', function() {
        if (_playerIsVisible()){
            _resumeAutoplay()
        }
    });

    window.addEventListener('blur', function() {
        if (_playerIsVisible()){
            _pauseAutoplay()
        }
    });

    //on player home click
    //show home content
    $("#app-player-btn-home").on("click",function(e){
            // alert("playing");
            // $("#app-home-player").fadeTo(200,0)
            // $("#app-home-content").fadeTo(500,1)
            // $("#app-home-player").hide()
            // $("#app-home-content").show()
            _stopAutoplayTimeout() 
            _pauseAutoplay()
            $("#app-home-player").fadeTo(100,0,function(){
                $("#app-home-player").css({"display": "none"});

                $("#app-home-content").css({"display": "initial","opacity":0});
                $("#app-home-content").fadeTo(100,1)

                //back navig
                window.history.replaceState("forward",null,"./")
            })
    })
    //on player fav click
    $("#app-player-btn-fav").on("click",async function(e){

        try{

            let userAuthStatus = await getUserAuthStatus()
            //don't enable fav functionality if no affirms
            if (_affirmArrIsEmpty()){
                return
            }

            let favContent= $(favSelector).html()


            if (favContent === favUnfav){
                let data = JSON.stringify({"affirm":getCurrentAffirmTextFromPlayer(),"fav":true})
                await setAffirmFav(data)

                $(favSelector).html(favFav)
                $(favSelector).removeClass("text-type-dark-transp")
                $(favSelector).addClass("text-logo-main-transp")
                $("#app-player-btn-fav").sparkleHover({
                    colors : ['#de6300' ],
                    num_sprites: 22,
                    lifespan: 350,
                    radius: 50,
                    sprite_size: 10,
                    shape: 'circle',
                    gravity: false,
                    offset_x:0,
                    offset_y:-10,
                })


            }else{
                let data = JSON.stringify({"affirm":getCurrentAffirmTextFromPlayer(),"fav":false})
                await setAffirmFav(data)

                $(favSelector).html(favUnfav)
                $(favSelector).removeClass("text-logo-main-transp")
                $(favSelector).addClass("text-type-dark-transp")
            }

        }catch(error){
            log.error({error:error})
            if (error["status"]===401){
                showLoginModal()
                hideLoader()
            }else{
                showErrorModal()
                hideLoader()
            }

        }

    })

}
//----------------------------------------------------------------
//----------------------------------------------------------------
