define ["app"], (app) ->
    app.controller "HomeCtrl", ['$scope', '$rootScope', ($scope, $rootScope) ->
        $scope.init = () ->
            $scope.loggedEmail = localStorage.getItem("email")

        $scope.login = () ->
            if $scope.email
                localStorage.setItem("email", $scope.email)
                $scope.loggedEmail = $scope.email
                $scope.email = ""
                $scope.$emit("login")

        $scope.init()
        $rootScope.$on('logout', $scope.init)
    ]
