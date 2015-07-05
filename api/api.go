package api

import (
	"github.com/didiercrunch/doormanserver/simpleapi"
	"net/http"
)

func Create(root string) http.Handler {
	serverSpecification := &simpleapi.Endpoint{
		"/server",
		"Get the server specification.",
		"GET",
		GetServerSpecification,
	}

	createDoorman := &simpleapi.Endpoint{
		"/doormen",
		"Create a new doorman.",
		"POST",
		CreateDoorman,
	}

	allDoorman := &simpleapi.Endpoint{
		"/doormen",
		"Get the list of all the public doormen.",
		"GET",
		GetAllDoormen,
	}

	doorman := &simpleapi.Endpoint{
		"/doormen/{id}",
		"Get a specific doorman by id.",
		"GET",
		GetDoorman,
	}

	updateDoorman := &simpleapi.Endpoint{
		"/doormen/{id}",
		`Update the doorman with new probabilities or owner emails.  The payload
		must be a valid doorman with the same id than the targeted doorman.
		The values of the PUT doorman will replace the values of the old
		ones.`,
		"PUT",
		UpdateDoorman,
	}

	getDoormanStatus := &simpleapi.Endpoint{
		"/doormen/{id}/status",
		`Returns the status of the doorman as a doorman client expect it.`,
		"GET",
		GetDoormanStatus,
	}

	api := simpleapi.New(
		root,
		serverSpecification,
		createDoorman,
		allDoorman,
		getDoormanStatus,
		doorman,
		updateDoorman,
	)
	api.EnableDocumentation("documentation")
	api.AddMiddlewares(AuthenticationMiddleware)
	api.InitRouter()

	return api
}
