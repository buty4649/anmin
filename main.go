package main

import (
    "fmt"
    "io/ioutil"
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

        if err:= exec(t); err != nil {
            panic(err)
            fin <- true
        }

        for {
            select {
            case <-sigCh:
                fin <- true
                return

            case <-ticker.C:
                if err:= exec(t); err != nil {
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

func exec(t *tsl256x.TSL256X) error{
    lux, err := t.ReadLux()
    if err != nil {
        return err
    }

    if lux > 30 {
        updateLEDStatus(0, "0",   "mmc0")
        updateLEDStatus(1, "255", "input")
    } else {
        updateLEDStatus(0, "0", "none")
        updateLEDStatus(1, "0", "none")
    }
    return nil
}

func updateLEDStatus(number int, brightness string, trigger string) error {
    base_path := fmt.Sprintf("/sys/class/leds/led%d/", number)

    err := ioutil.WriteFile(base_path + "brightness", []byte(brightness), os.ModePerm)
    if err != nil{
        return err
    }

    err = ioutil.WriteFile(base_path + "trigger", []byte(trigger), os.ModePerm)
    if err != nil{
        return err
    }

    return nil
}
