'use strict';

angular.module('Clients', ['ngRoute'])

    .config(['$routeProvider', function ($routeProvider) {
        $routeProvider.when('/clients/list', {
            templateUrl: 'clients/list.html',
            controller: 'ListClients'
        });
        $routeProvider.when('/clients/create', {
            templateUrl: 'clients/create.html',
            controller: 'CreateClients'
        });

        $routeProvider.when('/clients/edit/:id', {
            templateUrl: 'clients/create.html',
            controller: 'EditClients'
        });
    }])

    .controller('ListClients', ['$scope', '$http', 'toastr', function ($scope, $http, toastr) {
       $scope.clients = [];

        $scope.search=function (keywords){
            if (keywords==""){
                update();
                return

            }
            $http({
                method: 'GET',
                url: window.hostname + 'clients/search/'+keywords
            }).then(function successCallback(response) {
                $scope.clients = response.data;
                // this callback will be called asynchronously
                // when the response is available
            }, function errorCallback(response) {
                toastr.error('Erro!', response.data);
                // called asynchronously if an error occurs
                // or server returns response with an error status.
            });


        };
        var update = function () {
            $http({
                method: 'GET',
                url: window.hostname + 'clients/listall'
            }).then(function successCallback(response) {
               $scope.clients = response.data;
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
                url: window.hostname + 'clients/delete/' + id,
                data: $scope.client
            }).then(function successCallback(response) {
                update();
                toastr.success('Sucesso!', 'Cliente Eliminado');
                // this callback will be called asynchronously
                // when the response is available
            }, function errorCallback(response) {
                toastr.error('Erro!', response.data);
                // called asynchronously if an error occurs
                // or server returns response with an error status.
            });

        }


    }])
    .controller('CreateClients', ['$scope', '$http', 'toastr', function ($scope, $http, toastr) {
        $scope.client = {};

        $scope.submit = function () {
            console.log($scope.client);
            $http({
                method: 'POST',
                url: window.hostname + 'clients/create',
                data: $scope.client
            }).then(function successCallback(response) {
                toastr.success('Sucesso!', 'Colaborador Adicionado');
                // this callback will be called asynchronously
                // when the response is available
            }, function errorCallback(response) {
                toastr.error('Erro!', response.data);
                // called asynchronously if an error occurs
                // or server returns response with an error status.
            });
        }

    }])
    .controller('EditClients', ['$scope', '$http', '$routeParams', 'toastr', function ($scope, $http, $routeParams, toastr) {

        var id = $routeParams.id;
        $scope.client = {};

        $http({
            method: 'GET',
            url: window.hostname + 'clients/getid/' + id
        }).then(function successCallback(response) {
            $scope.client = response.data;
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
                url: window.hostname + 'clients/edit/' + id,
                data: $scope.client
            }).then(function successCallback(response) {
                toastr.success('Sucesso!', 'Cliente Modificado');
                // this callback will be called asynchronously
                // when the response is available
            }, function errorCallback(response) {
                toastr.error('Erro!', response.data);
                // called asynchronously if an error occurs
                // or server returns response with an error status.
            });
        }
    }]);