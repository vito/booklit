module github.com/vito/booklit

require (
	dagger.io/dagger v0.7.2
	github.com/agext/levenshtein v1.2.3
	github.com/alecthomas/chroma v0.10.0
	github.com/dagger/dagger v0.0.0-00010101000000-000000000000
	github.com/go-bindata/go-bindata v3.1.2+incompatible
	github.com/jessevdk/go-flags v1.5.0
	github.com/mna/pigeon v1.0.1-0.20200224192238-18953b277063
	github.com/onsi/ginkgo/v2 v2.9.1
	github.com/onsi/gomega v1.27.4
	github.com/segmentio/textio v1.2.0
	github.com/sirupsen/logrus v1.9.3
	golang.org/x/sync v0.3.0
)

require (
	github.com/99designs/gqlgen v0.17.31 // indirect
	github.com/Khan/genqlient v0.6.0 // indirect
	github.com/adrg/xdg v0.4.0 // indirect
	github.com/dlclark/regexp2 v1.9.0 // indirect
	github.com/go-logr/logr v1.2.4 // indirect
	github.com/go-task/slim-sprig v0.0.0-20210107165309-348f09dbbbc0 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/google/pprof v0.0.0-20230406165453-00490a63f317 // indirect
	github.com/iancoleman/strcase v0.2.0 // indirect
	github.com/kr/pretty v0.3.0 // indirect
	github.com/rogpeppe/go-internal v1.10.0 // indirect
	github.com/stretchr/testify v1.8.4 // indirect
	github.com/vektah/gqlparser/v2 v2.5.6 // indirect
	golang.org/x/mod v0.11.0 // indirect
	golang.org/x/net v0.11.0 // indirect
	golang.org/x/sys v0.9.0 // indirect
	golang.org/x/text v0.10.0 // indirect
	golang.org/x/tools v0.10.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

go 1.18

replace dagger.io/dagger => github.com/vito/dagger/sdk/go v0.0.0-20230713012757-685d8b7011ae

replace github.com/dagger/dagger => github.com/vito/dagger v0.0.0-20230713012757-685d8b7011ae
