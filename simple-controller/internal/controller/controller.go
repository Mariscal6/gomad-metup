/*
Copyright 2018 The Kubernetes Authors.

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
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// ReconcileDeployment reconciles Deployments
type ReconcileDeployment struct {
	// client can be used to retrieve objects from the APIServer.
	client client.Client
}

func NewReconciler(client client.Client) *ReconcileDeployment {
	return &ReconcileDeployment{client: client}
}

// Implement reconcile.Reconciler so the controller can reconcile objects
var _ reconcile.Reconciler = &ReconcileDeployment{}

func (r *ReconcileDeployment) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	// set up a convenient log object so we don't have to type request over and over again
	log := log.FromContext(ctx)
	log.Info("ns: %s name: %s\n", request.Namespace, request.Name)
	// Fetch the Deployment from the cache
	dep := &appsv1.Deployment{}
	err := r.client.Get(ctx, request.NamespacedName, dep)
	if errors.IsNotFound(err) {
		err = DeleteService(ctx, r.client, dep)
		return reconcile.Result{}, nil
	}
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("could not fetch Deployment: %+v", err)
	}

	if dep.Labels != nil {
		svc, ok := dep.Labels["gomad"]
		if !ok || svc == "" {
			return reconcile.Result{}, nil
		}
		err = CreateSvc(ctx, r.client, dep)
		if err != nil {
			return reconcile.Result{}, fmt.Errorf("could not create Service: %+v", err)
		}
	}
	return reconcile.Result{}, nil
}

func CreateSvc(ctx context.Context, c client.Client, dep *appsv1.Deployment) error {
	svc := &corev1.Service{}
	err := c.Get(ctx, client.ObjectKey{Namespace: dep.Namespace, Name: dep.Name}, svc)
	if errors.IsNotFound(err) {
		err := c.Create(ctx, &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      dep.Name,
				Namespace: dep.Namespace,
				Labels:    dep.Labels,
			},
			Spec: corev1.ServiceSpec{
				Ports: []corev1.ServicePort{
					{
						Port:     80,
						Name:     "http",
						Protocol: corev1.ProtocolTCP,
						TargetPort: intstr.IntOrString{
							Type:   intstr.Int,
							IntVal: 80,
						},
					},
				},
				Selector: dep.Labels,
				Type:     corev1.ServiceTypeClusterIP,
			},
		})
		if err != nil {
			return fmt.Errorf("could not create Service: %+v", err)
		}
	}
	return nil
}
func DeleteService(ctx context.Context, c client.Client, dep *appsv1.Deployment) error {
	svc := &corev1.Service{}
	err := c.Get(ctx, client.ObjectKey{Namespace: dep.Namespace, Name: dep.Name}, svc)
	if errors.IsNotFound(err) {
		fmt.Printf("Service %s not found\n", dep.Name)
		return nil
	}
	if err != nil {
		return fmt.Errorf("could not fetch Service: %+v", err)
	}
	err = c.Delete(ctx, svc)
	if err != nil {
		return fmt.Errorf("could not delete Service: %+v", err)
	}
	return nil
}
