package controllers

import (
	"context"

	appsv1alpha1 "github.com/majguo/visitors-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func mysqlDeploymentName(v *appsv1alpha1.VisitorsApp) string {
	return v.Name + "-mysql"
}

func mysqlServiceName(v *appsv1alpha1.VisitorsApp) string {
	return v.Name + "-mysql-service"
}

func mysqlAuthName(v *appsv1alpha1.VisitorsApp) string {
	return v.Name + "-mysql-auth"
}

func (r *VisitorsAppReconciler) mysqlAuthSecret(v *appsv1alpha1.VisitorsApp) *corev1.Secret {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      mysqlAuthName(v),
			Namespace: v.Namespace,
		},
		Type: "Opaque",
		StringData: map[string]string{
			"username": "visitors-user",
			"password": "visitors-pass",
		},
	}
	ctrl.SetControllerReference(v, secret, r.Scheme)
	return secret
}

func (r *VisitorsAppReconciler) mysqlDeployment(v *appsv1alpha1.VisitorsApp) *appsv1.Deployment {
	labels := labels(v, "mysql")
	size := int32(1)

	userSecret := &corev1.EnvVarSource{
		SecretKeyRef: &corev1.SecretKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{Name: mysqlAuthName(v)},
			Key:                  "username",
		},
	}

	passwordSecret := &corev1.EnvVarSource{
		SecretKeyRef: &corev1.SecretKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{Name: mysqlAuthName(v)},
			Key:                  "password",
		},
	}

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      mysqlDeploymentName(v),
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
						Image: "mysql:5.7",
						Name:  "visitors-mysql",
						Ports: []corev1.ContainerPort{{
							ContainerPort: 3306,
							Name:          "mysql",
						}},
						Env: []corev1.EnvVar{
							{
								Name:  "MYSQL_ROOT_PASSWORD",
								Value: "password",
							},
							{
								Name:  "MYSQL_DATABASE",
								Value: "visitors_db",
							},
							{
								Name:      "MYSQL_USER",
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

func (r *VisitorsAppReconciler) mysqlService(v *appsv1alpha1.VisitorsApp) *corev1.Service {
	labels := labels(v, "mysql")

	s := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      mysqlServiceName(v),
			Namespace: v.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Ports: []corev1.ServicePort{{
				Port: 3306,
			}},
			Type: "ClusterIP",
		},
	}

	ctrl.SetControllerReference(v, s, r.Scheme)
	return s
}

// Returns whether or not the MySQL deployment is running
func (r *VisitorsAppReconciler) isMysqlUp(ctx context.Context, v *appsv1alpha1.VisitorsApp) bool {
	logger := log.FromContext(ctx)

	deployment := &appsv1.Deployment{}

	err := r.Get(context.TODO(), types.NamespacedName{
		Name:      mysqlDeploymentName(v),
		Namespace: v.Namespace,
	}, deployment)

	if err != nil {
		logger.Error(err, "Deployment mysql not found")
		return false
	}

	if deployment.Status.ReadyReplicas == 1 {
		return true
	}

	return false
}
