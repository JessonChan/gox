/*
 Copyright 2021 The GoPlus Authors (goplus.org)
 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at
     http://www.apache.org/licenses/LICENSE-2.0
 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package packages

import (
	"errors"
	"go/types"
	"syscall"
	"testing"
)

// ----------------------------------------------------------------------------

func TestLoadDep(t *testing.T) {
	pkgs, err := loadDeps("./.gop", "fmt")
	if err != nil {
		t.Fatal("LoadDeps failed:", pkgs, err)
	}
	if _, ok := pkgs["runtime"]; !ok {
		t.Fatal("LoadDeps failed:", pkgs)
	}

	err = loadDepPkgsFrom(nil, " ")
	if err != nil {
		t.Fatal("LoadDeps: error?", err)
	}
}

func TestLoadDepErr(t *testing.T) {
	_, err := loadDeps("/.gop", "fmt")
	if err == nil {
		t.Fatal("LoadDeps: no error")
	}
}

// ----------------------------------------------------------------------------

func TestLoadPkgs(t *testing.T) {
	pkgs, err := LoadPkgs("", "fmt", "strings")
	if err != nil {
		t.Fatal("Load failed:", err)
	}
	if len(pkgs) != 2 {
		t.Log(pkgs)
	}
}

func TestLoadPkgsErr(t *testing.T) {
	{
		err := &ExecCmdError{Stderr: []byte("Hi")}
		if err.Error() != "Hi" {
			t.Fatal("ExecCmdError failed:", err)
		}

		err = &ExecCmdError{Err: errors.New("Hi")}
		if err.Error() != "Hi" {
			t.Fatal("ExecCmdError failed:", err)
		}
	}
	pkgs, err := LoadPkgs("", "?")
	if err == nil || err.Error() != `malformed import path "?": invalid char '?'
exit status 1` {
		t.Fatal("loadPkgs:", pkgs, err)
	}
}

func TestLoadPkgsFromErr(t *testing.T) {
	_, err := loadPkgsFrom(nil, []byte("{"))
	if err == nil {
		t.Fatal("loadPkgs no error?")
	}
	_, err = loadPkgsFrom(nil, []byte("{\n"))
	if err == nil {
		t.Fatal("loadPkgs no error?")
	}
	_, err = loadPkgsFrom(nil, []byte("{\n1\n}\n"))
	if err == nil {
		t.Fatal("loadPkgs no error?")
	}
}

// ----------------------------------------------------------------------------

func TestLoadErr(t *testing.T) {
	pkgs, err := Load(nil, "?")
	if err == nil || err.Error() != `exit status 1` {
		t.Fatal("Load:", pkgs, err)
	}

	/*	_, err = loadPkgExport("/not-found", nil, make(map[string]*types.Package), "fmt")
		if err == nil {
			t.Fatal("loadPkgExport no error?")
		}
		_, err = loadPkgExport("load.go", nil, make(map[string]*types.Package), "fmt")
		if err == nil {
			t.Fatal("loadPkgExport no error?")
		}
	*/
}

func TestLoadNoConf(t *testing.T) {
	pkgs, err := Load(nil, "fmt", "strings")
	if err != nil {
		t.Fatal("Load failed:", err)
	}
	if len(pkgs) != 2 {
		t.Log(pkgs)
	}
}

func TestLoadConf(t *testing.T) {
	conf := &Config{
		Loaded: make(map[string]*types.Package),
	}
	pkgs1, err := Load(conf, "fmt", "strings")
	if err != nil {
		t.Fatal("Load failed:", err)
	}
	if len(pkgs1) != 2 {
		t.Log(pkgs1)
	}

	pkgs2, err := Load(conf, "fmt", "strconv")
	if err != nil {
		t.Fatal("Load failed:", err)
	}
	if len(pkgs2) != 2 {
		t.Log(pkgs2)
	}

	if pkgs1[0] != pkgs2[0] {
		t.Fatal("Load failed: unmatched `fmt` pkg")
	}
}

func TestImporterNormal(t *testing.T) {
	conf := &Config{
		Loaded:  make(map[string]*types.Package),
		ModPath: "github.com/goplus/gox/packages",
	}
	p, _, err := NewImporter(conf, ".")
	if err != nil {
		t.Fatal("NewImporter failed:", err)
	}
	pkg, err := p.Import("fmt")
	if err != nil || pkg.Path() != "fmt" {
		t.Fatal("Import failed:", pkg, err)
	}
	if _, err = p.Import("not-found"); err != syscall.ENOENT {
		t.Fatal("Import not-found:", err)
	}
}

func TestImporterRecursive(t *testing.T) {
	conf := &Config{
		Loaded:  make(map[string]*types.Package),
		ModRoot: "..",
		ModPath: "github.com/goplus/gox",
	}
	p, pkgPaths, err := NewImporter(conf, "../internal/go/...")
	if err != nil {
		t.Fatal("NewImporter failed:", err)
	}
	if len(pkgPaths) != 2 {
		t.Fatal("NewImporter pkgPaths:", pkgPaths)
	}
	pkg, err := p.Import(pkgPaths[0])
	if err != nil || pkg.Path() != pkgPaths[0] {
		t.Fatal("Import failed:", pkg, pkgPaths, err)
	}
}

func TestImporterRecursiveErr(t *testing.T) {
	conf := &Config{
		Loaded:  make(map[string]*types.Package),
		ModPath: "github.com/goplus/gox/packages",
	}
	p, pkgPaths, err := NewImporter(conf, "/...")
	if err == nil || err.Error() != "directory `/` outside available modules" {
		t.Fatal("NewImporter failed:", p, pkgPaths, err)
	}
}

// ----------------------------------------------------------------------------
