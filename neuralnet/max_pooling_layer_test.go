package neuralnet

import (
	"math/rand"
	"testing"

	"github.com/unixpickle/autofunc"
	"github.com/unixpickle/autofunc/functest"
	"github.com/unixpickle/num-analysis/linalg"
	"github.com/unixpickle/serializer"
)

func TestMaxPoolingDimensions(t *testing.T) {
	layers := []*MaxPoolingLayer{
		{3, 3, 9, 9, 5},
		{2, 2, 9, 9, 7},
		{4, 10, 30, 51, 2},
	}
	outSizes := [][3]int{
		{3, 3, 5},
		{5, 5, 7},
		{8, 6, 2},
	}
	for i, layer := range layers {
		expOutSize := outSizes[i]
		if layer.OutputWidth() != expOutSize[0] ||
			layer.OutputHeight() != expOutSize[1] {
			t.Errorf("test %d gave downstream size %dX%dX%d (expected %dX%dX%d)",
				i, layer.OutputWidth(), layer.OutputHeight(), layer.InputDepth,
				expOutSize[0], expOutSize[1], expOutSize[2])
		}
	}
}

func TestMaxPoolingForward(t *testing.T) {
	layer := &MaxPoolingLayer{3, 3, 10, 11, 2}

	input := []float64{
		0.5305, 0.7935, 0.3718, 0.4026, 0.8246, 0.6875, 0.6069, 0.0399, 0.4759, 0.3548, 0.8465, 0.0479, 0.4841, 0.1277, 0.2060, 0.6833, 0.0844, 0.0793, 0.1564, 0.2891,
		0.9761, 0.1716, 0.2394, 0.6439, 0.2834, 0.5429, 0.5479, 0.6228, 0.3308, 0.4145, 0.4472, 0.8445, 0.1258, 0.9365, 0.8861, 0.5686, 0.7676, 0.5818, 0.8840, 0.4068,
		0.0427, 0.2888, 0.2321, 0.2350, 0.3702, 0.8161, 0.9992, 0.3097, 0.2996, 0.7116, 0.6126, 0.5868, 0.0587, 0.3701, 0.8875, 0.5653, 0.1161, 0.3778, 0.5768, 0.6405,

		0.2868, 0.2617, 0.6762, 0.9683, 0.7948, 0.8449, 0.7876, 0.3225, 0.0139, 0.2315, 0.5635, 0.5076, 0.8530, 0.4785, 0.8244, 0.0356, 0.1402, 0.8464, 0.6470, 0.5444,
		0.4489, 0.3268, 0.9251, 0.6568, 0.7592, 0.0223, 0.6244, 0.9696, 0.2035, 0.6457, 0.0505, 0.8712, 0.2836, 0.0689, 0.6179, 0.0421, 0.0373, 0.2316, 0.7921, 0.7195,
		0.7107, 0.7147, 0.3756, 0.0563, 0.3803, 0.4184, 0.2551, 0.7702, 0.8207, 0.9405, 0.4711, 0.1529, 0.1081, 0.6531, 0.5117, 0.1368, 0.2331, 0.7265, 0.0986, 0.7236,

		0.1467, 0.1398, 0.4580, 0.1640, 0.2878, 0.3895, 0.5600, 0.1037, 0.9899, 0.8434, 0.5762, 0.3068, 0.6564, 0.4465, 0.0134, 0.8445, 0.8760, 0.9951, 0.4819, 0.5924,
		0.2894, 0.4773, 0.0628, 0.3025, 0.2345, 0.9472, 0.7258, 0.2077, 0.3428, 0.6104, 0.0639, 0.0854, 1.0000, 0.0372, 0.3874, 0.6501, 0.6533, 0.2953, 0.5591, 0.9967,
		0.6510, 0.3776, 0.6511, 0.9123, 0.9738, 0.4100, 0.3743, 0.9791, 0.3929, 0.8278, 0.1919, 0.2566, 0.3484, 0.3768, 0.0108, 0.5234, 0.4480, 0.3097, 0.5598, 0.5840,

		0.0082, 0.5011, 0.3124, 0.8709, 0.6181, 0.1428, 0.7824, 0.7105, 0.0922, 0.5858, 0.1643, 0.3963, 0.1715, 0.2448, 0.7961, 0.1675, 0.2949, 0.3438, 0.4825, 0.8616,
		0.5648, 0.3950, 0.7001, 0.3238, 0.3235, 0.4789, 0.4206, 0.0502, 0.3165, 0.2146, 0.5393, 0.9277, 0.4361, 0.1530, 0.3192, 0.9463, 0.0317, 0.3078, 0.8892, 0.0508,
	}

	output := []float64{
		0.9761, 0.8161, 0.9992, 0.8445, 0.8875, 0.9365, 0.8840, 0.6405,
		0.9251, 0.9683, 0.8207, 0.9696, 0.8530, 0.8464, 0.7921, 0.7236,
		0.9738, 0.9472, 0.9899, 0.9791, 1.0000, 0.9951, 0.5598, 0.9967,
		0.7001, 0.8709, 0.7824, 0.9277, 0.7961, 0.9463, 0.8892, 0.8616,
	}

	result := layer.Apply(&autofunc.Variable{input}).Output()
	for i, x := range output {
		if actual := result[i]; actual != x {
			t.Errorf("expected output %d to be %f but got %f", i, x, actual)
		}
	}
}

