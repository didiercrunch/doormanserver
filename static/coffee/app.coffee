
define [
    "angularAMD"
    "angular-route"
    "foundation"
], (angularAMD) ->
    app = angular.module("webapp", ["ngRoute", "mm.foundation"])
    app.config ($routeProvider) ->
        $routeProvider.when("/", angularAMD.route(
            templateUrl: "partials/home.html"
            controller: "HomeCtrl"
            controllerUrl: "controllers/homeCtrl"
        )).when("/serverspec", angularAMD.route(
            templateUrl: "partials/serverspec.html"
            controller: "ServerspecCtrl"
            controllerUrl: "controllers/serverspecCtrl"
        )).when("/apispec", angularAMD.route(
            templateUrl: "partials/apispec.html"
            controller: "ApispecCtrl"
            controllerUrl: "controllers/apispecCtrl"
        )).when("/newdoorman", angularAMD.route(
            templateUrl: "partials/newdoorman.html"
            controller: "NewDoormanCtrl"
            controllerUrl: "controllers/newdoormanCtrl"
        )).when("/doormen", angularAMD.route(
            templateUrl: "partials/doormen.html"
            controller: "DoormenCtrl"
            controllerUrl: "controllers/doormenCtrl"
        )).when("/doormen/error/:error/:id", angularAMD.route(
            templateUrl: "partials/doormanError.html"
            controller: "DoormanErrorCtrl"
            controllerUrl: "controllers/doormanErrorCtrl"
        )).when("/doormen/:id", angularAMD.route(
            templateUrl: "partials/doorman.html"
            controller: "DoormanCtrl"
            controllerUrl: "controllers/doormanCtrl"
        )).otherwise redirectTo: "/"
        return

        serverspec

    angularAMD.bootstrap(app)
