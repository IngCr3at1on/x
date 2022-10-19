module github.com/ingcr3at1on/x/with-env

go 1.17

replace github.com/ingcr3at1on/x/lazyfstools => ../lazyfstools

require (
	github.com/ingcr3at1on/x/lazyfstools v0.0.0
	github.com/ingcr3at1on/x/sigctx v0.0.0-20211221234048-f0c5e7996872
	github.com/joho/godotenv v1.4.0
	github.com/spf13/afero v1.8.1
	github.com/spf13/cobra v1.3.0
	github.com/stretchr/testify v1.7.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/text v0.3.7 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)
