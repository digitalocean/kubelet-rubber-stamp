package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/klog"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
	"sigs.k8s.io/controller-runtime/pkg/metrics/server"

	"github.com/kontena/kubelet-rubber-stamp/pkg/apis"
	"github.com/kontena/kubelet-rubber-stamp/pkg/controller"
)

const (
	leaderElectionNamespace = "kube-system"
	leaderElectionConfigMap = "kubelet-rubber-stamp-leader-election"
)

func printVersion() {
	klog.V(2).Infof("Go Version: %s", runtime.Version())
	klog.V(2).Infof("Go OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH)
}

func main() {
	var metricsAddr string
	var leaderElect bool

	klog.InitFlags(nil)
	flag.StringVar(&metricsAddr, "metrics-addr", "", fmt.Sprintf("The address the metric endpoint binds to, or \"0\" to disable (default: %s)", server.DefaultBindAddress))
	flag.BoolVar(&leaderElect, "leader-elect", false, "Enable leader election")
	flag.Set("logtostderr", "true")
	flag.Set("v", "2")
	flag.Parse()

	printVersion()

	var namespace string
	namespace, hasErr := os.LookupEnv("WATCH_NAMESPACE")
	if hasErr {
		klog.Warning("failed to get watch namespace")
		namespace = ""
	}

	// Get a config to talk to the apiserver
	cfg, err := config.GetConfig()
	if err != nil {
		klog.Fatal(err)
	}

	switch metricsAddr {
	case "":
		klog.V(2).Infof("Exposing metrics on %s", server.DefaultBindAddress)
	case "0":
		klog.V(2).Info("Disabling metrics endpoint")
	default:
		klog.V(2).Infof("Exposing metrics on %s", metricsAddr)
	}

	// Create a new Cmd to provide shared dependencies and start components
	mgr, err := manager.New(cfg, manager.Options{
		Cache: cache.Options{
			DefaultNamespaces: map[string]cache.Config{
				namespace: {},
			},
		},
		Metrics: server.Options{
			BindAddress: metricsAddr,
		},
		LeaderElection:          leaderElect,
		LeaderElectionNamespace: leaderElectionNamespace,
		LeaderElectionID:        leaderElectionConfigMap,
	})
	if err != nil {
		klog.Fatal(err)
	}

	klog.V(2).Info("Registering Components.")

	// Setup Scheme for all resources
	if err := apis.AddToScheme(mgr.GetScheme()); err != nil {
		klog.Fatal(err)
	}

	// Setup all Controllers
	if err := controller.AddToManager(mgr); err != nil {
		klog.Fatal(err)
	}

	klog.V(2).Info("Starting the Cmd.")

	// Start the Cmd
	klog.Fatal(mgr.Start(signals.SetupSignalHandler()))
}
