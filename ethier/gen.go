package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/ethereum/go-ethereum/common/compiler"
	"github.com/spf13/cobra"
	"golang.org/x/tools/go/ast/astutil"

	_ "embed"
)

const srcMapFlag = "experimental_src_map"

func init() {
	cmd := &cobra.Command{
		Use:   "gen",
		Short: "Compiles Solidity contracts to generate Go ABI bindings with go:generate",
		RunE:  gen,
		Args: func(_ *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("no source files provided")
			}
			for _, a := range args {
				if !strings.HasSuffix(a, ".sol") {
					return fmt.Errorf("non-Solidity file %q", a)
				}
			}
			return nil
		},
	}

	cmd.Flags().Bool(srcMapFlag, false, "Generate source maps to determine Solidity code location from EVM traces")

	rootCmd.AddCommand(cmd)
}

// gen runs `solc | abigen` on the Solidity source files passed as the args.
// TODO: support wildcard / glob matching of files.
func gen(cmd *cobra.Command, args []string) (retErr error) {
	pwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("os.Getwd(): %v", err)
	}
	// The Go package for abigen.
	pkg := filepath.Base(pwd)
	log.Printf("Generating package %q: %s", pkg, args)

	defer func() {
		if retErr != nil {
			retErr = fmt.Errorf("generating %q: %w", pkg, retErr)
		}
	}()

	// solc requires a base-path within which absolute includes are found. We
	// define this as the base path of the Go module.
	basePath := pwd
	for ; ; basePath = filepath.Join(basePath, "..") {
		if _, err := os.Stat(filepath.Join(basePath, "go.mod")); !errors.Is(err, os.ErrNotExist) {
			break
		}
	}

	args = append(
		args,
		"--base-path", basePath,
		"--include-path", filepath.Join(basePath, "node_modules"),
		"--combined-json", "abi,bin,bin-runtime,hashes,metadata,srcmap-runtime",
	)
	solc := exec.Command("solc", args...)
	solc.Stderr = os.Stderr

	// TODO: use bind.Bind() directly, instead of piping to abigen, which
	// requires that it's installed and within PATH. Blocked by
	// https://github.com/ethereum/go-ethereum/issues/23939 for which we've
	// submitted a fix.
	abigen := exec.Command(
		"abigen",
		"--combined-json", "/dev/stdin",
		"--pkg", pkg,
	)
	abigen.Stderr = os.Stderr

	r, w := io.Pipe()
	solc.Stdout = w
	combinedJSON := bytes.NewBuffer(nil)
	abigen.Stdin = io.TeeReader(r, combinedJSON)

	generated := bytes.NewBuffer(nil)
	abigen.Stdout = generated

	if err := solc.Start(); err != nil {
		return fmt.Errorf("start `solc`: %v", err)
	}
	if err := abigen.Start(); err != nil {
		return fmt.Errorf("start `abigen`: %v", err)
	}
	if err := solc.Wait(); err != nil {
		w.Close()
		return fmt.Errorf("`solc` returned: %v", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("close write-half of pipe from solc to abigen: %v", err)
	}
	if err := abigen.Wait(); err != nil {
		return fmt.Errorf("`abigen` returned: %v", err)
	}
	if err := r.Close(); err != nil {
		return fmt.Errorf("close read-half of pipe from solc to abigen: %v", err)
	}

	extend, err := cmd.Flags().GetBool(srcMapFlag)
	if err != nil {
		return fmt.Errorf("%T.Flags().GetBool(%q): %v", cmd, srcMapFlag, err)
	}
	if !extend {
		return os.WriteFile("generated.go", generated.Bytes(), 0644)
	}

	out, err := extendGeneratedCode(generated, combinedJSON, basePath)
	if err != nil {
		return err
	}
	return os.WriteFile("generated.go", out, 0644)
}

