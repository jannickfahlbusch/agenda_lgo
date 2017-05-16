package main

import (
	"flag"
	"log"

	"gitlab.com/jannickfahlbusch/agenda_lgo"
)

var (
	authFilePath string
	out          string
)

func init() {
	flag.StringVar(&authFilePath, "a", ".auth", "Path to the authentication-file")
	flag.StringVar(&out, "o", "out", "Path to the directory where the files should be stored, must exist")
}

func main() {
	flag.Parse()

	lgo := agenda_lgo.NewLGO(authFilePath, out)
	err := lgo.Login()
	if err != nil {
		log.Fatal(err)
	}

	documentList, err := lgo.FetchDocumentList()

	for _, doc := range documentList {
		lgo.SaveDocument(doc)
	}

}
