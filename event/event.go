package event

import (
	"fmt"
	"itc/bit"
)

// Pacote responsavel pelas operacoes envolvendo eventos
type Event struct {
	Valor       uint32
	Fesq, Fdir *Event
	Folha      bool
}

const (
	zero  = uint32(0)
	um   = uint32(1)
	dois   = uint32(2)
	tres = uint32(3)
)

// Cria novo evento
func New() *Event {
	return NovaFolha(zero)
}

func NovaFolha(value uint32) *Event {
	return &Event{Valor: value, Folha: true}
}

// Cria um node  vazio
func NovoNodeVazio(value uint32) *Event {
	return &Event{Valor: value, Folha: false, Fesq: New(), Fdir: New()}
}

// Cria um no
func NovoNode(value, left, right uint32) *Event {
	return &Event{Valor: value, Folha: false, Fesq: NovaFolha(left), Fdir: NovaFolha(right)}
}

// Clona um evento, eh necessario em fork para clonar a parte causal de uma stamp
func (e *Event) Clona() *Event {
	clone := New()
	clone.Folha = e.Folha
	clone.Valor = e.Valor
	if e.Fesq != nil {
		clone.Fesq = e.Fesq.Clona()
	}
	if e.Fdir != nil {
		clone.Fdir = e.Fdir.Clona()
	}
	return clone
}

///////// Operacoes entre eventos
// Verifica se dois eventos sao iguais
func (e *Event) Equals(o *Event) bool {
	return (e == nil && o == nil) ||
		((e.Folha == o.Folha) &&
			(e.Valor == o.Valor) &&
			e.Fesq.Equals(o.Fesq) &&
			e.Fdir.Equals(o.Fdir))
}

// Coloca o evento na forma normal
func (e *Event) Normal() *Event {
	if e.Folha {
		return e
	}
	if e.Fesq.Folha && e.Fdir.Folha && e.Fesq.Valor == e.Fdir.Valor {
		return NovaFolha(e.Valor + e.Fesq.Valor)
	}
	m := Menor(e.Fesq.Min(), e.Fdir.Min())
	e.Fesq = e.Fesq.Normal().sink(m)
	e.Fdir = e.Fdir.Normal().sink(m)
	return e.lift(m)
}

// Realiza a operacao lift
func (e *Event) lift(m uint32) *Event {
	result := e.Clona()
	result.Valor += m
	return result
}

// Realiza a operacao sink
func (e *Event) sink(m uint32) *Event {
	result := e.Clona()
	result.Valor -= m
	return result
}

// Pega o valor maximo do evento
func (e *Event) Max() uint32 {
	if e.Folha {
		return e.Valor
	}
	return e.Valor + Maior(e.Fesq.Max(), e.Fdir.Max())
}

// Pega o valor minimo do evento
func (e *Event) Min() uint32 {
	if e.Folha {
		return e.Valor
	}
	return e.Valor + Menor(e.Fesq.Min(), e.Fdir.Min())
}

// Escreve o evento
func (e *Event) String() string {
	if e.Folha {
		return fmt.Sprintf("%d", e.Valor)
	}
	return fmt.Sprintf("(%d, %s, %s)", e.Valor, e.Fesq, e.Fdir)
}

// Junta duas entidades somando o id e fazendo join dos eventos
func Join(e1, e2 *Event) *Event {
	if e1.Folha && e2.Folha {
		return NovaFolha(Maior(e1.Valor, e2.Valor))
	}
	if e1.Folha {
		return Join(NovoNodeVazio(e1.Valor), e2)
	}
	if e2.Folha {
		return Join(e1, NovoNodeVazio(e2.Valor))
	}
	if e1.Valor > e2.Valor {
		return Join(e2, e1)
	}
	e := NovoNodeVazio(e1.Valor)
	e.Fesq = Join(e1.Fesq, e2.Fesq.lift(e2.Valor-e1.Valor))
	e.Fdir = Join(e1.Fdir, e2.Fdir.lift(e2.Valor-e1.Valor))
	return e.Normal()
}

// Verifica se e1 tem valor menor ou igual a e2
func LEQ(e1, e2 *Event) bool {
	if e1.Folha {
		return e1.Valor <= e2.Valor
	}
	if e2.Folha {
		return (e1.Valor <= e2.Valor) &&
			LEQ(e1.Fesq.lift(e1.Valor), e2) &&
			LEQ(e1.Fdir.lift(e1.Valor), e2)
	}
	return (e1.Valor <= e2.Valor) &&
		LEQ(e1.Fesq.lift(e1.Valor), e2.Fesq.lift(e2.Valor)) &&
		LEQ(e1.Fdir.lift(e1.Valor), e2.Fdir.lift(e2.Valor))
}

// Pega o valor maximo entre dois inteiros
func Maior(m1, m2 uint32) uint32 {
	if m1 > m2 {
		return m1
	}
	return m2
}
// Pega o valor minimo entre dois inteiros
func Menor(m1, m2 uint32) uint32 {
	if m1 < m2 {
		return m1
	}
	return m2
}

/// Insere os dados do evento no pacote
func (e Event) Pacote(notation *bit.Pacote) {
	if e.Folha {
		notation.Push(um, um)
		bit.Encode(uint32(e.Valor), dois, notation)
		return
	}

	notation.Push(zero, um)
	if e.Valor == 0 {
		if e.Fesq.Folha && e.Fesq.Valor == 0 {
			notation.Push(zero, dois)
			e.Fdir.Pacote(notation)
			return
		}
		if e.Fdir.Folha && e.Fdir.Valor == 0 {
			notation.Push(um, dois)
			e.Fesq.Pacote(notation)
			return
		}
		notation.Push(dois, dois)
		e.Fesq.Pacote(notation)
		e.Fdir.Pacote(notation)
		return
	}

	notation.Push(tres, dois)
	if e.Fesq.Folha && e.Fesq.Valor == 0 {
		notation.Push(zero, um)
		notation.Push(zero, um)
		NovaFolha(e.Valor).Pacote(notation)
		e.Fdir.Pacote(notation)
		return
	}
	if e.Fdir.Folha && e.Fdir.Valor == 0 {
		notation.Push(zero, um)
		notation.Push(um, um)
		NovaFolha(e.Valor).Pacote(notation)
		e.Fesq.Pacote(notation)
		return
	}
	notation.Push(um, um)
	NovaFolha(e.Valor).Pacote(notation)
	e.Fesq.Pacote(notation)
	e.Fdir.Pacote(notation)
	return
}

/// Retira os dados do pacote
func Informacao(bits *bit.Informacao) *Event {
	if bits.Pop(um) == 1 {
		return NovaFolha(bit.Decode(dois, bits))
	}
	e := NovoNode(zero, zero, zero)
	switch bits.Pop(dois) {
	case 0:
		e.Fdir = Informacao(bits)
	case 1:
		e.Fesq = Informacao(bits)
	case 2:
		e.Fesq = Informacao(bits)
		e.Fdir = Informacao(bits)
	case 3:
		if bits.Pop(um) == 0 {
			if bits.Pop(um) == 0 {
				bits.Pop(um)
				e.Valor = bit.Decode(dois, bits)
				e.Fdir = Informacao(bits)
			} else {
				bits.Pop(um)
				e.Valor = bit.Decode(dois, bits)
				e.Fesq = Informacao(bits)
			}
		} else {
			bits.Pop(um)
			e.Valor = bit.Decode(dois, bits)
			e.Fesq = Informacao(bits)
			e.Fdir = Informacao(bits)
		}
	}
	return e
}
