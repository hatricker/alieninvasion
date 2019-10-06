package games

import (
	"testing"

	"github.com/hatricker/alieninvasion/generators"
	"github.com/stretchr/testify/assert"
)

type fakeZeroGen struct {
}

type fakeNumArrayGen struct {
}

var (
	fakeZeroGenerator = &fakeZeroGen{}
	fakeArrGenerator  = &fakeNumArrayGen{}
	east              = generators.East
	west              = generators.West
	north             = generators.North
	south             = generators.South
	testingMasks      = [][]int{
		{0, 0, 0},
		{0, east | south, west | south},
		{0, east | north, west | north},
	}
	testingAlien     = generators.AlienNames[0]
	testingCityNames = generators.CityNames[:4]
)

func (fn *fakeZeroGen) GenerateNum(_ int) int {
	return 0
}

func (fna *fakeNumArrayGen) GenerateNums(num int) []int {
	result := make([]int, num)
	for i := 0; i < num; i++ {
		result[i] = i
	}
	return result
}

func generateCityMap() map[string]*generators.CityNode {
	return generators.GenerateCityMap(testingMasks, testingCityNames)
}

func generateGame() *Game {
	cityMap := generateCityMap()
	alienLocations := map[string]string{testingAlien: testingCityNames[0]}
	cityMap[testingCityNames[0]].Aliens = append(cityMap[testingCityNames[0]].Aliens, testingAlien)
	return &Game{CityMap: cityMap, AlienLocations: alienLocations, randGen: fakeZeroGenerator}

}

func TestNewGame(t *testing.T) {
	assert := assert.New(t)

	cityMap := generateCityMap()
	aliens := []string{testingAlien}

	game := NewGame(aliens, cityMap, fakeZeroGenerator)

	assert.Equal(1, len(game.AlienLocations))
	assert.Equal((len(testingMasks)-1)*(len(testingMasks[0])-1), len(game.CityMap))

	cn := game.AlienLocations[testingAlien]
	assert.Equal(testingAlien, game.CityMap[cn].Aliens[0])
}

func TestMove(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		move map[string]int
		cn   string
	}{
		{
			map[string]int{testingAlien: east},
			testingCityNames[1],
		},
		{
			map[string]int{testingAlien: south},
			testingCityNames[3],
		},
		{
			map[string]int{testingAlien: west},
			testingCityNames[2],
		},
		{
			map[string]int{testingAlien: north},
			testingCityNames[0],
		},
	}
	game := generateGame()

	for _, tt := range tests {
		game.MakeMove(tt.move)
		assert.Equal(tt.cn, game.AlienLocations[testingAlien])
	}
}

func TestCheckAndDestroy(t *testing.T) {
	assert := assert.New(t)

	game := generateGame()
	anotherAlien := generators.AlienNames[1]
	game.AlienLocations[anotherAlien] = testingCityNames[3]
	game.CityMap[testingCityNames[3]].Aliens = append(game.CityMap[testingCityNames[3]].Aliens, anotherAlien)

	move := map[string]int{testingAlien: east, anotherAlien: north}
	game.MakeMove(move)
	game.CheckAndDestroy()
	assert.Equal(0, len(game.AlienLocations))
}
