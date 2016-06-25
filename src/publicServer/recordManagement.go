package main

import (
    "io/ioutil"
    "os"
    "encoding/json"
    "net/url"
    "fmt"
    "sync"
)

//Pay attention that filename is a relative path
func ReadFileToJSON(filename string) (map[string]string, error) {
    var err error

    var res []byte
    res, err=ioutil.ReadFile(filename)
    if err!=nil {
        return nil, err
    }

    var ret map[string]string
    err=json.Unmarshal(res, &ret)
    if err!=nil {
        return nil, err
    }

    return ret, nil
}

var globalMap map[string]*url.URL
const metadataFilename="data/metadata.json"
var mapLock=sync.RWMutex{}

func init() {
    var tmp, _=ReadFileToJSON(metadataFilename)
    if tmp==nil {
        tmp=make(map[string]string)
    }

    globalMap=make(map[string]*url.URL)
    for k, v:=range tmp {
        var tv, err=url.Parse(v)
        if err==nil {
            globalMap[k]=tv
        }
    }
}

func WriteBack(filename string) error {

    mapLock.RLock()

    var tmp=make(map[string]string)
    for k, v:=range globalMap {
        tmp[k]=v.String()
    }

    mapLock.RUnlock()

    var data, err=json.Marshal(tmp)
    if err!=nil {
        return err
    }
    return ioutil.WriteFile(filename, data, os.ModePerm)
}

func MapGet(key string) *url.URL {
    mapLock.RLock()
    defer mapLock.RUnlock()

    return globalMap[key]
}

func MapSet(key string, val *url.URL) {
    mapLock.Lock()
    if val==nil {
        delete(globalMap, key)
    } else {
        globalMap[key]=val
    }
    mapLock.Unlock()

    if err:=WriteBack(metadataFilename); err!=nil {
        fmt.Println("Writing error:", err);
    }
}

func MapDump() map[string]string {
    mapLock.RLock()

    var tmp=make(map[string]string)
    for k, v:=range globalMap {
        tmp[k]=v.String()
    }

    mapLock.RUnlock()

    return tmp
}
