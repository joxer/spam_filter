package common

import (
	"encoding/csv"
	"strings"
	"sort"
	"os"
	"math"
	"github.com/jdkato/prose/tokenize"
	"github.com/jdkato/prose/tag"
	"fmt"
)

const (
	SPAM                = "spam"
	HAM                 = "ham"
	MOST_FREQUENT_WORDS = 300
)

func RemoveDuplicates(class_1 []string, class_2 []string){
	ret := make([]string, len(class_1)+len(class_2))
	for i := 0; i < len(class_1); i++ {
		ret[i] = class_1[i]
	}
	for i := len(class_1); i < len(class_1)+len(class_2); i++ {
		ret[i] = class_2[i-len(class_1)]
	}
	sort.Strings(ret)
	unique := make([]string, int(math.Max(float64(len(class_1)), float64(len(class_2)))))
	j := 0
	for i:= 0 ; i  < len(ret)-1;i+=1 {
		if(ret[i] == ret[i+1]){
			if !Contains(unique, ret[i]) {
				unique[j] = ret[i]
				j+=1
			}
		}
	}

	for _,value := range(unique) {
		if value != "" {
			i := ContainsWithIdx(class_1, value)
			if i != -1 && i != len(class_1){
				Remove(class_1, value)
			}


			i = ContainsWithIdx(class_2, value)
			if i != -1 && i != len(class_2){
				Remove(class_2, value)
			}

		}
	}

}


func GetNames(data [][]string) ([]string, []string) {

	words_map := make(map[string]int)
	words_class := make(map[string][]string)
	words_class["ham"] = make([]string, 50*len(data))
	words_class["spam"] = make([]string, 50*len(data))
	for idx_line, data_line := range (data) {
		class := data_line[0]
		line := data_line[1]

		words := tokenize.NewTreebankWordTokenizer().Tokenize(line)
		tagger := tag.NewPerceptronTagger()
		for _, tok := range tagger.Tag(words) {
			if strings.Contains(tok.Tag, "NN") {
				key := string(tok.Text)
				if _, ok := words_map[key]; !ok {
					words_map[key] = 1
				} else {
					words_map[key] += 1
				}
				words_class[class] = append(words_class[class], key)

			}
		}
		fmt.Println("doing line: ", idx_line, ", remaning : ", len(data)-idx_line)

	}


	spam_words := GetMostFrequent(words_map, words_class["spam"])
	ham_words := GetMostFrequent(words_map, words_class["ham"])
	RemoveDuplicates(spam_words, ham_words)
	return spam_words, ham_words
}

func Remove(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}


func GetMostFrequent(words map[string]int, class []string) []string {
	words_clone := map[string]int{}
	for k, v := range words {
		words_clone[k] = v
	}
	var max string;
	max_value := 0
	ret := make([]string, MOST_FREQUENT_WORDS)
	for i := 0; i < MOST_FREQUENT_WORDS; i++ {

		for value := range (words_clone) {
			if Contains(class, value) {
				if (words_clone[value] > max_value) {
					max = value
					max_value = words_clone[value]
				}
			}
		}
		ret[i] = max
		words_clone[max] = -1
		max_value = 0
	}
	return ret
}

func OpenFile(file string) ([][]string, error) {

	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close() // this needs to be after the err check

	lines, err := csv.NewReader(f).ReadAll()
	ret := make([][]string, len(lines))
	if err != nil {
		return nil, err
	}

	for idx, line := range lines {
		data := make([]string, 2)
		data[0] = strings.ToLower(line[0])
		data[1] = line[1]
		ret[idx] = data
	}
	return ret, nil
}


func OpenFileWords(file string) ([]string, error) {

	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close() // this needs to be after the err check

	lines, err := csv.NewReader(f).ReadAll()
	ret := make([]string, len(lines))
	if err != nil {
		return nil, err
	}

	for idx, line := range lines {
		ret[idx] =  strings.ToLower(line[0])
	}
	return ret, nil
}


func WriteFile(file string, lines []string) {
	f, _ := os.Create(file)
	defer f.Close() // this needs to be after the err check

	csv_writer := csv.NewWriter(f)
	for _, line := range lines {
		csv_writer.Write([]string{line})
	}

	defer f.Sync()
	defer csv_writer.Flush()
}

func Contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}


func ContainsWithIdx(arr []string, str string) int {

	for idx, a := range arr {
		if a == str {
			return idx
		}
	}
	return 1
}