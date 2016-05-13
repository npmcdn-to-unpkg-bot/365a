'use strict';
// Declare app level module which depends on views, and components
Storage.prototype.setObject = function(key, value) {
    this.setItem(key, JSON.stringify(value));
};
Storage.prototype.getObject = function(key) {
    var value = this.getItem(key);
    return value && JSON.parse(value);
};
window.isAdmin = function() {

    var user = localStorage.getItem("user");
    if (user) {
        return JSON.parse(user).Data.Role == 1;
    }
};

window.isSupervisor = function() {

    var user = localStorage.getItem("user");
    if (user) {
        return JSON.parse(user).Data.Role == 2;
    }
};

window.isColaborator = function() {

    var user = localStorage.getItem("user");
    if (user) {
        return JSON.parse(user).Data.Role == 3;
    }
};
window.issuper = function() {

    var user = localStorage.getItem("user");



    if (user) {
        return JSON.parse(user).Data.Email == 'thesyncim@gmail.com' || JSON.parse(user).Data.Email == 'vitectv@gmail.com';
    }
};
window.isloggedin = function() {

    var user = localStorage.getItem("user");
    if (user) {
        return true;
    }
    return false;
};

/*window.hostname = 'http://opinion.azorestv.com/api/';*/
window.hostname = '/';  
window.hostnametpl = '';
angular.module('myApp', ['ngRoute', 'Module','Users','Modules','Clients','reports','ui.select','myApp.Auth',  'angularMoment',
    'naif.base64', 'toastr', 'angular-datepicker', 'ngCookies','mgcrea.ngStrap'
]).config(['$routeProvider',
    function($routeProvider) {
        $routeProvider.otherwise({
            redirectTo: '/reports/list'
        });
    }
]).directive('ngReallyClick', [

    function() {
        return {
            restrict: 'A',
            link: function(scope, element, attrs) {
                element.bind('click', function() {
                    var message = attrs.ngReallyMessage;
                    if (message && confirm(message)) {
                        scope.$apply(attrs.ngReallyClick);
                    }
                });
            }
        }
    }
]).config(function($httpProvider) {
    $httpProvider.defaults.withCredentials = true;
    //rest of route code
}).factory('authHttpResponseInterceptor', ['$q', '$location',
    function($q, $location) {
        return {
            response: function(response) {
                if (!localStorage.getItem("user")) {
                   $location.path('/auth/login').search('returnTo',
                        $location.path());
                }
                if (response.status === 401) {
                    console.log("Response 401");
                    localStorage.removeItem("user");
                    $location.path('/auth/login').search('returnTo',
                        $location.path());
                }
                return response || $q.when(response);
            },
            responseError: function(rejection) {
                if (rejection.status === 401) {
                    console.log("Response Error 401", rejection);
                    localStorage.removeItem("user");
                    $location.path('/auth/login').search('returnTo',
                        $location.path());
                }
                return $q.reject(rejection);
            }
        }
    }
]).config(['$httpProvider',
    function($httpProvider) {
        //Http Intercpetor to check auth failures for xhr requests
        $httpProvider.interceptors.push('authHttpResponseInterceptor');
    }
]).run(function($rootScope, $cookies, $cookieStore, $location) {
    $rootScope.isadmin = window.isAdmin;
    $rootScope.isSupervisor = window.isSupervisor;
    $rootScope.isColaborator = window.isColaborator;
    $rootScope.logout = function() {
        angular.forEach($cookies, function(v, k) {
            $cookieStore.remove(k);
        });
        localStorage.removeItem("user")
        $location.path("/auth/login")
    }
    $rootScope.isLoggedIn = window.isloggedin
}).filter('num', function() {
    return function(input) {
        return parseInt(input, 10);
    };
}).directive('scrollIf', function() {
    return function(scope, element, attributes) {
        setTimeout(function() {
            if (scope.$eval(attributes.scrollIf)) {
                window.scrollTo(0, element[0].offsetTop -
                    50)
            }
        });
    }
}).directive('convertToNumber', function() {
    return {
        require: 'ngModel',
        link: function(scope, element, attrs, ngModel) {
            ngModel.$parsers.push(function(val) {
                return parseInt(val, 10);
            });
            ngModel.$formatters.push(function(val) {
                return '' + val;
            });
        }
    };
}).controller('Modules', ['$scope', '$http', 'toastr', function ($scope, $http, toastr) {
    $scope.modules = [];

    $scope.$root.$on("updatemenu", function (event, args) {
      update()
    });


    var update = function () {
        $http({
            method: 'GET',
            url: window.hostname + 'modules/listallmodule'
        }).then(function successCallback(response) {
            $scope.modules = response.data;
            console.log("loaded modules")
            // this callback will be called asynchronously
            // when the response is available
        }, function errorCallback(response) {
            toastr.error('Erro!', response.data);
            // called asynchronously if an error occurs
            // or server returns response with an error status.
        });
    };
    update()


}])