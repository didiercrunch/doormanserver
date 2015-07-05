define ["app", "angularAMD"], (app, angularAMD) ->
    app.controller "ApispecCtrl",['$scope', '$http', ($scope, $http) ->
        $http.get("/api/documentation").then (res)->
            $scope.endpoints = res.data.endpoints

        $scope.methods =
            get: "GET"
            post: "POST"
    ]

    angularAMD.directive "wabMethod", () ->
        directive =
                restrict: 'E',
                scope:
                  method: '@'
                link: (scope, element) ->
                    scope.methods =
                        get: "GET"
                        put: "PUT"
                        delete: "DELETE"
                        post: "POST"

                template: '<code style="width:5em; display: block; text-align: center;">{{ methods[method] }}</code>'
        return directive

    angularAMD.directive "wabApi", () ->
        directive =
            restrict: 'E',
            scope:
                endpoint: '@'
            template: '<a>/api{{ endpoint }}</a>'
        return directive
