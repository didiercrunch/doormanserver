
define [
    "angularAMD"
    "lodash"
    "angular-route"
    "foundation"

], (angularAMD, _) ->
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
        )).otherwise({redirectTo: "/"})

        document.getElementsByTagName("body")[0].style.visibility = ""
        return

    app.controller "topbarCtrl", ["$scope", "$rootScope", "$location", ($scope, $rootScope, $location) ->

        $scope.init = () ->
            $scope.email = localStorage.getItem("email")
            $scope.loggedIn = !!$scope.email
            if !$scope.loggedIn
                $location.path('/').replace()

        $scope.logout = () ->
            localStorage.setItem("email", "")
            $scope.email = ""
            $scope.loggedIn = false
            $scope.$emit("logout")

        $scope.init()
        $rootScope.$on('login', $scope.init)
        return
    ]


    angularAMD.bootstrap(app)
