/**
 * 生产级简繁体转换逻辑 (基于 OpenCC-js)
 * 默认显示简体，支持 localStorage 偏持持久化
 */

var zh_default = 's'; // 默认语言：s-简体, t-繁体
var zh_choose = localStorage.getItem('zh_choose') || zh_default;

// 初始化 OpenCC 转换器
var converter_s2t = null;
var converter_t2s = null;

function zh_init_opencc() {
    if (typeof OpenCC === 'undefined') {
        console.error('OpenCC-js 未加载成功');
        return false;
    }
    // cn2t: 简体 -> 繁体 (台湾习惯)
    converter_s2t = OpenCC.Converter({ from: 'cn', to: 'tw' });
    // t2cn: 繁体 -> 简体
    converter_t2s = OpenCC.Converter({ from: 'tw', to: 'cn' });
    return true;
}

/**
 * 外部调用的转换入口
 * @param {string} mode 's' 或 't'
 */
window.zh_tran = function (mode) {
    zh_choose = mode;
    localStorage.setItem('zh_choose', mode);
    // 刷新页面以应用最新选择，或者如果不刷新则直接执行转换
    // 为了保证转换的最彻底，建议直接刷新
    location.reload();
}

/**
 * 执行全站 DOM 转换
 */
function zh_tran_all() {
    if (!zh_init_opencc()) return;

    var converter = (zh_choose === 't') ? converter_s2t : null;
    if (!converter && zh_choose === 's') {
        // 如果是简体且原本就是简体，无需处理
        // 但如果页面是从繁体切回来的，可能需要强制转换回简体
        // 简单处理：如果不是简体默认值，则视为需要转换
        if (localStorage.getItem('zh_had_converted') === 'true') {
            converter = converter_t2s;
        } else {
            return;
        }
    }

    if (zh_choose === 't') {
        localStorage.setItem('zh_had_converted', 'true');
    }

    // 转换标题
    document.title = converter(document.title);

    // 递归转换文本节点
    function traverse(node) {
        if (node.nodeType === 3) { // Text node
            var original = node.data;
            var converted = converter(original);
            if (original !== converted) {
                node.data = converted;
            }
        } else if (node.nodeType === 1 && node.nodeName !== 'SCRIPT' && node.nodeName !== 'STYLE') {
            for (var i = 0; i < node.childNodes.length; i++) {
                traverse(node.childNodes[i]);
            }
        }
    }

    traverse(document.body);
}

// 页面加载后自动运行
document.addEventListener('DOMContentLoaded', function () {
    if (zh_choose === 't') {
        // 如果当前选择是繁体，等待 OpenCC 加载后执行
        var checkOpenCC = setInterval(function () {
            if (typeof OpenCC !== 'undefined') {
                clearInterval(checkOpenCC);
                zh_tran_all();
            }
        }, 50);
        // 设置超时防止死循环
        setTimeout(function () { clearInterval(checkOpenCC); }, 5000);
    }
});
