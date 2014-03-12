'use strict';

/* App Module */

var app = angular.module('app', []);

app.controller('Ctrl', ['$scope', '$http', '$timeout',
	       function($scope, $http, $timeout) {
		       var timer;
		       $scope.channel = {};
		       $scope.programme = {};

		       $scope.channel.days = [
			       {name: "Today", value: moment().toJSON()},
			       {name: "Tomorrow", value: moment().add('d', 1).toJSON()}
		       ];

		       $scope.channel.list = [];
		       $scope.channel.map = {};
		       $scope.programme.list = [];

		       GetChannels(); 

		       $scope.programme.Fetch = function() {
			       $scope.programme.list = [];
			       if( angular.isDefined($scope.channel.selected) && angular.isDefined($scope.channel.selectedDay) ) {
				       $http.post('/api/programme',{channels: $scope.channel.selected, days: $scope.channel.selectedDay}).
					       success(function(data, status, headers, config) {
					       //console.log('programme success', data);
					       angular.forEach(data.data, function(v){
						       $scope.programme.list.push(v);
					       });
					       //console.log('programmes: ', $scope.programme.list);
				       }).error(function(data, status, headers, config) {
					       console.log('failure',data);
				       });
			       }
		       }

		       function GetChannels() {
			       $http.post('/api/channel',{channel_name: ''}).
				       success(function(data, status, headers, config) {
				       //console.log('channel success', data);
				       angular.forEach(data, function(v){
					       $scope.channel.list.push(v);
					       $scope.channel.map[v.id] = v.display_name.text;
				       });
			       }).error(function(data, status, headers, config) {
				       console.log('failure',data);
			       });
		       }

	       }]
	      );

app.filter('dateFilter', function() {
	return function(input) {
		var date = moment(input).calendar();
		//console.log('input, date:',input,date);
		return date;
	}
});

app.filter('channelFilter', function() {
	return function(input, str) {
		if( angular.isDefined(input) ) {
			var array = [];
			if ( angular.isDefined(str) ) {
				for( var i=0; i< input.length; i++) {
					if( input[i].display_name.text.toLowerCase().indexOf(str.toLowerCase()) != -1 ) {
						array.push(input[i]);
					}
				}
			} else {
				array = input;
			}
			array.sort(function(a, b){
				if( a.display_name.text.toLowerCase() < b.display_name.text.toLowerCase() ) return -1;
				if( a.display_name.text.toLowerCase() > b.display_name.text.toLowerCase() ) return 1;
				return 0;
			});
			return array;
		}
	};
});
