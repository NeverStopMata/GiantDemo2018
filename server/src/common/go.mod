module common

go 1.15

replace base => ../base

require (
	base v0.0.0-00010101000000-000000000000
	github.com/gogo/protobuf v1.3.2
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
