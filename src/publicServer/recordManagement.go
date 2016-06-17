package main

import (
    "io/ioutil"
    "os"
    "encoding/json"
    "net/url"
    "fmt"
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

    var tmp=make(map[string]string)
    for k, v:=range globalMap {
        tmp[k]=v.String()
    }

    var data, err=json.Marshal(tmp)
    if err!=nil {
        return err
    }
    return ioutil.WriteFile(filename, data, os.ModePerm)
}

func MapGet(key string) *url.URL {
    return globalMap[key]
}

func MapSet(key string, val *url.URL) {
    globalMap[key]=val
    if err:=WriteBack(metadataFilename); err!=nil {
        fmt.Println("Writing error:", err);
    }
}
