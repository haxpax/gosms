$(function() {
  var SMSStatus = ["Processed", "Pending", "Error"]
  $.ajax({
    url: "/api/logs/"       
  })
  .done(function(logs) {
    // SMS Log Table    
    var logTable = $('#smsdata').dataTable({
      "data": logs.messages,
      "bLengthChange": false,
      "oLanguage": { "sSearch": "" },
      "columns": [
          { "data": "mobile" },
          { "data": "body" },
          { "data": "status",
            "mRender": function( data, type, full ) {
              if(data == 0)
                return "Pending";
              else if(data == 1)
                return "Processed";
              return "Error";
            },
            bUseRendered: false
          }
      ]
    });
  })  
  .done(function(logs) {
    // Bar Chart
    var data = []
    var daycount = logs.daycount    
    for(day in daycount) {
      data.push([ day, daycount[day] ])
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
          show: true
        },
      },
      legend: {
        show: true,
      }
    });        
  })

  function labelFormatter(label, series) {
    return "<div style='font-size:8pt; text-align:center; padding:2px; color:white;'>" + label + "<br/>"
    + Math.round(series.percent) + "%</div>";
  }

  // Send Test SMS
  $("#testSMS").submit(function() {
    var url = $(this).attr('action');
    var formData = $(this).serialize();
    $.post(url, formData, function(resp) {
      // reload logs table					
      logTable.api().ajax.reload();
    });
    return false;
  });
			
});