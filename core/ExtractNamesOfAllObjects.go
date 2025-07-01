package core

import (
	"fmt"
	"strings"

	"github.com/TheManticoreProject/Manticore/logger"
	"github.com/TheManticoreProject/Manticore/network/ldap"
)

// ExtractNamesOfAllObjects extracts the names of all objects from the LDAP server.
// It uses the ldap.Session to query the LDAP server and get the names of all objects.
// It then adds the names to the wordlist.
func ExtractNamesOfAllObjects(ldapSession ldap.Session, wordlist *Wordlist) {
	logger.Info("Extracting names of all objects from LDAP...")

	query := "(objectClass=*)"
	ldapResults, err := ldapSession.QueryWholeSubtree("", query, []string{"name"})
	if err != nil {
		logger.Error(fmt.Sprintf("Error querying LDAP: %s", err))
		return
	}

	candidates := make([]string, 0)
	for _, entry := range ldapResults {
		// Process name
		if len(entry.GetEqualFoldAttributeValues("name")) > 0 { // If name is a slice
			names := make([]string, 0)
			for _, name := range entry.GetEqualFoldAttributeValues("name") {
				if len(name) > 0 {
					names = append(names, name)
				}
			}
			words := strings.Split(strings.Join(names, " "), " ")
			candidates = append(candidates, words...)
		} else {
			// If name is a single string
			words := strings.Split(entry.GetEqualFoldAttributeValues("name")[0], " ")
			candidates = append(candidates, words...)
		}
	}

	nwords := wordlist.AddUniqueWords(candidates)

	logger.Info(fmt.Sprintf(" └──[+] Added %d unique words to wordlist.", nwords))
}
