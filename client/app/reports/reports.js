'use strict';

angular.module('reports', ['ngRoute'])

    .config(['$routeProvider', function ($routeProvider) {
        $routeProvider.when('/reports/list', {
            templateUrl: 'reports/list.html',
            controller: 'Listreports'
        });
        $routeProvider.when('/reports/create', {
            templateUrl: 'reports/create.html',
            controller: 'Createreports'
        });

        $routeProvider.when('/reports/edit/:id', {
            templateUrl: 'reports/create.html',
            controller: 'Editreports'
        });
    }])

    .controller('Listreports', ['$scope', '$http', 'toastr', function ($scope, $http, toastr) {
       $scope.reports = [];

        $scope.search=function (keywords){
            if (keywords==""){
                update();
                return

            }
            $http({
                method: 'GET',
                url: window.hostname + 'reports/search/'+keywords
            }).then(function successCallback(response) {

                $scope.reports = response.data;

                // this callback will be called asynchronously
                // when the response is available
            }, function errorCallback(response) {
                console.log(response);
                // called asynchronously if an error occurs
                // or server returns response with an error status.
            });


        };

        var update = function () {
            $http({
                method: 'GET',
                url: window.hostname + 'reports/listall'
            }).then(function successCallback(response) {
                console.log(response);
               $scope.reports = response.data;
                // this callback will be called asynchronously
                // when the response is available
            }, function errorCallback(response) {
                toastr.error('Erro!', response.data);
                // called asynchronously if an error occurs
                // or server returns response with an error status.
            });
        };
        update();
        $scope.delete = function (id) {

            $http({
                method: 'POST',
                url: window.hostname + 'reports/delete/' + id,
                data: $scope.report
            }).then(function successCallback(response) {
                update();
                toastr.success('Sucesso!', 'reporte Eliminado');
                // this callback will be called asynchronously
                // when the response is available
            }, function errorCallback(response) {
                toastr.error('Erro!', response.data);
                // called asynchronously if an error occurs
                // or server returns response with an error status.
            });

        }


    }])
    .controller('Createreports', ['$scope', '$http', 'toastr', function ($scope, $http, toastr) {
        $scope.report = {};
        $scope.clients = [];
        $scope.users= [];
        $scope.report.Attachments=[];


        $scope.addAttach=function(){
            $scope.report.Attachments.push({})
        };

        $scope.deleteAttach=function(i){
            $scope.report.Attachments.splice(i, 1);
        };
        function update(){
            $http({
                method: 'GET',
                url: window.hostname + 'clients/listall'
            }).then(function successCallback(response) {
              $scope.clients=response.data;

                // this callback will be called asynchronously
                // when the response is available
            }, function errorCallback(response) {
                toastr.error('Erro!', response.data);
                // called asynchronously if an error occurs
                // or server returns response with an error status.
            });

            $http({
                method: 'GET',
                url: window.hostname + 'users/listall'
            }).then(function successCallback(response) {
                $scope.users=response.data;
                // this callback will be called asynchronously
                // when the response is available
            }, function errorCallback(response) {
                toastr.error('Erro!', response.data);
                // called asynchronously if an error occurs
                // or server returns response with an error status.
            });
            
        }
        update();

        $scope.submit = function () {
            $http({
                method: 'POST',
                url: window.hostname + 'reports/create',
                data: $scope.report
            }).then(function successCallback(response) {
                console.log(response)
                toastr.success('Sucesso!', 'Relatório Adicionado');
                // this callback will be called asynchronously
                // when the response is available
            }, function errorCallback(response) {
                toastr.error('Erro!', response.data);
                // called asynchronously if an error occurs
                // or server returns response with an error status.
            });
        }

    }])
    .controller('Editreports', ['$scope', '$http', '$routeParams', 'toastr', function ($scope, $http, $routeParams, toastr) {
        $scope.report = {};
        $scope.clients = [];
        $scope.users= [];
        $scope.report.Attachments=[];
        var id = $routeParams.id;


        $scope.addAttach=function(){
            $scope.report.Attachments.push({})
        };

        $scope.deleteAttach=function(i){
            $scope.report.Attachments.splice(i, 1);
        };

        $http({
            method: 'GET',
            url: window.hostname + 'reports/getid/' + id
        }).then(function successCallback(response) {
            $scope.report = response.data;
            // this callback will be called asynchronously
            // when the response is available
        }, function errorCallback(response) {
            toastr.error('Erro!', response.data);
            // called asynchronously if an error occurs
            // or server returns response with an error status.
        });

        function update(){
            $http({
                method: 'GET',
                url: window.hostname + 'clients/listall'
            }).then(function successCallback(response) {
                $scope.clients=response.data;

                // this callback will be called asynchronously
                // when the response is available
            }, function errorCallback(response) {
                toastr.error('Erro!', response.data);
                // called asynchronously if an error occurs
                // or server returns response with an error status.
            });

            $http({
                method: 'GET',
                url: window.hostname + 'users/listall'
            }).then(function successCallback(response) {
                $scope.users=response.data;
                // this callback will be called asynchronously
                // when the response is available
            }, function errorCallback(response) {
                toastr.error('Erro!', response.data);
                // called asynchronously if an error occurs
                // or server returns response with an error status.
            });

        }
        update();



        $scope.submit = function () {
            $http({
                method: 'POST',
                url: window.hostname + 'reports/edit/' + id,
                data: $scope.report
            }).then(function successCallback(response) {
               
                toastr.success('Sucesso!', 'relatório Modificado');
                // this callback will be called asynchronously
                // when the response is available
            }, function errorCallback(response) {
                toastr.error('Erro!', response.data);

                // called asynchronously if an error occurs
                // or server returns response with an error status.
            });
        }
    }]);