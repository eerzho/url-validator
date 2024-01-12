package url

type Service interface {
	Validate(domain string, urls []string) map[string]int
}
