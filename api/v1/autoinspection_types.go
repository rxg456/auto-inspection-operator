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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// AutoInspectionSpec defines the desired state of AutoInspection.
type AutoInspectionSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// 定义巡检任务和调度时间
	Jobs []Job `json:"jobs"`

	// 定义邮件服务器配置
	SMTP SMTP `json:"smtp"`

	// 定义接收通知的邮件地址
	NotifyTo []string `json:"notifyTo"`

	// 定义Prometheus API地址
	PrometheusURL string `json:"prometheusURL"`

	// 定义巡检对象（业务的主机）
	InspectionObject InspectionObject `json:"inspectionObject"`
}

// Job定义巡检任务和调度
type Job struct {
	// 任务名称
	Name string `json:"name"`
	// 任务Cron表达式
	Schedule string `json:"schedule"`
}

// SMTP定义邮件服务器配置
type SMTP struct {
	// 服务器地址
	Server string `json:"server"`
	// 服务器端口
	Port int `json:"port"`
	// 发件人地址
	From string `json:"from"`
	// 用户名
	Username string `json:"username"`
	// 密码
	Password string `json:"password"`
}

// 巡检对象
type InspectionObject struct {
	// 业务名称
	Business string `json:"business"`
	// 主机名称
	Hosts Hosts `json:"hosts"`
}

// 主机信息
type Hosts struct {
	// 标签
	Labels map[string]string `json:"labels,omitempty"`
	// 主机列表
	Nodes []string `json:"nodes,omitempty"`
}

// AutoInspectionStatus defines the observed state of AutoInspection.
type AutoInspectionStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// 最后巡检时间
	LastInspectionTime *metav1.Time `json:"lastInspectionTime,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:printcolumn:name="Last Inspection",type="date",JSONPath=".status.lastInspectionTime",description="Last inspection time"

// AutoInspection is the Schema for the autoinspections API.
type AutoInspection struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AutoInspectionSpec   `json:"spec,omitempty"`
	Status AutoInspectionStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// AutoInspectionList contains a list of AutoInspection.
type AutoInspectionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AutoInspection `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AutoInspection{}, &AutoInspectionList{})
}
