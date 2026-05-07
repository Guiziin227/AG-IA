package genetico

import "math/rand"

// GerarPopulacaoInicial cria indivíduos permitindo repetição (genes independentes).
func GerarPopulacaoInicial(r *rand.Rand, tamanho, numCidades int) [][]int {
	pop := make([][]int, tamanho)
	for i := 0; i < tamanho; i++ {
		ind := make([]int, numCidades)
		for j := 0; j < numCidades; j++ {
			ind[j] = r.Intn(numCidades) + 1 // valores 1..numCidades, com repetição
		}
		pop[i] = ind
	}

	return pop
}

func copiar(v []int) []int {
	c := make([]int, len(v))
	copy(c, v)
	return c
}
