package main

import (
	"net/http"

	"github.com/Martins-Iroka/MyGallery-Backend/internal/util"
)

// healthcheckHandler godoc
//
//	@Summary		Healthcheck
//	@Description	Healthcheck endpoint
//	@Tags			ops
//	@Produce		json
//	@Success		200	{object}	string	"ok"
//	@Router			/health [get]
func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":  "ok",
		"version": version,
	}

	if err := util.JsonResponse(w, http.StatusOK, data); err != nil {
		util.InternalServerErrorResponse(w, r, err, app.logger)
	}
}
