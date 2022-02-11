/*
Copyright 2022.

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
	"github.com/go-logr/logr"
	"github.com/vbouchaud/factory-operator/internal/ldap"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	appv1 "github.com/vbouchaud/factory-operator/api/v1"
)

// TeamReconciler reconciles a Team object
type TeamReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Ldap   *ldap.Client
}

const teamFinalizer = "app.heidrun.bouchaud.org/team-finalizer"

//+kubebuilder:rbac:groups=app.heidrun.bouchaud.org,resources=teams,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=app.heidrun.bouchaud.org,resources=teams/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=app.heidrun.bouchaud.org,resources=teams/finalizers,verbs=update

func (r *TeamReconciler) addCondition(l logr.Logger, b *appv1.Team, t string, s metav1.ConditionStatus) {
	l.Info("Setting condition", "status", t, "condition", s)

	meta.SetStatusCondition(&b.Status.Conditions, metav1.Condition{
		Type:   t,
		Status: s,
	})
}

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Team object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *TeamReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	logger.Info("Reconciling Team.")

	// Fetch the Team instance
	team := &appv1.Team{}
	err := r.Get(ctx, req.NamespacedName, team)
	if err != nil {
		if errors.IsNotFound(err) {
			logger.Info("Team resource not found. Ignoring since object must be deleted.")
			return ctrl.Result{}, nil
		}
		logger.Error(err, "Failed to get Team.")
		return ctrl.Result{}, err
	}

	// Team deletion
	isTeamMarkedToBeDeleted := team.GetDeletionTimestamp() != nil
	if isTeamMarkedToBeDeleted {
		if controllerutil.ContainsFinalizer(team, teamFinalizer) {
			if err := r.Ldap.DeleteGroup(team.Name); err != nil {
				if ldap.IsNotFound(err) {
					logger.Info("Ldap group not found.", "ldap-group", team.Name)
				} else {
					logger.Error(err, "Error while removing group.", "ldap-group", team.Name)
					return ctrl.Result{}, err
				}
			}

			controllerutil.RemoveFinalizer(team, teamFinalizer)
			if err = r.Update(ctx, team); err != nil {
				logger.Error(err, "Failed to remove finalizer.", "ldap-group", team.Name)
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	// Team Initialization
	if !controllerutil.ContainsFinalizer(team, teamFinalizer) {
		controllerutil.AddFinalizer(team, teamFinalizer)
		r.addCondition(logger, team, conditionInitialized, metav1.ConditionTrue)
		if err = r.Update(ctx, team); err != nil {
			logger.Error(err, "Failed to initialize Team status.", "ldap-group", team.Name)
			return ctrl.Result{}, err
		}
	}

	// Team update
	var changed bool
	team.Status.DistinguishedName, err, changed = r.Ldap.ReconcileGroup(team.Name, team.Spec.Comment, team.Spec.Subjects)
	if err != nil {
		logger.Error(err, "Failed to crupdate Team resource.", "ldap-group", team.Name)
		return ctrl.Result{}, err
	}

	if changed {
		r.addCondition(logger, team, conditionConfigured, metav1.ConditionFalse)
		if err = r.Update(ctx, team); err != nil {
			logger.Error(err, "Failed to update Team status.", "ldap-group", team.Name)
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *TeamReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appv1.Team{}).
		Complete(r)
}
