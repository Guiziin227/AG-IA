package main

import (
	"ag/genetico"
	"fmt"
	"math/rand"
	"time"
)

func main() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	resultado := genetico.Executar(
		r,
		9,  // numCidades
		6,  // tamanhoPopulacao
		10, // geracoes
		3,  // torneio
		0.7, // taxaCruzamento
		0.2, // taxaMutacao
	)

	fmt.Println("Melhor rota encontrada:", resultado.Rota)
}
