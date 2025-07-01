package core

import (
	"fmt"
	"strings"

	"github.com/TheManticoreProject/Manticore/logger"
	"github.com/TheManticoreProject/Manticore/network/ldap"
)

// ExtractADSites extracts the AD Sites from the LDAP server.
// It uses the ldap.Session to query the LDAP server and get the AD Sites.
// It then adds the AD Sites to the wordlist.
func ExtractADSites(ldapSession ldap.Session, wordlist *Wordlist) {
	// Extract AD Sites

	logger.Info("Extracting AD Sites from LDAP...")

	query := "(objectClass=site)"
	ldapResults, err := ldapSession.QueryWholeSubtree("", query, []string{"name", "description"})
	if err != nil {
		logger.Error(fmt.Sprintf("Error querying LDAP: %s", err))
		return
	}

	candidates := make([]string, 0)
	for _, entry := range ldapResults {
		// Process description
		if len(entry.GetEqualFoldAttributeValues("description")) > 0 {
			// If description is a slice
			words := strings.Split(strings.Join(entry.GetEqualFoldAttributeValues("description"), " "), " ")
			candidates = append(candidates, words...)
		} else {
			// If description is a single string
			words := strings.Split(entry.GetEqualFoldAttributeValues("description")[0], " ")
			candidates = append(candidates, words...)
		}

		// Process name
		if len(entry.GetEqualFoldAttributeValues("name")) > 0 {
			// If name is a slice
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
