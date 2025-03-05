package inspection

import (
	"context"
	"fmt"
	"math"
	"time"

	devopsv1 "github.com/rxg456/auto-inspection-operator/api/v1"
	"github.com/rxg456/auto-inspection-operator/internal/controller/mail"
	"github.com/rxg456/auto-inspection-operator/internal/controller/prometheus"
	"github.com/rxg456/auto-inspection-operator/internal/controller/report"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// Inspector 巡检器
type Inspector struct {
	prometheusClient *prometheus.Client
	reportGenerator  *report.Generator
	mailSender       *mail.Sender
	inspection       *devopsv1.AutoInspection
}

// NewInspector 创建巡检器
func NewInspector(inspection *devopsv1.AutoInspection) (*Inspector, error) {
	// 创建Prometheus客户端
	prometheusClient := prometheus.NewClient(inspection.Spec.PrometheusURL)

	// 创建报告生成器
	reportGenerator := report.NewGenerator()

	// 创建邮件发送器
	mailSender := mail.NewSender(inspection.Spec.SMTP)

	return &Inspector{
		prometheusClient: prometheusClient,
		reportGenerator:  reportGenerator,
		mailSender:       mailSender,
		inspection:       inspection,
	}, nil
}

// RunInspection 执行巡检
func (i *Inspector) RunInspection(ctx context.Context) error {
	logger := log.FromContext(ctx)
	logger.Info("开始执行巡检")

	nodeMetrics := []report.NodeMetric{}

	// 当前时间
	now := time.Now()
	// 24小时前
	dayAgo := now.Add(-24 * time.Hour)

	// 获取需要巡检的节点列表
	nodes := i.inspection.Spec.InspectionObject.Hosts.Nodes
	labels := i.inspection.Spec.InspectionObject.Hosts.Labels

	// 当nodes为空但labels不为空时，通过labels查询节点
	if len(nodes) == 0 && len(labels) > 0 {
		logger.Info("节点列表为空但有标签，尝试通过标签查询节点", "labels", labels)
		var err error
		nodes, err = i.prometheusClient.GetNodesByLabels(ctx, labels)
		if err != nil {
			logger.Error(err, "通过标签查询节点失败")
			return fmt.Errorf("通过标签查询节点失败: %w", err)
		}
		logger.Info("通过标签查询到节点", "count", len(nodes), "nodes", nodes)
	}

	// 如果节点列表仍为空，则返回错误
	if len(nodes) == 0 {
		logger.Error(fmt.Errorf("节点列表为空"), "无法执行巡检")
		return fmt.Errorf("节点列表为空，无法执行巡检")
	}

	// 获取所有主机的指标
	for _, node := range nodes {
		logger.Info("巡检主机", "node", node)

		metrics, err := i.collectNodeMetrics(ctx, node, labels, now, dayAgo)
		if err != nil {
			logger.Error(err, "获取主机指标失败", "node", node)
			continue
		}

		// 检查阈值并设置状态
		report.CheckThresholds(metrics)
		nodeMetrics = append(nodeMetrics, *metrics)
	}

	// 生成报告数据
	reportData, err := report.GenerateReport(i.inspection.Spec.InspectionObject.Business, nodeMetrics)
	if err != nil {
		return fmt.Errorf("生成报告数据失败: %w", err)
	}

	// 生成HTML报告
	htmlReport, err := i.reportGenerator.GenerateHTML(reportData)
	if err != nil {
		return fmt.Errorf("生成HTML报告失败: %w", err)
	}

	// 发送邮件
	subject := fmt.Sprintf("%s业务系统巡检报告 - %s", i.inspection.Spec.InspectionObject.Business, time.Now().Format("2006-01-02"))
	err = i.mailSender.SendMail(i.inspection.Spec.NotifyTo, subject, htmlReport)
	if err != nil {
		return fmt.Errorf("发送邮件失败: %w", err)
	}

	logger.Info("巡检任务完成")
	return nil
}

// collectNodeMetrics 收集节点指标
func (i *Inspector) collectNodeMetrics(
	ctx context.Context,
	node string,
	labels map[string]string,
	now, dayAgo time.Time,
) (*report.NodeMetric, error) {
	metrics := &report.NodeMetric{
		Name: node,
	}

	// 采集CPU使用率
	cpuQuery := prometheus.CPUUsageQuery(node, labels)
	cpuResult, err := i.prometheusClient.Query(ctx, cpuQuery, now)
	if err != nil {
		return nil, fmt.Errorf("查询CPU使用率失败: %w", err)
	}
	cpuValue, err := prometheus.ParseValue(cpuResult)
	if err != nil {
		return nil, fmt.Errorf("解析CPU使用率失败: %w", err)
	}
	metrics.CPUNow = math.Round(cpuValue*100) / 100

	// 采集24小时前的CPU使用率
	cpuOffsetResult, err := i.prometheusClient.Query(ctx, cpuQuery, dayAgo)
	if err != nil {
		return nil, fmt.Errorf("查询24小时前CPU使用率失败: %w", err)
	}
	cpuOffsetValue, err := prometheus.ParseValue(cpuOffsetResult)
	if err != nil {
		return nil, fmt.Errorf("解析24小时前CPU使用率失败: %w", err)
	}
	metrics.CPUOffset = math.Round(cpuOffsetValue*100) / 100
	metrics.CPURate = math.Round((metrics.CPUNow-metrics.CPUOffset)*100) / 100

	// 采集内存使用率
	memoryQuery := prometheus.MemoryUsageQuery(node, labels)
	memoryResult, err := i.prometheusClient.Query(ctx, memoryQuery, now)
	if err != nil {
		return nil, fmt.Errorf("查询内存使用率失败: %w", err)
	}
	memoryValue, err := prometheus.ParseValue(memoryResult)
	if err != nil {
		return nil, fmt.Errorf("解析内存使用率失败: %w", err)
	}
	metrics.MemNow = math.Round(memoryValue*100) / 100

	// 采集24小时前的内存使用率
	memoryOffsetResult, err := i.prometheusClient.Query(ctx, memoryQuery, dayAgo)
	if err != nil {
		return nil, fmt.Errorf("查询24小时前内存使用率失败: %w", err)
	}
	memoryOffsetValue, err := prometheus.ParseValue(memoryOffsetResult)
	if err != nil {
		return nil, fmt.Errorf("解析24小时前内存使用率失败: %w", err)
	}
	metrics.MemOffset = math.Round(memoryOffsetValue*100) / 100
	metrics.MemRate = math.Round((metrics.MemNow-metrics.MemOffset)*100) / 100

	// 采集硬盘使用率 (取使用率最大的分区)
	diskQuery := prometheus.DiskUsageQuery(node, labels)
	diskResult, err := i.prometheusClient.Query(ctx, diskQuery, now)
	if err != nil {
		return nil, fmt.Errorf("查询硬盘使用率失败: %w", err)
	}
	diskValue, err := prometheus.ParseValue(diskResult)
	if err != nil {
		return nil, fmt.Errorf("解析硬盘使用率失败: %w", err)
	}
	metrics.DiskNow = math.Round(diskValue*100) / 100

	// 采集24小时前的硬盘使用率
	diskOffsetResult, err := i.prometheusClient.Query(ctx, diskQuery, dayAgo)
	if err != nil {
		return nil, fmt.Errorf("查询24小时前硬盘使用率失败: %w", err)
	}
	diskOffsetValue, err := prometheus.ParseValue(diskOffsetResult)
	if err != nil {
		return nil, fmt.Errorf("解析24小时前硬盘使用率失败: %w", err)
	}
	metrics.DiskOffset = math.Round(diskOffsetValue*100) / 100
	metrics.DiskRate = math.Round((metrics.DiskNow-metrics.DiskOffset)*100) / 100

	// 采集inode使用率
	inodeQuery := prometheus.InodeUsageQuery(node, labels)
	inodeResult, err := i.prometheusClient.Query(ctx, inodeQuery, now)
	if err != nil {
		return nil, fmt.Errorf("查询inode使用率失败: %w", err)
	}
	inodeValue, err := prometheus.ParseValue(inodeResult)
	if err != nil {
		return nil, fmt.Errorf("解析inode使用率失败: %w", err)
	}
	metrics.InodeNow = math.Round(inodeValue*100) / 100

	// 采集24小时前的inode使用率
	inodeOffsetResult, err := i.prometheusClient.Query(ctx, inodeQuery, dayAgo)
	if err != nil {
		return nil, fmt.Errorf("查询24小时前inode使用率失败: %w", err)
	}
	inodeOffsetValue, err := prometheus.ParseValue(inodeOffsetResult)
	if err != nil {
		return nil, fmt.Errorf("解析24小时前inode使用率失败: %w", err)
	}
	metrics.InodeOffset = math.Round(inodeOffsetValue*100) / 100
	metrics.InodeRate = math.Round((metrics.InodeNow-metrics.InodeOffset)*100) / 100

	return metrics, nil
}
