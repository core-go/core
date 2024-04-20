package compare

func Contains[T comparable](elems []T, v T) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}

func ListToMap[K comparable, V any](l []V, mapFn func(V) K) map[K]V {
	res := make(map[K]V)
	for _, item := range l {
		res[mapFn(item)] = item
	}
	return res
}

func First[V any](l []V, findFn func(V) bool) (*V, bool) {
	for _, item := range l {
		if findFn(item) {
			return &item, true
		}
	}
	return nil, false
}

func PointerOf[V any](v V) *V {
	a := v
	return &a
}

func Unique[K comparable](l []K) []K {
	visit := make(map[K]bool)
	res := make([]K, len(l), cap(l))
	counter := 0
	for _, item := range l {
		if _, ok := visit[item]; !ok {
			visit[item] = true
			res[counter] = item
			counter++
		}
	}
	return res
}

func Map[I any, O any](l []I, mapFn func(I) O) []O {
	res := make([]O, len(l), cap(l))
	for i := range l {
		res[i] = mapFn(l[i])
	}
	return res
}

func Filter[T any](l []T, filterFn func(T) bool) []T {
	var res = make([]T, 0)
	for _, item := range l {
		if filterFn(item) {
			res = append(res, item)
		}
	}
	return res
}
