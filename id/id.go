package id

import (
	"fmt"
	"itc/bit"
	"log"
)

type ID struct {
	Valor       uint32
	Fesq, Fdir *ID
	Folha      bool
}

const (
	zero  = uint32(0)
	um   = uint32(1)
	dois   = uint32(2)
	tres = uint32(3)
)

// Cria nova ID com valor padrao (1)
func New() *ID {
	return NovaID(um)
}
// Cria nova ID com valor
func NovaID(value uint32) *ID {
	return &ID{Valor: value, Folha: true}
}
// Cria nova ID coomo folha
func (i *ID) IDfolha(value uint32) *ID {
	i.Valor = value
	i.Folha = true
	i.Fesq = nil
	i.Fdir = nil
	return i
}

// Cria nova ID como node da arvore
func (i *ID) Node(left, right uint32) *ID {
	return i.IDnode(NovaID(left), NovaID(right))
}

// Popular informacoes de ID node
func (i *ID) IDnode(left, right *ID) *ID {
	i.Valor = 0
	i.Folha = false
	i.Fesq = left
	i.Fdir = right
	return i
}

// Deixar a id na forma normal
func (i *ID) Normal() *ID {
	if i.Folha || !i.Fesq.Folha || !i.Fdir.Folha || i.Fesq.Valor != i.Fdir.Valor {
		return i
	}
	return i.IDfolha(i.Fesq.Valor)
}

// Fazer fork nos eventos
// No fork o evento eh mantido e a id eh separada em duas partes
func (i *ID) Split() (i1, i2 *ID) {

	i1 = New()
	i2 = New()

	if i.Folha && i.Valor == 0 {
		// split(0) = (0, 0)
		i1.Valor = 0
		i2.Valor = 0
		return
	}
	if i.Folha && i.Valor == 1 {
		// split(1) = ((1,0), (0,1))
		i1.Node(um, zero)
		i2.Node(zero, um)
		return
	}
	if (i.Fesq.Folha && i.Fesq.Valor == 0) && (!i.Fdir.Folha || i.Fdir.Valor == 1) {
		// split((0, i)) = ((0, i1), (0, i2)), where (i1, i2) = split(i)
		r1, r2 := i.Fdir.Split()
		i1.IDnode(NovaID(zero), r1)
		i2.IDnode(NovaID(zero), r2)
		return
	}
	if (!i.Fesq.Folha || i.Fesq.Valor == 1) && (i.Fdir.Folha && i.Fdir.Valor == 0) {
		// split((i, 0)) = ((i1, 0), (i2, 0)), where (i1, i2) = split(i)
		l1, l2 := i.Fesq.Split()
		i1.IDnode(l1, NovaID(zero))
		i2.IDnode(l2, NovaID(zero))
		return
	}
	if (!i.Fesq.Folha || i.Fesq.Valor == 1) && (!i.Fdir.Folha || i.Fdir.Valor == 1) {
		// split((i1, i2)) = ((i1, 0), (0, i2))
		i1.IDnode(i.Fesq, NovaID(zero))
		i2.IDnode(NovaID(zero), i.Fdir)
		return
	}
	log.Fatalf("Nao consegui dividir ID: %s", i.String())
	return
}

// Escreve ID
func (i *ID) String() string {
	if i.Folha {
		return fmt.Sprintf("%d", i.Valor)
	}
	return fmt.Sprintf("(%s, %s)", i.Fesq, i.Fdir)
}

// Soma o intervalo de duas ids
func (i *ID) Soma(i1, i2 *ID) *ID {
	if i1.Folha && i1.Valor == 0 {
		i.Valor = i2.Valor
		i.Fesq = i2.Fesq
		i.Fdir = i2.Fdir
		i.Folha = i2.Folha
		return i
	}
	if i2.Folha && i2.Valor == 0 {
		i.Valor = i1.Valor
		i.Fesq = i1.Fesq
		i.Fdir = i1.Fdir
		i.Folha = i1.Folha
		return i
	}
	return i.IDnode(New().Soma(i1.Fesq, i2.Fesq), New().Soma(i1.Fdir, i2.Fdir)).Normal()
}

// Trabalha com o pacote de ID
func (i *ID) Pacote(notation *bit.Pacote) {
	if i.Folha {
		notation.Push(zero, dois)
		notation.Push(i.Valor, um)
		return
	}
	if i.Fesq.Folha && i.Fesq.Valor == 0 {
		notation.Push(um, dois)
		i.Fdir.Pacote(notation)
		return
	}
	if i.Fdir.Folha && i.Fdir.Valor == 0 {
		notation.Push(dois, dois)
		i.Fesq.Pacote(notation)
		return
	}
	notation.Push(tres, dois)
	i.Fesq.Pacote(notation)
	i.Fdir.Pacote(notation)
	return
}

// Retira os dados de ID do pacote
func Informacao(bits *bit.Informacao) *ID {
	i := New()
	switch bits.Pop(dois) {
	case 0:
		i.IDfolha(bits.Pop(um))
	case 1:
		newID := Informacao(bits)
		i.IDnode(NovaID(zero), newID)
	case 2:
		newID := Informacao(bits)
		i.IDnode(newID, NovaID(zero))
	case 3:
		newLeft := Informacao(bits)
		newRight := Informacao(bits)
		i.IDnode(newLeft, newRight)
	}
	return i
}
