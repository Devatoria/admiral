package jobs

var Jobs map[string]func() error

func init() {
	Jobs = make(map[string]func() error)
	Jobs["synchronize"] = SynchronizeCatalog
}
