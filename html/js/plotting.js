function plotter(container, placeholder, data) {
	this.container = container;
	this.placeholder = placeholder;

	this.data = null;
	this.postObj = new Object();
	this.plot = null;
	this.updateInterval = 60000;

	this.go = function() {

		var livePlot = (this.end === null || this.end === undefined);

		this.plot = $.plot(this.container, [ this.getData() ], {
			series: {
				shadowSize: 0,	// Drawing is faster without shadows
				color: "rgb(1,169,220)",
				lines: { show: true },
			},
			yaxis: {
				min: 50,
				max: 80,
				tickFormatter: function(tickVal, axis) { return tickVal + "&deg;F"; }
			},
			xaxis: {
				mode: "time",
				timezone: "browser",
				ticks: 5,
				twelveHourClock: true
			}
		});

		var oldHeight = parseInt($(this.container).css("height").replace("px", ""));
		var heightAddition = 0;//parseInt($(container+" .plot-title").css("height").replace("px", ""));
		var newHeight = 10 + oldHeight + heightAddition + "px";

		$(this.placeholder).css("height", newHeight);
		
		if(livePlot)
			this.update();
	}

	this.getData = function() {
		
		if(!this.postObj.after) {

			this.postObj.after = this.start;
			if(this.end) {
				this.postObj.before = this.end;
			}

			$.ajax({
  			url: "gettemperature.php",
  			type: "POST",
  			dataType: "json",
  			async: false,
  			data: this.postObj,
  			context: this
			}).done(function(jsonData) {
				
				this.data = new Array();

				if(jsonData.length <= 1) {
					if(this.statstable)
						$(this.statstable).hide();
					return;
				}

				if(this.statstable) {
					var stats = jsonData[0];
					$(this.statstable).show();
					$(this.statstable+' .min-temp').html(parseFloat(stats.min).toFixed(3));
					$(this.statstable+' .max-temp').html(parseFloat(stats.max).toFixed(3));
					$(this.statstable+' .avg-temp').html(parseFloat(stats.avg).toFixed(3));
				}

				jsonData = jsonData.slice(1);

				this.postObj.after = jsonData[0].time;
				
				jsonData.reverse();

				if(jsonData.length <= 25000 ||
					(navigator.userAgent.indexOf("Safari") < 0 || navigator.userAgent.indexOf("Chrome") >= 0)) {
					var dontlosescope = this;

					$.each(jsonData, function(i, item) {
	  					dontlosescope.data.push([fixDate(item.time), item.temp]);
						});
				}
			});
		}
		else {
			$.ajax({
  			url: "gettemperature.php",
  			type: "POST",
  			dataType: "json",
  			async: false,
  			data: this.postObj,
  			context: this
			}).done(function(jsonData) {
				
				if(jsonData.length <= 1) return;

				jsonData = jsonData.slice(1);

				this.postObj.after = jsonData[0].time;
				
				jsonData.reverse();

				var offset = this.data[this.data.length-1];
				if((this.data + jsonData.length) <= 25000 ||
					(navigator.userAgent.indexOf("Safari") < 0 || navigator.userAgent.indexOf("Chrome") >= 0)) {
					var dontlosescope = this;

					$.each(jsonData, function(i, item) {
	  				dontlosescope.data.push([fixDate(item.time), item.temp]);
					});
				}
				
				if(this.statstable) {
					$(this.statstable).show();
					var min = 10000;
					var max = 0;
					var avg = 0;

					$.each(this.data, function(i, item) {
						var temp = parseFloat(item[1]);
						if(temp < min) min = temp;
						if(temp > max) max = temp;
						avg += temp;
					});

					avg /= this.data.length;

					$(this.statstable+' .min-temp').html(parseFloat(min).toFixed(3));
					$(this.statstable+' .max-temp').html(parseFloat(max).toFixed(3));
					$(this.statstable+' .avg-temp').html(parseFloat(avg).toFixed(3));
				}

			});
		}

		return this.data;
	}

	this.update = function(useThis) {

		if(!useThis)
			useThis = this;

		useThis.plot.setData([ useThis.getData() ]);

		useThis.plot.setupGrid();

		useThis.plot.draw();
			
		setTimeout(function() { useThis.update(useThis); }, useThis.updateInterval);
	}
}