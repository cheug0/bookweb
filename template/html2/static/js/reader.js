// 字体设置 - Toggles the native html2 settings menu
function cog() {
    var menu = $('.mlfy_main_sz');
    if (menu.is(':visible')) {
        menu.hide();
    } else {
        menu.show();
    }
}

// Bind events when document is ready
$(document).ready(function () {
    // Close button
    $('.mlfy_main_sz .close').click(function () {
        $('.mlfy_main_sz').hide();
    });

    // Font size
    $('.dxl').click(function () { changeSize('min'); });
    $('.dxr').click(function () { changeSize('plus'); });

    // Theme Switching
    // Themes correspond to .c1 -> .bg1, .c2 -> .bg2, ... .c8 -> .bg8 (c7 is night)
    $('.mlfy_main_sz .c1, .mlfy_main_sz .c2, .mlfy_main_sz .c3, .mlfy_main_sz .c4, .mlfy_main_sz .c5, .mlfy_main_sz .c6, .mlfy_main_sz .c7, .mlfy_main_sz .c8').click(function () {
        var className = $(this).attr('class');
        // Extract the number from c1, c2, etc. (e.g. "c1" -> "1")
        var match = className.match(/c(\d+)/);
        if (match) {
            var themeId = match[1];
            setTheme(themeId);
        }
    });

    // Font Family Switching
    // zt span click. The HTML has <span class="zt">...</span>
    // We map them by index: 0->Default, 1->SimSun(zt2), 2->Kaiti(zt3), etc.
    // Based on CSS: .zt2=SimSun, .zt3=Kaiti, .zt4=Qiti, .zt5=Source, .zt6=PingFang
    // The clickable spans are: 1:Yahei(default/zt1?), 2:Songti, 3:Kaiti, 4:Qiti, 5:Source, 6:Pingfang
    $('.mlfy_main_sz .zt').click(function () {
        var index = $('.mlfy_main_sz .zt').index(this);
        // index 0 is Yahei (default), so let's use index+1 as the id, but css starts zt2.
        // Let's assume:
        // Index 0 (Yahei) -> Remove zt class / Default
        // Index 1 (Songti) -> zt2
        // Index 2 (Kaiti) -> zt3
        // Index 3 (Qiti) -> zt4
        // Index 4 (Source) -> zt5
        // Index 5 (Pingfang) -> zt6
        setFont(index);
    });

    // Initialize from cookies
    var savedTheme = $.cookie('theme_id');
    var savedFont = $.cookie('font_id');

    // Compatibility: check old isnight cookie
    if (!savedTheme && $.cookie('isnight') === '1') {
        savedTheme = '7';
    }

    if (savedTheme) {
        setTheme(savedTheme);
    } else {
        // Default theme usually bg6 (based on html template default class on body)
        // But body has class "bg6" in HTML, so we don't need to force it if not set.
        // Just mark the menu active state
        $('.mlfy_main_sz .c6').addClass('hover');
    }

    if (savedFont) {
        setFont(parseInt(savedFont));
    } else {
        $('.mlfy_main_sz .zt').eq(0).addClass('hover');
    }
});

// Theme Setter
function setTheme(id) {
    // Update Body Class
    var body = $('body');
    // Remove all bg classes
    body.removeClass('bg1 bg2 bg3 bg4 bg5 bg6 bg7 bg8');
    body.addClass('bg' + id);

    // Update Menu Selection
    $('.mlfy_main_sz i[class^="c"]').removeClass('hover');
    $('.mlfy_main_sz .c' + id).addClass('hover');

    // Save Cookie
    $.cookie('theme_id', id, { path: '/' });
    $.cookie('isnight', (id === '7' ? '1' : '0'), { path: '/' }); // Sync night cookie
}

// Font Setter
function setFont(index) {
    var article = $('#mlfy_main_text');
    // Remove all zt classes
    article.removeClass('zt1 zt2 zt3 zt4 zt5 zt6');

    // Map index to css class
    // 0: Yahei -> No class or zt1 (CSS doesn't explicitly have zt1, likely default)
    // 1: Songti -> zt2
    // 2: Kaiti -> zt3
    // 3: Qiti -> zt4
    // 4: Source -> zt5
    // 5: Pingfang -> zt6
    if (index > 0) {
        article.addClass('zt' + (index + 1));
    }

    // Update Menu Selection
    $('.mlfy_main_sz .zt').removeClass('hover');
    $('.mlfy_main_sz .zt').eq(index).addClass('hover');

    // Save Cookie
    $.cookie('font_id', index, { path: '/' });
}

// 字号调整
function changeSize(type) {
    var article = $('#mlfy_main_text');
    var currentSize = parseInt(article.css('fontSize'));
    if (!currentSize) currentSize = 20;

    if (type === 'plus') {
        article.css('fontSize', (currentSize + 2) + 'px');
    } else if (type === 'min') {
        article.css('fontSize', (currentSize - 2) + 'px');
    }
}

// 极简模式
function ismini() {
    $('.top, .fl, .mlfy_main_l i, .mlfy_page').toggle();
}

// Reading History
var lastread = {
    set: function (bookUrl, chapterUrl, bookName, chapterName, author, cover) {
        var history = JSON.parse(localStorage.getItem('history') || '[]');
        history = history.filter(function (h) { return h.bookName !== bookName; });
        history.unshift({
            bookUrl: bookUrl,
            chapterUrl: chapterUrl,
            bookName: bookName,
            chapterName: chapterName,
            author: author,
            cover: cover,
            time: new Date().getTime()
        });
        if (history.length > 20) history.pop();
        localStorage.setItem('history', JSON.stringify(history));
    }
};

function tips(bookname) {
    document.write('提示：按回车[Enter]键返回书目，按←键返回上一页， 按→键进入下一页。');
}
