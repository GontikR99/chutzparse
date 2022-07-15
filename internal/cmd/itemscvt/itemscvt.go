package main

import (
	"compress/gzip"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/gontikr99/chutzparse/internal/eqspec"
	"os"
	"sort"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <trie/listing> <filename>\n", os.Args[0])
		return
	}
	dbFile, err := os.Open(os.Args[2])
	if err != nil {
		panic(err)
	}
	defer dbFile.Close()
	lineFile, err := gzip.NewReader(dbFile)
	if err != nil {
		panic(err)
	}
	defer lineFile.Close()
	csvReader := csv.NewReader(lineFile)
	csvReader.Read() // get rid of header
	switch os.Args[1] {
	case "listing":
		var itemNames []string
		for fields, _ := csvReader.Read(); fields != nil; fields, _ = csvReader.Read() {
			itemName := fields[1]
			itemNames = append(itemNames, itemName)
		}
		sort.Sort(eqspec.LexOrderIgnoreCase(itemNames))
		fmt.Println("package eqspec")
		fmt.Println()
		fmt.Println("var everquestItems=[]string{")
		for _, name := range itemNames {
			v, _ := json.Marshal(name)
			fmt.Println("    " + string(v) + ",")
		}
		fmt.Println("}")
	case "trie":
		itemTrie := eqspec.NewItemTrie()
		for fields, _ := csvReader.Read(); fields != nil; fields, _ = csvReader.Read() {
			itemName := fields[1]
			itemTrie = itemTrie.With(itemName)
		}
		cTrie := itemTrie.Compress()

		fmt.Println("package eqspec")
		fmt.Println("var BuiltTrie=CompressedItemTrie{")
		fmt.Print("    Transitions: CompressedItemTrieTransitions{")
		outIdx := 0
		for _, val := range cTrie.Transitions {
			if outIdx%8 == 0 {
				fmt.Println()
				fmt.Print("        ")
			}
			outIdx++
			fmt.Printf("0x%x, ", uint64(val))
		}
		fmt.Println("    },")
		fmt.Print("    Accepts: []int{")
		outIdx = 0
		for _, val := range cTrie.Accepts {
			if outIdx%8 == 0 {
				fmt.Println()
				fmt.Print("        ")
			}
			outIdx++
			fmt.Printf("0x%x, ", val)
		}
		fmt.Println("    },")
		fmt.Println("}")
	default:
		fmt.Fprintln(os.Stderr, "Unknown subcommand "+os.Args[1])
	}
}
