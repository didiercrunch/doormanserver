package api

import (
	"github.com/didiercrunch/doormanserver/simpleapi"
	"net/http"
)

func Create() http.Handler {
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
		"Get the doorman.",
		"GET",
		GetDoorman,
	}

	updateDoorman := &simpleapi.Endpoint{
		"/doormen/{id}",
		"Update the doorman with new probabilities.",
		"PUT",
		UpdateDoorman,
	}

	api := simpleapi.New(
		"/api",
		serverSpecification,
		createDoorman,
		allDoorman,
		doorman,
		updateDoorman,
	)

	return api
}
