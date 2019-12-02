package itc

import (
	"fmt"
)

func Example() {
	fmt.Printf("Exemplo da apresentacao\n")
	fmt.Printf("Criacao do seed a\n")
	a := NewStamp()
	fmt.Printf("   a: %s\n", a)
	fmt.Printf("Fork de a\n")
	b := a.Fork()
	fmt.Printf("   b: %s\n", a)
	fmt.Printf("Evento em b\n")
	a.Event()
	fmt.Printf("   c: %s\n", a)
	b.Event()
	b.Event()
	fmt.Printf("Apos dois eventos no filho direito de a\n")
	fmt.Printf("   d: %s\n", b)
	c := a.Fork()
	a.Event()
	b.Join(c)
	fmt.Printf("Apos fork em c e join de d com o filho direito de c\n")
	fmt.Printf("   e: %s\n", b)
	c = b.Fork()
	fmt.Printf("Fork de e\n")
	fmt.Printf("   f: %s\n", c)
	fmt.Printf("Join entre c após evento e filho esquerdo de e\n")
	a.Join(b)
	fmt.Printf("   g: %s\n", a)
	fmt.Printf("Evento em g\n")
	a.Event()
	fmt.Printf("   h: %s\n", a)
	// Output:
	// Exemplo da apresentacao
	// Criacao do seed a
	//    a: (1, 0)
	// Fork de a
	//    b: ((1, 0), 0)
	// Evento em b
	//    c: ((1, 0), (0, 1, 0))
	// Apos dois eventos no filho direito de a
	//    d: ((0, 1), (0, 0, 2))
	// Apos fork em c e join de d com o filho direito de c
	//    e: (((0, 1), 1), (1, 0, 1))
	// Fork de e
	//    f: ((0, 1), (1, 0, 1))
	// Join entre c após evento e filho esquerdo de e
	//    g: ((1, 0), (1, (0, 1, 0), 1))
	// Evento em g
	//    h: ((1, 0), 2)
}

func CriarSeed() {
	a := NewStamp()
	fmt.Printf("Criacao de seed a\n")
	fmt.Printf("a: %s\n", a)
	// Output:
	// Criacao de seed a
	// a: (1, 0)
}

func TesteFork() {
	a := NewStamp()
	b := a.Fork()
	fmt.Printf("Fork em seed a\n")
	fmt.Printf("a: %s\n", a)
	fmt.Printf("b: %s\n", b)

	// Output:
	// Fork em seed a
	// a: ((1, 0), 0)
	// b: ((0, 1), 0)
}

func TesteForkEventJoin() {
	a := NewStamp()
	b := a.Fork()
	a.Event()
	b.Event()
	c := a.Fork()
	b.Event()
	a.Event()
	b.Join(c)
	b.Event()
	fmt.Printf("Teste de fork e join\n")
	fmt.Printf("a: %s\n", a)
	fmt.Printf("b: %s\n", b)

	// Output:
	// Teste de fork e join
	// a: (((1, 0), 0), (0, (1, 1, 0), 0))
	// b: (((0, 1), 1), (1, 0, 2))
}
