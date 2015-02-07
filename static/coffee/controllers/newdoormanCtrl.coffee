epsilon = 0.0001

define ["app"], (app) ->
    app.controller "NewDoormanCtrl", ["$scope", "$http", "$location", ($scope, $http, $location) ->
        $scope.error = ""
        $scope.value = {}
        $scope.payload =
            name: ""
            values: []

        $scope.isValidValue = () ->
            return $scope.value.name and $scope.value.probability

        $scope.addNewValue = () ->
            if !$scope.isValidValue()
                return
            $scope.payload.values.push($scope.value)
            $scope.value = {name: "", probability: 0}

        $scope.removeValue = (idx) ->
            $scope.payload.values.splice(idx, 1);

        $scope.getProbabilityRemaining = () ->
            ret = 0
            for v in $scope.payload.values
                ret += v.probability
            return 1 - ret

        $scope.redirectToDoorman = (doormanLocation) ->
            $http.get(doormanLocation).then (res) ->
                id = res.data.id
                return  $location.url("/doormen/#{id}")

        $scope.createDoorman = () ->
            if not $scope.isDoormanValid()
                return
            $http.post("/api/doormen", $scope.payload).
            then((res, s) ->
                $scope.redirectToDoorman(res.headers("location"))
            ).catch((res) ->
                $scope.error = res.data
            )

        $scope.isDoormanValid = () ->
            if not $scope.payload.name
                return false
            if $scope.payload.values.length < 2
                return false
            if Math.abs($scope.getProbabilityRemaining()) > epsilon
                return false
            return true


        return
    ]

    return
