package main

import "fmt"

func merge(ms []map[string]interface{}) map[string]interface{} {
	res := map[string][]interface{}{}
	for _, m := range ms {
		for k, v := range m {
			res[k] = append(res[k], v)
		}
	}
	nameRes := make([]interface{}, 0, len(res["name"]))
	pidRes := make([]interface{}, 0, len(res["pId"]))
	temp := map[interface{}]struct{}{}
	for _, item := range res["name"] {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			nameRes = append(nameRes, item)
		}
	}
	for _, item := range res["pId"] {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			pidRes = append(pidRes, item)
		}
	}

	resMap := map[string]interface{}{"name": nameRes[0], "pId": pidRes[0], "children": res["children"]}
	return resMap
}

func removeDuplicateValues(interfaceList []interface{}) []interface{} {
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

func main() {
	var arrayMap = []map[string]interface{}{
		{"pId": 1, "name": "group", "children": map[string]interface{}{"id": 1, "name": "查看用户组"}},
		{"pId": 1, "name": "group", "children": map[string]interface{}{"id": 2, "name": "添加用户组"}},
		{"pId": 1, "name": "group", "children": map[string]interface{}{"id": 3, "name": "修改用户组"}},
		{"pId": 1, "name": "group", "children": map[string]interface{}{"id": 4, "name": "删除用户组"}},
		{"pId": 2, "name": "permission", "children": map[string]interface{}{"id": 1, "name": "查看权限"}},
		{"pId": 2, "name": "permission", "children": map[string]interface{}{"id": 2, "name": "添加权限"}},
		{"pId": 2, "name": "permission", "children": map[string]interface{}{"id": 3, "name": "修改权限"}},
		{"pId": 2, "name": "permission", "children": map[string]interface{}{"id": 4, "name": "删除权限"}},
	}
	//fmt.Printf("%T\n", merge(arrayMap)["pId"])
	fmt.Println(merge(arrayMap))

	//fmt.Println(res)
	//name := removeDuplicateValues(res["name"])
	//pid := removeDuplicateValues(res["pId"])
	//resMsg := map[string][]interface{}{"name": name, "pId": pid, "children": res["children"]}
	//fmt.Println(resMsg)
}
