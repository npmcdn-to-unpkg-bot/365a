'use strict';

angular.module('Users', ['ngRoute'])

    .config(['$routeProvider', function ($routeProvider) {
        $routeProvider.when('/users/list', {
            templateUrl: 'users/list.html',
            controller: 'ListUsers'
        });
        $routeProvider.when('/users/create', {
            templateUrl: 'users/create.html',
            controller: 'CreateUsers'
        });

        $routeProvider.when('/users/edit/:id', {
            templateUrl: 'users/create.html',
            controller: 'EditUsers'
        });
    }])

    .controller('ListUsers', ['$scope', '$http', 'toastr', function ($scope, $http, toastr) {

 
        $scope.usertype=function (type){
            switch (type){
                case 1:
                    return "Administrador";
                case 2:
                    return "Supervisor";
                case 3:
                    return "Colaborador";
            }
            return "Invalido"
        };
       $scope.users = [];

        $scope.search=function (keywords){
            if (keywords==""){
                update();
                return

            }
            $http({
                method: 'GET',
                url: window.hostname + 'users/search/'+keywords
            }).then(function successCallback(response) {

                $scope.users = response.data;

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
                url: window.hostname + 'users/listall'
            }).then(function successCallback(response) {
               $scope.users = response.data;
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
                url: window.hostname + 'users/delete/' + id,
                data: $scope.user
            }).then(function successCallback(response) {
                update();
                toastr.success('Sucesso!', 'Colaborador Eliminado');
                // this callback will be called asynchronously
                // when the response is available
            }, function errorCallback(response) {
                toastr.error('Erro!', response.data);
                // called asynchronously if an error occurs
                // or server returns response with an error status.
            });

        }


    }])
    .controller('CreateUsers', ['$scope', '$http', 'toastr', function ($scope, $http, toastr) {
        $scope.user = {};
        $scope.user.Role=2;

        $scope.submit = function () {
            console.log($scope.user)
            $http({
                method: 'POST',
                url: window.hostname + 'users/create',
                data: $scope.user
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
    .controller('EditUsers', ['$scope', '$http', '$routeParams', 'toastr', function ($scope, $http, $routeParams, toastr) {

        var user = localStorage.getItem("user");

            $scope.current=JSON.parse(user).Data
        var id = $routeParams.id;
        $scope.user = {};

        $http({
            method: 'GET',
            url: window.hostname + 'users/getid/' + id
        }).then(function successCallback(response) {
            $scope.user = response.data;
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
                url: window.hostname + 'users/edit/' + id,
                data: $scope.user
            }).then(function successCallback(response) {
             
                toastr.success('Sucesso!', 'Publisher Edited');
                // this callback will be called asynchronously
                // when the response is available
            }, function errorCallback(response) {
                toastr.error('Erro!', response.data);

                // called asynchronously if an error occurs
                // or server returns response with an error status.
            });
        }
    }]);