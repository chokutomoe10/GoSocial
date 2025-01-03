package main

import (
	"net/http"
	"p1/internal/store"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type userKey string

const userCtx userKey = "user"

// GetUser godoc
//
//	@Summary		Fetches a user profile
//	@Description	Fetches a user profile by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object}	store.User
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/{id}	[get]
func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "userID")

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil || id < 1 {
		app.badRequest(w, r, err)
		return
	}

	user, err := app.getUser(r.Context(), id)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFound(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}

	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// FollowUser godoc
//
//	@Summary		Follows a user
//	@Description	Follows a user by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int		true	"User ID"
//	@Success		204		{string}	string	"User followed"
//	@Failure		400		{object}	error	"User payload missing"
//	@Failure		404		{object}	error	"User not found"
//	@Failure		409		{object}	error	"User already followed"
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/{userID}/follow	[put]
func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	followerUser := getUserFromCtx(r)
	followedID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	if err := app.store.Followers.Follow(r.Context(), followerUser.ID, followedID); err != nil {
		switch err {
		case store.ErrConflict:
			app.conflict(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}

	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// UnfollowUser godoc
//
//	@Summary		Unfollows a user
//	@Description	Unfollows a user by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int		true	"User ID"
//	@Success		204		{string}	string	"User unfollowed"
//	@Failure		400		{object}	error	"User payload missing"
//	@Failure		404		{object}	error	"User not found"
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/{userID}/unfollow	[put]
func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	unfollowedUser := getUserFromCtx(r)

	followerID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	if err := app.store.Followers.Unfollow(r.Context(), unfollowedUser.ID, followerID); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// ActivateUser godoc
//
//	@Summary		Activates/Register a user
//	@Description	Activates/Register a user by invitation token
//	@Tags			users
//	@Produce		json
//	@Param			token	path		string	true	"Invitation token"
//	@Success		204		{string}	string	"User activated"
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/activate/{token}	[put]
func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")

	if err := app.store.Users.Activate(r.Context(), token); err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFound(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusNoContent, ""); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// func (app *application) userContextMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		idParam := chi.URLParam(r, "userID")

// 		id, err := strconv.ParseInt(idParam, 10, 64)
// 		if err != nil {
// 			app.badRequest(w, r, err)
// 			return
// 		}

// 		ctx := r.Context()

// 		user, err := app.store.Users.GetByID(ctx, id)
// 		if err != nil {
// 			switch err {
// 			case store.ErrNotFound:
// 				app.notFound(w, r, err)
// 				return
// 			default:
// 				app.internalServerError(w, r, err)
// 				return
// 			}
// 		}

// 		ctx = context.WithValue(ctx, userCtx, user)

// 		next.ServeHTTP(w, r.WithContext(ctx))
// 	})
// }

func getUserFromCtx(r *http.Request) *store.User {
	user, _ := r.Context().Value(userCtx).(*store.User)
	return user
}
