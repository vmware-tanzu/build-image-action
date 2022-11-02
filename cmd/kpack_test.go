package cmd_test

import (
	"context"
	"github.com/google/go-cmp/cmp"
	"github.com/vmware-tanzu/build-image-action/cmd"
	"github.com/vmware-tanzu/build-image-action/cmd/cmdfakes"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"testing"
)

type TestCase struct {
	Name           string
	Config         *cmd.Config
	ExpectedCreate *unstructured.Unstructured
}

func (tc *TestCase) Run(t *testing.T) {
	t.Run(tc.Name, func(t *testing.T) {
		fakeClient := &cmdfakes.FakeClient{}

		err := tc.Config.Build(context.Background(), fakeClient)
		if err != nil {
			t.Fatalf("err: %+v", err)
		}

		if fakeClient.CreateCallCount() != 1 {
			t.Fatalf("got number of create calls: %d; expected 1", fakeClient.CreateCallCount())
		}

		_, obj, _ := fakeClient.CreateArgsForCall(0)
		if !equality.Semantic.DeepEqual(obj, tc.ExpectedCreate) {
			t.Fatal(cmp.Diff(tc.ExpectedCreate, obj))
		}
	})
}

func TestNewCmdKpack(t *testing.T) {
	tests := map[string]TestCase{
		"create kpack build": {
			Config: &cmd.Config{
				Namespace: "my-ns",
			},
			ExpectedCreate: &unstructured.Unstructured{Object: map[string]interface{}{
				"apiVersion": "kpack.io/v1alpha2",
				"kind":       "Build",
				"metadata": map[string]interface{}{
					"namespace": "my-ns",
				},
			}},
		},
	}

	for name, tt := range tests {
		tt.Name = name
		tt.Run(t)
	}
}
