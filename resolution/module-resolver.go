package resolution

import (
	"path"
	"path/filepath"
	"strings"

	"github.com/vknabel/lithia/ast"
	"github.com/vknabel/lithia/world"
)

type ModuleResolver struct {
	// each root contains multiple packages
	externalImportRoots []string
	// defaults to Potfile
	defaultPackageName  string
	manifestName        string
	manifestSearchPaths []string
	defaultSrcDir       string
	lithiaSourceGlob    string
}

func NewDefaultModuleResolver(importRoots ...string) ModuleResolver {
	return ModuleResolver{
		externalImportRoots: defaultImportRoots(importRoots...),
		defaultPackageName:  "root",
		manifestName:        "Potfile",
		manifestSearchPaths: []string{".", "..", "../..", "../../..", "../../../.."},
		defaultSrcDir:       "src",
		lithiaSourceGlob:    "*.lithia",
	}
}

type ResolvedPackage struct {
	Name     string
	Path     string
	Manifest *PackageManifest
}

type ResolvedModule struct {
	packageRef *ResolvedPackage
	// all modules of this package are relative to this path
	// might contain the package manifest file
	relativeName ast.ModuleName
	Path         string
	Files        []string
}

type PackageManifest struct {
	// a Potfile-file‚
	// the package module path will be derived from this location
	Path string
}

func defaultImportRoots(importRoots ...string) []string {
	roots := importRoots
	if path, ok := world.Current.Env.LookupEnv("LITHIA_LOCALS"); ok {
		roots = append(roots, path)
	}
	if path, ok := world.Current.Env.LookupEnv("LITHIA_PACKAGES"); ok {
		roots = append(roots, path)
	}
	if path, ok := world.Current.Env.LookupEnv("LITHIA_STDLIB"); ok {
		roots = append(roots, path)
	} else {
		roots = append(roots, "/usr/local/opt/lithia/stdlib")
	}
	absoluteImportRoots := make([]string, len(roots))
	for i, root := range roots {
		absolute, err := filepath.Abs(root)
		if err == nil {
			absoluteImportRoots[i] = absolute
		} else {
			absoluteImportRoots[i] = root
		}
	}
	return absoluteImportRoots
}

func (mr *ModuleResolver) ResolvePackageForReferenceFile(referenceFile string) ResolvedPackage {
	referenceFile = removingFilePrefix(referenceFile)
	for _, candidate := range mr.manifestSearchPaths {
		manifestPath := filepath.Join(path.Dir(referenceFile), candidate, mr.manifestName)
		if _, err := world.Current.FS.Stat(manifestPath); err == nil {
			packagePath := path.Dir(manifestPath)
			packageName := path.Base(packagePath)
			return ResolvedPackage{
				Name: packageName,
				Path: packagePath,
				Manifest: &PackageManifest{
					Path: manifestPath,
				},
			}
		}
	}
	dir := path.Dir(referenceFile)
	return ResolvedPackage{Name: mr.defaultPackageName, Path: dir}
}

func (mr *ModuleResolver) ResolvePackageAndModuleForReferenceFile(referenceFile string) ResolvedModule {
	referenceFile = removingFilePrefix(referenceFile)
	pkg := mr.ResolvePackageForReferenceFile(referenceFile)
	if pkg.Manifest == nil {
		return mr.CreateSingleFileModule(pkg, referenceFile)
	}
	relativeFile, err := filepath.Rel(pkg.Path, referenceFile)
	if err != nil {
		return mr.CreateSingleFileModule(pkg, referenceFile)
	}
	moduleFilepath := filepath.Dir(relativeFile)
	moduleParts := strings.Split(moduleFilepath, string(filepath.Separator))
	for i := len(moduleParts) - 1; i >= 0; i-- {
		if moduleParts[i] == "." {
			moduleParts = append(moduleParts[:i], moduleParts[i+1:]...)
		}
	}
	if len(moduleParts) >= 1 && moduleParts[0] == mr.defaultSrcDir {
		moduleParts = moduleParts[1:]
	}
	resolvedModule, err := mr.resolveModuleWithinPackage(pkg, moduleParts)
	if err != nil {
		return mr.CreateSingleFileModule(pkg, referenceFile)
	}
	return resolvedModule
}

