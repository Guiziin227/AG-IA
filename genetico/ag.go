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

	fmt.Println("POPULAÇÃO INICIAL")
	for i, ind := range populacao {
		fmt.Printf("Indivíduo %02d -> %v | Fitness: %d\n", i+1, ind, fitness.Avaliar(ind))
	}

	for g := 0; g < geracoes; g++ {

		fmt.Printf("GERAÇÃO %d\n", g+1)

		novaPop := make([][]int, 0, tamanhoPopulacao)
		filhoNum := 1
		usadosNaGeracao := make(map[string]bool)

		// Produz dois filhos por par de pais: seleciona os pais uma vez por par
		for len(novaPop) < tamanhoPopulacao {
			candidatosPai := candidatosDisponiveis(populacao, usadosNaGeracao)
			if len(candidatosPai) == 0 {
				usadosNaGeracao = make(map[string]bool)
				candidatosPai = indicesPopulacao(len(populacao))
			} else if len(candidatosPai) < torneioK {
				// reabastecer usadosNaGeracao quando não houver candidatos suficientes
				usadosNaGeracao = make(map[string]bool)
				candidatosPai = candidatosDisponiveis(populacao, usadosNaGeracao)
				if len(candidatosPai) == 0 {
					candidatosPai = indicesPopulacao(len(populacao))
				}
			}

			var paiRes, maeRes TorneioResultado
			if len(candidatosPai) >= 2 {
				// amostrar uma vez e usar a mesma amostra para escolher pai e mãe
				// usar K mínimo = 3 para amostragem, mas não maior que candidatos disponíveis
				k := torneioK
				if k < 3 {
					k = 3
				}
				if k > len(candidatosPai) {
					k = len(candidatosPai)
				}
				perm := r.Perm(len(candidatosPai))
				amostra := make([]int, k)
				for i := 0; i < k; i++ {
					amostra[i] = candidatosPai[perm[i]]
				}

				// primeiro torneio para escolher o pai
				p := Torneio(populacao, amostra, k, r)
				paiRes = p
				usadosNaGeracao[ChaveIndividuo(p.Individuo)] = true

				// preparar amostra para a mãe removendo o índice do pai
				amostraMae := make([]int, 0, k)
				paiKey := ChaveIndividuo(p.Individuo)
				// primeiro, adicione os que já estavam na amostra original (exceto o pai)
				for _, idx := range amostra {
					if ChaveIndividuo(populacao[idx]) != paiKey {
						amostraMae = append(amostraMae, idx)
					}
				}
				// se ainda faltar, use os elementos restantes da permutação para completar até k
				if len(amostraMae) < k {
					for i := k; i < len(candidatosPai) && len(amostraMae) < k; i++ {
						idx := candidatosPai[perm[i]]
						if ChaveIndividuo(populacao[idx]) == paiKey {
							continue
						}
						// evitar duplicatas
						dup := false
						for _, v := range amostraMae {
							if v == idx {
								dup = true
								break
							}
						}
						if !dup {
							amostraMae = append(amostraMae, idx)
						}
					}
				}

				if len(amostraMae) == 0 {
					// se nada sobrou, escolha a partir dos disponíveis (exceto o pai)
					candidatosMae := candidatosDisponiveis(populacao, usadosNaGeracao)
					if len(candidatosMae) == 0 {
						candidatosMae = indicesPopulacao(len(populacao))
					}
					m := Torneio(populacao, candidatosMae, k, r)
					maeRes = m
					usadosNaGeracao[ChaveIndividuo(m.Individuo)] = true
				} else {
					// use k when possible; Torneio will clamp if amostraMae shorter
					m := Torneio(populacao, amostraMae, k, r)
					maeRes = m
					usadosNaGeracao[ChaveIndividuo(m.Individuo)] = true
				}
			} else {
				p := Torneio(populacao, candidatosPai, torneioK, r)
				usadosNaGeracao[ChaveIndividuo(p.Individuo)] = true
				candidatosMae := candidatosDisponiveis(populacao, usadosNaGeracao)
				if len(candidatosMae) == 0 {
					usadosNaGeracao = make(map[string]bool)
					usadosNaGeracao[ChaveIndividuo(p.Individuo)] = true
					candidatosMae = candidatosDisponiveis(populacao, usadosNaGeracao)
					if len(candidatosMae) == 0 {
						candidatosMae = indicesPopulacao(len(populacao))
					}
				}
				m := Torneio(populacao, candidatosMae, torneioK, r)
				usadosNaGeracao[ChaveIndividuo(m.Individuo)] = true
				paiRes = p
				maeRes = m
			}

			// Gerar até dois filhos a partir deste par
			for k := 0; k < 2 && len(novaPop) < tamanhoPopulacao; k++ {
				var filho []int
				var ponto int

				if r.Float64() < taxaCruzamento {
					filho, ponto = Crossover(paiRes.Individuo, maeRes.Individuo, r)
				} else {
					filho = copiar(paiRes.Individuo)
					ponto = -1
				}

				fmt.Printf("Filho %02d\n", filhoNum)
				fmt.Printf("  Pai  -> %v | Fitness: %d\n", paiRes.Individuo, paiRes.Fitness)
				fmt.Printf("  Mãe  -> %v | Fitness: %d\n", maeRes.Individuo, maeRes.Fitness)
				if ponto >= 0 {
					fmt.Printf("  Crossover single-point -> corte %d\n", ponto)
				} else {
					fmt.Println("  Sem crossover -> cópia do pai")
				}
				fmt.Printf("  Filho antes da mutação -> %v\n", filho)

				mutou := false
				iMut, old := -1, -1
				if r.Float64() < taxaMutacao {
					mutou, iMut, old = Mutacao(filho, r, numCidades)
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
