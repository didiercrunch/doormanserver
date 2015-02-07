define ["app"], (app) ->
    app.controller "DoormenCtrl",[ "$scope", "$http", ($scope, $http) ->
        $scope.doormen = []
        $http.get("/api/doormen").then (res) ->
            $scope.doormen = res.data.doormen
        return
    ]

    return