func TestMaxPoolingBackward(t *testing.T) {
	layer := &MaxPoolingLayer{3, 3, 10, 11, 2}

	input := []float64{
		0.5305, 0.7935, 0.3718, 0.4026, 0.8246, 0.6875, 0.6069, 0.0399, 0.4759, 0.3548, 0.8465, 0.0479, 0.4841, 0.1277, 0.2060, 0.6833, 0.0844, 0.0793, 0.1564, 0.2891,
		0.9761, 0.1716, 0.2394, 0.6439, 0.2834, 0.5429, 0.5479, 0.6228, 0.3308, 0.4145, 0.4472, 0.8445, 0.1258, 0.9365, 0.8861, 0.5686, 0.7676, 0.5818, 0.8840, 0.4068,
		0.0427, 0.2888, 0.2321, 0.2350, 0.3702, 0.8161, 0.9992, 0.3097, 0.2996, 0.7116, 0.6126, 0.5868, 0.0587, 0.3701, 0.8875, 0.5653, 0.1161, 0.3778, 0.5768, 0.6405,

		0.2868, 0.2617, 0.6762, 0.9683, 0.7948, 0.8449, 0.7876, 0.3225, 0.0139, 0.2315, 0.5635, 0.5076, 0.8530, 0.4785, 0.8244, 0.0356, 0.1402, 0.8464, 0.6470, 0.5444,
		0.4489, 0.3268, 0.9251, 0.6568, 0.7592, 0.0223, 0.6244, 0.9696, 0.2035, 0.6457, 0.0505, 0.8712, 0.2836, 0.0689, 0.6179, 0.0421, 0.0373, 0.2316, 0.7921, 0.7195,
		0.7107, 0.7147, 0.3756, 0.0563, 0.3803, 0.4184, 0.2551, 0.7702, 0.8207, 0.9405, 0.4711, 0.1529, 0.1081, 0.6531, 0.5117, 0.1368, 0.2331, 0.7265, 0.0986, 0.7236,

		0.1467, 0.1398, 0.4580, 0.1640, 0.2878, 0.3895, 0.5600, 0.1037, 0.9899, 0.8434, 0.5762, 0.3068, 0.6564, 0.4465, 0.0134, 0.8445, 0.8760, 0.9951, 0.4819, 0.5924,
		0.2894, 0.4773, 0.0628, 0.3025, 0.2345, 0.9472, 0.7258, 0.2077, 0.3428, 0.6104, 0.0639, 0.0854, 1.0000, 0.0372, 0.3874, 0.6501, 0.6533, 0.2953, 0.5591, 0.9967,
		0.6510, 0.3776, 0.6511, 0.9123, 0.9738, 0.4100, 0.3743, 0.9791, 0.3929, 0.8278, 0.1919, 0.2566, 0.3484, 0.3768, 0.0108, 0.5234, 0.4480, 0.3097, 0.5598, 0.5840,

		0.0082, 0.5011, 0.3124, 0.8709, 0.6181, 0.1428, 0.7824, 0.7105, 0.0922, 0.5858, 0.1643, 0.3963, 0.1715, 0.2448, 0.7961, 0.1675, 0.2949, 0.3438, 0.4825, 0.8616,
		0.5648, 0.3950, 0.7001, 0.3238, 0.3235, 0.4789, 0.4206, 0.0502, 0.3165, 0.2146, 0.5393, 0.9277, 0.4361, 0.1530, 0.3192, 0.9463, 0.0317, 0.3078, 0.8892, 0.0508,
	}

	downstreamGrad := NewTensor3(4, 4, 2)
	for i := range downstreamGrad.Data {
		downstreamGrad.Data[i] = rand.Float64()*2 - 1
	}

	gradientMask := []int{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 1, 0,
		0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1,

		0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0,
		0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,

		0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0,
		0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1,
		0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0,

		0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1,
		0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 1, 0,
	}

	inputVar := &autofunc.Variable{input}
	g := autofunc.NewGradient([]*autofunc.Variable{inputVar})
	downstreamCopy := make(linalg.Vector, len(downstreamGrad.Data))
	copy(downstreamCopy, downstreamGrad.Data)
	layer.Apply(inputVar).PropagateGradient(downstreamCopy, g)
	actualGrad := g[inputVar]

	idx := 0
	for y := 0; y < 11; y++ {
		for x := 0; x < 10; x++ {
			for z := 0; z < 2; z++ {
				isChosen := gradientMask[z+x*2+y*20] == 1
				gradValue := actualGrad[idx]
				outputGrad := downstreamGrad.Get(x/3, y/3, z)
				idx++
				if !isChosen && gradValue != 0 {
					t.Errorf("expected gradient at %d,%d,%d to be 0, but got %f",
						x, y, z, gradValue)
				} else if isChosen && gradValue != outputGrad {
					t.Errorf("expected gradient at %d,%d,%d to be %f, but got %f",
						x, y, z, outputGrad, gradValue)
				}
			}
		}
	}
}

