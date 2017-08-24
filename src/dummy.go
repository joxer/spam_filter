package main

import (
	"fmt"
	"./common"
	"github.com/jdkato/prose/tokenize"
	"strings"
)
func main() {

		testset, _ := common.OpenFile("/home/joxer/code/go/spam_filter/data/testset.csv")

	tokenizer := tokenize.NewTreebankWordTokenizer()
	correct := 0
	incorrect := 0
	for _, line := range (testset) {
		words := (tokenizer.Tokenize(line[1]))
		class := get_class(words)
		fmt.Println(class,line)
		if ((line[0] == "ham") && ("ham" == class)) {
			correct += 1
		} else if ((line[0] == "spam") && "spam" == class) {
			correct += 1
		} else {
			incorrect += 1
		}

	}
	fmt.Println(correct)
	fmt.Println(incorrect)
}

func get_class(data []string) string{

	ret := "ham"
	for _, v:= range(data){
		vl := strings.ToLower(v)
		vl = strings.Replace(vl,".","",-1)
		vl = strings.Replace(vl,",","",-1)
		vl = strings.Replace(vl,"'","",-1)
		if  vl == "sex" || vl == "credit" || vl == "free" ||
			vl == "call" || vl == "love" || vl =="kiss" ||
			vl == "phone" || vl == "time" || vl == "sir" ||
			vl == "box" || vl == "cellphone" || vl == "home"{
			ret = "spam"
			break
		}
	}
	return ret

}