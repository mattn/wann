package wann

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
)

// ScorePopulation evaluates a population, given a slice of input numbers.
// It returns a map with scores, together with the sum of scores.
func ScorePopulation(population []Network, inputData [][]float64, correctOutputMultipliers []float64) (map[int]float64, float64) {
	// Use a shared random weight
	w := rand.Float64()

	scoreMap := make(map[int]float64)
	scoreSum := 0.0

	for i := 0; i < len(population); i++ {
		net := population[i]

		net.checkInputNeurons()

		net.SetWeight(w)
		result := 0.0
		for i := 0; i < len(inputData); i++ {
			result += net.Evaluate(inputData[i]) * correctOutputMultipliers[i]
		}
		score := result / net.Complexity()
		scoreSum += score
		scoreMap[i] = score
	}
	return scoreMap, scoreSum
}

// Evolve evolves a neural network, given a slice of training data and a slice of correct output values.
// Will overwrite config.Inputs.
func (config *Config) Evolve(inputData [][]float64, correctOutputMultipliers []float64) (*Network, error) {

	inputLength := len(inputData)

	if inputLength == 0 {
		return nil, errors.New("no input data")
	}

	if len(correctOutputMultipliers) == 1 && inputLength != 1 {
		// Assume the first slice of floats in the input data is the correct and that the rest are examples of being wrong
		for i := 1; i < inputLength; i++ {
			correctOutputMultipliers = append(correctOutputMultipliers, -1.0)
		}
	} else if inputLength != len(correctOutputMultipliers) {
		// Assume that the list of correct output multipliers should match the length of the float64 slices in inputData
		return nil, errors.New("the length of the input data and the slice of output multipliers differs")
	}

	config.Inputs = len(inputData[0])

	population := make([]Network, config.PopulationSize)

	// Initialize the population
	for i := 0; i < config.PopulationSize; i++ {
		population[i] = NewNetwork(config)
	}

	var bestNetwork *Network

	// Keep track of the best scores
	bestScore := 0.0
	lastBestScore := 0.0
	noImprovementOfBestScoreCounter := 0

	// Keep track of the average scores
	averageScore := 0.0
	lastAverageScore := 0.0
	noImprovementOfAverageScoreCounter := 0

	// Keep track of the worst scores
	worstScore := 0.0
	lastWorstScore := 0.0
	noImprovementOfWorstScoreCounter := 0

	// For each generation, evaluate and modify the networks
	for j := 0; j < config.Generations; j++ {

		bestNetwork = nil

		// Initialize the scores with unlikely values
		// TODO: Use the first network in the population for initializing these instead
		bestScore = -9999.0
		averageScore = 0.0
		worstScore = 9999.0

		if config.Verbose {
			fmt.Println("------ generation " + strconv.Itoa(j) + ", population size " + strconv.Itoa(len(population)))
		}

		// The scores for this generation (using a random shared weight within ScorePopulation).
		// CorrectOutputMultipliers gives weight to the "correct" or "wrong" results, with the same index as the inputData
		// Score each network in the population.
		scoreMap, scoreSum := ScorePopulation(population, inputData, correctOutputMultipliers)

		// Sort by score
		scoreList := SortByValue(scoreMap)

		// Handle the best score stats
		lastBestScore = bestScore
		if scoreList[0].Value > bestScore {
			bestScore = scoreList[0].Value
			bestNetwork = &(population[scoreList[0].Key])
		}
		if bestScore >= lastBestScore {
			noImprovementOfBestScoreCounter++
		} else {
			noImprovementOfBestScoreCounter = 0
		}

		// Handle the average score stats
		lastAverageScore = averageScore
		averageScore = scoreSum / float64(config.PopulationSize)
		if averageScore >= lastAverageScore {
			noImprovementOfAverageScoreCounter++
		} else {
			noImprovementOfAverageScoreCounter = 0
		}

		// Handle the worst score stats
		lastWorstScore = worstScore
		if scoreList[len(scoreList)-1].Value < worstScore {
			worstScore = scoreList[len(scoreList)-1].Value
		}
		if worstScore >= lastWorstScore {
			noImprovementOfWorstScoreCounter++
		} else {
			noImprovementOfWorstScoreCounter = 0
		}

		if bestNetwork == nil {
			panic("implementation error: no best network")
		}

		if config.Verbose {
			fmt.Println("Best, average and worst score:", bestScore, averageScore, worstScore)
			fmt.Println("Best, average and worst improvement counters:", noImprovementOfBestScoreCounter, noImprovementOfAverageScoreCounter, noImprovementOfWorstScoreCounter)
		}

		bestThirdCountdown := len(population) / 3

		// Now loop over all networks, sorted by score (descending order)
		for _, p := range scoreList {
			networkIndex := p.Key
			networkScore := p.Value
			if bestThirdCountdown > 0 {
				bestThirdCountdown--
				// In the best third of the networks
				fmt.Println("BEST THIRD:", networkIndex, "score", networkScore)
			} else {
				fmt.Println("WORST TWO THIRDS:", networkIndex, "score", networkScore)
			}
		}

		// // TODO: Rewrite. The entire loop should loop from highest to lowest scoring.
		// for networkIndex := 0; networkIndex < config.PopulationSize; networkIndex++ {
		// 	networkScore := scoreMap[networkIndex]
		// 	// Is this network in the best half?
		// 	bestHalf := (networkScore >= averageScore) && (networkScore > 0)
		// 	if !bestHalf {
		// 		// Take a proper copy, not just the the pointers, because the nodes will be changed
		// 		// Assign it to the population, replacing the low-scoring ones
		// 		// TODO: Actually replace the low-scoring ones
		// 		newNetwork := bestNetwork.Copy()
		// 		newNetwork.Modify(100)
		// 		population[networkIndex] = newNetwork
		// 	}
		// }

	}
	if bestNetwork == nil {
		return nil, errors.New("the best network is nil")
	}
	return bestNetwork, nil
}

