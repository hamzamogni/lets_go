package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"net/http"
)

// The routes() method returns a servemux containing our application routes.
func (app *application) routes() http.Handler {
	// Initialize the router.
	router := httprouter.New()

	// Create a handler function which wraps our notFound() helper, and then
	// assign it as the custom handler for 404 Not Found responses. You can also
	// set a custom handler for 405 Method Not Allowed responses by setting
	// router.MethodNotAllowed in the same way too.
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	// Update the pattern for the route for the static files.
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))
	// And then create the routes using the appropriate methods, patterns and
	// handlers.
	router.HandlerFunc(http.MethodGet, "/", app.home)
	router.HandlerFunc(http.MethodGet, "/snippet/view/:id", app.snippetView)
	router.HandlerFunc(http.MethodGet, "/snippet/create", app.snippetCreateForm)
	router.HandlerFunc(http.MethodPost, "/snippet/create", app.snippetCreatePost)

	// Create a middleware chain containing our 'standard' middleware
	// which will be used for every request our application receives.
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	router.Handler(http.MethodGet, "/user/signup", standard.ThenFunc(app.userSignup))
	router.Handler(http.MethodPost, "/user/signup", standard.ThenFunc(app.userSignupPost))
	router.Handler(http.MethodGet, "/user/login", standard.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/login", standard.ThenFunc(app.userLoginPost))
	router.Handler(http.MethodPost, "/user/logout", standard.ThenFunc(app.userLogoutPost))

	return standard.Then(router)
}
