/*
逆水行舟 QQ群: 249310348
*/

//reflush authcode
function reloadcode() {
    newcode = $('#showcode').prop("src") + '?' + Math.random();
    $('#showcode').prop("src", newcode);
}


function register_check() {

    if ($('#regname').val() == '' || $('#regpass').val() == '' || $('#repass').val() == '' || $('#email').val() == '' || $('#authcode').val() == '') {
        alert('每项都必须填写');
        return false;
    };

    if ($('#regpass').val() !== $('#repass').val()) {
        alert('两次输入的密码不一致');
        return false;
    }
    return true;

}
//addbookcase
function addbookcase(aid, name, cid, cname) {
    if ($.cookie('ss_userid') && $.cookie('PHPSESSID') != -1) {
        rico_data = {
            articleid: aid,
            articlename: name,
            chapterid: cid,
            chaptername: cname
        }

        $.ajax({
            type: "post",
            url: "/addbookcase/",
            data: rico_data,
            success: function (data) {
                alert(data);
            }
        })

    } else {
        if (window.confirm("\n永久书架需要登录才能使用，转到登录页面吗？")) {
            window.location.href = "/login/";
        } else {
            return;
        }
    }
}

//delbookcase
function delbookcase(aid) {
    if (window.confirm("\n确定要删除吗？")) {
        if ($.cookie('ss_userid') && $.cookie('PHPSESSID') != -1) {
            rico_data = {
                articleid: aid,
            },

                $.ajax({
                    type: "post",
                    url: "/delbookcase/",
                    data: rico_data,
                    success: function (data) {
                        alert(data);
                        window.location.reload();
                    }
                })
        }
    } else {
        return;
    }
}


var checkbg = "#ccc";
//内容页用户设置
function nr_setbg(intype) {
    var huyandiv = document.getElementById("huyandiv");
    var light = document.getElementById("lightdiv");
    if (intype == "huyan") {
        if (huyandiv.style.backgroundColor == "") {
            set("light", "huyan");
            document.cookie = "light=huyan;path=/";
        } else {
            set("light", "no");
            document.cookie = "light=no;path=/";
        }
    }
    if (intype == "light") {
        if (light.innerHTML == "关灯") {
            set("light", "yes");
            document.cookie = "light=yes;path=/";
        } else {
            set("light", "no");
            document.cookie = "light=no;path=/";
        }
    }
    if (intype == "big") {
        set("font", "big");
        document.cookie = "font=big;path=/";
    }
    if (intype == "middle") {
        set("font", "middle");
        document.cookie = "font=middle;path=/";
    }
    if (intype == "small") {
        set("font", "small");
        document.cookie = "font=small;path=/";
    }
}

//内容页读取设置
function getset() {
    var strCookie = document.cookie;
    var arrCookie = strCookie.split("; ");
    var light;
    var font;

    for (var i = 0; i < arrCookie.length; i++) {
        var arr = arrCookie[i].split("=");
        if ("light" == arr[0]) {
            light = arr[1];
            break;
        }
    }
    for (var i = 0; i < arrCookie.length; i++) {
        var arr = arrCookie[i].split("=");
        if ("font" == arr[0]) {
            font = arr[1];
            break;
        }
    }

    //light
    if (light == "yes") {
        set("light", "yes");
    } else if (light == "no") {
        set("light", "no");
    } else if (light == "huyan") {
        set("light", "huyan");
    }
    //font
    if (font == "big") {
        set("font", "big");
    } else if (font == "middle") {
        set("font", "middle");
    } else if (font == "small") {
        set("font", "small");
    } else {
        set("", "");
    }
}

