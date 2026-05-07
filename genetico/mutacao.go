package genetico

import "math/rand"

// Mutacao substitui um gene por um valor aleatório (1..numCidades).
func Mutacao(ind []int, r *rand.Rand, numCidades int) (bool, int, int) {
	if len(ind) == 0 {
		return false, -1, -1
	}
	idx := r.Intn(len(ind))
	old := ind[idx]
	ind[idx] = r.Intn(numCidades) + 1
	return true, idx, old
}
