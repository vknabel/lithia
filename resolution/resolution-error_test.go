package resolution_test

import (
	"testing"

	"github.com/vknabel/lithia/resolution"
)

func TestModuleFoundError(t *testing.T) {
	e := resolution.ModuleNotFoundError{
		FromPackage: resolution.ResolvedPackage{
			Path: "foo/bar",
		},
		ModuleParts: []string{"foo", "bar"},
	}
	if e.Error() != "module foo.bar not found from package foo/bar" {
		t.Errorf("Expected error to be \"module foo.bar not found from package foo/bar\", got %s", e.Error())
	}
}
