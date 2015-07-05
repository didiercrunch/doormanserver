epsilon = 0.0001

define ["app", "lodash", "../filters", "directives/ratinput"], (app, _, filters, ratinput) ->
    app.controller "NewDoormanCtrl", ["$scope", "$http", "$location", "rationals", ($scope, $http, $location, rationals) ->
        $scope.error = ""
        $scope.user = localStorage.getItem("email")
        $scope.value = {probability: "1", name: ""}
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
            $scope.value =
                name: ""
                probability: $scope.getProbabilityRemaining()

        $scope.removeValue = (idx) ->
            $scope.payload.values.splice(idx, 1);

        $scope.getProbabilityRemaining = () ->
            sum = rationals.sum(_.pluck($scope.payload.values, "probability")...)
            return rationals.minus("1", sum)


        $scope.redirectToDoorman = (doormanLocation) ->
            $http.get(doormanLocation).then (res) ->
                id = res.data.id
                return  $location.url("/doormen/#{id}")

        $scope.getPayload = () ->
            payload = _.cloneDeep($scope.payload)
            for value in payload.values
                value.probability = value.probability
            return payload

        $scope.createDoorman = () ->
            if not $scope.isDoormanValid()
                return
            params = {params: {user: $scope.user}}
            $http.post("/api/doormen", $scope.getPayload(), params).
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
            rationals.equal($scope.getProbabilityRemaining(), "0" )


        $scope.addEmail = () ->
            if not $scope.newEmail
                return
            $scope.payload.emails.push($scope.newEmail)
            $scope.newEmail = ""

        $scope.removeEmail = (idx) ->
            if $scope.payload.emails[idx] == $scope.user
                return
            $scope.payload.emails.splice(idx, 1);

        $scope.getSum = () ->
            probs = _.pluck($scope.payload.values, "probability")
            return rationals.pp(rationals.sum(probs...))



        return
    ]

    return