var (
	//go:embed gen_extra.go.tmpl
	extraCode string

	// extraTemplate is the template for use by extendGeneratedCode().
	extraTemplate = template.Must(
		template.New("extra").
			Funcs(template.FuncMap{
				"quote": func(s interface{}) string {
					return fmt.Sprintf("%q", s)
				},
				"soliditySourceSlice": func(as, bs []string) (string, error) {
					if n, m := len(as), len(bs); n != m {
						return "", fmt.Errorf("can't zip strings of different lengths %d and %d", n, m)
					}
					q := make([]string, len(as))
					for i := range as {
						q[i] = fmt.Sprintf("{File: %q, Code: %q}", as[i], bs[i])
					}
					return fmt.Sprintf("[]*solidity.Source{\n%s,\n}", strings.Join(q, ",\n")), nil
				},
				"contract": func(s string) (string, error) {
					parts := strings.Split(s, ".sol:")
					if len(parts) != 2 {
						return "", fmt.Errorf("invalid contract name %q must have format path/to/file.sol:ContractName", s)
					}
					return parts[1], nil
				},
			}).
			Parse(extraCode),
	)

	// deployedRegex matches abigen-generated return statements for Deploy*()
	// functions to allow their addresses to be stored in a map to be passed to
	// solidity.NewSourceMap().
	deployedRegexp = regexp.MustCompile(`^\s*return address, tx, &(.+?)\{.*Transactor.*\}, nil\s*$`)
	// TODO(aschlosberg) replace regular expressions with a more explicit
	// approach once https://github.com/ethereum/go-ethereum/pull/25055 has been
	// merged.
)

// extendGeneratedCode adds ethier-specific functionality to code generated by
// abigen, allowing for interoperability with the ethier/solidity package for
// source-map interpretation at runtime.
func extendGeneratedCode(generated, combinedJSON *bytes.Buffer, basePath string) ([]byte, error) {
	meta := struct {
		SourceList []string `json:"sourceList"`
		SourceCode []string
		Version    string `json:"version"`

		Contracts    map[string]*compiler.Contract
		CombinedJSON string
	}{CombinedJSON: combinedJSON.String()}

	if err := json.Unmarshal(combinedJSON.Bytes(), &meta); err != nil {
		return nil, fmt.Errorf("json.Unmarshal([solc output], %T): %v", &meta, err)
	}

	// TODO(aschlosberg) move the source code into generated_test.go.
	meta.SourceCode = make([]string, len(meta.SourceList))
	for i, src := range meta.SourceList {
		f := filepath.Join(basePath, src)
		buf, err := os.ReadFile(f)
		if os.IsNotExist(err) {
			// As this file isn't in the basePath, it's from a library, so we
			// don't care about its code being monitored.
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("os.ReadFile(%q): %v", f, err)
		}
		meta.SourceCode[i] = string(buf)
	}

	cs, err := compiler.ParseCombinedJSON(combinedJSON.Bytes(), "", "", meta.Version, "")
	if err != nil {
		return nil, fmt.Errorf("compiler.ParseCombinedJSON(): %v", err)
	}
	meta.Contracts = cs
	for k, c := range meta.Contracts {
		if c.RuntimeCode == "0x" {
			delete(meta.Contracts, k)
		}
	}

	if err := extraTemplate.Execute(generated, meta); err != nil {
		return nil, fmt.Errorf("%T.Execute(): %v", extraTemplate, err)
	}

	// When using vm.Config.Trace, the only contract-identifying information is
	// the address to which the transaction was sent. We must therefore modify
	// every DeployFoo() function to save the address(es) at which the contract
	// is deployed.
	lines := strings.Split(generated.String(), "\n")
	for i, l := range lines {
		matches := deployedRegexp.FindStringSubmatch(l)
		if len(matches) == 0 {
			continue
		}
		lines[i] = fmt.Sprintf(
			`deployedContracts[address] = %q // Added by ethier gen
			%s`,
			matches[1], l,
		)
	}

	// Effectively the same as running goimports on the (ugly) generated code.
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "generated.go", strings.Join(lines, "\n"), parser.ParseComments|parser.AllErrors)
	if err != nil {
		return nil, fmt.Errorf("parser.ParseFile(%T, …): %v", fset, err)
	}
	for _, pkg := range []string{
		"github.com/ethereum/go-ethereum/common/compiler",
		"github.com/divergencetech/ethier/solidity",
	} {
		if !astutil.AddImport(fset, f, pkg) {
			return nil, fmt.Errorf("add import %q to generated Go: %v", pkg, err)
		}
	}

	buf := bytes.NewBuffer(nil)
	if err := format.Node(buf, fset, f); err != nil {
		return nil, fmt.Errorf("format.Node(%T, %T, %T): %v", buf, fset, f, err)
	}
	return buf.Bytes(), nil
}
