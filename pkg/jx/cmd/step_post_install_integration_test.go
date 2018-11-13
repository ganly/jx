// +build integration

package cmd_test

import (
	"github.com/jenkins-x/golang-jenkins"
	"github.com/jenkins-x/jx/pkg/gits"
	"github.com/jenkins-x/jx/pkg/helm"
	"github.com/jenkins-x/jx/pkg/jenkins/fake"
	"github.com/jenkins-x/jx/pkg/jx/cmd"
	"github.com/jenkins-x/jx/pkg/kube"
	"github.com/jenkins-x/jx/pkg/testkube"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"os"
	"testing"
)

func TestStepPostInstall(t *testing.T) {
	t.Parallel()

	dev := kube.CreateDefaultDevEnvironment("jx")
	testOrg := "mytestorg"
	testRepo := "mytestrepo"
	staging := kube.NewPermanentEnvironmentWithGit("staging", "https://fake.git/"+testOrg+"/"+testRepo+".git")

	o := cmd.StepPostInstallOptions{
		StepOptions: cmd.StepOptions{
			CommonOptions: cmd.CommonOptions{
				In:  os.Stdin,
				Out: os.Stdout,
				Err: os.Stderr,
			},
		},
	}
	cmd.ConfigureTestOptionsWithResources(&o.CommonOptions,
		[]runtime.Object{
			&corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      kube.ConfigMapJenkinsX,
					Namespace: "jx",
				},
				Data: map[string]string{},
			},
			testkube.CreateFakeGitSecret(),
		},
		[]runtime.Object{
			dev,
			staging,
		},
		gits.NewGitCLI(),
		helm.NewHelmCLI("helm", helm.V2, "", true),
	)

	o.BatchMode = true
	jenkinsClient := fake.NewFakeJenkins()
	o.SetJenkinsClient(jenkinsClient)
	o.GitClient = &gits.GitFake{}

	err := o.Run()
	require.NoError(t, err, "failed to run jx step post install")
	
	// assert we have a jenkins job for the staging env repo
	AssertJenkinsJobExists(t, jenkinsClient, testOrg, testRepo)

	// TODO assert we have a webhook for the staging env repo

}

// AssertJenkinsJobExists asserts that the job exists for the given organisation and repo
func AssertJenkinsJobExists(t *testing.T, jenkinsClient *fake.FakeJenkins, testOrg string, testRepo string) {
	job, err := jenkinsClient.GetJobByPath(testOrg, testRepo)
	if !assert.NoError(t, err, "failed to query Jenkins Job for %s/%s", testOrg, testRepo) {
		DumpJenkinsJobs(t, jenkinsClient)
		return
	}
	if !assert.Equal(t, job.Name, testRepo, "job.Name") {
		DumpJenkinsJobs(t, jenkinsClient)
		return
	}

	t.Logf("Found Jenkins Job at URL: %s\n", job.Url)
}

// DumpJenkinsJobs dumps the current jenkins jobs in the given client to aid debugging a failing test
func DumpJenkinsJobs(t *testing.T, jenkinsClient gojenkins.JenkinsClient) {
	jobs, err := jenkinsClient.GetJobs()
	require.NoError(t, err, "failed to get jobs")

	for _, job := range jobs {
		t.Logf("Jenkins Job: %s at %s\n", job.Name, job.Url)
		for _, cj := range job.Jobs {
			t.Logf("\t child Job: %s at %s\n", cj.Name, cj.Url)
		}
	}
}
