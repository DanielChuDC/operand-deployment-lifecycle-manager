//
// Copyright 2020 IBM Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package operandbindinfo

import (
	"context"
	"reflect"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/klog"

	operatorv1alpha1 "github.com/IBM/operand-deployment-lifecycle-manager/pkg/apis/operator/v1alpha1"
)

func (r *ReconcileOperandBindInfo) updateBindInfoPhase(cr *operatorv1alpha1.OperandBindInfo, phase operatorv1alpha1.BindInfoPhase, requestNamespaces []operatorv1alpha1.ReconcileRequest) error {
	if err := wait.PollImmediate(time.Second*20, time.Minute*10, func() (done bool, err error) {
		bindInfoInstance, err := r.getBindInfoInstance(cr.Name, cr.Namespace)
		if err != nil {
			return false, err
		}
		var requestNsList []string
		for _, ns := range requestNamespaces {
			if ns.Namespace == bindInfoInstance.Namespace {
				continue
			}
			requestNsList = append(requestNsList, ns.Namespace)
		}
		requestNsList = unique(requestNsList)
		if bindInfoInstance.Status.Phase == phase && reflect.DeepEqual(requestNsList, bindInfoInstance.Status.RequestNamespaces) {
			return true, nil
		}
		if len(requestNsList) != 0 {
			bindInfoInstance.Status.RequestNamespaces = requestNsList
		}
		bindInfoInstance.Status.Phase = phase
		if err := r.client.Status().Update(context.TODO(), bindInfoInstance); err != nil {
			klog.V(3).Info("Waiting for OperandBindInfo instance status ready ...")
			return false, nil
		}
		return true, nil
	}); err != nil {
		return err
	}
	return nil
}

func unique(stringSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
