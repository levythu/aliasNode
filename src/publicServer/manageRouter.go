package main

import (
    . "github.com/levythu/gurgling"
    "github.com/levythu/gurgling/midwares/cookie"
    "github.com/levythu/gurgling/midwares/staticfs"
    "github.com/levythu/gurgling/midwares/bodyparser"
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
    // if r, e:=url.Parse(VERIFY_URL); e==nil {
    //     if strings.ToLower(r.Scheme)=="https" {
    //         client=piptClientSSL
    //     }
    // }
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
        fmt.Println(err)
        return false
    }

    var proxyResponse, err2=client.Do(proxyRequest)
    if err2!=nil {
        fmt.Println(err2)
        return false
    }

    var content, err3=ioutil.ReadAll(proxyResponse.Body)
    defer proxyResponse.Body.Close()
    if err3!=nil {
        fmt.Println(err3)
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

    r.Get("/list", func(res Response) {
        res.JSON(MapDump())
    })

    // auth
    r.Use(func(req Request, res Response) bool {
        var session=req.F()["session"].(map[string]string)
        if _, ok:=session["user"]; !ok {
            var parameters=url.Values{}
	        parameters.Add("destination", HOSTNAME_AND_SCHEMA+"/admin/callback")
            res.Redirect(CREATE_URL+"?"+parameters.Encode())
            return false
        } else {
            return true
        }
    })

    r.Use(staticfs.AStaticfs("public"))

    r.Get("/", func(req Request, res Response) {
        res.SendFile("./public/manage.general.html")
    })

    r.Use("/modify", bodyparser.ABodyParser())
    r.Post("/modify", func(req Request, res Response) {
        var bd, ok=req.Body().(url.Values)
        if (!ok || len(bd["oldk"])==0 || len(bd["k"])==0 || len(bd["v"])==0) {
            res.SendCode(403)
            return
        }
        var oldk=bd["oldk"][0]
        var k=bd["k"][0]
        var v=bd["v"][0]
        if oldk=="" {
            if k!="" && v!="" && MapGet(k)==nil {
                var v2, err=url.Parse(v)
                if err==nil {
                    MapSet(k, v2)
                    res.Send("OK")
                    return
                }
            }
            res.SendCode(403)
            return
        } else {
            MapSet(oldk, nil)
            if k!="" && MapGet(k)==nil {
                var v2, err=url.Parse(v)
                if err==nil {
                    MapSet(k, v2)
                    res.Send("OK")
                    return
                }
            }
            res.SendCode(204)
            return
        }
    })

    return r
}
