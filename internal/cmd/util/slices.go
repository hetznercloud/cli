package util

func Batches[T any](all []T, size int) (batches [][]T) {
	for size < len(all) {
		all, batches = all[size:], append(batches, all[:size])
	}
	return append(batches, all)
}
