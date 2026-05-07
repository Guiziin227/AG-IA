package modelo

// Rota representa uma solução candidata do problema de roteamento.
type Rota []int

// Config agrupa os parâmetros do AG.
type Config struct {
	NumCidades  int
	Populacao   int
	Geracoes    int
	TorneioK    int
	TaxaMutacao float64
}

