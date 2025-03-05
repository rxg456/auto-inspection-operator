package report

// ReportTemplate 包含HTML报告模板的内容
const ReportTemplate = `<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8" />
    <meta http-equiv="X-UA-Compatible" content="IE=edge" />
    <title>{{ .Metadata.Title }}</title>
    <meta name="description" content />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <link rel="stylesheet" href="https://cdn.staticfile.org/twitter-bootstrap/5.3.2/css/bootstrap.min.css">
    <style type="text/css">
      body {
        font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
        color: #333;
        background-color: #f8f9fa;
        line-height: 1.6;
      }
      .container {
        background-color: #fff;
        box-shadow: 0 0 15px rgba(0, 0, 0, 0.1);
        border-radius: 8px;
        padding: 25px;
        margin-top: 20px;
        margin-bottom: 20px;
      }
      .report-header {
        padding: 20px;
        background: linear-gradient(135deg, #0d6efd, #0a58ca);
        color: white;
        border-radius: 6px;
        margin-bottom: 25px;
        text-align: center;
      }
      .report-header h2 {
        margin-bottom: 5px;
        font-weight: 600;
      }
      .report-header p {
        opacity: 0.9;
        margin-bottom: 0;
      }
      .report-section {
        margin-bottom: 30px;
        border: 1px solid #e7e7e7;
        border-radius: 6px;
        padding: 20px;
        background-color: #fdfdfd;
      }
      .section-title {
        border-bottom: 2px solid #0d6efd;
        padding-bottom: 10px;
        margin-bottom: 20px;
        color: #0a58ca;
        font-weight: 500;
      }
      .table-container {
        width: 100%;
        margin-bottom: 20px;
        overflow-x: auto;
      }
      .table-responsive {
        /* 移除高度限制 */
        border-radius: 6px;
        overflow: visible;
      }
      .table {
        margin-bottom: 0;
        width: 100%;
        border-collapse: collapse;
        font-size: 12px; /* 减小字体大小以适应更多内容 */
      }
      .table thead {
        background-color: #f1f5fd;
        /* 移除sticky定位，提高邮件兼容性 */
        /* position: sticky; */
        /* top: 0; */
        /* z-index: 10; */
      }
      .table thead td {
        font-weight: 600;
        color: #495057;
        border-bottom: 2px solid #dee2e6;
        padding: 8px;
        text-align: center;
      }
      .table tbody td {
        padding: 6px;
        text-align: center;
        vertical-align: middle;
      }
      .table tbody tr:nth-of-type(odd) {
        background-color: rgba(0, 0, 0, 0.02);
      }
      .table tbody tr:hover {
        background-color: rgba(13, 110, 253, 0.05);
      }
      .badge {
        font-size: 0.9em;
        padding: 6px 10px;
        border-radius: 4px;
        font-weight: 500;
      }
      .text-bg-default {
        background-color: #f8f9fa;
        color: #212529;
      }
      .text-bg-success {
        background-color: #28a745;
      }
      .text-bg-danger {
        background-color: #dc3545;
      }
      .info-card {
        background-color: #f8f9fa;
        border-radius: 6px;
        padding: 15px;
        margin-bottom: 20px;
      }
      .info-section {
        border-left: 4px solid #0d6efd;
        padding-left: 12px;
        margin-bottom: 15px;
      }
      .status-item {
        display: flex;
        align-items: center;
        gap: 8px;
      }
      .definition-title {
        color: #0a58ca;
        font-weight: 600;
        margin-top: 12px;
        margin-bottom: 8px;
      }
      .collapsible-content {
        max-height: 0;
        overflow: hidden;
        transition: max-height 0.3s ease;
      }
      .collapsible-content.expanded {
        max-height: 1000px;
      }
      .report-footer {
        text-align: center;
        padding-top: 20px;
        color: #6c757d;
        font-size: 0.9em;
      }
    </style>
  </head>

  <body>
    <div class="container">
      <div class="report-header">
        <h2>{{ .Metadata.Business }}业务系统巡检报告</h2>
        <p>生成时间: {{ .Metadata.Date }} {{ .Metadata.Timestamp }}</p>
      </div>
      
      <div class="report-section">
        <h4 class="section-title">巡检说明</h4>
        
        <div class="info-card">
          <div class="info-section">
            <div class="definition-title">巡检报表取值说明</div>
            <p>以巡检时间为参考，巡检报表中的各项指标（硬盘使用率、inode使用率、CPU使用率、内存使用率）的取值规则如下：</p>
            
            <div class="definition-title">节点选择机制</div>
            <p><strong>两种方式:</strong> 系统支持通过明确的节点列表或标签自动发现进行巡检</p>
            <ul>
              <li><strong>节点列表:</strong> 直接指定IP地址列表，如<code>["192.18.0.1:9100", "192.18.0.2:9100"]</code></li>
              <li><strong>标签自动发现:</strong> 通过标签自动查找匹配的节点，如<code>{"business": "CRM", "env": "prod"}</code>，系统将查询所有具有这些标签的节点</li>
              <li><strong>混合模式:</strong> 同时支持指定节点列表和标签，将对两者找到的所有节点进行巡检</li>
            </ul>
            
            <div class="definition-title">1. 当前值（最近5分钟内的最新值）</div>
            <p><strong>定义:</strong> 当前值表示系统在最近 5 分钟内的最新状态。</p>
            <p><strong>来源:</strong> 从 Prometheus 监控系统获取的实时查询结果。</p>
            <p><strong>指标计算:</strong></p>
            <ul>
              <li><strong>CPU使用率:</strong> 计算过去60分钟内非空闲CPU时间的平均百分比，公式：<code>(1 - avg(irate(node_cpu_seconds_total{mode="idle"}[60m])) by (instance))*100</code></li>
              <li><strong>内存使用率:</strong> 计算已使用内存占总内存的百分比，公式：<code>(1 - (node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes))*100</code></li>
              <li><strong>硬盘使用率:</strong> 计算已使用硬盘空间使用率最大的分区，公式：<code>max(100 - ((node_filesystem_avail_bytes / node_filesystem_size_bytes) * 100))</code></li>
              <li><strong>inode使用率:</strong> 计算已使用inode使用率最大的分区，公式：<code>max(100 - ((node_filesystem_files_free / node_filesystem_files)*100))</code></li>
            </ul>
            
            <div class="definition-title">2. 24小时值（前24小时的一个值）</div>
            <p><strong>定义:</strong> 表示系统在前 24 小时的一个采样点数据。</p>
            <p><strong>来源:</strong> 从 Prometheus 监控系统中查询 24 小时前的历史记录，使用相同的查询方式但指定了时间偏移。</p>
            <p><strong>指标:</strong> 与当前值使用相同的计算公式，只是时间点不同。</p>
          
            <div class="definition-title">3. 差值（增减率百分比）</div>
            <p><strong>公式:</strong> 当前值% - 24小时前值% = 差值</p>
            <p><strong>含义:</strong> 反映当前系统状态相较于前 24 小时是否发生显著变化，并用百分比表示增减幅度。</p>
          </div>
        </div>
      </div>
      
      <div class="report-section">
        <h4 class="section-title">服务器巡检结果</h4>
          <div class="status-item">
            <span class="badge text-bg-success">正常</span>
            <span>巡检使用率正常</span>
          </div>
          <div class="status-item">
            <span class="badge text-bg-danger">异常</span>
            <span>硬盘高于80% or inode使用率高于60% or CPU平均使用率高于60% or 内存平均使用率高于80% or 波动幅度大于10%</span>
          </div>
        <div class="table-container">
          <table class="table table-bordered">
            <thead>
              <tr>
                <td style="width: 500px">主机IP</td>
                <td style="width: 150px">硬盘使用率</td>
                <td style="width: 150px">硬盘使用率前24h</td>
                <td style="width: 150px">硬盘使用率差值</td>
                <td style="width: 150px">inode使用率</td>
                <td style="width: 150px">inode使用率前24h</td>
                <td style="width: 150px">inode使用率差值</td>
                <td style="width: 150px">CPU使用率</td>
                <td style="width: 150px">CPU使用率前24h</td>
                <td style="width: 150px">CPU使用率差值</td>
                <td style="width: 150px">内存使用率</td>
                <td style="width: 150px">内存使用率前24h</td>
                <td style="width: 150px">内存使用率差值</td>
                <td style="width: 100px">状态</td>
              </tr>
            </thead>
            <tbody>
              {{ range .Inspection.Node }}
              <tr>
                <td>
                  <span class="badge text-bg-default">
                    {{ .Name }}
                  </span>
                </td>
                <td>
                  <span class="badge {{ if eq .DiskNowStatus 1 }}text-bg-danger{{ else }}text-bg-default{{ end }}">
                    {{ .DiskNow }}%
                  </span>
                </td>
                <td>
                  <span class="badge text-bg-default">
                    {{ .DiskOffset }}%
                  </span>
                </td>
                <td>
                  <span class="badge {{ if eq .DiskRateStatus 1 }}text-bg-danger{{ else }}text-bg-default{{ end }}">
                    {{ .DiskRate }}%
                  </span>
                </td>
                <td>
                  <span class="badge {{ if eq .InodeNowStatus 1 }}text-bg-danger{{ else }}text-bg-default{{ end }}">
                    {{ .InodeNow }}%
                  </span>
                </td>
                <td>
                  <span class="badge text-bg-default">
                    {{ .InodeOffset }}%
                  </span>
                </td>
                <td>
                  <span class="badge {{ if eq .InodeRateStatus 1 }}text-bg-danger{{ else }}text-bg-default{{ end }}">
                    {{ .InodeRate }}%
                  </span>
                </td>
                <td>
                  <span class="badge {{ if eq .CPUNowStatus 1 }}text-bg-danger{{ else }}text-bg-default{{ end }}">
                    {{ .CPUNow }}%
                  </span>
                </td>
                <td>
                  <span class="badge text-bg-default">
                    {{ .CPUOffset }}%
                  </span>
                </td>
                <td>
                  <span class="badge {{ if eq .CPURateStatus 1 }}text-bg-danger{{ else }}text-bg-default{{ end }}">
                    {{ .CPURate }}%
                  </span>
                </td>
                <td>
                  <span class="badge {{ if eq .MemNowStatus 1 }}text-bg-danger{{ else }}text-bg-default{{ end }}">
                    {{ .MemNow }}%
                  </span>
                </td>
                <td>
                  <span class="badge text-bg-default">
                    {{ .MemOffset }}%
                  </span>
                </td>
                <td>
                  <span class="badge {{ if eq .MemRateStatus 1 }}text-bg-danger{{ else }}text-bg-default{{ end }}">
                    {{ .MemRate }}%
                  </span>
                </td>
                <td>
                  <span class="badge {{ if eq .Status 1 }}text-bg-danger{{ else }}text-bg-success{{ end }}">
                    {{ if eq .Status 1 }}异常{{ else }}正常{{ end }}
                  </span>
                </td>
              </tr>
              {{ end }}
            </tbody>
          </table>
        </div>
      </div>
      
      <div class="report-footer">
        此报告由自动巡检系统生成 &copy; {{ .Metadata.Date }}
      </div>
    </div>
  </body>
</html>`
