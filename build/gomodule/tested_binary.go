package gomodule

import (
	"fmt"
	"path"

	"github.com/google/blueprint"
	"github.com/roman-mazur/bood"
)

var (
	// Package context used to define Ninja build rules.
	pctx = blueprint.NewPackageContext("github.com/Encelad/ArchLab2/build/gomodule")

	// Ninja rule to execute go build.
	goBuild = pctx.StaticRule("binaryBuild", blueprint.RuleParams{
		Command:     "cd $workDir && go build -o $outputPath $pkg",
		Description: "build go command $pkg",
	}, "workDir", "outputPath", "pkg")

	// Ninja rule to execute go test.
	goTest = pctx.StaticRule("gotest", blueprint.RuleParams{
		Command:     "cd $workDir && go test -v -bench=. -benchtime=100x $pkgTest > $outputFile",
		Description: "go test $pkgTest",
	}, "workDir", "pkgTest", "outputFile")

	// Ninja rule to execute go mod vendor.
	goVendor = pctx.StaticRule("vendor", blueprint.RuleParams{
		Command:     "cd $workDir && go mod vendor",
		Description: "vendor dependencies of $name",
	}, "workDir", "name")

	config *(bood.Config)
)

// goBinaryModuleType implements the simplest Go binary build without running tests for the target Go package
type goBinaryModuleType struct {
	blueprint.SimpleName

	properties struct {
		// Go package name to build as a command with "go build".
		Pkg         string
		TestPkg     string
		OutTestFile string
		// List of source files.
		Srcs []string
		// Exclude patterns.
		SrcsExclude []string
		Optional    bool
		// If to call vendor command.
		VendorFirst bool
	}
}

func (bm *goBinaryModuleType) GenerateBuildActions(ctx blueprint.ModuleContext) {
	name := ctx.ModuleName()
	config = bood.ExtractConfig(ctx)
	config.Debug.Printf("Adding build actions for go binary module '%s'", name)

	outputPath := path.Join(config.BaseOutputDir, "bin", name)
	outputFile := path.Join(config.BaseOutputDir, "bin", bm.properties.OutTestFile)

	var inputs []string
	inputErors := false
	for _, src := range bm.properties.Srcs {
		if matches, err := ctx.GlobWithDeps(src, bm.properties.SrcsExclude); err == nil {
			inputs = append(inputs, matches...)
		} else {
			ctx.PropertyErrorf("srcs", "Cannot resolve files that match pattern %s", src)
			inputErors = true
		}
	}
	if inputErors {
		return
	}

	var inputsTest []string
	for _, src := range bm.properties.Srcs {
		if matches, err := ctx.GlobWithDeps(src, nil); err == nil {
			inputsTest = append(inputsTest, matches...)
		}
	}

	if bm.properties.VendorFirst {
		vendorDirPath := path.Join(ctx.ModuleDir(), "vendor")
		ctx.Build(pctx, blueprint.BuildParams{
			Description: fmt.Sprintf("Vendor dependencies of %s", name),
			Rule:        goVendor,
			Outputs:     []string{vendorDirPath},
			Implicits:   []string{path.Join(ctx.ModuleDir(), "go.mod")},
			Optional:    true,
			Args: map[string]string{
				"workDir": ctx.ModuleDir(),
				"name":    name,
			},
		})
		inputs = append(inputs, vendorDirPath)
	}

	if len(inputs) != 0 {
		ctx.Build(pctx, blueprint.BuildParams{
			Description: fmt.Sprintf("Build %s as Go binary", name),
			Rule:        goBuild,
			Outputs:     []string{outputPath},
			Implicits:   inputs,
			Args: map[string]string{
				"outputPath": outputPath,
				"workDir":    ctx.ModuleDir(),
				"pkg":        bm.properties.Pkg,
			},
		})
	}

	ctx.Build(pctx, blueprint.BuildParams{
		Description: fmt.Sprintf("Checking my module"),
		Rule:        goTest,
		Outputs:     []string{outputFile},
		Implicits:   inputsTest,
		Optional:    bm.properties.Optional,
		Args: map[string]string{
			"outputFile": outputFile,
			"workDir":    ctx.ModuleDir(),
			"pkgTest":    bm.properties.TestPkg,
		},
	})
}

func (bm *goBinaryModuleType) Outputs() []string {
	return []string{path.Join(config.BaseOutputDir, "bin", bm.Name())}
}

func SimpleBinFactory() (blueprint.Module, []interface{}) {
	mType := &goBinaryModuleType{}
	return mType, []interface{}{&mType.SimpleName.Properties, &mType.properties}
}
