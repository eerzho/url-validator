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

func (this *Url) Validate(domain string, urls []string) map[string]int32 {
	const op = "service.url.validate"

	domain = strings.TrimRight(domain, "/")

	log := this.log.With(
		slog.String("op", op),
		slog.String("domain", domain),
	)

	var wg sync.WaitGroup
	var validated sync.Map

	for _, originalUrl := range urls {
		url := "/" + strings.Trim(originalUrl, "/")

		if _, ok := validated.Load(url); !ok {

			log.Info("start validate " + url)
			wg.Add(1)

			go func(url string) {
				defer wg.Done()

				status, err := this.fetchURLStatus(domain + url)
				if err != nil {
					log.Error("failed validate "+url, sl.Err(err))
					validated.Store(url, -1)
					return
				}

				validated.Store(url, status)
				log.Info("finish validate " + url)
			}(url)
		}
	}

	wg.Wait()

	results := make(map[string]int32)
	validated.Range(func(key, value any) bool {
		results[key.(string)] = value.(int32)

		return true
	})

	return results
}

func (this *Url) fetchURLStatus(fullUrl string) (int, error) {
	resp, err := http.Get(fullUrl)

	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()

	return resp.StatusCode, nil
}
