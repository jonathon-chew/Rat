# CHANGELOG

	## v0.2.2 

		### UPDATES
		1. wasn't dealing with files in caps as a previous command was forcing it to lower, so file could not be found successfully
		1. matches are now 1 or more and IPs are cyan to seperate from RAG rating time periods
		1. README clean up

	## v0.2.1 

		### NEW
		1. ping now reads from stdout in real time instead of waiting for the command to finish

	## v0.2.0 

		### NEW
		1. adding traceroute, ip works, time needs some work

	## v0.1.6 

		### NEW
		1. can now run ping and arp with coloured outputs for ipaddress, time (coloured coded by length of time of ping) and name, ipadress and macaddress repsectively

	## v0.1.5 

		### NEW
		1. converting from os.ReadFile to reading a file buffer
		### UPDATES
		1. upgrading Aphrodite so only outputs ANSI if output is a terminal
		1. scripts to match other repos for some uniformality
		### DELETES
		1. using file pointer as spanning over boundaries can mean tokens won't work as well, reversing previous push

	## v0.1.4 

		### UPDATES
		1. minor performance updates

	## v0.1.3 

		### UPDATES
		1. version number

	## v0.1.2 

		### REFACTORS
		1. remove fmt.Sprintf a[And now takes 33% of the time!

	## v0.1.1 

		### REFACTORS
		1. nonPossessiveComman (apporstophe) and logic when calling it; new: python and js tokens added null/import/from

	## v0.1.0 

		### NEW
		1. flush buffer function to reduce print calls to speed up the progress of the programme; update: tokens for langauges; refactor: large refactor to increase performance and better handle different variations of comment
		### DELETES
		1. achieved todos

	## v0.0.2 

		### NEW
		1. CICD script for updating and checking!
		### UPDATES
		1. add golang to list of acceptable languages prior to implimentation; update: CICD script to better reflect processes
		1. version number; update script to make work with cli to check all flags are captured for help file