define ["app"], (app) ->
    app.controller "DoormanErrorCtrl", ["$scope","$routeParams", ($scope, $routeParams) ->
        $scope.id = $routeParams.id
        $scope.error = Number($routeParams.error)
    ]
    return