func (mr *ModuleResolver) AddRootImportPath(path string) {
	mr.externalImportRoots = append([]string{path}, mr.externalImportRoots...)
}

func (mr *ModuleResolver) CreateSingleFileModule(pkg ResolvedPackage, file string) ResolvedModule {
	file = removingFilePrefix(file)
	trimmed := strings.TrimSuffix(filepath.Base(file), ".lithia")
	uniform := strings.ReplaceAll(trimmed, ".", "_")
	return ResolvedModule{
		packageRef:   &pkg,
		relativeName: ast.ModuleName(uniform),
		Path:         file,
		Files:        []string{file},
	}
}

func (mr *ModuleResolver) ResolveModuleFromPackage(pkg ResolvedPackage, moduleName ast.ModuleName) (ResolvedModule, error) {
	moduleParts := strings.Split(string(moduleName), ".")
	if len(moduleParts) == 0 {
		return mr.resolveModuleWithinPackage(pkg, moduleParts)
	}
	packageName := moduleParts[0]
	packageLevelModuleParts := moduleParts[1:]
	if packageName == mr.defaultPackageName {
		return mr.resolveModuleWithinPackage(pkg, packageLevelModuleParts)
	}

	searchPaths := append([]string{pkg.Path}, mr.externalImportRoots...)
	for _, searchPath := range searchPaths {
		packagePath := path.Join(searchPath, packageName)
		if info, err := world.Current.FS.Stat(packagePath); err == nil && info.IsDir() {
			var match ResolvedPackage
			manifestPath := path.Join(packagePath, mr.manifestName)
			if _, err := world.Current.FS.Stat(manifestPath); err == nil && !info.IsDir() {
				match = ResolvedPackage{
					Name: packageName,
					Path: packagePath,
					Manifest: &PackageManifest{
						Path: manifestPath,
					},
				}
			} else {
				match = ResolvedPackage{Name: packageName, Path: packagePath}
			}
			return mr.resolveModuleWithinPackage(match, packageLevelModuleParts)
		}
	}
	return ResolvedModule{}, ModuleNotFoundError{ModuleParts: moduleParts, FromPackage: pkg}
}

func (mr *ModuleResolver) resolveModuleWithinPackage(pkg ResolvedPackage, moduleParts []string) (ResolvedModule, error) {
	if len(moduleParts) == 0 {
		files, err := world.Current.FS.Glob(path.Join(pkg.Path, mr.defaultSrcDir, mr.lithiaSourceGlob))
		if len(files) > 0 {
			return ResolvedModule{
				packageRef: &pkg,
				Path:       path.Join(pkg.Path, mr.defaultSrcDir),
				Files:      files,
			}, err
		}
		files, err = world.Current.FS.Glob(path.Join(pkg.Path, mr.lithiaSourceGlob))
		return ResolvedModule{
			packageRef: &pkg,
			Path:       pkg.Path,
			Files:      files,
		}, err
	}
	pathElems := append([]string{pkg.Path}, moduleParts...)
	modulePath := path.Join(pathElems...)
	files, err := world.Current.FS.Glob(path.Join(modulePath, mr.lithiaSourceGlob))
	return ResolvedModule{
		packageRef:   &pkg,
		relativeName: ast.ModuleName(strings.Join(moduleParts, ".")),
		Path:         modulePath,
		Files:        files,
	}, err
}

func (mod ResolvedModule) Package() ResolvedPackage {
	return *mod.packageRef
}

func (mod ResolvedModule) AbsoluteModuleName() ast.ModuleName {
	if mod.relativeName == "" {
		return ast.ModuleName(mod.Package().Name)
	} else {
		return ast.ModuleName(mod.packageRef.Name) + "." + mod.relativeName
	}
}

func removingFilePrefix(str string) string {
	return strings.TrimPrefix(str, "file://")
}
