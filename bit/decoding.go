package bit

import (
	"encoding/binary"
)

type Informacao struct {
	index  uint32
	packed []byte
}

func NovaInformacao(packed []byte) *Informacao {
	return &Informacao{packed: packed}
}

func (bits *Informacao) Pop(size uint32) (value uint32) {
	byteIndex := (bits.index/8)/4
	offset := bits.index % 32
	value = binary.BigEndian.Uint32(bits.packed[byteIndex : byteIndex+4])
	shift := 32 - size
	value = uint32(value<<(offset)) >> shift
	if offset > shift {
		remain := offset - shift
		low := binary.BigEndian.Uint32(bits.packed[byteIndex+4:byteIndex+8]) >> (32 - remain)
		value = value | low
	}
	bits.index += size
	return
}

func Decode(B uint32, unpacker *Informacao) uint32 {
	max := uint32(1) << uint32(B)
	if unpacker.Pop(uint32(1)) == 0 {
		return unpacker.Pop(B)
	}
	return max + Decode(B+1, unpacker)
}
