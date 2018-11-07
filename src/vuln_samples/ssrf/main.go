package main

import (
    "net/http"
    "io/ioutil"
    "fmt"
    "os"
)

func usage() {
    fmt.Printf("usage: %s URL", os.Args[0])
    os.Exit(0)
}

func main() {
    if len(os.Args) < 2 {
        usage()
    }

    url := os.Args[1]
    resp, err := http.Get(url)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
            panic(err)
    }
    fmt.Printf("%s", body)
}