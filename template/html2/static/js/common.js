function count() {
    //JS 统计代码
}
function gotop() { $('body,html').animate({ scrollTop: 0 }, 600); }
function gofooter() { $('body,html').animate({ scrollTop: $(document).height() }, 600); }
function lazy() { $("img.lazy").lazyload({ effect: "fadeIn" }) }

function desc(obj) {
    $(obj).text() == '倒序 ↑' ? $(obj).text('正序 ↓') : $(obj).text('倒序 ↑');
    let lis = $("#chapterList").children();
    $("#chapterList").empty();
    for (let i = lis.length - 1; i >= 0; i--) {
        $("#chapterList").append(lis.eq(i).clone())
    }
}

function addbookcase(articleid, articlename, chapterid, chaptername) {
    if (chapterid && chaptername) {
        // Add bookmark
        $.ajax({
            url: "/bookmark/add",
            type: "POST",
            data: { articleid: articleid, chapterid: chapterid },
            dataType: "json",
            success: function (res) {
                alert(res.message);
            },
            error: function () {
                alert("请求失败，请稍后重试");
            }
        });
    } else {
        // Add to bookshelf
        $.ajax({
            url: "/bookcase/add",
            type: "POST",
            data: { articleid: articleid },
            dataType: "json",
            success: function (res) {
                alert(res.message);
            },
            error: function () {
                alert("请求失败，请稍后重试");
            }
        });
    }
}

function click_fav() {
    var url = window.location.href;
    var title = document.title;
    try {
        window.external.addFavorite(url, title);
    } catch (e) {
        try {
            window.sidebar.addPanel(title, url, "");
        } catch (e) {
            alert("加入收藏失败，请使用Ctrl+D进行添加");
        }
    }
}