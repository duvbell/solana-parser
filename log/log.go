package log

import (
	"fmt"
	"github.com/hashicorp/go-hclog"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"time"
)

var (
	Logger = hclog.Default()
)

func init() {
	fileName := fmt.Sprintf("logs/%s.log.%%Y%%m%%d", "parser")
	rl, err := rotatelogs.New(fileName, rotatelogs.WithMaxAge(120*24*time.Hour), rotatelogs.WithClock(rotatelogs.UTC))
	if err != nil {
		panic(err)
	}
	Logger = hclog.New(&hclog.LoggerOptions{
		Output: rl,
		Level:  hclog.Info,
		Name:   "parser",
	})
}
