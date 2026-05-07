package genetico

import (
	"math/rand"
)

// Crossover aplica crossover de ponto unico aleatorio.
// Filho = prefixo do pai + sufixo da mae.
func Crossover(pai, mae []int, r *rand.Rand) ([]int, int) {
	n := len(pai)
	filho := make([]int, n)
	if n == 0 {
		return filho, 0
	}

	ponto := r.Intn(n-1) + 1

	// copia prefixo do pai
	for i := 0; i < ponto; i++ {
		filho[i] = pai[i]
	}
	// copia sufixo da mae
	for i := ponto; i < n; i++ {
		filho[i] = mae[i]
	}

	return filho, ponto
}

// CrossoverMetadeFixa aplica crossover com corte fixo na metade.
// Filho = primeira metade do pai + segunda metade da mae.
func CrossoverMetadeFixa(pai, mae []int) ([]int, int) {
	n := len(pai)
	filho := make([]int, n)
	if n == 0 {
		return filho, 0
	}

	ponto := n / 2

	for i := 0; i < ponto; i++ {
		filho[i] = pai[i]
	}
	for i := ponto; i < n; i++ {
		filho[i] = mae[i]
	}

	return filho, ponto
}
