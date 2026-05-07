package fitness

// Avaliar calcula a penalidade de uma rota.
// Quanto menor o valor, melhor a rota.
func Avaliar(rota []int) int {
	penalidade := 0

	// Regra 1: se uma cidade maior vier antes de uma menor, soma 10.
	for i := 0; i < len(rota)-1; i++ {
		if rota[i] > rota[i+1] {
			penalidade += 10
		}
	}

	// Regra 2: para cada par de ocorrências iguais, soma 20.
	for i := 0; i < len(rota); i++ {
		for j := i + 1; j < len(rota); j++ {
			if rota[i] == rota[j] {
				penalidade += 20
			}
		}
	}

	return penalidade
}
