package main

import (
    "fmt"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/buty4649/anmin/driver"
)

func main() {

    t, err := tsl256x.Open("/dev/i2c-1")
    if err != nil {
        panic(err)
    }

    fin := make(chan bool)
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, syscall.SIGINT, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

    go func() {
        ticker := time.NewTicker(1 * time.Second)

        if err:= printLux(t); err != nil {
            panic(err)
            fin <- true
        }

        for {
            select {
            case <-sigCh:
                fin <- true
                return

            case <-ticker.C:
                if err:= printLux(t); err != nil {
                    panic(err)
                    fin <- true
                }
            }
        }
    }()
    <-fin

    t.Close()
    os.Exit(0)
}

func printLux(t *tsl256x.TSL256X) error{
    lux, err := t.ReadLux()
    if err != nil {
        return err
    }
    fmt.Printf("%.1f\n", lux)
    return nil
}
