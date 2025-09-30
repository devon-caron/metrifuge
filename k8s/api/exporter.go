package api

type Exporter interface {
	GetMetadata() Metadata
	AddRule(rule *Rule)
}
