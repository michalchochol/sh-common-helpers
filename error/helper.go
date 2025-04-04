package error

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

func Chk(err error, info string) {
	if err != nil {
		log.Error().Msg(fmt.Sprintf("%s: %s", info, err))
		return
	}
}

func Fatal(err error, info string) {
	if err != nil {
		log.Fatal().Msg(fmt.Sprintf("%s: %s", info, err))
		return
	}
}
