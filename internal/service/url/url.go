package url

import (
	"golang.org/x/exp/slog"
	"io"
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
	const op = "service.url.Validate"

	domain = strings.TrimRight(domain, "/")

	log := this.log.With(
		slog.String("op", op),
		slog.String("domain", domain),
	)

	log.Info("starting validate urls")

	var wg sync.WaitGroup
	var validated sync.Map // add only string => int
	results := make(map[string]int)

	for _, url := range urls {
		url = "/" + strings.Trim(url, "/")

		if _, ok := results[url]; !ok {
			results[url] = -1

			wg.Add(1)

			go func(url string) {
				defer wg.Done()

				status, err := this.fetchUrlStatus(domain + url)
				if err != nil {
					log.Error("failed validate "+url, sl.Err(err))
					validated.Store(url, -1)
					return
				}

				validated.Store(url, status)
			}(url)
		}
	}

	wg.Wait()

	log.Info("finished validate urls")

	validated.Range(func(key, value any) bool {
		results[key.(string)] = value.(int)

		return true
	})

	return results
}

func (this *Url) fetchUrlStatus(fullUrl string) (int, error) {
	const op = "service.url.fetchUrlStatus"

	log := this.log.With(
		slog.String("op", op),
		slog.String("fullUrl", fullUrl),
	)

	resp, err := http.Get(fullUrl)

	if err != nil {
		return 0, err
	}

	func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Error("failed close body", sl.Err(err))
		}
	}(resp.Body)

	return resp.StatusCode, nil
}
