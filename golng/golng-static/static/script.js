var last = 0;
var items = [];
var position = "left";

$(document).ready(function() {
    RequestData();

    $("#load").click(RequestData);
});

Element.prototype.remove = function() {
    this.parentElement.removeChild(this);
}

NodeList.prototype.remove = HTMLCollection.prototype.remove = function() {
    for (var i = this.length - 1; i >= 0; i--) {
        if (this[i] && this[i].parentElement) {
            this[i].parentElement.removeChild(this[i]);
        }
    }
}

function timeConverter(UNIX_timestamp) {
    var a = new Date(UNIX_timestamp * 1000);
    var months = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];
    var year = a.getFullYear();
    var month = months[a.getMonth()];
    var date = a.getDate();
    var time = date + ' ' + month + ' ' + year;
    return time;
}

String.prototype.replaceAll = function(search, replacement) {
    var target = this;
    return target.split(search).join(replacement);
};

function RequestData() {
    var url = "https://golng.ml/api/" + last;
    jQuery.ajax({
        url: url,
        type: "GET",
        crossDomain: true,
        beforeSend: function() {
            $("#load").hide();
            $("#preloader").show();
        },
        success: function(data) {
            data = JSON.parse(data);
            if (data.errorcode == 0) {
                var result = data.posts;
                for (var x = 0; x < result.length; x++) {
                    if (items.indexOf(result[x].id) > 0) continue;
                    items.push(result[x].id);

                    result[x].content = result[x].content.replaceAll("\n","<br/>");
                    var div = document.createElement('div');
                    div.setAttribute("class", "timeline-item");
                    var icon = document.createElement('div');
                    icon.setAttribute("class", "timeline-icon");
                    div.appendChild(icon);
                    var content = document.createElement('div');
                    content.setAttribute("class", "timeline-content " + position)
                    content.innerHTML = "<h2>The Daily Golang</h2>" + "<p>" + result[x].content + "</p><hr/><h5>" + timeConverter(result[x].time) + "</h5>";
                    content.innerHTML += "<a href='/"+ result[x].id + "'> Link</a>";
                    div.appendChild(content);
                    document.getElementById("timeline").appendChild(div);
                    if (position == "right") position = "left";
                    else position = "right";
                    last = result[x].time;
                }
            } else {
                alert("error");
            }
            $("#load").show();
            $("#preloader").hide();
        },
        error: function(XMLHttpRequest, textStatus, errorThrown) {
            alert("error");
            $("#load").hide();
            $("#preloader").hide();
        }
    });
}
