//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright The Kubernetes Authors.

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

// Code generated by deepcopy-gen. DO NOT EDIT.

package v1beta1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FooServer) DeepCopyInto(out *FooServer) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FooServer.
func (in *FooServer) DeepCopy() *FooServer {
	if in == nil {
		return nil
	}
	out := new(FooServer)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *FooServer) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FooServerDeploymentSpec) DeepCopyInto(out *FooServerDeploymentSpec) {
	*out = *in
	if in.PodReplicas != nil {
		in, out := &in.PodReplicas, &out.PodReplicas
		*out = new(int32)
		**out = **in
	}
	if in.PodContainerPort != nil {
		in, out := &in.PodContainerPort, &out.PodContainerPort
		*out = new(int32)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FooServerDeploymentSpec.
func (in *FooServerDeploymentSpec) DeepCopy() *FooServerDeploymentSpec {
	if in == nil {
		return nil
	}
	out := new(FooServerDeploymentSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FooServerList) DeepCopyInto(out *FooServerList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]FooServer, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FooServerList.
func (in *FooServerList) DeepCopy() *FooServerList {
	if in == nil {
		return nil
	}
	out := new(FooServerList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *FooServerList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FooServerSecretSpec) DeepCopyInto(out *FooServerSecretSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FooServerSecretSpec.
func (in *FooServerSecretSpec) DeepCopy() *FooServerSecretSpec {
	if in == nil {
		return nil
	}
	out := new(FooServerSecretSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FooServerServiceSpec) DeepCopyInto(out *FooServerServiceSpec) {
	*out = *in
	if in.ServicePort != nil {
		in, out := &in.ServicePort, &out.ServicePort
		*out = new(int)
		**out = **in
	}
	if in.ServiceTargetPort != nil {
		in, out := &in.ServiceTargetPort, &out.ServiceTargetPort
		*out = new(int)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FooServerServiceSpec.
func (in *FooServerServiceSpec) DeepCopy() *FooServerServiceSpec {
	if in == nil {
		return nil
	}
	out := new(FooServerServiceSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FooServerSpec) DeepCopyInto(out *FooServerSpec) {
	*out = *in
	in.DeploymentSpec.DeepCopyInto(&out.DeploymentSpec)
	in.ServiceSpec.DeepCopyInto(&out.ServiceSpec)
	out.SecretSpec = in.SecretSpec
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FooServerSpec.
func (in *FooServerSpec) DeepCopy() *FooServerSpec {
	if in == nil {
		return nil
	}
	out := new(FooServerSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FooServerStatus) DeepCopyInto(out *FooServerStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FooServerStatus.
func (in *FooServerStatus) DeepCopy() *FooServerStatus {
	if in == nil {
		return nil
	}
	out := new(FooServerStatus)
	in.DeepCopyInto(out)
	return out
}
