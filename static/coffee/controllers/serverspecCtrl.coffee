define ["app"], (app) ->
    app.controller "ServerspecCtrl",["$scope", "$http", ($scope, $http) ->
        $scope.server = {}
        $http.get("/api/server").then (res) -> $scope.server = res.data
        return
    ]
    return
