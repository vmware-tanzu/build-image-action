package cmd_test

import (
	"context"
	"encoding/json"
	"github.com/google/go-cmp/cmp"
	"github.com/pivotal/kpack/pkg/apis/build/v1alpha2"
	"github.com/pivotal/kpack/pkg/apis/core/v1alpha1"
	"github.com/vmware-tanzu/build-image-action/cmd"
	"github.com/vmware-tanzu/build-image-action/cmd/cmdfakes"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"testing"
)

type TestCase struct {
	Name           string
	Config         *cmd.Config
	ExpectedCreate *unstructured.Unstructured
	ExpectedGets   []*unstructured.Unstructured
}

func (tc *TestCase) Run(t *testing.T) {
	t.Run(tc.Name, func(t *testing.T) {
		fakeClient := &cmdfakes.FakeClient{}

		existingBuilder := &v1alpha2.ClusterBuilder{
			Status: v1alpha2.BuilderStatus{
				Stack: v1alpha1.BuildStack{
					RunImage: "index.docker.io/paketobuildpacks/run@sha256:00aa8f",
				},
				LatestImage: "gcr.io/my-project/my-builder/clusterbuilder@sha256:95ab59",
			},
		}

		count := 0

		fakeClient.GetStub = func(ctx context.Context, key types.NamespacedName, obj client.Object, _ ...client.GetOption) error {

			if count > len(tc.ExpectedGets) {
				t.Fatalf("VERY BAD")
				return nil
			}

			currentGet := tc.ExpectedGets[count]

			if key.Name == "my-builder" {
				bytes, _ := json.Marshal(existingBuilder)
				_ = json.Unmarshal(bytes, obj)
			} else {
				t.Fatalf("unimplemented call for get")
			}

			return nil
		}

		err := tc.Config.Build(context.Background(), fakeClient)
		if err != nil {
			t.Fatalf("err: %+v", err)
		}

		if tc.ExpectedCreate != nil {
			if fakeClient.CreateCallCount() != 1 {
				t.Fatalf("got number of create calls: %d; expected 1", fakeClient.CreateCallCount())
			}

			_, obj, _ := fakeClient.CreateArgsForCall(0)
			if !equality.Semantic.DeepEqual(obj, tc.ExpectedCreate) {
				t.Fatal(cmp.Diff(tc.ExpectedCreate, obj))
			}
		}

		if fakeClient.GetCallCount() != len(tc.ExpectedGets) {
			t.Fatalf("got number of get calls: %d; expected %d", fakeClient.GetCallCount(), len(tc.ExpectedGets))
		}

		//pod name

		// pod status

		for _, expectedGet := range tc.ExpectedGets {
			_, namespacedName, _, _ := fakeClient.GetArgsForCall(0)
			if namespacedName.Name != expectedGet.GetName() {
				t.Fatalf("got builder: %s; expected %s", namespacedName.Name, expectedGet.GetName())
			}

		}

	})
}

func TestNewCmdKpack(t *testing.T) {
	tests := map[string]TestCase{
		"create kpack build": {
			Config: &cmd.Config{
				Namespace:          "my-ns",
				CaCert:             "abc123",
				Server:             "https://my-server.com",
				Token:              "efg456",
				GitServer:          "https://github.com",
				GitRepo:            "my-org/my-repo",
				GitSha:             "xyz890",
				Tag:                "gcr.io/my-project/my-app",
				Env:                "BP_JAVA_VERSION=17",
				ServiceAccountName: "my-sa",
				GithubOutput:       "output.txt",
			},
			ExpectedCreate: &unstructured.Unstructured{Object: map[string]interface{}{
				"apiVersion": "kpack.io/v1alpha2",
				"kind":       "Build",
				"metadata": map[string]interface{}{
					"generateName": "my-org-my-repo-",
					"namespace":    "my-ns",
					"annotations": map[string]interface{}{
						"app.kubernetes.io/managed-by": "vmware-tanzu/build-image-action ",
					},
				},
				"spec": map[string]interface{}{
					"builder": map[string]interface{}{
						"image": "gcr.io/my-project/my-builder/clusterbuilder@sha256:95ab59",
					},
					"runImage": map[string]interface{}{
						"image": "index.docker.io/paketobuildpacks/run@sha256:00aa8f",
					},
					"serviceAccountName": "my-sa",
					"source": map[string]interface{}{
						"git": map[string]interface{}{
							"url":      "https://github.com/my-org/my-repo",
							"revision": "xyz890",
						},
					},
					"tags": []string{
						"gcr.io/my-project/my-app",
					},
					"env": []map[string]string{{"name": "BP_JAVA_VERSION", "value": "17"}},
				},
			},
			},
			ExpectedGets: []*unstructured.Unstructured{
				{
					Object: map[string]interface{}{
						"apiVersion": "kpack.io/v1alpha2",
						"kind":       "ClusterBuilder",
						"metadata": map[string]interface{}{
							"name": "my-builder",
						},
						"spec": map[string]interface{}{
							"tag": "gcr.io/my-project/my-builder",
						},
						"status": map[string]interface{}{
							"latestImage": "gcr.io/my-project/my-builder/clusterbuilder@sha256:95ab59",
							"stack": map[string]interface{}{
								"runImage": "index.docker.io/paketobuildpacks/run@sha256:00aa8f",
							},
						},
					},
				},
				//{
				//	Object: map[string]interface{}{
				//		"apiVersion": "kpack.io/v1alpha2",
				//		"kind":       "Build",
				//		"metadata": map[string]interface{}{
				//			"generateName": "my-org-my-repo-",
				//			"namespace":    "my-ns",
				//			"annotations": map[string]interface{}{
				//				"app.kubernetes.io/managed-by": "vmware-tanzu/build-image-action ",
				//			},
				//		},
				//		"spec": map[string]interface{}{
				//			"serviceAccountName": "my-sa",
				//			"source": map[string]interface{}{
				//				"git": map[string]interface{}{
				//					"url":      "https://github.com/my-org/my-repo",
				//					"revision": "xyz890",
				//				},
				//			},
				//			"tags": []string{
				//				"gcr.io/my-project/my-app",
				//			},
				//			"env": []map[string]string{{"name": "BP_JAVA_VERSION", "value": "17"}},
				//		},
				//	},
				//},
			},
		},
	}

	for name, tt := range tests {
		tt.Name = name
		tt.Run(t)
	}
}

//
//
//Happiest Path:
//	Expect Create -> happy create
//    Expected Gets:
//		- builder
// 		- logs
//  		- build - with pod
//        - build - complete happy
//
//	get builder
//	create build
//    get build - returns pod
//	watch pod
//    get build - returns happy
