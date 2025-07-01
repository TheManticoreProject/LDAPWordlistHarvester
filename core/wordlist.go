package core

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/TheManticoreProject/Manticore/logger"
)

// Wordlist is a struct that contains the path to the wordlist file and the list of words.
type Wordlist struct {
	Path     string
	Wordlist []string
}

// NewWordlist creates a new Wordlist struct with the given path.
// It initializes the Wordlist struct with the given path and an empty list of words.
func NewWordlist(wordlistPath string) *Wordlist {
	return &Wordlist{
		Path:     wordlistPath,
		Wordlist: make([]string, 0),
	}
}

// AddUniqueWords adds unique words from the provided candidates slice to the Wordlist.
//
// This method iterates over each word in the candidates slice and checks if it already exists
// in the Wordlist. If a word does not exist in the Wordlist, it is added to the Wordlist.
//
// Parameters:
//
//	candidates []string: A slice of words to be added to the Wordlist.
//
// Example usage:
//
//	wordlist := NewWordlist("path/to/wordlist.txt")
//	wordlist.Wordlist = []string{"existingWord1", "existingWord2"}
//	candidates := []string{"newWord1", "existingWord1", "newWord2"}
//	wordlist.AddUniqueWords(candidates)
//	// wordlist.Wordlist now contains: ["existingWord1", "existingWord2", "newWord1", "newWord2"]
//
// Note:
//   - This method does not add duplicate words to the Wordlist.
//   - The comparison is case-sensitive, meaning "Word" and "word" are considered different.
func (w *Wordlist) AddUniqueWords(candidates []string) int {
	addedWords := 0
	for _, word := range candidates {
		found := false
		for _, existing := range w.Wordlist {
			if word == existing {
				found = true
				break
			}
		}
		if !found {
			for _, existing := range w.Wordlist {
				if strings.EqualFold(word, existing) {
					found = true
					break
				}
			}
		}
		if !found {
			w.Wordlist = append(w.Wordlist, word)
			addedWords++
		}
	}
	return addedWords
}

// Write writes the wordlist to the file specified by the Path field of the Wordlist struct.
// It creates a new file if it does not exist, or truncates the file if it does exist.
//
// Returns an error if the file cannot be created or if there is an error during writing.
//
// Example usage:
//
//	wordlist := NewWordlist("path/to/wordlist.txt")
//	wordlist.Wordlist = []string{"word1", "word2", "word3"}
//	err := wordlist.Write()
//	if err != nil {
//	    fmt.Printf("Error writing wordlist: %v\n", err)
//	}
func (w *Wordlist) Write() error {
	file, err := os.Create(w.Path)
	if err != nil {
		return fmt.Errorf("failed to create wordlist file: %v", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, word := range w.Wordlist {
		_, err := writer.WriteString(word + "\n")
		if err != nil {
			return fmt.Errorf("failed to write to wordlist: %v", err)
		}
	}
	writer.Flush()
	logger.Info(fmt.Sprintf("Wordlist written to: %s (%d words)", w.Path, len(w.Wordlist)))

	return nil
}
