package operadatatypes

import (
	"encoding/binary"
	"fmt"
	"io"
)

const DATA_TYPE_SPS30 = "S"

type Sps30Data struct {
	Pm1                 float32 `json:"pm1" binding:"required"`
	Pm2p5               float32 `json:"pm2p5" binding:"required"`
	Pm4                 float32 `json:"pm4" binding:"required"`
	Pm10                float32 `json:"pm10" binding:"required"`
	Pn0p5               float32 `json:"pm0p5" binding:"required"`
	Pn1                 float32 `json:"pn1" binding:"required"`
	Pn2p5               float32 `json:"pn2p5" binding:"required"`
	Pn4                 float32 `json:"pn4" binding:"required"`
	Pn10                float32 `json:"pn10" binding:"required"`
	TypicalParticleSize float32 `json:"typical_particle_size" binding:"required"`
}

func (d *Sps30Data) String() string {
	return fmt.Sprintf("[SPS30 Data| PM 1: %.1f, 2.5: %.1f, 4: %.1f, 10: %.1f | PN .5: %.1f, 1: %.1f, 2.5: %.1f, 4: %.1f | 10: %.1f]",
		d.Pm1, d.Pm2p5, d.Pm4, d.Pm10, d.Pn0p5, d.Pn1, d.Pn2p5, d.Pn4, d.Pn10)
}

func (d *Sps30Data) SendGob(unixSocketPath string) error {
	return sendStructGob(d, DATA_TYPE_SPS30, unixSocketPath)
}

func (d *Sps30Data) Pack(w io.Writer) {
	for _, val := range []float32{
		d.Pm1, d.Pm2p5, d.Pm4, d.Pm10, d.Pn0p5, d.Pn1, d.Pn2p5, d.Pn4, d.Pn10, d.TypicalParticleSize,
	} {
		binary.Write(w, binary.LittleEndian, val)
	}
}

func (d *Sps30Data) Unpack(r io.Reader) error {
	for _, val := range []*float32{
		&d.Pm1, &d.Pm2p5, &d.Pm4, &d.Pm10, &d.Pn0p5, &d.Pn1, &d.Pn2p5, &d.Pn4, &d.Pn10, &d.TypicalParticleSize,
	} {
		if err := binary.Read(r, binary.LittleEndian, val); err != nil {
			return err
		}
	}
	return nil
}
