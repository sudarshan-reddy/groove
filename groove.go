package groove

import (
	"fmt"
	"time"
	"unsafe"

	"github.com/mrmorphic/hwio"
)

const (
	digitalRead  = 1
	digitalWrite = 2
	analogRead   = 3
	analogWrite  = 4
	pinMode      = 5
	dhtRead      = 0x44

	modeOutput    = "output"
	moduleI2C     = "i2c"
	responseDelay = 100 * time.Second
)

// Handler is the abstracted interface if someone wants to
// customize/ reimplement this
type Handler interface {
	AnalogRead(byte) (int, error)
	DigitalRead(byte) (byte, error)
	DigitalWrite(byte, byte) error
	ReadDHT(pin byte) (float32, float32, error)
	PinMode(byte, string) error
	Close()
}

// Groove holds the hwio i2c Ops
type Groove struct {
	i2cmodule hwio.I2CModule
	i2cDevice hwio.I2CDevice
}

// InitGroove initialises a new instance of Handlers to be used
func InitGroove(address int) (Handler, error) {
	module, err := hwio.GetModule(moduleI2C)
	if err != nil {
		return nil, err
	}
	i2cmodule := module.(hwio.I2CModule)
	i2cmodule.Enable()
	i2cDevice := i2cmodule.GetDevice(address)

	return &Groove{i2cmodule: i2cmodule, i2cDevice: i2cDevice}, nil
}

// Close disables the module effectively closing the device
// and can be a defer call
func (g *Groove) Close() {
	g.i2cmodule.Disable()
}

// AnalogRead reads the pin that's sent in the input
// And outputs a numerical number if there are no errors.
func (g *Groove) AnalogRead(pin byte) (int, error) {
	buffer := []byte{analogRead, pin, 0, 0}
	err := g.i2cDevice.Write(1, buffer)
	if err != nil {
		return 0, err
	}
	time.Sleep(responseDelay)
	g.i2cDevice.ReadByte(1)
	val, err := g.i2cDevice.Read(1, 4)
	if err != nil {
		return 0, err
	}
	return ((int(val[1]) * 256) + int(val[2])), nil
}

// DigitalRead reads value for the input pin.
func (g *Groove) DigitalRead(pin byte) (byte, error) {
	buffer := []byte{digitalRead, pin, 0, 0}
	err := g.i2cDevice.Write(1, buffer)
	if err != nil {
		return 0, err
	}
	time.Sleep(100 * time.Millisecond)
	val, err := g.i2cDevice.ReadByte(1)
	if err != nil {
		return 0, err
	}
	return val, nil
}

// DigitalWrite writes the input value to the input pin.
func (g *Groove) DigitalWrite(pin byte, val byte) error {
	buffer := []byte{digitalWrite, pin, val, 0}
	err := g.i2cDevice.Write(1, buffer)
	time.Sleep(100 * time.Millisecond)
	if err != nil {
		return err
	}
	return nil
}

// PinMode is to set the mode to output or input usually
// used before a digital write or read. It currently
// takes in output as 1 and assumes the rest as 0.
func (g *Groove) PinMode(pin byte, mode string) error {
	var buffer []byte
	if mode == modeOutput {
		buffer = []byte{pinMode, pin, 1, 0}
	} else {
		buffer = []byte{pinMode, pin, 0, 0}
	}
	err := g.i2cDevice.Write(1, buffer)
	time.Sleep(100 * time.Millisecond)
	if err != nil {
		return err
	}
	return nil
}

// ReadDHT reads raw data from the DHT sensors and
// parses them and returns a temperature value, a
// humidity value and an error
func (g *Groove) ReadDHT(pin byte) (float32, float32, error) {
	b := []byte{dhtRead, pin, 1, 0}
	rawdata, err := g.readDHTRawData(b)
	if err != nil {
		return 0, 0, err
	}
	fmt.Println("rawdata -> ", rawdata)
	temperatureData := rawdata[1:5]

	tInt := int32(temperatureData[0]) | int32(temperatureData[1])<<8 | int32(temperatureData[2])<<16 | int32(temperatureData[3])<<24
	t := (*(*float32)(unsafe.Pointer(&tInt)))

	humidityData := rawdata[5:9]
	humInt := int32(humidityData[0]) | int32(humidityData[1])<<8 | int32(humidityData[2])<<16 | int32(humidityData[3])<<24
	h := (*(*float32)(unsafe.Pointer(&humInt)))
	return t, h, nil
}

func (g *Groove) readDHTRawData(buffer []byte) ([]byte, error) {
	err := g.i2cDevice.Write(0x02C, buffer)
	if err != nil {
		return nil, err
	}
	time.Sleep(600 * time.Millisecond)
	g.i2cDevice.ReadByte(1)
	time.Sleep(100 * time.Millisecond)
	raw, err := g.i2cDevice.Read(0x00, 9)
	if err != nil {
		return nil, err
	}
	return raw, nil
}
