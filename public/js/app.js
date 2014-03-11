'use strict';

/* App Module */

var app = angular.module('app', []);

app.controller('Ctrl', ['$scope', '$http', '$timeout',
	       function($scope, $http, $timeout) {
		       var timer;
		       var channels = [];
		       var programmes = [];
		       $scope.channel = {};
		       $scope.programme = {};

		       $scope.channel.days = [
			       {name: "Today", value: moment().toJSON()},
			       {name: "Tomorrow", value: moment().add('d', 1).toJSON()}
		       ];

		       $scope.channel.selectedDay = $scope.channel.days[0];

		       $scope.programme.list = [];


		       GetChannels(); 

		       $scope.channel.Select = function() {
			       console.log("select ", $scope.channel.selected);
		       }
		       $scope.channel.SelectDay = function() {
			       programmes = [];
			       console.log("select day", $scope.channel.selectedDay);

			       $http.post('/api/programme',{channels: $scope.channel.selected, days: $scope.channel.selectedDay}).
				       success(function(data, status, headers, config) {
				       console.log('programme success');
				       //console.log(data);
				       angular.forEach(data.data, function(v){
					       	programmes.push(v);
				       });
				       $scope.programme.list = programmes;
			       }).error(function(data, status, headers, config) {
				       console.log('failure',data);
			       });
		       }

		       function GetChannels() {
			       $http.post('/api/channel',{channel_name: ''}).
				       success(function(data, status, headers, config) {
				       console.log('channel success');
				       angular.forEach(data, function(v){
					       	channels.push(v);
				       });
				       $scope.channel.list = channels;
			       }).error(function(data, status, headers, config) {
				       console.log('failure',data);
			       });
		       }

	       }]
	      );

app.filter('channelFilter', function() {
	return function(input, str) {
		if( angular.isDefined(input) ) {
			var array = [];
			if ( angular.isDefined(str) ) {
				for( var i=0; i< input.length; i++) {
					if( input[i].DisplayName.Text.toLowerCase().indexOf(str.toLowerCase()) != -1 ) {
						array.push(input[i]);
					}
				}
			} else {
				array = input;
			}
			array.sort(function(a, b){
				if( a.DisplayName.Text.toLowerCase() < b.DisplayName.Text.toLowerCase() ) return -1;
				if( a.DisplayName.Text.toLowerCase() > b.DisplayName.Text.toLowerCase() ) return 1;
				return 0;
			});
			return array;
		}
	};
});
