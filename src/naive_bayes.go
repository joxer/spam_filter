package main

import (
	"fmt"
	"github.com/jdkato/prose/tokenize"

	"./bayes"
	"./common"
)

func main() {

	data, _ := common.OpenFile("/home/joxer/code/go/spam_filter/data/dataset.csv")
	testset, _ := common.OpenFile("/home/joxer/code/go/spam_filter/data/testset.csv")
	spam, ham := common.GetNames(data[0:200])
	fmt.Println(spam)
	fmt.Println(ham)
	classifier := bayes.NewClassifier("Good", "Bad")
	classifier.Learn(ham, "Good")
	classifier.Learn(spam, "Bad")
	tokenizer := tokenize.NewTreebankWordTokenizer()
	correct := 0
	incorrect := 0
	for _, line := range (testset) {
		words := (tokenizer.Tokenize(line[1]))
		_ , class, _,_  := classifier.SafeProbScores(words)
		if((line[0] == "ham") && (0 == class)) {
			correct += 1
		} else if ((line[0] == "spam") && 1 == class) {
			correct += 1
		} else {
			incorrect += 1
		}

	}
	fmt.Println(correct)
	fmt.Println(incorrect)

}

