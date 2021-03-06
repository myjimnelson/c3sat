package luaciv3

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	prompt "github.com/c-bata/go-prompt"
)

func myCompleter(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

// I'm thinking of just using tests while developing the package, so will
//   try to put stuff here instead of a new executable
// Go will cache test results; run `go test` with `-count=1` to skip caching the interactive input/output
func disabledTestWhatevs(t *testing.T) {

	doStrings := []struct {
		desc string
		cmd  string
	}{
		/*
				{"Print Hi", `print("Hi from lua")`},
					{"Pring _G k/v pairs", `for k, v in pairs(_G) do
						print(k, v)
						end`},
				{"Print _Version", `print(_VERSION)`},
				{"test(5, 2)", `print(test(5, 2))`},
			{"Print civ3", `print(civ3)`},
		*/
		{"Print install_path", `print(install_path)`},
		{"Print bic.loadDefault()", `print(bic.load_default())`},
		{"Print sav.load(<path>)", `print(sav.load(install_path .. "/Saves/Auto/Conquests Autosave 4000 BC.SAV"))`},
		/*
			{"Print sav k/v pairs", `for k, v in pairs(sav) do
				print(k, v)
				end`},
			{"Print bic k/v pairs", `for k, v in pairs(bic) do
					print(k, v)
					end`},
		*/
		{"Print civ3 k/v pairs", `for k, v in pairs(civ3) do
			print(k, v)
			end`},
		{"Print save_path", `print(save_path)`},
		{"Print save_name", `print(save_name)`},
		{"Print bic_path", `print(bic_path)`},
		{"Print bic_name", `print(bic_name)`},
		{"Print sav.dump()", `print(sav.dump())`},
		{"Print bic.dump()", `print(bic.dump())`},
		// {"", ``},
		// {"", ``},
	}
	L := NewState()
	defer L.Close()
	for _, o := range doStrings {
		fmt.Printf("=== %s ===\n", o.desc)
		if err := L.DoString(o.cmd); err != nil {
			t.Error("DoString: ", err.Error())
		}

	}
	/*
		l := prompt.Input("> ", myCompleter)
		if err := L.DoString(l); err != nil {
			t.Error("DoString: ", err.Error())
		}
	*/
}

// go-prompt example
/*
func completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "users", Description: "Store the username and age"},
		{Text: "articles", Description: "Store the article text posted by user"},
		{Text: "comments", Description: "Store the text commented to articles"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func TestGoPrompt(t *testing.T) {
	fmt.Println("Please select table.")
	tbl := prompt.Input("> ", completer)
	fmt.Println("You selected " + tbl)
}
*/

func TestLuaScript(t *testing.T) {
	L := NewState()
	// get filename of current file; will use relative path from here for test data input
	_, filename, _, _ := runtime.Caller(0)
	luaFile := filepath.Dir(filename) + "/test_scripts/luaciv3test.lua"
	if err := L.DoFile(luaFile); err != nil {
		t.Fatal("Error running lua script: ", err)
	}
}
