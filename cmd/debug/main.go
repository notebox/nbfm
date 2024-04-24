package main

import (
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/notebox/nbfm/pkg/nav"
)

func main() {
	wd, err := nav.NewWorkingDir()
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}

	fmt.Printf("%+v", wd)
	dir, err := nav.ReadDirFiles(wd.Path, true)
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}

	fmt.Printf("%+#v", dir)
}
