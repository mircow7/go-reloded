package main

import (
	"fmt"
	"os"
)

func main() {
	// Il faut exactement 2 arguments : le fichier d'entrée et le fichier de sortie
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "usage: go run . sample.txt result.txt")
		os.Exit(1)
	}

	// On lit le fichier d'entrée
	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Erreur de lecture :", err)
		os.Exit(1)
	}

	// On corrige le texte
	result := ProcessText(string(data))

	// On écrit le résultat dans le fichier de sortie
	err = os.WriteFile(os.Args[2], []byte(result), 0644)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Erreur d'écriture :", err)
		os.Exit(1)
	}
}
