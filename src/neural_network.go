package main

import (
     "fmt"
     "os"
     "./common"
     "github.com/jdkato/prose/tokenize"
     "github.com/gonum/matrix/mat64"
     "./neural"
     "./neural/config"

)

func main(){
     LEN_TRAIN := 10
     LEN_TEST := 200
     var spam []string
     var ham []string
     s_data, s_err := common.OpenFileWords("spam.txt")
     h_data, h_err := common.OpenFileWords("ham.txt")
     data, _ := common.OpenFile("/home/joxer/code/go/spam_filter/data/dataset.csv")
     if( s_err != nil && h_err != nil) {
          spam, ham = common.GetNames(data[0:LEN_TRAIN])
          common.WriteFile("spam.txt", spam)
          common.WriteFile("ham.txt", ham)
     } else{
          spam, ham = s_data, h_data
     }
     fmt.Println(spam)
     fmt.Println(ham)


     testset, _ := common.OpenFile("/home/joxer/code/go/spam_filter/data/testset.csv")
     tokenizer := tokenize.NewTreebankWordTokenizer()
     indices_strings,keys := generateIndices(spam,ham)
     features := make([]float64, 0)
     classes := make([]float64, 0)
     var array_len int
     for _, value := range(data[0:LEN_TRAIN]) {
          class, str := value[0], value[1]
          words := (tokenizer.Tokenize(str))
          array := (generateArrayBasedOnIndices(words, indices_strings, keys))
          array_len = len(array)
          if(class == common.HAM){
               classes = append(classes, 1.0)
          }else {
               classes = append(classes,2.0)
          }
          for _, val := range(array){
               features = append(features, val)
          }
     }
     testarray := make([]float64,0)
     for _, value := range(testset[0:LEN_TEST]) {
          str := value[1]
          words := (tokenizer.Tokenize(str))
          array := (generateArrayBasedOnIndices(words, indices_strings, keys))

           for _, val := range(array){
               testarray = append(testarray, val)
          }
     }


     net_config := getConfig(array_len)
     net, err := neural.NewNetwork(net_config)
     c := &config.TrainConfig{
          Kind:   "backprop",
          Cost:   "xentropy",
          Lambda: 1.0,
          Optimize: &config.OptimConfig{
               Method:     "bfgs",
               Iterations: 80,
          },
     }

     var inMx      *mat64.Dense
     var labelsVec *mat64.Vector
     inMx = mat64.NewDense(LEN_TRAIN, array_len, features)
     labelsVec = mat64.NewVector(len(classes), classes)

     if err != nil {
          fmt.Printf("Error creating network: %s\n", err)
          os.Exit(1)
     }
     fmt.Printf("Created new neural network: %v\n", net)
     net.Train(c,inMx, labelsVec)
     var testarray_dense  *mat64.Dense
     testarray_dense = mat64.NewDense(LEN_TEST, array_len, testarray)
     r,_ := net.Classify(testarray_dense)
     var right int
     var nright int
     for i:=0; i < LEN_TEST;i++{
          var rt string
          if r.At(i,0) > r.At(i,1) {
               rt = common.HAM
          } else {
               rt = common.SPAM
          }

          if testset[i][0] == rt {
               right += 1
          }else{
               nright += 1
          }

     }
     fmt.Println(right, nright)



}

func generateArrayBasedOnIndices(words []string, indicesArray map[string]int, keys []string) []float64 {

     ret := make([]float64, common.MOST_FREQUENT_WORDS*2)
     for _, value := range(words) {
          if(common.Contains(keys, value)){
               ret[indicesArray[value]] = 1.0
          }
     }
     return ret
}

func generateIndices(spam []string, ham []string) (map[string]int, []string) {

     ret := make(map[string]int)
     j := 0
     var i int
     for  i = 0; i < len(spam) && spam[i] != "";i++{
          if _, ok := ret[spam[i]]; !ok {
               ret[spam[i]] = j
               j += 1
          }
     }
     for i = 0; i < len(ham) && ham[i] != "";i++{
          if _, ok := ret[ham[i]]; !ok {
               ret[ham[i]] = j
               j += 1
          }
     }

     var keys []string
     for key, _ := range ret {
          keys = append(keys, key)
     }

     return ret,keys
}

func getConfig(len int) *config.NetConfig {
     return &config.NetConfig{
          Kind: "feedfwd",
          Arch: &config.NetArch{
               Input: &config.LayerConfig{
                    Kind: "input",
                    Size: len,
               },
               Hidden: []*config.LayerConfig{
                    &config.LayerConfig{
                         Kind: "hidden",
                         Size: 20,
                         NeurFn: &config.NeuronConfig{
                              Activation: "tanh",
                         },
                    },

               },
               Output: &config.LayerConfig{
                    Kind: "output",
                    Size: 2,
                    NeurFn: &config.NeuronConfig{
                         Activation: "softmax",
                    },
               },
          },
     }
}

