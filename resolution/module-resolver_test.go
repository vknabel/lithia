package resolution_test

import (
	"testing"

	"github.com/vknabel/lithia/resolution"
	"github.com/vknabel/lithia/testing/worldtest"
	"github.com/vknabel/lithia/world"
)

func TestDefaultModuleResolverAgainstNonSrcRootModule(t *testing.T) {
	world.Current = worldtest.NewTestWorld(
		map[string]string{},
		map[string][]byte{
			"foo/proj/Potfile": []byte(`
			import pot.cmds
			cmds.run "test", "cmd/test.lithia"
			`),
			"foo/proj/main.lithia": []byte(`
			package main
			print "hello"
			`),
		},
	)
	resolver := resolution.NewDefaultModuleResolver()
	mod := resolver.ResolvePackageAndModuleForReferenceFile("foo/proj/main.lithia")
	if mod.Package().Name != "proj" {
		t.Errorf("expected proj package, got %s", mod.Package().Name)
	}
	if mod.AbsoluteModuleName() != "proj" {
		t.Errorf("expected proj module, got %s", mod.AbsoluteModuleName())
	}
	if len(mod.Files) != 1 {
		t.Errorf("expected 1 file, got %d. %q", len(mod.Files), mod.Files)
	}
	if len(mod.Files) > 0 && mod.Files[0] != "foo/proj/main.lithia" {
		t.Errorf("expected main.lithia, got %s", mod.Files)
	}
	if mod.Package().Manifest == nil {
		t.Errorf("expected manifest, got nil")
	}
	if mod.Package().Manifest.Path != "foo/proj/Potfile" {
		t.Errorf("expected manifest at foo/proj/Potfile, got %s", mod.Package().Manifest.Path)
	}
}

func TestDefaultModuleResolverAgainstSrcRootModule(t *testing.T) {
	world.Current = worldtest.NewTestWorld(
		map[string]string{},
		map[string][]byte{
			"foo/proj/Potfile": []byte(`
			import pot.cmds
			cmds.run "test", "cmd/test.lithia"
			`),
			"foo/proj/src/main.lithia": []byte(`
			package main
			print "hello"
			`),
		},
	)
	resolver := resolution.NewDefaultModuleResolver()
	mod := resolver.ResolvePackageAndModuleForReferenceFile("foo/proj/src/main.lithia")
	if mod.Package().Name != "proj" {
		t.Errorf("expected proj package, got %s", mod.Package().Name)
	}
	if mod.AbsoluteModuleName() != "proj" {
		t.Errorf("expected proj module, got %s", mod.AbsoluteModuleName())
	}
	if len(mod.Files) != 1 {
		t.Errorf("expected 1 file, got %d. %q", len(mod.Files), mod.Files)
	}
	if len(mod.Files) > 0 && mod.Files[0] != "foo/proj/src/main.lithia" {
		t.Errorf("expected main.lithia, got %s", mod.Files)
	}
	if mod.Package().Manifest == nil {
		t.Errorf("expected manifest, got nil")
	}
	if mod.Package().Manifest.Path != "foo/proj/Potfile" {
		t.Errorf("expected manifest at foo/proj/Potfile, got %s", mod.Package().Manifest.Path)
	}
}

func TestDefaultModuleResolverAgainstStandaloneFile(t *testing.T) {
	world.Current = worldtest.NewTestWorld(
		map[string]string{},
		map[string][]byte{
			"foo/proj/src/main.lithia": []byte(`
			package main
			print "hello"
			`),
		},
	)
	resolver := resolution.NewDefaultModuleResolver()
	mod := resolver.ResolvePackageAndModuleForReferenceFile("foo/proj/src/main.lithia")
	if mod.Package().Name != "root" {
		t.Errorf("expected root package, got %s", mod.Package().Name)
	}
	if mod.AbsoluteModuleName() != "root.main" {
		t.Errorf("expected root.main module, got %s", mod.AbsoluteModuleName())
	}
	if len(mod.Files) != 1 {
		t.Errorf("expected 1 file, got %d. %q", len(mod.Files), mod.Files)
	}
	if len(mod.Files) > 0 && mod.Files[0] != "foo/proj/src/main.lithia" {
		t.Errorf("expected main.lithia, got %s", mod.Files)
	}
	if mod.Package().Manifest != nil {
		t.Errorf("expected no manifest, got %s", mod.Package().Manifest.Path)
	}
}