// Modify this network a bit
func (net *Network) Modify(maxIterations int) {

	//fmt.Println("A")
	//net.checkInputNeurons()

	// Use method 0, 1 or 2
	method := rand.Intn(3) // up to and not including 3
	//method := 0
	// TODO: Perform a modfification, using one of the three methods outlined in the paper
	switch method {
	case 0:
		//fmt.Println("Modifying the network using method 1 - insert node")

		// It's important that GetRandomNeuron is used before NewRandomNeuron is called
		nodeA, nodeB := net.GetRandomNode(), net.GetRandomNode()

		//fmt.Println("MODIFY METHOD 0, START, MAX ITERATIONS:", maxIterations)
		_, newNodeIndex := net.NewNeuron()
		//fmt.Println("NEW NEURON AT INDEX", newNodeIndex)

		//fmt.Println("USING NODE A AND B:", nodeA, nodeB)

		// A bit risky, time-wise, but continue finding random neurons until they work out
		// Insert a new node with a random activation function
		counter := 0

		// InsertNode adds the new node to net.AllNodes
		err := net.InsertNode(nodeA, nodeB, newNodeIndex)

		if err != nil {
			//fmt.Println("INSERT NODE ERROR: " + err.Error())
		}

		if !net.AllNodes[net.OutputNode].InputNeuronsAreGood() {
			//panic("implementation error: Modify: input neurons are not good")
		}

		for err != nil {
			//(fmt.Println("COUNTER", counter)
			nodeA, nodeB = net.GetRandomNode(), net.GetRandomNode()
			counter++
			//fmt.Println("COUNTER", counter, "MAX ITERATIONS", maxIterations)
			if maxIterations > 0 && counter > maxIterations {
				// Could not add a new node. This may happen if the network is only input nodes and one output node
				//panic("implementation error: could not a add a new node, even after " + strconv.Itoa(maxIterations) + " iterations: " + err.Error())
				// Add a node between a random input node and the output node
				err = net.InsertNode(net.GetRandomInputNode(), net.OutputNode, newNodeIndex)

				if err != nil {
					//fmt.Println("INSERT NODE, LAST DITCH ERROR: " + err.Error())
				}
				// if the randomly chosen input node already connects to the output node, then that's fine, let`s move on
				return
			}
			err = net.InsertNode(nodeA, nodeB, newNodeIndex)
			//if err != nil {
			//	fmt.Println("INSERT NODE ERROR: " + err.Error())
			//}

		}
		if err != nil {
			// This should never happen, since adding a node between an input node and the output node should always work
			//panic("implementation error : " + err.Error())
		}

	case 1:
		//fmt.Println("Modifying the network using method 2 - add connection")

		nodeA, nodeB := net.GetRandomNode(), net.GetRandomNode()
		// A bit risky, time-wise, but continue finding random neurons until they work out
		// Create a new connection
		counter := 0
		for net.AddConnection(nodeA, nodeB) != nil {
			nodeA, nodeB = net.GetRandomNode(), net.GetRandomNode()
			counter++
			if maxIterations > 0 && counter > maxIterations {
				// Could not add a connection. The possibilities for connections might be saturated.
				return
			}
		}
	case 2:
		//fmt.Println("Modifying the network using method 3 - change activation")
		// Change the activation function
		net.RandomizeActivationFunctionForRandomNeuron()
	default:
		panic("implementation error: invalid method number: " + strconv.Itoa(method))
	}

	//fmt.Println("B")
	//net.checkInputNeurons()
}
