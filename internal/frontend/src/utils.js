import {MAX_AFFIRMATIONS} from './const.js'

//globals
let cardId=0


export function fxFlash(){
    $("html").fadeOut(198).fadeIn(400)
}

export function showLoader(){
    // $("#cover-spin").show(0)
    $("#cover-spin").fadeIn(200)
}

export function showErrorModal(){
    // $("#modal-error").fadeIn(200)
    // $("#modal-error").modal('open')
    var elem = document.querySelector('#modal-error');
    var instance = M.Modal.getInstance(elem)
    instance.open()
    $("#modal-error").css("z-index",99999)
}

export function showLoginModal(){
    // $("#modal-error").fadeIn(200)
    // $("#modal-error").modal('open')
    var elem = document.querySelector('#modal-login');
    var instance = M.Modal.getInstance(elem)
    instance.open()
    $("#modal-login").css("z-index",99999)
}

export function showSettingsModal(){
    // $("#modal-error").fadeIn(200)
    // $("#modal-error").modal('open')
    var elem = document.querySelector('#modal-settings');
    var instance = M.Modal.getInstance(elem)
    instance.open()
    // $("#modal-login").css("z-index",99999)
}


export function showEditAffirmModal(){
    var elem = document.querySelector('#modal-edit-affirm');
    var instance = M.Modal.getInstance(elem)
    instance.open()

}

export function getSelectedCardId(){
    return cardId
}

export function setSelectedCardId(_cardId){
    cardId=_cardId
}


export function hideLoader(){
    // $("#cover-spin").hide(0)
    $("#cover-spin").fadeOut(200)
}

export function populateEditAffirmModal(affirmArr){
    //remove old text, clear the fields
    for (let i=0; i<MAX_AFFIRMATIONS; i++){
        let sel ="#modal-edit-affirm > div.modal-content > div > div:nth-child(" + (i+1) +") textarea" 
        $(sel).val("")

    }
    //put in new field vals
    for (let i=0; i<affirmArr.length; i++){
        let sel ="#modal-edit-affirm > div.modal-content > div > div:nth-child(" + (i+1) +") textarea" 
        $(sel).val(affirmArr[i].content)
    }
    M.updateTextFields()
}

export function getSettingsFromModal(){
    let randomAffirm = $("#modal-settings #modal-settings-check-rand-affirm").prop("checked");
    let autoplay = $("#modal-settings #modal-settings-check-autoplay").prop("checked");
    let autoplayDuration =$("#modal-settings-select-autoplay-duration").val()
    let settings ={
        "autoplay":autoplay,
        "autoplayDuration":autoplayDuration,
        "randomAffirm":randomAffirm
    }
    return settings
}

export function getCardFavStatusFromHtml(cardId){
    let elemText = $("[data-cardid="+cardId+"] i.card-fav-btn").text()
    if (elemText === "favorite"){
        return true
    }
    return false
}

export function getCurrentAffirmTextFromPlayer(){
    let elemVal = $("h1.app-affirm-text").text()
    return elemVal
}


export function getAffirmArrFromModal(){
    let affirmArr = []
    for (let i=0; i<MAX_AFFIRMATIONS; i++){
        let sel ="#modal-edit-affirm > div.modal-content > div > div:nth-child(" + (i+1) +") textarea" 
        // $(sel).val("")
        let val =$(sel).val()
        if (val.length !==0){
            affirmArr.push(val)
        }
    }
    return {"affirmations":affirmArr}
    // return affirmArr
}




export function populateSettingsModal(settings){
    $("#modal-settings #modal-settings-check-rand-affirm").prop("checked",settings["randomAffirm"]);
    $("#modal-settings #modal-settings-check-autoplay").prop("checked",settings["autoplay"]);

    if (settings["autoplayDuration"]===10){
        $("#modal-settings-select-autoplay-duration").val("10")
    } else if (settings["autoplayDuration"]===15){
        $("#modal-settings-select-autoplay-duration").val("15")
    }else if (settings["autoplayDuration"]===25){
        $("#modal-settings-select-autoplay-duration").val("25")
    }else if (settings["autoplayDuration"]===40){
        $("#modal-settings-select-autoplay-duration").val("40")
    }else if (settings["autoplayDuration"]===60){
        $("#modal-settings-select-autoplay-duration").val("60")
    }else if (settings["autoplayDuration"]===120){
        $("#modal-settings-select-autoplay-duration").val("120")
    }

    $("#modal-settings-select-autoplay-duration").formSelect();

}
