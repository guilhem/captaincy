// +build !ignore_autogenerated

/*
Copyright 2017 The etcd-operator Authors

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

// This file was autogenerated by deepcopy-gen. Do not edit it manually!

package v1beta2

import (
	v1 "k8s.io/api/core/v1"
	conversion "k8s.io/apimachinery/pkg/conversion"
	runtime "k8s.io/apimachinery/pkg/runtime"
	reflect "reflect"
)

// GetGeneratedDeepCopyFuncs returns the generated funcs, since we aren't registering them.
//
// Deprecated: deepcopy registration will go away when static deepcopy is fully implemented.
func GetGeneratedDeepCopyFuncs() []conversion.GeneratedDeepCopyFunc {
	return []conversion.GeneratedDeepCopyFunc{
		{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*BackupCRStatus).DeepCopyInto(out.(*BackupCRStatus))
			return nil
		}, InType: reflect.TypeOf(&BackupCRStatus{})},
		{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*BackupSpec).DeepCopyInto(out.(*BackupSpec))
			return nil
		}, InType: reflect.TypeOf(&BackupSpec{})},
		{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*BackupStorageSource).DeepCopyInto(out.(*BackupStorageSource))
			return nil
		}, InType: reflect.TypeOf(&BackupStorageSource{})},
		{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*ClusterCondition).DeepCopyInto(out.(*ClusterCondition))
			return nil
		}, InType: reflect.TypeOf(&ClusterCondition{})},
		{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*ClusterSpec).DeepCopyInto(out.(*ClusterSpec))
			return nil
		}, InType: reflect.TypeOf(&ClusterSpec{})},
		{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*ClusterStatus).DeepCopyInto(out.(*ClusterStatus))
			return nil
		}, InType: reflect.TypeOf(&ClusterStatus{})},
		{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*EtcdBackup).DeepCopyInto(out.(*EtcdBackup))
			return nil
		}, InType: reflect.TypeOf(&EtcdBackup{})},
		{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*EtcdBackupList).DeepCopyInto(out.(*EtcdBackupList))
			return nil
		}, InType: reflect.TypeOf(&EtcdBackupList{})},
		{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*EtcdCluster).DeepCopyInto(out.(*EtcdCluster))
			return nil
		}, InType: reflect.TypeOf(&EtcdCluster{})},
		{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*EtcdClusterList).DeepCopyInto(out.(*EtcdClusterList))
			return nil
		}, InType: reflect.TypeOf(&EtcdClusterList{})},
		{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*EtcdRestore).DeepCopyInto(out.(*EtcdRestore))
			return nil
		}, InType: reflect.TypeOf(&EtcdRestore{})},
		{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*EtcdRestoreList).DeepCopyInto(out.(*EtcdRestoreList))
			return nil
		}, InType: reflect.TypeOf(&EtcdRestoreList{})},
		{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*MemberSecret).DeepCopyInto(out.(*MemberSecret))
			return nil
		}, InType: reflect.TypeOf(&MemberSecret{})},
		{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*MembersStatus).DeepCopyInto(out.(*MembersStatus))
			return nil
		}, InType: reflect.TypeOf(&MembersStatus{})},
		{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*PodPolicy).DeepCopyInto(out.(*PodPolicy))
			return nil
		}, InType: reflect.TypeOf(&PodPolicy{})},
		{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*RestoreSource).DeepCopyInto(out.(*RestoreSource))
			return nil
		}, InType: reflect.TypeOf(&RestoreSource{})},
		{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*RestoreSpec).DeepCopyInto(out.(*RestoreSpec))
			return nil
		}, InType: reflect.TypeOf(&RestoreSpec{})},
		{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*RestoreStatus).DeepCopyInto(out.(*RestoreStatus))
			return nil
		}, InType: reflect.TypeOf(&RestoreStatus{})},
		{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*S3RestoreSource).DeepCopyInto(out.(*S3RestoreSource))
			return nil
		}, InType: reflect.TypeOf(&S3RestoreSource{})},
		{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*S3Source).DeepCopyInto(out.(*S3Source))
			return nil
		}, InType: reflect.TypeOf(&S3Source{})},
		{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*SelfHostedPolicy).DeepCopyInto(out.(*SelfHostedPolicy))
			return nil
		}, InType: reflect.TypeOf(&SelfHostedPolicy{})},
		{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*StaticTLS).DeepCopyInto(out.(*StaticTLS))
			return nil
		}, InType: reflect.TypeOf(&StaticTLS{})},
		{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*TLSPolicy).DeepCopyInto(out.(*TLSPolicy))
			return nil
		}, InType: reflect.TypeOf(&TLSPolicy{})},
	}
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BackupCRStatus) DeepCopyInto(out *BackupCRStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BackupCRStatus.
func (in *BackupCRStatus) DeepCopy() *BackupCRStatus {
	if in == nil {
		return nil
	}
	out := new(BackupCRStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BackupSpec) DeepCopyInto(out *BackupSpec) {
	*out = *in
	in.BackupStorageSource.DeepCopyInto(&out.BackupStorageSource)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BackupSpec.
func (in *BackupSpec) DeepCopy() *BackupSpec {
	if in == nil {
		return nil
	}
	out := new(BackupSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BackupStorageSource) DeepCopyInto(out *BackupStorageSource) {
	*out = *in
	if in.S3 != nil {
		in, out := &in.S3, &out.S3
		if *in == nil {
			*out = nil
		} else {
			*out = new(S3Source)
			**out = **in
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BackupStorageSource.
func (in *BackupStorageSource) DeepCopy() *BackupStorageSource {
	if in == nil {
		return nil
	}
	out := new(BackupStorageSource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterCondition) DeepCopyInto(out *ClusterCondition) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterCondition.
func (in *ClusterCondition) DeepCopy() *ClusterCondition {
	if in == nil {
		return nil
	}
	out := new(ClusterCondition)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterSpec) DeepCopyInto(out *ClusterSpec) {
	*out = *in
	if in.Pod != nil {
		in, out := &in.Pod, &out.Pod
		if *in == nil {
			*out = nil
		} else {
			*out = new(PodPolicy)
			(*in).DeepCopyInto(*out)
		}
	}
	if in.SelfHosted != nil {
		in, out := &in.SelfHosted, &out.SelfHosted
		if *in == nil {
			*out = nil
		} else {
			*out = new(SelfHostedPolicy)
			**out = **in
		}
	}
	if in.TLS != nil {
		in, out := &in.TLS, &out.TLS
		if *in == nil {
			*out = nil
		} else {
			*out = new(TLSPolicy)
			(*in).DeepCopyInto(*out)
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterSpec.
func (in *ClusterSpec) DeepCopy() *ClusterSpec {
	if in == nil {
		return nil
	}
	out := new(ClusterSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterStatus) DeepCopyInto(out *ClusterStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]ClusterCondition, len(*in))
		copy(*out, *in)
	}
	in.Members.DeepCopyInto(&out.Members)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterStatus.
func (in *ClusterStatus) DeepCopy() *ClusterStatus {
	if in == nil {
		return nil
	}
	out := new(ClusterStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EtcdBackup) DeepCopyInto(out *EtcdBackup) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EtcdBackup.
func (in *EtcdBackup) DeepCopy() *EtcdBackup {
	if in == nil {
		return nil
	}
	out := new(EtcdBackup)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *EtcdBackup) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	} else {
		return nil
	}
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EtcdBackupList) DeepCopyInto(out *EtcdBackupList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]EtcdBackup, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EtcdBackupList.
func (in *EtcdBackupList) DeepCopy() *EtcdBackupList {
	if in == nil {
		return nil
	}
	out := new(EtcdBackupList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *EtcdBackupList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	} else {
		return nil
	}
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EtcdCluster) DeepCopyInto(out *EtcdCluster) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EtcdCluster.
func (in *EtcdCluster) DeepCopy() *EtcdCluster {
	if in == nil {
		return nil
	}
	out := new(EtcdCluster)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *EtcdCluster) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	} else {
		return nil
	}
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EtcdClusterList) DeepCopyInto(out *EtcdClusterList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]EtcdCluster, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EtcdClusterList.
func (in *EtcdClusterList) DeepCopy() *EtcdClusterList {
	if in == nil {
		return nil
	}
	out := new(EtcdClusterList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *EtcdClusterList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	} else {
		return nil
	}
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EtcdRestore) DeepCopyInto(out *EtcdRestore) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EtcdRestore.
func (in *EtcdRestore) DeepCopy() *EtcdRestore {
	if in == nil {
		return nil
	}
	out := new(EtcdRestore)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *EtcdRestore) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	} else {
		return nil
	}
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EtcdRestoreList) DeepCopyInto(out *EtcdRestoreList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]EtcdRestore, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EtcdRestoreList.
func (in *EtcdRestoreList) DeepCopy() *EtcdRestoreList {
	if in == nil {
		return nil
	}
	out := new(EtcdRestoreList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *EtcdRestoreList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	} else {
		return nil
	}
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MemberSecret) DeepCopyInto(out *MemberSecret) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MemberSecret.
func (in *MemberSecret) DeepCopy() *MemberSecret {
	if in == nil {
		return nil
	}
	out := new(MemberSecret)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MembersStatus) DeepCopyInto(out *MembersStatus) {
	*out = *in
	if in.Ready != nil {
		in, out := &in.Ready, &out.Ready
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Unready != nil {
		in, out := &in.Unready, &out.Unready
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MembersStatus.
func (in *MembersStatus) DeepCopy() *MembersStatus {
	if in == nil {
		return nil
	}
	out := new(MembersStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PodPolicy) DeepCopyInto(out *PodPolicy) {
	*out = *in
	if in.Labels != nil {
		in, out := &in.Labels, &out.Labels
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.NodeSelector != nil {
		in, out := &in.NodeSelector, &out.NodeSelector
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	in.Resources.DeepCopyInto(&out.Resources)
	if in.Tolerations != nil {
		in, out := &in.Tolerations, &out.Tolerations
		*out = make([]v1.Toleration, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.EtcdEnv != nil {
		in, out := &in.EtcdEnv, &out.EtcdEnv
		*out = make([]v1.EnvVar, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.AutomountServiceAccountToken != nil {
		in, out := &in.AutomountServiceAccountToken, &out.AutomountServiceAccountToken
		if *in == nil {
			*out = nil
		} else {
			*out = new(bool)
			**out = **in
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PodPolicy.
func (in *PodPolicy) DeepCopy() *PodPolicy {
	if in == nil {
		return nil
	}
	out := new(PodPolicy)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RestoreSource) DeepCopyInto(out *RestoreSource) {
	*out = *in
	if in.S3 != nil {
		in, out := &in.S3, &out.S3
		if *in == nil {
			*out = nil
		} else {
			*out = new(S3RestoreSource)
			**out = **in
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RestoreSource.
func (in *RestoreSource) DeepCopy() *RestoreSource {
	if in == nil {
		return nil
	}
	out := new(RestoreSource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RestoreSpec) DeepCopyInto(out *RestoreSpec) {
	*out = *in
	in.ClusterSpec.DeepCopyInto(&out.ClusterSpec)
	in.RestoreSource.DeepCopyInto(&out.RestoreSource)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RestoreSpec.
func (in *RestoreSpec) DeepCopy() *RestoreSpec {
	if in == nil {
		return nil
	}
	out := new(RestoreSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RestoreStatus) DeepCopyInto(out *RestoreStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RestoreStatus.
func (in *RestoreStatus) DeepCopy() *RestoreStatus {
	if in == nil {
		return nil
	}
	out := new(RestoreStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *S3RestoreSource) DeepCopyInto(out *S3RestoreSource) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new S3RestoreSource.
func (in *S3RestoreSource) DeepCopy() *S3RestoreSource {
	if in == nil {
		return nil
	}
	out := new(S3RestoreSource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *S3Source) DeepCopyInto(out *S3Source) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new S3Source.
func (in *S3Source) DeepCopy() *S3Source {
	if in == nil {
		return nil
	}
	out := new(S3Source)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SelfHostedPolicy) DeepCopyInto(out *SelfHostedPolicy) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SelfHostedPolicy.
func (in *SelfHostedPolicy) DeepCopy() *SelfHostedPolicy {
	if in == nil {
		return nil
	}
	out := new(SelfHostedPolicy)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StaticTLS) DeepCopyInto(out *StaticTLS) {
	*out = *in
	if in.Member != nil {
		in, out := &in.Member, &out.Member
		if *in == nil {
			*out = nil
		} else {
			*out = new(MemberSecret)
			**out = **in
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StaticTLS.
func (in *StaticTLS) DeepCopy() *StaticTLS {
	if in == nil {
		return nil
	}
	out := new(StaticTLS)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TLSPolicy) DeepCopyInto(out *TLSPolicy) {
	*out = *in
	if in.Static != nil {
		in, out := &in.Static, &out.Static
		if *in == nil {
			*out = nil
		} else {
			*out = new(StaticTLS)
			(*in).DeepCopyInto(*out)
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TLSPolicy.
func (in *TLSPolicy) DeepCopy() *TLSPolicy {
	if in == nil {
		return nil
	}
	out := new(TLSPolicy)
	in.DeepCopyInto(out)
	return out
}
