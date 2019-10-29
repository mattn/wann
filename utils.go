package wann

import (
	"fmt"
	"sort"
)

// Pair is used for sorting dictionaries by value.
// Thanks https://stackoverflow.com/a/18695740/131264
type Pair struct {
	Key   int
	Value float64
}

// PairList is a slice of Pair
type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// SortByValue sorts a map[int]float64 by value
func SortByValue(m map[int]float64) PairList {
	pl := make(PairList, len(m))
	i := 0
	for k, v := range m {
		pl[i] = Pair{k, v}
		i++
	}
	sort.Sort(sort.Reverse(pl))
	return pl
}

func (neuron *Neuron) checkInputNeurons() {
	for _, inputNeuronIndex := range neuron.InputNodes {
		fmt.Println("Network of neuron ", neuron.neuronIndex, ":")
		fmt.Println(neuron.Net)
		if int(inputNeuronIndex) >= len(neuron.Net.AllNodes) {
			msg := fmt.Sprintf("at %d which has input index %d", neuron.neuronIndex, inputNeuronIndex)
			panic("input neuron index is pointing out of bounds: " + msg)
		}
	}
}

func (net *Network) checkInputNeurons() {
	for _, n := range net.All() {
		n.checkInputNeurons()
	}
}
