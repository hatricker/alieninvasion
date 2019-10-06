package generators

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math/rand"
	"strings"
	"time"
)

var (
	//CityNames holds a set of city names that the program is going to use.
	CityNames = []string{"Canandaigua", "PaintedHills", "Eschbach", "Saybrook", "Hahira", "Oostburg", "Kappa", "GranvilleSouth", "Manahawkin", "LeRaysville",
		"RanchoMirage", "Zena", "Manhattan", "Wildomar", "Burnettsville", "Protivin", "Pioneer", "Whitakers", "GrandMound", "Greilickville", "Kipton",
		"Sedan", "Edinburg", "Antelope", "Tecumseh", "ConesusHamlet", "LakeCassidy", "Zelienople", "LeadvilleNorth", "MapleHeights", "Syosset",
		"SixShooterCanyon", "Gonvick", "Westmont", "AltoBonitoHeights", "Abeytas", "Herington", "Fulford", "Fairgarden", "ArrowheadSprings",
		"ElkHorn", "MorroBay", "RockCity", "Alamo", "ClearviewAcres", "JAARS", "Taunton", "Funkley", "BessemerBend", "RiverRidge", "Yadao",
		"Dagangjiu", "Hongpan Xiang", "Jincun", "Ganbao", "Dahuang", "Yecitang", "Sanxingzhen", "Xinzhuangcun", "Donghuanglou", "Xiaojiaao", "Qiaodong",
		"Yangqiaocun", "Dawan", "Shengdianzhuang", "Lüdazhuang", "Langdaling", "Minzhu", "Mirantarium", "Jishui", "Maidigou", "Lengshuiwan", "Fuche",
		"Bamian", "Guanyinmiao", "Shuangpeixia", "Dushanwu", "Baiyang", "Daoziba", "FujiaBeigou", "Ganzhou", "Xishan", "Dabali", "Maykhutu", "Huakengwu",
		"Miaobei", "Shangyang", "Xiaoqiao", "Huangyaoguan", "Tongziwo", "Wutongmiao", "Huangtuling", "Xintan", "Ziyipu", "Xizhuangtou", "Tongjunzhuang",
		"Loufanggou", "Zhengzhuang", "Peitaiho", "Hengkeng"}

	//AlienNames hodls a set of alien names that the program is going to use
	AlienNames = []string{"Yalmimin", "Raxomalik", "Degir", "Cfuujaban", "Borger",
		"Oleniflax", "Vuludha", "NeFlav Yucholl", "Mane", "Nbaaana",
		"Ruavu Strogonar", "Araime Fallapadax", "Proog Wontwoon", "Salah", "Ndidi",
		"Olgsivoor", "Ghuyot", "Kragitur", "Zumbal", "Zidane",
		"Luvendav", "Tamer", "Ruavu", "Ofnatsuza", "Cleayomaar"}

	//ErrReqTooLarge is returned when the requsted number of aliens or cities is beyond the size of predefined list
	ErrReqTooLarge = fmt.Errorf("requested size is too large")
	//ErrInvalidInput is returned when the input is invalid
	ErrInvalidInput = fmt.Errorf("invalid input")

	East  = 1
	West  = 2
	North = 4
	South = 8
	//DirectionBitMap is a map of 4 directions
	DirectionBitMap = []int{East, West, North, South}

	RandNumGenerator    = &RandNumGen{}
	RandNumArrGenerator = &RandNumArrayGen{}
)

//NumGen defines the interface of number generator
type NumGen interface {
	//n is range
	GenerateNum(n int) int
}

//RandNumGen implements NumGen interface
type RandNumGen struct {
}

//GenerateNum is the detailed implementation of RandNumGen
func (rn *RandNumGen) GenerateNum(n int) int {
	return rand.Intn(n)
}

//NumArrayGen is the common interface of number generator
type NumArrayGen interface {
	GenerateNums(num int) []int
}

//RandNumArrayGen provides functionality of generating a list of random numbers
type RandNumArrayGen struct {
}

//CityNode defines a city node in the whole map
type CityNode struct {
	Name                     string
	East, West, North, South *CityNode
	Aliens                   []string
}

//NewRandNumGen returns a RandNumArrayGen object
func NewRandNumGen() *RandNumArrayGen {
	return &RandNumArrayGen{}
}

//GenerateNums implements NumArrayGen interface
func (rna *RandNumArrayGen) GenerateNums(num int) []int {
	if num <= 0 {
		return nil
	}
	rand.Seed(time.Now().UnixNano())
	return rand.Perm(num)
}

//GenerateNames returns numbers of random names from the predefined list
func GenerateNames(generator NumArrayGen, nameList []string, num int) ([]string, error) {
	if num > len(nameList) {
		return nil, ErrReqTooLarge
	}

	nums := generator.GenerateNums(num)
	names := make([]string, 0, num)
	for _, num := range nums {
		names = append(names, nameList[num])
	}
	return names, nil
}

//GenerateCityNames returns a list of city names
func GenerateCityNames(generator NumArrayGen, num int) ([]string, error) {
	return GenerateNames(generator, CityNames, num)
}

//GenerateAlienNames returns a list of alien names
func GenerateAlienNames(generator NumArrayGen, num int) ([]string, error) {
	return GenerateNames(generator, AlienNames, num)
}

//Generate and encode the 4-direction into an integer. It needs to consider
//its neighbors' values as well. One of the input is neighbors' values in
//the direction of {"west", "east", "south", "north"}, because the node's direction
//values are encoded in the order of {"east", "west", "north", "south"}
//If its neighbor's value is already set, honor that value. If not, randomly decide
//whether to have a path to its neighbor
func getDirectionValue(neighbors [4]int, rg NumGen) int {
	value := 0
	for i, val := range neighbors {
		//-1 means randomly deciding having path or not to that direction
		if val == -1 {
			if rg.GenerateNum(2) != 0 {
				value = value | DirectionBitMap[i]
			}
		} else {
			value = value | val
		}
	}
	return value
}

