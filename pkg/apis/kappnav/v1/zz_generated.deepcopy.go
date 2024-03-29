// +build !ignore_autogenerated

// Code generated by operator-sdk. DO NOT EDIT.

package v1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Environment) DeepCopyInto(out *Environment) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Environment.
func (in *Environment) DeepCopy() *Environment {
	if in == nil {
		return nil
	}
	out := new(Environment)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Kappnav) DeepCopyInto(out *Kappnav) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Kappnav.
func (in *Kappnav) DeepCopy() *Kappnav {
	if in == nil {
		return nil
	}
	out := new(Kappnav)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Kappnav) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KappnavContainerConfiguration) DeepCopyInto(out *KappnavContainerConfiguration) {
	*out = *in
	if in.Resources != nil {
		in, out := &in.Resources, &out.Resources
		*out = new(KappnavResourceConstraints)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KappnavContainerConfiguration.
func (in *KappnavContainerConfiguration) DeepCopy() *KappnavContainerConfiguration {
	if in == nil {
		return nil
	}
	out := new(KappnavContainerConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KappnavImageConfiguration) DeepCopyInto(out *KappnavImageConfiguration) {
	*out = *in
	if in.PullSecrets != nil {
		in, out := &in.PullSecrets, &out.PullSecrets
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KappnavImageConfiguration.
func (in *KappnavImageConfiguration) DeepCopy() *KappnavImageConfiguration {
	if in == nil {
		return nil
	}
	out := new(KappnavImageConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KappnavList) DeepCopyInto(out *KappnavList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Kappnav, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KappnavList.
func (in *KappnavList) DeepCopy() *KappnavList {
	if in == nil {
		return nil
	}
	out := new(KappnavList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *KappnavList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KappnavResourceConstraints) DeepCopyInto(out *KappnavResourceConstraints) {
	*out = *in
	if in.Requests != nil {
		in, out := &in.Requests, &out.Requests
		*out = new(Resources)
		**out = **in
	}
	if in.Limits != nil {
		in, out := &in.Limits, &out.Limits
		*out = new(Resources)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KappnavResourceConstraints.
func (in *KappnavResourceConstraints) DeepCopy() *KappnavResourceConstraints {
	if in == nil {
		return nil
	}
	out := new(KappnavResourceConstraints)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KappnavSpec) DeepCopyInto(out *KappnavSpec) {
	*out = *in
	if in.AppNavAPI != nil {
		in, out := &in.AppNavAPI, &out.AppNavAPI
		*out = new(KappnavContainerConfiguration)
		(*in).DeepCopyInto(*out)
	}
	if in.AppNavController != nil {
		in, out := &in.AppNavController, &out.AppNavController
		*out = new(KappnavContainerConfiguration)
		(*in).DeepCopyInto(*out)
	}
	if in.AppNavUI != nil {
		in, out := &in.AppNavUI, &out.AppNavUI
		*out = new(KappnavContainerConfiguration)
		(*in).DeepCopyInto(*out)
	}
	if in.ExtensionContainers != nil {
		in, out := &in.ExtensionContainers, &out.ExtensionContainers
		*out = make(map[string]*KappnavContainerConfiguration, len(*in))
		for key, val := range *in {
			var outVal *KappnavContainerConfiguration
			if val == nil {
				(*out)[key] = nil
			} else {
				in, out := &val, &outVal
				*out = new(KappnavContainerConfiguration)
				(*in).DeepCopyInto(*out)
			}
			(*out)[key] = outVal
		}
	}
	if in.Image != nil {
		in, out := &in.Image, &out.Image
		*out = new(KappnavImageConfiguration)
		(*in).DeepCopyInto(*out)
	}
	if in.Env != nil {
		in, out := &in.Env, &out.Env
		*out = new(Environment)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KappnavSpec.
func (in *KappnavSpec) DeepCopy() *KappnavSpec {
	if in == nil {
		return nil
	}
	out := new(KappnavSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KappnavStatus) DeepCopyInto(out *KappnavStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]StatusCondition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KappnavStatus.
func (in *KappnavStatus) DeepCopy() *KappnavStatus {
	if in == nil {
		return nil
	}
	out := new(KappnavStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Resources) DeepCopyInto(out *Resources) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Resources.
func (in *Resources) DeepCopy() *Resources {
	if in == nil {
		return nil
	}
	out := new(Resources)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StatusCondition) DeepCopyInto(out *StatusCondition) {
	*out = *in
	if in.LastTransitionTime != nil {
		in, out := &in.LastTransitionTime, &out.LastTransitionTime
		*out = (*in).DeepCopy()
	}
	in.LastUpdateTime.DeepCopyInto(&out.LastUpdateTime)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StatusCondition.
func (in *StatusCondition) DeepCopy() *StatusCondition {
	if in == nil {
		return nil
	}
	out := new(StatusCondition)
	in.DeepCopyInto(out)
	return out
}
