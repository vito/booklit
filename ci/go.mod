module ci

go 1.20

replace dagger.io/dagger => github.com/vito/dagger/sdk/go v0.0.0-20230827015344-eb8dc6ffbff9

replace github.com/dagger/dagger => github.com/vito/dagger v0.0.0-20230827015344-eb8dc6ffbff9

require (
	dagger.io/dagger v0.7.2
	github.com/Khan/genqlient v0.6.0
	github.com/iancoleman/strcase v0.3.0
	github.com/vektah/gqlparser/v2 v2.5.14
)

require (
	github.com/99designs/gqlgen v0.17.31 // indirect
	golang.org/x/sync v0.3.0 // indirect
)
