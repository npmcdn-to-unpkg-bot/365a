'use strict';

angular.module('myApp.Auth', ['ngRoute'])

    .config(['$routeProvider', function ($routeProvider) {
        $routeProvider.when('/auth/login', {
            templateUrl: 'auth/login.html',
            controller: 'Login'
        });
    }])

    .controller('Login', ['$scope', '$http', '$location', 'toastr',function ($scope, $http, $location,toastr) {
        $scope.auth = {};

        $scope.submit = function () {
            $http({
                method: 'POST',
                url: window.hostname + 'auth/login',
                data: $scope.auth
            }).then(function successCallback(response) {
                console.log(response.data)

                localStorage.setObject("user", response.data)

                $location.path('/')
                // this callback will be called asynchronously
                // when the response is available
            }, function errorCallback(response) {
                console.log(response)
                toastr.error('Erro!', response.data);
                // called asynchronously if an error occurs
                // or server returns response with an error status.
            });

        }


    }]);