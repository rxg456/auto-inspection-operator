package report

import (
	"bytes"
	"fmt"
	"html/template"
	"time"
)

// NodeMetric 表示主机指标
type NodeMetric struct {
	Name string
	// 硬盘使用率
	DiskNow        float64
	DiskOffset     float64
	DiskRate       float64
	DiskNowStatus  int
	DiskRateStatus int
	// Inode使用率
	InodeNow        float64
	InodeOffset     float64
	InodeRate       float64
	InodeNowStatus  int
	InodeRateStatus int
	// CPU使用率
	CPUNow        float64
	CPUOffset     float64
	CPURate       float64
	CPUNowStatus  int
	CPURateStatus int
	// 内存使用率
	MemNow        float64
	MemOffset     float64
	MemRate       float64
	MemNowStatus  int
	MemRateStatus int
	// 整体状态
	Status int
}

// InspectionData 巡检数据
type InspectionData struct {
	Node []NodeMetric
}

// ReportMetadata 报告元数据
type ReportMetadata struct {
	Business  string
	Title     string
	Date      string
	Timestamp string
}

// ReportData 报告数据
type ReportData struct {
	Metadata   ReportMetadata
	Inspection InspectionData
}

// Generator 报告生成器
type Generator struct {
	// 不再需要TemplatePath字段
}

// NewGenerator 创建新的报告生成器
func NewGenerator() *Generator {
	return &Generator{}
}

// GenerateHTML 生成HTML报告
func (g *Generator) GenerateHTML(data *ReportData) (string, error) {
	// 使用内置的模板变量中读取
	tmpl, err := template.New("report").Parse(ReportTemplate)
	if err != nil {
		return "", fmt.Errorf("解析模板失败: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("执行模板失败: %w", err)
	}

	return buf.String(), nil
}

// GenerateReport 根据指标数据生成报告
func GenerateReport(business string, nodes []NodeMetric) (*ReportData, error) {
	now := time.Now()
	data := &ReportData{
		Metadata: ReportMetadata{
			Business:  business,
			Title:     fmt.Sprintf("%s业务系统巡检报告", business),
			Date:      now.Format("2006-01-02"),
			Timestamp: now.Format("2006-01-02 15:04:05"),
		},
		Inspection: InspectionData{},
	}

	data.Inspection.Node = nodes

	return data, nil
}

// CheckThresholds 检查阈值并设置状态
func CheckThresholds(node *NodeMetric) {
	// 检查硬盘使用率
	if node.DiskNow > 80 {
		node.DiskNowStatus = 1
		node.Status = 1
	}

	// 检查硬盘使用率波动
	if node.DiskRate > 10 {
		node.DiskRateStatus = 1
		node.Status = 1
	}

	// 检查inode使用率
	if node.InodeNow > 60 {
		node.InodeNowStatus = 1
		node.Status = 1
	}

	// 检查inode使用率波动
	if node.InodeRate > 10 {
		node.InodeRateStatus = 1
		node.Status = 1
	}

	// 检查CPU使用率
	if node.CPUNow > 60 {
		node.CPUNowStatus = 1
		node.Status = 1
	}

	// 检查CPU使用率波动
	if node.CPURate > 10 {
		node.CPURateStatus = 1
		node.Status = 1
	}

	// 检查内存使用率
	if node.MemNow > 80 {
		node.MemNowStatus = 1
		node.Status = 1
	}

	// 检查内存使用率波动
	if node.MemRate > 10 {
		node.MemRateStatus = 1
		node.Status = 1
	}
}
