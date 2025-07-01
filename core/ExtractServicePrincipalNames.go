package core

import (
	"fmt"
	"strings"

	"github.com/TheManticoreProject/Manticore/logger"
	"github.com/TheManticoreProject/Manticore/network/ldap"
)

// ExtractServicePrincipalNames extracts the service principal names from the LDAP server.
// It uses the ldap.Session to query the LDAP server and get the service principal names.
// It then adds the service principal names to the wordlist.
func ExtractServicePrincipalNames(ldapSession ldap.Session, wordlist *Wordlist) {
	logger.Info("Extracting service principal names from LDAP...")

	query := "(objectClass=*)"
	ldapResults, err := ldapSession.QueryWholeSubtree("", query, []string{"servicePrincipalName"})
	if err != nil {
		logger.Error(fmt.Sprintf("Error querying LDAP: %s", err))
		return
	}

	candidates := make([]string, 0)
	for _, entry := range ldapResults {
		// Process service principal name
		servicePrincipalNames := entry.GetEqualFoldAttributeValues("servicePrincipalName")
		if len(servicePrincipalNames) > 0 {
			// If service principal name is a slice
			words := strings.Split(strings.Join(servicePrincipalNames, "/"), "/")
			candidates = append(candidates, words...)
		} else {
			// If service principal name is a single string
			for _, servicePrincipalName := range servicePrincipalNames {
				words := strings.Split(servicePrincipalName, "/")
				candidates = append(candidates, words...)
			}
		}
	}

	nwords := wordlist.AddUniqueWords(candidates)

	logger.Info(fmt.Sprintf(" └──[+] Added %d unique words to wordlist.", nwords))
}
