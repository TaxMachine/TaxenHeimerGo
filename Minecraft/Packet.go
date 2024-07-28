package Minecraft

import (
	"bytes"
	"fmt"
	"os"
)

const (
	SegmentBits            uint32 = 0x7F
	SegmentContinuationBit uint32 = 0x80
)

type Packet struct {
	buffer bytes.Buffer
}

func NewMinecraftPacket() *Packet {
	return &Packet{}
}

func (mp *Packet) WriteVarInt(value uint32) error {
	var remaining uint32 = value
	for i := 0; i < 5; i++ {
		if remaining&^SegmentBits == 0 {
			mp.buffer.WriteByte(byte(remaining))
			return nil
		}
		mp.buffer.WriteByte(byte(remaining&SegmentBits | SegmentContinuationBit))
		remaining >>= 7
	}
	return fmt.Errorf("VarInt too big")
}

func (mp *Packet) WriteString(value string) (err error) {
	err = mp.WriteVarInt(uint32(len(value)))
	mp.buffer.WriteString(value)
	return nil
}

func (mp *Packet) WriteShort(value uint16) {
	mp.buffer.WriteByte(byte((value >> 8) & 0xFF))
	mp.buffer.WriteByte(byte(value & 0xFF))
}

func (mp *Packet) WriteLong(value uint64) error {
	var remaining uint64 = value
	for i := 0; i < 5; i++ {
		if (remaining & (uint64(SegmentBits) ^ 0xffffffffffffffff)) == 0 {
			mp.buffer.WriteByte(byte(remaining))
			return nil
		}
		mp.buffer.WriteByte(byte((remaining & uint64(SegmentBits)) | uint64(SegmentContinuationBit)))
		remaining >>= 7
	}
	return fmt.Errorf("VarLong too big")
}

func (mp *Packet) WriteBuffer(buffer []byte) {
	mp.buffer.Write(buffer)
}

func (mp *Packet) GetBuffer() []byte {
	return mp.buffer.Bytes()
}

func (mp *Packet) GetBufferStream() (file *os.File, err error) {
	file, err = os.CreateTemp("", "buffer")
	if err != nil {
		panic(err)
	}
	_, err = file.Write(mp.buffer.Bytes())
	_, err = file.Seek(0, 0)
	return
}

func (mp *Packet) Size() int {
	return mp.buffer.Len()
}

func (mp *Packet) PrintHex() {
	data := mp.buffer.Bytes()
	for _, b := range data {
		fmt.Printf("%02x ", b)
	}
	fmt.Println()
}
