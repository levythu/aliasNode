package main

import (
    . "github.com/levythu/gurgling"
    "github.com/levythu/gurgling/midwares/cookie"
    "math/rand"
    "time"
    "crypto/tls"
    "net/http"
    "io/ioutil"
    "net/url"
    "fmt"
    "strings"
)

var pipeClient=&http.Client{}
var piptClientSSL=&http.Client{Transport: &http.Transport{
	TLSClientConfig: &tls.Config{},
}}
var client *http.Client

func init() {
    rand.Seed(time.Now().UnixNano())
    client=pipeClient
    if r, e:=url.Parse(VERIFY_URL); e==nil {
        if strings.ToLower(r.Scheme)=="https" {
            client=piptClientSSL
        }
    }
}


const strFarm="qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM1234567890"
func generateRandStr(length int) string {
    var ret=make([]byte, length)
    for i:=0; i<length; i++ {
        ret[i]=strFarm[rand.Intn(len(strFarm))]
    }

    return string(ret)
}

func VerifyRequest(token string) bool {
    var proxyRequest, err=http.NewRequest("GET", VERIFY_URL+"?token="+token, strings.NewReader(""))
    if err!=nil {
        return false
    }

    var proxyResponse, err2=client.Do(proxyRequest)
    if err2!=nil {
        return false
    }

    var content, err3=ioutil.ReadAll(proxyResponse.Body)
    defer proxyResponse.Body.Close()
    if err3!=nil {
        return false
    }

    if string(content)=="TRUE" {
        return true
    }
    return false
}

func CreateManageRouter() Router {
    var r=ARouter().Use(cookie.ASession(generateRandStr(50)))
    r.Get("/callback", func(req Request, res Response) {
        if VerifyRequest(req.Query()["token"]) {
            var session=req.F()["session"].(map[string]string)
            session["user"]="authedUser"
            res.Redirect("/admin")
        } else {
            res.Redirect("/admin/fail")
        }
    })
    r.Get("/fail", func(req Request, res Response) {
        res.Send("Auth fail.")
    })
    r.Use(func(req Request, res Response) bool {
        var session=req.F()["session"].(map[string]string)
        if _, ok:=session["user"]; !ok {
            var parameters=url.Values{}
	        parameters.Add("destination", HOSTNAME_AND_SCHEMA+"/admin/callback")
            res.Redirect(CREATE_URL+"?"+parameters.Encode())
            return false
        } else {
            res.Send("Hello!")
            return false
        }
    })

    return r
}
