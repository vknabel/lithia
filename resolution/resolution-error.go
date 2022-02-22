package resolution

import (
	"fmt"
	"strings"
)

type ModuleNotFoundError struct {
	FromPackage ResolvedPackage
	ModuleParts []string
}

func (e ModuleNotFoundError) Error() string {
	return fmt.Sprintf("module %s not found from package %s", strings.Join(e.ModuleParts, "."), e.FromPackage.Path)
}
