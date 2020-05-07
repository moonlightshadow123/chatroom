var msgDiv_str          = '<div class="row"/>';     

var usrColDiv_str       = '<div class="col-1"/>';
var usernameSpan_str    = '<span style="margin-right: 15px;font-weight: 700;overflow: hidden;text-align: right;font-size:20px;"/>';

var msgColDiv_str       = '<div class="col-9"/>';
var msgbodySpan_str     = '<p style="color: gray;font-size:20px; display: block;word-wrap: break-word;"/>';

var cltColDiv_str       = '<div class="col-10" style="text-align:right"/>'
var cltinfoSpan_str     = '<span style="color:#999999;font-size:small">';

var dateColDiv_str      = '<div class="col-2" style="text-align:right"/>';
var dateSpan_str        = '<span style="color: gray;font-size:small"/>';

    // 通过一个hash函数得到用户名的颜色
function getUsernameColor (username) {
    var COLORS = [
        '#e21400', '#91580f', '#f8a700', '#f78b00',
        '#58dc00', '#287b00', '#a8f07a', '#4ae8c4',
        '#3b88eb', '#3824aa', '#a700ff', '#d300e7'
    ];
    // Compute hash code
    var hash = 7;
    for (var i = 0; i < username.length; i++) {
        hash = username.charCodeAt(i) + (hash << 5) - hash;
    }
    // Calculate color
    var index = Math.abs(hash % COLORS.length);
    return COLORS[index];
}

function genMsgDiv(data){
    var name = data.name;
    var type = data.type;
    var msg = data.message;
    var date = data.date;
    var id = data.id;

    // type为0表示有人发消息
    var msgDiv;
    if (type == 0) {
        // UsrColDiv
        var $usernameSpan = $(usernameSpan_str).text(name);
        $usernameSpan.css('color', getUsernameColor(name));
        var $usrColDiv = $(usrColDiv_str).append($usernameSpan);
        // MsgColDiv
        var $msgColDiv = $(msgColDiv_str).append($(msgbodySpan_str).text(msg));
        // DateColDiv
        var $dateColDiv = $(dateColDiv_str).append($(dateSpan_str).text(date));
        // MsgDiv
        $msgDiv = $(msgDiv_str).addClass("msgDiv").attr("data-id", id.toString()).append($usrColDiv, $msgColDiv, $dateColDiv);
    }
    // type为1或2表示有人加入或退出
    else {
        // cltColDiv
        var $cltColDiv = $(cltColDiv_str).append($(cltinfoSpan_str).text(msg));
        // DateColDiv
        var $dateColDiv = $(dateColDiv_str).append($(dateSpan_str).text(date));
        // MsgDiv
        $msgDiv = $(msgDiv_str).attr("id", id.toString()).append($cltColDiv, $dateColDiv);
    }
    return $msgDiv;
}

var iconStr = '<i data-id="good" style="color: #999999; float:left" class="fa fa-circle" >';
var uSpanStr = '<span/>';
var cltDivStr = '<div class="cltDiv" style="text-align:center"></div>';
var liStr = '<li></li>';

function genLi(){
    var map = umap;
    $ul = $("#cltUl")
    for(var key in map){
        // icon
        $icon = $(iconStr).attr("data-id", key);
        if(map[key]==true) {$icon.css("color", "#339533");}
        // name
        $userSpan = $(uSpanStr).text(key);//$(usernameSpan_str).text(key);
        //$userSpan.css('color', getUsernameColor(key));
        // li
        $li = $(liStr).append($(cltDivStr).append($icon, $userSpan));
        $ul.append($li);
    }
}

var lasttimeDive_str = '<div class="row"></div>';

function getHistroyList(num){
    $.getJSON(window.location.origin + "/hist/" + window.stamp + "/" + num.toString(), function(data){
        if(data["stamp"] != "") window.stamp = data["stamp"];
        msglist = data["msglist"];
        msglist.forEach(function(msg){ 
            $msgDiv = genMsgDiv(msg)
            $("#msgContainer").children(':eq(0)').after($msgDiv);
        });
    });
}


$(function() {

    //=========================初始化====================================
    var $window = $(window);
    var $msgContainer = $('#msgContainer');       // 消息显示的区域
    var $inputArea = $('#inputArea');           // 输入消息的区域
    $inputArea.focus();  // 首先聚焦到输入框
    var connected = false;                      // 用来判断是否连接的标志

    console.log(umap);

    // 初始化显示的名字颜色
    var name = $("#name").text();
    var nameColor = getUsernameColor(name);
    $("#name").css('color', nameColor);
    genLi();
    $("#histLink").click(function(){
        //alert("Hist!");
        getHistroyList(10);
    });


    //====================webSocket连接======================
    // 创建一个webSocket连接
    var socket = new WebSocket('ws://'+window.location.host+'/chatRoom/WS?name=' + $('#name').text());

    // 当webSocket连接成功的回调函数
    socket.onopen = function () {
        console.log("webSocket open");
        connected = true;
    };

    // 断开webSocket连接的回调函数
    socket.onclose = function () {
        console.log("webSocket close");
        connected = false;
    };

    //=======================接收消息并显示===========================
    // 接受webSocket连接中，来自服务端的消息
    socket.onmessage = function(event) {
        // 将服务端发送来的消息进行json解析
        var data = JSON.parse(event.data);
        console.log("revice:" , data);

        $msgDiv = genMsgDiv(data)
        if(data["type"] == 1){ //join
            activateIcon(data["name"]);
        }else if(data["type"] == 2){ //leave
            deactivateIcon(data["name"]);
        }else if(data["type"] == 3){ //recall
            var id = data["related"];
            $("div.msgDiv[data-id='" + id.toString() + "']").remove();
        }

        $msgContainer.append($msgDiv);
        $msgContainer[0].scrollTop = $msgContainer[0].scrollHeight;   // 让屏幕滚动
        
    }

    //========================发送消息==========================
    // 当回车键敲下
    $window.keydown(function (event) {
        // 13是回车的键位
        if (event.which === 13) {
            sendMessage();
            typing = false;
        }
    });

    // 发送按钮点击事件
    $("#sendBtn").click(function () {
        sendMessage();
    });

    // suppress the right-click menu
    $($msgContainer).on('contextmenu', 'div.msgDiv', function(evt) {
        evt.preventDefault();
    });

    $($msgContainer).on('mouseup', 'div.msgDiv', function(evt) {
      if (evt.which === 3) { // right-click
        /* if you wanted to be less strict about what
           counts as a double click you could use
           evt.originalEvent.detail > 1 instead */
           //alert($(this).attr("data-id"));
        recallMsg($(this).attr("data-id"));
      }
    });

    // 通过webSocket发送消息到服务端
    function sendMessage (){
        var inputMessage = $inputArea.val();  // 获取输入框的值

        var data = {type:0, message:inputMessage}

        if (inputMessage && connected) {
            $inputArea.val('');      // 清空输入框的值
            socket.send(JSON.stringify(data));  // 基于WebSocket连接发送消息
            console.log("send message:" + inputMessage);
        }
    }
    function recallMsg(id){
        console.log(id);
        var data = {type:3, message:name+" recalled a message.", related:id,name:name};
        socket.send(JSON.stringify(data));  // 基于WebSocket连接发送消息
        console.log(data);
    }
});

function activateIcon(name){
    $("i[data-id='" + name + "']").css("color", "#339533");
}

function deactivateIcon(name){
    $("i[data-id='" + name + "']").css("color", "#999999");
}

//function 

