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

const backendPort = 8000
const defaultBackendImage = "majguo/visitors-service"

func backendDeploymentName(v *appsv1alpha1.VisitorsApp) string {
	return v.Name + "-backend"
}

func backendServiceName(v *appsv1alpha1.VisitorsApp) string {
	return v.Name + "-backend-service"
}

func getBackendImage(v *appsv1alpha1.VisitorsApp) string {
	if v.Spec.BackendImage != "" {
		return v.Spec.BackendImage
	} else {
		return defaultBackendImage
	}
}

func (r *VisitorsAppReconciler) backendDeployment(v *appsv1alpha1.VisitorsApp) *appsv1.Deployment {
	labels := labels(v, "backend")
	size := v.Spec.Size
	backendImage := getBackendImage(v)

	userSecret := &corev1.EnvVarSource{
		SecretKeyRef: &corev1.SecretKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{Name: mysqlAuthName(v)},
			Key:                  mysqlUsernameKey,
		},
	}

	passwordSecret := &corev1.EnvVarSource{
		SecretKeyRef: &corev1.SecretKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{Name: mysqlAuthName(v)},
			Key:                  mysqlPasswordKey,
		},
	}

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      backendDeploymentName(v),
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
						Image:           backendImage,
						ImagePullPolicy: corev1.PullAlways,
						Name:            "visitors-service",
						Ports: []corev1.ContainerPort{{
							ContainerPort: backendPort,
							Name:          "visitors",
						}},
						Env: []corev1.EnvVar{
							{
								Name:  "MYSQL_DATABASE",
								Value: mysqlDBName,
							},
							{
								Name:  "MYSQL_SERVICE_HOST",
								Value: mysqlServiceName(v),
							},
							{
								Name:  "MYSQL_SERVICE_PORT",
								Value: strconv.Itoa(mysqlPort),
							},
							{
								Name:      "MYSQL_USERNAME",
								ValueFrom: userSecret,
							},
							{
								Name:      "MYSQL_PASSWORD",
								ValueFrom: passwordSecret,
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

func (r *VisitorsAppReconciler) backendService(v *appsv1alpha1.VisitorsApp) *corev1.Service {
	labels := labels(v, "backend")

	s := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      backendServiceName(v),
			Namespace: v.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Ports: []corev1.ServicePort{{
				Protocol:   corev1.ProtocolTCP,
				Port:       backendPort,
				TargetPort: intstr.FromInt(backendPort),
			}},
			Type: corev1.ServiceTypeClusterIP,
		},
	}

	ctrl.SetControllerReference(v, s, r.Scheme)
	return s
}

func (r *VisitorsAppReconciler) updateBackendStatus(v *appsv1alpha1.VisitorsApp) error {
	v.Status.BackendImage = getBackendImage(v)
	err := r.Status().Update(context.TODO(), v)
	return err
}

func (r *VisitorsAppReconciler) handleBackendChanges(ctx context.Context, v *appsv1alpha1.VisitorsApp) (*ctrl.Result, error) {
	logger := log.FromContext(ctx)

	found := &appsv1.Deployment{}
	err := r.Get(context.TODO(), types.NamespacedName{
		Name:      backendDeploymentName(v),
		Namespace: v.Namespace,
	}, found)
	if err != nil {
		// The deployment may not have been created yet, so requeue
		return &ctrl.Result{RequeueAfter: 5 * time.Second}, err
	}

	size := v.Spec.Size
	existingBackendImage := found.Spec.Template.Spec.Containers[0].Image
	specBackendImage := getBackendImage(v)

	if size != *found.Spec.Replicas || specBackendImage != existingBackendImage {
		found.Spec.Replicas = &size
		found.Spec.Template.Spec.Containers[0].Image = specBackendImage

		logger.Info("Updating an existing backend Deployment", "Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name,
			"Spec.Replicas", size, "Spec.Template.Spec.Containers[0].Image", specBackendImage)
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
