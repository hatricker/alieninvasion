package generators

import (
	"bufio"
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type fakeOneGen struct {
}

type fakeZeroGen struct {
}

var (
	fakeOneGenerator  = &fakeOneGen{}
	fakeZeroGenerator = &fakeZeroGen{}
	fakeArrGenerator  = &fakeNumArrayGen{}
)

func (fn *fakeOneGen) GenerateNum(_ int) int {
	return 1
}
func (fn *fakeZeroGen) GenerateNum(_ int) int {
	return 0
}

type fakeNumArrayGen struct {
}

func (fna fakeNumArrayGen) GenerateNums(num int) []int {
	result := make([]int, num)
	for i := 0; i < num; i++ {
		result[i] = i
	}
	return result
}

func TestGenerateCitiesOrAlienNames(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		input     int
		isCity    bool
		result    []string
		resultErr error
	}{
		{
			20,
			true,
			CityNames[:20],
			nil,
		},
		{
			0,
			true,
			[]string{},
			nil,
		},
		{
			len(CityNames) + 1,
			true,
			nil,
			ErrReqTooLarge,
		},
		{
			2,
			false,
			AlienNames[:2],
			nil,
		},
		{
			0,
			false,
			[]string{},
			nil,
		},
		{
			len(AlienNames) + 1,
			false,
			nil,
			ErrReqTooLarge,
		},
	}
	for _, tt := range tests {
		var result []string
		var err error
		if tt.isCity {
			result, err = GenerateCityNames(fakeArrGenerator, tt.input)

		} else {
			result, err = GenerateAlienNames(fakeArrGenerator, tt.input)
		}
		assert.Equal(tt.result, result)
		assert.Equal(tt.resultErr, err)
	}
}

func TestGenerateDirectionMask(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		x, y      int
		ng        NumGen
		result    [][]int
		resultErr error
	}{
		{
			0, 0, fakeOneGenerator,
			nil, ErrInvalidInput,
		},
		{
			1, 0, fakeOneGenerator,
			nil, ErrInvalidInput,
		},
		{
			0, 1, fakeOneGenerator,
			nil, ErrInvalidInput,
		},
		{
			100, 10, fakeOneGenerator,
			nil, ErrReqTooLarge,
		},
		{
			3, 3, fakeOneGenerator,
			[][]int{
				{0, 0, 0, 0},
				{0, East | South, East | West | South, West | South},
				{0, East | North | South, East | West | North | South, West | North | South},
				{0, East | North, East | West | North, West | North},
			}, nil,
		},
		{
			2, 2, fakeZeroGenerator,
			[][]int{
				{0, 0, 0},
				{0, 0, 0},
				{0, 0, 0},
			}, nil,
		},
	}

	for _, tt := range tests {
		result, err := GenerateDirectionMask(tt.x, tt.y, tt.ng)
		assert.Equal(tt.result, result)
		assert.Equal(tt.resultErr, err)
	}
}

func TestGenerateCityMapNeg(t *testing.T) {
	assert := assert.New(t)

	negTests := []struct {
		mask      [][]int
		cityNames []string
	}{
		{
			[][]int{{}},
			[]string{"Fremont"},
		},
		{
			[][]int{{1}, {1}},
			[]string{},
		},
	}

	for _, tt := range negTests {
		result := GenerateCityMap(tt.mask, tt.cityNames)
		assert.Empty(result)
	}
}

func TestGenerateCityMap(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		masks [][]int
	}{
		{
			[][]int{
				{0, 0, 0, 0},
				{0, East, West, South},
				{0, East | South, West, North | South},
				{0, East | North, West | East, West | North},
			},
		},
		{
			[][]int{
				{0, 0, 0, 0},
				{0, East, West, South},
				{0, South, 0, North | South},
				{0, East | North, West | East, West | North},
			},
		},
		{
			[][]int{
				{0, 0, 0},
				{0, East, West | South},
				{0, South, North | South},
				{0, East | North, West | North},
			},
		},
	}

	for _, tt := range tests {
		masks := tt.masks
		x, y := len(masks)-1, len(masks[0])-1
		cityNames, _ := GenerateCityNames(fakeArrGenerator, x*y)
		cityMap := GenerateCityMap(masks, cityNames)

		for i := 1; i <= x; i++ {
			for j := 1; j <= y; j++ {
				nameIndex := (i-1)*y + j - 1
				cn := cityNames[nameIndex]
				node, ok := cityMap[cn]
				assert.Equal(cn, node.Name)
				assert.Equal(true, ok)

				if masks[i][j]&East > 0 {
					assert.NotNil(node.East)
					eastNeighbor := cityMap[cityNames[nameIndex+1]]
					assert.Equal(node.East.Name, eastNeighbor.Name)
				}
				if masks[i][j]&West > 0 {
					assert.NotNil(node.West)
					westNeighbor := cityMap[cityNames[nameIndex-1]]
					assert.Equal(node.West.Name, westNeighbor.Name)
				}
				if masks[i][j]&North > 0 {
					assert.NotNil(node.North)
					northNeighbor := cityMap[cityNames[nameIndex-y]]
					assert.Equal(node.North.Name, northNeighbor.Name)
				}
				if masks[i][j]&South > 0 {
					assert.NotNil(node.South)
					southNeighbor := cityMap[cityNames[nameIndex+y]]
					assert.Equal(node.South.Name, southNeighbor.Name)
				}
			}
		}

	}
}

func TestGenerateCityMapFromSteam(t *testing.T) {
	assert := assert.New(t)

	input := "Foo,north=Bar,west=Baz,south=Qu-ux Bar,south=Foo,west=Bee"
	scanner := bufio.NewScanner(strings.NewReader(input))
	scanner.Split(bufio.ScanWords)

	cityMap := GenerateCityMapFromSteam(scanner, ',')

	assert.Equal("Bar", cityMap["Foo"].North.Name)
	assert.Equal("Baz", cityMap["Foo"].West.Name)
	assert.Equal("Qu-ux", cityMap["Foo"].South.Name)

	assert.Equal("Bee", cityMap["Bar"].West.Name)
	assert.Equal("Foo", cityMap["Bar"].South.Name)
}

func TestGenerateMapFile(t *testing.T) {
	assert := assert.New(t)
	var b bytes.Buffer

	input := "Foo,north=Bar,west=Baz,south=Qu-ux,east=Bee"
	scanner := bufio.NewScanner(strings.NewReader(input))
	scanner.Split(bufio.ScanWords)

	cityMap := GenerateCityMapFromSteam(scanner, ',')

	GenerateMapFile(cityMap, &b)
	assert.Equal("Foo east=Bee west=Baz north=Bar south=Qu-ux \n", b.String())
}
