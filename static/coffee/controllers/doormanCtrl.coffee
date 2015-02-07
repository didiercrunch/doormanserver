define ["app"], (app) ->
    app.controller "DoormanCtrl", ["$scope", "$http","$routeParams","$location", ($scope, $http, $routeParams, $location) ->
        $scope.doorman = {values: []}
        $scope.apiurl = "/api/doormen/#{$routeParams.id}"
        $http.get($scope.apiurl).then((res) ->
            $scope.doorman = res.data
        ).catch((res) ->
            if res.status == 404
                $location.path("/doormen/error/#{res.status}/#{$routeParams.id}")
        )

        $scope.deleteDoorman = () ->
            console.log("not yet implemented")

        $scope.validate = ()->
            ret = 0
            for val in $scope.doorman.values
                ret += val.probability
            return Math.abs(ret - 1.0) < 0.000001

        $scope.updateDoorman = () ->
            if !$scope.validate()
                return
            $http.put($scope.apiurl, $scope.doorman).then((ret)->
                console.log "update"
            ).catch((res) ->
                d = 9
            )
    ]
    return
