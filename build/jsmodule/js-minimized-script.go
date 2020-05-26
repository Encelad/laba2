package jsmodule

import (
	"fmt"
	"path"
	"strings"

	"github.com/google/blueprint"
	"github.com/roman-mazur/bood"
)

var (
	// Package context used to define Ninja build rules.
	pctx = blueprint.NewPackageContext("github.com/Encelad/ArchLab2/build/jsmodule")

	obfuscMinimizeRule = pctx.StaticRule("obfuscMinim", blueprint.RuleParams{
		Command:     "npx webpack --mode=production ${input} -o ${output} > ${log} ",
		Description: "Obfuscate and minimize js script ",
	}, "workDir", "input", "output", "log")

	minimizeRule = pctx.StaticRule("minim", blueprint.RuleParams{
		Command:     "cd ${workDir} && npx webpack --mode=none ${input} -o ${output} > ${log}",
		Description: "Minimize js script ",
	}, "workDir", "input", "output", "log")
)

type jsScriptModule struct {
	blueprint.SimpleName

	properties struct {
		Name string
		// List of source files.
		Srcs []string
		// Exclude patterns.
		SrcsExclude []string
		Obfuscate   bool
	}
}

func (jm *jsScriptModule) GenerateBuildActions(ctx blueprint.ModuleContext) {

	name := ctx.ModuleName()
	config := bood.ExtractConfig(ctx)
	config.Debug.Printf("Adding minimizing actions for js script '%s'", name)

	outputFile := path.Join(config.BaseOutputDir, "js", name+".js")
	logFile := path.Join(config.BaseOutputDir, "js", "log.txt")

	var inputs []string

	for _, src := range jm.properties.Srcs {
		if matches, err := ctx.GlobWithDeps(src, jm.properties.SrcsExclude); err == nil {
			inputs = append(inputs, matches...)
		} else {
			ctx.PropertyErrorf("srcs", "Cannot resolve files that match pattern %s", src)
			return
		}
	}

	input := strings.Join(inputs[:], " ")

	if jm.properties.Obfuscate && jm.properties.Obfuscate == true {
		ctx.Build(pctx, blueprint.BuildParams{
			Description: fmt.Sprintf("Minimize and obfuscate script %s", name),
			Rule:        obfuscMinimizeRule,
			Outputs:     []string{outputFile},
			Implicits:   inputs,
			Args: map[string]string{
				"workDir": ctx.ModuleDir(),
				"input":   input,
				"output":  outputFile,
				"log":     logFile,
			},
		})

	} else {

		ctx.Build(pctx, blueprint.BuildParams{
			Description: fmt.Sprintf("Minimize script %s", name),
			Rule:        minimizeRule,
			Outputs:     []string{outputFile},
			Implicits:   inputs,
			Args: map[string]string{
				"workDir": ctx.ModuleDir(),
				"input":   input,
				"output":  outputFile,
				"log":     logFile,
			},
		})
	}
}

func JsMinimizedScriptFactory() (blueprint.Module, []interface{}) {
	mType := &jsScriptModule{}
	return mType, []interface{}{&mType.SimpleName.Properties, &mType.properties}
}
