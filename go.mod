module crow/oraiplayground

go 1.22.3

replace crow/orai => /home/crow/repos/orai

require (
	crow/orai v0.0.0-00010101000000-000000000000
	crow/webutil v0.0.0-00010101000000-000000000000
	github.com/gorilla/mux v1.8.1
)

require gopkg.in/yaml.v3 v3.0.1

replace crow/webutil => /home/crow/repos/webutil
