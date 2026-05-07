package genetico

import (
	"ag/fitness"
	"fmt"
	"math/rand"
)

// Resultado guarda a melhor rota e sua penalidade.
type Resultado struct {
	Rota    []int
	Fitness int
}

func indicesPopulacao(tamanho int) []int {
	indices := make([]int, tamanho)
	for i := 0; i < tamanho; i++ {
		indices[i] = i
	}
	return indices
}

func candidatosDisponiveis(populacao [][]int, usados map[string]bool) []int {
	candidatos := make([]int, 0, len(populacao))
	for i, ind := range populacao {
		if !usados[ChaveIndividuo(ind)] {
			candidatos = append(candidatos, i)
		}
	}
	return candidatos
}

// Executar roda o AG e, quando verbose for true, imprime o processo completo.
func Executar(r *rand.Rand, numCidades, tamanhoPopulacao, geracoes, torneioK int, taxaCruzamento, taxaMutacao float64) Resultado {
	populacao := GerarPopulacaoInicial(r, tamanhoPopulacao, numCidades)

	melhor := copiar(populacao[0])
	melhorFitness := fitness.Avaliar(melhor)
	geracaoMelhor := 0

	fmt.Println("=== POPULAÇÃO INICIAL ===")
	for i, ind := range populacao {
		fmt.Printf("Indivíduo %02d -> %v | Fitness: %d\n", i+1, ind, fitness.Avaliar(ind))
	}

	for g := 0; g < geracoes; g++ {

		fmt.Printf("=== GERAÇÃO %d ===\n", g+1)

		novaPop := make([][]int, 0, tamanhoPopulacao)
		filhoNum := 1
		usadosNaGeracao := make(map[string]bool)

		for len(novaPop) < tamanhoPopulacao {
			candidatosPai := candidatosDisponiveis(populacao, usadosNaGeracao)
			if len(candidatosPai) == 0 {
				usadosNaGeracao = make(map[string]bool)
				candidatosPai = indicesPopulacao(len(populacao))
			}
			pai := TorneioComPool(populacao, candidatosPai, torneioK, r)
			usadosNaGeracao[ChaveIndividuo(pai.Individuo)] = true

			candidatosMae := candidatosDisponiveis(populacao, usadosNaGeracao)
			if len(candidatosMae) == 0 {
				usadosNaGeracao = make(map[string]bool)
				usadosNaGeracao[ChaveIndividuo(pai.Individuo)] = true
				candidatosMae = candidatosDisponiveis(populacao, usadosNaGeracao)
				if len(candidatosMae) == 0 {
					candidatosMae = indicesPopulacao(len(populacao))
				}
			}
			mae := TorneioComPool(populacao, candidatosMae, torneioK, r)
			usadosNaGeracao[ChaveIndividuo(mae.Individuo)] = true
			var filho []int
			var ponto int

			if r.Float64() < taxaCruzamento {
				filho, ponto = CrossoverSinglePoint(pai.Individuo, mae.Individuo, r)
			} else {
				filho = copiar(pai.Individuo)
				ponto = -1
			}

			fmt.Printf("Filho %02d\n", filhoNum)
			fmt.Printf("  Pai  -> %v | Fitness: %d\n", pai.Individuo, pai.Fitness)
			fmt.Printf("  Mãe  -> %v | Fitness: %d\n", mae.Individuo, mae.Fitness)
			if ponto >= 0 {
				fmt.Printf("  Crossover single-point -> corte %d\n", ponto)
			} else {
				fmt.Println("  Sem crossover -> cópia do pai")
			}
			fmt.Printf("  Filho antes da mutação -> %v\n", filho)

			mutou := false
			iMut, old := -1, -1
			if r.Float64() < taxaMutacao {
				mutou, iMut, old = MutacaoRandomReset(filho, r, numCidades)
			}

			if mutou {
				fmt.Printf("  Mutação aplicada -> posição %d (antes=%d)\n", iMut, old)
				fmt.Printf("  Filho após mutação -> %v\n", filho)
			} else {
				fmt.Println("  Mutação não aplicada")
			}
			fmt.Printf("  Fitness do filho -> %d\n", fitness.Avaliar(filho))
			fmt.Println()

			novaPop = append(novaPop, filho)
			filhoNum++
		}

		populacao = novaPop
		melhorDaGeracao := copiar(populacao[0])
		fitnessDaGeracao := fitness.Avaliar(melhorDaGeracao)

		for _, ind := range populacao {
			f := fitness.Avaliar(ind)
			if f < fitnessDaGeracao {
				fitnessDaGeracao = f
				melhorDaGeracao = copiar(ind)
			}
			if f < melhorFitness {
				melhorFitness = f
				melhor = copiar(ind)
				geracaoMelhor = g + 1
			}
		}

		fmt.Printf("Melhor da geração %d -> %v | Fitness: %d\n", g+1, melhorDaGeracao, fitnessDaGeracao)
		fmt.Printf("Melhor global até agora -> %v | Fitness: %d | Encontrado na geração %d\n", melhor, melhorFitness, geracaoMelhor)
		fmt.Println()

		if melhorFitness == 0 {
			fmt.Println("Fitness 0 atingido. Melhor rota encontrada com sucesso.")
			break
		}
	}

	fmt.Println("FINAL:")
	fmt.Printf("Melhor rota encontrada: %v\n", melhor)
	fmt.Printf("Fitness: %d\n", melhorFitness)

	return Resultado{Rota: melhor, Fitness: melhorFitness}
}
