/*
Copyright 2021.

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

package controllers

import (
	"context"
	"time"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	appv1 "github.com/vbouchaud/factory-operator/api/v1"
)

// ProjectReconciler reconciles a Project object
type ProjectReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

const projectFinalizer = "app.heidrun.bouchaud.org/finalizer"

const (
	projectConditionInitialized = "Initialized"
	projectConditionConfigured  = "Configured"
)

//+kubebuilder:rbac:groups=app.heidrun.bouchaud.org,resources=projects,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=app.heidrun.bouchaud.org,resources=projects/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=app.heidrun.bouchaud.org,resources=projects/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Project object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.9.2/pkg/reconcile
func (r *ProjectReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	logger.Info("Reconciling Project.")

	// Fetch the Project instance
	project := &appv1.Project{}
	err := r.Get(ctx, req.NamespacedName, project)
	if err != nil {
		if errors.IsNotFound(err) {
			logger.Info("Project resource not found. Ignoring since object must be deleted.")
			return ctrl.Result{}, nil
		}
		logger.Error(err, "Failed to get Project.")
		return ctrl.Result{}, err
	}

	isProjectMarkedToBeDeleted := project.GetDeletionTimestamp() != nil
	if isProjectMarkedToBeDeleted {
		if controllerutil.ContainsFinalizer(project, projectFinalizer) {
			if err := r.finalizeProject(logger, project); err != nil {
				return ctrl.Result{}, err
			}

			controllerutil.RemoveFinalizer(project, projectFinalizer)
			err := r.Update(ctx, project)
			if err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	if !controllerutil.ContainsFinalizer(project, projectFinalizer) {
		controllerutil.AddFinalizer(project, projectFinalizer)
		addCondition(logger, project, projectConditionInitialized, metav1.ConditionTrue)
		addCondition(logger, project, projectConditionConfigured, metav1.ConditionFalse)
		err = r.Update(ctx, project)
		if err != nil {
			return ctrl.Result{}, err
		}
	}

	// Here handle the creation, update, etc. of external resources needed for this project

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ProjectReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appv1.Project{}).
		Complete(r)
}

func (r *ProjectReconciler) finalizeProject(l logr.Logger, p *appv1.Project) error {
	// TODO(user): Add the cleanup steps that the operator
	// needs to do before the CR can be deleted. Examples
	// of finalizers include performing backups and deleting
	// resources that are not owned by this CR, like a PVC.
	l.Info("Successfully finalized project.")
	return nil
}

func addCondition(l logr.Logger, p *appv1.Project, c string, s metav1.ConditionStatus) {
	l.Info("Setting condition %s to %s.", c, s)

	meta.SetStatusCondition(&p.Status.Conditions, metav1.Condition{
		Type:   c,
		Status: s,
		LastTransitionTime: metav1.Time{
			Time: time.Now(),
		},
	})
}
