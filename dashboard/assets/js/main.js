$(function() {
  var SMSStatus = ["Pending", "Processed", "Error"]

  // SMS Log Table
  var logTable = $('#smsdata').dataTable({
    "data": [],
    "iDisplayLength": 5,
    "bLengthChange": false,
    "oLanguage": { "sSearch": "" },
    "columns": [
        { "data": "mobile" },
        { "data": "body" },
        { "data": "status",
          "mRender": function( data, type, full ) {
            return SMSStatus[data];
          },
          bUseRendered: false
        }
    ]
  });
  
  var loadData = function() {
    $.ajax({
      url: "/api/logs/"
    })
    .done(function(logs) {
      if(!logs.messages) {
        return
      }
      logTable.fnClearTable(logs.messages);
      logTable.fnAddData(logs.messages);
    })
    .done(function(logs) {
      // Bar Chart
      var data = []
      var daycount = logs.daycount
      for(dt in daycount) {
        var day = moment(dt, "YYYY-MM-DD").format("ddd");
        data.push([ day, daycount[dt] ])
      }
      var plot = $.plot("#barChart", [ data ], {
        series: {
          bars: {
            show: true,
            barWidth: 0.4,
            align: "center"
          }
        },
        xaxis: {
          mode: "categories",
          tickLength: 0
        }
      });
    })
    .done(function(logs) {
      // Pie Chart
      var status = []
      var summary = logs.summary;
      for(var i = 0;i < summary.length;i++) {
        status.push({ label: SMSStatus[i], data: summary[i] })
      }
      $.plot("#pieChart", status, {
        series: {
          pie: {
            radius: 1,
            innerRadius: 0.5,
            show: true,
            label: {
              radius: 3/4,
              show: true,
              formatter: labelFormatter
            }
          },
        },
        legend: {
          show: true,
        }
      });
    })
  }

  // Function to format pie chart labels
  function labelFormatter(label, series) {
    return "<div style='font-size:8pt; text-align:center; padding:2px; color: #333;'>" + Math.round(series.data[0][1]) + "</div>";
  }

  // Send Test SMS
  $("#testSMS").submit(function() {
    var url = $(this).attr('action');
    var formData = $(this).serialize();
    $.post(url, formData, function(resp) {
      // reload logs table					
      loadData();
    });
    return false;
  });
  
  loadData();

});