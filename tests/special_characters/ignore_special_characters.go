package main

import (
	"fmt"
	"regexp"
	"strings"
)

func main() {
	chars := []string{"]", "^", "\\\\", "[", ".", "(", ")", "-"}
	r := strings.Join(chars, "")
	fmt.Println(r)
	s := "[Some]-\\(string)^."
	re := regexp.MustCompile("[" + r + "]+")
	s = re.ReplaceAllString(s, "")
	fmt.Println(s)

	s1 := "dejun.test"
	s1 = re.ReplaceAllString(s1, " ")
	fmt.Println(s1)
}
