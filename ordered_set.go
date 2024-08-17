package reggi

type OrderedSet[T comparable] struct {
	set  map[T]int
	next int
}

func (o *OrderedSet[T]) add(values ...T) {
	if o.set == nil {
		o.set = make(map[T]int)
	}

	for _, v := range values {
		if !o.has(v) {
			o.set[v] = o.next
			o.next++
		}
	}
}

func (o *OrderedSet[T]) has(v T) bool {
	_, ok := o.set[v]
	return ok
}

func (o *OrderedSet[T]) list() []T {
	ret := make([]T, len(o.set))
	for k, v := range o.set {
		ret[v] = k
	}

	return ret
}

func (o *OrderedSet[T]) index(v T) int {
	return o.set[v]
}
