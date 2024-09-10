package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sync"
)

// read a file from a filepath and return a slice of bytes
func readFile(filePath string) ([]byte, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", filePath, err)
		return nil, err
	}
	return data, nil
}

// sum all bytes of a file
func sum(filePath string) (int, error) {
	data, err := readFile(filePath)
	if err != nil {
		return 0, err
	}

	_sum := 0
	for _, b := range data {
		_sum += int(b)
	}

	return _sum, nil
}
 

func chamadaConcorrente(path string, grupando *sync.WaitGroup, exclusao *sync.Mutex, totalSum *int64, sums map[int][]string) {
	defer grupando.Done()		// é chamada quanto a função terminar e libera

	_sum, err := sum(path)
	
	if err != nil {
		return
	}

	exclusao.Lock()
	defer exclusao.Unlock()			// libera quando terminar

	*totalSum += int64(_sum)
	sums[_sum] = append(sums[_sum], path)

}
// print the totalSum for all files and the files with equal sum
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <file1> <file2> ...")
		return
	}

	//parte em serial
	var totalSum int64
	sums := make(map[int][]string)

	//declarando
	grupando sync.WaitGroup   	// esperar todas terminarem antes de imprimir
	exclusao sync.Mutex			// garante exclusao mutua

	for _, path := range os.Args[1:] {
		grupando.Add(1)
		go chamadaConcorrente(path, grupando, exclusao, totalSum, sums)
	}

	grupando.Wait()

	fmt.Println(totalSum)


	for sum, files := range sums {
		if len(files) > 1 {
			fmt.Printf("Sum %d: %v\n", sum, files)
		}
	}
}