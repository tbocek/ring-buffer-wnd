module tomtp

go 1.22.0

toolchain go1.22.2

replace github.com/MatusOllah/slogcolor v1.1.0 => ../slogcolor

require (
	filippo.io/edwards25519 v1.1.0
	github.com/MatusOllah/slogcolor v1.1.0
	github.com/fatih/color v1.16.0
	github.com/stretchr/testify v1.9.0
	golang.org/x/crypto v0.23.0
	golang.org/x/sys v0.20.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rogpeppe/go-internal v1.12.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
