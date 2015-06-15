define ["app"], (app) ->
    app.controller "HomeCtrl", ($scope) ->
        $scope.message = "message from home ctrl"
        $scope.loggedEmail = localStorage.getItem("email")
        $scope.login = () ->
            if $scope.email
                localStorage.setItem("email", $scope.email)
                $scope.loggedEmail = $scope.email
                $scope.email = ""

    return
