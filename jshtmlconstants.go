package main

const (
	Header = `<!DOCTYPE html PUBLIC "-//W3C//DTD HTML 4.01//EN" "http://www.w3.org/TR/html4/strict.dtd">
<html>
<head>
	<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
	<title>RUPS</title>

	<link href="css/examples.css" rel="stylesheet" type="text/css">

	<!--[if lte IE 8]><script language="javascript" type="text/javascript" src="js/excanvas.min.js"></script><![endif]-->
	<script language="javascript" type="text/javascript" src="js/moment.min.js"></script>
	<script language="javascript" type="text/javascript" src="js/jquery.js"></script>
	<script language="javascript" type="text/javascript" src="js/jquery.flot.js"></script>
	<script language="javascript" type="text/javascript" src="js/jquery.flot.time.min.js"></script>
	<script type="text/javascript">

	$(function () {
`

	PlotClose = `		], {
			series: {
				lines: { show: true },
			},
			xaxis: {
				mode: "time",
				timezone: "browser",
				ticks: 5,
				twelveHourClock: true
			},
			yaxis: {
				min: -1,
				max: 2.,
			},
			grid: {
				backgroundColor: { colors: [ "#fff", "#eee" ] },
				borderWidth: {
					top: 1,
					right: 1,
					bottom: 2,
					left: 2
				}
			},
			legend: {
				position: "sw"
			}
		});
`

	Middle = `
		// Add the Flot version string to the footer

		$("#footer").append("Flot " + $.plot.version);
	});

	</script>
</head>
<body>
	<div id="content">
`

	Footer = `
	</div>

	<div id="footer">Powered by Go and </div>

</body>
</html>
`
)