'use strict';

angular.module('Module', ['ngRoute'])

    .config(['$routeProvider', function ($routeProvider) {
        $routeProvider.when('/module/:moduleid/list', {
            templateUrl: 'module/list.html',
            controller: 'Listmodules'
        });
        $routeProvider.when('/module/:moduleid/create', {
            templateUrl: 'module/create.html',
            controller: 'Createmodules'
        });

        $routeProvider.when('/module/:moduleid/edit/:id', {
            templateUrl: 'module/create.html',
            controller: 'Editmodules'
        });
    }])

    .controller('Listmodules', ['$scope', '$http', 'toastr','$routeParams', function ($scope, $http, toastr,$routeParams) {
       $scope.modules = [];
        $scope.schema={};

        var moduleid = $routeParams.moduleid;

        $scope.search=function (keywords){
            if (keywords==""){
                update();
                return

            }

            if (keywords.indexOf(suffix, keywords.length - suffix.length) !== -1){
                return
            }


            $http({
                method: 'GET',
                url: window.hostname + 'modules/search/'+moduleid+'/'+keywords
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

        var update = function () {
            $http({
                method: 'GET',
                url: window.hostname + 'modules/getmoduleschema/'+moduleid
            }).then(function successCallback(response) {
                $scope.schema=response.data;


                // this callback will be called asynchronously
                // when the response is available
            }, function errorCallback(response) {
                toastr.error('Erro!', response.data);
                // called asynchronously if an error occurs
                // or server returns response with an error status.
            });
            $http({
                method: 'GET',
                url: window.hostname + 'modules/listall/'+moduleid
            }).then(function successCallback(response) {
                console.log(response);
               $scope.modules = response.data;
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
                url: window.hostname + 'modules/delete/'+moduleid+'/' + id,
                data: $scope.module
            }).then(function successCallback(response) {
                update();
                toastr.success('Sucesso!', 'Eliminado');
                // this callback will be called asynchronously
                // when the response is available
            }, function errorCallback(response) {
                toastr.error('Erro!', response.data);
                // called asynchronously if an error occurs
                // or server returns response with an error status.
            });

        }


    }])
    .controller('Createmodules', ['$scope', '$http', 'toastr','$routeParams', function ($scope, $http, toastr,$routeParams) {

        $scope.update=false;

        var moduleid = $routeParams.moduleid;
        $scope.module = {};
        $scope.clients = [];
        $scope.users= [];
        $scope.schema={};
        $scope.module.Attachments=[];
 

        $scope.addAttach=function(){
            $scope.module.Attachments.push({})
        };

        $scope.deleteAttach=function(i){
            $scope.module.Attachments.splice(i, 1);
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
                url: window.hostname + 'modules/getmoduleschema/'+moduleid
            }).then(function successCallback(response) {
                $scope.schema=response.data;
                $scope.module.ExtraFields=$scope.schema.ExtraFieldsHeader

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
                url: window.hostname + 'modules/create/'+moduleid,
                data: $scope.module
            }).then(function successCallback(response) {
                console.log(response)
                toastr.success('Sucesso!', 'Adicionado');
                // this callback will be called asynchronously
                // when the response is available
            }, function errorCallback(response) {
                toastr.error('Erro!', response.data);
                // called asynchronously if an error occurs
                // or server returns response with an error status.
            });
        }

    }])
    .controller('Editmodules', ['$scope', '$http', '$routeParams', 'toastr', function ($scope, $http, $routeParams, toastr) {
        $scope.update=true;
        var moduleid = $routeParams.moduleid;
        $scope.module = {};
        $scope.clients = [];
        $scope.users= [];
        $scope.schema={};
        $scope.module.Attachments=[];
        var id = $routeParams.id;


        $scope.addAttach=function(){
            $scope.module.Attachments.push({})
        };

        $scope.deleteAttach=function(i){
            $scope.module.Attachments.splice(i, 1);
        };

        $http({
            method: 'GET',
            url: window.hostname + 'modules/getid/'+moduleid+'/'+ id
        }).then(function successCallback(response) {
            $scope.module = response.data;
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
                url: window.hostname + 'modules/getmoduleschema/'+moduleid
            }).then(function successCallback(response) {
                $scope.schema=response.data;
                if  ($scope.module.ExtraFields.length < 1){
                    $scope.module.ExtraFields=$scope.schema.ExtraFieldsHeader
                }


                // this callback will be called asynchronously
                // when the response is available
            }, function errorCallback(response) {
                toastr.error('Erro!', response.data);
                // called asynchronously if an error occurs
                // or server returns response with an error status.
            });



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
                url: window.hostname + 'modules/edit/'+moduleid+'/'+ id,
                data: $scope.module
            }).then(function successCallback(response) {
               
                toastr.success('Sucesso!', 'Modificado');
                // this callback will be called asynchronously
                // when the response is available
            }, function errorCallback(response) {
                toastr.error('Erro!', response.data);

                // called asynchronously if an error occurs
                // or server returns response with an error status.
            });
        }
    }]);