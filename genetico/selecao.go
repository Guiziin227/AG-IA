package genetico

import (
	"ag/fitness"
	"fmt"
	"math/rand"
)

// TorneioResultado guarda o indivíduo escolhido e sua aptidão.
type TorneioResultado struct {
	Individuo []int
	Fitness   int
}

// ChaveIndividuo converte uma rota em uma chave estável para mapa.
func ChaveIndividuo(individuo []int) string {
	return fmt.Sprint(individuo)
}

// TorneioComPool seleciona o melhor indivíduo apenas entre os candidatos informados.
func TorneioComPool(populacao [][]int, candidatos []int, torneioK int, r *rand.Rand) TorneioResultado {
	if len(populacao) == 0 {
		return TorneioResultado{}
	}
	if len(candidatos) == 0 {
		candidatos = make([]int, len(populacao))
		for i := range populacao {
			candidatos[i] = i
		}
	}
	if torneioK < 1 {
		torneioK = 1
	}
	if torneioK > len(candidatos) {
		torneioK = len(candidatos)
	}

	perm := r.Perm(len(candidatos))

	melhorIdx := candidatos[perm[0]]
	melhor := populacao[melhorIdx]
	melhorFitness := fitness.Avaliar(melhor)

	fmt.Println("=== TORNEIO ===")
	fmt.Printf("Candidato 1 -> idx %d %v | Fitness: %d\n", melhorIdx, melhor, melhorFitness)

	for i := 1; i < torneioK; i++ {
		idx := candidatos[perm[i]]
		candidato := populacao[idx]
		f := fitness.Avaliar(candidato)
		fmt.Printf("Candidato %d -> idx %d %v | Fitness: %d\n", i+1, idx, candidato, f)

		if f < melhorFitness {
			melhor = candidato
			melhorFitness = f
			melhorIdx = idx
		}
	}

	fmt.Printf("Vencedor do torneio -> idx %d %v | Fitness: %d\n", melhorIdx, melhor, melhorFitness)
	fmt.Println()

	return TorneioResultado{
		Individuo: melhor,
		Fitness:   melhorFitness,
	}
}
