package url

import (
	"golang.org/x/exp/slog"
	"net/http"
	"strings"
	"sync"
	"url-validator/internal/lib/logger/sl"
)

type Url struct {
	log *slog.Logger
}

func New(log *slog.Logger) *Url {
	return &Url{
		log: log,
	}
}

func (this *Url) Validate(domain string, urls []string) map[string]int {
	const op = "service.url.validate"

	domain = strings.TrimRight(domain, "/")

	log := this.log.With(
		slog.String("op", op),
		slog.String("domain", domain),
	)

	var wg sync.WaitGroup
	mutex := sync.Mutex{}

	validated := make(map[string]int)

	for _, url := range urls {
		url = "/" + strings.Trim(url, "/")

		mutex.Lock()
		_, exists := validated[url]
		mutex.Unlock()

		if !exists {
			wg.Add(1)

			go func(url string) {

				log.Info("starting validate " + url)
				status, err := this.fetchURLStatus(domain + url)

				mutex.Lock()
				if err != nil {
					log.Error("failed validate "+url, sl.Err(err))
				} else {
					validated[url] = status
				}
				mutex.Unlock()

				log.Info("finished validate " + url)
				wg.Done()

			}(url)
		}
	}

	wg.Wait()

	return validated
}

func (this *Url) fetchURLStatus(fullUrl string) (int, error) {
	resp, err := http.Get(fullUrl)

	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()

	return resp.StatusCode, nil
}
