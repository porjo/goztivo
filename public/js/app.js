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
					       console.log('programme success');
					       //console.log(data);
					       angular.forEach(data.data, function(v){
						       $scope.programme.list.push(v);
					       });
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
					       $scope.channel.map[v.Id] = v.DisplayName.Text;
				       });
			       }).error(function(data, status, headers, config) {
				       console.log('failure',data);
			       });
		       }

	       }]
	      );

	      /*
	      app.filter('channelFilter', function() {
		      return function(items, str) {
			      if( !angular.isDefined(items) ) {
				      return
			      }
			      if ( !angular.isDefined(str) ) {
				      return items
			      }
			      var result = {};
			      var sortable = [];
			      angular.forEach(items, function(value, key) {
				      if( value.DisplayName.Text.toLowerCase().indexOf(str.toLowerCase()) != -1 ) {
					      sortable.push(value);
				      }
			      });
			      //console.log("before: ",sortable);
			      sortable.sort(function(a, b) {
				      if( a.DisplayName.Text.toLowerCase() < b.DisplayName.Text.toLowerCase() ) return -1;
				      if( a.DisplayName.Text.toLowerCase() > b.DisplayName.Text.toLowerCase() ) return 1;
				      return 0;
			      });
			      //console.log("after: ",sortable);
			      for(var i in sortable){
				      var v = sortable[i];
				      result[v.Id] = v;
			      }
			      return result;
		      };
	      });
	     */

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
