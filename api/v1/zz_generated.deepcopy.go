//go:build !ignore_autogenerated

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

// Code generated by controller-gen. DO NOT EDIT.

package v1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AutoInspection) DeepCopyInto(out *AutoInspection) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AutoInspection.
func (in *AutoInspection) DeepCopy() *AutoInspection {
	if in == nil {
		return nil
	}
	out := new(AutoInspection)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *AutoInspection) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AutoInspectionList) DeepCopyInto(out *AutoInspectionList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]AutoInspection, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AutoInspectionList.
func (in *AutoInspectionList) DeepCopy() *AutoInspectionList {
	if in == nil {
		return nil
	}
	out := new(AutoInspectionList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *AutoInspectionList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AutoInspectionSpec) DeepCopyInto(out *AutoInspectionSpec) {
	*out = *in
	if in.Jobs != nil {
		in, out := &in.Jobs, &out.Jobs
		*out = make([]Job, len(*in))
		copy(*out, *in)
	}
	out.SMTP = in.SMTP
	if in.NotifyTo != nil {
		in, out := &in.NotifyTo, &out.NotifyTo
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	in.InspectionObject.DeepCopyInto(&out.InspectionObject)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AutoInspectionSpec.
func (in *AutoInspectionSpec) DeepCopy() *AutoInspectionSpec {
	if in == nil {
		return nil
	}
	out := new(AutoInspectionSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AutoInspectionStatus) DeepCopyInto(out *AutoInspectionStatus) {
	*out = *in
	if in.LastInspectionTime != nil {
		in, out := &in.LastInspectionTime, &out.LastInspectionTime
		*out = (*in).DeepCopy()
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AutoInspectionStatus.
func (in *AutoInspectionStatus) DeepCopy() *AutoInspectionStatus {
	if in == nil {
		return nil
	}
	out := new(AutoInspectionStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Hosts) DeepCopyInto(out *Hosts) {
	*out = *in
	if in.Labels != nil {
		in, out := &in.Labels, &out.Labels
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Nodes != nil {
		in, out := &in.Nodes, &out.Nodes
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Hosts.
func (in *Hosts) DeepCopy() *Hosts {
	if in == nil {
		return nil
	}
	out := new(Hosts)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InspectionObject) DeepCopyInto(out *InspectionObject) {
	*out = *in
	in.Hosts.DeepCopyInto(&out.Hosts)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InspectionObject.
func (in *InspectionObject) DeepCopy() *InspectionObject {
	if in == nil {
		return nil
	}
	out := new(InspectionObject)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Job) DeepCopyInto(out *Job) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Job.
func (in *Job) DeepCopy() *Job {
	if in == nil {
		return nil
	}
	out := new(Job)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SMTP) DeepCopyInto(out *SMTP) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SMTP.
func (in *SMTP) DeepCopy() *SMTP {
	if in == nil {
		return nil
	}
	out := new(SMTP)
	in.DeepCopyInto(out)
	return out
}
