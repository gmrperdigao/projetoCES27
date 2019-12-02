package itc

import (
	"fmt"
	"itc/bit"
	"itc/event"
	"itc/id"
)

// Stamp é o relógio de cada maquina
type Stamp struct {
	event *event.Event
	id    *id.ID
}

// Criacao de uma stamp raiz (1, 0)
func NewStamp() *Stamp {
	return &Stamp{event: event.New(), id: id.New()}
}

// A existencia de um novo evento adiciona esse evento no componete: (i, e') vem de (i, e), com  e < e'.
func (s *Stamp) Event() {
	antes := s.event.Clona()
	depois := s.fill()
	if depois.Equals(antes) {
		s.event, _ = s.grow()
	} else {
		s.event = depois
	}
}

// Realiza a operacao fill:
func (s *Stamp) fill() *event.Event {
	return fill(s.id, s.event)
}

// Realiza a operacao grow se nao for possivel simplificar a event tree
func (s *Stamp) grow() (*event.Event, int) {
	return grow(s.id, s.event)
}

// Fork clona os componentes passados, resultando em elementos com mesmos eventos e ids diferentes
func (s *Stamp) Fork() *Stamp {
	st := NewStamp()
	id1, id2 := s.id.Split()
	s.id = id1
	st.id = id2
	st.event = s.event.Clona()
	return st
}

// Junta duas stamps produzindo uma nova
func (s *Stamp) Join(other *Stamp) {
	s.id = id.New().Soma(s.id, other.id)
	s.event = event.Join(s.event, other.event)
}

// Tipo especial de fork
func (s *Stamp) Peek() *Stamp {
	peekSt := NewStamp()
	peekSt.event = s.event.Clona()
	return peekSt
}

// Compara duas stamps, retornando verdade de menor ou igual que
func (s *Stamp) LEQ(other *Stamp) bool {
	return event.LEQ(s.event, other.event)
}

// Escreve a stamp e retorna uma string com seu valor
func (s *Stamp) String() string {
	return fmt.Sprintf("(%s, %s)", s.id, s.event)
}

// Calculo de fill
func fill(i *id.ID, e *event.Event) *event.Event {
	if i.Folha {
		if i.Valor == 0 {
			return e
		}
		return event.NovaFolha(e.Max())
	}
	if e.Folha {
		return e
	}
	r := event.NovoNodeVazio(e.Valor)
	if i.Fesq.Folha && i.Fesq.Valor == 1 {
		r.Fdir = fill(i.Fdir, e.Fdir)
		r.Fesq = event.NovaFolha(event.Maior(e.Fesq.Max(), r.Fdir.Min()))
	} else if i.Fdir.Folha && i.Fdir.Valor == 1 {
		r.Fesq = fill(i.Fesq, e.Fesq)
		r.Fdir = event.NovaFolha(event.Maior(e.Fdir.Max(), r.Fesq.Min()))
	} else {
		r.Fesq = fill(i.Fesq, e.Fesq)
		r.Fdir = fill(i.Fdir, e.Fdir)
	}
	return r.Normal()
}

// Calculo de grow
func grow(i *id.ID, e *event.Event) (*event.Event, int) {
	if e.Folha {
		if i.Folha && i.Valor == 1 {
			return event.NovaFolha(e.Valor + 1), 0
		}
		ex, c := grow(i, event.NovoNodeVazio(e.Valor))
		return ex, c + 99999
	}
	if i.Fesq.Folha && i.Fesq.Valor == 0 {
		exr, cr := grow(i.Fdir, e.Fdir)
		r := event.NovoNodeVazio(e.Valor)
		r.Fesq = e.Fesq
		r.Fdir = exr
		return r, cr + 1
	}
	if i.Fdir.Folha && i.Fdir.Valor == 0 {
		exl, cl := grow(i.Fesq, e.Fesq)
		r := event.NovoNodeVazio(e.Valor)
		r.Fesq = exl
		r.Fdir = e.Fdir
		return r, cl + 1
	}
	exl, cl := grow(i.Fesq, e.Fesq)
	exr, cr := grow(i.Fdir, e.Fdir)
	if cl < cr {
		r := event.NovoNodeVazio(e.Valor)
		r.Fesq = exl
		r.Fdir = e.Fdir
		return r, cl + 1
	}
	r := event.NovoNodeVazio(e.Valor)
	r.Fesq = e.Fesq
	r.Fdir = exr
	return r, cr + 1
}

// Trabalhar com pacote e pegar informacao das stamps
func (s *Stamp) Pacote(p *bit.Pacote) {
	s.id.Pacote(p)
	s.event.Pacote(p)
}

func (s *Stamp) Informacao(bits *bit.Informacao) {
	s.id = id.Informacao(bits)
	s.event = event.Informacao(bits)
}

// UnmarshalBinary decodifica a stamp
func (s *Stamp) UnmarshalBinary(data []byte) error {
	bits := bit.NovaInformacao(data)
	s.Informacao(bits)
	return nil
}

// MarshalBinary para codificar a stamp para ser enviada por conexao
func (s *Stamp) MarshalBinary() ([]byte, error) {
	notation := bit.NovoPacote()
	s.Pacote(notation)
	return notation.Pacote(), nil
}