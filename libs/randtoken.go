package libs

import (
	"fmt"
	"math/rand"
	"time"
)

func RandToken(len int) string {
	b := make([]byte, len)
	rand.Seed(time.Now().UTC().UnixNano())
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
