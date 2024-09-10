package main

import (
	"fmt"
	"io"
	"os"
)

// read a file from a filepath and return a slice of fingerprints (partial sums)
func readFile(filePath string) ([]int64, error) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error reading file %s: %v", filePath, err)
		return nil, err
	}
	defer file.Close()

	var chunks []int64
	buffer := make([]byte, 100) // Buffer de 100 bytes, como no código em Java
	for {
		bytesRead, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			return nil, err
		}
		if bytesRead == 0 {
			break
		}
		chunkSum := sum(buffer[:bytesRead])
		chunks = append(chunks, chunkSum)
	}
	return chunks, nil
}

// sum all bytes of a buffer
func sum(buffer []byte) int64 {
	var _sum int64
	for _, b := range buffer {
		_sum += int64(b)
	}
	return _sum
}

// calculate the similarity between two fingerprints
func similarity(base, target []int64) float64 {
	counter := 0
	targetCopy := make([]int64, len(target))
	copy(targetCopy, target)

	for _, v := range base {
		for i, t := range targetCopy {
			if v == t {
				counter++
				targetCopy = append(targetCopy[:i], targetCopy[i+1:]...) // Remove o valor encontrado
				break
			}
		}
	}

	return float64(counter) / float64(len(base))
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <file1> <file2> ...")
		return
	}

	// Cria um mapa para armazenar os fingerprints de cada arquivo
	fileFingerprints := make(map[string][]int64)

	// Calcula o fingerprint de cada arquivo
	for _, path := range os.Args[1:] {
		fingerprint, err := readFile(path)
		if err != nil {
			fmt.Printf("Error processing file %s: %v\n", path, err)
			continue
		}
		fileFingerprints[path] = fingerprint
	}

	// Compara cada par de arquivos
	for i := 0; i < len(os.Args[1:]); i++ {
		for j := i + 1; j < len(os.Args[1:]); j++ {
			file1 := os.Args[i+1]
			file2 := os.Args[j+1]
			fingerprint1 := fileFingerprints[file1]
			fingerprint2 := fileFingerprints[file2]
			similarityScore := similarity(fingerprint1, fingerprint2)
			fmt.Printf("Similarity between %s and %s: %.2f%%\n", file1, file2, similarityScore*100)
		}
	}
}
