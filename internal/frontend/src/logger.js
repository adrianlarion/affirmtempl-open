import {AFFIRMTEMPL_PROD} from './const.js'
var bunyan = require("bunyan")
var log = bunyan.createLogger({name:"affirmtempl", src:true})
//set this to FATAL in production
if(AFFIRMTEMPL_PROD==="no"){
    log.level(bunyan.TRACE)
}else{
    log.level(bunyan.FATAL)
}
export var log
