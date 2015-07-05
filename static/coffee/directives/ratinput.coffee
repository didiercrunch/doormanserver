template = """
<div class="row rat-input">
    <div class="row" ng-show="!edit" >
        <div class="columns large-12" style="text-align: center;">
            <span class="fraction-box" ng-click="enableEdit()">
                <span  ng-bind="nominator"></span> &frasl; <span ng-bind="denominator" ></span>
                <i class="fa fa-pencil-square-o" ng-show="!ngDisabled"></i>
            </span>
        </div>
    </div>
    <div  ng-show="edit" class="edit">
        <form  ng-submit="edit=false">
            <div class="large-12 columns" style="text-align: center;">
                <input ng-model="nominator" type="text" ng-change="changeModel()"  ng-disabled="ngDisabled">
                <span style="text-align: center;">&frasl;<span>
                <input ng-model="denominator" type="text" ng-change="changeModel()" ng-disabled="ngDisabled">
                <i class="fa fa-check" ng-click="edit=false"></i>
                <input type="submit" style="display:none;">
            </div>
        </form>
    </div>
</div>
"""


define ["app", "angularAMD", "lodash", "services/rationals"], (app, angularAMD, _, rationals) ->

    angularAMD.directive "ratInput", () ->
        directive =
                restrict: 'E',
                template: template
                scope:
                    ngModel: "="
                    ngDisabled: "="
                controller: [ '$scope', "rationals", ($scope, rationals) ->

                    $scope.getNominator = ()->
                        if _.contains($scope.ngModel, "/")
                            return Number($scope.ngModel.substr(0, $scope.ngModel.indexOf("/")))
                        Number($scope.ngModel)

                    $scope.getDenominator = () ->
                        if _.contains($scope.ngModel, "/")
                            return Number($scope.ngModel.substr($scope.ngModel.indexOf("/") + 1))
                        1

                    $scope.changeModel = ()->
                        $scope.ngModel = rationals.newRat($scope.nominator, $scope.denominator)

                    $scope.enableEdit = () ->
                        if not $scope.ngDisabled
                            $scope.edit = true

                    $scope.init = () ->
                        $scope.nominator = rationals.getNominator($scope.ngModel)
                        $scope.denominator = $scope.getDenominator()

                    $scope.$watch 'ngModel', () ->
                        $scope.init()

                    $scope.init()

                ]
        return directive
