'use strict';

/* App Module */

var app = angular.module('app', []);

app.controller('Ctrl', ['$scope', '$http', '$timeout',
	       function($scope, $http, $timeout) {
		       var timer;
		       $scope.channel = {};

		       Update(''); 

		       $scope.channel.Select = function(selected) {
			       console.log("select ", selected);
		       }

		       $scope.channel.Update = function(query) {
			       console.log("search ", query);
			       // Delay submit
			       if(timer){
				       $timeout.cancel(timer)
			       }  
			       timer = $timeout(function(){
				       Update(query); 
			       },500)
		       }


		       function Update(query) {
			       $http.post('/api/channel',{channel_name: query, contains: true}).
				       success(function(data, status, headers, config) {
				       console.log('success');
				       $scope.channel.list = data;
			       }).
				       error(function(data, status, headers, config) {
				       console.log('failure',data);
			       });
		       }

	       }]
	      );

