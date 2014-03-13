'use strict';

/* App Module */

var app = angular.module('app', []);

app.controller('Ctrl', ['$scope', '$http', '$timeout',
	       function($scope, $http, $timeout) {
		       var timer;
		       $scope.channel = {};
		       $scope.programme = {};

		       var m = moment();
		       $scope.channel.days = [
			       {text: "Today", value: m.toJSON()},
			       {text: "Tomorrow", value: m.add('d', 1).toJSON()}
		       ];

		       BuildHourList();

		       $scope.channel.list = [];
		       $scope.channel.selectedChannel = [];
		       $scope.channel.selectedDay = [];
		       $scope.channel.map = {};
		       $scope.programme.list = [];

		       GetChannels(); 

		       $scope.programme.Fetch = function() {
			       if($scope.channel.selectedChannel.length == 0 || $scope.channel.selectedDay.length == 0) {
				       return;
			       }
			       if( angular.isDefined($scope.channel.selectedChannel) && angular.isDefined($scope.channel.selectedDay) ) {
				       $http.post('/api/programme',{channels: $scope.channel.selectedChannel, days: $scope.channel.selectedDay}).
					       success(function(data, status, headers, config) {
					       $scope.programme.list = [];
					       //console.log('programme success', data);
					       angular.forEach(data.data, function(v){
						       $scope.programme.list.push(v);
					       });
					       console.log('programmes: ', $scope.programme.list);
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

		       function BuildHourList() {
			       var mdt = moment().hour(0).minute(0);
			       $scope.channel.hours = [];
			       for( var i=0; i<24; i++) {
				       var tmp = mdt;
				       var text = mdt.format("HH:mm");
				       var hour = {}
				       hour.value = mdt.toJSON()
				       mdt.add('h',1).toJSON();
				       text += "-" + mdt.format("HH:mm");
				       hour.text = text;
				       $scope.channel.hours.push(hour);
				       if(tmp.hour() == moment().hour()){
					       $scope.channel.selectedHour = hour;
				       }
			       }
		       }

	       }]
	      );

app.filter('dateFilter', function() {
	return function(input,format) {
		console.log("date:",input,format);
		var date = moment(input).format(format);
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

app.filter('truncateFilter', function() {
	return function(input, count) {
		if( angular.isDefined(input) ) {
			if( input.length > 100 ) {
				return input.substring(0,100) + "...";
			}
			return input;
		}
	}
});

app.filter('programmeFilter', function() {
	return function(input, hours) {
		if( !angular.isDefined(input) ) {
			return 
		}
		if( angular.isDefined(hours) ) {
			if(hours.length == 0) {	return input; }
			var pmap = {};
			for( var i in input) {
				var pm = moment(input[i].start_time);
				for(var j=0; j<hours.length; j++) {
					var hm = moment(hours[j]);
					if(pm.hour() == hm.hour()) {
						pmap[input[i].title] = input[i];
					}
				}
			}
			var programmes = [];
			for( var p in pmap) {
				programmes.push(pmap[p]);
			}
			return programmes;
		}
	}
});
