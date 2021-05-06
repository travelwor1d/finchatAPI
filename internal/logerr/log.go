package logerr

import (
	"log"

	"cloud.google.com/go/errorreporting"
	"github.com/finchatapp/finchat-api/internal/appconfig"
	"github.com/valyala/fasthttp"
)

type Logerr struct {
	r *errorreporting.Client
}

func New(r *errorreporting.Client) *Logerr {
	return &Logerr{r}
}

func (c *Logerr) LogError(err error, req *fasthttp.Request) {
	if appconfig.Config.ErrorReporting {
		c.r.Report(errorreporting.Entry{
			Error: err,
		})
	}
	log.Print(err)
}
