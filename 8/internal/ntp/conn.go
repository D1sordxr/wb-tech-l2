package ntp

import (
	"time"

	"github.com/beevik/ntp"
)

const url = "0.beevik-ntp.pool.ntp.org"

type Service struct{}

func (*Service) GetTime() (time.Time, error) {
	return ntp.Time(url)
}
