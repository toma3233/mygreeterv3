package logattrs

import (
	log "log/slog"

	<<- if contains .envInformation.goModuleNamePrefix "go.goms.io">>
	// TODO: Internal AKS pkg, use generic pkg in future
	"go.goms.io/aks/rp/toolkit/aksbinversion"
	<<end>>
)

var attrs []log.Attr

func init() {
	<<- if contains .envInformation.goModuleNamePrefix "go.goms.io">>
	attrs = append(attrs,
		log.String("version", aksbinversion.GetVersion()),
		log.String("branch", aksbinversion.GetGitBranch()),
	)
	<<end>>
}

func GetAttrs() []log.Attr {
	return attrs
}
