package apps

type Metadata struct {
	Namespace  string
	Id         string
	Env        string
	TimeFormat string
}

func NewMetadata() *Metadata {
	m := &Metadata{}
	m.TimeFormat = "2006-01-02T15:04:05.999999"
	return m
}
