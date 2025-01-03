package main

import (
	"net/http"
	"p1/internal/store"
)

// GetUserFeedHandler godoc
//
//	@Summary		Fetches the user feed
//	@Description	Fetches the user feed
//	@Tags			feed
//	@Accept			json
//	@Produce		json
//	@Param			since	query		string	false	"Since"
//	@Param			until	query		string	false	"Until"
//	@Param			limit	query		int		false	"Limit"
//	@Param			offset	query		int		false	"Offset"
//	@Param			sort	query		string	false	"Sort"
//	@Param			tags	query		string	false	"Tags"
//	@Param			search	query		string	false	"Search"
//	@Success		200		{object}	[]store.PostWithMetadata
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/feed	[get]
func (app *application) getFeedHandler(w http.ResponseWriter, r *http.Request) {
	pf := store.PaginatedFeedQuery{
		Limit:  20,
		Offset: 0,
		Sort:   "desc",
	}

	pf, err := pf.Parse(r)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	if err := Validate.Struct(pf); err != nil {
		app.badRequest(w, r, err)
		return
	}

	feed, err := app.store.Posts.GetFeed(r.Context(), int64(2), pf)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, feed); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
