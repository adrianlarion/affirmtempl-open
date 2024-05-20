import {resetPlayer} from './player.js'
import {showSettingsModal, showLoginModal, showErrorModal,getSelectedCardId,setSelectedCardId,showLoader, hideLoader,populateSettingsModal,showEditAffirmModal,populateEditAffirmModal, getSettingsFromModal, getAffirmArrFromModal,getCardFavStatusFromHtml} from './utils.js'

import {logoutUser,getUserAuthStatus,getDefaultAffirmArr, getSettings,setSettings,settingsNeedsBackendData,getAffirmArr,affirmArrNeedsBackendData,setAffirmArr, setCardFav} from './data.js'
import {log} from './logger.js'

var cardId

export function homeFunctionality(){

    //on reload remember scroll position
     if (localStorage.getItem("easyaffirm") != null) {
        $(window).scrollTop(localStorage.getItem("easyaffirm-quote-scroll"));
    }

    $(window).on("scroll", function() {
        localStorage.setItem("easyaffirm-quote-scroll", $(window).scrollTop());
    });

    //change cursor to hand on hover over edit affirm and fav buttons
    $("i.card-fav-btn").hover(function(){
        $(this).css('cursor','pointer')
    })
    $("i.card-edit-affirm-btn").hover(function(){
        $(this).css('cursor','pointer')
    })


    //logout btn
    $(".header-logout-btn").on("click",async function(e){
        try{
            showLoader()
            await logoutUser()
            hideLoader()
            location.reload()
        }catch(error){
            log.error({error:error})
            hideLoader()
            showErrorModal()
        }
    })


    //on error modal refresh button click refresh page
    $("#modal-error-refresh-btn").on("click",function(e){
        location.reload()
    })

    //on licking the settings button
    $(".header-settings-btn").on("click",async function(e){
        try{
            let userAuthStatus = await getUserAuthStatus()

            let needsBackendData = settingsNeedsBackendData()
            //only show loader if we need to retrieve data from backend
            if (needsBackendData){
                showLoader()
            }
            let settings = await getSettings()
            populateSettingsModal(settings)
            showSettingsModal()

            if (needsBackendData){
                hideLoader()
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

    // on closing the settings modal
    //on closing settings modal
    var settingsElem = document.querySelector('#modal-settings')
    M.Modal.init(settingsElem,{onCloseStart:async function(){
        let settings=getSettingsFromModal()
        let settings_json = JSON.stringify(settings)

        showLoader()
        await setSettings(settings_json)
        hideLoader()

    }})

    //on closing the edit affirm modal
    var editAffirmModal = document.querySelector('#modal-edit-affirm')
    M.Modal.init(editAffirmModal,{onCloseStart:async function(){
        let affirmArr=getAffirmArrFromModal()
        let affirmArrJson = JSON.stringify(affirmArr)

        showLoader()
        await setAffirmArr(affirmArrJson,cardId)
        hideLoader()

    }})

    //on pressing the reset affirm button in edit affirm modal
    $(".edit-affirm-arr-reset-btn").on("click",async function(e){
        showLoader()
        //get default affirm arr
        let defaultAffirmArr=await getDefaultAffirmArr(cardId)
        populateEditAffirmModal(defaultAffirmArr)

        //if you uncomment this, we will update data on server on button click
        //else we update server on edit affimr arr exit, but we'll do this twice since we also send data on popup close
        //
        //send this new data to the server 
        //
        // let affirmArr=getAffirmArrFromModal()
        // let affirmArrJson = JSON.stringify(affirmArr)
        // await setAffirmArr(affirmArrJson,cardId)

        hideLoader()
    })



    //Show player on card click
    $(".card").on("click",async function(e){
        setSelectedCardId(this.dataset.cardid)
        cardId=getSelectedCardId()
        // cardId=getSelectedCardId()

        if (e.target.tagName === 'I'){
            //on edit affirm click
            if (e.target.classList.contains("card-edit-affirm-btn")){

                try{
                    //check if user logged in, else show login popup
                    let userAuthStatus = await getUserAuthStatus()

                    let needBackendDataAffirmArr=affirmArrNeedsBackendData(cardId)
                    if (needBackendDataAffirmArr){
                        showLoader()
                    }
                    let affirmArr = await getAffirmArr(cardId)
                    populateEditAffirmModal(affirmArr)
                    showEditAffirmModal(cardId)

                    if (needBackendDataAffirmArr){
                        hideLoader()
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



            }else if (e.target.classList.contains("card-fav-btn")){
                //do backend stuff

                try{
                    let userAuthStatus = await getUserAuthStatus()
                    //backend call
                    showLoader()
                    let data = JSON.stringify({"fav":!getCardFavStatusFromHtml(cardId)})
                    await setCardFav(data,cardId)

                    //animation
                    let favSelector = "[data-cardid="+cardId+"] i.card-fav-btn"
                    let favUnfav = 'favorite_border'
                    let favFav = 'favorite'
                    let favContent= $(favSelector).html()
                    if (favContent === favUnfav){
                        //call backend with fav
                        $(favSelector).removeClass("text-type-dark-transp")
                        $(favSelector).addClass("text-logo-main-transp")

                        $(favSelector).html(favFav)
                        $(favSelector).sparkleHover({
                            colors : ['#de6300' ],
                            num_sprites: 22,
                            lifespan: 350,
                            radius: 25,
                            sprite_size: 5,
                            shape: 'circle',
                            gravity: false,
                            offset_x:0,
                            offset_y:0,
                        })
                        //wait for a bit to let animation finish
                        await new Promise(r => setTimeout(r, 300));

                    }else{
                        //call backend with unfav
                        //wait for a bit to let animation finish
                        await new Promise(r => setTimeout(r, 200));
                        $(favSelector).removeClass("text-logo-main-transp")
                        $(favSelector).addClass("text-type-dark-transp")
                        $(favSelector).html(favUnfav)
                    }

                    location.reload();

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

            }
            e.stopPropagation();
        }else{
            //open affirmation playkk
            // alert("playing");
            // $("#app-home-content").css({"display": "none"});
            // $("#app-home-player").css({"display": "initial"});
            // $("#app-home-content").fadeTo(200,0)
            // $("#app-home-player").fadeTo(500,1)
            // $("#app-home-content").hide()
            // $("#app-home-player").show()
            resetPlayer(cardId)
            $("#app-home-content").fadeTo(100,0,function(){
                $("#app-home-content").css({"display": "none"});
                //set height manually because we're using absolute position for child elements
                $("#app-home-player").css({"display": "block","opacity":0,"height":$(window).height()+"px"});
                $("#app-home-player").fadeTo(100,1);

                //back navigation
                window.history.pushState('forward', null, './#forward');

            })



        }
    })

}
