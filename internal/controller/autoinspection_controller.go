/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	devopsv1 "github.com/rxg456/auto-inspection-operator/api/v1"
	"github.com/rxg456/auto-inspection-operator/internal/controller/inspection"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AutoInspectionReconciler reconciles a AutoInspection object
type AutoInspectionReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=devops.rxg98.cn,resources=autoinspections,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=devops.rxg98.cn,resources=autoinspections/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=devops.rxg98.cn,resources=autoinspections/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *AutoInspectionReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	logger.Info("Reconciling AutoInspection", "namespacedName", req.NamespacedName)

	// 获取AutoInspection对象
	var app devopsv1.AutoInspection
	if err := r.Get(ctx, req.NamespacedName, &app); err != nil {
		// 如果找不到资源，可能已被删除
		logger.Error(err, "Unable to fetch AutoInspection")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// 检查每个任务是否应该执行
	shouldRun := false
	var jobToRun *devopsv1.Job
	for i, job := range app.Spec.Jobs {
		canRun, err := inspection.ShouldRunNow(job, app.Status.LastInspectionTime)
		if err != nil {
			logger.Error(err, "检查任务执行时间失败", "job", job.Name)
			continue
		}

		if canRun {
			shouldRun = true
			jobToRun = &app.Spec.Jobs[i]
			break
		}
	}

	// 如果没有任务需要执行，则计算下一次检查时间
	if !shouldRun || jobToRun == nil {
		// 默认10分钟后重新检查
		return ctrl.Result{RequeueAfter: 10 * time.Minute}, nil
	}

	logger.Info("开始执行巡检任务", "job", jobToRun.Name)

	inspector, err := inspection.NewInspector(&app)
	if err != nil {
		logger.Error(err, "创建巡检器失败")
		return ctrl.Result{RequeueAfter: 5 * time.Minute}, nil
	}

	// 执行巡检
	if err := inspector.RunInspection(ctx); err != nil {
		logger.Error(err, "执行巡检失败")
		return ctrl.Result{RequeueAfter: 5 * time.Minute}, nil
	}

	// 更新状态
	now := metav1.NewTime(time.Now())
	app.Status.LastInspectionTime = &now

	if err := r.Status().Update(ctx, &app); err != nil {
		logger.Error(err, "更新AutoInspection状态失败")
		return ctrl.Result{RequeueAfter: 1 * time.Minute}, nil
	}

	// 计算下一次运行时间
	nextRun, err := inspection.GetNextRunTime(*jobToRun, app.Status.LastInspectionTime)
	if err != nil {
		logger.Error(err, "计算下一次运行时间失败")
		return ctrl.Result{RequeueAfter: 10 * time.Minute}, nil
	}

	// 计算需要等待的时间
	waitDuration := time.Until(nextRun)
	if waitDuration < 0 {
		waitDuration = 1 * time.Minute // 如果下次运行时间已过，1分钟后重试
	}

	logger.Info("巡检任务执行完成，等待下次执行", "nextRun", nextRun, "waitDuration", waitDuration)
	return ctrl.Result{RequeueAfter: waitDuration}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AutoInspectionReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&devopsv1.AutoInspection{}).
		Named("autoinspection").
		Complete(r)
}
