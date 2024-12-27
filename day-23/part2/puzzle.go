package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"slices"
	"strings"
)

func contains2D(slice [][]int, target []int) bool {
	for _, s := range slice {
		if reflect.DeepEqual(s, target) {
			return true
		}
	}
	return false
}

func FindCliques(prevCliques [][]int, cliqueSize, verticies int, adjacent map[int][]int) [][]int {
	// filter out prevCliques
	cliqueCandidates := make([][]int, 0, len(prevCliques))
	for _, prevClique := range prevCliques {
		cliqueCandidate := make([]int, 0)
		for _, vertex := range prevClique {
			if len(adjacent[vertex])+1 < cliqueSize {
				log.Printf("Clique candidate %v has not enough adjacent verticies: %d\n", prevClique, len(adjacent[vertex])+1)
			}
			cliqueCandidate = append(cliqueCandidate, vertex)
		}
		if len(cliqueCandidate)+1 == cliqueSize {
			cliqueCandidates = append(cliqueCandidates, cliqueCandidate)
		}
	}

	log.Printf("-----------Clique candidates-----------\n")
	log.Printf("FindClique(%d)\n", cliqueSize)
	log.Printf("Amount: %d\n", len(cliqueCandidates))
	for _, clique := range cliqueCandidates {
		log.Printf("Clique: %v\n", clique)
	}

	results := make([][]int, 0)
	for _, candidateClique := range cliqueCandidates {
		// try to expand prev clique by adding one node to it
		for vertex := range verticies {
			// Skip if already in the clique
			if slices.Contains(candidateClique, vertex) {
				continue
			}

			// Check if adjacent to every other vertex in a clique
			adjacent2clique := true
			for _, cVertex := range candidateClique {
				if !slices.Contains(adjacent[cVertex], vertex) {
					adjacent2clique = false
					break
				}
			}
			if !adjacent2clique {
				continue
			}

			candidateClique = append(candidateClique, vertex)
			break
		}
		slices.Sort(candidateClique)
		if contains2D(results, candidateClique) {
			continue
		}
		results = append(results, candidateClique)
	}

	return results
}

func main() {

	verbose := flag.Bool("verbose", false, "Print verbose output")
	flag.Parse()

	log.SetFlags(0)
	if !*verbose {
		log.SetOutput(io.Discard)
	}

	file, err := os.Open(flag.Arg(0))
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	name2id := make(map[string]int, 0)
	id2name := make(map[int]string, 0)
	adjacent := make(map[int][]int)
	edges := make([][2]int, 0)

	nextId := 0
	for scanner.Scan() {
		line := scanner.Text()
		vName, uName, found := strings.Cut(line, "-")
		if !found {
			panic("Invalid input")
		}

		v, exists := name2id[vName]
		if !exists {
			v = nextId
			nextId++

			name2id[vName] = v
			id2name[v] = vName
		}

		u, exists := name2id[uName]
		if !exists {
			u = nextId
			nextId++

			name2id[uName] = u
			id2name[u] = uName
		}

		adjacent[v] = append(adjacent[v], u)
		adjacent[u] = append(adjacent[u], v)
		edges = append(edges, [2]int{u, v})
	}
	verticies := nextId
	log.Printf("Number of verticies: %d\n", verticies)

	// Find 3 interconnected nodes
	setsOf3 := make(map[[3]int]bool)
	for _, edge := range edges {
		u, v := edge[0], edge[1]
		for _, q := range adjacent[u] {
			for _, neighborOfQ := range adjacent[q] {
				if neighborOfQ == v {
					tripleS := []int{u, v, q}
					slices.Sort(tripleS)

					triple := [3]int{tripleS[0], tripleS[1], tripleS[2]}

					setsOf3[triple] = true
				}
			}
		}
	}

	prevCliques := make([][]int, 0, len(setsOf3))
	for set := range setsOf3 {
		u, v, q := set[0], set[1], set[2]
		sortedSet := []int{u, v, q}
		slices.Sort(sortedSet)
		prevCliques = append(prevCliques, sortedSet)
	}

	log.Printf("Cliques of size 3: %d\n", len(prevCliques))
	if *verbose {
		for _, clique := range prevCliques {
			log.Printf("Clique: %v\n", clique)
		}
	}

	cliqueSize := 4
	for true {
		nextCliques := FindCliques(prevCliques, cliqueSize, verticies, adjacent)
		log.Printf("Cliques of size %d: %v\n", cliqueSize, nextCliques)
		if len(nextCliques) == 0 {
			break
		}

		prevCliques = nextCliques
		cliqueSize++
	}

	for _, prevClique := range prevCliques {
		names := make([]string, len(prevClique))
		for i, id := range prevClique {
			names[i] = id2name[id]
		}

		slices.Sort(names)

		password := strings.Join(names, ",")
		fmt.Printf("Len(names): %d\n", len(names))
		fmt.Printf("Password: %s\n", password)
	}
}
