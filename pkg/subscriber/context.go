package subscriber

import "github.com/gocraft/work"

const (
	// RedisNamespace is redis namespace
	RedisNamespace = "saverbate_downloads"
	// JobName is job name for periodically enqueuing
	JobName = "follow_performers"
)

// Context is context of subscriber service
type Context struct{}

// FollowPerformers opens performer page and click Follow button
func (ctx *Context) FollowPerformers(job *work.Job) error {
	return nil
}
