package util

func TernaryOperator[A any](statement bool, valueIfTrue A, valueIfFalse A) A {
	if statement {
		return valueIfTrue
	}

	return valueIfFalse
}

func NewPointer[A any](a A) *A {
	return &a
}

func UnwrapPointer[A any](a *A) A {
	if a == nil {
		return *new(A)
	}

	return *a
}

func Map[A any, B any](f func(a A) B, amap []A) []B {
	bmap := make([]B, len(amap))
	for k, v := range amap {
		bmap[k] = f(v)
	}
	return bmap
}

func MapPointer[A any, B any](f func(a A) B, pointer *A) *B {
	if pointer == nil {
		return nil
	}

	result := f(*pointer)
	return &result
}

func Member[A comparable](a A, amap []A) bool {
	for _, v := range amap {
		if v == a {
			return true
		}
	}
	return false
}

func Deduplicate[A comparable](alist []A) []A {
	allKeys := make(map[A]bool)
	list := []A{}
	for _, item := range alist {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func GetDuplicateItems[A comparable](alist []A) []A {
	allKeys := make(map[A]bool)
	list := []A{}
	for _, item := range alist {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			continue
		}

		list = append(list, item)
	}
	return list
}

func Filter[A any](f func(a A) bool, amap []A) []A {
	res := make([]A, 0)
	for _, v := range amap {
		if f(v) {
			res = append(res, v)
		}
	}
	return res
}
