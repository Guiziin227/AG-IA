# AG de Roteamento em Go

Refatoração do AG em packages separados para resolver o problema de roteamento com 9 cidades.

## Estrutura

- `main.go`: ponto de entrada
- `modelo/`: tipos compartilhados
- `fitness/`: cálculo da aptidão
- `genetico/`: população, seleção por torneio, crossover, mutação e execução do AG

## Como rodar

```powershell
go run .
```

O `main.go` está configurado para executar em modo detalhado, imprimindo seleção por torneio, crossover, mutação e fitness por geração até encontrar a melhor rota.

## Regra de aptidão

- `+10` quando uma cidade maior aparece antes de uma menor
- `+20` para cada par de ocorrências repetidas
- `+20` para cidades fora do intervalo `1..9`

A melhor solução possível é `1 2 3 4 5 6 7 8 9`, com fitness `0`.