//内容页应用设置
function set(intype, p) {
    var nr_body = document.getElementById("nr_body"); //页面body
    var huyandiv = document.getElementById("huyandiv"); //护眼div
    var lightdiv = document.getElementById("lightdiv"); //灯光div
    var fontfont = document.getElementById("fontfont"); //字体div
    var fontbig = document.getElementById("fontbig"); //大字体div
    var fontmiddle = document.getElementById("fontmiddle"); //中字体div
    var fontsmall = document.getElementById("fontsmall"); //小字体div
    var nr1 = document.getElementById("nr1"); //内容div
    var nr_title = document.getElementById("nr_title"); //文章标题
    var nr_title = document.getElementById("nr_title"); //文章标题
    var shuqian_2 = document.getElementById("shuqian"); //书签链接

    var pt_mulu = document.getElementById("pt_mulu");
    var pt_index = document.getElementById("pt_index");
    var pt_bookcase = document.getElementById("pt_bookcase");
    var pb_prev = document.getElementById("pb_prev");
    var pb_mulu = document.getElementById("pb_mulu");
    var pb_next = document.getElementById("pb_next");

    //初始化
    //document.getElementsByName("page_nr")[1].style.color = "#000";

    //灯光
    if (intype == "light") {
        if (p == "yes") {
            var cssText = "background-color:#666;color:#ccc;border:1px solid #444";
            //关灯
            lightdiv.innerHTML = "开灯";
            nr_body.style.backgroundColor = "#444";
            huyandiv.style.backgroundColor = "";
            nr_title.style.color = "#ddd";
            nr1.style.color = "#999";
            shuqian_2.style.cssText = cssText;
            pt_mulu.style.cssText = cssText;
            pt_index.style.cssText = cssText;
            pt_bookcase.style.cssText = cssText;
            pb_prev.style.cssText = cssText;
            pb_mulu.style.cssText = cssText;
            pb_next.style.cssText = cssText;
        } else if (p == "no") {
            //开灯
            lightdiv.innerHTML = "关灯";
            nr_body.style.backgroundColor = "rgb(251, 246, 236)";
            nr1.style.color = "#000";
            nr_title.style.color = "#000";
            huyandiv.style.backgroundColor = "";
            shuqian_2.style.cssText = "";
            pt_mulu.style.cssText = "";
            pt_index.style.cssText = "";
            pt_bookcase.style.cssText = "";
            pb_prev.style.cssText = "";
            pb_mulu.style.cssText = "";
            pb_next.style.cssText = "";
        } else if (p == "huyan") {
            var cssText = "background-color:rgb(204, 226, 191);color:green;border:1px solid rgb(187, 214, 170)";
            //护眼
            lightdiv.innerHTML = "关灯";
            huyandiv.style.backgroundColor = checkbg;
            nr_body.style.backgroundColor = "rgb(220, 236, 210)";
            nr1.style.color = "#333";
            shuqian_2.style.cssText = cssText;
            pt_mulu.style.cssText = cssText;
            pt_index.style.cssText = cssText;
            pt_bookcase.style.cssText = cssText;
            pb_prev.style.cssText = cssText;
            pb_mulu.style.cssText = cssText;
            pb_next.style.cssText = cssText;
        }
    }
    //字体
    if (intype == "font") {
        //alert(p);
        fontbig.style.backgroundColor = "";
        fontmiddle.style.backgroundColor = "";
        fontsmall.style.backgroundColor = "";
        if (p == "big") {
            fontbig.style.backgroundColor = checkbg;
            nr1.style.fontSize = "24px";
        }
        if (p == "middle") {
            fontmiddle.style.backgroundColor = checkbg;
            nr1.style.fontSize = "20px";
        }
        if (p == "small") {
            fontsmall.style.backgroundColor = checkbg;
            nr1.style.fontSize = "16px";
        }
    }
}

//分类显示隐藏
function toggleSort() {
    $('#submenu').slideToggle();
}

// 跳转书架（带登录检查）
// 原生获取 Cookie 函数
function getCookie(name) {
    var arr, reg = new RegExp("(^| )" + name + "=([^;]*)(;|$)");
    if (arr = document.cookie.match(reg))
        return unescape(arr[2]);
    else
        return null;
}

// 跳转书架（带登录检查）
function toBookcase() {
    var username = getCookie('username');
    if (username) {
        window.location.href = "/bookcase";
    } else {
        if (window.confirm("\n我的书架需要登录才能使用，转到登录页面吗？")) {
            window.location.href = "/login";
        }
    }
}

// 更新底部用户信息 (在 foot.html 中调用)
function updateFooter() {
    var user = getCookie('username');
    var html = "";
    if (user) {
        html = '<a href="/user" rel="nofollow">用户中心(' + user + ')</a>';
    } else {
        html = '<a href="/login" rel="nofollow">登录</a>';
    }
    //同时保留原有的“用户中心”链接做备份或者替换
    //逻辑：foot.html 中放置一个容器 id="foot_user"
    if ($("#foot_user").length > 0) {
        $("#foot_user").html(html);
    }
}