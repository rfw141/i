package slicex

func StrRemove(arr []string, elem string) []string {
	for i, e := range arr {
		if e == elem {
			return append(arr[:i], arr[i+1:]...)
		}
	}
	return arr
}

func StrContains(arr []string, elem string) bool {
	for _, e := range arr {
		if e == elem {
			return true
		}
	}
	return false
}

func StrUniq(arr []string) []string {
	uniqueMap := make(map[string]bool)
	var uniqueArr []string

	for _, val := range arr {
		if _, ok := uniqueMap[val]; !ok {
			uniqueMap[val] = true
			uniqueArr = append(uniqueArr, val)
		}
	}

	return uniqueArr
}
