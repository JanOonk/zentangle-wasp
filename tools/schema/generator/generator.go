// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package generator

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// TODO nested structs
// TODO handle case where owner is type AgentID[]

type Generator interface {
	init(s *Schema)
	funcName(f *Func) string
	generateLanguageSpecificFiles() error
	setFieldKeys(pad bool)
	setFuncKeys()
}

type GenBase struct {
	currentField    *Field
	currentFunc     *Func
	currentStruct   *Struct
	emitters        map[string]func(g *GenBase)
	extension       string
	file            *os.File
	folder          string
	funcRegexp      *regexp.Regexp
	gen             Generator
	keys            map[string]string
	language        string
	maxCamelFuncLen int
	maxSnakeFuncLen int
	maxCamelFldLen  int
	maxSnakeFldLen  int
	newTypes        map[string]bool
	rootFolder      string
	s               *Schema
	templates       map[string]string
}

const spaces = "                                             "

func (g *GenBase) init(s *Schema, templates []map[string]string) {
	g.s = s
	g.emitters = map[string]func(g *GenBase){}
	g.newTypes = map[string]bool{}
	g.keys = map[string]string{}
	g.setCommonKeys()
	g.templates = map[string]string{}
	g.addTemplates(commonTemplates)
	for _, template := range templates {
		g.addTemplates(template)
	}
}

func (g *GenBase) addTemplates(t map[string]string) {
	for k, v := range t {
		g.templates[k] = v
	}
}

func (g *GenBase) close() {
	_ = g.file.Close()
}

func (g *GenBase) createFile(path string, overwrite bool, generator func()) (err error) {
	if !overwrite && g.exists(path) == nil {
		return nil
	}
	g.file, err = os.Create(path)
	if err != nil {
		return err
	}
	defer g.close()
	generator()
	return nil
}

// TODO take copyright from schema?
func (g *GenBase) createSourceFile(name string, condition bool) error {
	if !condition {
		return nil
	}
	return g.createFile(g.folder+name+g.extension, true, func() {
		g.emit("copyright")
		g.emit("warning")
		g.emit(name + g.extension)
	})
}

func (g *GenBase) error(what string) {
	g.println("???:" + what)
}

func (g *GenBase) exists(path string) (err error) {
	_, err = os.Stat(path)
	return err
}

func (g *GenBase) funcName(f *Func) string {
	return f.Kind + capitalize(f.Name)
}

func (g *GenBase) Generate(s *Schema) error {
	g.gen.init(s)

	g.folder = g.rootFolder + "/"
	if g.rootFolder != "src" {
		module := strings.ReplaceAll(moduleCwd, "\\", "/")
		module = module[strings.LastIndex(module, "/")+1:]
		g.folder += module + "/"
	}
	if g.s.CoreContracts {
		g.folder += g.s.Name + "/"
	}

	err := os.MkdirAll(g.folder, 0o755)
	if err != nil {
		return err
	}
	info, err := os.Stat(g.folder + "consts" + g.extension)
	if err == nil && info.ModTime().After(s.SchemaTime) {
		fmt.Printf("skipping %s code generation\n", g.language)
		return nil
	}

	fmt.Printf("generating %s code\n", g.language)
	err = g.generateCode()
	if err != nil {
		return err
	}
	if !g.s.CoreContracts {
		err = g.generateTests()
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *GenBase) generateCode() error {
	err := g.createSourceFile("consts", true)
	if err != nil {
		return err
	}
	err = g.createSourceFile("keys", !g.s.CoreContracts)
	if err != nil {
		return err
	}
	err = g.createSourceFile("structs", len(g.s.Structs) != 0)
	if err != nil {
		return err
	}
	err = g.createSourceFile("typedefs", len(g.s.Typedefs) != 0)
	if err != nil {
		return err
	}
	err = g.createSourceFile("params", len(g.s.Params) != 0)
	if err != nil {
		return err
	}
	err = g.createSourceFile("results", len(g.s.Results) != 0)
	if err != nil {
		return err
	}
	err = g.createSourceFile("state", !g.s.CoreContracts)
	if err != nil {
		return err
	}
	err = g.createSourceFile("contract", true)
	if err != nil {
		return err
	}
	err = g.createSourceFile("lib", !g.s.CoreContracts)
	if err != nil {
		return err
	}
	if !g.s.CoreContracts {
		err = g.generateFuncs()
		if err != nil {
			return err
		}
	}

	return g.gen.generateLanguageSpecificFiles()
}

func (g *GenBase) generateFuncs() error {
	scFileName := g.folder + g.s.Name + g.extension
	if g.exists(scFileName) != nil {
		// generate initial SC function file
		return g.createFile(scFileName, false, func() {
			g.emit("copyright")
			g.emit("funcs" + g.extension)
		})
	}

	// append missing SC functions to existing code file

	// scan existing file for function names
	existing := make(StringMap)
	lines := make([]string, 0)
	err := g.scanExistingCode(scFileName, &existing, &lines)
	if err != nil {
		return err
	}

	// save old one from overwrite
	scOriginal := g.folder + g.s.Name + ".bak"
	err = os.Rename(scFileName, scOriginal)
	if err != nil {
		return err
	}

	err = g.createFile(scFileName, false, func() {
		// make copy of original file
		for _, line := range lines {
			g.println(line)
		}

		// append any new funcs
		for _, g.currentFunc = range g.s.Funcs {
			if existing[g.gen.funcName(g.currentFunc)] == "" {
				g.setFuncKeys()
				g.emit("funcSignature")
			}
		}
	})
	if err != nil {
		return err
	}
	return os.Remove(scOriginal)
}

func (g *GenBase) generateTests() error {
	err := os.MkdirAll("test", 0o755)
	if err != nil {
		return err
	}

	// do not overwrite existing file
	name := strings.ToLower(g.s.Name)
	filename := "test/" + name + "_test.go"
	return g.createFile(filename, false, func() {
		g.emit("test.go")
	})
}

func (g *GenBase) openFile(path string, processor func() error) (err error) {
	g.file, err = os.Open(path)
	if err != nil {
		return err
	}
	defer g.close()
	return processor()
}

func (g *GenBase) println(a ...interface{}) {
	_, _ = fmt.Fprintln(g.file, a...)
}

func (g *GenBase) scanExistingCode(path string, existing *StringMap, lines *[]string) error {
	return g.openFile(path, func() error {
		scanner := bufio.NewScanner(g.file)
		for scanner.Scan() {
			line := scanner.Text()
			matches := g.funcRegexp.FindStringSubmatch(line)
			if matches != nil {
				(*existing)[matches[1]] = line
			}
			*lines = append(*lines, line)
		}
		return scanner.Err()
	})
}
