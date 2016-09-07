package jobs

// Jobs map containing all runnable jobs
var Jobs map[string]func([]string) error

func init() {
	Jobs = make(map[string]func([]string) error)
	Jobs["synchronize"] = SynchronizeCatalog
	Jobs["create_admin"] = CreateAdminUser
}
