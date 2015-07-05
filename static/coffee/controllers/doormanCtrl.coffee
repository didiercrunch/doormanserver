define ["app", "lodash", "directives/ratinput"], (app, _, ratinput) ->

    app.controller "DoormanCtrl", ["$scope", "$http","$routeParams","$location", ($scope, $http, $routeParams, $location) ->
        $scope.doorman = {values: []}
        $scope.somevalue = "2/6"
        $scope.apiurl = "/api/doormen/#{$routeParams.id}"
        $scope.user = localStorage.getItem("email")

        $http.get($scope.apiurl).then((res) ->
            $scope.doorman = res.data
            $scope.isAdmin = _.contains($scope.doorman.emails, $scope.user)
            $scope.initEmails()

        ).catch((res) ->
            if res.status == 404
                $location.path("/doormen/error/#{res.status}/#{$routeParams.id}")
        )

        $scope.deleteDoorman = () ->
            console.log("not yet implemented")

        $scope.initEmails = () ->
            [$scope.leftEmailColumns, $scope.rightEmailColumns] = _.chunk($scope.doorman.emails, 2)

        $scope.rationalToFloat = (rat) ->
            if _.contains(rat, "/")
                num = Number(rat.substr(0, rat.indexOf("/")))
                denum = Number(rat.substr(rat.indexOf("/") + 1))
                return num / denum
            Number(rat)


        $scope.validate = ()->
            ret = 0
            for val in $scope.doorman.values
                ret += $scope.rationalToFloat(val.probability)
            return Math.abs(ret - 1.0) < 0.000001

        $scope.addNewAuthorEmail = () ->
            if not $scope.canAddNewEmail()
                return
            $scope.doorman.emails.push($scope.newEmail)
            $scope.updateDoorman()
            $scope.newEmail = ""
            $scope.initEmails()

        $scope.removeEmail = (email) ->
            $scope.doorman.emails = _.reject($scope.doorman.emails , (email_) -> email_ == email)
            $scope.initEmails()
            $scope.updateDoorman()

        $scope.canAddNewEmail = ()->
            if !$scope.newEmail or $scope.newEmail.length < 5
                return false
            not _.contains($scope.doorman.emails, $scope.newEmail)


        $scope.updateDoorman = () ->
            if !$scope.validate()
                return
            params = {params: {user: $scope.user}}
            $http.put($scope.apiurl, $scope.doorman, params).then((ret)->
                console.log "update"
            ).catch((res) ->
                $scope.error =
                    message: res.data
                    status : res.status
            )
    ]
    return
