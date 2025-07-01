package main

import (
	"fmt"

	"github.com/TheManticoreProject/LDAPWordlistHarvester/core"
	"github.com/TheManticoreProject/Manticore/logger"
	"github.com/TheManticoreProject/Manticore/network/ldap"
	"github.com/TheManticoreProject/Manticore/windows/credentials"
	"github.com/TheManticoreProject/goopts/parser"
)

var (
	// Configuration
	debug bool

	// Authentication
	authDomain   string
	authUsername string
	authPassword string
	authHashes   string

	// LDAP Connection Settin	mode string
	domainController string
	ldapPort         int
	useLdaps         bool
	useKerberos      bool
	outputFile       string
)

func parseArgs() {
	ap := parser.ArgumentsParser{
		Banner: "LDAPWordlistHarvester - by Remi GASCOU (Podalirius) @ TheManticoreProject - v1.0.0",
	}
	ap.SetOptShowBannerOnHelp(true)
	ap.SetOptShowBannerOnRun(true)

	// Configuration flags
	group_config, err := ap.NewArgumentGroup("Configuration")
	if err != nil {
		fmt.Printf("[error] Error creating ArgumentGroup: %s\n", err)
	} else {
		group_config.NewBoolArgument(&debug, "", "--debug", false, "Debug mode.")
		group_config.NewStringArgument(&outputFile, "-o", "--output", "wordlist.txt", false, "Output file.")
	}
	// LDAP Connection Settings
	group_ldapSettings, err := ap.NewArgumentGroup("LDAP Connection Settings")
	if err != nil {
		fmt.Printf("[error] Error creating ArgumentGroup: %s\n", err)
	} else {
		group_ldapSettings.NewStringArgument(&domainController, "-dc", "--dc-ip", "", true, "IP Address of the domain controller or KDC (Key Distribution Center) for Kerberos. If omitted, it will use the domain part (FQDN) specified in the identity parameter.")
		group_ldapSettings.NewTcpPortArgument(&ldapPort, "-lp", "--ldap-port", 389, false, "Port number to connect to LDAP server.")
		group_ldapSettings.NewBoolArgument(&useLdaps, "-L", "--use-ldaps", false, "Use LDAPS instead of LDAP.")
		group_ldapSettings.NewBoolArgument(&useKerberos, "-k", "--use-kerberos", false, "Use Kerberos instead of NTLM.")
	}
	// Authentication flags
	group_auth, err := ap.NewArgumentGroup("Authentication")
	if err != nil {
		fmt.Printf("[error] Error creating ArgumentGroup: %s\n", err)
	} else {
		group_auth.NewStringArgument(&authDomain, "-d", "--domain", "", true, "Active Directory domain to authenticate to.")
		group_auth.NewStringArgument(&authUsername, "-u", "--username", "", true, "User to authenticate as.")
		group_auth.NewStringArgument(&authPassword, "-p", "--password", "", false, "Password to authenticate with.")
		group_auth.NewStringArgument(&authHashes, "-H", "--hashes", "", false, "NT/LM hashes, format is LMhash:NThash.")
	}

	ap.Parse()
}

func main() {
	parseArgs()

	creds, err := credentials.NewCredentials(authDomain, authUsername, authPassword, authHashes)
	if err != nil {
		fmt.Printf("[error] Error creating credentials: %s\n", err)
		return
	}

	ldapSession := ldap.Session{}
	ldapSession.InitSession(domainController, ldapPort, creds, useLdaps, useKerberos)
	success, err := ldapSession.Connect()
	if !success {
		logger.Warn(fmt.Sprintf("Error performing LDAP search: %s\n", err))
		return
	}

	wordlist := core.NewWordlist(outputFile)

	core.ExtractADSites(ldapSession, wordlist)

	core.ExtractNamesOfAllObjects(ldapSession, wordlist)

	core.ExtractDescriptionsOfAllObjects(ldapSession, wordlist)

	core.ExtractServicePrincipalNames(ldapSession, wordlist)

	wordlist.Write()

	logger.Print("Done")
}
