package url

import (
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"golang.org/x/exp/slog"
	"io"
	"net/http"
	resp "url-validator/internal/lib/api/response"
	"url-validator/internal/lib/logger/sl"
)

type Url struct {
	log     *slog.Logger
	service Service
}

func New(log *slog.Logger, service Service) *Url {
	return &Url{log: log, service: service}
}

// Validate godoc
// @Summary validate urls
// @Description validate the given urls
// @Tags Url
// @Accept json
// @Produce json
// @Param validateRequest body ValidateRequest true "Validate Request"
// @Success 200 {object} ValidateResponse
// @Router /api/urls/validate [post]
func (this *Url) Validate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.url.validate"

		log := this.log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req ValidateRequest
		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")

			render.JSON(w, r, resp.Error("empty request"))

			return
		}

		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Info("request body decode", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))

			render.JSON(w, r, resp.ValidationError(validateErr))

			return
		}

		render.JSON(w, r, ValidateResponse{
			Response:  resp.OK(),
			Validated: this.service.Validate(req.Domain, req.Urls),
		})
	}
}
