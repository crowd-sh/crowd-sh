if(window.location.href.match("localhost")) {
    WorkMachine = {
        //server: "http://localhost:5000",
        server: "http://britishlibraryapi.workmachine.us:5000",
    }
} else if(window.location.href.match("workmachine.us")) {
    WorkMachine = {
        server: "http://britishlibraryapi.workmachine.us:5000",
    }
}

angular.module('workmachine', ['ngRoute', 'ngTouch'])

    .config(function($routeProvider) {
        $routeProvider
            .when('/', {
                controller:'WorkCtrl',
                templateUrl:'templates/work.html'
            })
            .otherwise({
                redirectTo:'/'
            });
    })

    .factory('WorkMachineService', function($http) {
        return {
            getAssignment: function() {
                return $http.get(WorkMachine.server + "/v1/assignment")
                    .then(function(response) {
                        console.log("Assignment: " + response.data);
                        return response.data;
                    });
            },

            postAssignment: function(data) {
                return $http.post(
                    WorkMachine.server + "/v1/assignment",
                    $.param(data),
                    { headers: { 'Content-Type': 'application/x-www-form-urlencoded' } }
                    ).then(function(response) {
                        console.log("Assignment: " + response.data);
                        return response.data;
                    });
            },
        }
    })

    .controller('WorkCtrl', function($scope, $sce, WorkMachineService) {
        $scope.work = {};

        WorkMachineService.getAssignment().then(function(data) {
            $scope.assignment = data;
            $scope.assignment.job.info.description = $sce.trustAsHtml(data.job.info.description);
        });

        $scope.submitCheckbox = function(value) {
            $scope.work[$scope.assignment.job.output.Id] = value;

            $scope.postAssignment();
        };

        $scope.postAssignment = function() {
            $scope.work['id'] = $scope.assignment.id;
            WorkMachineService.postAssignment($scope.work);

            if(typeof Android != 'undefined') {
                Android.incrementTask();
            }

            window.location = "#/new";

            return true;
        };
    })

    .controller('LeaderBoardCtrl', function($scope, WorkMachineService) {
        $scope.leaderboard = [];

        if(typeof Android != 'undefined') {
            $scope.leaderboard = Android.getLeaderboard();
        }

    });
