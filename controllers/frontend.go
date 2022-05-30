package controllers

import (
	"context"
	"strconv"
	"time"

	appsv1alpha1 "github.com/majguo/visitors-operator/api/v1alpha1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const frontendPort = 3000
const frontendImage = "majguo/visitors-webui:1.0.0"

func frontendDeploymentName(v *appsv1alpha1.VisitorsApp) string {
	return v.Name + "-frontend"
}

func frontendServiceName(v *appsv1alpha1.VisitorsApp) string {
	return v.Name + "-frontend-service"
}

func (r *VisitorsAppReconciler) frontendDeployment(v *appsv1alpha1.VisitorsApp) *appsv1.Deployment {
	labels := labels(v, "frontend")
	size := v.Spec.Size

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      frontendDeploymentName(v),
			Namespace: v.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &size,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image:           frontendImage,
						ImagePullPolicy: corev1.PullAlways,
						Name:            "visitors-webui",
						Ports: []corev1.ContainerPort{{
							ContainerPort: frontendPort,
							Name:          "visitors",
						}},
						Env: []corev1.EnvVar{
							{
								Name:  "REACT_APP_TITLE",
								Value: v.Spec.Title,
							},
							{
								Name:  "REACT_APP_BACKEND_HOST",
								Value: backendServiceName(v),
							},
							{
								Name:  "REACT_APP_BACKEND_PORT",
								Value: strconv.Itoa(backendPort),
							},
						},
					}},
				},
			},
		},
	}

	ctrl.SetControllerReference(v, dep, r.Scheme)
	return dep
}

func (r *VisitorsAppReconciler) frontendService(v *appsv1alpha1.VisitorsApp) *corev1.Service {
	labels := labels(v, "frontend")

	s := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      frontendServiceName(v),
			Namespace: v.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Ports: []corev1.ServicePort{{
				Protocol:   corev1.ProtocolTCP,
				Port:       frontendPort,
				TargetPort: intstr.FromInt(frontendPort),
			}},
			Type: corev1.ServiceTypeLoadBalancer,
		},
	}

	ctrl.SetControllerReference(v, s, r.Scheme)
	return s
}

func (r *VisitorsAppReconciler) updateFrontendStatus(v *appsv1alpha1.VisitorsApp) error {
	v.Status.FrontendImage = frontendImage
	err := r.Status().Update(context.TODO(), v)
	return err
}

func (r *VisitorsAppReconciler) handleFrontendChanges(ctx context.Context, v *appsv1alpha1.VisitorsApp) (*ctrl.Result, error) {
	logger := log.FromContext(ctx)

	found := &appsv1.Deployment{}
	err := r.Get(context.TODO(), types.NamespacedName{
		Name:      frontendDeploymentName(v),
		Namespace: v.Namespace,
	}, found)
	if err != nil {
		// The deployment may not have been created yet, so requeue
		return &ctrl.Result{RequeueAfter: 5 * time.Second}, err
	}

	title := v.Spec.Title
	existing := found.Spec.Template.Spec.Containers[0].Env[0].Value
	size := v.Spec.Size

	if title != existing || size != *found.Spec.Replicas {
		found.Spec.Template.Spec.Containers[0].Env[0].Value = title
		found.Spec.Replicas = &size
		logger.Info("Updating an existing frontend Deployment", "Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name, "Spec.Replicas", size, "Spec.Template.Spec.Containers[0].Env[0]", title)
		err = r.Update(context.TODO(), found)
		if err != nil {
			logger.Error(err, "Failed to update Deployment.", "Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name)
			return &ctrl.Result{}, err
		}
		// Spec updated - return and requeue
		return &ctrl.Result{Requeue: true}, nil
	}

	return nil, nil
}
