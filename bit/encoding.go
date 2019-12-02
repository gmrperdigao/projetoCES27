package bit

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

/// Função responsavel por trabalhar com os pacotes de informacao do clock, como colocar no pacote, encoding, etc
// Exemplo: <<2:3,0:1,1:2>> representa 0100001, x:y representa o valor x codificado em y bits
type Entrada struct {
	Value uint32
	Size uint32
}

type Pacote struct {
	Comprimento uint32
	Entradas   []Entrada
	packed    []byte
}

func NovoPacote() *Pacote {
	return &Pacote{Entradas: make([]Entrada, 0), packed: make([]byte, 4)}
}

func (p *Pacote) Push(value, size uint32) {
	p.Entradas = append(p.Entradas, Entrada{Value: value, Size: size})
	bitsLivres := uint32(8*len(p.packed)) - p.Comprimento
	index := p.Comprimento/32
	if bitsLivres >= size {
		shift := bitsLivres - size
		v := binary.BigEndian.Uint32(p.packed[index:]) | (value << shift)
		binary.BigEndian.PutUint32(p.packed[index:], v)
	} else {
		buf := make([]byte, 4)
		low := value
		if bitsLivres > 0 {
			high := value>>size - bitsLivres
			v := binary.BigEndian.Uint32(p.packed[index:]) | high
			binary.BigEndian.PutUint32(p.packed[index:], v)
			low = (value << bitsLivres) >> bitsLivres
		} else {
			shift := 32 - size
			low = value<<shift
		}
		binary.BigEndian.PutUint32(buf, low)
		p.packed = append(p.packed[:], buf[:]...)
	}
	p.Comprimento += size
}

func (p *Pacote) Pacote() []byte {
	return p.packed
}

func (p *Pacote) PackedString() string {
	var b bytes.Buffer
	for i, w := range p.packed {
		remainingBits := int32(p.Comprimento) - int32(i*8)
		if remainingBits < 0 {
			break
		}
		if remainingBits < 8 {
			b.WriteString(fmt.Sprintf("%0*b", int(remainingBits), w>>uint(8-remainingBits)))
		} else {
			b.WriteString(fmt.Sprintf("%0*b", int(8), w))
		}
	}
	return b.String()
}

func (p *Pacote) String() string {
	var buf bytes.Buffer
	buf.WriteString("<<")
	for i, entry := range p.Entradas {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(fmt.Sprintf("%d:%d", entry.Value, entry.Size))
	}
	buf.WriteString(">>")
	return buf.String()
}

func Encode(n, B uint32, packer *Pacote) *Pacote {
	max := uint32(1) << B
	if n < max {
		packer.Push(uint32(0), uint32(1))
		packer.Push(n, B)
	} else {
		packer.Push(uint32(1), uint32(1))
		Encode(n-max, B+1, packer)
	}
	return packer
}
