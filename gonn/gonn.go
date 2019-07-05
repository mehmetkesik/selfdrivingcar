//golang basic neural network library

package gonn

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"time"
)

type MyNetwork struct {
	nn     *Network
	isInit bool
}

func (self *MyNetwork) Init(neurons []int) {
	var err error
	self.nn, err = NewNetwork(neurons)
	if err != nil {
		panic(err)
	}
	self.isInit = true
}

func (self *MyNetwork) Train(trainingData [][][]float64, epochs int, learningRate float64, debug bool) {
	if !self.isInit {
		fmt.Println("Must be initial. Call Inıt() function..")
		return
	}
	if epochs < 1 {
		epochs = 1
	}
	if learningRate == 0 {
		learningRate = 0.05
	}
	self.nn.Train(trainingData, epochs, learningRate, debug)
}

func (self *MyNetwork) SingleTrain(inputData []float64, outputData []float64) {
	if !self.isInit {
		fmt.Println("Must be initial. Call Inıt() function..")
		return
	}
	var trainingData = [][][]float64{[][]float64{inputData, outputData,},}
	self.nn.Train(trainingData, 1, 0.05, false)
}

func (self *MyNetwork) Predict(inputData []float64) []float64 {
	if !self.isInit {
		fmt.Println("Must be initial. Call Inıt() function..")
		return nil
	}
	return self.nn.Predict(inputData)
}

func (self *MyNetwork) ArealMutation(areaRate, mutationRate float64) {
	w := self.nn.GetWeights()
	rand.Seed(time.Now().UTC().UnixNano())
	nodeCount := self.NodeCount()
	for m := 0; m < len(w); m++ {
		for i := 0; i < len(w[m]); i++ {
			for j := 0; j < len(w[m][i]); j++ {
				rast := rand.Intn(nodeCount)
				if rast > int(areaRate*float64(nodeCount)) {
					continue
				}
				w[m][i][j] += (((rand.Float64() * mutationRate) * 2) - mutationRate)
				/*if w[m][i][j] > 1 {
					w[m][i][j] = 1
				}
				if w[m][i][j] < 0 {
					w[m][i][j] = 0
				}*/
			}
		}
	}
	self.nn.SetWeights(w)
}

func (self *MyNetwork) NodeCount() (count int) {
	w := self.nn.GetWeights()
	for m := 0; m < len(w); m++ {
		for i := 0; i < len(w[m]); i++ {
			for j := 0; j < len(w[m][i]); j++ {
				count++
			}
		}
	}
	return
}

func (self *MyNetwork) GetWeights() [][][]float64 {
	return self.nn.GetWeights()
}

func (self *MyNetwork) SetWeights(nw [][][]float64) {
	self.nn.SetWeights(nw)
}

func (self *MyNetwork) Save(filename string) {
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	err = self.nn.Export(file)
	if err != nil {
		panic(err)
	}
}

func (self *MyNetwork) Load(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	b, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	net, err := Load(b)
	if err != nil {
		panic(err)
	}
	self.nn = net
}
