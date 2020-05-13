module github.com/vito/booklit

require (
	github.com/agext/levenshtein v1.2.3
	github.com/alecthomas/chroma v0.7.3
	github.com/alecthomas/repr v0.0.0-20181024024818-d37bc2a10ba1 // indirect
	github.com/fsnotify/fsnotify v1.4.7 // indirect
	github.com/go-bindata/go-bindata v3.1.2+incompatible
	github.com/golang/protobuf v1.3.1 // indirect
	github.com/hpcloud/tail v1.0.0 // indirect
	github.com/jessevdk/go-flags v1.4.0
	github.com/konsorten/go-windows-terminal-sequences v1.0.2 // indirect
	github.com/mna/pigeon v1.0.1-0.20200224192238-18953b277063
	github.com/onsi/ginkgo v1.6.0
	github.com/onsi/gomega v1.4.1
	github.com/russross/blackfriday/v2 v2.0.1
	github.com/segmentio/textio v1.2.0
	github.com/shurcooL/sanitized_anchor_name v1.0.0 // indirect
	github.com/sirupsen/logrus v1.4.1
	github.com/yuin/goldmark v1.1.30
	golang.org/x/tools v0.0.0-20200505023115-26f46d2f7ef8 // indirect
	gopkg.in/fsnotify.v1 v1.4.7 // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
	gopkg.in/yaml.v2 v2.2.8 // indirect
)

replace github.com/yuin/goldmark => ./goldmark

go 1.13
