package main

import (
	"embed"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sort"
	"strings"
	"text/template"

	sprig "github.com/go-task/slim-sprig"
	"github.com/jessevdk/go-flags"
	"github.com/pb33f/libopenapi"
	validator "github.com/pb33f/libopenapi-validator"
	"github.com/pb33f/libopenapi/datamodel"
	"github.com/pb33f/libopenapi/datamodel/high/base"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/pb33f/libopenapi/renderer"
	"golang.org/x/exp/slog"
)

//go:embed templates/*
var templates embed.FS

type endpoint struct {
	Name     string
	Resource *v3.PathItem
	Weight   int
}

type schemaData struct {
	Name   string
	Schema *base.Schema
	Weight int
}

type config struct {
	Dir        string `long:"dir" short:"d" description:"The directory in which to place the cfg.Dir files" required:"true"`
	Spec       string `long:"spec" short:"s" description:"The OpenAPI spec file to laod" required:"true"`
	SinglePage bool   `short:"1" description:"Create a single output page of the entire API"`
	Github     bool   `short:"g" description:"Create Github compatible markdown"`
}

func main() {
	cfg := &config{}
	_, err := flags.Parse(cfg)
	if err != nil {
		log.Fatalf("Cannot parse command arguments: %v", err)
	}

	fInfo, err := os.Stat(cfg.Dir)
	if err != nil {
		log.Fatalf("Error getting output dir info: %v", err)
	}
	if !fInfo.IsDir() {
		log.Fatalf("%s is not a directory", cfg.Dir)
	}

	prefix := ""
	if !cfg.Github {
		prefix = "hugo"
	}

	mg := renderer.NewMockGenerator(renderer.JSON)
	mg.SetPretty()

	indextmpl, err := template.New("base").Funcs(sprig.FuncMap()).Funcs(
		template.FuncMap{
			"renderExample": func(obj any) string {
				defer func() {
					if r := recover(); r != nil {
						slog.Error(fmt.Sprintf("Could not render example: %v", r))
					}
				}()
				result, err := mg.GenerateMock(obj, "")
				if err != nil {
					return fmt.Sprintf("Could not render example: %v", err)
				}
				return string(result)
			},
		},
	).ParseFS(templates, "templates/*.tmpl")
	if err != nil {
		panic(err)
	}

	fd, err := ioutil.ReadFile(cfg.Spec)
	if err != nil {
		panic(err)
	}

	apicfg := datamodel.DocumentConfiguration{
		AllowFileReferences:   true,
		AllowRemoteReferences: true,
	}

	document, err := libopenapi.NewDocumentWithConfiguration(fd, &apicfg)
	if err != nil {
		panic(err)
	}

	highLevelValidator, validatorErrs := validator.NewValidator(document)
	if validatorErrs != nil {
		slog.Error("Please fix validation errors:")
		for _, ve := range validatorErrs {
			slog.Error(ve.Error())
		}
		log.Fatal("")
	}
	if ok, vErrs := highLevelValidator.ValidateDocument(); !ok {
		slog.Error("OpenAPI document is not valid:")
		for _, ve := range vErrs {
			slog.Error(fmt.Sprintf("%s: %s", ve.Message, ve.HowToFix))
		}
		log.Fatal()
	}

	oapi, errors := document.BuildV3Model()

	// if anything went wrong when building the v3 model, a slice of errors will be returned
	if len(errors) > 0 {
		for i := range errors {
			fmt.Printf("error: %e\n", errors[i])
		}
		panic(fmt.Sprintf("cannot create v3 model from document: %d errors reported",
			len(errors)))
	}
	indexName := "_index.md"
	if cfg.SinglePage || cfg.Github {
		indexName = "README.md"
	}
	out, err := os.Create(path.Join(cfg.Dir, indexName))
	err = indextmpl.ExecuteTemplate(out, fmt.Sprintf("%s%s", prefix, "index"), oapi.Model)
	if err != nil {
		panic(err)
	}
	if !cfg.SinglePage {
		out.Close()
	}
	if oapi.Model.Paths != nil {
		paths := make([]string, 0, len(oapi.Model.Paths.PathItems))
		for k := range oapi.Model.Paths.PathItems {
			paths = append(paths, k)
		}
		sort.Strings(paths)
		for x, name := range paths {
			resource := oapi.Model.Paths.PathItems[name]
			if !cfg.SinglePage {
				out, err = os.Create(path.Join(cfg.Dir, fileName(name)))
				if err != nil {
					panic(err)
				}
			}
			endpoint := &endpoint{Name: name, Resource: resource, Weight: (x + 1) * 10}
			err = indextmpl.ExecuteTemplate(out, fmt.Sprintf("%s%s", prefix, "resource"), endpoint)
			if err != nil {
				slog.Error(fmt.Sprintf("Error rendering %s%s (%s): %v", prefix, "resource", endpoint.Name, err))
			}
			if !cfg.SinglePage {
				out.Close()
			}
		}
	}
	/*  For generating separate schema pages
	if oapi.Model.Components != nil && oapi.Model.Components.Schemas != nil {
		if !cfg.SinglePage {
			out, err = os.Create(path.Join(cfg.Dir, "schemas", indexName))
			if err != nil {
				panic(err)
			}
		}
		err = indextmpl.ExecuteTemplate(out, fmt.Sprintf("%sschema%s", prefix, "index"), oapi.Model)
		if err != nil {
			panic(err)
		}
		if !cfg.SinglePage {
			out.Close()
		}
		schemas := make([]string, 0, len(oapi.Model.Components.Schemas))
		for k := range oapi.Model.Components.Schemas {
			schemas = append(schemas, k)
		}
		sort.Strings(schemas)
		for x, sName := range schemas {
			schema := oapi.Model.Components.Schemas[sName]
			data := &schemaData{Name: sName, Schema: schema.Schema(), Weight: (x + 1) * 10}
			if !cfg.SinglePage {
				out, err = os.Create(path.Join(cfg.Dir, "schemas", fileName(sName)))
				if err != nil {
					panic(err)
				}
			}
			err = indextmpl.ExecuteTemplate(out, fmt.Sprintf("%s%s", prefix, "schema"), data)
			if err != nil {
				slog.Error(fmt.Sprintf("Error rendering %s%s (%s): %v", prefix, "schema", sName, err))
			}
			if !cfg.SinglePage {
				out.Close()
			}
		}
	}
	*/
	out.Close()
}

func fileName(name string) string {
	return strings.Replace(name, "/", "_", -1) + ".md"
}
