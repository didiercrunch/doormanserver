epsilon = 0.0001

define ["app", "../filters"], (app, filters) ->
    app.controller "NewDoormanCtrl", ["$scope", "$http", "$location",  ($scope, $http, $location) ->
        $scope.error = ""
        $scope.user = localStorage.getItem("email")
        $scope.value = {probability: 0, name: ""}
        $scope.payload =
            name: ""
            values: []
            emails: [$scope.user]

        $scope.isValidValue = () ->
            return $scope.value.name

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
            params = {params: {user: $scope.user}}
            $http.post("/api/doormen", $scope.payload, params).
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

        $scope.addEmail = () ->
            if not $scope.newEmail
                return
            $scope.payload.emails.push($scope.newEmail)
            $scope.newEmail = ""

        $scope.removeEmail = (idx) ->
            if $scope.payload.emails[idx] == $scope.user
                return
            $scope.payload.emails.splice(idx, 1);



        return
    ]

    return
