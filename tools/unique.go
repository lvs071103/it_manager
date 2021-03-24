package tools

func Merge(ms []map[string]interface{}) map[string][]interface{} {
	res := map[string][]interface{}{}
	for _, m := range ms {
		for k, v := range m {
			res[k] = append(res[k], v)
		}
	}
	return res
}

func RemoveDuplicateValues(interfaceList []interface{}) []interface{} {
	result := make([]interface{}, 0, len(interfaceList))
	temp := map[interface{}]struct{}{}
	for _, item := range interfaceList {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}
