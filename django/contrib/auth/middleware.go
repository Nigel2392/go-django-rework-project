package auth

import (
	"net/http"

	models "github.com/Nigel2392/django/contrib/auth/auth-models"
	"github.com/Nigel2392/django/core/except"
	"github.com/Nigel2392/mux"
	"github.com/Nigel2392/mux/middleware/authentication"
	"github.com/Nigel2392/mux/middleware/sessions"
)

const SESSION_COOKIE_NAME = "user_authentication"

func UnAuthenticatedUser() *models.User {
	return &models.User{
		IsLoggedIn: false,
	}
}

// Get the user from a request.
func UserFromRequest(r *http.Request) *models.User {

	var u = authentication.Retrieve(r)
	if u != nil {
		return u.(*models.User)
	}

	var session = sessions.Retrieve(r)
	except.Assert(
		session != nil,
		http.StatusInternalServerError,
		"Session must exist in the request",
	)

	var userID = session.Get(SESSION_COOKIE_NAME)
	if userID == nil {
		return UnAuthenticatedUser()
	}

	var uidInt, ok = userID.(uint64)
	if !ok {
		return UnAuthenticatedUser()
	}
	var user, err = Auth.Queries.GetUserById(r.Context(), uidInt)
	if err != nil || len(user) == 0 {
		return UnAuthenticatedUser()
	}

	var (
		loggedInUserData = user[0]
		loggedInuser     = loggedInUserData.User
	)
	loggedInuser.IsLoggedIn = true
	return &loggedInuser
}

func UserFromRequestPure(r *http.Request) authentication.User {
	return UserFromRequest(r)
}

// Add a user to a request, if one exists in the session.
func AddUserMiddleware() mux.Middleware {
	return authentication.AddUserMiddleware(UserFromRequestPure)
}
