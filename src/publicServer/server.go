package main

import (
    . "github.com/levythu/gurgling"
    "fmt"
)

func main() {
    fmt.Println("Server running at port 2335")
    var r=ARouter()

    r.Use("/admin", CreateManageRouter())

    r.UseSpecified("/", "GET", func(req Request, res Response) bool {
        var desURL=*(req.R().URL)
        var oriPath=desURL.Path
        if (len(oriPath)<=1) {
            return true
        }
        var p=1
        for (p<len(oriPath) && oriPath[p]!='/') {
            p++
        }
        var flag=oriPath[1:p]
        if p==len(oriPath) {
            desURL.Path=""
        } else {
            desURL.Path=oriPath[p:]
        }

        if result:=MapGet(flag); result!=nil {
            desURL.Host=result.Host
            desURL.Scheme=result.Scheme
            res.Redirect((&desURL).String())
            return false
        } else {
            return true
        }
    }, false)

    // catch all
    r.Use("/", func(req Request, res Response) {
        var desURL=*(req.R().URL)
        if result:=MapGet("."); result!=nil {
            desURL.Host=result.Host
            desURL.Scheme=result.Scheme
            res.Redirect((&desURL).String())
        } else {
            res.Status("NOT FOUND.", 404)
        }
    })

    r.Launch(":2335")
}
