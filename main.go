package main

import (
    "fmt"
    "github.com/buty4649/anmin/driver"
)

func main() {
    t, err := tsl256x.Open("/dev/i2c-1")
    if err != nil {
        panic(err)
    }

    lux, err := t.ReadLux()
    if err != nil {
        panic(err)
    }

    fmt.Printf("%.1f\n", lux)
    t.Close()
}
