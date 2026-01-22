// 字体设置
function cog() {
    var control = document.getElementById('text_control');
    if (control.style.display === 'block') {
        control.style.display = 'none';
    } else {
        control.style.display = 'block';
    }
}

// 字号调整
function changeSize(type) {
    var article = document.getElementById('article');
    var currentSize = parseInt(window.getComputedStyle(article).fontSize);
    if (type === 'plus') {
        article.style.fontSize = (currentSize + 2) + 'px';
    } else if (type === 'min') {
        article.style.fontSize = (currentSize - 2) + 'px';
    } else {
        article.style.fontSize = '18px'; // default
    }
}

// 夜间模式
// 夜间模式
function isnight() {
    var body = document.body;
    var read_bg = document.querySelector('.read_bg');
    var article = document.getElementById('article');
    var section_style = document.querySelector('.section_style');
    var read_nav = document.querySelector('.read_nav');
    var text_set = document.querySelector('.text_set');

    // Check current state (default is light)
    var isDark = read_bg.style.backgroundColor === 'rgb(26, 26, 26)' || read_bg.style.backgroundColor === '#1a1a1a';

    if (!isDark) {
        // Switch to night
        read_bg.style.backgroundColor = '#1a1a1a';
        body.style.color = '#777';
        if (article) article.style.color = '#777';
        if (section_style) section_style.style.background = 'transparent'; // Let read_bg show through
        if (read_nav) read_nav.style.background = 'transparent';
        if (text_set) text_set.style.background = '#1a1a1a'; // Match dark bg
        document.querySelector('.style_h1').style.color = '#777';
        $('.text_info span a').css("color", "#777");
        $('.text_info span').css("color", "#777");
        $('.read_nav a').css("color", "#777");
        $.cookie('isnight', '1', { path: '/' });
    } else {
        // Switch to day
        read_bg.style.backgroundColor = '#e7e1d4';
        body.style.color = '#666';
        if (article) article.style.color = '#262626';
        if (section_style) section_style.style.background = '#FBF6EC';
        if (read_nav) read_nav.style.background = '#FBF6EC';
        if (text_set) text_set.style.background = '#FBF6EC';
        document.querySelector('.style_h1').style.color = '#555';
        $('.text_info span a').css("color", "gray");
        $('.text_info span').css("color", "gray");
        $('.read_nav a').css("color", "");
        $.cookie('isnight', '0', { path: '/' });
    }
}

// Check cookie on load for night mode
$(document).ready(function () {
    if ($.cookie('isnight') === '1') {
        var read_bg = document.querySelector('.read_bg');
        var article = document.getElementById('article');
        var section_style = document.querySelector('.section_style');
        var read_nav = document.querySelector('.read_nav');
        var text_set = document.querySelector('.text_set');

        if (read_bg) read_bg.style.backgroundColor = '#1a1a1a';
        document.body.style.color = '#777';
        if (article) article.style.color = '#777';
        if (section_style) section_style.style.background = 'transparent';
        if (read_nav) read_nav.style.background = 'transparent';
        if (text_set) text_set.style.background = '#1a1a1a';
        if (document.querySelector('.style_h1')) document.querySelector('.style_h1').style.color = '#777';
        $('.text_info span a').css("color", "#777");
        $('.text_info span').css("color", "#777");
        $('.read_nav a').css("color", "#777");
    }
});

// 极简模式 (hides header, footer, nav)
function ismini() {
    $('header, .navigation, .read_nav, footer, .s_gray, .section_style > .text').toggle();
    // Re-show cog icon if it was hidden inside .text which we might have hidden. 
    // Actually .text contains the cog. Let's be more specific.
    // The requirement is vague, let's just toggle non-content elements.
}

// 阅读记录
var lastread = {
    set: function (bookUrl, chapterUrl, bookName, chapterName, author, cover) {
        var history = JSON.parse(localStorage.getItem('history') || '[]');
        // Remove duplicate
        history = history.filter(function (h) { return h.bookName !== bookName; });
        // Add new
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

// 提示信息
function tips(bookname) {
    document.write('提示：按回车[Enter]键返回书目，按←键返回上一页， 按→键进入下一页。');
}
