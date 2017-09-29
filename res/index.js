function setIpAllow() {
    minAjax({
        url: "/ipAllow",
        type: "POST",
        data: {
            envs: getCheckedValues('env'),
            officeIp: $('myip').innerText
        },
        success: function (redirectUrl) {
            window.location = redirectUrl
        }
    })
}

function getCheckedValues(checkboxClass) {
    var checkedValue = []
    var inputElements = document.getElementsByClassName(checkboxClass)
    for (var i = 0; inputElements[i]; ++i) {
        if (inputElements[i].checked) {
            checkedValue.push(inputElements[i].value)
        }
    }
    return checkedValue.join(',')
}

function $(id) {
    return document.getElementById(id)
}

/*|--(A Minimalistic Pure JavaScript Header for Ajax POST/GET Request )--|
  |--Author : flouthoc (gunnerar7@gmail.com)(http://github.com/flouthoc)--|
*/

function initXMLhttp() {
    if (window.XMLHttpRequest) { // code for IE7,firefox chrome and above
        return new XMLHttpRequest()
    } else { // code for Internet Explorer
        return new ActiveXObject("Microsoft.XMLHTTP")
    }
}

function minAjax(config) {
    /*
        Config Structure
        url:"reqesting URL"
        type:"GET or POST"
        async: "(OPTIONAL) True for async and False for Non-async | By default its Async"
        data: "(OPTIONAL) another Nested Object which should contains reqested Properties in form of Object Properties"
        success: "(OPTIONAL) Callback function to process after response | function(data,status)"
    */

    config.async = config.async || true

    var xmlhttp = initXMLhttp()
    xmlhttp.onreadystatechange = function () {
        if (xmlhttp.readyState == 4 && xmlhttp.status == 200) {
            config.success(xmlhttp.responseText, xmlhttp.readyState)
        }
    }

    var sendString = [], sendData = config.data
    if (typeof sendData === "string") {
        var tmpArr = String.prototype.split.call(sendData, '&')
        for (var i = 0, j = tmpArr.length; i < j; i++) {
            var datum = tmpArr[i].split('=')
            sendString.push(encodeURIComponent(datum[0]) + "=" + encodeURIComponent(datum[1]))
        }
    } else if (typeof sendData === 'object' && !( sendData instanceof String || (FormData && sendData instanceof FormData) )) {
        for (var k in sendData) {
            var datum = sendData[k]
            if (Object.prototype.toString.call(datum) == "[object Array]") {
                for (var i = 0, j = datum.length; i < j; i++) {
                    sendString.push(encodeURIComponent(k) + "[]=" + encodeURIComponent(datum[i]))
                }
            } else {
                sendString.push(encodeURIComponent(k) + "=" + encodeURIComponent(datum))
            }
        }
    }
    sendString = sendString.join('&')

    if (config.type == "GET") {
        xmlhttp.open("GET", config.url + "?" + sendString, config.async)
        xmlhttp.send()
    } else if (config.type == "POST") {
        xmlhttp.open("POST", config.url, config.async)
        xmlhttp.setRequestHeader("Content-type", "application/x-www-form-urlencoded")
        xmlhttp.send(sendString)
    }
}

minAjax({
    url: "http://icanhazip.com",
    type: "GET",
    data: {},
    success: function (data) {
        $('myip').innerText = data.trim()
    }
})

/*.ALERTS*/