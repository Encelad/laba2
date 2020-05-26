package jsmodule

import (
	"bytes"
	"github.com/google/blueprint"
	"github.com/roman-mazur/bood"
	"strings"
	"testing"
)

var (
	fileSystems = []map[string][]byte{
		{
			"Blueprints": []byte(`
			js_bundle {
			  name: "package-out",
			  srcs: ["main.js", "second.js"],
              obfuscate: true,
			}
		`),
			"main.js":   nil,
			"second.js": nil,
		},
		{
			"Blueprints": []byte(`
			js_bundle {
			  name: "new-name",
			  srcs: ["./new.js", "examples/second.js"],
              obfuscate: false,
			}
		`),
			"new.js":             nil,
			"examples/second.js": nil,
		},
	}

	expectedOutput = [][]string{
		{
			"build out/js/package-out.js",
			"g.jsmodule.obfuscMinim",
			"main.js second.js",
		},
		{
			"build out/js/new-name.js:",
			" g.jsmodule.minim",
			"new.js examples/second.js",
		},
	}
)

func TestJsMinimizedScriptFactory(t *testing.T) {
	for index, fileSystem := range fileSystems {
		t.Run(string(index), func(t *testing.T) {
			ctx := blueprint.NewContext()

			ctx.MockFileSystem(fileSystem)

			ctx.RegisterModuleType("js_bundle", JsMinimizedScriptFactory)

			cfg := bood.NewConfig()

			_, errs := ctx.ParseBlueprintsFiles(".", cfg)
			if len(errs) != 0 {
				t.Fatalf("Syntax errors in the test blueprint file: %s", errs)
			}

			_, errs = ctx.PrepareBuildActions(cfg)
			if len(errs) != 0 {
				t.Errorf("Unexpected errors while preparing build actions: %s", errs)
			}
			buffer := new(bytes.Buffer)
			if err := ctx.WriteBuildFile(buffer); err != nil {
				t.Errorf("Error writing ninja file: %s", err)
			} else {
				text := buffer.String()
				//t.Logf("Gennerated ninja build file:\n%s", text)
				for _, expectedStr := range expectedOutput[index] {
					//build rule
					if strings.Contains(text, expectedStr) != true {
						t.Errorf("Generated ninja file does not have expected string `%s`", expectedStr)
					}
				}
			}

		})
	}
}
