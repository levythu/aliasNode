package main

import (
    . "github.com/levythu/gurgling"
    "github.com/levythu/gurgling/midwares/cookie"
    "math/rand"
    "time"
    "net/url"
)

func init() {
    rand.Seed(time.Now().UnixNano())
}
const strFarm="qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM1234567890"
func generateRandStr(length int) string {
    var ret=make([]byte, length)
    for i:=0; i<length; i++ {
        ret[i]=strFarm[rand.Intn(len(strFarm))]
    }

    return string(ret)
}

func CreateManageRouter() Router {
    var r=ARouter().Use(cookie.ASession(generateRandStr(50)))
    r.Get("/callback", func(req Request, res Response) {
        res.Send("123")
    })
    r.Use(func(req Request, res Response) bool {
        var session=req.F()["session"].(map[string]string)
        if _, ok:=session["user"]; !ok {
            var parameters=url.Values{}
	        parameters.Add("destination", HOSTNAME_AND_SCHEMA+"/admin/callback")
            res.Redirect(VERIFY_URL+"?"+parameters.Encode())
            return false
        } else {
            res.Send("Hello!")
            return false
        }
    })

    return r
}
