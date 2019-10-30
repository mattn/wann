# WANN [![Build Status](https://travis-ci.org/xyproto/wann.svg?branch=master)](https://travis-ci.org/xyproto/wann) [![Go Report Card](https://goreportcard.com/badge/github.com/xyproto/wann)](https://goreportcard.com/report/github.com/xyproto/wann) [![GoDoc](https://godoc.org/github.com/xyproto/wann?status.svg)](https://godoc.org/github.com/xyproto/wann)

<img alt=Network src=img/after.svg width=128 />

Weight Agnostic Neural Networks, implemented in Go, using the techniques outlined in the paper:

*"Weight Agnostic Neural Networks" by Adam Gaier and David Ha*. ([PDF](https://arxiv.org/pdf/1906.04358.pdf) | [Interactive version](https://weightagnostic.github.io/) | [Google AI blog post](https://ai.googleblog.com/2019/08/exploring-weight-agnostic-neural.html))

## Features and limitations

* Neural networks can be trained and used, but I have only tried this on very simple training data and there is surely a lot of room for improvement, both in term of benchmarking/profiling and controlling the rate of mutation.
* A random weight is chosen when training, instead of looping over the range of the weight. The paper describes both methods.
* After the network has been trained, the optimal weight is found by looping over all weights (with a step size of `0.0001`).
* Increased complexity counts negatively when evolving networks. A quick benchmark of all available activation functions at the start of the program determines which activation function us more complex. This optimizes not only for less complex networks, but also for execution speed.
* The diagram drawing routine plots the activation functions directly onto the nodes, together with a label. This can be saved as an SVG file.

## Example program

This is a simple example, for creating a network that can recognize one of four shapes:

```go
package main

import (
	"fmt"
	"os"

	"github.com/xyproto/wann"
)

func main() {
	// Here are four shapes, representing: up, down, left and right:

	up := []float64{
		0.0, 1.0, 0.0, //  o
		1.0, 1.0, 1.0} // ooo

	down := []float64{
		1.0, 1.0, 1.0, // ooo
		0.0, 1.0, 0.0} //  o

	left := []float64{
		1.0, 1.0, 1.0, // ooo
		0.0, 0.0, 1.0} //   o

	right := []float64{
		1.0, 1.0, 1.0, // ooo
		0.1, 0.0, 0.0} // o

	// Prepare the input data as a 2D slice
	inputData := [][]float64{
		up,
		down,
		left,
		right,
	}

	// Which of the elements in the input data are we trying to identify?
	correctResultsForUp := []float64{1.0, -1.0, -1.0, -1.0}

	// Prepare a neural network configuration struct
	config := &wann.Config{
		InitialConnectionRatio: 0.1,
		Generations:            2000,
		PopulationSize:         200,
		Verbose:                true,
	}

	// Evolve a network, using the input data and the sought after results
	trainedNetwork, err := config.Evolve(inputData, correctResultsForUp)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}

	// Now to test the trained network on 4 different inputs and see if it passes the test
	upScore := trainedNetwork.Evaluate(up)
	downScore := trainedNetwork.Evaluate(down)
	leftScore := trainedNetwork.Evaluate(left)
	rightScore := trainedNetwork.Evaluate(right)

	if config.Verbose {
		if upScore > downScore && upScore > leftScore && upScore > rightScore {
			fmt.Println("Network training complete, the results are good.")
		} else {
			fmt.Println("Network training complete, but the results did not pass the test.")
		}
	}

	// Save the trained network as an SVG image
	if config.Verbose {
		fmt.Print("Writing network.svg...")
	}
	if err := trainedNetwork.WriteSVG("network.svg"); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
	if config.Verbose {
		fmt.Println("ok")
	}
}
```

Here's the resulting network generated by the above program (the results will very, since there is randomness involved):

<img alt=Network src=img/labels.svg width=256 />

Currently, the input node is not labeled with which input number it uses, but I plan to implement this.

## Quick start

This requires Go 1.11 or later.

Clone the repository:

    git clone https://github.com/xyproto/wann

Enter the `cmd/evolve` directory:

    cd wann/cmd/evolve

Build and run the example:

    go build && ./evolve

Take a look at the best network for judging if a set of numbers that are either 0 or 1 are of one category:

    xdg-open network.svg

(If needed, use your favorite SVG viewer instead of the `xdg-open` command).

## Ideas

* Adding convolution nodes might give interesting results.

## Generating Go code from a trained network

This is an experimental feature and a work in progress!

The idea is to generate one large expression from all the expressions that each node in the network represents.

Right now, his only works for networks that has a depth of 1.

For example, adding these two lines to `cmd/evolve/main.go`:

```go
// Output a Go function for this network
fmt.Println(trainedNetwork.GoFunction())
```

Produces this output:

```go
func f(x float64) float64 { return -x }
```

The plan is to output a function that takes the input data instead, and refers to the input data by index. Support for deeper networks also needs to be added.

There is a complete example for outputting Go code in `cmd/gofunction`.

## General info

* Version: 0.2.0
* License: MIT
* Author: Alexander F. Rødseth &lt;xyproto@archlinux.org&gt;
