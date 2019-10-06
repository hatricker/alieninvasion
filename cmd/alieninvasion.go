package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/hatricker/alieninvasion/games"
	"github.com/hatricker/alieninvasion/generators"
)

func main() {
	var (
		numAliens   = flag.Int("na", 2, "Number of Aliens")
		numMoves    = flag.Int("nm", 10000, "Number of Moves")
		cityMatrixX = flag.Int("mx", 0, "City Map matrix x")
		cityMatrixY = flag.Int("my", 0, "City Map matrix x")
		mapFile     = flag.String("mapfile", "", "Map of the world")
		outputFile  = flag.String("output", "", "output file to dump the map info")
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: go run cmd/alieninvasion.go [-na <number of aliens> -mx <X> -my <Y> -mapfile <input map file> -output <output map file>]\n")

		flag.PrintDefaults()
	}
	flag.Parse()

	//when outputFile is given, just dump the generated city map
	if *outputFile != "" {
		cityMap, err := generateMap(*cityMatrixX, *cityMatrixY)
		if err != nil {
			log.Fatalf("cannot generate map, %v", err)
		}
		dumpMapIntoFile(cityMap, *outputFile)
		return
	}
	if (*numMoves) <= 0 || (*numMoves) > 10000 {
		log.Fatalln("Number of Moves must be within 1-10000")
	}

	if (*numAliens) <= 0 {
		log.Fatalln("Number of Aliens must be greater than 0")
	}

	if err := playGame(*mapFile, *numAliens, *numMoves, *cityMatrixX, *cityMatrixY); err != nil {
		log.Fatalf("Error happened when running the game: %v", err)
	}
}

func generateMap(x, y int) (map[string]*generators.CityNode, error) {
	if x == 0 || y == 0 {
		return nil, fmt.Errorf("need to provide both city matrix x and y")
	}
	masks, err := generators.GenerateDirectionMask(x, y, generators.RandNumGenerator)
	if err != nil {
		return nil, fmt.Errorf("cannot generate city map matrix masks, %v", err)
	}
	cityNames, err := generators.GenerateCityNames(generators.RandNumArrGenerator, x*y)
	if err != nil {
		return nil, fmt.Errorf("cannot generate city names, %v", err)
	}
	return generators.GenerateCityMap(masks, cityNames), nil
}

func dumpMapIntoFile(cm map[string]*generators.CityNode, fileName string) error {
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	generators.GenerateMapFile(cm, f)
	return nil
}

//Obtain the map either by generating it on the fly or taking from a local file,
//then start the game
func playGame(mapFile string, numAliens, numMoves, x, y int) error {
	var (
		cityMap map[string]*generators.CityNode
		err     error
	)

	//if no map file is provided, generate a map automatically
	if mapFile == "" {
		if cityMap, err = generateMap(x, y); err != nil {
			return fmt.Errorf("cannot generate map, %v", err)
		}
	} else {
		//read from the input file
		f, err := os.Open(mapFile)
		if err != nil {
			return err
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		cityMap = generators.GenerateCityMapFromSteam(scanner, ' ')
	}

	log.Println("Obtained city map...")
	printCityMap(cityMap)

	aliens, err := generators.GenerateAlienNames(generators.RandNumArrGenerator, numAliens)
	if err != nil {
		return fmt.Errorf("cannot generate alien names, %v", err)
	}

	log.Printf("Generated aliens: %s", strings.Join(aliens, " "))

	g := games.NewGame(aliens, cityMap, generators.RandNumGenerator)
	log.Println("Game starting...")
	g.StartGame(numMoves)

	log.Println("City map at the end of game...")
	printCityMap(g.CityMap)
	return nil
}

//print city map to the stdout
func printCityMap(cm map[string]*generators.CityNode) {
	var b bytes.Buffer
	generators.GenerateMapFile(cm, &b)
	fmt.Printf("%s", b.String())
}
