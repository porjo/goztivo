'use strict';

/* App Module */

var app = angular.module('app', []);

app.controller('Ctrl', ['$scope', '$http', '$timeout',
	       function($scope, $http, $timeout) {
		       var timer;
		       var channels = [];
		       $scope.channel = {};

		       GetChannels(); 

		       $scope.channel.Select = function(selected) {
			       console.log("select ", selected);
		       }

		       function GetChannels() {
			       $http.post('/api/channel',{channel_name: ''}).
				       success(function(data, status, headers, config) {
				       console.log('success');
				       angular.forEach(data, function(v){
					       	channels.push(v);
				       });
				       $scope.channel.list = channels;
				       //console.log($scope.channel.list);
				       //console.log(data);
				       console.log(channels);
			       }).
				       error(function(data, status, headers, config) {
				       console.log('failure',data);
			       });
		       }

	       }]
	      );