func TestMaxPoolingBatch(t *testing.T) {
	layer := &MaxPoolingLayer{
		XSpan:       5,
		YSpan:       4,
		InputWidth:  17,
		InputHeight: 19,
		InputDepth:  3,
	}

	n := 3
	batchInput := make(linalg.Vector, n*layer.InputWidth*layer.InputHeight*layer.InputDepth)
	for i := range batchInput {
		batchInput[i] = rand.NormFloat64()
	}
	batchRes := &autofunc.Variable{Vector: batchInput}

	testBatcher(t, layer, batchRes, n, []*autofunc.Variable{batchRes})
}

func TestMaxPoolingBatchR(t *testing.T) {
	layer := &MaxPoolingLayer{
		XSpan:       5,
		YSpan:       4,
		InputWidth:  17,
		InputHeight: 19,
		InputDepth:  3,
	}

	n := 3
	batchInput := make(linalg.Vector, n*layer.InputWidth*layer.InputHeight*layer.InputDepth)
	for i := range batchInput {
		batchInput[i] = rand.NormFloat64()
	}
	batchRes := &autofunc.Variable{Vector: batchInput}

	rVec := autofunc.RVector{
		batchRes: make(linalg.Vector, len(batchInput)),
	}
	for i := range rVec[batchRes] {
		rVec[batchRes][i] = rand.NormFloat64()
	}

	testRBatcher(t, rVec, layer, autofunc.NewRVariable(batchRes, rVec),
		n, []*autofunc.Variable{batchRes})
}

func TestMaxPoolingSerialize(t *testing.T) {
	layer := &MaxPoolingLayer{3, 3, 10, 11, 2}
	encoded, err := layer.Serialize()
	if err != nil {
		t.Fatal(err)
	}
	layerType := layer.SerializerType()
	decoded, err := serializer.GetDeserializer(layerType)(encoded)
	if err != nil {
		t.Fatal(err)
	}
	layer, ok := decoded.(*MaxPoolingLayer)
	if !ok {
		t.Fatalf("expected *MaxPoolingLayer but got %T", decoded)
	}
}

func TestMaxPoolingRProp(t *testing.T) {
	layer := &MaxPoolingLayer{3, 3, 10, 11, 2}
	input := make(linalg.Vector, 10*11*2)
	inputR := make(linalg.Vector, len(input))
	for i := range input {
		input[i] = rand.Float64()*2 - 1
		inputR[i] = rand.Float64()*2 - 1
	}

	inputVar := &autofunc.Variable{input}
	rVector := autofunc.RVector{inputVar: inputR}

	funcTest := &functest.RFuncChecker{
		F:     layer,
		Vars:  []*autofunc.Variable{inputVar},
		Input: inputVar,
		RV:    rVector,
	}
	funcTest.FullCheck(t)
}
