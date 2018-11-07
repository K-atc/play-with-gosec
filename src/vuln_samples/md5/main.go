package main

import (
    "fmt"
    "crypto/md5"
)

func main(){

    x := 1
    y := 2

    if x < y {
        h := md5.New()
        fmt.Printf("%x", h.Sum(nil))
    }
}