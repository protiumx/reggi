package reggi

type Set[T comparable] map[T]struct{}

func NewSet[T comparable](items ...T) Set[T] {
	set := Set[T](make(map[T]struct{}, len(items)))
	for _, item := range items {
		set.add(item)
	}

	return set
}

func (s *Set[T]) add(t T) {
	(*s)[t] = struct{}{}
}

func (s *Set[T]) remove(t T) {
	delete(*s, t)
}

func (s *Set[T]) has(t T) bool {
	_, ok := (*s)[t]
	return ok
}

func (s *Set[T]) list() []T {
	set := *s
	list := make([]T, 0, len(set))
	for item := range set {
		list = append(list, item)
	}

	return list
}

func (s *Set[T]) size() int {
	return len(*s)
}
