package prometheus

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client 是Prometheus API客户端
type Client struct {
	URL    string
	Client *http.Client
}

// NewClient 创建一个新的Prometheus客户端
func NewClient(prometheusURL string) *Client {
	return &Client{
		URL: prometheusURL,
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// QueryResult 表示Prometheus查询结果
type QueryResult struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric map[string]string `json:"metric"`
			Value  []interface{}     `json:"value"`
		} `json:"result"`
	} `json:"data"`
}

// QueryRange 表示Prometheus范围查询结果
type QueryRangeResult struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric map[string]string `json:"metric"`
			Values [][]interface{}   `json:"values"`
		} `json:"result"`
	} `json:"data"`
}

// Query 执行Prometheus即时查询
func (c *Client) Query(ctx context.Context, query string, timestamp time.Time) (*QueryResult, error) {
	u, err := url.Parse(fmt.Sprintf("%s/api/v1/query", c.URL))
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("query", query)
	if !timestamp.IsZero() {
		q.Set("time", timestamp.Format(time.RFC3339))
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var result QueryResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// 获取CPU使用率查询
func CPUUsageQuery(instance string, labels map[string]string) string {
	var conditions []string

	// 添加instance条件（如果提供）
	if instance != "" {
		conditions = append(conditions, fmt.Sprintf(`instance="%s"`, instance))
	}

	// 添加mode=idle条件（必需的）
	conditions = append(conditions, `mode="idle"`)

	// 添加labels条件（如果提供）
	if labels != nil {
		for k, v := range labels {
			conditions = append(conditions, fmt.Sprintf(`%s="%s"`, k, v))
		}
	}

	// 用逗号连接所有条件
	selectorStr := strings.Join(conditions, ",")

	return fmt.Sprintf(`(1 - avg(irate(node_cpu_seconds_total{%s}[60m])) by (instance))*100`, selectorStr)
}

// 获取内存使用率查询
func MemoryUsageQuery(instance string, labels map[string]string) string {
	var conditions []string

	// 添加instance条件（如果提供）
	if instance != "" {
		conditions = append(conditions, fmt.Sprintf(`instance="%s"`, instance))
	}

	// 添加labels条件（如果提供）
	if labels != nil {
		for k, v := range labels {
			conditions = append(conditions, fmt.Sprintf(`%s="%s"`, k, v))
		}
	}

	// 用逗号连接所有条件
	selectorStr := strings.Join(conditions, ",")

	return fmt.Sprintf(`(1 - (node_memory_MemAvailable_bytes{%s} / node_memory_MemTotal_bytes{%s}))*100`, selectorStr, selectorStr)
}

// 获取硬盘使用率查询
func DiskUsageQuery(instance string, labels map[string]string) string {
	var conditions []string

	// 添加instance条件（如果提供）
	if instance != "" {
		conditions = append(conditions, fmt.Sprintf(`instance="%s"`, instance))
	}

	// 添加labels条件（如果提供）
	if labels != nil {
		for k, v := range labels {
			conditions = append(conditions, fmt.Sprintf(`%s="%s"`, k, v))
		}
	}

	// 用逗号连接所有条件
	selectorStr := strings.Join(conditions, ",")

	// 如果条件为空，添加一个fstype!=条件以排除特殊文件系统
	if len(conditions) == 0 {
		selectorStr = `fstype!="tmpfs",fstype!="rootfs"`
	}

	return fmt.Sprintf(`max(100 - ((node_filesystem_avail_bytes{%s} / node_filesystem_size_bytes{%s}) * 100))`, selectorStr, selectorStr)
}

// 获取inode使用率查询
func InodeUsageQuery(instance string, labels map[string]string) string {
	var conditions []string

	// 添加instance条件（如果提供）
	if instance != "" {
		conditions = append(conditions, fmt.Sprintf(`instance="%s"`, instance))
	}

	// 添加labels条件（如果提供）
	if labels != nil {
		for k, v := range labels {
			conditions = append(conditions, fmt.Sprintf(`%s="%s"`, k, v))
		}
	}

	// 用逗号连接所有条件
	selectorStr := strings.Join(conditions, ",")

	// 如果条件为空，添加一个fstype!=条件以排除特殊文件系统
	if len(conditions) == 0 {
		selectorStr = `fstype!="tmpfs",fstype!="rootfs"`
	}

	return fmt.Sprintf(`max(100 - ((node_filesystem_files_free{%s} / node_filesystem_files{%s})*100))`, selectorStr, selectorStr)
}

// ParseValue 从查询结果中解析浮点值
func ParseValue(result *QueryResult) (float64, error) {
	if result == nil || len(result.Data.Result) == 0 {
		return 0, fmt.Errorf("no results returned")
	}

	// 获取值
	value := result.Data.Result[0].Value[1]
	if strValue, ok := value.(string); ok {
		var floatValue float64
		_, err := fmt.Sscanf(strValue, "%f", &floatValue)
		if err != nil {
			return 0, err
		}
		return floatValue, nil
	}

	return 0, fmt.Errorf("failed to parse value: %v", value)
}

// GetNodesByLabels 根据标签查询所有符合条件的节点
func (c *Client) GetNodesByLabels(ctx context.Context, labels map[string]string) ([]string, error) {
	if labels == nil || len(labels) == 0 {
		return nil, fmt.Errorf("labels不能为空")
	}

	// 构建标签条件
	var conditions []string
	for k, v := range labels {
		conditions = append(conditions, fmt.Sprintf(`%s="%s"`, k, v))
	}
	labelSelector := strings.Join(conditions, ",")

	// 使用up指标查询符合标签条件的所有节点
	// up指标通常存在于所有节点，用于表示节点是否在线
	query := fmt.Sprintf(`up{%s}`, labelSelector)

	result, err := c.Query(ctx, query, time.Time{})
	if err != nil {
		return nil, fmt.Errorf("查询节点失败: %w", err)
	}

	if result == nil || len(result.Data.Result) == 0 {
		return nil, fmt.Errorf("未找到符合条件的节点")
	}

	// 提取所有唯一的instance值
	nodeMap := make(map[string]struct{})
	for _, item := range result.Data.Result {
		if instance, ok := item.Metric["instance"]; ok {
			nodeMap[instance] = struct{}{}
		}
	}

	// 转换为切片
	nodes := make([]string, 0, len(nodeMap))
	for node := range nodeMap {
		nodes = append(nodes, node)
	}

	return nodes, nil
}
