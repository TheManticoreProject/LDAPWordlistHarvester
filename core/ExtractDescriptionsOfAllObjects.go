package core

import (
	"fmt"
	"strings"

	"github.com/TheManticoreProject/Manticore/logger"
	"github.com/TheManticoreProject/Manticore/network/ldap"
)

// ExtractDescriptionsOfAllObjects extracts the descriptions of all objects from the LDAP server.
// It uses the ldap.Session to query the LDAP server and get the descriptions of all objects.
// It then adds the descriptions to the wordlist.
func ExtractDescriptionsOfAllObjects(ldapSession ldap.Session, wordlist *Wordlist) {
	logger.Info("Extracting descriptions of all objects from LDAP...")

	query := "(objectClass=*)"
	ldapResults, err := ldapSession.QueryWholeSubtree("", query, []string{"description"})
	if err != nil {
		logger.Error(fmt.Sprintf("Error querying LDAP: %s", err))
		return
	}

	candidates := make([]string, 0)
	for _, entry := range ldapResults {
		// Process description
		descriptions := entry.GetEqualFoldAttributeValues("description")
		if len(descriptions) > 0 {
			// If description is a slice
			words := strings.Split(strings.Join(descriptions, " "), " ")
			candidates = append(candidates, words...)
		} else {
			// If description is a single string
			for _, description := range descriptions {
				words := strings.Split(description, " ")
				candidates = append(candidates, words...)
			}
		}
	}

	nwords := wordlist.AddUniqueWords(candidates)

	logger.Info(fmt.Sprintf(" └──[+] Added %d unique words to wordlist.", nwords))
}
