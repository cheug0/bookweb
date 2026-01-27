// plugin.go (db_optimizer)
// 数据库优化插件
// 提供数据库连接检查及手动优化功能
package db_optimizer

import (
	"bookweb/utils"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

// Plugin 数据库优化插件
type Plugin struct {
	config map[string]interface{}
}

// New 创建插件实例
func New() *Plugin {
	return &Plugin{}
}

// Name 返回插件名称
func (p *Plugin) Name() string {
	return "db_optimizer"
}

// Init 初始化
func (p *Plugin) Init(cfg map[string]interface{}) error {
	p.config = cfg
	return nil
}

// GetRoutes 获取路由
func (p *Plugin) GetRoutes() map[string]http.HandlerFunc {
	return map[string]http.HandlerFunc{
		"/admin/plugin/db_optimizer":          p.Index,
		"/admin/plugin/db_optimizer/check":    p.Check,
		"/admin/plugin/db_optimizer/optimize": p.Optimize,
	}
}

// Shutdown 关闭
func (p *Plugin) Shutdown() error {
	return nil
}

// Index 页面
func (p *Plugin) Index(w http.ResponseWriter, r *http.Request) {
	// 简单的内嵌模板
	const tplStr = `
<!DOCTYPE html>
<html>
<head>
    <title>数据库优化</title>
    <style>
        body { font-family: sans-serif; padding: 20px; background: #f0f2f5; }
        .container { background: #fff; padding: 20px; border-radius: 8px; box-shadow: 0 2px 12px 0 rgba(0,0,0,0.1); max-width: 800px; margin: 0 auto; }
        h2 { margin-top: 0; color: #1f2d3d; border-bottom: 1px solid #eee; padding-bottom: 15px; }
        .btn { display: inline-block; padding: 10px 20px; font-size: 14px; cursor: pointer; text-align: center; text-decoration: none; outline: none; color: #fff; background-color: #409eff; border: none; border-radius: 4px; box-shadow: 0 2px 0 #2b85e4; }
        .btn:active { box-shadow: 0 1px 0 #2b85e4; transform: translateY(1px); }
        .btn-success { background-color: #67c23a; box-shadow: 0 2px 0 #529b2e; }
        .btn-success:active { box-shadow: 0 1px 0 #529b2e; }
        .status-box { margin-top: 20px; padding: 15px; border: 1px solid #EBEEF5; background-color: #F2F6FC; border-radius: 4px; }
        .status-item { padding: 8px 0; border-bottom: 1px dashed #ddd; }
        .status-item:last-child { border-bottom: none; }
        .missing { color: #F56C6C; font-weight: bold; }
        .ok { color: #67c23a; font-weight: bold; }
        .loading { color: #909399; }
    </style>
</head>
<body>
    <div class="container">
        <h2>此工具将检测并修复数据库索引缺失问题</h2>
        <div>
            <button class="btn" onclick="checkDB()">开始检测</button>
            <button class="btn btn-success" id="fixBtn" onclick="optimizeDB()" style="display:none; margin-left: 10px;">一键修复</button>
        </div>
        <div id="status" class="status-box" style="display:none;"></div>
    </div>

    <script>
        function checkDB() {
            const statusDiv = document.getElementById('status');
            const fixBtn = document.getElementById('fixBtn');
            statusDiv.style.display = 'block';
            statusDiv.innerHTML = '<div class="loading">正在检测数据库索引...</div>';
            fixBtn.style.display = 'none';

            fetch('/admin/plugin/db_optimizer/check')
                .then(res => res.json())
                .then(data => {
                    let html = '';
                    let hasMissing = false;
                    data.results.forEach(item => {
                        const stateClass = item.missing ? 'missing' : 'ok';
                        const stateText = item.missing ? '缺失 (建议添加)' : '正常';
                        if (item.missing) hasMissing = true;
                        html += '<div class="status-item">' +
                                '<strong>表:</strong> ' + item.table + ' | ' +
                                '<strong>索引:</strong> ' + item.index + ' (' + item.columns + ') ' +
                                '<span class="' + stateClass + '" style="float:right">' + stateText + '</span>' +
                                '</div>';
                    });
                    
                    if (data.results.length === 0) {
                        html = '<div class="status-item">未定义检测规则</div>';
                    }

                    statusDiv.innerHTML = html;
                    if (hasMissing) {
                        fixBtn.style.display = 'inline-block';
                    }
                })
                .catch(err => {
                    statusDiv.innerHTML = '<div class="missing">检测失败: ' + err + '</div>';
                });
        }

        function optimizeDB() {
            const statusDiv = document.getElementById('status');
            const fixBtn = document.getElementById('fixBtn');
            const originalContent = statusDiv.innerHTML;
            
            if (!confirm('确定要执行优化操作吗？这将对数据库进行 ALTER TABLE 操作。')) return;

            statusDiv.innerHTML += '<div class="loading" style="margin-top:10px; border-top:1px solid #ccc; padding-top:10px;">正在执行修复，请稍候...</div>';
            fixBtn.disabled = true;

            fetch('/admin/plugin/db_optimizer/optimize', { method: 'POST' })
                .then(res => res.json())
                .then(data => {
                    alert(data.message);
                    checkDB(); // 重新检测
                    fixBtn.disabled = false;
                })
                .catch(err => {
                    alert('修复失败: ' + err);
                    statusDiv.innerHTML = originalContent;
                    fixBtn.disabled = false;
                });
        }
    </script>
</body>
</html>
`
	t, _ := template.New("index").Parse(tplStr)
	t.Execute(w, nil)
}

// Check 检测索引
func (p *Plugin) Check(w http.ResponseWriter, r *http.Request) {
	results := []map[string]interface{}{}

	// 定义需要检查的索引
	indices := []struct {
		Table   string
		Index   string
		Columns string // 用于显示
	}{
		{"jieqi_article_chapter", "idx_article_chapterorder", "articleid, chapterorder"},
		{"jieqi_article_article", "idx_allvisit", "allvisit"},
		{"jieqi_article_article", "idx_sort_display_visit", "sortid, display, allvisit"},
		{"jieqi_article_article", "idx_sort_lastupdate", "sortid, lastupdate"},
		{"sort", "idx_weight", "weight"},
	}

	for _, idx := range indices {
		exists, err := p.checkIndexExists(idx.Table, idx.Index)
		missing := false
		if err != nil {
			utils.LogError("DBOptimizer", "Check index error: %v", err)
			// 出错视为不确定，或者缺失
			missing = true
		} else if !exists {
			missing = true
		}

		results = append(results, map[string]interface{}{
			"table":   idx.Table,
			"index":   idx.Index,
			"columns": idx.Columns,
			"missing": missing,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"results": results})
}

// Optimize 修复索引
func (p *Plugin) Optimize(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// 执行创建索引 SQL
	sqls := []string{
		"CREATE INDEX idx_article_chapterorder ON jieqi_article_chapter(articleid, chapterorder)",
		"CREATE INDEX idx_allvisit ON jieqi_article_article(allvisit DESC)",
		"CREATE INDEX idx_sort_display_visit ON jieqi_article_article(sortid, display, allvisit DESC)",
		"CREATE INDEX idx_sort_lastupdate ON jieqi_article_article(sortid, lastupdate DESC)",
		"CREATE INDEX idx_weight ON sort(weight ASC)",
	}

	successCount := 0
	for _, sqlStr := range sqls {
		// 简单起见，忽略已存在错误（MySQL CREATE INDEX IF NOT EXISTS 并不是所有版本都支持，或者语法较长）
		// 我们直接执行，如果报错（比如 Duplicate key name）则忽略
		_, err := utils.Db.Exec(sqlStr)
		if err == nil {
			successCount++
		} else {
			// 检查是否是索引已存在错误
			//Error 1061: Duplicate key name
			if strings.Contains(err.Error(), "Duplicate key name") || strings.Contains(err.Error(), "already exists") {
				// 忽略
			} else {
				utils.LogError("DBOptimizer", "Failed to execute SQL: %s, Error: %v", sqlStr, err)
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("优化完成，尝试执行了 %d 个操作", len(sqls)),
	})
}

func (p *Plugin) checkIndexExists(table, indexName string) (bool, error) {
	// MySQL 特定查询
	query := fmt.Sprintf("SHOW INDEX FROM %s WHERE Key_name = ?", table)
	rows, err := utils.Db.Query(query, indexName)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	return rows.Next(), nil
}
