package fc

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFC_Panic(t *testing.T) {
	t.Run("invalid_num_layers", func(t *testing.T) {
		defer func(t *testing.T) {
			if r := recover(); r == nil {
				t.Fatal("expecting panic")
			}
		}(t)
		NewFC(1)
	})

	t.Run("invalid_layer_size", func(t *testing.T) {
		defer func(t *testing.T) {
			if r := recover(); r == nil {
				t.Fatal("expecting panic")
			}
		}(t)
		NewFC(1, 0)
	})
}

func TestNewFC_Success(t *testing.T) {
	t.Run("3_2_1_network", func(t *testing.T) {
		fc := NewFC(3, 2, 1)
		assert.Len(t, fc.Layers, 2)

		// 2 nodes for layer 0 with 3 weights each
		assert.Len(t, fc.Layers[0].Nodes, 2)
		assert.Len(t, fc.Layers[0].Nodes[0].Weights, 3)
		assert.Len(t, fc.Layers[0].Nodes[1].Weights, 3)

		// 1 node for layer 1 with 2 weights each
		assert.Len(t, fc.Layers[1].Nodes, 1)
		assert.Len(t, fc.Layers[1].Nodes[0].Weights, 2)
	})
}

func TestPredict_Success(t *testing.T) {
	t.Run("2_2_2_network", func(t *testing.T) {
		fc := example_2_2_2()

		prediction := fc.Predict([]float64{.05, .10})

		assert.InDelta(t, 0.470281212, prediction[0], 0.000000001)
		assert.InDelta(t, 0.529718787, prediction[1], 0.000000001)
	})
}

func TestTrain_Success(t *testing.T) {
	t.Run("10_8_4_2_network", func(t *testing.T) {
		fc := NewFC(10, 8, 4, 2)

		const epochs = 10000
		for i := 0; i < epochs; i++ {
			fc.Train([][]float64{{.05, .22, .51, .77, .41, .94, .29, .0, .04, .01}}, [][]float64{{0.01, .99}}, 0.5)
		}

		prediction := fc.Predict([]float64{.05, .22, .51, .77, .41, .94, .29, .0, .04, .01})

		assert.InDelta(t, 0.01, prediction[0], 0.002)
		assert.InDelta(t, 0.99, prediction[1], 0.002)
	})
}

func TestSaveLoad_Success(t *testing.T) {
	t.Run("2_2_2_network", func(t *testing.T) {
		fc := NewFC(2, 2, 2)

		const epochs = 10000
		for i := 0; i < epochs; i++ {
			fc.Train([][]float64{{.05, .10}}, [][]float64{{0.01, .99}}, 0.5)
		}

		err := fc.Save("./test.gob")
		if err != nil {
			t.Fatal(err)
		}

		fc2, err := LoadFromFile("./test.gob")
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, fc.Sizes, fc2.Sizes)
		assert.Equal(t, fc.Layers[0].Nodes[0].Weights, fc2.Layers[0].Nodes[0].Weights)
		assert.Equal(t, fc.Layers[0].Nodes[0].Bias, fc2.Layers[0].Nodes[0].Bias)

		assert.Equal(t, fc.Layers[0].Nodes[1].Weights, fc2.Layers[0].Nodes[1].Weights)
		assert.Equal(t, fc.Layers[0].Nodes[1].Bias, fc2.Layers[0].Nodes[1].Bias)

		assert.Equal(t, fc.Layers[1].Nodes[0].Weights, fc2.Layers[1].Nodes[0].Weights)
		assert.Equal(t, fc.Layers[1].Nodes[0].Bias, fc2.Layers[1].Nodes[0].Bias)

		assert.Equal(t, fc.Layers[1].Nodes[1].Weights, fc2.Layers[1].Nodes[1].Weights)
		assert.Equal(t, fc.Layers[1].Nodes[1].Bias, fc2.Layers[1].Nodes[1].Bias)
	})
}

func nodesOutput(nodes []*fcNode) string {
	var s string

	for i, node := range nodes {
		s += fmt.Sprintf(`
		*****Node %d*****
		weights: %v
		delta: %v
		netErr: %v
		output: %v
		*****************
		`, i, node.Weights, node.delta, node.netErr, node.output)
	}

	return s
}

// https://mattmazur.com/2015/03/17/a-step-by-step-backpropagation-example
func example_2_2_2() *FC {
	fc := NewFC(2, 2, 2)
	// hidden layer
	fc.Layers[0].Nodes[0].Weights[0] = .15
	fc.Layers[0].Nodes[0].Weights[1] = .20
	fc.Layers[0].Nodes[1].Weights[0] = .25
	fc.Layers[0].Nodes[1].Weights[1] = .30
	fc.Layers[0].Nodes[0].Bias = .35
	fc.Layers[0].Nodes[1].Bias = .35

	// output layer
	fc.Layers[1].Nodes[0].Weights[0] = .40
	fc.Layers[1].Nodes[0].Weights[1] = .45
	fc.Layers[1].Nodes[1].Weights[0] = .50
	fc.Layers[1].Nodes[1].Weights[1] = .55
	fc.Layers[1].Nodes[0].Bias = .60
	fc.Layers[1].Nodes[1].Bias = .60

	return fc
}
