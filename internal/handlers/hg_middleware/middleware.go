package hg_middleware

import (
	"github.com/turnerbenjamin/heterogen-go/internal/httpErrors"
)

type Middleware func(httpErrors.ReqHandler) httpErrors.ReqHandler
