package main

import (
    "github.com/levythu/gurgling/routers/simplefsserver"
    "fmt"
)

func main() {
    fmt.Println("Server running at port 2333")
    simplefsserver.ASimpleFSServer(".").Launch(":2333")
}
