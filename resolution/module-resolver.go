package resolution

import (
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/vknabel/lithia/ast"
)

type ModuleResolver struct {
	// each root contains multiple packages
	externalImportRoots []string
	// defaults to Potfile
	defaultPackageName  string
	manifestName        string
	manifestSearchPaths []string
	defaultSrcDir       string
}

func DefaultModuleResolver() ModuleResolver {
	return ModuleResolver{
		externalImportRoots: defaultImportRoots(),
		defaultPackageName:  "package",
		manifestName:        "Potfile",
		manifestSearchPaths: []string{".", "..", "../..", "../../..", "../../../.."},
		defaultSrcDir:       "src",
	}
}

type ResolvedPackage struct {
	Name     string
	Path     string
	manifest *packageManifest
}

type ResolvedModule struct {
	packageRef *ResolvedPackage
	// all modules of this package are relative to this path
	// might contain the package manifest file
	relativeName ast.ModuleName
	Path         string
	Files        []string
}

type packageManifest struct {
	// a Potfile-file‚
	// the package module path will be derived from this location
	path string
}

func defaultImportRoots() []string {
	roots := []string{}
	if path, ok := os.LookupEnv("LITHIA_LOCALS"); ok {
		roots = append(roots, path)
	}
	if path, ok := os.LookupEnv("LITHIA_PACKAGES"); ok {
		roots = append(roots, path)
	}
	if path, ok := os.LookupEnv("LITHIA_STDLIB"); ok {
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
	for _, candidates := range mr.manifestSearchPaths {
		manifestPath := filepath.Join(path.Dir(referenceFile), candidates, mr.manifestName)
		if _, err := os.Stat(manifestPath); err == nil {
			modulePath := path.Dir(manifestPath)
			return ResolvedPackage{
				Name: mr.defaultPackageName,
				Path: modulePath,
				manifest: &packageManifest{
					path: manifestPath,
				},
			}
		}
	}
	dir, err := os.Getwd()
	if err != nil {
		dir = path.Dir(referenceFile)
	}
	return ResolvedPackage{Name: mr.defaultPackageName, Path: dir}
}

func (mr *ModuleResolver) AddRootImportPath(path string) {
	mr.externalImportRoots = append([]string{path}, mr.externalImportRoots...)
}

func (mr *ModuleResolver) CreateSingleFileModule(pkg ResolvedPackage, file string) ResolvedModule {
	return ResolvedModule{
		packageRef:   &pkg,
		relativeName: ast.ModuleName(strings.ReplaceAll(strings.TrimSuffix(filepath.Base(file), ".lithia"), ".", "_")),
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
		if info, err := os.Stat(packagePath); err == nil && info.IsDir() {
			var match ResolvedPackage
			manifestPath := path.Join(packagePath, mr.manifestName)
			if _, err := os.Stat(manifestPath); err == nil && !info.IsDir() {
				match = ResolvedPackage{
					Name: packageName,
					Path: packagePath,
					manifest: &packageManifest{
						path: manifestPath,
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
		files, err := filepath.Glob(path.Join(pkg.Path, mr.defaultSrcDir, "*.lithia"))
		if len(files) > 0 {
			return ResolvedModule{
				packageRef: &pkg,
				Path:       pkg.Path,
				Files:      files,
			}, err
		}
		files, err = filepath.Glob(path.Join(pkg.Path, "*.lithia"))
		return ResolvedModule{
			packageRef: &pkg,
			Path:       pkg.Path,
			Files:      files,
		}, err
	}
	pathElems := append([]string{pkg.Path}, moduleParts...)
	modulePath := path.Join(pathElems...)
	files, err := filepath.Glob(path.Join(modulePath, "*.lithia"))
	return ResolvedModule{
		packageRef:   &pkg,
		relativeName: ast.ModuleName(strings.Join(moduleParts, ".")),
		Path:         pkg.Path,
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
