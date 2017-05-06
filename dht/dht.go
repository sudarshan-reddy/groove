package dht

import (
	"time"

	smbus "github.com/corrupt/go-smbus"
)

// SmbData is the structure returned by Read
type SmbData struct {
	CelsiusTemp   int32
	FarenheitTemp int32
	Humidity      int32
}

// ReadDHT reads the I2C bus and responds with calculated
// temperature in celsius, farenheit and humidity.
func ReadDHT() (*SmbData, error) {
	smb := &smbus.SMBus{}
	err := smb.Bus_open(1)
	if err != nil {
		return nil, err
	}
	defer smb.Bus_close()
	err = smb.Set_addr(0x44)
	if err != nil {
		return nil, err
	}
	buf := make([]byte, 6)
	err = smb.Write_byte_data(0x2C, 0x06)
	if err != nil {
		return nil, err
	}
	time.Sleep(1 * time.Second)
	_, err = smb.Read_i2c_block_data(0x00, buf)
	if err != nil {
		return nil, err
	}
	c, f := toTemp(buf)
	h := toHumidity(buf)
	return &SmbData{CelsiusTemp: c, FarenheitTemp: f, Humidity: h}, nil
}

func toTemp(buf []byte) (int, int) {
	temp := (int32(buf[0])*256 + int32(buf[1]))
	c := -45 + (175 * temp / 65535.0)
	f := -49 + (315 * temp / 65535.0)
	return c, f
}

func toHumidity(buf []byte) int {
	return 100 * (int32(buf[3])*256 + int32(buf[4])) / 65535.0
}
