function setIpAllow() {
    if (!$('myip').value.match(ipRegex)) {
        alert('请在输入正确的IP后，再设置！')
        $('myip').focus()
        return
    }
    minAjax({
        url: "/ipAllow",
        type: "POST",
        data: {
            envs: getCheckedValues('env'),
            officeIp: $('myip').value
        },
        success: function (redirectUrl) {
            if (redirectUrl.indexOf('http') == 0) {
                window.location = redirectUrl
            } else {
                alert(redirectUrl)
            }
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
    } else if (typeof sendData === 'object' && !(sendData instanceof String || (FormData && sendData instanceof FormData))) {
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

//this function will work cross-browser for loading scripts asynchronously
function loadScript(src, callback) {
    var r = false
    var s = document.createElement('script')
    s.async = true
    s.type = 'text/javascript'
    s.src = src
    s.onload = s.onreadystatechange = function () {
        // console.log( this.readyState ); // uncomment this line to see which ready states are called.
        if (!r && (!this.readyState || this.readyState == 'complete')) {
            r = true
            if (callback) callback()
        }
    }
    var t = document.getElementsByTagName('script')[0]
    t.parentNode.insertBefore(s, t)
}

var ipRegex = /^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$/

function setIp(json) {
    console.log('开始设置IP:' + JSON.stringify(json))
    var ip = json.ip || json.query
    if ($('myip').value.match(ipRegex)) {
        console.log('已经设置，跳过')
        return
    }

    if (ip && ip.length > 0) {
        $('myip').value = ip
    } else {
        $('myip').value = '识别失败，请手工输入IP或者拷贝下面的IP后设置。'
    }
}

loadScript("http://pv.sohu.com/cityjson/getip.aspx", function () {
    setIp({ip: returnCitySN.cip})
})

// {"ip":"180.111.235.50"}
loadScript("https://api.ipify.org?format=jsonp&callback=setIp")

// {"as":"AS4134 No.31,Jin-rong Street","city":"Nanjing","country":"China","countryCode":"CN","isp":"China Telecom jiangsu",
// "lat":32.0617,"lon":118.7778,"org":"China Telecom jiangsu","query":"180.111.235.50","region":"32","regionName":"Jiangsu",
// "status":"success","timezone":"Asia/Shanghai","zip":""}
loadScript("http://ip-api.com/json/?callback=setIp")

/*.ALERTS*/