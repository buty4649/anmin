package tsl256x

import (
    "golang.org/x/exp/io/i2c"
    "math"
    "time"
)

const (
    DEV_ADDR = 0x39

    REG_CONTROL = 0x80
    REG_TIMING  = 0x81
    REG_ID      = 0x8A
    REG_CH0     = 0xAC
    REG_CH1     = 0xAE

    TIMING_1X  = 0x02 // 402ms low gain(1X)

    CONTROL_POWER_ON  = 0x03
    CONTROL_POWER_OFF = 0x00

    ID_TSL2560CS    = 0x00
    ID_TSL2561CS    = 0x01
    ID_TSL2560TFNCL = 0x04
    ID_TSL2561TFNCL = 0x05
)

type TSL256X struct {
    Device i2c.Device
    PartNo byte
    RevNo  byte
}

func Open(path string) (*TSL256X, error) {
    device, err := i2c.Open(&i2c.Devfs{Dev: path}, DEV_ADDR)
    if err != nil {
        return nil, err
    }

    buffer := make([]byte, 1)
    err = device.ReadReg(REG_ID, buffer)

    return &TSL256X{
        Device: *device,
        PartNo: buffer[0] >> 4,
        RevNo:  buffer[0] & 0x0f,
    },nil
}

func (t *TSL256X) ReadLux() (float64, error) {
    if err:= t.init(); err != nil {
        return 0.0, err
    }

    data,err := t.readData()
    if err != nil {
        return 0.0, err
    }

    return t.calcLux(data), nil
}

func (t *TSL256X) Close() error {
    return t.Device.Close()
}

func (t *TSL256X) init() (error) {
    err := t.Device.WriteReg(REG_TIMING, []byte{TIMING_1X})
    if err != nil {
        return err
    }

    err = t.Device.WriteReg(REG_CONTROL, []byte{CONTROL_POWER_ON})
    if err != nil {
        return err
    }
    time.Sleep(403 * time.Millisecond)

    return nil;
}

func (t *TSL256X) readData() ([]uint16, error) {
    buffer := make([]byte, 2)

    err := t.Device.ReadReg(REG_CH0, buffer)
    if err != nil {
        return nil, err
    }
    ch0 := uint16(buffer[0]) + uint16(buffer[1]) << 8

    err = t.Device.ReadReg(REG_CH1, buffer)
    if err != nil {
        return nil, err
    }
    ch1 := uint16(buffer[0]) + uint16(buffer[1]) << 8

    return []uint16{ch0, ch1}, nil
}

func (t *TSL256X) calcLux(data []uint16) (float64) {
    // scale 1X to 16
    ch0 := float64(data[0] * 16)
    ch1 := float64(data[1] * 16)

    lux := 0.0
    if t.PartNo < 4 {
        switch x := ch1 / ch0; {
            case x <= 0.52: lux = 0.0315  * ch0 - 0.0593  * ch0 * math.Pow(ch1/ch0, 1.4)
            case x <= 0.65: lux = 0.0229  * ch0 - 0.0291  * ch1
            case x <= 0.80: lux = 0.0157  * ch0 - 0.0180  * ch1
            case x <= 1.30: lux = 0.00338 * ch0 - 0.00260 * ch1
        }
    } else {
        switch x := ch1 / ch0; {
            case x <= 0.50: lux = 0.0304  * ch0 - 0.062   * ch0 * math.Pow(ch1/ch0, 1.4)
            case x <= 0.61: lux = 0.0224  * ch0 - 0.031   * ch1
            case x <= 0.80: lux = 0.0128  * ch0 - 0.0153  * ch1
            case x <= 1.30: lux = 0.00146 * ch0 - 0.00112 * ch1
        }
    }

    return lux
}
