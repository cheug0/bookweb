function count(){
    //JS 统计代码
}
function gotop(){$('body,html').animate({scrollTop:0},600);}
function gofooter(){$('body,html').animate({scrollTop:$(document).height()},600);}
function lazy(){$("img.lazy").lazyload({effect : "fadeIn"})}

function desc(obj){
    $(obj).text() == '倒序 ↑' ? $(obj).text('正序 ↓') : $(obj).text('倒序 ↑');
    let lis =  $("#chapterList").children();
    $("#chapterList").empty();
    for(let i = lis.length-1; i >= 0; i--) {
        $("#chapterList").append(lis.eq(i).clone())
    }
}