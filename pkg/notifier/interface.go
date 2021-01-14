package notifier

// Notifier is the interface to comment on the PR page, could be implemented by different
// provider, like GihHub, Gitlab, and etc.
type Notifier interface {
	Comment(owner string, repo string, number int, msg string) error
}
