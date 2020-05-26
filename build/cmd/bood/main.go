package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/Encelad/ArchLab2/build/gomodule"
	jsmodule "github.com/Encelad/ArchLab2/build/jsmodule" //esli zabil
	"github.com/google/blueprint"
	boodmain "github.com/roman-mazur/bood"
)

var (
	dryRun  = flag.Bool("dry-run", false, "Generate ninja build file but don't start the build")
	verbose = flag.Bool("v", false, "Display debugging logs")
)

func NewContext() *blueprint.Context {
	ctx := boodmain.PrepareContext()
	ctx.RegisterModuleType("go_binary", gomodule.SimpleBinFactory)
	ctx.RegisterModuleType("js_bundle", jsmodule.JsMinimizedScriptFactory)
	return ctx
}

func main() {
	flag.Parse()

	config := boodmain.NewConfig()
	if !*verbose {
		config.Debug = log.New(ioutil.Discard, "", 0)
	}
	ctx := NewContext()

	ninjaBuildPath := boodmain.GenerateBuildFile(config, ctx)

	if !*dryRun {

		config.Info.Println("Starting the build now")
		cmd := exec.Command("ninja", append([]string{"-f", ninjaBuildPath}, flag.Args()...)...)
		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			config.Info.Fatal("Error invoking ninja build. See logs above.")
		}
	}
}
