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

package main

import (
	"os"

	simplecon "github.com/Mariscal6/gomad-metup/simple-controller/internal/controller"
	appsv1 "k8s.io/api/apps/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

func init() {
	log.SetLogger(zap.New())
}

func main() {
	entryLog := log.Log.WithName("entrypoint")
	// Setup a Manager
	// Manager is required for creating a Controller and provides
	// the Controller shared dependencies such as clients, caches, schemes, etc.
	mgr, err := manager.New(config.GetConfigOrDie(), manager.Options{})
	if err != nil {
		entryLog.Error(err, "unable to set up overall controller manager")
		os.Exit(1)
	}

	// Setup a new controller to reconcile.Deployments
	// Controller implements a Kubernetes API by responding to events (object Create, Update, Delete)
	// and ensuring that the state specified in the Spec of the object matches the state of the system.
	// This is called a reconcile. If they do not match, the Controller will create / update / delete objects
	// as needed to make them match.
	entryLog.Info("Setting up controller")
	c, err := controller.New("gomad-simple-controller", mgr, controller.Options{
		Reconciler: simplecon.NewReconciler(mgr.GetClient()),
	})
	if err != nil {
		entryLog.Error(err, "unable to set up individual controller")
		os.Exit(1)
	}

	// Watch.Deployments and enqueue.Deployment object key
	if err := c.Watch(source.Kind(mgr.GetCache(), &appsv1.Deployment{}, &handler.TypedEnqueueRequestForObject[*appsv1.Deployment]{})); err != nil {
		entryLog.Error(err, "unable to watch.Deployments")
		os.Exit(1)
	}

	entryLog.Info("starting manager")
	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		entryLog.Error(err, "unable to run manager")
		os.Exit(1)
	}
}