//GenerateDirectionMask generates a matrix which represents the city map
//The value at coordinate (i, j) show its neighbors at four directions
//Check the predefines variable DirectionBitMap on top to see how the directions
//are represented. For example, when value is 9 (East|South), it has neighbors
//in the east and south. The generated masks is like 2-dimension array below
//{0,     0,            0,             0},
//{0,     East,         West,          South},
//{0,     South,        0,             North | South},
//{0,     East | North, West | East,   West | North},
//Note, the first row and first column are there to help generate direction masks
func GenerateDirectionMask(x, y int, rg NumGen) ([][]int, error) {
	if x*y > len(CityNames) {
		return nil, ErrReqTooLarge
	}
	if x == 0 || y == 0 {
		return nil, ErrInvalidInput
	}
	//Make an extra rown and an extra column in order to make
	//direction calculation easier. The real matrix starts at [1,1]
	m := make([][]int, x+1)
	m[0] = make([]int, y+1)
	for i := 1; i <= x; i++ {
		m[i] = make([]int, y+1)
		for j := 1; j <= y; j++ {
			w, n := 0, 0
			//has neighbor in the west?
			if m[i][j-1]&East > 0 {
				w = West
			}
			//has neighbor in the north?
			if m[i-1][j]&South > 0 {
				n = North
			}

			neighbors := [4]int{-1, w, n, -1}
			//last column does not have east neighbor
			if j == y {
				neighbors[0] = 0
			}
			//last row does not have south neighbor
			if i == x {
				neighbors[3] = 0
			}
			m[i][j] = getDirectionValue(neighbors, rg)
		}
	}

	return m, nil
}

//GenerateCityMap returns city nodes map
//It takes a mask generated by GenerateDirectionMask above and list of city names
//to build a map which represents the map in memory
func GenerateCityMap(mask [][]int, cityNames []string) map[string]*CityNode {
	if len(mask) == 0 {
		return nil
	}
	x, y := len(mask)-1, len(mask[0])-1

	// provided city names must no less than the matrix(map) element
	if len(cityNames) < x*y {
		return nil
	}

	cm := make(map[string]*CityNode)

	for i := 0; i < x*y; i++ {
		cn := &CityNode{Name: cityNames[i], Aliens: make([]string, 0, 20)}
		cm[cityNames[i]] = cn
	}
	//coordinate starts at [1:1]
	for i := 1; i <= x; i++ {
		for j := 1; j <= y; j++ {
			maskVal := mask[i][j]
			if maskVal == 0 {
				continue
			}

			currNameInd := (i-1)*y + j - 1
			currNode := cm[cityNames[currNameInd]]
			if maskVal&East > 0 {
				ec := cityNames[currNameInd+1]
				currNode.East = cm[ec]
			}
			if maskVal&West > 0 {
				wc := cityNames[currNameInd-1]
				currNode.West = cm[wc]
			}
			if maskVal&North > 0 {
				nc := cityNames[currNameInd-y]
				currNode.North = cm[nc]
			}
			if maskVal&South > 0 {
				sc := cityNames[currNameInd+y]
				currNode.South = cm[sc]
			}
		}
	}
	return cm
}

//GenerateCityMapFromSteam reads city map from stream
//It can be from a real file, or a string stream for testing purpose
func GenerateCityMapFromSteam(scanner *bufio.Scanner, splitter rune) map[string]*CityNode {
	cm := make(map[string]*CityNode)

	for scanner.Scan() {
		line := scanner.Text()

		//when input is string, it's by default tokenizing by whitespace
		//so, add a special splitter helper to allow specifying other character
		splitFn := func(c rune) bool {
			return c == splitter
		}

		tokens := strings.FieldsFunc(line, splitFn)
		if len(tokens) > 0 {
			cityName := tokens[0]
			if _, ok := cm[cityName]; !ok {
				cm[cityName] = &CityNode{Name: cityName, Aliens: make([]string, 0, 20)}
			}

			for _, token := range tokens[1:] {
				directStrs := strings.Split(token, "=")
				if len(directStrs) != 2 {
					log.Panic("invalid direction map")
				}
				city, ok := cm[directStrs[1]]
				if !ok {
					city = &CityNode{Name: directStrs[1], Aliens: make([]string, 0, 20)}
					cm[directStrs[1]] = city
				}
				switch directStrs[0] {
				case "east":
					cm[cityName].East = city
				case "west":
					cm[cityName].West = city
				case "north":
					cm[cityName].North = city
				case "south":
					cm[cityName].South = city
				}
			}

		}
	}
	if err := scanner.Err(); err != nil {
		log.Panic("cannot parse the input stream properly")
	}
	return cm
}

//GenerateMapFile writes the map info into output source
func GenerateMapFile(cm map[string]*CityNode, w io.Writer) {
	bufWriter := bufio.NewWriter(w)

	for city, node := range cm {
		coordinates := make([]string, 0, 6)
		coordinates = append(coordinates, city)
		if node.East != nil {
			coordinates = append(coordinates, "east="+node.East.Name)
		}
		if node.West != nil {
			coordinates = append(coordinates, "west="+node.West.Name)
		}
		if node.North != nil {
			coordinates = append(coordinates, "north="+node.North.Name)
		}
		if node.South != nil {
			coordinates = append(coordinates, "south="+node.South.Name)
		}
		if len(coordinates) == 1 {
			continue
		}
		coordinates = append(coordinates, "\n")
		bufWriter.WriteString(strings.Join(coordinates, " "))
	}
	bufWriter.Flush()
}
