package sen54

import (
	"encoding/binary"
	"fmt"
	"periph.io/x/conn/v3/i2c"
	"sync"
	"time"
)

const (
	SensorAddr               = 0x69
	StartPeriodicMeasurement = 0x0021
	ReadMeasurement          = 0x03C4
	StopPeriodicMeasurement  = 0x0104
	CRC8_POLYNOMIAL          = byte(0x31)
	CRC8_INIT                = byte(0xff)
)

// Supported commands
var StartCommand = Command{
	cmd:       StartPeriodicMeasurement,
	respBytes: 0,
	delay:     100 * time.Millisecond,
	desc:      "start periodic measurements",
}
var StopCommand = Command{
	cmd:       StopPeriodicMeasurement,
	respBytes: 0,
	delay:     250 * time.Millisecond,
	desc:      "stop periodic measurements",
}
var MeasureCommand = Command{
	cmd:       ReadMeasurement,
	respBytes: 24,
	delay:     70 * time.Millisecond,
	desc:      "read sensor metrics",
}

// The sensor can't handle mulitple commands at once
var mu sync.Mutex

type Sen54 struct {
	dev *i2c.Dev
}

func New(b i2c.Bus) *Sen54 {
	return &Sen54{
		dev: &i2c.Dev{Addr: SensorAddr, Bus: b},
	}
}

type SensorData struct {
	PM1_0    float64
	PM2_5    float64
	PM4      float64
	PM10     float64
	Humidity float64
	Temp     float64 // deg C
	VOC      float64
}

type Command struct {
	cmd       uint16        // hex code from data sheet
	respBytes uint16        // expected response size (typically 0, 3, or 9)
	delay     time.Duration // time to sleep after cmd
	desc      string        // useful description for error messages
}

type Response struct {
	data []byte // expected to be two bytes
	crc  byte   // CRC8 sent by sensor of previous two bytes
}

func (r Response) CrcMatch() bool {
	return crc8(r.data, uint16(len(r.data))) == r.crc
}

func (r Response) GetData() uint16 {
	return binary.BigEndian.Uint16(r.data)
}

func (s *Sen54) StartMeasurements() error {
	mu.Lock()
	defer mu.Unlock()
	if err := s.sendCommand(StartCommand); err != nil {
		return err
	}
	return nil
}

func (s *Sen54) StopMeasurements() error {
	mu.Lock()
	defer mu.Unlock()
	if err := s.sendCommand(StopCommand); err != nil {
		return err
	}
	return nil
}

func (s *Sen54) ReadMeasurement() (SensorData, error) {
	mu.Lock()
	defer mu.Unlock()
	var result SensorData
	resp, err := s.readCommand(MeasureCommand)
	if err != nil {
		return result, err
	}
	// check CRCs
	for _, r := range resp {
		if !r.CrcMatch() {
			return result, fmt.Errorf("measuerment CRC mismatch")
		}
	}
	result = SensorData{
		PM1_0:    float64(resp[0].GetData()) / 10,
		PM2_5:    float64(resp[1].GetData()) / 10,
		PM4:      float64(resp[2].GetData()) / 10,
		PM10:     float64(resp[3].GetData()) / 10,
		Humidity: 100 * (float64(int16(resp[4].GetData())) / 100),
		Temp:     float64(int16(resp[5].GetData())) / 200,
		VOC:      float64(int16(resp[6].GetData())) / 10,
	}
	return result, nil
}

// Adapted from the C/C++ example in the SDC4x data sheet
func crc8(data []byte, count uint16) byte {
	crc := CRC8_INIT
	for currentByte := uint16(0); currentByte < count; currentByte++ {
		crc ^= data[currentByte]
		for crcBit := 8; crcBit > 0; crcBit-- {
			if crc&0x80 != 0 {
				crc = (crc << 1) ^ CRC8_POLYNOMIAL
			} else {
				crc = crc << 1
			}
		}
	}
	return crc
}

func (s *Sen54) sendCommand(cmd Command) error {
	c := make([]byte, 2)
	binary.BigEndian.PutUint16(c, cmd.cmd)
	if err := s.dev.Tx(c, nil); err != nil {
		return fmt.Errorf("error while %s: %v", cmd.desc, err)
	}
	if cmd.delay > 0 {
		time.Sleep(cmd.delay)
	}
	return nil
}

func (s *Sen54) readCommand(cmd Command) ([]Response, error) {
	c := make([]byte, 2)
	r := make([]byte, cmd.respBytes)
	binary.BigEndian.PutUint16(c, cmd.cmd)
	if err := s.dev.Tx(c, r); err != nil {
		return nil, fmt.Errorf("error while %s: %v", cmd.desc, err)
	}
	resp := []Response{}
	for i := 0; i < int(cmd.respBytes)-2; i += 3 {
		j := Response{data: r[i : i+2], crc: r[i+2]}
		resp = append(resp, j)
	}
	if cmd.delay > 0 {
		time.Sleep(cmd.delay)
	}
	return resp, nil
}
