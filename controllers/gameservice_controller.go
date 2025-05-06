package controllers

import (
	"context"
	"fmt"

	gamesv1alpha1 "github.com/creativewild/gameservice-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	portStart = 30000
	portEnd   = 32767
)

type GameServiceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *GameServiceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	var gs gamesv1alpha1.GameService
	if err := r.Get(ctx, req.NamespacedName, &gs); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		logger.Error(err, "unable to fetch GameService")
		return ctrl.Result{}, err
	}

	if gs.Status.Assigned {
		return ctrl.Result{}, nil
	}

	var gsList gamesv1alpha1.GameServiceList
	if err := r.List(ctx, &gsList); err != nil {
		logger.Error(err, "failed to list GameServices")
		return ctrl.Result{}, err
	}

	usedPorts := make(map[int]bool)
	for _, item := range gsList.Items {
		if item.Status.Assigned {
			usedPorts[item.Spec.MappedPort] = true
		}
	}

	var assignedPort int
	for port := portStart; port <= portEnd; port++ {
		if !usedPorts[port] {
			assignedPort = port
			break
		}
	}
	if assignedPort == 0 {
		return ctrl.Result{}, fmt.Errorf("no available ports in range %d–%d", portStart, portEnd)
	}

	// Assign values
	gs.Spec.MappedPort = assignedPort
	gs.Status.Assigned = true
	gs.Status.Endpoints.IPv4 = fmt.Sprintf("%s:%d", gs.Spec.SharedIPv4, assignedPort)
	gs.Status.Endpoints.IPv6 = gs.Spec.IPv6Address

	if err := r.Status().Update(ctx, &gs); err != nil {
		logger.Error(err, "failed to update GameService status")
		return ctrl.Result{}, err
	}

	logger.Info("Assigned GameService", "name", gs.Name, "port", assignedPort)
	return ctrl.Result{}, nil
}

func (r *GameServiceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&gamesv1alpha1.GameService{}).
		Complete(r)
}
