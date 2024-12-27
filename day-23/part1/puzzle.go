package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
)

func FindCliques(prevCliques []map[int]bool, cliqueSize, verticies int, adjacent map[int][]int) []map[int]bool {
	// filter out prevCliques
	cliqueCandidates := make([]map[int]bool, len(prevCliques))
	for _, prevClique := range prevCliques {
		cliqueCandidate := make(map[int]bool)
		for vertex := range prevClique {
			if len(adjacent[vertex])+1 < cliqueSize {
				fmt.Printf("Clique candidate %v has not enough adjacent verticies: %d\n", prevClique, len(adjacent[vertex])+1)
			}
			cliqueCandidate[vertex] = true
		}
		if len(cliqueCandidate)+1 == cliqueSize {
			cliqueCandidates = append(cliqueCandidates, cliqueCandidate)
		}
	}

	for _, candidateClique := range cliqueCandidates {
		// try to expand prev clique by adding one node to it
		for vertex := range verticies {
			// Skip if already in the clique
			if candidateClique[vertex] {
				continue
			}

			// Check if adjacent to every other vertex in a clique
			adjacent2clique := true
			for cVertex := range candidateClique {
				if !slices.Contains(adjacent[cVertex], vertex) {
					adjacent2clique = false
					break
				}
			}
			if !adjacent2clique {
				continue
			}

			candidateClique[vertex] = true
		}
	}

	// Filter out candidateCliques which haven't reached the cliqueSize
	result := make([]map[int]bool, 0)
	for _, candidateClique := range cliqueCandidates {
		if len(candidateClique) == cliqueSize {
			result = append(result, candidateClique)
		}
	}

	return result
}

func main() {
	file, err := os.Open(os.Args[1])
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
	fmt.Printf("Number of verticies: %d\n", verticies)

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

	// Filter out all the triple that don't contain a node with a name that starts with 't'
	setsOf3WithTs := make(map[[3]int]bool)
	for triple := range setsOf3 {
		found := false
		for _, id := range triple {
			if strings.HasPrefix(id2name[id], "t") {
				found = true
				break
			}
		}
		if found {
			setsOf3WithTs[triple] = true
		}
	}

	// List all triples
	for triple := range setsOf3WithTs {
		fmt.Printf("%s %s %s\n", id2name[triple[0]], id2name[triple[1]], id2name[triple[2]])
	}
	fmt.Printf("There are %d triples that contain a node that starts with t\n", len(setsOf3WithTs))

	// Find the maxiumum possible maximum clique size
	maxCliqueSize := 4
	candidateNodes := []int{}
	for true {
		fmt.Printf("Trying clique size %d\n", maxCliqueSize)
		candidateNodes = []int{}
		for v := range verticies {
			if len(adjacent[v]) >= maxCliqueSize-1 {
				candidateNodes = append(candidateNodes, v)
			}
		}

		// Clique with maxCliqueSize impossible
		if len(candidateNodes) < maxCliqueSize {
			maxCliqueSize--
			break
		}

		maxCliqueSize++
	}
	fmt.Printf("The maximum clique size is %d\n", maxCliqueSize)

	// Check all sets of maxCliqueSize nodes from candidateNodes
	for cliqueSize := maxCliqueSize; cliqueSize > 1; cliqueSize-- {
	}

}
