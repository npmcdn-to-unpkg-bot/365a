'use strict';

angular.module('Modules', ['ngRoute'])

    .config(['$routeProvider', function ($routeProvider) {
        $routeProvider.when('/modules/list', {
            templateUrl: 'modules/list.html',
            controller: 'ListModules'
        });
        $routeProvider.when('/modules/create', {
            templateUrl: 'modules/create.html',
            controller: 'CreateModules'
        });

        $routeProvider.when('/modules/edit/:id', {
            templateUrl: 'modules/create.html',
            controller: 'EditModules'
        });
    }])

    .controller('ListModules', ['$scope', '$http', 'toastr', function ($scope, $http, toastr) {
       $scope.modules = [];

        
        var update = function () {
            $http({
                method: 'GET',
                url: window.hostname + 'modules/listallmodule'
            }).then(function successCallback(response) {
               $scope.modules = response.data;
                // this callback will be called asynchronously
                // when the response is available
            }, function errorCallback(response) {
                toastr.error('Erro!', response.data);
                // called asynchronously if an error occurs
                // or server returns response with an error status.
            });
        };
        update()
        $scope.delete = function (id) {

            $http({
                method: 'POST',
                url: window.hostname + 'modules/deletemodule/' + id,
                data: $scope.module
            }).then(function successCallback(response) {
                $scope.$root.$broadcast("updatemenu");

                update();
                toastr.success('Sucesso!', 'Modulo Eliminado');
                // this callback will be called asynchronously
                // when the response is available
            }, function errorCallback(response) {
                toastr.error('Erro!', response.data);
                // called asynchronously if an error occurs
                // or server returns response with an error status.
            });

        }


    }])
    .controller('CreateModules', ['$scope', '$http', 'toastr', function ($scope, $http, toastr) {
        $scope.module = {};
        $scope.module.Fields=[];


        $scope.addField=function(){
            $scope.module.Fields.push({})
        };

        $scope.deleteField=function(i){
            $scope.module.Fields.splice(i, 1);
        };
        $scope.submit = function () {
            console.log($scope.client);
            $http({
                method: 'POST',
                url: window.hostname + 'modules/createmodule',
                data: $scope.module
            }).then(function successCallback(response) {
                toastr.success('Sucesso!', 'Modulo Adicionado');
                $scope.$root.$broadcast("updatemenu");
                // this callback will be called asynchronously
                // when the response is available
            }, function errorCallback(response) {
                toastr.error('Erro!', response.data);
                // called asynchronously if an error occurs
                // or server returns response with an error status.
            });
        }

    }])
    .controller('EditModules', ['$scope', '$http', '$routeParams', 'toastr', function ($scope, $http, $routeParams, toastr) {

        var id = $routeParams.id;
        $scope.module = {};
        $scope.module.Fields=[];


        $scope.addField=function(){
            $scope.module.Fields.push({})
        };

        $scope.deleteField=function(i){
            $scope.module.Fields.splice(i, 1);
        };

        $http({
            method: 'GET',
            url: window.hostname + 'modules/getidmodule/' + id
        }).then(function successCallback(response) {
            $scope.module = response.data;
           /* if (!$scope.module.Fields){
                $scope.module.Fields=[];
            }*/


            // this callback will be called asynchronously
            // when the response is available
        }, function errorCallback(response) {
            toastr.error('Erro!', response.data);
            // called asynchronously if an error occurs
            // or server returns response with an error status.
        });

        $scope.submit = function () {
            $http({
                method: 'POST',
                url: window.hostname + 'modules/editmodule/' + id,
                data: $scope.module
            }).then(function successCallback(response) {
                toastr.success('Sucesso!', 'Modulo Modificado');
                $scope.$root.$broadcast("updatemenu");
                // this callback will be called asynchronously
                // when the response is available
            }, function errorCallback(response) {
                toastr.error('Erro!', response.data);
                // called asynchronously if an error occurs
                // or server returns response with an error status.
            });
        }
    }]);