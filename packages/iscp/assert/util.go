package assert

import (
	"fmt"
	"strings"

	"github.com/iotaledger/wasp/packages/iscp"
)

type Assert struct {
	log    iscp.LogInterface
	prefix string
}

func NewAssert(log iscp.LogInterface, name ...string) Assert {
	p := "assertion failed: "
	if len(name) > 0 {
		p = fmt.Sprintf("assertion failed (%s): ", name[0])
	}
	return Assert{
		log:    log,
		prefix: p,
	}
}

//nolint:goprintffuncname
func (a Assert) Require(cond bool, format string, args ...interface{}) {
	if cond {
		return
	}
	if a.log == nil {
		panic(fmt.Sprintf(a.prefix+format, args...))
	}
	a.log.Panicf(a.prefix+format, args...)
}

func (a Assert) RequireNoError(err error, str ...string) {
	a.Require(err == nil, fmt.Sprintf(a.prefix+"%s %v", strings.Join(str, " "), err))
}

func (a Assert) RequireChainOwner(ctx iscp.Sandbox, name ...string) {
	if !ctx.ChainOwnerID().Equals(ctx.Caller()) {
		if len(name) > 0 {
			a.log.Panicf(a.prefix+"unauthorized access: %s", name[0])
		} else {
			a.log.Panicf(a.prefix + "unauthorized access")
		}
	}
}
