alieninvasion is a small program simulatoring Alien Invasion Game.
# Get the program
```
go get -v github.com/hatricker/alieninvasion
```

# Compile the program
```
cd $GOPATH/src/github.com/hatricker/alieninvasion
make alien
```
It will generate a binary file called *alieninvasion* under *./bin* folder

# Run the program
```
Usage: ./bin/alieninvasion [-na <number of aliens> -mx <X> -my <Y> -mapfile <input map file> -output <output map file>]
  -mapfile string
    	Input map file
  -mx int
    	size of x-coordinate of map matrix
  -my int
    	size of y-coordinate of map matrix
  -na int
    	Number of Aliens (default 2)
  -nm int
    	Number of Moves (default 10000)
  -output string
    	output file to dump the map info
```

### Explanation about the flags

* -mapfile : input map file which defines all the cities and paths connecting each other
* -mx : when **mapfile** is not provided, the program will generate a matrix map automatically. *mx* defines size of x-coordinate
* -my : when **mapfile** is not provided, the program will generate a matrix map automatically. *my* defines size of y-coordinate
* -na : number of aliens in the game.
* -nm : maximum possible moves in the game. The default value is 10000
* -output : output file name where the generated map is dumped to

### A few examples

- To generate a map file
```
./bin/alieninvasion -mx 8 -my 8 -output worldmap.txt  

# The generated map is a 8x8 matrix. Each node on the matrix could have neighbor(s) in four directions
```

- To run the game with a pre-generated map. Note, there are a few pre-generated maps under *./map* folder.
```
./bin/alieninvasion -mapfile worldmap.txt -na 10 -nm 100 

# There will be 10 aliens scatted on the map randomly and the maximum possible moves are 100
```

- To run the game with an automatically generated map
```
./bin/alieninvasion -mx 7 -my 6 -na 10 -nm 100 

# There will be 10 aliens randomly scatted on a 7x6 automatically created map. The maximum possible moves are 100
```

# Run tests
```
make test
```


