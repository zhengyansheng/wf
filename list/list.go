package list

import (
	"context"
	"sort"
	"strings"

	workflowpkg "github.com/argoproj/argo-workflows/v3/pkg/apiclient/workflow"
	wfv1 "github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	argotime "github.com/argoproj/pkg/time"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

type workflowFlags struct {
	namespace      string
	status         []string
	completed      bool
	running        bool
	resubmitted    bool
	prefix         string
	createdSince   string
	finishedBefore string
	chunkSize      int64
	noHeaders      bool
	labels         string
	fields         string
}

func ListWorkflows(ctx context.Context, serviceClient workflowpkg.WorkflowServiceClient) (wfv1.Workflows, error) {
	var flags = workflowFlags{}
	listOpts := &metav1.ListOptions{
		Limit: 0,
	}
	labelSelector, err := labels.Parse(flags.labels)
	if err != nil {
		return nil, err
	}

	listOpts.LabelSelector = labelSelector.String()
	listOpts.FieldSelector = flags.fields
	var workflows wfv1.Workflows
	for {
		log.WithField("listOpts", listOpts).Debug()
		wfList, err := serviceClient.ListWorkflows(ctx, &workflowpkg.WorkflowListRequest{
			Namespace:   flags.namespace,
			ListOptions: listOpts,
			//Fields:      flags.displayFields(),
		})
		if err != nil {
			return nil, err
		}
		workflows = append(workflows, wfList.Items...)
		if wfList.Continue == "" {
			break
		}
		listOpts.Continue = wfList.Continue
	}
	workflows = workflows.
		Filter(func(wf wfv1.Workflow) bool {
			return strings.HasPrefix(wf.ObjectMeta.Name, flags.prefix)
		})
	if flags.createdSince != "" && flags.finishedBefore != "" {
		startTime, err := argotime.ParseSince(flags.createdSince)
		if err != nil {
			return nil, err
		}
		endTime, err := argotime.ParseSince(flags.finishedBefore)
		if err != nil {
			return nil, err
		}
		workflows = workflows.Filter(wfv1.WorkflowRanBetween(*startTime, *endTime))
	} else {
		if flags.createdSince != "" {
			t, err := argotime.ParseSince(flags.createdSince)
			if err != nil {
				return nil, err
			}
			workflows = workflows.Filter(wfv1.WorkflowCreatedAfter(*t))
		}
		if flags.finishedBefore != "" {
			t, err := argotime.ParseSince(flags.finishedBefore)
			if err != nil {
				return nil, err
			}
			workflows = workflows.Filter(wfv1.WorkflowFinishedBefore(*t))
		}
	}
	sort.Sort(workflows)
	return workflows, nil
}
